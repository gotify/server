import {WithStyles} from "material-ui";
import Delete from 'material-ui-icons/Delete';
import IconButton from 'material-ui/IconButton';
import {withStyles} from 'material-ui/styles';
import Typography from 'material-ui/Typography';
import React, {Component} from 'react';
import TimeAgo from 'react-timeago';
import Container from './Container';

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

type Style = WithStyles<'header' | 'headerTitle' | 'trash' | 'wrapperPadding' | 'messageContentWrapper' | 'image' | 'imageWrapper'>;

interface IProps {
    title: string
    image?: string
    date: string
    content: string
    fDelete: VoidFunction
}

class Message extends Component<IProps & Style> {
    public render() {
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

export default withStyles(styles)<IProps>(Message);
