import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import axios from 'axios';
import {getToken} from './defaultAxios';
import {snack} from './GlobalAction';
import * as UserAction from './UserAction';

/** Fetches all messages from the current user. */
export function fetchMessages() {
    axios.get(config.get('url') + 'message').then((resp) => {
        dispatcher.dispatch({type: 'UPDATE_MESSAGES', payload: resp.data});
    });
}

/** Deletes all messages from the current user. */
export function deleteMessages() {
    axios.delete(config.get('url') + 'message').then(fetchMessages).then(() => snack('Messages deleted'));
}

/**
 * Deletes all messages from the current user and an application.
 * @param {int} id the application id
 */
export function deleteMessagesByApp(id) {
    axios.delete(config.get('url') + 'application/' + id + '/message').then(fetchMessages)
        .then(() => snack('Deleted all messages from the application'));
}

/**
 * Deletes a message by id.
 * @param {int} id the message id
 */
export function deleteMessage(id) {
    axios.delete(config.get('url') + 'message/' + id).then(fetchMessages).then(() => snack('Message deleted'));
}

let wsActive = false;

/**
 * Starts listening to the stream for new messages.
 */
export function listenToWebSocket() {
    if (!getToken() || wsActive) {
        return;
    }
    wsActive = true;

    const wsUrl = config.get('url').replace('http', 'ws').replace('https', 'wss');
    const ws = new WebSocket(wsUrl + 'stream?token=' + getToken());

    ws.onerror = (e) => {
        wsActive = false;
        console.log('WebSocket connection errored', e);
    };

    ws.onmessage = (data) => dispatcher.dispatch({type: 'ONE_MESSAGE', payload: JSON.parse(data.data)});

    ws.onclose = () => {
        wsActive = false;
        UserAction.tryAuthenticate().then(() => {
            snack('WebSocket connection closed, trying again in 30 seconds.');
            setTimeout(listenToWebSocket, 30000);
        });
    };
}
