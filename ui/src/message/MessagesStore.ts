import {BaseStore} from '../common/BaseStore';
import {action, IObservableArray, observable, reaction, transaction} from 'mobx';
import axios, {AxiosResponse} from 'axios';
import * as config from '../config';
import {chunkProcessor, createTransformer} from 'mobx-utils';
import {SnackReporter} from '../snack/SnackManager';

const AllMessages = -1;

interface MessagesState {
    messages: IObservableArray<IMessage>;
    hasMore: boolean;
    nextSince: number;
    loaded: boolean;
}

export class MessagesStore {
    @observable
    private state: Record<number, MessagesState> = {};

    @observable
    private readMessageQueue: number[] = [];

    @observable
    private newMessageQueue: IMessage[] = [];

    private loading = false;

    public constructor(
        private readonly appStore: BaseStore<IApplication>,
        private readonly snack: SnackReporter
    ) {
        reaction(() => appStore.getItems(), this.createEmptyStatesForApps);
        chunkProcessor(this.readMessageQueue, this.markMessagesAsReadRemote, 1000);
        chunkProcessor(this.newMessageQueue, (messages) => messages.map(this.showNewMessage), 200);
    }

    private markMessagesAsReadRemote = (ids: number[]) => {
        axios.post(config.get('url') + 'message/read?' + ids.map((id) => `id=${id}`).join('&'));
    };

    private stateOf = (appId: number, create = true) => {
        if (this.state[appId] || !create) {
            return this.state[appId] || this.emptyState();
        }
        return (this.state[appId] = this.emptyState());
    };

    public canLoadMore = (appId: number) => this.stateOf(appId, /*create*/ false).hasMore;

    @action
    public loadMore = async (appId: number) => {
        const state = this.stateOf(appId);
        if (!state.hasMore || this.loading) {
            return Promise.resolve();
        }
        this.loading = true;

        const pagedResult = await this.fetchMessages(appId, state.nextSince).then(
            (resp) => resp.data
        );
        transaction(() => {
            state.loaded = true;
            state.nextSince = pagedResult.paging.since || 0;
            state.hasMore = 'next' in pagedResult.paging;
            state.messages.replace([...state.messages, ...pagedResult.messages]);
            this.loading = false;
        });
        return Promise.resolve();
    };

    @action
    public markAsRead = (message: IMessage) => {
        console.log('MARK AS READ', this.stateOf(-1).loaded);

        if (this.exists(AllMessages)) {
            console.log('ALl messages');

            this.markAsReadWithMessages(this.state[AllMessages].messages, message);
        }
        if (this.exists(message.appid)) {
            this.markAsReadWithMessages(this.state[message.appid].messages, message);
        }
        this.readMessageQueue.push(message.id);
    };

    private markAsReadWithMessages = (messages: IMessage[], toUpdate: IMessage) => {
        const foundMessage = messages.find((message) => toUpdate.id === message.id);
        console.log('UPDATE MESSAGE', foundMessage);

        if (foundMessage) {
            foundMessage.read = true;
        }
    };

    @action
    public publishSingleMessage = (message: IMessage) => {
        this.newMessageQueue.push(message);
    };

    @action
    public showNewMessage = (message: IMessage) => {
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
            await axios.delete(config.get('url') + 'message');
            this.snack('Deleted all messages');
            this.clearAll();
        } else {
            await axios.delete(config.get('url') + 'application/' + appId + '/message');
            this.snack(`Deleted all messages from ${this.appStore.getByID(appId).name}`);
            this.clear(AllMessages);
            this.clear(appId);
        }
        await this.loadMore(appId);
    };

    @action
    public removeSingle = async (message: IMessage) => {
        await axios.delete(config.get('url') + 'message/' + message.id);
        if (this.exists(AllMessages)) {
            this.removeFromList(this.state[AllMessages].messages, message);
        }
        if (this.exists(message.appid)) {
            this.removeFromList(this.state[message.appid].messages, message);
        }
        this.snack('Message deleted');
    };

    @action
    public clearAll = () => {
        this.state = {};
        this.createEmptyStatesForApps(this.appStore.getItems());
    };

    public exists = (id: number) => this.stateOf(id).loaded;

    private removeFromList(messages: IMessage[], messageToDelete: IMessage) {
        if (messages) {
            const index = messages.findIndex((message) => message.id === messageToDelete.id);
            if (index !== -1) {
                messages.splice(index, 1);
            }
        }
    }

    private clear = (appId: number) => (this.state[appId] = this.emptyState());

    private fetchMessages = (
        appId: number,
        since: number
    ): Promise<AxiosResponse<IPagedMessages>> => {
        if (appId === AllMessages) {
            return axios.get(config.get('url') + 'message?since=' + since);
        } else {
            return axios.get(
                config.get('url') + 'application/' + appId + '/message?since=' + since
            );
        }
    };

    private getUnCached = (appId: number): Array<IMessage & {image: string}> => {
        const appToImage = this.appStore
            .getItems()
            .reduce((all, app) => ({...all, [app.id]: app.image}), {});

        return this.stateOf(appId, false).messages.map((message: IMessage) => {
            return {
                ...message,
                image: appToImage[message.appid] || 'still loading',
            };
        });
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
