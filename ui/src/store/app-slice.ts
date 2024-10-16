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
} = {
    items: [],
    selectedItem: initialSelectedItemState,
};

const appSlice = createSlice({
    name: 'apps',
    initialState: initialAppState,
    reducers: {
        set(state, action: PayloadAction<IApplication[]>) {
            state.items = action.payload;
        },
        add(state, action: PayloadAction<IApplication>) {
            state.items.push(action.payload);
        },
        replace(state, action: PayloadAction<IApplication>) {
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
        select(state, action: PayloadAction<IApplication | null>) {
            if (action.payload === null) {
                state.selectedItem = initialSelectedItemState;
            } else {
                state.selectedItem = action.payload;
            }
        }
    }
});

export const appActions = appSlice.actions;

export default appSlice.reducer;
