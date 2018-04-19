import axios, {AxiosResponse} from 'axios';
import {detect} from 'detect-browser';
import * as config from '../config';
import ClientStore from '../stores/ClientStore';
import dispatcher from '../stores/dispatcher';
import {getToken, setAuthorizationToken} from './defaultAxios';
import * as GlobalAction from './GlobalAction';
import {snack} from './GlobalAction';

/**
 * Login the user.
 * @param {string} username
 * @param {string} password
 */
export function login(username: string, password: string) {
    const browser = detect();
    const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
    authenticating();
    axios.create().request({
        url: config.get('url') + 'client',
        method: 'POST',
        data: {name},
        auth: {username, password},
    }).then((resp) => {
        snack(`A client named '${name}' was created for your session.`);
        setAuthorizationToken(resp.data.token);
        tryAuthenticate().then(GlobalAction.initialLoad)
            .catch(() => console.log('create client succeeded, but authenticated with given token failed'));
    }).catch(() => {
        snack('Login failed');
        noAuthentication();
    });
}

/** Log the user out. */
export function logout() {
    const token = getToken();
    if (token !== null) {
        axios.delete(config.get('url') + 'client/' + ClientStore.getIdByToken(token)).then(() => {
            setAuthorizationToken(null);
            noAuthentication();
        });
    }
}

export function tryAuthenticate() {
    return axios.create().get(config.get('url') + 'current/user', {headers: {'X-Gotify-Key': getToken()}}).then((resp) => {
        dispatcher.dispatch({type: 'AUTHENTICATED', payload: resp.data});
        return resp;
    }).catch((resp) => {
        if (getToken()) {
            setAuthorizationToken(null);
            snack('Authentication failed, try to re-login. (client or user was deleted)');
        }
        noAuthentication();
        return Promise.reject(resp);
    });
}

export function checkIfAlreadyLoggedIn() {
    const token = getToken();
    if (token) {
        setAuthorizationToken(token);
        tryAuthenticate().then(GlobalAction.initialLoad);
    } else {
        noAuthentication();
    }
}

function noAuthentication() {
    dispatcher.dispatch({type: 'NO_AUTHENTICATION'});
}

function authenticating() {
    dispatcher.dispatch({type: 'AUTHENTICATING'});
}

/**
 * Changes the current user.
 * @param {string} pass
 */
export function changeCurrentUser(pass: string) {
    axios.post(config.get('url') + 'current/user/password', {pass}).then(() => snack('Password changed'));
}

/** Fetches all users. */
export function fetchUsers() {
    axios.get(config.get('url') + 'user').then((resp: AxiosResponse<IUser[]>) => {
        dispatcher.dispatch({type: 'UPDATE_USERS', payload: resp.data});
    });
}

/**
 * Delete a user.
 * @param {int} id the user id
 */
export function deleteUser(id: number) {
    axios.delete(config.get('url') + 'user/' + id).then(fetchUsers).then(() => snack('User deleted'));
}

/**
 * Create a user.
 * @param {string} name
 * @param {string} pass
 * @param {bool} admin if true, the user is an administrator
 */
export function createUser(name: string, pass: string, admin: boolean) {
    axios.post(config.get('url') + 'user', {name, pass, admin}).then(fetchUsers).then(() => snack('User created'));
}

/**
 * Update a user by id.
 * @param {int} id
 * @param {string} name
 * @param {string} pass empty if no change
 * @param {bool} admin if true, the user is an administrator
 */
export function updateUser(id: number, name: string, pass: string | null, admin: boolean) {
    axios.post(config.get('url') + 'user/' + id, {name, pass, admin}).then(() => {
        fetchUsers();
        tryAuthenticate(); // try authenticate updates the current user
        snack('User updated');
    });
}
