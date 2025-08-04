import {Button, Theme} from '@mui/material';
import IconButton from '@mui/material/IconButton';
import {makeStyles} from 'tss-react/mui';
import Typography from '@mui/material/Typography';
import {ExpandLess, ExpandMore} from '@mui/icons-material';
import Delete from '@mui/icons-material/Delete';
import React from 'react';
import TimeAgo from 'react-timeago';
import Container from '../common/Container';
import {Markdown} from '../common/Markdown';
import * as config from '../config';
import {IMessageExtras} from '../types';
import {contentType, RenderMode} from './extras';
import {makeIntlFormatter} from 'react-timeago/defaultFormatter';

const PREVIEW_LENGTH = 500;

const useStyles = makeStyles()((theme: Theme) => ({
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
        marginBottom: 12,
    },
    messageContentWrapper: {
        minWidth: 200,
    },
    image: {
        marginRight: 15,
        [theme.breakpoints.down('md')]: {
            width: 32,
            height: 32,
        },
        [theme.breakpoints.down('sm')]: {
            display: 'none',
        },
    },
    appName: {
        opacity: 0.7,
    },
    date: {
        [theme.breakpoints.down('md')]: {
            opacity: 0.7,
        },
    },
    messageWrapper: {
        display: 'grid',
        gridTemplateColumns: '70px 1fr',
        [theme.breakpoints.down('sm')]: {
            gridTemplateColumns: '1fr',
        },
        gap: 16,
        width: '100%',
    },
    actionsWrapper: {
        flex: 1,
        display: 'flex',
        justifyContent: 'flex-end',
        flexDirection: 'row',
        alignItems: 'flex-start',
    },
    imageWrapper: {
        display: 'flex',
    },
    plainContent: {
        whiteSpace: 'pre-wrap',
    },
    content: {
        maxHeight: PREVIEW_LENGTH,
        wordBreak: 'break-all',
        overflowY: 'hidden',
        '&.expanded': {
            maxHeight: 'none',
        },
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
}));

interface IProps {
    title: string;
    image?: string;
    date: string;
    content: string;
    priority: number;
    appName: string;
    fDelete: VoidFunction;
    extras?: IMessageExtras;
    expanded: boolean;
    onExpand: (expand: boolean) => void;
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

const Message = ({
    fDelete,
    title,
    date,
    image,
    priority,
    content,
    extras,
    appName,
    onExpand,
    expanded: initialExpanded,
}: IProps) => {
    const [previewRef, setPreviewRef] = React.useState<HTMLDivElement | null>(null);
    const {classes} = useStyles();
    const [expanded, setExpanded] = React.useState(initialExpanded);
    const [isOverflowing, setOverflowing] = React.useState(false);

    React.useEffect(() => {
        setOverflowing(!!previewRef && previewRef.scrollHeight > previewRef.clientHeight);
    }, [previewRef]);

    React.useEffect(() => void onExpand(expanded), [expanded]);

    const togglePreviewHeight = () => setExpanded((b) => !b);

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
        <div className={`${classes.wrapperPadding} message`}>
            <Container
                style={{
                    display: 'flex',
                    flexWrap: 'wrap',
                    borderLeftColor: priorityColor(priority),
                    borderLeftWidth: 6,
                    borderLeftStyle: 'solid',
                }}>
                <div className={classes.messageWrapper}>
                    {image !== null ? (
                        <img
                            src={config.get('url') + image}
                            alt={`${appName} logo`}
                            width="70"
                            height="70"
                            className={classes.image}
                        />
                    ) : null}
                    <div style={{width: '100%'}}>
                        <div style={{display: 'flex', width: '100%'}}>
                            <div className={classes.messageContentWrapper}>
                                <div className={classes.header}>
                                    <Typography
                                        className={`${classes.headerTitle} title`}
                                        variant="h5">
                                        {title}
                                    </Typography>
                                </div>
                                <div className={classes.appName}>
                                    <Typography variant="body1" className={classes.date}>
                                        {appName}
                                    </Typography>
                                </div>
                            </div>
                            <div className={classes.actionsWrapper}>
                                <Typography variant="body1" className={classes.date}>
                                    <TimeAgo
                                        date={date}
                                        formatter={makeIntlFormatter({
                                            style: 'narrow',
                                        })}
                                    />
                                </Typography>
                                <IconButton
                                    onClick={fDelete}
                                    className={`${classes.trash} delete`}
                                    size="large">
                                    <Delete />
                                </IconButton>
                            </div>
                        </div>
                        <div>
                            <Typography
                                component="div"
                                ref={setPreviewRef}
                                className={`${classes.content} content ${
                                    isOverflowing && expanded ? 'expanded' : ''
                                }`}>
                                {renderContent()}
                            </Typography>
                        </div>
                        {isOverflowing && (
                            <Button
                                style={{marginTop: 16}}
                                onClick={togglePreviewHeight}
                                variant="contained"
                                color="primary"
                                size="large"
                                fullWidth={true}
                                startIcon={expanded ? <ExpandLess /> : <ExpandMore />}>
                                {expanded ? 'Read Less' : 'Read More'}
                            </Button>
                        )}
                    </div>
                </div>
            </Container>
        </div>
    );
};

export default Message;
