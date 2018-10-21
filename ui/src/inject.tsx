import * as React from 'react';
import {UserStore} from './stores/UserStore';
import {SnackManager} from './snack/SnackManager';
import {MessagesStore} from './stores/MessagesStore';
import {CurrentUser} from './stores/CurrentUser';
import {ClientStore} from './stores/ClientStore';
import {AppStore} from './application/AppStore';
import {inject as mobxInject, Provider} from 'mobx-react';
import {WebSocketStore} from './stores/WebSocketStore';

export interface StoreMapping {
    userStore: UserStore;
    snackManager: SnackManager;
    messagesStore: MessagesStore;
    currentUser: CurrentUser;
    clientStore: ClientStore;
    appStore: AppStore;
    wsStore: WebSocketStore;
}

export type AllStores = Extract<keyof StoreMapping, string>;
export type Stores<T extends AllStores> = Pick<StoreMapping, T>;

export const inject = <I extends AllStores>(...stores: I[]) => {
    return <P extends {}>(
        node: React.ComponentType<P>
    ): React.ComponentType<Pick<P, Exclude<keyof P, I>>> => {
        return mobxInject(...stores)(node);
    };
};

export const InjectProvider: React.SFC<{stores: StoreMapping}> = ({children, stores}) => {
    return <Provider {...stores}>{children}</Provider>;
};
