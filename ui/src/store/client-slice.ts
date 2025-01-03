import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IClient} from '../types.ts';

interface ClientState {
    items: IClient[];
}

const initialClientState: ClientState = {
    items: [],
}

export const clientSlice = createSlice({
    name: "client",
    initialState: initialClientState,
    reducers: {
        set(state, action: PayloadAction<IClient[]>) {
            state.items = action.payload;
        },
        add(state, action: PayloadAction<IClient>) {
            state.items.push(action.payload);
        },
        replace(state, action: PayloadAction<IClient>) {
            const itemIndex = state.items.findIndex((item) => item.id === action.payload.id);

            if (itemIndex !== -1) {
                state.items[itemIndex] = action.payload;
            }
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id !== action.payload);
        },
        clear(state) {
            state.items = [];
        },
    }
});

export const clientActions = clientSlice.actions;

export default clientSlice.reducer;
