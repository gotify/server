import AppBar from '@mui/material/AppBar';
import Button, {ButtonProps} from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import {Theme} from '@mui/material/styles';
import {withStyles} from 'tss-react/mui';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import AccountCircle from '@mui/icons-material/AccountCircle';
import Chat from '@mui/icons-material/Chat';
import DevicesOther from '@mui/icons-material/DevicesOther';
import ExitToApp from '@mui/icons-material/ExitToApp';
import Highlight from '@mui/icons-material/Highlight';
import GitHubIcon from '@mui/icons-material/GitHub';
import MenuIcon from '@mui/icons-material/Menu';
import Apps from '@mui/icons-material/Apps';
import SupervisorAccount from '@mui/icons-material/SupervisorAccount';
import React, {Component, CSSProperties} from 'react';
import {Link} from 'react-router-dom';
import {observer} from 'mobx-react';
import {useMediaQuery} from '@mui/material';

const styles = (theme: Theme) =>
    ({
        appBar: {
            zIndex: theme.zIndex.drawer + 1,
            [theme.breakpoints.down('sm')]: {
                paddingBottom: 10,
            },
        },
        toolbar: {
            justifyContent: 'space-between',
            [theme.breakpoints.down('sm')]: {
                flexWrap: 'wrap',
            },
        },
        menuButtons: {
            display: 'flex',
            [theme.breakpoints.down('md')]: {
                flex: 1,
            },
            justifyContent: 'center',
            [theme.breakpoints.down('sm')]: {
                flexBasis: '100%',
                marginTop: 5,
                order: 1,
                height: 50,
                justifyContent: 'space-between',
                alignItems: 'center',
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
    } as const);

interface IProps {
    loggedIn: boolean;
    name: string;
    admin: boolean;
    version: string;
    classes?: Partial<Record<keyof ReturnType<typeof styles>, string>>;
    toggleTheme: VoidFunction;
    showSettings: VoidFunction;
    logout: VoidFunction;
    style: CSSProperties;
    setNavOpen: (open: boolean) => void;
}

@observer
class Header extends Component<IProps> {
    public render() {
        const {version, name, loggedIn, admin, toggleTheme, logout, style, setNavOpen} = this.props;

        const classes = withStyles.getClasses(this.props);

        return (
            <AppBar
                sx={{position: {xs: 'sticky', sm: 'fixed'}}}
                style={style}
                className={classes.appBar}>
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
                    {loggedIn && this.renderButtons(name, admin, logout, setNavOpen)}
                    <div>
                        <IconButton onClick={toggleTheme} color="inherit" size="large">
                            <Highlight />
                        </IconButton>

                        <a
                            href="https://github.com/gotify/server"
                            className={classes.link}
                            target="_blank"
                            rel="noopener noreferrer">
                            <IconButton color="inherit" size="large">
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
        setNavOpen: (open: boolean) => void
    ) {
        const classes = withStyles.getClasses(this.props);
        const {showSettings} = this.props;
        return (
            <div className={classes.menuButtons}>
                <ResponsiveButton
                    sx={{display: {sm: 'none', xs: 'block'}}}
                    icon={<MenuIcon />}
                    onClick={() => setNavOpen(true)}
                    label="menu"
                    color="inherit"
                />
                {admin && (
                    <Link className={classes.link} to="/users" id="navigate-users">
                        <ResponsiveButton
                            icon={<SupervisorAccount />}
                            label="users"
                            color="inherit"
                        />
                    </Link>
                )}
                <Link className={classes.link} to="/applications" id="navigate-apps">
                    <ResponsiveButton icon={<Chat />} label="apps" color="inherit" />
                </Link>
                <Link className={classes.link} to="/clients" id="navigate-clients">
                    <ResponsiveButton icon={<DevicesOther />} label="clients" color="inherit" />
                </Link>
                <Link className={classes.link} to="/plugins" id="navigate-plugins">
                    <ResponsiveButton icon={<Apps />} label="plugins" color="inherit" />
                </Link>
                <ResponsiveButton
                    icon={<AccountCircle />}
                    label={name}
                    onClick={showSettings}
                    id="changepw"
                    color="inherit"
                />
                <ResponsiveButton
                    icon={<ExitToApp />}
                    label="Logout"
                    onClick={logout}
                    id="logout"
                    color="inherit"
                />
            </div>
        );
    }
}

const ResponsiveButton: React.FC<{
    color: 'inherit';
    sx?: ButtonProps['sx'];
    label: string;
    id?: string;
    onClick?: () => void;
    icon: React.ReactNode;
}> = ({icon, label, ...rest}) => {
    const matches = useMediaQuery('(max-width:1000px)');
    if (matches) {
        return (
            <IconButton {...rest} size="large">
                {icon}
            </IconButton>
        );
    }
    return (
        <Button startIcon={icon} {...rest}>
            {label}
        </Button>
    );
};

export default withStyles(Header, styles);
