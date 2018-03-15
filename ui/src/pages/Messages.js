import React, {Component} from 'react';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';
import Message from '../component/Message';
import MessageStore from '../stores/MessageStore';
import AppStore from '../stores/AppStore';
import * as MessageAction from '../actions/MessageAction';
import DefaultPage from '../component/DefaultPage';

class Messages extends Component {
    constructor() {
        super();
        this.state = {appId: -1, messages: [], name: 'unknown'};
    }

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
        const messages = MessageStore.getForAppId(appId);
        const name = AppStore.getName(appId);
        this.setState({appId: appId, messages, name});
    };

    updateAll = () => this.updateAllWithProps(this.props);

    static appId(props) {
        if (props === undefined) {
            return -1;
        }
        const {match} = props;
        return match.params.id !== undefined ? parseInt(match.params.id) : -1;
    }

    render() {
        const {name, messages, appId} = this.state;
        const fDelete = appId === -1 ? MessageAction.deleteMessages : MessageAction.deleteMessagesByApp.bind(this, appId);

        const noMessages = (
            <Grid item xs={12}><Typography variant="caption" gutterBottom align="center">No messages</Typography></Grid>
        );

        return (
            <DefaultPage title={name} buttonTitle="Delete All" fButton={fDelete} buttonDisabled={messages.length === 0}>
                {messages.length === 0 ? noMessages : messages.map((message) => {
                    return (
                        <Grid item xs={12} key={message.id}>
                            <Message fDelete={() => MessageAction.deleteMessage(message.id)} title={message.title}
                                     date={message.date} content={message.message}/>
                        </Grid>
                    );
                })}
            </DefaultPage>
        );
    }
}

export default Messages;
