import axios from 'axios';
import {snack} from './GlobalAction';
import {tryAuthenticate} from './UserAction';

const tokenKey = 'gotify-login-key';

/**
 * Set the authorization token for the next requests.
 * @param {string|null} token the gotify application token
 */
export function setAuthorizationToken(token: string | null) {
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

    const status = error.response.status;

    if (status === 401) {
        tryAuthenticate().then(() => snack('Could not complete request.'));
    }

    if (status === 400) {
        snack(error.response.data.error + ': ' + error.response.data.errorDescription);
    }

    return Promise.reject(error);
});

/**
 * @return {string} the application token
 */
export function getToken(): string | null {
    return localStorage.getItem(tokenKey);
}
