import React, {useEffect, useState} from 'react';
import {Link} from 'react-router-dom';
import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import {makeStyles} from 'tss-react/mui';
import {mayAllowPermission, requestPermission} from '../snack/browserNotification';
import {
    Button,
    IconButton,
    Typography,
    ListItemButton,
    ListItemText,
    ListItemAvatar,
    Avatar,
    Box,
} from '@mui/material';
import {DrawerProps} from '@mui/material/Drawer/Drawer';
import CloseIcon from '@mui/icons-material/Close';
import {useAppDispatch, useAppSelector} from '../store';
import {fetchApps} from '../store/app-actions.ts';
import {appActions} from '../store/app-slice.ts';
import {uiActions} from '../store/ui-slice.ts';
import * as config from '../config';

const useStyles = makeStyles()(() => {
    return {
        root: {
            height: '100%',
        },
        drawerPaper: {
            position: 'relative',
            width: 250,
            minHeight: '100%',
            height: '100vh',
        },
        // toolbar: theme.mixins.toolbar,
        link: {
            color: 'inherit',
            textDecoration: 'none',
        },
    };
});

const Navigation = () => {
    const [showRequestNotification, setShowRequestNotification] = useState(mayAllowPermission());
    const apps = useAppSelector((state) => state.app.items);
    const {classes} = useStyles();
    const dispatch = useAppDispatch();
    const loggedIn = useAppSelector((state) => state.auth.loggedIn);

    useEffect(() => {
        if (loggedIn) {
            dispatch(fetchApps());
        }
    }, [dispatch, loggedIn]);

    const userApps =
        apps.length === 0
            ? null
            : apps.map((app) => (
                  <Link
                      onClick={() => {
                          dispatch(uiActions.setNavOpen(false));
                          // TODO: make sure we can select the app 'All Messages', does not work yet
                          dispatch(appActions.select(app));
                      }}
                      className={`${classes.link} item`}
                      to={'/messages/' + app.id}
                      key={app.id}>
                      <ListItemButton>
                          <ListItemAvatar style={{minWidth: 42}}>
                              <Avatar
                                  style={{width: 32, height: 32}}
                                  src={config.get('url') + app.image}
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

    // TODO: on mobile (size maybe md or xs) the app list on left should be hidden by default
    return (
        <ResponsiveDrawer
            classes={{root: classes.root, paper: classes.drawerPaper}}
            id="message-navigation">
            {/*<div className={classes.toolbar} />*/}
            <div className="toolbar" />
            <Link
                className={classes.link}
                to="/"
                onClick={() => {
                    dispatch(uiActions.setNavOpen(false));
                    dispatch(appActions.select(null));
                }}>
                <ListItemButton disabled={!loggedIn} className="all">
                    <ListItemText primary="All Messages" />
                </ListItemButton>
            </Link>
            <Divider />
            <div>{loggedIn ? userApps : placeholderItems}</div>
            <Divider />
            <Typography align="center" style={{marginTop: 10}}>
                {showRequestNotification ? (
                    <Button
                        onClick={() => {
                            requestPermission();
                            setShowRequestNotification(false);
                        }}>
                        Enable Notifications
                    </Button>
                ) : null}
            </Typography>
        </ResponsiveDrawer>
    );
};

const ResponsiveDrawer: React.FC<DrawerProps & {}> = ({children, ...rest}) => {
    const dispatch = useAppDispatch();
    const navOpen = useAppSelector((state) => state.ui.navOpen);

    return (
        <>
            <Box sx={{display: {xs: 'block', sm: 'none'}}}>
                <Drawer variant="temporary" open={navOpen} {...rest}>
                    <IconButton onClick={() => dispatch(uiActions.setNavOpen(false))}>
                        <CloseIcon />
                    </IconButton>
                    {children}
                </Drawer>
            </Box>
            <Box sx={{display: {xs: 'none', sm: 'block'}}}>
                <Drawer variant="permanent" {...rest}>
                    {children}
                </Drawer>
            </Box>
        </>
    );
};

export default Navigation;
