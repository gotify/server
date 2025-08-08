import * as React from 'react';
import {UserStore} from './user/UserStore';
import {SnackManager} from './snack/SnackManager';
import {MessagesStore} from './message/MessagesStore';
import {CurrentUser} from './CurrentUser';
import {ClientStore} from './client/ClientStore';
import {AppStore} from './application/AppStore';
import {WebSocketStore} from './message/WebSocketStore';
import {PluginStore} from './plugin/PluginStore';

export interface StoreMapping {
    userStore: UserStore;
    snackManager: SnackManager;
    messagesStore: MessagesStore;
    currentUser: CurrentUser;
    clientStore: ClientStore;
    appStore: AppStore;
    pluginStore: PluginStore;
    wsStore: WebSocketStore;
}

export const StoreContext = React.createContext<StoreMapping | undefined>(undefined);

export const useStores = (): StoreMapping => {
    const mapping = React.useContext(StoreContext);
    if (!mapping) throw new Error('uninitialized');
    return mapping;
};
