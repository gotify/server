import {BaseStore} from '../common/BaseStore';
import {action, IObservableArray, observable, reaction, runInAction} from 'mobx';
import * as config from '../config';
import {createTransformer} from 'mobx-utils';
import {SnackReporter} from '../snack/SnackManager';
import {IApplication, IMessage, IPagedMessages} from '../types';
import {closeSnackbar, SnackbarKey} from 'notistack';
import {identityTransform, jsonBody, jsonTransform} from '../fetchUtils';
import {CurrentUser} from '../CurrentUser';

const AllMessages = -1;

interface MessagesState {
    messages: IObservableArray<IMessage>;
    hasMore: boolean;
    nextSince: number;
    loaded: boolean;
}

interface PendingDelete {
    key: SnackbarKey;
    message: IMessage;
}

export class MessagesStore {
    @observable private accessor state: Record<string, MessagesState> = {};
    @observable private accessor pendingDeletes: Map<number, PendingDelete> = observable.map();

    private loading = false;

    public constructor(
        private readonly currentUser: CurrentUser,
        private readonly appStore: BaseStore<IApplication>,
        private readonly snack: SnackReporter
    ) {
        reaction(() => appStore.getItems(), this.createEmptyStatesForApps);
    }

    private stateOf = (appId: number, create = true) => {
        if (!this.state[appId] && create) {
            this.state[appId] = this.emptyState();
        }
        return this.state[appId] || this.emptyState();
    };

    public loaded = (appId: number) => this.stateOf(appId, /*create*/ false).loaded;

    public canLoadMore = (appId: number) => this.stateOf(appId, /*create*/ false).hasMore;

    @action
    public loadMore = async (appId: number) => {
        const state = this.stateOf(appId);
        if (!state.hasMore || this.loading) {
            return Promise.resolve();
        }
        this.loading = true;

        try {
            const pagedResult = await this.fetchMessages(appId, state.nextSince);
            runInAction(() => {
                state.messages.replace([...state.messages, ...pagedResult.messages]);
                state.nextSince = pagedResult.paging.since ?? 0;
                state.hasMore = 'next' in pagedResult.paging;
                state.loaded = true;
            });
        } finally {
            this.loading = false;
        }

        return Promise.resolve();
    };

    @action
    public publishSingleMessage = (message: IMessage) => {
        if (this.exists(AllMessages)) {
            this.stateOf(AllMessages).messages.unshift(message);
        }
        if (this.exists(message.appid)) {
            this.stateOf(message.appid).messages.unshift(message);
        }
    };

    @action
    public removeByApp = async (appId: number) => {
        if (appId === AllMessages) {
            await this.currentUser
                .authenticatedFetch(
                    config.get('url') + 'message',
                    {
                        method: 'DELETE',
                    },
                    identityTransform
                )
                .then(() => this.snack('Deleted all messages'));
            this.clearAll();
        } else {
            await this.currentUser
                .authenticatedFetch(
                    config.get('url') + 'application/' + appId + '/message',
                    {
                        method: 'DELETE',
                    },
                    identityTransform
                )
                .then(() =>
                    this.snack(`Deleted all messages from ${this.appStore.getByID(appId).name}`)
                );
            this.clear(AllMessages);
            this.clear(appId);
        }
        await this.loadMore(appId);
    };

    @action
    public addPendingDelete = (pending: PendingDelete) =>
        this.pendingDeletes.set(pending.message.id, pending);

    @action
    public cancelPendingDelete = (message: IMessage): boolean => {
        const pending = this.pendingDeletes.get(message.id);
        if (pending) {
            this.pendingDeletes.delete(message.id);
            closeSnackbar(pending.key);
        }
        return !!pending;
    };

    @action
    public executePendingDeletes = () =>
        Array.from(this.pendingDeletes.values()).forEach(({message}) => this.removeSingle(message));

    public visible = (message: number): boolean => !this.pendingDeletes.has(message);

    @action
    public removeSingle = async (message: IMessage) => {
        if (!this.pendingDeletes.has(message.id)) {
            return;
        }

        await this.currentUser
            .authenticatedFetch(
                config.get('url') + 'message/' + message.id,
                {
                    method: 'DELETE',
                    keepalive: true,
                },
                identityTransform
            )
            .then(() => this.snack(`Deleted message ${message.id}`));
        if (this.exists(AllMessages)) {
            this.removeFromList(this.state[AllMessages].messages, message);
        }
        if (this.exists(message.appid)) {
            this.removeFromList(this.state[message.appid].messages, message);
        }
        this.cancelPendingDelete(message);
    };

    public sendMessage = async (
        appId: number,
        message: string,
        title: string,
        priority: number
    ): Promise<void> => {
        const app = this.appStore.getByID(appId);
        const payload: Pick<IMessage, 'title' | 'message' | 'priority'> = {
            message,
            priority,
            title,
        };

        const fetchInit = jsonBody(payload);
        fetchInit.headers = new Headers(fetchInit.headers);
        fetchInit.headers.set('X-Gotify-Key', app.token);
        await this.currentUser
            .authenticatedFetch(config.get('url') + 'message', fetchInit, jsonTransform)
            .then(() => this.snack(`Message sent to ${app.name}`));
    };

    @action
    public clearAll = () => {
        this.state = {};
        this.createEmptyStatesForApps(this.appStore.getItems());
    };

    @action
    public refreshByApp = async (appId: number) => {
        this.clearAll();
        this.loadMore(appId);
    };

    public exists = (id: number) => this.stateOf(id).loaded;

    @action
    private removeFromList(messages: IMessage[], messageToDelete: IMessage): false | number {
        if (messages) {
            const index = messages.findIndex((message) => message.id === messageToDelete.id);
            if (index !== -1) {
                messages.splice(index, 1);
                return index;
            }
        }
        return false;
    }

    @action
    private clear = (appId: number) => (this.state[appId] = this.emptyState());

    private fetchMessages = (appId: number, since: number): Promise<IPagedMessages> => {
        if (appId === AllMessages) {
            return this.currentUser.authenticatedFetch(
                config.get('url') + 'message?since=' + since,
                {},
                jsonTransform<IPagedMessages>
            );
        } else {
            return this.currentUser.authenticatedFetch(
                config.get('url') + 'application/' + appId + '/message?since=' + since,
                {},
                jsonTransform<IPagedMessages>
            );
        }
    };

    private getUnCached = (appId: number): Array<IMessage> => {
        const appToImage: Partial<Record<string, string>> = this.appStore
            .getItems()
            .reduce((all, app) => ({...all, [app.id]: app.image}), {});

        return this.stateOf(appId, false)
            .messages.filter((message) => !this.pendingDeletes.has(message.id))
            .map((message: IMessage): IMessage => ({...message, image: appToImage[message.appid]}));
    };

    public get = createTransformer(this.getUnCached);

    private clearCache = () => (this.get = createTransformer(this.getUnCached));

    private createEmptyStatesForApps = (apps: IApplication[]) => {
        apps.map((app) => app.id).forEach((id) => this.stateOf(id, /*create*/ true));
        this.clearCache();
    };

    private emptyState = (): MessagesState => ({
        messages: observable.array(),
        hasMore: true,
        nextSince: 0,
        loaded: false,
    });
}
