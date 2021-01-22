import Fab from '@material-ui/core/Fab';
import KeyboardArrowUp from '@material-ui/icons/KeyboardArrowUp';
import React, {Component} from 'react';

class ScrollUpButton extends Component {
    state = {
        display: 'none',
        opacity: 0,
    };
    componentDidMount() {
        if (typeof window !== 'undefined') {
            window.addEventListener('scroll', () => {
                let currentScrollPos = window.pageYOffset;
                if (currentScrollPos > 0) {
                    this.setState({display: 'inherit'});
                    this.setState({opacity: currentScrollPos / 500});
                } else {
                    this.setState({display: 'none'});
                    this.setState({opacity: 0});
                }
            });
        }
    }
    public render() {
        return (
            <Fab
                color="primary"
                style={{
                    position: 'fixed',
                    bottom: '30px',
                    right: '30px',
                    zIndex: 100000,
                    display: `${this.state.display}`,
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
