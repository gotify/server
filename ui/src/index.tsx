import * as React from 'react';
import * as ReactDOM from 'react-dom';
import 'typeface-roboto';
import {initAxios} from './apiAuth';
import * as config from './config';
import Layout from './layout/Layout';
import registerServiceWorker from './registerServiceWorker';
import * as Notifications from './snack/browserNotification';
import {CurrentUser} from './CurrentUser';
import {AppStore} from './application/AppStore';
import {reaction} from 'mobx';
import {WebSocketStore} from './message/WebSocketStore';
import {SnackManager} from './snack/SnackManager';
import {InjectProvider, StoreMapping} from './inject';
import {UserStore} from './user/UserStore';
import {MessagesStore} from './message/MessagesStore';
import {ClientStore} from './client/ClientStore';
import {PluginStore} from './plugin/PluginStore';

const defaultDevConfig = {
    url: 'http://localhost:80/',
};

const {port, hostname, protocol} = window.location;
const slashes = protocol.concat('//');
const url = slashes.concat(port ? hostname.concat(':', port) : hostname);
const urlWithSlash = url.endsWith('/') ? url : url.concat('/');

const defaultProdConfig = {
    url: urlWithSlash,
};

declare global {
    // tslint:disable-next-line
    interface Window {
        config: config.IConfig;
    }
}

const initStores = (): StoreMapping => {
    const snackManager = new SnackManager();
    const appStore = new AppStore(snackManager.snack);
    const userStore = new UserStore(snackManager.snack);
    const messagesStore = new MessagesStore(appStore, snackManager.snack);
    const currentUser = new CurrentUser(snackManager.snack);
    const clientStore = new ClientStore(snackManager.snack);
    const wsStore = new WebSocketStore(snackManager.snack, currentUser);
    const pluginStore = new PluginStore(snackManager.snack);
    appStore.onDelete = () => messagesStore.clearAll();

    return {
        appStore,
        snackManager,
        userStore,
        messagesStore,
        currentUser,
        clientStore,
        wsStore,
        pluginStore,
    };
};

(function clientJS() {
    Notifications.requestPermission();
    if (process.env.NODE_ENV === 'production') {
        config.set(window.config || defaultProdConfig);
    } else {
        config.set(window.config || defaultDevConfig);
    }
    const stores = initStores();
    initAxios(stores.currentUser, stores.snackManager.snack);

    reaction(
        () => stores.currentUser.loggedIn,
        (loggedIn) => {
            if (loggedIn) {
                stores.wsStore.listen((message) => {
                    stores.messagesStore.publishSingleMessage(message);
                    Notifications.notifyNewMessage(message);
                });
                stores.appStore.refresh();
            } else {
                stores.messagesStore.clearAll();
                stores.appStore.clear();
                stores.clientStore.clear();
                stores.userStore.clear();
                stores.wsStore.close();
            }
        }
    );

    stores.currentUser.tryAuthenticate().catch(() => {});

    ReactDOM.render(
        <InjectProvider stores={stores}>
            <Layout />
        </InjectProvider>,
        document.getElementById('root')
    );
    registerServiceWorker();
})();
