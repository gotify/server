import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IUser} from '../types.ts';

const initialAuthState = {
    loggedIn: false,
    authenticating: false,
    user: {
        name: 'unknown',
        admin: false,
        id: -1,
    },
};

export const authSlice = createSlice({
    name: 'auth',
    initialState: initialAuthState,
    reducers: {
        login: (state, action: PayloadAction<IUser>) => {
            state.user = action.payload;
            state.loggedIn = true;
        },
        logout: (state) => {
            // TODO: does return undefined maybe work to reset the state?
            state.loggedIn = false;
            state.authenticating = false;
            state.user.name = 'unknown';
            state.user.admin = false;
            state.user.id = -1;
            // TODO: maybe we need to clear the complete store to not leak messages?
        },
        isAuthenticating: (state, action: PayloadAction<boolean>) => {
            state.authenticating = action.payload;
        }
    },
});

export const authActions = authSlice.actions;

export default authSlice.reducer;
