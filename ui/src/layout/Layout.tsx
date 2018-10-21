import {createMuiTheme, MuiThemeProvider, Theme, WithStyles, withStyles} from '@material-ui/core';
import CssBaseline from '@material-ui/core/CssBaseline';
import axios, {AxiosResponse} from 'axios';
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
import Login from '../pages/Login';
import Messages from '../message/Messages';
import Users from '../user/Users';
import {observer} from 'mobx-react';
import {observable} from 'mobx';
import {inject, Stores} from '../inject';

const lightTheme = createMuiTheme({
    palette: {
        type: 'light',
    },
});
const darkTheme = createMuiTheme({
    palette: {
        type: 'dark',
    },
});

const styles = (theme: Theme) => ({
    content: {
        margin: '0 auto',
        marginTop: 64,
        padding: theme.spacing.unit * 4,
        width: '100%',
    },
});

@observer
class Layout extends React.Component<WithStyles<'content'> & Stores<'currentUser'>> {
    private static defaultVersion = '0.0.0';

    @observable
    private darkThemeVisible = true;
    @observable
    private showSettings = false;
    @observable
    private version = Layout.defaultVersion;

    public componentDidMount() {
        if (this.version === Layout.defaultVersion) {
            axios.get(config.get('url') + 'version').then((resp: AxiosResponse<IVersion>) => {
                this.setState({...this.state, version: resp.data.version});
            });
        }
    }

    public render() {
        const {version, showSettings, darkThemeVisible} = this;
        const {
            classes,
            currentUser: {
                loggedIn,
                authenticating,
                user: {name, admin},
                logout,
            },
        } = this.props;
        const theme = darkThemeVisible ? darkTheme : lightTheme;
        const loginRoute = () => (loggedIn ? <Redirect to="/" /> : <Login />);
        return (
            <MuiThemeProvider theme={theme}>
                <HashRouter>
                    <div style={{display: 'flex'}}>
                        <CssBaseline />
                        <Header
                            admin={admin}
                            name={name}
                            version={version}
                            loggedIn={loggedIn}
                            toggleTheme={() => (this.darkThemeVisible = !this.darkThemeVisible)}
                            showSettings={() => (this.showSettings = true)}
                            logout={logout}
                        />
                        <Navigation loggedIn={loggedIn} />

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
                                <Route exact path="/messages/:id" component={Messages} />
                                <Route exact path="/applications" component={Applications} />
                                <Route exact path="/clients" component={Clients} />
                                <Route exact path="/users" component={Users} />
                            </Switch>
                        </main>
                        {showSettings && (
                            <SettingsDialog fClose={() => (this.showSettings = false)} />
                        )}
                        <ScrollUpButton />
                        <SnackBarHandler />
                    </div>
                </HashRouter>
            </MuiThemeProvider>
        );
    }
}

export default withStyles(styles, {withTheme: true})<{}>(inject('currentUser')(Layout));
