import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IUser} from '../types.ts';

interface UserState {
    items: IUser[];
}

const initialUserState: UserState = {
    items: [],
};

const userSlice = createSlice({
    name: 'users',
    initialState: initialUserState,
    reducers: {
        set(state, action: PayloadAction<IUser[]>) {
            state.items = action.payload;
        },
        add(state, action: PayloadAction<IUser>) {
            state.items.push(action.payload);
        },
        replace(state, action: PayloadAction<IUser>) {
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
    },
});

export const userActions = userSlice.actions;

export default userSlice.reducer;
