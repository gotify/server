import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IMessage, IPagedMessages} from '../types.ts';
import {AllMessages} from './message-actions.ts';

interface MessagesState {
    items: IMessage[];
    hasMore: boolean;
    nextSince: number;
    loaded: boolean;
}

const initialMessageState: MessagesState = {
    items: [],
    hasMore: false,
    nextSince: 0,
    loaded: false,
}

export const messageSlice = createSlice({
    name: 'message',
    initialState: initialMessageState,
    reducers: {
        set(state, action: PayloadAction<{appId: number; pagedMessages: IPagedMessages;}>) {
            if (action.payload.appId === AllMessages) {
                state.items = action.payload.pagedMessages.messages;
            } else {
                const allMessages = state.items = [
                    ...state.items.filter(item => item.appid !== action.payload.appId),
                    ...action.payload.pagedMessages.messages,
                ];
                // keep the messages sorted from newest to oldest message
                state.items = allMessages.sort((a, b) => new Date(b.date).getTime() - new Date(a.date).getTime());
            }
            // TODO: paging functionality is missing (do not have a test case to test it with)
            // TODO: we maybe need to add here an additional state array object that is holding the information of more messages are existing
            // state.nextSince = action.payload.paging.since ?? 0;
            state.loaded = true;
        },
        add(state, action: PayloadAction<IMessage>) {
            state.items.unshift(action.payload);
            state.loaded = true;
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id !== action.payload);
            state.loaded = true;
        },
        removeByAppId(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.appid !== action.payload);
            state.loaded = true;
        },
        clear(state) {
            state.items = [];
            state.loaded = true;
        },
        loading(state, action: PayloadAction<boolean>) {
            state.loaded = !action.payload;
        },
    },
});

export const messageActions = messageSlice.actions;

export default messageSlice.reducer;
