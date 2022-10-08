import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import {createStyles, Theme, WithStyles, withStyles} from '@material-ui/core/styles';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import AccountCircle from '@material-ui/icons/AccountCircle';
import Chat from '@material-ui/icons/Chat';
import DevicesOther from '@material-ui/icons/DevicesOther';
import ExitToApp from '@material-ui/icons/ExitToApp';
import Highlight from '@material-ui/icons/Highlight';
import GitHubIcon from '@material-ui/icons/GitHub';
import MenuIcon from '@material-ui/icons/Menu';
import Apps from '@material-ui/icons/Apps';
import SupervisorAccount from '@material-ui/icons/SupervisorAccount';
import React, {Component, CSSProperties} from 'react';
import {Link} from 'react-router-dom';
import {observer} from 'mobx-react';
import {Hidden, PropTypes, withWidth} from '@material-ui/core';
import {Breakpoint} from '@material-ui/core/styles/createBreakpoints';

const styles = (theme: Theme) =>
    createStyles({
        appBar: {
            zIndex: theme.zIndex.drawer + 1,
            [theme.breakpoints.down('xs')]: {
                paddingBottom: 10,
            },
        },
        toolbar: {
            justifyContent: 'space-between',
            [theme.breakpoints.down('xs')]: {
                flexWrap: 'wrap',
            },
        },
        menuButtons: {
            display: 'flex',
            [theme.breakpoints.down('sm')]: {
                flex: 1,
            },
            justifyContent: 'center',
            [theme.breakpoints.down('xs')]: {
                flexBasis: '100%',
                marginTop: 5,
                order: 1,
                justifyContent: 'space-between',
            },
        },
        title: {
            [theme.breakpoints.up('md')]: {
                flex: 1,
            },
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

type Styles = WithStyles<'link' | 'menuButtons' | 'toolbar' | 'titleName' | 'title' | 'appBar'>;

interface IProps extends Styles {
    loggedIn: boolean;
    name: string;
    admin: boolean;
    version: string;
    toggleTheme: VoidFunction;
    showSettings: VoidFunction;
    logout: VoidFunction;
    style: CSSProperties;
    width: Breakpoint;
    setNavOpen: (open: boolean) => void;
}

@observer
class Header extends Component<IProps> {
    public render() {
        const {
            classes,
            version,
            name,
            loggedIn,
            admin,
            toggleTheme,
            logout,
            style,
            setNavOpen,
            width,
        } = this.props;

        const position = width === 'xs' ? 'sticky' : 'fixed';

        return (
            <AppBar position={position} style={style} className={classes.appBar}>
                <Toolbar className={classes.toolbar}>
                    <div className={classes.title}>
                        <Link to="/" className={classes.link}>
                            <Typography variant="h5" className={classes.titleName} color="inherit">
                                Gotify
                            </Typography>
                        </Link>
                        <a
                            href={'https://github.com/gotify/server/releases/tag/v' + version}
                            className={classes.link}>
                            <Typography variant="button" color="inherit">
                                @{version}
                            </Typography>
                        </a>
                    </div>
                    {loggedIn && this.renderButtons(name, admin, logout, width, setNavOpen)}
                    <div>
                        <IconButton onClick={toggleTheme} color="inherit">
                            <Highlight />
                        </IconButton>

                        <a
                            href="https://github.com/gotify/server"
                            className={classes.link}
                            target="_blank"
                            rel="noopener noreferrer">
                            <IconButton color="inherit">
                                <GitHubIcon />
                            </IconButton>
                        </a>
                    </div>
                </Toolbar>
            </AppBar>
        );
    }

    private renderButtons(
        name: string,
        admin: boolean,
        logout: VoidFunction,
        width: Breakpoint,
        setNavOpen: (open: boolean) => void
    ) {
        const {classes, showSettings} = this.props;
        return (
            <div className={classes.menuButtons}>
                <Hidden smUp implementation="css">
                    <ResponsiveButton
                        icon={<MenuIcon />}
                        onClick={() => setNavOpen(true)}
                        label="menu"
                        width={width}
                        color="inherit"
                    />
                </Hidden>
                {admin && (
                    <Link className={classes.link} to="/users" id="navigate-users">
                        <ResponsiveButton
                            icon={<SupervisorAccount />}
                            label="users"
                            width={width}
                            color="inherit"
                        />
                    </Link>
                )}
                <Link className={classes.link} to="/applications" id="navigate-apps">
                    <ResponsiveButton icon={<Chat />} label="apps" width={width} color="inherit" />
                </Link>
                <Link className={classes.link} to="/clients" id="navigate-clients">
                    <ResponsiveButton
                        icon={<DevicesOther />}
                        label="clients"
                        width={width}
                        color="inherit"
                    />
                </Link>
                <Link className={classes.link} to="/plugins" id="navigate-plugins">
                    <ResponsiveButton
                        icon={<Apps />}
                        label="plugins"
                        width={width}
                        color="inherit"
                    />
                </Link>
                <ResponsiveButton
                    icon={<AccountCircle />}
                    label={name}
                    onClick={showSettings}
                    id="changepw"
                    width={width}
                    color="inherit"
                />
                <ResponsiveButton
                    icon={<ExitToApp />}
                    label="Logout"
                    onClick={logout}
                    id="logout"
                    width={width}
                    color="inherit"
                />
            </div>
        );
    }
}

const ResponsiveButton: React.FC<{
    width: Breakpoint;
    color: PropTypes.Color;
    label: string;
    id?: string;
    onClick?: () => void;
    icon: React.ReactNode;
}> = ({width, icon, label, ...rest}) => {
    if (width === 'xs' || width === 'sm') {
        return <IconButton {...rest}>{icon}</IconButton>;
    }
    return (
        <Button startIcon={icon} {...rest}>
            {label}
        </Button>
    );
};

export default withWidth()(withStyles(styles, {withTheme: true})(Header));
