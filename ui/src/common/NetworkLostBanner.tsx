import React, { useEffect } from 'react';
import Button from '@material-ui/core/Button';
import Typography from '@material-ui/core/Typography';

interface NetworkLostBannerProps {
    height: number;
    retry: () => void;
}

export const NetworkLostBanner = ({height, retry}: NetworkLostBannerProps) => {
    useEffect(() => {
        const intervalId = setInterval(retry, 3000);
            
        return() => {
                clearInterval(intervalId);
            }
        });

    return (
        <div
            style={{
                backgroundColor: '#e74c3c',
                height,
                width: '100%',
                zIndex: 1300,
                position: 'relative',
            }}>
            <Typography align="center" variant="h6" style={{lineHeight: `${height}px`}}>
                No network connection.{' '}
                <Button variant="outlined" onClick={retry}>
                    Retry
                </Button>
            </Typography>
        </div>
    );

};
