import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import axios from 'axios';
import {getToken} from './defaultAxios';

/** Fetches all messages from the current user. */
export function fetchMessages() {
    axios.get(config.get('url') + 'message').then((resp) => {
        dispatcher.dispatch({type: 'UPDATE_MESSAGES', payload: resp.data});
    });
}

/** Deletes all messages from the current user. */
export function deleteMessages() {
    axios.delete(config.get('url') + 'message').then(fetchMessages);
}

/**
 * Deletes all messages from the current user and an application.
 * @param {int} id the application id
 */
export function deleteMessagesByApp(id) {
    axios.delete(config.get('url') + 'application/' + id + '/message').then(fetchMessages);
}

/**
 * Deletes a message by id.
 * @param {int} id the message id
 */
export function deleteMessage(id) {
    axios.delete(config.get('url') + 'message/' + id).then(fetchMessages);
}

/**
 * Starts listening to the stream for new messages.
 */
export function listenToWebSocket() {
    if (!getToken()) {
        return;
    }

    const ws = new WebSocket('ws://localhost:80/stream?token=' + getToken());

    ws.onerror = (e) => {
        console.log('WebSocket connection errored; trying again in 60 seconds', e);
        setTimeout(listenToWebSocket, 60000);
    };

    ws.onmessage = (data) => {
        dispatcher.dispatch({type: 'ONE_MESSAGE', payload: JSON.parse(data.data)});
    };
}
