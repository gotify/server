import Fab from '@material-ui/core/Fab';
import KeyboardArrowUp from '@material-ui/icons/KeyboardArrowUp';
import React, {Component} from 'react';

class ScrollUpButton extends Component {
    state = {
        display: 'none',
        opacity: 0,
    };
    componentDidMount() {
        window.addEventListener('scroll', this.scrollHandler);
    }

    componentWillUnmount() {
        window.removeEventListener('scroll', this.scrollHandler);
    }

    scrollHandler = () => {
        const currentScrollPos = window.pageYOffset;
        const opacity = Math.min(currentScrollPos / 500, 1);
        const nextState = {display: currentScrollPos > 0 ? 'inherit' : 'none', opacity};
        if (this.state.display !== nextState.display || this.state.opacity !== nextState.opacity) {
            this.setState(nextState);
        }
    };

    public render() {
        return (
            <Fab
                color="primary"
                style={{
                    position: 'fixed',
                    bottom: '30px',
                    right: '30px',
                    zIndex: 100000,
                    display: this.state.display,
                    opacity: this.state.opacity,
                }}
                onClick={this.scrollUp}>
                <KeyboardArrowUp />
            </Fab>
        );
    }

    private scrollUp = () => window.scrollTo(0, 0);
}

export default ScrollUpButton;
