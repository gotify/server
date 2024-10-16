import React from 'react';
import {Outlet} from 'react-router-dom';
import {createTheme, CssBaseline, Theme, ThemeProvider} from '@mui/material';
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {makeStyles} from 'tss-react/mui';
import LoadingSpinner from '../common/LoadingSpinner.tsx';
import ScrollUpButton from '../common/ScrollUpButton.tsx';
import SettingsDialog from '../common/SettingsDialog.tsx';
import * as config from '../config.ts';
import Header from '../layout/Header.tsx';
import Navigation from '../layout/Navigation.tsx';
import SnackBarHandler from '../snack/SnackBarHandler.tsx';

import {useAppDispatch, useAppSelector} from '../store';
import {ConnectionErrorBanner} from '../common/ConnectionErrorBanner.tsx';
import {ThemeKey, uiActions} from '../store/ui-slice.ts';

const useStyles = makeStyles()((theme) => {
    return {
        content: {
            margin: '0 auto',
            marginTop: 64,
            padding: theme.spacing(4),
            width: '100%',
            [theme.breakpoints.down('xs')]: {
                marginTop: 0,
            },
        },
    };
});

const themeMap: Record<ThemeKey, Theme> = {
    light: createTheme({
        palette: {
            mode: 'light',
        },
    }),
    dark: createTheme({
        palette: {
            mode: 'dark',
        },
    }),
};

const RootLayout = () => {
    const dispatch = useAppDispatch();
    const themeKey = useAppSelector((state) => state.ui.themeKey);
    const connectionErrorMessage = useAppSelector((state) => state.ui.connectionErrorMessage);
    const showSettings = useAppSelector((state) => state.ui.showSettings);
    const authenticating = useAppSelector((state) => state.auth.authenticating);
    const { classes } = useStyles();

    const theme = themeMap[themeKey];
    const { version } = config.get('version');

    return (
        <ThemeProvider theme={theme}>
            <div>
                {!connectionErrorMessage ? null : (
                    <ConnectionErrorBanner
                        height={64}
                        message={connectionErrorMessage}
                    />
                )}
                <div style={{display: 'flex', flexDirection: 'column'}}>
                    <CssBaseline />
                    <Header
                        style={{top: !connectionErrorMessage ? 0 : 64}}
                        version={version}
                    />
                    <div style={{display: 'flex'}}>
                        <Navigation />
                        <main className={classes.content}>
                            {authenticating ? <LoadingSpinner /> : <Outlet />}
                        </main>
                    </div>
                    {showSettings && (
                        <SettingsDialog fClose={() => dispatch(uiActions.setShowSettings(false))} />
                    )}
                    <ScrollUpButton />
                    <SnackBarHandler />
                </div>
            </div>
        </ThemeProvider>
    );
};

export default RootLayout;
