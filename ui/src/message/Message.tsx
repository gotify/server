import React, {useEffect, useRef} from 'react';
import IconButton from '@mui/material/IconButton';
import {Theme} from '@mui/material/styles';
import Typography from '@mui/material/Typography';
import DeleteIcon from '@mui/icons-material/Delete';
import TimeAgo from 'react-timeago';
import {makeStyles} from 'tss-react/mui';
import Container from '../common/Container';
import * as config from '../config';
import {Markdown} from '../common/Markdown';
import {RenderMode, contentType} from './extras';
import {IMessageExtras} from '../types';

const useStyles = makeStyles()((theme: Theme) => {
    return {
        header: {
            display: 'flex',
            flexWrap: 'wrap',
            marginBottom: 0,
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
            maxWidth: 585,
        },
        image: {
            marginRight: 15,
            [theme.breakpoints.down('sm')]: {
                width: 32,
                height: 32,
            },
        },
        date: {
            [theme.breakpoints.down('sm')]: {
                order: 1,
                flexBasis: '100%',
                opacity: 0.7,
            },
        },
        imageWrapper: {
            display: 'flex',
        },
        plainContent: {
            whiteSpace: 'pre-wrap',
        },
        content: {
            wordBreak: 'break-all',
            '& p': {
                margin: 0,
            },
            '& a': {
                color: '#ff7f50',
            },
            '& pre': {
                overflow: 'auto',
            },
            '& img': {
                maxWidth: '100%',
            },
        },
    };
});

interface IProps {
    title: string;
    image?: string;
    date: string;
    content: string;
    priority: number;
    fDelete: VoidFunction;
    extras?: IMessageExtras;
    height: (height: number) => void;
}

const priorityColor = (priority: number) => {
    if (priority >= 4 && priority <= 7) {
        return 'rgba(230, 126, 34, 0.7)';
    } else if (priority > 7) {
        return '#e74c3c';
    } else {
        return 'transparent';
    }
};

const Message = ({fDelete, title, date, image, priority, content, extras, height}: IProps) => {
    const {classes} = useStyles();
    const node = useRef<HTMLDivElement>(null);

    useEffect(() => {
        // TODO: fix this
        // height(node ? node.getBoundingClientRect().height : 0);
    }, []);

    const renderContent = () => {
        switch (contentType(extras)) {
            case RenderMode.Markdown:
                return <Markdown>{content}</Markdown>;
            case RenderMode.Plain:
            default:
                return <span className={classes.plainContent}>{content}</span>;
        }
    };

    return (
        <div className={`${classes.wrapperPadding} message`} ref={node}>
            <Container
                style={{
                    display: 'flex',
                    borderLeftColor: priorityColor(priority),
                    borderLeftWidth: 6,
                    borderLeftStyle: 'solid',
                }}>
                <div className={classes.imageWrapper}>
                    {image !== null ? (
                        <img
                            src={config.get('url') + image}
                            alt="app logo"
                            width="70"
                            height="70"
                            className={classes.image}
                        />
                    ) : null}
                </div>
                <div className={classes.messageContentWrapper}>
                    <div className={classes.header}>
                        <Typography className={`${classes.headerTitle} title`} variant="h5">
                            {title}
                        </Typography>
                        <Typography variant="body1" className={classes.date}>
                            <TimeAgo date={date} />
                        </Typography>
                        <IconButton onClick={fDelete} className={`${classes.trash} delete`}>
                            <DeleteIcon />
                        </IconButton>
                    </div>
                    <Typography component="div" className={`${classes.content} content`}>
                        {renderContent()}
                    </Typography>
                </div>
            </Container>
        </div>
    );
};

export default Message;
