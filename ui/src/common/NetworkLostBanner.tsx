import React from 'react';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

interface NetworkLostBannerProps {
    height: number;
    retry: () => void;
}

export const NetworkLostBanner = ({height, retry}: NetworkLostBannerProps) => {
    return (
        <div
            style={{
                backgroundColor: '#e74c3c',
                height,
                width: '100%',
                zIndex: 1300,
                position: 'relative',
            }}>
            <Typography align="center" variant="title" style={{lineHeight: `${height}px`}}>
                No network connection.{' '}
                <Button variant="outlined" onClick={retry}>
                    Retry
                </Button>
            </Typography>
        </div>
    );
};
