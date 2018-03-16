import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import {getToken, setAuthorizationToken} from './defaultAxios';
import * as GlobalAction from './GlobalAction';
import axios from 'axios';
import {detect} from 'detect-browser';
import ClientStore from '../stores/ClientStore';

/**
 * Login the user.
 * @param {string} username
 * @param {string} password
 */
export function login(username, password) {
    const browser = detect();
    const name = (browser && browser.name + ' ' + browser.version) || 'unknown browser';
    axios.request(config.get('url') + 'client', {
        method: 'POST',
        data: {name: name},
        auth: {username: username, password: password},
    }).then(function(resp) {
        setAuthorizationToken(resp.data.token);
        GlobalAction.initialLoad();
    }).catch(() => {
        dispatcher.dispatch({type: 'LOGIN_FAILED'});
    });
}

/** Log the user out. */
export function logout() {
    if (getToken() !== null) {
        axios.delete(config.get('url') + 'client/' + ClientStore.getIdByToken(getToken())).then(() => {
            setAuthorizationToken(null);
            dispatcher.dispatch({type: 'REMOVE_CURRENT_USER'});
        });
    }
}

/** Fetches the current user. */
export function fetchCurrentUser() {
    axios.get(config.get('url') + 'current/user').then(function(resp) {
        dispatcher.dispatch({type: 'SET_CURRENT_USER', payload: resp.data});
    });
}

/**
 * Changes the current user.
 * @param {string} pass
 */
export function changeCurrentUser(pass) {
    axios.post(config.get('url') + 'current/user/password', {pass});
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
    axios.delete(config.get('url') + 'user/' + id).then(function() {
        fetchUsers();
    });
}

/**
 * Create a user.
 * @param {string} name
 * @param {string} pass
 * @param {bool} admin if true, the user is an administrator
 */
export function createUser(name, pass, admin) {
    axios.post(config.get('url') + 'user', {name, pass, admin}).then(function() {
        fetchUsers();
    });
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
        fetchCurrentUser(); // just in case update current user
    });
}
