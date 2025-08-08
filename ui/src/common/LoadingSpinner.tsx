import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import React from 'react';
import DefaultPage from './DefaultPage';

export default function LoadingSpinner() {
    return (
        <DefaultPage title="" maxWidth={250}>
            <Grid size={{xs: 12}} style={{textAlign: 'center'}}>
                <CircularProgress size={40} />
            </Grid>
        </DefaultPage>
    );
}
