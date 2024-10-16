import React from 'react';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import {useAppDispatch} from '../store';
import {tryReconnect} from '../store/auth-actions.ts';

interface ConnectionErrorBannerProps {
    height: number;
    message: string;
}

export const ConnectionErrorBanner = ({height, message}: ConnectionErrorBannerProps) => {
    const dispatch = useAppDispatch();

    const handleRetry = async () => {
        await dispatch(tryReconnect());
    }

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
                {message}{' '}
                <Button variant="outlined" onClick={handleRetry}>
                    Retry
                </Button>
            </Typography>
        </div>
    );
};
