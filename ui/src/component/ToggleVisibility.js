import React, {Component} from 'react';
import Typography from 'material-ui/Typography';
import IconButton from 'material-ui/IconButton';
import VisibilityOff from 'material-ui-icons/VisibilityOff';
import Visibility from 'material-ui-icons/Visibility';
import PropTypes from 'prop-types';

class ToggleVisibility extends Component {
    static propTypes = {
        value: PropTypes.string.isRequired,
        style: PropTypes.object,
    };

    constructor() {
        super();
        this.state = {visible: false};
    }

    toggleVisibility = () => this.setState({visible: !this.state.visible});

    render() {
        const {value, style} = this.props;
        const text = this.state.visible ? value : '•••••••••••••••';
        return (
            <div style={style}>
                <IconButton onClick={this.toggleVisibility}>
                    {this.state.visible ? <VisibilityOff/> : <Visibility/>}
                </IconButton>
                <Typography style={{fontFamily: '\'Roboto Mono\', monospace'}}>
                    {text}
                </Typography>
            </div>
        );
    }
}


export default ToggleVisibility;
