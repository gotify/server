import IconButton from '@material-ui/core/IconButton';
import {createStyles, Theme, withStyles, WithStyles} from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Delete from '@material-ui/icons/Delete';
import React from 'react';
import TimeAgo from 'react-timeago';
import Container from '../common/Container';
import * as config from '../config';
import {Markdown} from '../common/Markdown';
import {RenderMode, contentType} from './extras';
import {IMessageExtras} from '../types';

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
        content: {
            whiteSpace: 'pre-wrap',
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
    private node: HTMLDivElement | null = null;

    public componentDidMount = () =>
        this.props.height(this.node ? this.node.getBoundingClientRect().height : 0);

    private renderContent = () => {
        const content = this.props.content;
        switch (contentType(this.props.extras)) {
            case RenderMode.Markdown:
                return <Markdown>{content}</Markdown>;
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
                        <Typography component="div" className={`${classes.content} content`}>
                            {this.renderContent()}
                        </Typography>
                    </div>
                </Container>
            </div>
        );
    }
}

export default withStyles(styles, {withTheme: true})(Message);
