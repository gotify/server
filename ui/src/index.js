import React from 'react';
import ReactDOM from 'react-dom';
import Layout from './Layout';
import registerServiceWorker from './registerServiceWorker';
import {checkIfAlreadyLoggedIn} from './actions/defaultAxios';
import config from 'react-global-configuration';
import 'typeface-roboto';
import 'typeface-roboto-mono';

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

(function clientJS() {
    if (process.env.NODE_ENV === 'production') {
        config.set(window.config || defaultProdConfig);
    } else {
        config.set(window.config || defaultDevConfig);
    }
    checkIfAlreadyLoggedIn();
    ReactDOM.render(<Layout/>, document.getElementById('root'));
    registerServiceWorker();
}());
