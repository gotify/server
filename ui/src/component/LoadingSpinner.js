import React, {Component} from 'react';
import {CircularProgress} from 'material-ui/Progress';
import DefaultPage from './DefaultPage';
import Grid from 'material-ui/Grid';

class LoadingSpinner extends Component {
    render() {
        return (
            <DefaultPage title="" maxWidth={250} hideButton={true}>
                <Grid item xs={12} style={{textAlign: 'center'}}>
                    <CircularProgress size={150}/>
                </Grid>
            </DefaultPage>
        );
    }
}

export default LoadingSpinner;
