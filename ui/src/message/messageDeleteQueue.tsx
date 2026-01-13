import Button from '@mui/material/Button';
import React from 'react';
import {closeSnackbar, SnackbarKey} from 'notistack';
import {SnackReporter} from '../snack/SnackManager';
import {IMessage} from '../types';
import {MessagesStore} from './MessagesStore';

const UndoAutoHideMs = 5000;

interface PendingDelete {
    message: IMessage;
    allIndex: false | number;
    appIndex: false | number;
    snackKey: SnackbarKey;
}

class MessageDeleteQueue {
    private pendingDeletes = new Map<number, PendingDelete>();
    private pagehideBound = false;

    public constructor(
        private readonly messagesStore: MessagesStore,
        private readonly snack: SnackReporter
    ) {}

    public requestDelete = (message: IMessage) => {
        this.ensurePagehideHandler();
        if (this.pendingDeletes.has(message.id)) {
            return;
        }
        const {allIndex, appIndex} = this.messagesStore.removeSingleLocal(message);
        if (allIndex === false && appIndex === false) {
            return;
        }
        this.messagesStore.markPendingDelete(message.id);

        const snackKey = this.snack('Message deleted', {
            action: (key) => (
                <Button
                    color="inherit"
                    size="small"
                    onClick={() => this.undoDelete(message.id, key)}>
                    Undo
                </Button>
            ),
            autoHideDuration: UndoAutoHideMs,
            onExited: () => {
                void this.finalizeDelete(message.id);
            },
        });

        this.pendingDeletes.set(message.id, {
            message,
            allIndex,
            appIndex,
            snackKey,
        });
    };

    public undoDelete = (messageId: number, snackKey: SnackbarKey) => {
        const pending = this.pendingDeletes.get(messageId);
        if (!pending) {
            return;
        }
        this.pendingDeletes.delete(messageId);
        this.messagesStore.clearPendingDelete(messageId);
        this.messagesStore.restoreSingleLocal(pending.message, pending.allIndex, pending.appIndex);
        closeSnackbar(snackKey);
        this.snack('Delete undone');
    };

    public finalizePendingDeletes = (targetAppId?: number) => {
        const pendingIds = Array.from(this.pendingDeletes.keys());
        pendingIds.forEach((messageId) => {
            const pending = this.pendingDeletes.get(messageId);
            if (!pending) {
                return;
            }
            if (
                targetAppId != null &&
                targetAppId !== -1 &&
                pending.message.appid !== targetAppId
            ) {
                return;
            }
            void this.finalizeDelete(messageId, {closeSnack: true});
        });
    };

    private ensurePagehideHandler = () => {
        if (this.pagehideBound || typeof window === 'undefined') {
            return;
        }
        this.pagehideBound = true;
        window.addEventListener('pagehide', this.handlePagehide);
        window.addEventListener('beforeunload', this.handlePagehide);
    };

    private handlePagehide = () => {
        this.finalizePendingDeletes();
    };

    private finalizeDelete = async (
        messageId: number,
        options?: {closeSnack?: boolean}
    ): Promise<void> => {
        const pending = this.pendingDeletes.get(messageId);
        if (!pending) {
            return;
        }
        this.pendingDeletes.delete(messageId);
        this.messagesStore.clearPendingDelete(messageId);
        if (options?.closeSnack) {
            closeSnackbar(pending.snackKey);
        }
        try {
            await this.messagesStore.removeSingle(pending.message, {keepalive: true});
        } catch {
            this.messagesStore.restoreSingleLocal(
                pending.message,
                pending.allIndex,
                pending.appIndex
            );
            this.snack('Delete failed, message restored');
        }
    };
}

let messageDeleteQueue: MessageDeleteQueue | null = null;

export const getMessageDeleteQueue = (
    messagesStore: MessagesStore,
    snack: SnackReporter
): MessageDeleteQueue => {
    if (!messageDeleteQueue) {
        messageDeleteQueue = new MessageDeleteQueue(messagesStore, snack);
    }
    return messageDeleteQueue;
};
