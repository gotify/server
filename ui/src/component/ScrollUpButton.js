import React, {Component} from 'react';
import Button from 'material-ui/Button';
import KeyboardArrowUp from 'material-ui-icons/KeyboardArrowUp';

class ScrollUpButton extends Component {
    render() {
        return (
            <Button variant="fab" color="primary"
                    style={{position: 'fixed', bottom: '30px', right: '30px', zIndex: 100000}}
                    onClick={() => window.scrollTo(0, 0)}>
                <KeyboardArrowUp/>
            </Button>
        );
    }
}

export default ScrollUpButton;
