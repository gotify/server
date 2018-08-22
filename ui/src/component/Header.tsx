import {Theme, WithStyles} from 'material-ui';
import AccountCircle from 'material-ui-icons/AccountCircle';
import Chat from 'material-ui-icons/Chat';
import DevicesOther from 'material-ui-icons/DevicesOther';
import ExitToApp from 'material-ui-icons/ExitToApp';
import LightbulbOutline from 'material-ui-icons/LightbulbOutline';
import SupervisorAccount from 'material-ui-icons/SupervisorAccount';
import AppBar from 'material-ui/AppBar';
import Button from 'material-ui/Button';
import IconButton from 'material-ui/IconButton';
import {withStyles} from 'material-ui/styles';
import Toolbar from 'material-ui/Toolbar';
import Typography from 'material-ui/Typography';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import * as UserAction from '../actions/UserAction';

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
}

class Header extends Component<IProps & Styles> {
    public render() {
        const {classes, version, name, loggedIn, admin, toggleTheme} = this.props;

        return (
            <AppBar position="absolute" className={classes.appBar}>
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
                    {loggedIn && this.renderButtons(name, admin)}
                    <IconButton onClick={toggleTheme} color="inherit">
                        <LightbulbOutline />
                    </IconButton>
                </Toolbar>
            </AppBar>
        );
    }

    private renderButtons(name: string, admin: boolean) {
        const {classes, showSettings} = this.props;
        return (
            <div>
                {admin ? (
                    <Link className={classes.link} to="/users">
                        <Button color="inherit">
                            <SupervisorAccount />
                            &nbsp;users
                        </Button>
                    </Link>
                ) : (
                    ''
                )}
                <Link className={classes.link} to="/applications">
                    <Button color="inherit">
                        <Chat />
                        &nbsp;apps
                    </Button>
                </Link>
                <Link className={classes.link} to="/clients">
                    <Button color="inherit">
                        <DevicesOther />
                        &nbsp;clients
                    </Button>
                </Link>
                <Button color="inherit" onClick={showSettings}>
                    <AccountCircle />
                    &nbsp;
                    {name}
                </Button>
                <Button color="inherit" onClick={UserAction.logout}>
                    <ExitToApp />
                    &nbsp;Logout
                </Button>
            </div>
        );
    }
}

export default withStyles(styles, {withTheme: true})<IProps>(Header);
