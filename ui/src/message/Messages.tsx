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
import {PushMessageDialog} from './PushMessageDialog';
import {closeSnackbar, SnackbarKey} from 'notistack';

const UndoTimeoutMs = 5000;

interface PendingDelete {
    message: IMessage;
    allIndex: false | number;
    appIndex: false | number;
    timeoutId: number;
    snackKey: SnackbarKey;
}

const Messages = observer(() => {
    const {id} = useParams<{id: string}>();
    const appId = id == null ? -1 : parseInt(id as string, 10);

    const [deleteAll, setDeleteAll] = React.useState(false);
    const [pushMessageOpen, setPushMessageOpen] = React.useState(false);
    const [isLoadingMore, setLoadingMore] = React.useState(false);
    const {messagesStore, appStore, snackManager} = useStores();
    const messages = messagesStore.get(appId);
    const hasMore = messagesStore.canLoadMore(appId);
    const name = appStore.getName(appId);
    const hasMessages = messages.length !== 0;
    const expandedState = React.useRef<Record<number, boolean>>({});
    const pendingDeletesRef = React.useRef<Map<number, PendingDelete>>(new Map());
    const app = appId === -1 ? undefined : appStore.getByIDOrUndefined(appId);

    const deleteMessage = (message: IMessage) => () => startDelete(message);

    React.useEffect(() => {
        if (!messagesStore.loaded(appId)) {
            messagesStore.loadMore(appId);
        }
    }, [appId]);

    React.useEffect(() => () => clearPendingDeletes(), []);

    const clearPendingDeletes = (targetAppId?: number) => {
        pendingDeletesRef.current.forEach((pending, messageId) => {
            if (
                targetAppId != null &&
                targetAppId !== -1 &&
                pending.message.appid !== targetAppId
            ) {
                return;
            }
            window.clearTimeout(pending.timeoutId);
            closeSnackbar(pending.snackKey);
            messagesStore.clearPendingDelete(messageId);
            pendingDeletesRef.current.delete(messageId);
        });
    };

    const undoDelete = (messageId: number, snackKey: SnackbarKey) => {
        const pending = pendingDeletesRef.current.get(messageId);
        if (!pending) {
            return;
        }
        window.clearTimeout(pending.timeoutId);
        pendingDeletesRef.current.delete(messageId);
        messagesStore.clearPendingDelete(messageId);
        messagesStore.restoreSingleLocal(pending.message, pending.allIndex, pending.appIndex);
        closeSnackbar(snackKey);
        snackManager.snack('Delete undone');
    };

    const startDelete = (message: IMessage) => {
        if (pendingDeletesRef.current.has(message.id)) {
            return;
        }
        const {allIndex, appIndex} = messagesStore.removeSingleLocal(message);
        if (allIndex === false && appIndex === false) {
            return;
        }
        messagesStore.markPendingDelete(message.id);
        const snackKey = snackManager.snack('Message deleted', {
            action: (key) => (
                <Button color="inherit" size="small" onClick={() => undoDelete(message.id, key)}>
                    Undo
                </Button>
            ),
            autoHideDuration: UndoTimeoutMs,
        });
        const timeoutId = window.setTimeout(async () => {
            const pending = pendingDeletesRef.current.get(message.id);
            if (!pending) {
                return;
            }
            pendingDeletesRef.current.delete(message.id);
            messagesStore.clearPendingDelete(message.id);
            try {
                await messagesStore.removeSingle(message);
            } catch {
                messagesStore.restoreSingleLocal(message, pending.allIndex, pending.appIndex);
                snackManager.snack('Delete failed, message restored');
            }
        }, UndoTimeoutMs);
        pendingDeletesRef.current.set(message.id, {
            message,
            allIndex,
            appIndex,
            timeoutId,
            snackKey,
        });
    };

    const renderMessage = (_index: number, message: IMessage) => (
        <Message
            key={message.id}
            fDelete={deleteMessage(message)}
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
            endReached={checkIfLoadMore}
            data={messages}
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
            title={name}
            rightControl={
                <div>
                    {app && (
                        <Button
                            id="push-message"
                            variant="contained"
                            color="primary"
                            onClick={() => setPushMessageOpen(true)}
                            style={{marginRight: 5}}>
                            Push Message
                        </Button>
                    )}
                    <Button
                        id="refresh-all"
                        variant="contained"
                        color="primary"
                        onClick={() => {
                            clearPendingDeletes(appId);
                            messagesStore.refreshByApp(appId);
                        }}
                        style={{marginRight: 5}}>
                        Refresh
                    </Button>
                    <Button
                        id="delete-all"
                        variant="contained"
                        disabled={!hasMessages}
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
                    fOnSubmit={() => {
                        clearPendingDeletes(appId);
                        messagesStore.removeByApp(appId);
                    }}
                />
            )}
            {pushMessageOpen && app && (
                <PushMessageDialog
                    appName={app.name}
                    defaultPriority={app.defaultPriority}
                    fClose={() => setPushMessageOpen(false)}
                    fOnSubmit={(message, title, priority) =>
                        messagesStore.sendMessage(app.id, message, title, priority)
                    }
                />
            )}
        </DefaultPage>
    );
});

export default Messages;
