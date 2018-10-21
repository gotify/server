import * as React from 'react';
import * as ReactDOM from 'react-dom';
import 'typeface-roboto';
import 'typeface-roboto-mono';
import './actions/defaultAxios';
import * as config from './config';
import Layout from './Layout';
import registerServiceWorker from './registerServiceWorker';
import * as Notifications from './stores/Notifications';
import {currentUser} from './stores/CurrentUser';
import AppStore from './stores/AppStore';
import {reaction} from 'mobx';
import {WebSocketStore} from './stores/WebSocketStore';
import SnackManager from './stores/SnackManager';

const defaultDevConfig = {
    url: 'http://localhost:80/',
};

const {port, hostname, protocol} = window.location;
const slashes = protocol.concat('//');
const url = slashes.concat(hostname.concat(':', port));
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

(function clientJS() {
    Notifications.requestPermission();
    if (process.env.NODE_ENV === 'production') {
        config.set(window.config || defaultProdConfig);
    } else {
        config.set(window.config || defaultDevConfig);
    }
    const ws = new WebSocketStore(SnackManager.snack);
    reaction(
        () => currentUser.loggedIn,
        (loggedIn) => {
            if (loggedIn) {
                ws.listen();
            } else {
                ws.close();
            }
            AppStore.refresh();
        }
    );

    currentUser.tryAuthenticate();
    ReactDOM.render(<Layout />, document.getElementById('root'));
    registerServiceWorker();
})();
