import {AxiosResponse} from "axios";
import dispatcher from '../stores/dispatcher';
import * as AppAction from './AppAction';
import * as ClientAction from './ClientAction';
import * as MessageAction from './MessageAction';
import * as UserAction from './UserAction';

export function initialLoad(resp: AxiosResponse<IUser>) {
    AppAction.fetchApps();
    MessageAction.listenToWebSocket();
    ClientAction.fetchClients();
    if (resp.data.admin) {
        UserAction.fetchUsers();
    }
}

export function snack(message: string) {
    dispatcher.dispatch({type: 'SNACK', payload: message});
}
