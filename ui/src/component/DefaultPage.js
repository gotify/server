import React, {Component} from 'react';
import Button from 'material-ui/Button';
import Grid from 'material-ui/Grid';
import Typography from 'material-ui/Typography';
import PropTypes from 'prop-types';

export default class DefaultPage extends Component {
    static defaultProps = {
        buttonDisabled: false,
        hideButton: false,
        maxWidth: 700,
    };

    static propTypes = {
        title: PropTypes.string.isRequired,
        buttonTitle: PropTypes.string,
        fButton: PropTypes.func,
        buttonDisabled: PropTypes.bool.isRequired,
        maxWidth: PropTypes.number.isRequired,
        hideButton: PropTypes.bool.isRequired,
        children: PropTypes.oneOfType([
            PropTypes.arrayOf(PropTypes.node),
            PropTypes.node,
        ]).isRequired,
    };

    render() {
        const {title, buttonTitle, fButton, buttonDisabled, maxWidth, hideButton, children} = this.props;
        return (
            <main style={{margin: '0 auto', maxWidth}}>
                <Grid container spacing={24}>
                    <Grid item xs={12} style={{display: 'flex'}}>
                        <Typography variant="display1" style={{flex: 1}}>
                            {title}
                        </Typography>
                        {hideButton ? null : <Button variant="raised" color="primary" disabled={buttonDisabled}
                                                     onClick={fButton}>{buttonTitle}</Button>}
                    </Grid>
                    {children}
                </Grid>
            </main>
        );
    }
}
