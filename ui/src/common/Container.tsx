import Paper from '@mui/material/Paper';
import React from 'react';
import {makeStyles} from 'tss-react/mui';

const useStyles = makeStyles()(() => {
    return {
        paper: {
            padding: 16,
        },
    };
});

interface IProps {
    children: React.ReactNode;
    style?: React.CSSProperties;
}

const Container = ({ children, style}: IProps) => {
    const { classes } = useStyles();

    return(
        <Paper elevation={6} className={classes.paper} style={style}>
            {children}
        </Paper>
    );
}

export default Container;
