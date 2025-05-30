import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IPlugin} from '../types.ts';

interface PluginState {
    items: IPlugin[];
    isLoading: boolean;
    displayText: string | null;
    currentConfig: string | null;
}

const initialPluginState: PluginState = {
    items: [],
    isLoading: true,
    displayText: null,
    currentConfig: null,
}

const pluginSlice = createSlice({
    name: 'plugins',
    initialState: initialPluginState,
    reducers: {
        set(state, action: PayloadAction<IPlugin[]>) {
            state.items = action.payload;
            state.isLoading = false;
        },
        add(state, action: PayloadAction<IPlugin>) {
            state.items.push(action.payload);
            state.isLoading = false;
        },
        remove(state, action: PayloadAction<number>) {
            state.items = state.items.filter((item) => item.id === action.payload);
            state.isLoading = false;
        },
        clear(state) {
            state.items = [];
            state.displayText = null;
            state.currentConfig = null;
            state.isLoading = false;
        },
        setDisplay(state, action: PayloadAction<string | null>) {
            state.displayText = action.payload;
            state.isLoading = false;
        },
        setCurrentConfig(state, action: PayloadAction<string | null>) {
            state.currentConfig = action.payload;
            state.isLoading = false;
        },
        loading(state, action: PayloadAction<boolean>) {
            state.isLoading = action.payload;
}
    },
});

export const pluginActions = pluginSlice.actions;

export default pluginSlice.reducer;
