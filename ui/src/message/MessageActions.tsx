import Button from '@mui/material/Button';
import CircularProgress from '@mui/material/CircularProgress';
import Tooltip from '@mui/material/Tooltip';
import React from 'react';
import {useStores} from '../stores';
import {IActionButton} from './extras';

interface IProps {
    buttons: IActionButton[];
}

export const executeAction = async ({http}: IActionButton): Promise<Response> =>
    fetch(http.url, {
        method: http.method,
        headers: http.headers,
        credentials: 'omit',
        referrerPolicy: 'no-referrer',
    });

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
                `Action '${button.label}' failed: request could not be sent (network or CORS).`
            );
        } finally {
            setRunning(false);
        }
    };

    return (
        <Tooltip title={`${button.http.method} ${new URL(button.http.url).host}`}>
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

const MessageActions = ({buttons}: IProps) => {
    if (buttons.length === 0) {
        return null;
    }
    return (
        <div className="actions" style={{display: 'flex', flexWrap: 'wrap', gap: 8, marginTop: 12}}>
            {buttons.map((button, index) => (
                <ActionButton key={index} button={button} />
            ))}
        </div>
    );
};

export default MessageActions;
