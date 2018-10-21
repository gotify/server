import IconButton from '@material-ui/core/IconButton';
import {withStyles, WithStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Delete from '@material-ui/icons/Delete';
import React from 'react';
import TimeAgo from 'react-timeago';
import Container from '../component/Container';

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

type Style = WithStyles<
    | 'header'
    | 'headerTitle'
    | 'trash'
    | 'wrapperPadding'
    | 'messageContentWrapper'
    | 'image'
    | 'imageWrapper'
>;

interface IProps {
    title: string;
    image?: string;
    date: string;
    content: string;
    fDelete: VoidFunction;
}

function Message({fDelete, classes, title, date, content, image}: IProps & Style) {
    return (
        <div className={`${classes.wrapperPadding} message`}>
            <Container style={{display: 'flex'}}>
                <div className={classes.imageWrapper}>
                    <img
                        src={image}
                        alt="app logo"
                        width="70"
                        height="70"
                        className={classes.image}
                    />
                </div>
                <div className={classes.messageContentWrapper}>
                    <div className={classes.header}>
                        <Typography className={`${classes.headerTitle} title`} variant="headline">
                            {title}
                        </Typography>
                        <Typography variant="body1" className="date">
                            <TimeAgo date={date} />
                        </Typography>
                        <IconButton onClick={fDelete} className={classes.trash}>
                            <Delete className="delete" />
                        </IconButton>
                    </div>
                    <Typography component="p" className="content">
                        {content}
                    </Typography>
                </div>
            </Container>
        </div>
    );
}

export default withStyles(styles)<IProps>(Message);
