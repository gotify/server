import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import {Theme, WithStyles, withStyles} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Chat from '@material-ui/icons/Chat';
import DevicesOther from '@material-ui/icons/DevicesOther';
import ExitToApp from '@material-ui/icons/ExitToApp';
import Highlight from '@material-ui/icons/Highlight';
import Apps from '@material-ui/icons/Apps';
import SupervisorAccount from '@material-ui/icons/SupervisorAccount';
import React, {Component, CSSProperties} from 'react';
import {Link} from 'react-router-dom';
import {observer} from 'mobx-react';

const styles = (theme: Theme) => ({
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
    },
    title: {
        flex: 1,
        display: 'flex',
        alignItems: 'center',
    },
    titleName: {
        paddingRight: 10,
    },
    link: {
        color: 'inherit',
        textDecoration: 'none',
    },
});

type Styles = WithStyles<'link' | 'titleName' | 'title' | 'appBar'>;

interface IProps {
    loggedIn: boolean;
    name: string;
    admin: boolean;
    version: string;
    toggleTheme: VoidFunction;
    showSettings: VoidFunction;
    logout: VoidFunction;
    style: CSSProperties;
}

@observer
class Header extends Component<IProps & Styles> {
    public render() {
        const {classes, version, name, loggedIn, admin, toggleTheme, logout, style} = this.props;

        return (
            <AppBar position="absolute" style={style} className={classes.appBar}>
                <Toolbar>
                    <div className={classes.title}>
                        <a href="https://github.com/gotify/server" className={classes.link}>
                            <Typography
                                variant="headline"
                                className={classes.titleName}
                                color="inherit">
                                Gotify
                            </Typography>
                        </a>
                        <a
                            href={'https://github.com/gotify/server/releases/tag/v' + version}
                            className={classes.link}>
                            <Typography variant="button" color="inherit">
                                @{version}
                            </Typography>
                        </a>
                    </div>
                    {loggedIn && this.renderButtons(name, admin, logout)}
                    <IconButton onClick={toggleTheme} color="inherit">
                        <Highlight />
                    </IconButton>
                </Toolbar>
            </AppBar>
        );
    }

    private renderButtons(name: string, admin: boolean, logout: VoidFunction) {
        const {classes, showSettings} = this.props;
        return (
            <div>
                {admin ? (
                    <Link className={classes.link} to="/users" id="navigate-users">
                        <Button color="inherit">
                            <SupervisorAccount />
                            &nbsp;users
                        </Button>
                    </Link>
                ) : (
                    ''
                )}
                <Link className={classes.link} to="/applications" id="navigate-apps">
                    <Button color="inherit">
                        <Chat />
                        &nbsp;apps
                    </Button>
                </Link>
                <Link className={classes.link} to="/clients" id="navigate-clients">
                    <Button color="inherit">
                        <DevicesOther />
                        &nbsp;clients
                    </Button>
                </Link>
                <Link className={classes.link} to="/plugins" id="navigate-plugins">
                    <Button color="inherit">
                        <Apps />
                        &nbsp;plugins
                    </Button>
                </Link>
                <Button color="inherit" onClick={showSettings} id="changepw">
                    <AccountCircle />
                    &nbsp;
                    {name}
                </Button>
                <Button color="inherit" onClick={logout} id="logout">
                    <ExitToApp />
                    &nbsp;Logout
                </Button>
            </div>
        );
    }
}

export default withStyles(styles, {withTheme: true})<IProps>(Header);
