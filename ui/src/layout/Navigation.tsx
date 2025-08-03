import Divider from '@mui/material/Divider';
import Drawer from '@mui/material/Drawer';
import {Theme} from '@mui/material/styles';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import {observer} from 'mobx-react';
import {inject, Stores} from '../inject';
import {mayAllowPermission, requestPermission} from '../snack/browserNotification';
import {
    Button,
    IconButton,
    Typography,
    ListItemText,
    ListItemAvatar,
    Avatar,
    ListItemButton,
} from '@mui/material';
import {DrawerProps} from '@mui/material/Drawer/Drawer';
import CloseIcon from '@mui/icons-material/Close';
import {withStyles} from 'tss-react/mui';

const styles = (theme: Theme) =>
    ({
        root: {
            height: '100%',
        },
        drawerPaper: {
            position: 'relative',
            width: 250,
            minHeight: '100%',
            height: '100vh',
        },
        // eslint-disable-next-line
        toolbar: theme.mixins.toolbar as any,
        link: {
            color: 'inherit',
            textDecoration: 'none',
        },
    } as const);

interface IProps {
    loggedIn: boolean;
    navOpen: boolean;
    classes?: Partial<Record<keyof ReturnType<typeof styles>, string>>;
    setNavOpen: (open: boolean) => void;
}

@observer
class Navigation extends Component<
    IProps & Stores<'appStore'>,
    {showRequestNotification: boolean}
> {
    public state = {showRequestNotification: mayAllowPermission()};

    public render() {
        const {loggedIn, appStore, navOpen, setNavOpen} = this.props;
        const classes = withStyles.getClasses(this.props);
        const {showRequestNotification} = this.state;
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

        return (
            <ResponsiveDrawer
                classes={{root: classes.root, paper: classes.drawerPaper}}
                navOpen={navOpen}
                setNavOpen={setNavOpen}
                id="message-navigation">
                <div className={classes.toolbar} />
                <Link className={classes.link} to="/" onClick={() => setNavOpen(false)}>
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
                                this.setState({showRequestNotification: false});
                            }}>
                            Enable Notifications
                        </Button>
                    ) : null}
                </Typography>
            </ResponsiveDrawer>
        );
    }
}

const ResponsiveDrawer: React.FC<
    DrawerProps & {navOpen: boolean; setNavOpen: (open: boolean) => void}
> = ({navOpen, setNavOpen, children, ...rest}) => (
    <>
        <Drawer
            sx={{display: {sm: 'none', xs: 'block'}}}
            variant="temporary"
            open={navOpen}
            {...rest}>
            <IconButton onClick={() => setNavOpen(false)} size="large">
                <CloseIcon />
            </IconButton>
            {children}
        </Drawer>
        <Drawer sx={{display: {xs: 'none', sm: 'block'}}} variant="permanent" {...rest}>
            {children}
        </Drawer>
    </>
);

export default withStyles(inject('appStore')(Navigation), styles);
