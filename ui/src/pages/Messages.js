import React, {Component} from 'react';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';
import Message from '../component/Message';
import MessageStore from '../stores/MessageStore';
import AppStore from '../stores/AppStore';
import * as MessageAction from '../actions/MessageAction';
import DefaultPage from '../component/DefaultPage';
import ReactList from 'react-list';
import {CircularProgress} from 'material-ui/Progress';

class Messages extends Component {
    state = {appId: -1, messages: [], name: 'unknown', hasMore: true, list: null};

    componentWillReceiveProps(nextProps) {
        this.updateAllWithProps(nextProps);
    }

    componentWillMount() {
        MessageStore.on('change', this.updateAll);
        AppStore.on('change', this.updateAll);
        this.updateAll();
    }

    componentWillUnmount() {
        MessageStore.removeListener('change', this.updateAll);
        AppStore.removeListener('change', this.updateAll);
    }

    updateAllWithProps = (props) => {
        const appId = Messages.appId(props);
        this.setState({...MessageStore.get(appId), appId, name: AppStore.getName(appId)});
        if (!MessageStore.exists(appId)) {
            MessageStore.loadNext(appId);
        }
    };

    updateAll = () => this.updateAllWithProps(this.props);

    static appId(props) {
        if (props === undefined) {
            return -1;
        }
        const {match} = props;
        return match.params.id !== undefined ? parseInt(match.params.id, 10) : -1;
    }

    renderMessage = (index, key) => {
        this.checkIfLoadMore();
        const message = this.state.messages[index];
        return (
            <Message key={key}
                     fDelete={() => MessageAction.deleteMessage(message)}
                     title={message.title}
                     date={message.date}
                     content={message.message}
                     image={message.image}/>
        );
    };

    checkIfLoadMore() {
        const {hasMore, messages, appId} = this.state;
        if (hasMore) {
            const [, maxRenderedIndex] = (this.list && this.list.getVisibleRange()) || [0, 0];
            if (maxRenderedIndex > (messages.length - 30)) {
                MessageStore.loadNext(appId);
            }
        }
    }

    label = (text) => (
        <Grid item xs={12}><Typography variant="caption" gutterBottom align="center">{text}</Typography></Grid>
    );

    render() {
        const {name, messages, hasMore, appId} = this.state;
        const hasMessages = messages.length !== 0;
        const deleteMessages = () => MessageAction.deleteMessagesByApp(appId);

        return (
            <DefaultPage title={name} buttonTitle="Delete All" fButton={deleteMessages} buttonDisabled={!hasMessages}>
                {hasMessages
                    ? (
                        <div style={{width: '100%'}}>
                            <ReactList ref={(el) => this.list = el}
                                       itemRenderer={this.renderMessage}
                                       length={messages.length}
                                       threshold={1000}
                                       pageSize={30}
                                       type='variable'
                            />
                            {hasMore
                                ? <Grid item xs={12} style={{textAlign: 'center'}}><CircularProgress size={100}/></Grid>
                                : this.label('You\'ve reached the end')}
                        </div>
                    )
                    : this.label('No messages')
                }
            </DefaultPage>
        );
    }
}

export default Messages;
