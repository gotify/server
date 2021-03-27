import Paper from '@material-ui/core/Paper';
import {withStyles, WithStyles} from '@material-ui/core/styles';
import * as React from 'react';

const styles = () => ({
    paper: {
        padding: 16,
    },
});

interface IProps extends WithStyles<'paper'> {
    style?: React.CSSProperties;
}

const Container: React.FC<IProps> = ({classes, children, style}) => (
    <Paper elevation={6} className={classes.paper} style={style}>
        {children}
    </Paper>
);

export default withStyles(styles)(Container);
