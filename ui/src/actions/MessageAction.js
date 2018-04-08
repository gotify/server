import dispatcher from '../stores/dispatcher';
import config from 'react-global-configuration';
import axios from 'axios';
import {getToken} from './defaultAxios';
import {snack} from './GlobalAction';
import * as UserAction from './UserAction';

export function fetchMessagesApp(id, since) {
    if (id === -1) {
        return axios.get(config.get('url') + 'message?since=' + since).then((resp) => {
            newMessages(-1, resp.data);
        });
    } else {
        return axios.get(config.get('url') + 'application/' + id + '/message?since=' + since).then((resp) => {
            newMessages(id, resp.data);
        });
    }
}

function newMessages(id, data) {
    dispatcher.dispatch({
        type: 'UPDATE_MESSAGES', payload: {
            messages: data.messages,
            hasMore: 'next' in data.paging,
            nextSince: data.paging.since,
            id,
        },
    });
}

/**
 * Deletes all messages from the current user and an application.
 * @param {int} id the application id
 */
export function deleteMessagesByApp(id) {
    if (id === -1) {
        axios.delete(config.get('url') + 'message').then(() => {
            dispatcher.dispatch({type: 'DELETE_MESSAGES', id: -1});
            snack('Messages deleted');
        });
    } else {
        axios.delete(config.get('url') + 'application/' + id + '/message')
            .then(() => {
                dispatcher.dispatch({type: 'DELETE_MESSAGES', id});
                snack('Deleted all messages from the application');
            });
    }
}

export function deleteMessage(msg) {
    axios.delete(config.get('url') + 'message/' + msg.id).then(() => {
        dispatcher.dispatch({type: 'DELETE_MESSAGE', payload: msg});
        snack('Message deleted');
    });
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
