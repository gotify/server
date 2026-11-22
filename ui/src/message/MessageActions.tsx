import OpenInNew from '@mui/icons-material/OpenInNew';
import Alert from '@mui/material/Alert';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Tooltip from '@mui/material/Tooltip';
import React from 'react';
import {makeStyles} from 'tss-react/mui';
import {useStores} from '../stores';
import {IMessageExtras} from '../types';
import {actionButtons, IHttpAction} from './extras';

const useStyles = makeStyles()((theme) => ({
    actions: {
        display: 'flex',
        flexWrap: 'wrap',
        gap: theme.spacing(1),
        marginTop: theme.spacing(1.5),
    },
    error: {
        marginTop: theme.spacing(1.5),
    },
    errorDetails: {
        whiteSpace: 'pre-wrap',
        margin: 0,
        marginTop: theme.spacing(1),
    },
}));

const executeAction = (http: IHttpAction): Promise<Response> =>
    fetch(http.url, {
        method: http.method,
        headers: http.headers,
        body: http.body,
        credentials: 'omit',
        referrerPolicy: 'no-referrer',
    });

const describeUrl = (url: string) => {
    const {host, pathname} = new URL(url);
    return `${host}${pathname}`;
};

const HttpActionButton = ({label, http}: {label: string; http: IHttpAction}) => {
    const {snackManager} = useStores();
    const [running, setRunning] = React.useState(false);

    const onClick = async () => {
        setRunning(true);
        try {
            const response = await executeAction(http);
            if (response.ok) {
                snackManager.snack(`Action '${label}' succeeded.`);
            } else {
                snackManager.snack(
                    `Action '${label}' failed: ${response.status} ${response.statusText}`.trim()
                );
            }
        } catch {
            snackManager.snack(
                `Action '${label}': no readable response (network error or CORS). The request may still have been processed.`
            );
        } finally {
            setRunning(false);
        }
    };

    return (
        <Tooltip title={`${http.method} ${describeUrl(http.url)}`}>
            <span>
                <Button
                    className="action"
                    variant="outlined"
                    size="small"
                    disabled={running}
                    onClick={onClick}
                    startIcon={running ? <CircularProgress size={16} /> : undefined}>
                    {label}
                </Button>
            </span>
        </Tooltip>
    );
};

const OpenActionButton = ({label, url}: {label: string; url: string}) => (
    <Tooltip title={`Open ${describeUrl(url)}`}>
        <Button
            className="action"
            variant="outlined"
            size="small"
            href={url}
            target="_blank"
            rel="noopener noreferrer"
            endIcon={<OpenInNew />}>
            {label}
        </Button>
    </Tooltip>
);

const ActionError = ({error}: {error: string}) => {
    const {classes} = useStyles();
    const [showDetails, setShowDetails] = React.useState(false);
    return (
        <Alert
            className={`${classes.error} action-error`}
            severity="warning"
            action={
                <Button
                    className="action-error-toggle"
                    color="inherit"
                    size="small"
                    onClick={() => setShowDetails((show) => !show)}>
                    {showDetails ? 'Hide details' : 'Details'}
                </Button>
            }>
            This message contains invalid action buttons.
            {showDetails && (
                <pre className={`${classes.errorDetails} action-error-details`}>{error}</pre>
            )}
        </Alert>
    );
};

export const MessageActions = ({extras}: {extras?: IMessageExtras}) => {
    const {classes} = useStyles();
    const {buttons, error} = React.useMemo(() => actionButtons(extras), [extras]);
    return (
        <>
            {buttons.length > 0 && (
                <div className={`${classes.actions} actions`}>
                    {buttons.map(({label, http, open}, index) =>
                        http ? (
                            <HttpActionButton key={index} label={label} http={http} />
                        ) : (
                            open && <OpenActionButton key={index} label={label} url={open.url} />
                        )
                    )}
                </div>
            )}
            {error && <ActionError error={error} />}
        </>
    );
};
