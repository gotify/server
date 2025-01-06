import AppBar from '@mui/material/AppBar';
import Button from '@mui/material/Button';
import IconButton from '@mui/material/IconButton';
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
import React, {CSSProperties} from 'react';
import {useNavigate} from 'react-router';
import {Link} from 'react-router-dom';
import {SxProps, useMediaQuery, useTheme} from '@mui/material';
import {makeStyles} from 'tss-react/mui';
import {useAppDispatch, useAppSelector} from '../store';
import {logout} from '../store/auth-actions.ts';
import {uiActions} from '../store/ui-slice.ts';
import {toggleTheme} from '../store/ui-actions.ts';

const useStyles = makeStyles()((theme) => {
    return {
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
    };
});

interface IProps {
    version: string;
    style: CSSProperties;
}

const Header = ({version, style}: IProps) => {
    const theme = useTheme();
    const {classes} = useStyles();
    const dispatch = useAppDispatch();

    const loggedIn = useAppSelector((state) => state.auth.loggedIn);

    return (
        <AppBar
            position={theme.breakpoints.down('sm') ? 'sticky' : 'fixed'}
            style={style}
            enableColorOnDark
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
                {loggedIn && <RenderButtons />}
                <div>
                    <IconButton onClick={() => dispatch(toggleTheme())} color="inherit">
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
};

const RenderButtons = () => {
    const dispatch = useAppDispatch();
    const {classes} = useStyles();
    const navigate = useNavigate();
    const admin = useAppSelector((state) => state.auth.user.admin);
    const name = useAppSelector((state) => state.auth.user.name);

    const handleLogout = async () => {
        await dispatch(logout());
        navigate('/login');
    };

    return (
        <div className={classes.menuButtons}>
            <ResponsiveButton
                sx={{display: {xl: 'none', xs: 'block'}}}
                icon={<MenuIcon />}
                onClick={() => dispatch(uiActions.setNavOpen(true))}
                label="menu"
                color="inherit"
            />
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
                onClick={() => dispatch(uiActions.setShowSettings(true))}
                id="changepw"
                color="inherit"
            />
            <ResponsiveButton
                icon={<ExitToApp />}
                label="Logout"
                id="logout"
                onClick={handleLogout}
                color="inherit"
            />
        </div>
    );
};

const ResponsiveButton: React.FC<{
    color: any;
    label: string;
    id?: string;
    onClick?: () => void;
    icon: React.ReactNode;
    sx?: SxProps;
}> = ({icon, label, ...rest}) => {
    const theme = useTheme();
    const smallerMd = useMediaQuery(theme.breakpoints.down('md'));

    if (smallerMd) {
        return <IconButton {...rest}>{icon}</IconButton>;
    }
    return (
        <Button startIcon={icon} {...rest}>
            {label}
        </Button>
    );
};

export default Header;
