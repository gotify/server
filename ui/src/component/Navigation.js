import React, {Component} from 'react';
import Divider from 'material-ui/Divider';
import Drawer from 'material-ui/Drawer';
import {ListItem, ListItemText} from 'material-ui/List';
import {withStyles} from 'material-ui/styles';
import PropTypes from 'prop-types';
import AppStore from '../stores/AppStore';
import {Link} from 'react-router-dom';

const styles = (theme) => ({
    drawerPaper: {
        position: 'relative',
        width: 250,
        minHeight: '100%',
    },
    toolbar: theme.mixins.toolbar,
    link: {
        color: 'inherit',
        textDecoration: 'none',
    },
});

class Navigation extends Component {
    static propTypes = {
        classes: PropTypes.object.isRequired,
        loggedIn: PropTypes.bool.isRequired,
    };

    constructor() {
        super();
        this.state = {apps: []};
    }

    componentWillMount() {
        AppStore.on('change', this.updateApps);
    }

    componentWillUnmount() {
        AppStore.removeListener('change', this.updateApps);
    }

    updateApps = () => this.setState({apps: AppStore.get()});

    render() {
        const {classes, loggedIn} = this.props;
        const {apps} = this.state;

        const empty = (<ListItem disabled>
            <ListItemText primary="you have no applications :("/>
        </ListItem>);

        const userApps = apps.length === 0 ? empty : apps.map(function(app) {
            return (
                <Link className={classes.link} to={'/messages/' + app.id} key={app.id}>
                    <ListItem button>
                        <ListItemText primary={app.name}/>
                    </ListItem>
                </Link>
            );
        });

        const placeholderItems = [
            <ListItem button disabled key={-1}>
                <ListItemText primary="Some Server"/>
            </ListItem>,
            <ListItem button disabled key={-2}>
                <ListItemText primary="A Raspberry PI"/>
            </ListItem>,
        ];

        return (
            <Drawer variant="permanent" classes={{paper: classes.drawerPaper}}>
                <div className={classes.toolbar}/>
                <Link className={classes.link} to="/">
                    <ListItem button disabled={!loggedIn}>
                        <ListItemText primary="All Messages"/>
                    </ListItem>
                </Link>
                <Divider/>
                <div>{loggedIn ? userApps : placeholderItems}</div>
                <Divider/>
            </Drawer>
        );
    }
}

export default withStyles(styles, {withTheme: true})(Navigation);
