import Alert from '@mui/material/Alert';
import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Tooltip from '@mui/material/Tooltip';
import React from 'react';
import {makeStyles} from 'tss-react/mui';
import {useStores} from '../stores';
import {IMessageExtras} from '../types';
import {actionButtons, IActionButton} from './extras';

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

const executeAction = ({http}: IActionButton): Promise<Response> =>
    fetch(http.url, {
        method: http.method,
        headers: http.headers,
        body: http.body,
        credentials: 'omit',
        referrerPolicy: 'no-referrer',
    });

const describeTarget = ({http}: IActionButton) => {
    const url = new URL(http.url);
    return `${http.method} ${url.host}${url.pathname}`;
};

const ActionButton = ({button}: {button: IActionButton}) => {
    const {snackManager} = useStores();
    const [running, setRunning] = React.useState(false);

    const onClick = async () => {
        setRunning(true);
        try {
            const response = await executeAction(button);
            if (response.ok) {
                snackManager.snack(`Action '${button.label}' succeeded.`);
            } else {
                snackManager.snack(
                    `Action '${button.label}' failed: ${response.status} ${response.statusText}`.trim()
                );
            }
        } catch {
            snackManager.snack(
                `Action '${button.label}': no readable response (network error or CORS). The request may still have been processed.`
            );
        } finally {
            setRunning(false);
        }
    };

    return (
        <Tooltip title={describeTarget(button)}>
            <span>
                <Button
                    className="action"
                    variant="outlined"
                    size="small"
                    disabled={running}
                    onClick={onClick}
                    startIcon={running ? <CircularProgress size={16} /> : undefined}>
                    {button.label}
                </Button>
            </span>
        </Tooltip>
    );
};

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
                    {buttons.map((button, index) => (
                        <ActionButton key={index} button={button} />
                    ))}
                </div>
            )}
            {error && <ActionError error={error} />}
        </>
    );
};
