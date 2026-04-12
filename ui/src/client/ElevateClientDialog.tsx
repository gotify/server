import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import {observer} from 'mobx-react-lite';
import {useStores} from '../stores';
import ElevationForm from '../common/ElevationForm';

interface IProps {
    clientName: string;
    clientId: number;
    fClose: VoidFunction;
}

const durationOptions = [
    {label: 'Cancel elevation', seconds: -1},
    {label: '1 hour', seconds: 60 * 60},
    {label: '1 day', seconds: 24 * 60 * 60},
    {label: '30 days', seconds: 30 * 24 * 60 * 60},
    {label: '1 year', seconds: 365 * 24 * 60 * 60},
];

const ElevateClientDialog = observer(({clientName, clientId, fClose}: IProps) => {
    const {elevateStore, clientStore, currentUser} = useStores();
    const [durationSeconds, setDurationSeconds] = useState(durationOptions[1].seconds);

    const needsElevation = !elevateStore.elevated;

    const handleConfirm = async () => {
        await clientStore.elevate(clientId, durationSeconds);
        if (clientId === currentUser.user.clientId) {
            currentUser.tryAuthenticate();
        }
        fClose();
    };

    const handleClose = () => {
        elevateStore.cleanupOidcElevate();
        fClose();
    };

    return (
        <Dialog open={true} onClose={handleClose} aria-labelledby="elevate-client-dialog-title">
            <DialogTitle id="elevate-client-dialog-title">Elevate Client: {clientName}</DialogTitle>
            <DialogContent>
                {needsElevation ? (
                    <ElevationForm />
                ) : (
                    <FormControl fullWidth style={{marginTop: 8}}>
                        <InputLabel id="elevate-duration-label">Duration</InputLabel>
                        <Select
                            labelId="elevate-duration-label"
                            label="Duration"
                            value={durationSeconds}
                            onChange={(e) => setDurationSeconds(e.target.value as number)}>
                            {durationOptions.map((opt) => (
                                <MenuItem key={opt.seconds} value={opt.seconds}>
                                    {opt.label}
                                </MenuItem>
                            ))}
                        </Select>
                    </FormControl>
                )}
            </DialogContent>
            <DialogActions>
                <Button onClick={handleClose}>Cancel</Button>
                {!needsElevation && (
                    <Button onClick={handleConfirm} autoFocus color="primary" variant="contained">
                        Elevate
                    </Button>
                )}
            </DialogActions>
        </Dialog>
    );
});

export default ElevateClientDialog;
