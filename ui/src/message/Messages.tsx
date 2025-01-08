import Grid from '@mui/material/Grid2';
import Typography from '@mui/material/Typography';
import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router';
import {Virtuoso} from 'react-virtuoso';
import DefaultPage from '../common/DefaultPage';
import Button from '@mui/material/Button';
import {useAppDispatch, useAppSelector} from '../store';
import {getAppName} from '../application/app-actions.ts';
import {uiActions} from '../store/ui-slice.ts';
import {IMessage} from '../types.ts';
import {fetchMessages, removeMessagesByApp, removeSingleMessage} from './message-actions.ts';
import Message from './Message';
import ConfirmDialog from '../common/ConfirmDialog';
import LoadingSpinner from '../common/LoadingSpinner';

const Messages = () => {
    const dispatch = useAppDispatch();
    const {id} = useParams();
    const appId = id !== undefined ? parseInt(id as string) : -1;

    const [toDeleteAll, setToDeleteAll] = useState<boolean>(false);

    const reloadRequired = useAppSelector((state) => state.ui.reloadRequired);
    const selectedApp = useAppSelector((state) => state.app.items.find((app) => app.id === appId));
    const apps = useAppSelector((state) => state.app.items);
    const messages = useAppSelector((state) =>
        appId === -1
            ? state.message.items
            : state.message.items.filter((item) => item.appid === appId)
    );
    const hasMore = useAppSelector((state) => state.message.hasMore);
    const name = dispatch(getAppName(appId));
    const messagesLoaded = useAppSelector((state) => state.message.loaded);
    const hasMessages = messages.length !== 0;

    // handle a requested reload
    useEffect(() => {
        if (reloadRequired) {
            dispatch(uiActions.setReloadRequired(false));
            dispatch(fetchMessages());
        }
    }, [dispatch, reloadRequired]);

    useEffect(() => {
        window.onscroll = () => {
            if (
                window.innerHeight + window.scrollY >=
                document.body.offsetHeight - window.innerHeight * 2
            ) {
                checkIfLoadMore();
            }
        };

        dispatch(fetchMessages());
    }, [dispatch]);

    const checkIfLoadMore = () => {
        console.log('checkIfLoadMore');
    };

    const label = (text: string) => (
        <Grid size={12}>
            <Typography variant="caption" component="div" gutterBottom align="center">
                {text}
            </Typography>
        </Grid>
    );

    const messageFooter = () => {
        if (hasMore) {
            return <LoadingSpinner />;
        }
        if (hasMessages) {
            return label("You've reached the end");
        }
        return null;
    };

    const renderMessages = () => (
        <Virtuoso
            id="messages"
            style={{width: '100%'}}
            useWindowScroll
            totalCount={messages.length}
            data={messages}
            itemContent={(index, message) => renderMessage(index, message)}
            components={{
                Footer: messageFooter,
                EmptyPlaceholder: () => label('No messages'),
            }}
        />
    );

    const renderMessage = (index: number, message: IMessage) => (
        <Message
            key={index}
            fDelete={() => dispatch(removeSingleMessage(message))}
            title={message.title}
            date={message.date}
            content={message.message}
            image={apps.find((app) => app.id == message.appid)?.image}
            extras={message.extras}
            priority={message.priority}
        />
    );

    return (
        <DefaultPage
            title={name}
            rightControl={
                <div>
                    <Button
                        id="refresh-all"
                        variant="contained"
                        color="primary"
                        onClick={() => dispatch(fetchMessages(appId))}
                        style={{marginRight: 5}}>
                        Refresh
                    </Button>
                    <Button
                        id="delete-all"
                        variant="contained"
                        disabled={!hasMessages}
                        color="primary"
                        onClick={() => setToDeleteAll(true)}>
                        Delete All
                    </Button>
                </div>
            }>
            {!messagesLoaded ? <LoadingSpinner /> : renderMessages()}

            {toDeleteAll && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete all messages?'}
                    fClose={() => setToDeleteAll(false)}
                    fOnSubmit={() => dispatch(removeMessagesByApp(selectedApp))}
                />
            )}
        </DefaultPage>
    );
};
/*
@observer
class Messages_old extends Component<IProps & Stores<'messagesStore' | 'appStore'>, IState> {
    private isLoadingMore = false;

    private updateAllWithProps = (props: IProps & Stores<'messagesStore'>) => {
        const appId = Messages.appId(props);
        this.setState({appId});
        if (!props.messagesStore.exists(appId)) {
            props.messagesStore.loadMore(appId);
        }
    };

    private checkIfLoadMore() {
        const {appId} = this.state;
        if (!this.isLoadingMore && this.props.messagesStore.canLoadMore(appId)) {
            this.isLoadingMore = true;
            this.props.messagesStore.loadMore(appId).then(() => (this.isLoadingMore = false));
        }
    }
}
*/

export default Messages;
