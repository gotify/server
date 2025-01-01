import {createSlice, PayloadAction} from '@reduxjs/toolkit';
import {IUser} from '../types.ts';

const initialAuthUserState = {
    name: 'unknown',
    admin: false,
    id: -1,
}

const initialAuthState = {
    loggedIn: false,
    authenticating: false,
    user: initialAuthUserState,
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
            state.loggedIn = false;
            state.authenticating = false;
            state.user = initialAuthUserState;
        },
        isAuthenticating: (state, action: PayloadAction<boolean>) => {
            state.authenticating = action.payload;
        }
    },
});

export const authActions = authSlice.actions;

export default authSlice.reducer;
