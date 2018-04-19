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

class Container extends React.Component<IProps & WithStyles<'paper'>, {}> {
    public render() {
        const {classes, children, style} = this.props;
        return (
            <Paper elevation={6} className={classes.paper} style={style}>
                {children}
            </Paper>
        );
    }
}

export default withStyles(styles)<IProps>(Container);
