import KeyboardArrowUp from 'material-ui-icons/KeyboardArrowUp';
import Button from 'material-ui/Button';
import React, {Component} from 'react';

class ScrollUpButton extends Component {
    public render() {
        return (
            <Button variant="fab" color="primary"
                    style={{position: 'fixed', bottom: '30px', right: '30px', zIndex: 100000}}
                    onClick={this.scrollUp}>
                <KeyboardArrowUp/>
            </Button>
        );
    }

    private scrollUp = () => window.scrollTo(0, 0);
}

export default ScrollUpButton;
