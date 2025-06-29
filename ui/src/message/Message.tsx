import {Button, Collapse} from '@material-ui/core';
import IconButton from '@material-ui/core/IconButton';
import {createStyles, Theme, withStyles, WithStyles} from '@material-ui/core/styles';
import {ClassNameMap} from '@material-ui/core/styles/withStyles';
import Typography from '@material-ui/core/Typography';
import {ExpandLess, ExpandMore} from '@material-ui/icons';
import Delete from '@material-ui/icons/Delete';
import React, {RefObject} from 'react';
import TimeAgo from 'react-timeago';
import Container from '../common/Container';
import {Markdown} from '../common/Markdown';
import * as config from '../config';
import {IMessageExtras} from '../types';
import {contentType, RenderMode} from './extras';

const PREVIEW_LENGTH = 500;
const ANIMATION_TIMEOUT_MS = 500;

const styles = (theme: Theme) =>
    createStyles({
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
            maxHeight: PREVIEW_LENGTH,
            wordBreak: 'break-all',
            '&.expanded, &.collapsed': {
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

interface IState {
    expanded: boolean;
    isOverflowing: boolean;
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

class Message extends React.PureComponent<IProps & WithStyles<typeof styles>> {
    private node: HTMLDivElement | null = null;
    public state: IState = {expanded: false, isOverflowing: false};
    previewRef: RefObject<HTMLDivElement>;

    constructor(props: IProps) {
        super(
            props as IProps & {
                classes: ClassNameMap<
                    | 'header'
                    | 'image'
                    | 'content'
                    | 'headerTitle'
                    | 'trash'
                    | 'wrapperPadding'
                    | 'messageContentWrapper'
                    | 'date'
                    | 'imageWrapper'
                    | 'plainContent'
                >;
            }
        );
        this.previewRef = React.createRef();
    }

    public componentDidMount = () => {
        if (this.previewRef.current) {
            this.setState({
                ...this.state,
                isOverflowing:
                    this.previewRef.current.scrollHeight > this.previewRef.current.clientHeight,
            });
        }
        return this.props.height(this.node ? this.node.getBoundingClientRect().height : 0);
    };

    public togglePreviewHeight = async () => {
        this.setState({...this.state, expanded: !this.state.expanded});
        await new Promise((resolve) => setTimeout(resolve, ANIMATION_TIMEOUT_MS + 1));
        this.props.height(this.node ? this.node.getBoundingClientRect().height : 0);
    };

    private renderContent = () => {
        const content = this.props.content;
        switch (contentType(this.props.extras)) {
            case RenderMode.Markdown:
                return <Markdown>{content}</Markdown>;
            case RenderMode.Plain:
            default:
                return <span className={this.props.classes.plainContent}>{content}</span>;
        }
    };

    public render(): React.ReactNode {
        const {fDelete, classes, title, date, image, priority} = this.props;

        return (
            <div className={`${classes.wrapperPadding} message`} ref={(ref) => (this.node = ref)}>
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
                                <Delete />
                            </IconButton>
                        </div>

                        {!this.state.isOverflowing ? (
                            <Typography
                                component="div"
                                className={`${classes.content} content}`}
                                ref={this.previewRef}>
                                {this.renderContent()}
                            </Typography>
                        ) : (
                            <Collapse
                                in={this.state.expanded}
                                timeout={ANIMATION_TIMEOUT_MS}
                                collapsedSize={PREVIEW_LENGTH}>
                                <Typography
                                    component="div"
                                    className={`${classes.content} content ${this.state.expanded ? 'expanded' : 'collapsed'}`}>
                                    {this.renderContent()}
                                </Typography>
                            </Collapse>
                        )}

                        {this.state.isOverflowing && (
                            <Button
                                onClick={() => this.togglePreviewHeight()}
                                color="primary"
                                size="small"
                                startIcon={this.state.expanded ? <ExpandLess /> : <ExpandMore />}>
                                {this.state.expanded ? 'Read Less' : 'Read More'}
                            </Button>
                        )}
                    </div>
                </Container>
            </div>
        );
    }
}

export default withStyles(styles, {withTheme: true})(Message);
