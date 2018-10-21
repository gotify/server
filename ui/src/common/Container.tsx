import Paper from '@material-ui/core/Paper';
import {withStyles, WithStyles} from '@material-ui/core/styles';
import * as React from 'react';

const styles = () => ({
    paper: {
        padding: 16,
    },
});

interface IProps {
    style?: object;
}

const Container: React.SFC<IProps & WithStyles<'paper'>> = ({classes, children, style}) => {
    return (
        <Paper elevation={6} className={classes.paper} style={style}>
            {children}
        </Paper>
    );
};

export default withStyles(styles)<IProps>(Container);
