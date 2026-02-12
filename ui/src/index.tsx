import * as React from 'react';
import {createRoot} from 'react-dom/client';
import 'typeface-roboto';
import * as config from './config';
import Layout from './layout/Layout';
import {unregister} from './registerServiceWorker';
import {CurrentUser} from './CurrentUser';
import {AppStore} from './application/AppStore';
import {WebSocketStore} from './message/WebSocketStore';
import {SnackManager} from './snack/SnackManager';
import {UserStore} from './user/UserStore';
import {MessagesStore} from './message/MessagesStore';
import {ClientStore} from './client/ClientStore';
import {PluginStore} from './plugin/PluginStore';
import {registerReactions} from './reactions';
import {StoreContext, StoreMapping} from './stores';

const {port, hostname, protocol, pathname} = window.location;
const slashes = protocol.concat('//');
const path = pathname.endsWith('/') ? pathname : pathname.substring(0, pathname.lastIndexOf('/'));
const url = slashes.concat(port ? hostname.concat(':', port) : hostname) + path;
const urlWithSlash = url.endsWith('/') ? url : url.concat('/');

const prodUrl = urlWithSlash;

const initStores = (): StoreMapping => {
    const snackManager = new SnackManager();
    const currentUser = new CurrentUser(snackManager.snack);
    const appStore = new AppStore(currentUser, snackManager.snack);
    const userStore = new UserStore(currentUser, snackManager.snack);
    const messagesStore = new MessagesStore(currentUser, appStore, snackManager.snack);
    const clientStore = new ClientStore(currentUser, snackManager.snack);
    const wsStore = new WebSocketStore(snackManager.snack, currentUser);
    const pluginStore = new PluginStore(currentUser, snackManager.snack);
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
    config.set('url', prodUrl);
    const stores = initStores();
    registerReactions(stores);

    stores.currentUser.tryAuthenticate().catch(() => {});

    window.onbeforeunload = () => {
        stores.wsStore.close();
    };

    createRoot(document.getElementById('root')!).render(
        <StoreContext.Provider value={stores}>
            <Layout />
        </StoreContext.Provider>
    );
    unregister();
})();
