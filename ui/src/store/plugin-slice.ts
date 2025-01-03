import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IPlugin} from '../types.ts';

interface PluginState {
    items: IPlugin[];
    displayText: string | null;
    currentConfig: string | null;
}

const initialPluginState: PluginState = {
    items: [],
    displayText: null,
    currentConfig: null,
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
            state.displayText = null;
            state.currentConfig = null;
        },
        setDisplay(state, action: PayloadAction<string | null>) {
            state.displayText = action.payload;
        },
        setCurrentConfig(state, action: PayloadAction<string | null>) {
            state.currentConfig = action.payload;
        }
    },
});

export const pluginActions = pluginSlice.actions;

export default pluginSlice.reducer;
