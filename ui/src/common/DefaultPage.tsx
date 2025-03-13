import Grid from '@mui/material/Grid2';
import Typography from '@mui/material/Typography';
import React from 'react';

interface IProps {
    children: React.ReactNode;
    title: string;
    rightControl?: React.ReactNode;
    maxWidth?: number;
}

const DefaultPage = ({title, rightControl, maxWidth = 700, children}: IProps) => (
    <main style={{margin: '0 auto', maxWidth}}>
        <Grid container spacing={4}>
            <Grid size={{ xs: 12 }} style={{display: 'flex', flexWrap: 'wrap'}}>
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
