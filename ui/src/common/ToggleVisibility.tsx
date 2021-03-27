import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';
import React, {Component, CSSProperties} from 'react';

interface IProps {
    value: string;
    style?: CSSProperties;
}

interface IState {
    visible: boolean;
}

class ToggleVisibility extends Component<IProps, IState> {
    public state = {visible: false};

    public render() {
        const {value, style} = this.props;
        const text = this.state.visible ? value : '•••••••••••••••';
        return (
            <div style={style}>
                <IconButton onClick={this.toggleVisibility} className="toggle-visibility">
                    {this.state.visible ? <VisibilityOff /> : <Visibility />}
                </IconButton>
                <Typography style={{fontFamily: 'monospace', fontSize: 16}}>{text}</Typography>
            </div>
        );
    }

    private toggleVisibility = () => this.setState({visible: !this.state.visible});
}

export default ToggleVisibility;
