import React, {useState} from 'react';
import Button from '@mui/material/Button';
import TextField from '@mui/material/TextField';
import Typography from '@mui/material/Typography';
import {observer} from 'mobx-react-lite';
import {useStores} from '../stores';
import * as config from '../config';
import CircularProgress from '@mui/material/CircularProgress';
import {Box, Divider} from '@mui/material';

const ElevateDuration = 60 * 60;

const ElevationForm = observer(() => {
    const {elevateStore} = useStores();
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');

    const oidcEnabled = config.get('oidc');
    const oidcPending = elevateStore.oidcElevatePending;

    const handleLocalElevate = async () => {
        try {
            await elevateStore.localElevate(password, ElevateDuration);
        } catch {
            setError('Elevation failed. Check your password.');
        }
    };

    if (oidcPending) {
        return (
            <Box sx={{textAlign: 'center', my: 2}}>
                <CircularProgress sx={{mb: 2}} />
                <Typography sx={{mb: 1}}>Waiting for OIDC sign-in...</Typography>
                <Typography variant="body2" color="textSecondary" sx={{mb: 1}}>
                    Complete sign-in in the new tab, then close it to continue.
                </Typography>
                <Button
                    className="elevation-oidc-cancel"
                    variant="outlined"
                    fullWidth
                    onClick={() => elevateStore.cleanupOidcElevate()}>
                    Cancel OIDC Login
                </Button>
            </Box>
        );
    }

    return (
        <>
            <Typography>This action requires re-authentication.</Typography>
            <form
                onSubmit={(e) => {
                    e.preventDefault();
                    handleLocalElevate();
                }}>
                <TextField
                    autoFocus
                    margin="dense"
                    type="password"
                    label="Password"
                    className="elevation-password"
                    value={password}
                    onChange={(e) => {
                        setPassword(e.target.value);
                        setError('');
                    }}
                    fullWidth
                    error={!!error}
                    helperText={error}
                />
                <Button
                    type="submit"
                    className="elevation-submit"
                    disabled={password.length === 0}
                    color="primary"
                    variant="contained"
                    fullWidth>
                    Elevate with Password
                </Button>
            </form>

            {oidcEnabled && (
                <>
                    <Divider sx={{my: 2}}>or</Divider>
                    <Button
                        className="elevation-oidc"
                        variant="contained"
                        color="primary"
                        fullWidth
                        onClick={() => elevateStore.oidcElevate(ElevateDuration)}>
                        Elevate via OIDC
                    </Button>
                </>
            )}
        </>
    );
});

export default ElevationForm;
