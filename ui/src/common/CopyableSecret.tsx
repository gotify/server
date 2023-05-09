import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import Visibility from '@material-ui/icons/Visibility';
import Copy from '@material-ui/icons/FileCopyOutlined';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
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
                <IconButton onClick={this.copyToClipboard} title="Copy to clipboard">
                    <Copy />
                </IconButton>
                <IconButton onClick={this.toggleVisibility} className="toggle-visibility">
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
