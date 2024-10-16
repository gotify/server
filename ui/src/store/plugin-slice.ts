import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IPlugin} from '../types.ts';

interface PluginState {
    items: IPlugin[];
}

const initialPluginState: PluginState = {
    items: [],
}

const pluginSlice = createSlice({
    name: 'plugins',
    initialState: initialPluginState,
    reducers: {
        set(state, action: PayloadAction<IPlugin[]>) {
            state.items = action.payload;
        },
        add(state, action: PayloadAction<IPlugin>) {
            state.items.push(action.payload);
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id === action.payload);
        },
        clear(state) {
            state.items = [];
        },
    },
});

export const pluginActions = pluginSlice.actions;

export default pluginSlice.reducer;
