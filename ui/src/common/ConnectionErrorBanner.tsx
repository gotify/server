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
            minHeight: height,
            width: '100%',
            zIndex: 1300,
            position: 'relative',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: '8px 16px',
            boxSizing: 'border-box',
        }}>
        <Typography align="center" variant="h6" style={{lineHeight: 1.4}}>
            {message}{' '}
            <Button variant="outlined" onClick={retry}>
                Retry
            </Button>
        </Typography>
    </div>
);
