import React, {Component} from 'react';
import {withStyles} from 'material-ui/styles';
import Typography from 'material-ui/Typography';
import IconButton from 'material-ui/IconButton';
import PropTypes from 'prop-types';
import Container from './Container';
import TimeAgo from 'react-timeago';
import Delete from 'material-ui-icons/Delete';

const styles = () => ({
    header: {
        display: 'flex',
    },
    headerTitle: {
        flex: 1,
    },
    trash: {
        marginTop: -15,
        marginRight: -15,
    },
    wrapperPadding: {
        padding: 12,
    },
    messageContentWrapper: {
        width: '100%',
    },
    image: {
        marginRight: 15,
    },
    imageWrapper: {
        display: 'flex',
    },
});

class Message extends Component {
    static propTypes = {
        classes: PropTypes.object.isRequired,
        title: PropTypes.string.isRequired,
        image: PropTypes.string,
        date: PropTypes.string.isRequired,
        content: PropTypes.string.isRequired,
        fDelete: PropTypes.func.isRequired,
    };
    render() {
        const {fDelete, classes, title, date, content, image} = this.props;

        return (
            <div className={classes.wrapperPadding}>
                <Container style={{display: 'flex'}}>
                    <div className={classes.imageWrapper}>
                        <img src={image} alt="app logo" width="70" height="70" className={classes.image}/>
                    </div>
                    <div className={classes.messageContentWrapper}>
                        <div className={classes.header}>
                            <Typography className={classes.headerTitle} variant="headline">
                                {title}
                            </Typography>
                            <Typography variant="body1">
                                <TimeAgo date={date}/>
                            </Typography>
                            <IconButton onClick={fDelete} className={classes.trash}><Delete/></IconButton>
                        </div>
                        <Typography component="p">{content}</Typography>
                    </div>
                </Container>
            </div>
        );
    }
}

export default withStyles(styles)(Message);
