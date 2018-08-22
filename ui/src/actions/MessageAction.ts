import axios, {AxiosResponse} from 'axios';
import * as config from '../config';
import dispatcher from '../stores/dispatcher';
import {getToken} from './defaultAxios';
import {snack} from './GlobalAction';
import * as UserAction from './UserAction';

export function fetchMessagesApp(id: number, since: number) {
    if (id === -1) {
        return axios
            .get(config.get('url') + 'message?since=' + since)
            .then((resp: AxiosResponse<IPagedMessages>) => {
                newMessages(-1, resp.data);
            });
    } else {
        return axios
            .get(config.get('url') + 'application/' + id + '/message?since=' + since)
            .then((resp: AxiosResponse<IPagedMessages>) => {
                newMessages(id, resp.data);
            });
    }
}

function newMessages(id: number, data: IPagedMessages) {
    dispatcher.dispatch({
        type: 'UPDATE_MESSAGES',
        payload: {
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
export function deleteMessagesByApp(id: number) {
    if (id === -1) {
        axios.delete(config.get('url') + 'message').then(() => {
            dispatcher.dispatch({type: 'DELETE_MESSAGES', payload: -1});
            snack('Messages deleted');
        });
    } else {
        axios.delete(config.get('url') + 'application/' + id + '/message').then(() => {
            dispatcher.dispatch({type: 'DELETE_MESSAGES', payload: id});
            snack('Deleted all messages from the application');
        });
    }
}

export function deleteMessage(msg: IMessage) {
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

    const wsUrl = config
        .get('url')
        .replace('http', 'ws')
        .replace('https', 'wss');
    const ws = new WebSocket(wsUrl + 'stream?token=' + getToken());

    ws.onerror = (e) => {
        wsActive = false;
        console.log('WebSocket connection errored', e);
    };

    ws.onmessage = (data) =>
        dispatcher.dispatch({type: 'ONE_MESSAGE', payload: JSON.parse(data.data) as IMessage});

    ws.onclose = () => {
        wsActive = false;
        UserAction.tryAuthenticate().then(() => {
            snack('WebSocket connection closed, trying again in 30 seconds.');
            setTimeout(listenToWebSocket, 30000);
        });
    };
}
