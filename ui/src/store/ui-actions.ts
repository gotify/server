import {Middleware} from '@reduxjs/toolkit';
import {AppDispatch, RootState} from './index.ts';
import {localStorageThemeKey, ThemeKey, uiActions} from './ui-slice.ts';

export const toggleTheme = () => {
    return async (dispatch: AppDispatch, getState: () => RootState) => {
        const currentTheme = getState().ui.themeKey;
        const newTheme: ThemeKey = currentTheme === 'dark' ? 'light' : 'dark';

        dispatch(uiActions.setTheme(newTheme));
        localStorage.setItem(localStorageThemeKey, newTheme);
    }
}

const isThemeKey = (value: string | null): value is ThemeKey =>
    value === 'light' || value === 'dark';

export const loadStoredTheme = () => {
    return async (dispatch: AppDispatch, getState: () => RootState) => {

        const localStorageTheme = window.localStorage.getItem(localStorageThemeKey);
        if (isThemeKey(localStorageTheme)) {
            dispatch(uiActions.setTheme(localStorageTheme));
        } else {
            window.localStorage.setItem(localStorageThemeKey, getState().ui.themeKey);
        }

    }
}

// middleware that triggers a reload of messages as the websocket connection is affected by an unavailability of the server
export const connectionErrorMiddleware: Middleware<{}, RootState> = (store) => (next) => (action) => {
    const prevErrorMessage = store.getState().ui.connectionErrorMessage;
    const result = next(action);
    const nextErrorMessage = store.getState().ui.connectionErrorMessage;

    if (action.type.startsWith('ui/')) {
        if (prevErrorMessage !== null && nextErrorMessage === null) {
            store.dispatch(uiActions.addSnackMessage('Connection to server restored, re-fetching data...'));
            store.dispatch(uiActions.setReloadRequired(true));
        }
    }
    return result;
}
