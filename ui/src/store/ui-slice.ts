import {createSlice, PayloadAction} from '@reduxjs/toolkit';

interface UiState {
    themeKey: ThemeKey;
    connectionErrorMessage: string | null;
    navOpen: boolean;
    showSettings: boolean;
    snack: {
        messages: string[];
        message: string | null;
    }
}

const initialUiState: UiState = {
    themeKey: 'dark',
    connectionErrorMessage: null,
    navOpen: false,
    showSettings: false,
    snack: {
        messages: [],
        message: null,
    },
}
export type ThemeKey = 'dark' | 'light';

export const uiSlice = createSlice({
    name: 'ui',
    initialState: initialUiState,
    reducers: {
        setThemeKey: (state, action: PayloadAction<ThemeKey>) => {
            state.themeKey = action.payload;
        },
        setTheme: (state, action: PayloadAction<ThemeKey>) => {
            state.themeKey = action.payload;
        },
        setConnectionErrorMessage: (state, action: PayloadAction<string | null>) => {
            state.connectionErrorMessage = action.payload;
        },
        setNavOpen: (state, action: PayloadAction<boolean>) => {
            state.navOpen = action.payload;
        },
        setShowSettings: (state, action: PayloadAction<boolean>) => {
            state.showSettings = action.payload;
        },
        addSnackMessage: (state, action: PayloadAction<string>) => {
            state.snack.messages.push(action.payload);
        },
        nextSnackMessage: (state) => {
            if (state.snack.messages.length === 0) {
                throw new Error('There is no snack message');
            }
            // TODO check if this is working with state
            state.snack.message = state.snack.messages[0];
            state.snack.messages.splice(0, 1);
        }
    }
});

export const localStorageThemeKey = 'gotify-theme';

export const uiActions = uiSlice.actions;

export default uiSlice.reducer;
