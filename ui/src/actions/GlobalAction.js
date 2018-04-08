import * as AppAction from './AppAction';
import * as UserAction from './UserAction';
import * as MessageAction from './MessageAction';
import * as ClientAction from './ClientAction';
import dispatcher from '../stores/dispatcher';

export function initialLoad(resp) {
    AppAction.fetchApps();
    MessageAction.listenToWebSocket();
    ClientAction.fetchClients();
    if (resp.data.admin) {
        UserAction.fetchUsers();
    }
}

export function snack(message) {
    dispatcher.dispatch({type: 'SNACK', payload: message});
}
