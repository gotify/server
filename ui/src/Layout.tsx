import axios, {AxiosResponse} from 'axios';
import CssBaseline from 'material-ui/CssBaseline';
import {createMuiTheme, MuiThemeProvider, Theme, WithStyles, withStyles} from 'material-ui/styles';
import * as React from 'react';
import {HashRouter, Redirect, Route, Switch} from 'react-router-dom';
import Header from './component/Header';
import LoadingSpinner from './component/LoadingSpinner';
import Navigation from './component/Navigation';
import ScrollUpButton from './component/ScrollUpButton';
import SettingsDialog from './component/SettingsDialog';
import SnackBarHandler from './component/SnackBarHandler';
import * as config from './config';
import Applications from './pages/Applications';
import Clients from './pages/Clients';
import Login from './pages/Login';
import Messages from './pages/Messages';
import Users from './pages/Users';
import GlobalStore from './stores/GlobalStore';

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

interface IState {
    darkTheme: boolean;
    redirect: boolean;
    showSettings: boolean;
    loggedIn: boolean;
    admin: boolean;
    name: string;
    authenticating: boolean;
    version: string;
}

class Layout extends React.Component<WithStyles<'content'>, IState> {
    private static defaultVersion = '0.0.0';

    public state = {
        admin: GlobalStore.isAdmin(),
        authenticating: GlobalStore.authenticating(),
        darkTheme: true,
        loggedIn: GlobalStore.isLoggedIn(),
        name: GlobalStore.getName(),
        redirect: false,
        showSettings: false,
        version: Layout.defaultVersion,
    };

    public componentDidMount() {
        if (this.state.version === Layout.defaultVersion) {
            axios.get(config.get('url') + 'version').then((resp: AxiosResponse<IVersion>) => {
                this.setState({...this.state, version: resp.data.version});
            });
        }
    }

    public componentWillMount() {
        GlobalStore.on('change', this.updateUser);
    }

    public componentWillUnmount() {
        GlobalStore.removeListener('change', this.updateUser);
    }

    public toggleTheme = () => this.setState({...this.state, darkTheme: !this.state.darkTheme});

    public updateUser = () => {
        this.setState({
            ...this.state,
            admin: GlobalStore.isAdmin(),
            authenticating: GlobalStore.authenticating(),
            loggedIn: GlobalStore.isLoggedIn(),
            name: GlobalStore.getName(),
        });
    };

    public hideSettings = () => this.setState({...this.state, showSettings: false});
    public showSettings = () => this.setState({...this.state, showSettings: true});

    public render() {
        const {name, admin, version, loggedIn, showSettings, authenticating} = this.state;
        const {classes} = this.props;
        const theme = this.state.darkTheme ? darkTheme : lightTheme;
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
                            toggleTheme={this.toggleTheme}
                            showSettings={this.showSettings}
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
                        {showSettings && <SettingsDialog fClose={this.hideSettings} />}
                        <ScrollUpButton />
                        <SnackBarHandler />
                    </div>
                </HashRouter>
            </MuiThemeProvider>
        );
    }
}

export default withStyles(styles, {withTheme: true})<{}>(Layout);
