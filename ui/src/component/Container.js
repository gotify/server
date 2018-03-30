import React, {Component} from 'react';
import {withStyles} from 'material-ui/styles';
import Paper from 'material-ui/Paper';
import PropTypes from 'prop-types';

const styles = () => ({
    paper: {
        padding: 16,
    },
});

class Container extends Component {
    static propTypes = {
        classes: PropTypes.object.isRequired,
        children: PropTypes.node,
        style: PropTypes.object,
    };

    render() {
        const {classes, children, style} = this.props;
        return (
            <Paper elevation={6} className={classes.paper} style={style}>
                {children}
            </Paper>
        );
    }
}

export default withStyles(styles)(Container);
