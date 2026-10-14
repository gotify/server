// Regression test for gotify/server#1068:
// "WebUI sends duplicate DELETE /message/{id} when leaving the page with
// pending deletes".
//
// Deleting a message in the WebUI opens an undo snackbar (`Messages.tsx`
// `deleteMessage` -> `enqueueSnackbar` + `addPendingDelete`).  While that
// snackbar is still open, leaving the page must issue exactly ONE DELETE per
// message id.
//
// `reactions.ts:8-9` registers `messagesStore.executePendingDeletes` on BOTH
// `pagehide` and `beforeunload`, and browsers fire both when a page is left.
// `MessagesStore.removeSingle` checked `pendingDeletes.has(id)` *before*
// awaiting the request and only removed the entry *after* the await resolved,
// so both handlers passed the guard and each sent a DELETE; the second one
// 404s with "message does not exist".
import axios from 'axios';
import {afterEach, describe, expect, it, vi} from 'vitest';

vi.mock('notistack', () => ({closeSnackbar: () => undefined}));

const {MessagesStore} = await import('../message/MessagesStore');
import type {IApplication, IMessage} from '../types';

const deleteUrls = (): string[] => {
    const calls = (axios.delete as unknown as {mock: {calls: unknown[][]}}).mock.calls;
    return calls.map((call) => String(call[0]));
};

const app = (id: number): IApplication => ({
    id,
    name: `app-${id}`,
    description: '',
    image: '',
    internal: false,
    token: `token-${id}`,
    sortKey: '',
    defaultPriority: 0,
    lastUsed: '2026-10-14T00:00:00Z',
    createdAt: '2026-10-14T00:00:00Z',
});

const message = (id: number): IMessage => ({
    id,
    appid: 1,
    message: `body-${id}`,
    title: `title-${id}`,
    priority: 1,
    date: '2026-10-14T00:00:00Z',
});

const newStore = () =>
    new MessagesStore({getItems: () => [app(1)], getByID: () => app(1)} as never, () => undefined);

describe('pending deletes on page leave (#1068)', () => {
    afterEach(() => {
        vi.restoreAllMocks();
    });

    it('sends exactly one DELETE per id when leaving the page with a pending undo', async () => {
        vi.spyOn(axios, 'delete').mockResolvedValue({status: 200} as never);

        const store = newStore();
        store.addPendingDelete({key: 'undo' as never, message: message(150242)});

        // A browser fires pagehide and beforeunload while the previous async
        // handler is still in flight, so both run concurrently.
        await Promise.all([store.executePendingDeletes(), store.executePendingDeletes()]);

        expect(deleteUrls()).toEqual(['unsetmessage/150242']);
        expect(axios.delete).toHaveBeenCalledTimes(1);
    });

    it('does not re-send a DELETE for a message already being deleted', async () => {
        vi.spyOn(axios, 'delete').mockResolvedValue({status: 200} as never);

        const store = newStore();
        store.addPendingDelete({key: 'undo' as never, message: message(7)});

        // The snackbar's `onExited` (Messages.tsx) and a page-leave event can
        // both fire for the same message.
        const leaving = store.executePendingDeletes();
        const expired = store.removeSingle(message(7));
        await Promise.all([leaving, expired]);

        expect(axios.delete).toHaveBeenCalledTimes(1);
    });
});
