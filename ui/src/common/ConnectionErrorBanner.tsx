import React from 'react';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';

interface ConnectionErrorBannerProps {
    height: number;
    retry: () => void;
    message: string;
}

export const ConnectionErrorBanner = ({height, retry, message}: ConnectionErrorBannerProps) => (
    <div
        style={{
            backgroundColor: '#e74c3c',
            height,
            width: '100%',
            zIndex: 1300,
            position: 'relative',
        }}>
        <Typography align="center" variant="h6" style={{lineHeight: `${height}px`}}>
            {message}{' '}
            <Button variant="outlined" onClick={retry}>
                Retry
            </Button>
        </Typography>
    </div>
);
