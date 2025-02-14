import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IApplication} from '../types.ts';

const initialSelectedItemState: IApplication = {
    id: -1,
    name: 'All Messages',
    image: '',
    internal: false,
    defaultPriority: 0,
    description: '',
    lastUsed: null,
    token: ''
}

const initialAppState: {
    items: IApplication[];
    selectedItem: IApplication;
    isLoading: boolean;
} = {
    items: [],
    selectedItem: initialSelectedItemState,
    isLoading: true,
};

const appSlice = createSlice({
    name: 'apps',
    initialState: initialAppState,
    reducers: {
        set(state, action: PayloadAction<IApplication[]>) {
            state.items = action.payload;
            state.isLoading = false;
        },
        add(state, action: PayloadAction<IApplication>) {
            state.items.push(action.payload);
            state.isLoading = false;
        },
        replace(state, action: PayloadAction<IApplication>) {
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
        select(state, action: PayloadAction<IApplication | null>) {
            if (action.payload === null) {
                state.selectedItem = initialSelectedItemState;
            } else {
                state.selectedItem = action.payload;
            }
            state.isLoading = false;
        },
        loading(state, action: PayloadAction<boolean>) {
            state.isLoading = action.payload;
        }
    }
});

export const appActions = appSlice.actions;

export default appSlice.reducer;
