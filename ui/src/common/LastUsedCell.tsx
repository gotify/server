import {Typography} from '@mui/material';
import React from 'react';
import TimeAgo from 'react-timeago';

export const LastUsedCell: React.FC<{lastUsed: string | null}> = ({lastUsed}) => {
    if (lastUsed === null) {
        return <Typography>Never</Typography>;
    }

    if (+new Date(lastUsed) + 300000 > Date.now()) {
        return <Typography title={lastUsed}>Recently</Typography>;
    }

    return <TimeAgo date={lastUsed} />;
};
