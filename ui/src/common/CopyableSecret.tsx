import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Visibility from '@mui/icons-material/Visibility';
import Copy from '@mui/icons-material/FileCopyOutlined';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import React, {CSSProperties, useState} from 'react';
import {useAppDispatch} from '../store';
import {uiActions} from '../store/ui-slice.ts';

interface IProps {
    value: string;
    style?: CSSProperties;
}

const CopyableSecret = ({value, style}: IProps) => {
    const dispatch = useAppDispatch();
    const [ visible, setVisible ] = useState(false);
    const text = visible ? value : '•••••••••••••••';

    const toggleVisibility = () => setVisible((prevState) => !prevState);
    const copyToClipboard = async () => {
        try {
            await navigator.clipboard.writeText(value);
            dispatch(uiActions.addSnackMessage('Copied to clipboard'));
        } catch (error) {
            console.error('Failed to copy to clipboard:', error);
            dispatch(uiActions.addSnackMessage('Failed to copy to clipboard'));
        }
    };

    return (
        <div style={style}>
            <IconButton onClick={copyToClipboard} title="Copy to clipboard">
                <Copy />
            </IconButton>
            <IconButton onClick={toggleVisibility} className="toggle-visibility">
                {visible ? <VisibilityOff /> : <Visibility />}
            </IconButton>
            <Typography style={{fontFamily: 'monospace', fontSize: 16}}>{text}</Typography>
        </div>
    );
};

export default CopyableSecret;
