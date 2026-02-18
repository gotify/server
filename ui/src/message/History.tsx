import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import React from 'react';
import {useParams} from 'react-router';
import DefaultPage from '../common/DefaultPage';
import Button from '@mui/material/Button';
import Message from './Message';
import {observer} from 'mobx-react-lite';
import {IMessage} from '../types';
import ConfirmDialog from '../common/ConfirmDialog';
import LoadingSpinner from '../common/LoadingSpinner';
import {useStores} from '../stores';
import {Virtuoso} from 'react-virtuoso';
import {enqueueSnackbar} from 'notistack';

const UndoAutoHideMs = 250;

/***************************************************************************
 * Page Name: Deleted Messages
 * Right Control Buttons: Refresh, Delete All
 * URL: ../history (routing is located in ui/src/layout/Layout.tsx)
 * 
 * Function:
 * The purpose of this file is to generate the UI regarding viewing deleted
 * messages (essentially implementing a soft-delete), with the capability to 
 * permanently delete, and set an expiration timer for auto-deletion 
 * (not currently supported)
****************************************************************************/

const History = observer(() => {
    const {id} = useParams<{id: string}>();
    const appId = id == null ? -1 : parseInt(id as string, 10);

    const [deleteAll, setDeleteAll] = React.useState(false);
    const [isLoadingMore, setLoadingMore] = React.useState(false);
    const {messagesStore, appStore} = useStores();
    const expandedState = React.useRef<Record<number, boolean>>({});
    const pendingMessages = messagesStore.getPending(appId);
    const hasPendingMessages = pendingMessages.length !== 0;
    const hasMore = messagesStore.canLoadMore(appId);           

    const deleteMessage = (message: IMessage) => {
        const key = enqueueSnackbar({
            message: 'Message deleted permanently',
            variant: 'info',
            disableWindowBlurListener: true,
            transitionDuration: {enter: 0, exit: 0},
            autoHideDuration: UndoAutoHideMs,
            onExited: () => messagesStore.removeSingle(message),
        });
    };

    React.useEffect(() => {
        if (!messagesStore.loaded(appId)) {
            messagesStore.loadMore(appId);
        }
    }, [appId]);

    const renderMessage = (_index: number, message: IMessage) => (
        <Message
            key={message.id}
            fDelete={() => deleteMessage(message)}
            fRestore={() => messagesStore.cancelPendingDelete(message)}
            pending={true}
            onExpand={(expanded) => (expandedState.current[message.id] = expanded)}
            title={message.title}
            date={message.date}
            appName={appStore.getName(message.appid)}
            expanded={expandedState.current[message.id] ?? false}
            content={message.message}
            image={message.image}
            extras={message.extras}
            priority={message.priority}
        />
    );

    const checkIfLoadMore = () => {
        if (!isLoadingMore && messagesStore.canLoadMore(appId)) {
            setLoadingMore(true);
            messagesStore.loadMore(appId).then(() => setLoadingMore(false));
        }
    };

    const messageFooter = () => {
        if (hasMore) {
            return <LoadingSpinner />;
        }
        if (hasPendingMessages) {
            return label("You've reached the end");
        }
        return null;
    };

    const renderMessages = () => (
        <Virtuoso
            id="messages"
            style={{width: '100%'}}
            useWindowScroll
            totalCount={pendingMessages.length}
            endReached={checkIfLoadMore}            // Need to make sure this doesn't cause problems
            data={pendingMessages}
            itemContent={renderMessage}
            components={{
                Footer: messageFooter,
                EmptyPlaceholder: () => label('No messages'),
            }}
        />
    );
    const label = (text: string) => (
        <Grid size={{xs: 12}}>
            <Typography variant="caption" component="div" gutterBottom align="center">
                {text}
            </Typography>
        </Grid>
    );
    return (
        <DefaultPage
            title="Deleted Messages"
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
                        disabled={!hasPendingMessages}
                        color="primary"
                        onClick={() => {
                            setDeleteAll(true);
                        }}>
                        Delete All
                    </Button>
                </div>
            }>
            {!messagesStore.loaded(appId) ? <LoadingSpinner /> : renderMessages()}

            {deleteAll && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete all messages?'}
                    fClose={() => setDeleteAll(false)}
                    fOnSubmit={() => messagesStore.removeByApp(appId)}
                />
            )}
        </DefaultPage>
    );
});

export default History;