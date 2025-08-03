import Fab from '@mui/material/Fab';
import KeyboardArrowUp from '@mui/icons-material/KeyboardArrowUp';
import React from 'react';

const ScrollUpButton = () => {
    const [state, setState] = React.useState({display: 'none', opacity: 0});
    React.useEffect(() => {
        const scrollHandler = () => {
            const currentScrollPos = window.pageYOffset;
            const opacity = Math.min(currentScrollPos / 500, 1);
            const nextState = {display: currentScrollPos > 0 ? 'inherit' : 'none', opacity};
            if (state.display !== nextState.display || state.opacity !== nextState.opacity) {
                setState(nextState);
            }
        };
        window.addEventListener('scroll', scrollHandler);
        return () => window.removeEventListener('scroll', scrollHandler);
    }, []);

    return (
        <Fab
            color="primary"
            style={{
                position: 'fixed',
                bottom: '30px',
                right: '30px',
                zIndex: 100000,
                display: state.display,
                opacity: state.opacity,
            }}
            onClick={() => window.scrollTo(0, 0)}>
            <KeyboardArrowUp />
        </Fab>
    );
};

export default ScrollUpButton;
