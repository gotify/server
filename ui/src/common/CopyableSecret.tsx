import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Visibility from '@mui/icons-material/Visibility';
import Copy from '@mui/icons-material/FileCopyOutlined';
import VisibilityOff from '@mui/icons-material/VisibilityOff';
import React, {Component, CSSProperties} from 'react';
import {Stores, inject} from '../inject';

interface IProps {
    value: string;
    style?: CSSProperties;
}

interface IState {
    visible: boolean;
}

class CopyableSecret extends Component<IProps & Stores<'snackManager'>, IState> {
    public state = {visible: false};

    public render() {
        const {value, style} = this.props;
        const text = this.state.visible ? value : '•••••••••••••••';
        return (
            <div style={style}>
                <IconButton onClick={this.copyToClipboard} title="Copy to clipboard" size="large">
                    <Copy />
                </IconButton>
                <IconButton
                    onClick={this.toggleVisibility}
                    className="toggle-visibility"
                    size="large">
                    {this.state.visible ? <VisibilityOff /> : <Visibility />}
                </IconButton>
                <Typography style={{fontFamily: 'monospace', fontSize: 16}}>{text}</Typography>
            </div>
        );
    }

    private toggleVisibility = () => this.setState({visible: !this.state.visible});
    private copyToClipboard = async () => {
        const {snackManager, value} = this.props;
        try {
            await navigator.clipboard.writeText(value);
            snackManager.snack('Copied to clipboard');
        } catch (error) {
            console.error('Failed to copy to clipboard:', error);
            snackManager.snack('Failed to copy to clipboard');
        }
    };
}

export default inject('snackManager')(CopyableSecret);
