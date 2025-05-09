import IconButton from '@mui/material/IconButton';
import Snackbar from '@mui/material/Snackbar';
import Close from '@mui/icons-material/Close';
import React, {useEffect, useState} from 'react';
import { useAppDispatch, useAppSelector} from '../store';
import {uiActions} from '../store/ui-slice.ts';

const MAX_VISIBLE_SNACK_TIME_IN_MS = 6000;
const MIN_VISIBLE_SNACK_TIME_IN_MS = 1000;

const SnackBarHandler = () => {
    const dispatch = useAppDispatch();
    const [open, setOpen] = useState(false);
    const [openWhen, setOpenWhen] = useState(0);
    const snackMessageCounter = useAppSelector(state => state.ui.snack.messages.length);
    const snackMessage = useAppSelector(state => state.ui.snack.message);

    const closeCurrentSnack = () => setOpen(false);

    useEffect(() => {
        // if (snackMessageCounter === 0) {
        //     setOpen(false);
        //     setOpenWhen(0);
        //     return;
        // }

        if (!open) {
            openNextSnack();
            return;
        }

        const snackOpenSince = Date.now() - openWhen;
        if (snackOpenSince > MIN_VISIBLE_SNACK_TIME_IN_MS) {
            closeCurrentSnack();
        } else {
            setTimeout(closeCurrentSnack, MIN_VISIBLE_SNACK_TIME_IN_MS - snackOpenSince);
        }
    }, [snackMessageCounter]);

    const openNextSnack = () => {
        if (snackMessageCounter > 0) {
            setOpen(true);
            setOpenWhen(Date.now());
            dispatch(uiActions.nextSnackMessage());
        }
    };

    const duration =
        snackMessageCounter > 1 ? MIN_VISIBLE_SNACK_TIME_IN_MS : MAX_VISIBLE_SNACK_TIME_IN_MS;

    return (
        <Snackbar
            anchorOrigin={{vertical: 'bottom', horizontal: 'left'}}
            open={open}
            autoHideDuration={duration}
            onClose={closeCurrentSnack}
            TransitionProps={{onExited: openNextSnack}}
            message={<span id="message-id">{snackMessage}</span>}
            action={
                <IconButton
                    key="close"
                    aria-label="Close"
                    color="inherit"
                    onClick={closeCurrentSnack}>
                    <Close />
                </IconButton>
            }
        />
    );
};

export default SnackBarHandler;
