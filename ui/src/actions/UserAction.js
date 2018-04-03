import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import {getToken, setAuthorizationToken} from './defaultAxios';
import * as GlobalAction from './GlobalAction';
import axios from 'axios';
import {detect} from 'detect-browser';
import ClientStore from '../stores/ClientStore';
import {snack} from './GlobalAction';

/**
 * Login the user.
 * @param {string} username
 * @param {string} password
 */
export function login(username, password) {
    const browser = detect();
    const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
    authenticating();
    axios.create().request(config.get('url') + 'client', {
        method: 'POST',
        data: {name: name},
        auth: {username: username, password: password},
    }).then(function(resp) {
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
    if (getToken() !== null) {
        axios.delete(config.get('url') + 'client/' + ClientStore.getIdByToken(getToken())).then(() => {
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
export function changeCurrentUser(pass) {
    axios.post(config.get('url') + 'current/user/password', {pass}).then(() => snack('Password changed'));
}

/** Fetches all users. */
export function fetchUsers() {
    axios.get(config.get('url') + 'user').then(function(resp) {
        dispatcher.dispatch({type: 'UPDATE_USERS', payload: resp.data});
    });
}

/**
 * Delete a user.
 * @param {int} id the user id
 */
export function deleteUser(id) {
    axios.delete(config.get('url') + 'user/' + id).then(fetchUsers).then(() => snack('User deleted'));
}

/**
 * Create a user.
 * @param {string} name
 * @param {string} pass
 * @param {bool} admin if true, the user is an administrator
 */
export function createUser(name, pass, admin) {
    axios.post(config.get('url') + 'user', {name, pass, admin}).then(fetchUsers).then(() => snack('User created'));
}

/**
 * Update a user by id.
 * @param {int} id
 * @param {string} name
 * @param {string} pass empty if no change
 * @param {bool} admin if true, the user is an administrator
 */
export function updateUser(id, name, pass, admin) {
    axios.post(config.get('url') + 'user/' + id, {name, pass, admin}).then(function() {
        fetchUsers();
        tryAuthenticate(); // try authenticate updates the current user
        snack('User updated');
    });
}
