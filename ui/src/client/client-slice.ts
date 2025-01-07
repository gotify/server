import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IClient} from '../types.ts';

interface ClientState {
    items: IClient[];
    isLoading: boolean;
}

const initialClientState: ClientState = {
    items: [],
    isLoading: true,
}

export const clientSlice = createSlice({
    name: "client",
    initialState: initialClientState,
    reducers: {
        set(state, action: PayloadAction<IClient[]>) {
            state.items = action.payload;
            state.isLoading = false;
        },
        add(state, action: PayloadAction<IClient>) {
            state.items.push(action.payload);
            state.isLoading = false;
        },
        replace(state, action: PayloadAction<IClient>) {
            const itemIndex = state.items.findIndex((item) => item.id === action.payload.id);

            if (itemIndex !== -1) {
                state.items[itemIndex] = action.payload;
            }
            state.isLoading = false;
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id !== action.payload);
            state.isLoading = false;
        },
        clear(state) {
            state.items = [];
            state.isLoading = false;
        },
        loading(state, action: PayloadAction<boolean>) {
            state.isLoading = action.payload;
        }
    }
});

export const clientActions = clientSlice.actions;

export default clientSlice.reducer;
