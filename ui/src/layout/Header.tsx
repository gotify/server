import AppBar from '@mui/material/AppBar';
import Button, { ButtonProps } from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
import { Theme } from '@mui/material/styles';
import { makeStyles } from 'tss-react/mui';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import AccountCircle from '@mui/icons-material/AccountCircle';
import Chat from '@mui/icons-material/Chat';
import DevicesOther from '@mui/icons-material/DevicesOther';
import ExitToApp from '@mui/icons-material/ExitToApp';
import Brightness4 from '@mui/icons-material/Brightness4';
import Brightness7 from '@mui/icons-material/Brightness7';
import BrightnessAuto from '@mui/icons-material/BrightnessAuto';
import GitHubIcon from '@mui/icons-material/GitHub';
import MenuIcon from '@mui/icons-material/Menu';
import Apps from '@mui/icons-material/Apps';
import SupervisorAccount from '@mui/icons-material/SupervisorAccount';
import React, { CSSProperties } from 'react';
import { Link } from 'react-router';
import { useMediaQuery } from '@mui/material';
import Tooltip from '@mui/material/Tooltip';
import { ThemeKey } from './theme';

const themeIcons: Record<ThemeKey, React.ReactElement> = {
    dark: <Brightness4 />,
    light: <Brightness7 />,
    system: <BrightnessAuto />,
};

const useStyles = makeStyles()((theme: Theme) => ({
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        [theme.breakpoints.down('sm')]: {
            paddingBottom: 10,
        },
    },
    toolbar: {
        justifyContent: 'space-between',
    },
    menuButtons: {
        display: 'flex',
        [theme.breakpoints.down('md')]: {
            flex: 1,
        },
        justifyContent: 'center',
        [theme.breakpoints.down('sm')]: {
            flex: 'none',
            height: 'auto',
            margin: 0,
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
}));

interface IProps {
    loggedIn: boolean;
    name: string;
    admin: boolean;
    version: string;
    themeMode: ThemeKey;
    toggleTheme: VoidFunction;
    showSettings: VoidFunction;
    logout: VoidFunction;
    style: CSSProperties;
    navOpen: boolean;
    setNavOpen: (open: boolean) => void;
}

const Header = ({
    version,
    name,
    loggedIn,
    admin,
    toggleTheme,
    logout,
    style,
    navOpen,
    setNavOpen,
    showSettings,
    themeMode,
}: IProps) => {
    const { classes } = useStyles();
    const themeLabel = `Toggle theme (current: ${themeMode})`;
    const themeIcon = themeIcons[themeMode];
    return (
        <AppBar
            sx={{ position: { xs: 'sticky', sm: 'fixed' } }}
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
                        href={
                            version.startsWith('master-')
                                ? `https://github.com/gotify/server/commit/${version.replace('master-', '')}`
                                : `https://github.com/gotify/server/releases/tag/v${version}`
                        }
                        className={classes.link}>
                        <Typography variant="button" color="inherit">
                            @{version}
                        </Typography>
                    </a>
                </div>
                {loggedIn && (
                    <Buttons
                        admin={admin}
                        name={name}
                        logout={logout}
                        navOpen={navOpen}
                        setNavOpen={setNavOpen}
                        showSettings={showSettings}
                    />
                )}
                <div>
                    <Tooltip title={themeLabel} arrow>
                        <IconButton
                            onClick={toggleTheme}
                            color="inherit"
                            size="large"
                            aria-label={themeLabel}>
                            {themeIcon}
                        </IconButton>
                    </Tooltip>

                    <Tooltip title="Gotify on GitHub" arrow>
                        <a
                            href="https://github.com/gotify/server"
                            className={classes.link}
                            target="_blank"
                            rel="noopener noreferrer"
                            aria-label="Gotify on GitHub">
                            <IconButton color="inherit" size="large" aria-label="Gotify on GitHub">
                                <GitHubIcon />
                            </IconButton>
                        </a>
                    </Tooltip>
                </div>
            </Toolbar>
        </AppBar>
    );
};

const Buttons = ({
    showSettings,
    name,
    admin,
    logout,
    navOpen,
    setNavOpen,
}: {
    name: string;
    admin: boolean;
    logout: VoidFunction;
    navOpen: boolean;
    setNavOpen: (open: boolean) => void;
    showSettings: VoidFunction;
}) => {
    const { classes } = useStyles();
    const mobile = useMediaQuery('(max-width:600px)');

    return (
        <div className={classes.menuButtons}>
            <ResponsiveButton
                sx={{ display: { sm: 'none', xs: 'block' } }}
                icon={<MenuIcon />}
                onClick={() => setNavOpen(!navOpen)}
                label="menu"
                color="inherit"
            />
            {!mobile && (
                <>
                    {admin && (
                        <Link className={classes.link} to="/users" id="navigate-users">
                            <ResponsiveButton icon={<SupervisorAccount />} label="users" color="inherit" />
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
                </>
            )}
        </div>
    );
};

const ResponsiveButton: React.FC<{
    color: 'inherit';
    sx?: ButtonProps['sx'];
    label: string;
    id?: string;
    onClick?: () => void;
    icon: React.ReactNode;
}> = ({ icon, label, ...rest }) => {
    const matches = useMediaQuery('(max-width:1000px)');
    if (matches) {
        return (
            <Tooltip title={label} arrow>
                <IconButton {...rest} size="large" aria-label={label}>
                    {icon}
                </IconButton>
            </Tooltip>
        );
    }
    return (
        <Button startIcon={icon} {...rest}>
            {label}
        </Button>
    );
};

export default Header;
