import Grid from 'material-ui/Grid';
import {CircularProgress} from 'material-ui/Progress';
import React from 'react';
import DefaultPage from './DefaultPage';

export default function LoadingSpinner() {
    return (
        <DefaultPage title="" maxWidth={250} hideButton={true}>
            <Grid item xs={12} style={{textAlign: 'center'}}>
                <CircularProgress size={150} />
            </Grid>
        </DefaultPage>
    );
}
