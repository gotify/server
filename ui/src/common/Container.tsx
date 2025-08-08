import Paper from '@mui/material/Paper';
import {makeStyles} from 'tss-react/mui';
import * as React from 'react';

const useStyles = makeStyles()(() => ({
    paper: {
        padding: 16,
    },
}));

interface IProps {
    style?: React.CSSProperties;
}

const Container: React.FC<React.PropsWithChildren<IProps>> = ({children, style}) => {
    const {classes} = useStyles();
    return (
        <Paper elevation={6} className={classes.paper} style={style}>
            {children}
        </Paper>
    );
};

export default Container;
