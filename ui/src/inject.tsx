import * as React from 'react';
import {UserStore} from './user/UserStore';
import {SnackManager} from './snack/SnackManager';
import {MessagesStore} from './message/MessagesStore';
import {CurrentUser} from './CurrentUser';
import {ClientStore} from './client/ClientStore';
import {AppStore} from './application/AppStore';
import {inject as mobxInject, Provider} from 'mobx-react';
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

export type AllStores = Extract<keyof StoreMapping, string>;
export type Stores<T extends AllStores> = Pick<StoreMapping, T>;

export const inject =
    <I extends AllStores>(...stores: I[]) =>
    // eslint-disable-next-line @typescript-eslint/ban-types
    <P extends {}>(
        node: React.ComponentType<P>
    ): React.ComponentType<Pick<P, Exclude<keyof P, I>>> =>
        // eslint-disable-next-line @typescript-eslint/no-explicit-any
        mobxInject(...stores)(node) as any;

export const InjectProvider: React.FC<{stores: StoreMapping}> = ({children, stores}) => (
    <Provider {...stores}>{children}</Provider>
);
