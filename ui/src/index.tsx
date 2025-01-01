import './init';

import React from 'react';
import ReactDOM from 'react-dom/client';

import {initAxios} from './apiAuth';
import App from './App.tsx';
import * as config from './config';
import {unregister} from './registerServiceWorker';
import {tryAuthenticate} from './store/auth-actions.ts';
import {loadStoredTheme} from './store/ui-actions.ts';

import {Provider, } from 'react-redux';
import store from './store/index';

// the development server of vite will proxy this to the backend
const devUrl = '/api/';

const {port, hostname, protocol, pathname} = window.location;
const slashes = protocol.concat('//');
const path = pathname.endsWith('/') ? pathname : pathname.substring(0, pathname.lastIndexOf('/'));
const url = slashes.concat(port ? hostname.concat(':', port) : hostname) + path;
const urlWithSlash = url.endsWith('/') ? url : url.concat('/');

const prodUrl = urlWithSlash;

const clientJS = async () => {

    if (import.meta.env.MODE === 'production') {
        config.set('url', prodUrl);
    } else {
        config.set('url', devUrl);
        config.set('register', true);
    }

    await store.dispatch(loadStoredTheme());
    try {
        await store.dispatch(tryAuthenticate());
    } catch (e) {
        // console.info('Automatic login failed, will forward later to login page.')
    }

    initAxios();

    const root = ReactDOM.createRoot(document.getElementById('root')!);
    root.render(
        // TODO: enable strict mode again
        // <React.StrictMode>
            <Provider store={store}>
                <App />
            </Provider>
        // </React.StrictMode>
    );
    unregister();
};

clientJS().then();
