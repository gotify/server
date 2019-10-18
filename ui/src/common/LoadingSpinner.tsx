import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import DefaultPage from './DefaultPage';

export default function LoadingSpinner() {
    return (
        <DefaultPage title="" maxWidth={250}>
            <Grid item xs={12} style={{textAlign: 'center'}}>
                <CircularProgress size={150} />
            </Grid>
        </DefaultPage>
    );
}
