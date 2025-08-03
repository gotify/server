import {createTheme, ThemeProvider, StyledEngineProvider, Theme} from '@mui/material';
import {withStyles} from 'tss-react/mui';
import CssBaseline from '@mui/material/CssBaseline';
import * as React from 'react';
import {HashRouter, Redirect, Route, Switch} from 'react-router-dom';
import Header from './Header';
import LoadingSpinner from '../common/LoadingSpinner';
import Navigation from './Navigation';
import ScrollUpButton from '../common/ScrollUpButton';
import SettingsDialog from '../common/SettingsDialog';
import SnackBarHandler from '../snack/SnackBarHandler';
import * as config from '../config';
import Applications from '../application/Applications';
import Clients from '../client/Clients';
import Plugins from '../plugin/Plugins';
import PluginDetailView from '../plugin/PluginDetailView';
import Login from '../user/Login';
import Messages from '../message/Messages';
import Users from '../user/Users';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';
import {ConnectionErrorBanner} from '../common/ConnectionErrorBanner';

const styles = (theme: Theme) => ({
    content: {
        margin: '0 auto',
        marginTop: 64,
        padding: theme.spacing(4),
        width: '100%',
        [theme.breakpoints.down('sm')]: {
            marginTop: 0,
        },
    },
});

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

interface LayoutProps {
    classes?: Partial<Record<keyof ReturnType<typeof styles>, string>>;
}

@observer
class Layout extends React.Component<LayoutProps & Stores<'currentUser' | 'snackManager'>> {
    @observable
    private currentTheme: ThemeKey = 'dark';
    @observable
    private showSettings = false;
    @observable
    private navOpen = false;

    private setNavOpen(open: boolean) {
        this.navOpen = open;
    }

    public componentDidMount() {
        const localStorageTheme = window.localStorage.getItem(localStorageThemeKey);
        if (isThemeKey(localStorageTheme)) {
            this.currentTheme = localStorageTheme;
        } else {
            window.localStorage.setItem(localStorageThemeKey, this.currentTheme);
        }
    }

    public render() {
        const {showSettings, currentTheme} = this;
        const {
            currentUser: {
                loggedIn,
                authenticating,
                user: {name, admin},
                logout,
                tryReconnect,
                connectionErrorMessage,
            },
        } = this.props;
        const classes = withStyles.getClasses(this.props);
        const theme = themeMap[currentTheme];
        const loginRoute = () => (loggedIn ? <Redirect to="/" /> : <Login />);
        const {version} = config.get('version');
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
                                    style={{top: !connectionErrorMessage ? 0 : 64}}
                                    admin={admin}
                                    name={name}
                                    version={version}
                                    loggedIn={loggedIn}
                                    toggleTheme={this.toggleTheme.bind(this)}
                                    showSettings={() => (this.showSettings = true)}
                                    logout={logout}
                                    setNavOpen={this.setNavOpen.bind(this)}
                                />
                                <div style={{display: 'flex'}}>
                                    <Navigation
                                        loggedIn={loggedIn}
                                        navOpen={this.navOpen}
                                        setNavOpen={this.setNavOpen.bind(this)}
                                    />
                                    <main className={classes.content}>
                                        <Switch>
                                            {authenticating ? (
                                                <Route path="/">
                                                    <LoadingSpinner />
                                                </Route>
                                            ) : null}
                                            <Route exact path="/login" render={loginRoute} />
                                            {loggedIn ? null : <Redirect to="/login" />}
                                            <Route exact path="/" component={Messages} />
                                            <Route
                                                exact
                                                path="/messages/:id"
                                                component={Messages}
                                            />
                                            <Route
                                                exact
                                                path="/applications"
                                                component={Applications}
                                            />
                                            <Route exact path="/clients" component={Clients} />
                                            <Route exact path="/users" component={Users} />
                                            <Route exact path="/plugins" component={Plugins} />
                                            <Route
                                                exact
                                                path="/plugins/:id"
                                                component={PluginDetailView}
                                            />
                                        </Switch>
                                    </main>
                                </div>
                                {showSettings && (
                                    <SettingsDialog fClose={() => (this.showSettings = false)} />
                                )}
                                <ScrollUpButton />
                                <SnackBarHandler />
                            </div>
                        </div>
                    </HashRouter>
                </ThemeProvider>
            </StyledEngineProvider>
        );
    }

    private toggleTheme() {
        this.currentTheme = this.currentTheme === 'dark' ? 'light' : 'dark';
        localStorage.setItem(localStorageThemeKey, this.currentTheme);
    }
}

export default withStyles(inject('currentUser', 'snackManager')(Layout), styles);
