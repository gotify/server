import axios from 'axios';
import {snack} from './GlobalAction';
import {tryAuthenticate} from './UserAction';

const tokenKey = 'gotify-login-key';

/**
 * Set the authorization token for the next requests.
 * @param {string|null} token the gotify application token
 */
export function setAuthorizationToken(token) {
    if (token) {
        localStorage.setItem(tokenKey, token);
        axios.defaults.headers.common['X-Gotify-Key'] = token;
    } else {
        localStorage.removeItem(tokenKey);
        delete axios.defaults.headers.common['X-Gotify-Key'];
    }
}

axios.interceptors.response.use(undefined, (error) => {
    if (!error.response) {
        snack('Gotify server is not reachable, try refreshing the page.');
        return Promise.reject(error);
    }

    if (error.response.status === 401) {
        tryAuthenticate().then(() => snack('Could not complete request.'));
    }

    return Promise.reject(error);
});

/**
 * @return {string} the application token
 */
export function getToken() {
    return localStorage.getItem(tokenKey);
}
