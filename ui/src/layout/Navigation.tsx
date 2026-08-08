import Divider from '@mui/material/Divider';
import Drawer, {DrawerProps} from '@mui/material/Drawer';
import {Theme} from '@mui/material/styles';
import React from 'react';
import {Link} from 'react-router';
import {observer} from 'mobx-react-lite';
import {mayAllowPermission, requestPermission} from '../snack/browserNotification';
import {
    Button,
    Typography,
    ListItemText,
    ListItemAvatar,
    Avatar,
    ListItemButton,
} from '@mui/material';
import AccountCircle from '@mui/icons-material/AccountCircle';
import Apps from '@mui/icons-material/Apps';
import Chat from '@mui/icons-material/Chat';
import DevicesOther from '@mui/icons-material/DevicesOther';
import ExitToApp from '@mui/icons-material/ExitToApp';
import SupervisorAccount from '@mui/icons-material/SupervisorAccount';
import {makeStyles} from 'tss-react/mui';
import {useStores} from '../stores';

const useStyles = makeStyles()((theme: Theme) => ({
    root: {
        height: '100%',
    },
    drawerPaper: {
        position: 'fixed',
        width: 250,
        height: '100vh',
        display: 'flex',
        flexDirection: 'column',
    },
    drawerContent: {
        display: 'flex',
        flexDirection: 'column',
        minHeight: 0,
        height: '100%',
    },
    scrollArea: {
        minHeight: 0,
        overflowY: 'auto',
        flex: 1,
        overscrollBehavior: 'contain',
        WebkitOverflowScrolling: 'touch',
    },
    mobileOnly: {
        [theme.breakpoints.up('sm')]: {
            display: 'none',
        },
    },
    // eslint-disable-next-line
    toolbar: theme.mixins.toolbar as any,
    mobileToolbar: {
        [theme.breakpoints.down('sm')]: {
            display: 'none',
        },
    },
    link: {
        color: 'inherit',
        textDecoration: 'none',
    },
}));

interface IProps {
    loggedIn: boolean;
    admin: boolean;
    name: string;
    logout: VoidFunction;
    showSettings: VoidFunction;
    bannerVisible: boolean;
    navOpen: boolean;
    setNavOpen: (open: boolean) => void;
}

const Navigation = observer(
    ({loggedIn, admin, name, logout, showSettings, bannerVisible, navOpen, setNavOpen}: IProps) => {
        const [showRequestNotification, setShowRequestNotification] =
            React.useState(mayAllowPermission);
        const {classes} = useStyles();
        const {appStore} = useStores();
        const apps = appStore.getItems();

        const userApps =
            apps.length === 0
                ? null
                : apps.map((app) => (
                      <Link
                          onClick={() => setNavOpen(false)}
                          className={`${classes.link} item`}
                          to={'/messages/' + app.id}
                          key={app.id}>
                          <ListItemButton>
                              <ListItemAvatar style={{minWidth: 42}}>
                                  <Avatar
                                      style={{width: 32, height: 32}}
                                      src={app.image}
                                      variant="square"
                                  />
                              </ListItemAvatar>
                              <ListItemText primary={app.name} />
                          </ListItemButton>
                      </Link>
                  ));

        const placeholderItems = [
            <ListItemButton disabled key={-1}>
                <ListItemText primary="Some Server" />
            </ListItemButton>,
            <ListItemButton disabled key={-2}>
                <ListItemText primary="A Raspberry PI" />
            </ListItemButton>,
        ];

        const mobileAction = (action: VoidFunction) => {
            action();
            setNavOpen(false);
        };

        return (
            <ResponsiveDrawer
                classes={{root: classes.root, paper: classes.drawerPaper}}
                navOpen={navOpen}
                setNavOpen={setNavOpen}
                bannerVisible={bannerVisible}
                id="message-navigation">
                <div className={classes.drawerContent}>
                    <div className={`${classes.toolbar} ${classes.mobileToolbar}`} />
                    <div className={classes.scrollArea}>
                        <Link className={classes.link} to="/" onClick={() => setNavOpen(false)}>
                            <ListItemButton disabled={!loggedIn} className="all">
                                <ListItemText primary="All Messages" />
                            </ListItemButton>
                        </Link>
                        <Divider />
                        {loggedIn && (
                            <div className={classes.mobileOnly}>
                                <Link
                                    className={classes.link}
                                    to="/applications"
                                    onClick={() => setNavOpen(false)}>
                                    <ListItemButton>
                                        <ListItemAvatar>
                                            <Chat />
                                        </ListItemAvatar>
                                        <ListItemText primary="Applications" />
                                    </ListItemButton>
                                </Link>
                                <Link
                                    className={classes.link}
                                    to="/clients"
                                    onClick={() => setNavOpen(false)}>
                                    <ListItemButton>
                                        <ListItemAvatar>
                                            <DevicesOther />
                                        </ListItemAvatar>
                                        <ListItemText primary="Clients" />
                                    </ListItemButton>
                                </Link>
                                <Link
                                    className={classes.link}
                                    to="/plugins"
                                    onClick={() => setNavOpen(false)}>
                                    <ListItemButton>
                                        <ListItemAvatar>
                                            <Apps />
                                        </ListItemAvatar>
                                        <ListItemText primary="Plugins" />
                                    </ListItemButton>
                                </Link>
                                {admin && (
                                    <Link
                                        className={classes.link}
                                        to="/users"
                                        onClick={() => setNavOpen(false)}>
                                        <ListItemButton>
                                            <ListItemAvatar>
                                                <SupervisorAccount />
                                            </ListItemAvatar>
                                            <ListItemText primary="Users" />
                                        </ListItemButton>
                                    </Link>
                                )}
                                <Divider />
                                <ListItemButton onClick={() => mobileAction(showSettings)}>
                                    <ListItemAvatar>
                                        <AccountCircle />
                                    </ListItemAvatar>
                                    <ListItemText primary={name} />
                                </ListItemButton>
                                <ListItemButton onClick={() => mobileAction(logout)}>
                                    <ListItemAvatar>
                                        <ExitToApp />
                                    </ListItemAvatar>
                                    <ListItemText primary="Logout" />
                                </ListItemButton>
                                <Divider />
                            </div>
                        )}
                        {loggedIn ? userApps : placeholderItems}
                    </div>
                    <Divider />
                    <Typography align="center" style={{marginTop: 10}}>
                        {showRequestNotification ? (
                            <Button
                                onClick={() => {
                                    requestPermission();
                                    setShowRequestNotification(false);
                                    setNavOpen(false);
                                }}>
                                Enable Notifications
                            </Button>
                        ) : null}
                    </Typography>
                </div>
            </ResponsiveDrawer>
        );
    }
);

const ResponsiveDrawer: React.FC<
    DrawerProps & {
        navOpen: boolean;
        setNavOpen: (open: boolean) => void;
        bannerVisible: boolean;
    }
> = ({navOpen, setNavOpen, bannerVisible, children, ...rest}) => (
    <>
        <Drawer
            sx={{
                display: {sm: 'none', xs: 'block'},
                '& .MuiDrawer-paper': {
                    top: bannerVisible ? 128 : 64,
                    height: bannerVisible ? 'calc(100vh - 128px)' : 'calc(100vh - 64px)',
                },
            }}
            variant="temporary"
            open={navOpen}
            onClose={() => setNavOpen(false)}
            {...rest}>
            {children}
        </Drawer>
        <Drawer
            sx={{
                display: {xs: 'none', sm: 'block'},
                '& .MuiDrawer-paper': {
                    top: bannerVisible ? 64 : 0,
                    height: bannerVisible ? 'calc(100vh - 64px)' : '100vh',
                },
            }}
            variant="permanent"
            {...rest}>
            {children}
        </Drawer>
    </>
);

export default Navigation;
