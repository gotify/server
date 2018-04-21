import {WithStyles} from "material-ui";
import Paper from 'material-ui/Paper';
import {withStyles} from 'material-ui/styles';
import * as React from 'react';

const styles = () => ({
    paper: {
        padding: 16,
    },
});

interface IProps {
    style?: object,
}

const Container: React.SFC<IProps & WithStyles<'paper'>> = ({classes, children, style}) => {
    return (
        <Paper elevation={6} className={classes.paper} style={style}>
            {children}
        </Paper>
    );
};

export default withStyles(styles)<IProps>(Container);
