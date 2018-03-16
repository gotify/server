import axios from 'axios';
import dispatcher from '../stores/dispatcher';
import * as GlobalAction from './GlobalAction';

let currentToken = null;
const tokenKey = 'gotify-login-key';

/**
 * Set the authorization token for the next requests.
 * @param {string} token the gotify application token
 */
export function setAuthorizationToken(token) {
    currentToken = token;
    if (token) {
        localStorage.setItem(tokenKey, token);
        axios.defaults.headers.common['X-Gotify-Key'] = token;
    } else {
        localStorage.removeItem(tokenKey);
        delete axios.defaults.headers.common['X-Gotify-Key'];
    }
}

axios.interceptors.response.use(null, (error) => {
    if (error.response.status === 401) {
        setAuthorizationToken(null);
        dispatcher.dispatch({type: 'REMOVE_CURRENT_USER'});
    }

    return Promise.reject(error);
});

/**
 * @return {string} the application token
 */
export function getToken() {
    return currentToken;
}

/** Checks if the current user is logged, if so update the state. */
export function checkIfAlreadyLoggedIn() {
    const key = localStorage.getItem(tokenKey);
    if (key) {
        setAuthorizationToken(key);
        GlobalAction.initialLoad();
    }
}
