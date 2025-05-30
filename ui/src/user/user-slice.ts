import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IUser} from '../types.ts';

interface UserState {
    items: IUser[];
    isLoading: boolean;
}

const initialUserState: UserState = {
    items: [],
    isLoading: true,
};

const userSlice = createSlice({
    name: 'users',
    initialState: initialUserState,
    reducers: {
        set(state, action: PayloadAction<IUser[]>) {
            state.items = action.payload;
            state.isLoading = false;
        },
        add(state, action: PayloadAction<IUser>) {
            state.items.push(action.payload);
            state.isLoading = false;
        },
        replace(state, action: PayloadAction<IUser>) {
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
    },
});

export const userActions = userSlice.actions;

export default userSlice.reducer;
