import {createTheme, ThemeProvider, StyledEngineProvider, Theme} from '@mui/material';
import {makeStyles} from 'tss-react/mui';
import CssBaseline from '@mui/material/CssBaseline';
import * as React from 'react';
import {HashRouter, Navigate, Route, Routes} from 'react-router-dom';
import Header from './Header';
import Navigation from './Navigation';
import ScrollUpButton from '../common/ScrollUpButton';
import SettingsDialog from '../common/SettingsDialog';
import * as config from '../config';
import Applications from '../application/Applications';
import Clients from '../client/Clients';
import Plugins from '../plugin/Plugins';
import PluginDetailView from '../plugin/PluginDetailView';
import Login from '../user/Login';
import Messages from '../message/Messages';
import Users from '../user/Users';
import {observer} from 'mobx-react';
import {ConnectionErrorBanner} from '../common/ConnectionErrorBanner';
import {useStores} from '../stores';
import {SnackbarProvider} from 'notistack';

const useStyles = makeStyles()((theme: Theme) => ({
    content: {
        margin: '0 auto',
        marginTop: 64,
        padding: theme.spacing(4),
        width: '100%',
        [theme.breakpoints.down('sm')]: {
            marginTop: 0,
        },
    },
}));

const localStorageThemeKey = 'gotify-theme';
type ThemeKey = 'dark' | 'light';
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

const isThemeKey = (value: string | null): value is ThemeKey =>
    value === 'light' || value === 'dark';

const Layout = observer(() => {
    const {
        currentUser: {
            loggedIn,
            user: {name, admin},
            logout,
            tryReconnect,
            connectionErrorMessage,
        },
    } = useStores();
    const {classes} = useStyles();
    const [currentTheme, setCurrentTheme] = React.useState<ThemeKey>(() => {
        const stored = window.localStorage.getItem(localStorageThemeKey);
        return isThemeKey(stored) ? stored : 'dark';
    });
    const theme = themeMap[currentTheme];
    const {version} = config.get('version');
    const [navOpen, setNavOpen] = React.useState(false);
    const [showSettings, setShowSettings] = React.useState(false);

    const toggleTheme = () => {
        const next = currentTheme === 'dark' ? 'light' : 'dark';
        setCurrentTheme(next);
        localStorage.setItem(localStorageThemeKey, next);
    };

    const authed = (children: React.ReactNode) => (
        <RequireAuth loggedIn={loggedIn}>{children}</RequireAuth>
    );

    return (
        <StyledEngineProvider injectFirst>
            <ThemeProvider theme={theme}>
                <HashRouter>
                    <div>
                        {!connectionErrorMessage ? null : (
                            <ConnectionErrorBanner
                                height={64}
                                retry={() => tryReconnect()}
                                message={connectionErrorMessage}
                            />
                        )}
                        <div style={{display: 'flex', flexDirection: 'column'}}>
                            <CssBaseline />
                            <Header
                                admin={admin}
                                name={name}
                                style={{top: !connectionErrorMessage ? 0 : 64}}
                                version={version}
                                loggedIn={loggedIn}
                                toggleTheme={toggleTheme}
                                showSettings={() => setShowSettings(true)}
                                logout={logout}
                                setNavOpen={setNavOpen}
                            />
                            <div style={{display: 'flex'}}>
                                <Navigation
                                    loggedIn={loggedIn}
                                    navOpen={navOpen}
                                    setNavOpen={setNavOpen}
                                />
                                <main className={classes.content}>
                                    <Routes>
                                        <Route path="/login" element={<Login />} />
                                        <Route path="/" element={authed(<Messages />)} />
                                        <Route
                                            path="/messages/:id"
                                            element={authed(<Messages />)}
                                        />
                                        <Route
                                            path="/applications"
                                            element={authed(<Applications />)}
                                        />
                                        <Route path="/clients" element={authed(<Clients />)} />
                                        <Route path="/users" element={authed(<Users />)} />
                                        <Route path="/plugins" element={authed(<Plugins />)} />
                                        <Route
                                            path="/plugins/:id"
                                            element={authed(<PluginDetailView />)}
                                        />
                                    </Routes>
                                </main>
                            </div>
                            {showSettings && (
                                <SettingsDialog fClose={() => setShowSettings(false)} />
                            )}
                            <ScrollUpButton />
                            <SnackbarProvider />
                        </div>
                    </div>
                </HashRouter>
            </ThemeProvider>
        </StyledEngineProvider>
    );
});

const RequireAuth: React.FC<React.PropsWithChildren<{loggedIn: boolean}>> = ({
    children,
    loggedIn,
}) => {
    return loggedIn ? <>{children}</> : <Navigate replace={true} to="/login" />;
};

export default Layout;
