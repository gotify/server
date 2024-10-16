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


