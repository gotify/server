import Grid from '@material-ui/core/Grid';
import Typography from '@material-ui/core/Typography';
import React, {SFC} from 'react';

interface IProps {
    title: string;
    rightControl?: React.ReactNode;
    maxWidth?: number;
}

const DefaultPage: SFC<IProps> = ({title, rightControl, maxWidth = 700, children}) => (
    <main style={{margin: '0 auto', maxWidth}}>
        <Grid container spacing={4}>
            <Grid item xs={12} style={{display: 'flex'}}>
                <Typography variant="h4" style={{flex: 1}}>
                    {title}
                </Typography>
                {rightControl}
            </Grid>
            {children}
        </Grid>
    </main>
);
export default DefaultPage;
