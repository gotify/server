import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import {Theme, WithStyles} from '@material-ui/core/styles';
import {withStyles} from '@material-ui/core/styles';
import React, {Component} from 'react';
import {Link} from 'react-router-dom';
import AppStore from '../stores/AppStore';

const styles = (theme: Theme) => ({
    drawerPaper: {
        position: 'relative' as 'relative',
        width: 250,
        minHeight: '100%',
        height: '100vh',
    },
    toolbar: theme.mixins.toolbar as any,
    link: {
        color: 'inherit',
        textDecoration: 'none',
    },
});

type Styles = WithStyles<'drawerPaper' | 'toolbar' | 'link'>;

interface IProps {
    loggedIn: boolean;
}

interface IState {
    apps: IApplication[];
}

class Navigation extends Component<IProps & Styles, IState> {
    public state: IState = {apps: []};

    public componentWillMount() {
        AppStore.on('change', this.updateApps);
    }

    public componentWillUnmount() {
        AppStore.removeListener('change', this.updateApps);
    }

    public render() {
        const {classes, loggedIn} = this.props;
        const {apps} = this.state;

        const userApps =
            apps.length === 0
                ? null
                : apps.map((app) => {
                      return (
                          <Link className={classes.link} to={'/messages/' + app.id} key={app.id}>
                              <ListItem button>
                                  <ListItemText primary={app.name} />
                              </ListItem>
                          </Link>
                      );
                  });

        const placeholderItems = [
            <ListItem button disabled key={-1}>
                <ListItemText primary="Some Server" />
            </ListItem>,
            <ListItem button disabled key={-2}>
                <ListItemText primary="A Raspberry PI" />
            </ListItem>,
        ];

        return (
            <Drawer variant="permanent" classes={{paper: classes.drawerPaper}}>
                <div className={classes.toolbar} />
                <Link className={classes.link} to="/">
                    <ListItem button disabled={!loggedIn}>
                        <ListItemText primary="All Messages" />
                    </ListItem>
                </Link>
                <Divider />
                <div>{loggedIn ? userApps : placeholderItems}</div>
                <Divider />
            </Drawer>
        );
    }

    private updateApps = () => this.setState({apps: AppStore.get()});
}

export default withStyles(styles, {withTheme: true})<IProps>(Navigation);
