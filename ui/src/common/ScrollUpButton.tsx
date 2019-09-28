import Fab from '@material-ui/core/Fab';
import KeyboardArrowUp from '@material-ui/icons/KeyboardArrowUp';
import React, {Component} from 'react';

class ScrollUpButton extends Component {
    public render() {
        return (
            <Fab
                color="primary"
                style={{position: 'fixed', bottom: '30px', right: '30px', zIndex: 100000}}
                onClick={this.scrollUp}>
                <KeyboardArrowUp />
            </Fab>
        );
    }

    private scrollUp = () => window.scrollTo(0, 0);
}

export default ScrollUpButton;
