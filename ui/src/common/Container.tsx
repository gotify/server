import Paper from '@mui/material/Paper';
import {withStyles} from 'tss-react/mui';
import * as React from 'react';

const styles = () =>
    ({
        paper: {
            padding: 16,
        },
    } as const);

interface IProps {
    style?: React.CSSProperties;
    classes?: Partial<Record<keyof ReturnType<typeof styles>, string>>;
}

const Container: React.FC<IProps> = ({children, style, ...props}) => {
    const classes = withStyles.getClasses(props);
    return (
        <Paper elevation={6} className={classes.paper} style={style}>
            {children}
        </Paper>
    );
};

export default withStyles(Container, styles);
