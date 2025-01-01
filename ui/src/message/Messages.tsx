import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import React, {Component} from 'react';
import {RouteComponentProps} from 'react-router';
import DefaultPage from '../common/DefaultPage';
import Button from '@material-ui/core/Button';
import Message from './Message';
import {observer} from 'mobx-react';
import {inject, Stores } from '../inject';
import { action, makeObservable, observable } from 'mobx';
import ReactInfinite from 'react-infinite';
import { IMessage } from '../types';
import ConfirmDialog from '../common/ConfirmDialog';
import LoadingSpinner from '../common/LoadingSpinner';

type IProps = RouteComponentProps<{ id: string }>;

interface IState {
    appId: number;
}

@observer
class Messages extends Component<IProps & Stores<'messagesStore' | 'appStore'>, IState> {
    @observable
    private heights: Record<string, number> = {};
    @observable
    private deleteAll = false;

    constructor(props: any) {
        super(props);
        makeObservable(this);
    }

    @action
    private setHeight(id: string, height: number) {
        this.heights[id] = height;
    }

    @action
    private setDeleteAll(deleteAll: boolean) {
        this.deleteAll = deleteAll;
    }

    private static appId(props: IProps) {
        if (props === undefined) {
            return -1;
        }
        const {match} = props;
        return match.params.id !== undefined ? parseInt(match.params.id, 10) : -1;
    }

    public state = {appId: -1};

    private isLoadingMore = false;

    public componentWillReceiveProps(nextProps: IProps & Stores<'messagesStore' | 'appStore'>) {
        this.updateAllWithProps(nextProps);
    }

    public componentWillMount() {
        window.onscroll = () => {
            if (
                window.innerHeight + window.pageYOffset >=
                document.body.offsetHeight - window.innerHeight * 2
            ) {
                this.checkIfLoadMore();
            }
        };
        this.updateAll();
    }

    public render() {
        const {appId} = this.state;
        const {messagesStore, appStore} = this.props;
        const messages = messagesStore.get(appId);
        const hasMore = messagesStore.canLoadMore(appId);
        const name = appStore.getName(appId);
        const hasMessages = messages.length !== 0;

        return (
            <DefaultPage
                title={name}
                rightControl={
                    <div>
                        <Button
                            id="refresh-all"
                            variant="contained"
                            color="primary"
                            onClick={() => messagesStore.refreshByApp(appId)}
                            style={{marginRight: 5}}>
                            Refresh
                        </Button>
                        <Button
                            id="delete-all"
                            variant="contained"
                            disabled={!hasMessages}
                            color="primary"
                            onClick={() => {
                                this.setDeleteAll(true);
                            }}>
                            Delete All
                        </Button>
                    </div>
                }>
                {!messagesStore.loaded(appId) ? (
                    <LoadingSpinner />
                ) : hasMessages ? (
                    <div style={{width: '100%'}} id="messages">
                        <ReactInfinite
                            key={appId}
                            useWindowAsScrollContainer
                            preloadBatchSize={window.innerHeight * 3}
                            elementHeight={messages.map((m) => this.heights[m.id] || 1)}>
                            {messages.map(this.renderMessage)}
                        </ReactInfinite>

                        {hasMore ? <LoadingSpinner /> : this.label("You've reached the end")}
                    </div>
                ) : (
                    this.label('No messages')
                )}

                {this.deleteAll && (
                    <ConfirmDialog
                        title="Confirm Delete"
                        text={'Delete all messages?'}
                        fClose={() => (this.setDeleteAll(false))}
                        fOnSubmit={() => messagesStore.removeByApp(appId)}
                    />
                )}
            </DefaultPage>
        );
    }

    private updateAllWithProps = (props: IProps & Stores<'messagesStore'>) => {
        const appId = Messages.appId(props);
        this.setState({appId});
        if (!props.messagesStore.exists(appId)) {
            props.messagesStore.loadMore(appId);
        }
    };

    private updateAll = () => this.updateAllWithProps(this.props);

    private deleteMessage = (message: IMessage) => () =>
        this.props.messagesStore.removeSingle(message);

    private renderMessage = (message: IMessage) => (
        <Message
            key={message.id}
            height={(height: number) => {
                if (!this.heights[message.id]) {
                    this.setHeight(message.id, height);
                }
            }}
            fDelete={this.deleteMessage(message)}
            title={message.title}
            date={message.date}
            content={message.message}
            image={message.image}
            extras={message.extras}
            priority={message.priority}
        />
    );

    private checkIfLoadMore() {
        const {appId} = this.state;
        if (!this.isLoadingMore && this.props.messagesStore.canLoadMore(appId)) {
            this.isLoadingMore = true;
            this.props.messagesStore.loadMore(appId).then(() => (this.isLoadingMore = false));
        }
    }

    private label = (text: string) => (
        <Grid item xs={12}>
            <Typography variant="caption" component="div" gutterBottom align="center">
                {text}
            </Typography>
        </Grid>
    );
}

export default inject('messagesStore', 'appStore')(Messages);
