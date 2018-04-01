import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import axios from 'axios';
import {getToken} from './defaultAxios';
import {snack} from './GlobalAction';

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

/**
 * Starts listening to the stream for new messages.
 */
export function listenToWebSocket() {
    if (!getToken()) {
        return;
    }
    const wsUrl = config.get('url').replace('http', 'ws').replace('https', 'wss');
    const ws = new WebSocket(wsUrl + 'stream?token=' + getToken());

    ws.onerror = (e) => {
        console.log('WebSocket connection errored; trying again in 60 seconds', e);
        snack('Could not connect to the web socket, trying again in 60 seconds.');
        setTimeout(listenToWebSocket, 60000);
    };

    ws.onmessage = (data) => dispatcher.dispatch({type: 'ONE_MESSAGE', payload: JSON.parse(data.data)});

    ws.onclose = (data) => console.log('WebSocket closed, this normally means the client was deleted.', data);
}
