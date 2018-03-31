import React, {Component} from 'react';
import Reboot from 'material-ui/Reboot';
import ScrollUpButton from './component/ScrollUpButton';
import Header from './component/Header';
import Navigation from './component/Navigation';
import Messages from './pages/Messages';
import Login from './pages/Login';
import axios from 'axios';
import {createMuiTheme, MuiThemeProvider, withStyles} from 'material-ui/styles';
import config from 'react-global-configuration';
import CurrentUserStore from './stores/CurrentUserStore';
import {HashRouter, Redirect, Route, Switch} from 'react-router-dom';
import {getToken} from './actions/defaultAxios';
import Applications from './pages/Applications';
import Clients from './pages/Clients';
import Users from './pages/Users';
import PropTypes from 'prop-types';
import SettingsDialog from './component/SettingsDialog';
import SnackBarHandler from './component/SnackBarHandler';

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

const styles = (theme) => ({
    content: {
        marginTop: 64,
        padding: theme.spacing.unit * 4,
        margin: '0 auto',
        width: '100%',
    },
});

class Layout extends Component {
    static propTypes = {
        classes: PropTypes.object.isRequired,
    };

    static defaultVersion = '0.0.0';

    state = {
        darkTheme: true,
        redirect: false,
        showSettings: false,
        loggedIn: CurrentUserStore.isLoggedIn(),
        admin: CurrentUserStore.isAdmin(),
        name: CurrentUserStore.getName(),
        version: Layout.defaultVersion,
    };

    componentDidMount() {
        if (this.state.version === Layout.defaultVersion) {
            axios.get(config.get('url') + 'version').then((resp) => {
                this.setState({...this.state, version: resp.data.version});
            });
        }
    }

    componentWillMount() {
        CurrentUserStore.on('change', this.updateUser);
    }

    componentWillUnmount() {
        CurrentUserStore.removeListener('change', this.updateUser);
    }

    toggleTheme = () => this.setState({...this.state, darkTheme: !this.state.darkTheme});

    updateUser = () => {
        this.setState({
            ...this.state,
            loggedIn: CurrentUserStore.isLoggedIn(),
            admin: CurrentUserStore.isAdmin(),
            name: CurrentUserStore.getName(),
        });
    };

    hideSettings = () => this.setState({...this.state, showSettings: false});
    showSettings = () => this.setState({...this.state, showSettings: true});

    render() {
        const {name, admin, version, loggedIn, showSettings} = this.state;
        const {classes} = this.props;
        const theme = this.state.darkTheme ? darkTheme : lightTheme;
        return (
            <MuiThemeProvider theme={theme}>
                <HashRouter>

                    <div style={{display: 'flex'}}>

                        <Reboot/>
                        <Header admin={admin} name={name} version={version} loggedIn={loggedIn}
                                toggleTheme={this.toggleTheme} showSettings={this.showSettings}/>
                        <Navigation loggedIn={loggedIn}/>

                        <main className={classes.content}>
                            <Switch>
                                <Route exact path="/login" render={() =>
                                    (loggedIn ? (<Redirect to="/"/>) : (<Login/>))}/>
                                {(loggedIn || getToken() != null) ? null : <Redirect to="/login"/>}
                                <Route exact path="/" component={Messages}/>
                                <Route exact path="/messages/:id" component={Messages}/>
                                <Route exact path="/applications" component={Applications}/>
                                <Route exact path="/clients" component={Clients}/>
                                <Route exact path="/users" component={Users}/>
                            </Switch>
                        </main>
                        {showSettings && <SettingsDialog fClose={this.hideSettings}/>}
                        <ScrollUpButton/>
                        <SnackBarHandler/>
                    </div>
                </HashRouter>

            </MuiThemeProvider>
        );
    }
}

export default withStyles(styles, {withTheme: true})(Layout);
