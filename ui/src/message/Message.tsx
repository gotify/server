import IconButton from '@material-ui/core/IconButton';
import {withStyles, WithStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Delete from '@material-ui/icons/Delete';
import React from 'react';
import TimeAgo from 'react-timeago';
import Container from '../common/Container';
import * as config from '../config';
import {StyleRulesCallback} from '@material-ui/core/styles/withStyles';
import ReactMarkdown from 'react-markdown';
import {RenderMode, contentType} from './extras';

const styles: StyleRulesCallback = () => ({
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
        maxWidth: 585,
    },
    image: {
        marginRight: 15,
    },
    imageWrapper: {
        display: 'flex',
    },
    content: {
        whiteSpace: 'pre-wrap',
        '& p': {
            margin: 0,
        },
        '& a': {
            color: '#ff7f50',
        },
        '& pre': {
            overflow: 'auto',
        },
    },
});

interface IProps {
    title: string;
    image?: string;
    date: string;
    content: string;
    fDelete: VoidFunction;
    extras?: IMessageExtras;
    height: (height: number) => void;
}

class Message extends React.PureComponent<IProps & WithStyles<typeof styles>> {
    private node: HTMLDivElement | null;

    public componentDidMount = () =>
        this.props.height(this.node ? this.node.getBoundingClientRect().height : 0);

    private renderContent = () => {
        const content = this.props.content;
        switch (contentType(this.props.extras)) {
            case RenderMode.Markdown:
                return <ReactMarkdown source={content} escapeHtml={true} />;
            case RenderMode.Plain:
            default:
                return content;
        }
    };

    public render(): React.ReactNode {
        const {fDelete, classes, title, date, image} = this.props;

        return (
            <div className={`${classes.wrapperPadding} message`} ref={(ref) => (this.node = ref)}>
                <Container style={{display: 'flex'}}>
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
                            <Typography
                                className={`${classes.headerTitle} title`}
                                variant="headline">
                                {title}
                            </Typography>
                            <Typography variant="body1" className="date">
                                <TimeAgo date={date} />
                            </Typography>
                            <IconButton onClick={fDelete} className={`${classes.trash} delete`}>
                                <Delete />
                            </IconButton>
                        </div>
                        <Typography component="div" className={`${classes.content} content`}>
                            {this.renderContent()}
                        </Typography>
                    </div>
                </Container>
            </div>
        );
    }
}

export default withStyles(styles)<IProps>(Message);
