import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IMessage, IPagedMessages} from '../types.ts';

interface MessagesState {
    items: IMessage[];
    hasMore: boolean;
    nextSince: number;
    loaded: boolean;
    isLoading: boolean;
}

const initialMessageState: MessagesState = {
    items: [],
    hasMore: false,
    nextSince: 0,
    loaded: false,
    isLoading: false,
}

export const messageSlice = createSlice({
    name: 'message',
    initialState: initialMessageState,
    reducers: {
        set(state, action: PayloadAction<IPagedMessages>) {
            state.items = action.payload.messages;
            state.nextSince = action.payload.paging.since ?? 0;
            state.loaded = true;
            // state.items = action.payload;
        },
        add(state, action: PayloadAction<IMessage>) {
            state.items.unshift(action.payload);
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id !== action.payload);
        },
        removeByAppId(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.appid !== action.payload);
        },
        clear(state) {
            state.items = [];
        }
    }
});

export const messageActions = messageSlice.actions;

export default messageSlice.reducer;
