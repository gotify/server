import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import React, {Component} from 'react';
import {RouteComponentProps} from 'react-router';
import * as MessageAction from '../actions/MessageAction';
import DefaultPage from '../component/DefaultPage';
import ReactList from '../component/FixedReactList';
import Message from '../component/Message';
import AppStore from '../stores/AppStore';
import MessageStore from '../stores/MessageStore';

interface IProps extends RouteComponentProps<{id: string}> {}

interface IState {
    appId: number;
    messages: IMessage[];
    name: string;
    hasMore: boolean;
    nextSince?: number;
    id?: number;
}

class Messages extends Component<IProps, IState> {
    private static appId(props: IProps) {
        if (props === undefined) {
            return -1;
        }
        const {match} = props;
        return match.params.id !== undefined ? parseInt(match.params.id, 10) : -1;
    }

    public state = {appId: -1, messages: [], name: 'unknown', hasMore: true};

    private list: ReactList | null = null;

    public componentWillReceiveProps(nextProps: IProps) {
        this.updateAllWithProps(nextProps);
    }

    public componentWillMount() {
        MessageStore.on('change', this.updateAll);
        AppStore.on('change', this.updateAll);
        this.updateAll();
    }

    public componentWillUnmount() {
        MessageStore.removeListener('change', this.updateAll);
        AppStore.removeListener('change', this.updateAll);
    }

    public render() {
        const {name, messages, hasMore, appId} = this.state;
        const hasMessages = messages.length !== 0;
        const deleteMessages = () => MessageAction.deleteMessagesByApp(appId);

        return (
            <DefaultPage
                title={name}
                buttonTitle="Delete All"
                fButton={deleteMessages}
                buttonDisabled={!hasMessages}>
                {hasMessages ? (
                    <div style={{width: '100%'}}>
                        <ReactList
                            key={appId}
                            ref={(el: ReactList) => (this.list = el)}
                            itemRenderer={this.renderMessage}
                            length={messages.length}
                            threshold={1000}
                            pageSize={30}
                            type="variable"
                        />
                        {hasMore ? (
                            <Grid item xs={12} style={{textAlign: 'center'}}>
                                <CircularProgress size={100} />
                            </Grid>
                        ) : (
                            this.label("You've reached the end")
                        )}
                    </div>
                ) : (
                    this.label('No messages')
                )}
            </DefaultPage>
        );
    }

    private updateAllWithProps = (props: IProps) => {
        const appId = Messages.appId(props);

        const reset = MessageStore.shouldReset(appId);
        if (reset !== false && this.list) {
            this.list.clearCacheFromIndex(reset);
        }

        this.setState({...MessageStore.get(appId), appId, name: AppStore.getName(appId)});
        if (!MessageStore.exists(appId)) {
            MessageStore.loadNext(appId);
        }
    };

    private updateAll = () => this.updateAllWithProps(this.props);

    private deleteMessage = (message: IMessage) => () => MessageAction.deleteMessage(message);

    private renderMessage = (index: number, key: string) => {
        this.checkIfLoadMore();
        const message: IMessage = this.state.messages[index];
        return (
            <Message
                key={key}
                fDelete={this.deleteMessage(message)}
                title={message.title}
                date={message.date}
                content={message.message}
                image={message.image}
            />
        );
    };

    private checkIfLoadMore() {
        const {hasMore, messages, appId} = this.state;
        if (hasMore) {
            const [, maxRenderedIndex] = (this.list && this.list.getVisibleRange()) || [0, 0];
            if (maxRenderedIndex > messages.length - 30) {
                MessageStore.loadNext(appId);
            }
        }
    }

    private label = (text: string) => (
        <Grid item xs={12}>
            <Typography variant="caption" gutterBottom align="center">
                {text}
            </Typography>
        </Grid>
    );
}

export default Messages;
