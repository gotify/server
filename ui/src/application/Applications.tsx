import React, {ChangeEvent, useEffect, useRef, useState} from 'react';
import Grid from '@mui/material/Grid';
import IconButton from '@mui/material/IconButton';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Delete from '@mui/icons-material/Delete';
import Edit from '@mui/icons-material/Edit';
import CloudUpload from '@mui/icons-material/CloudUpload';
import Button from '@mui/material/Button';

import ConfirmDialog from '../common/ConfirmDialog';
import DefaultPage from '../common/DefaultPage';
import CopyableSecret from '../common/CopyableSecret';
import {AddApplicationDialog} from './AddApplicationDialog';
import * as config from '../config';
import {UpdateApplicationDialog} from './UpdateApplicationDialog';
import {IApplication} from '../types';
import {LastUsedCell} from '../common/LastUsedCell';
import {useStores} from '../stores';
import {observer} from 'mobx-react-lite';
import {makeStyles} from 'tss-react/mui';
import {ButtonBase, Tooltip} from '@mui/material';

const useStyles = makeStyles()((theme) => ({
    imageContainer: {
        '&::after': {
            content: '"×"',
            position: 'absolute',
            top: 0,
            left: 0,
            width: 40,
            height: 40,
            background: theme.palette.error.main,
            color: theme.palette.getContrastText(theme.palette.error.main),
            fontSize: 40,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            opacity: 0,
        },
        '&:hover::after': {opacity: 1},
    },
}));

const Applications = observer(() => {
    const {appStore} = useStores();
    const apps = appStore.getItems();
    const [toDeleteApp, setToDeleteApp] = useState<IApplication>();
    const [toDeleteImage, setToDeleteImage] = useState<IApplication>();
    const [toUpdateApp, setToUpdateApp] = useState<IApplication>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);

    const fileInputRef = useRef<HTMLInputElement>(null);
    const uploadId = useRef(-1);

    useEffect(() => void appStore.refresh(), []);

    const validExtensions = ['.gif', '.png', '.jpg', '.jpeg'];

    const handleImageUploadClick = (id: number) => {
        uploadId.current = id;
        if (fileInputRef.current) {
            fileInputRef.current.click();
        }
    };

    const onUploadImage = (e: ChangeEvent<HTMLInputElement>) => {
        const file = e.target.files?.[0];
        if (!file) {
            return;
        }
        appStore.uploadImage(uploadId.current, file);
    };

    return (
        <DefaultPage
            title="Applications"
            rightControl={
                <Button
                    id="create-app"
                    variant="contained"
                    color="primary"
                    onClick={() => setCreateDialog(true)}>
                    Create Application
                </Button>
            }
            maxWidth={1000}>
            <Grid size={12}>
                <Paper elevation={6} style={{overflowX: 'auto'}}>
                    <Table id="app-table">
                        <TableHead>
                            <TableRow>
                                <TableCell padding="checkbox" style={{width: 80}} />
                                <TableCell>Name</TableCell>
                                <TableCell>Token</TableCell>
                                <TableCell>Description</TableCell>
                                <TableCell>Priority</TableCell>
                                <TableCell>Last Used</TableCell>
                                <TableCell />
                                <TableCell />
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {apps.map((app: IApplication) => (
                                <Row
                                    key={app.id}
                                    description={app.description}
                                    defaultPriority={app.defaultPriority}
                                    image={app.image}
                                    name={app.name}
                                    value={app.token}
                                    lastUsed={app.lastUsed}
                                    fUpload={() => handleImageUploadClick(app.id)}
                                    fDeleteImage={() => setToDeleteImage(app)}
                                    fDelete={() => setToDeleteApp(app)}
                                    fEdit={() => setToUpdateApp(app)}
                                    noDelete={app.internal}
                                />
                            ))}
                        </TableBody>
                    </Table>
                    <input
                        ref={fileInputRef}
                        type="file"
                        accept={validExtensions.join(',')}
                        style={{display: 'none'}}
                        onChange={onUploadImage}
                    />
                </Paper>
            </Grid>
            {createDialog && (
                <AddApplicationDialog
                    fClose={() => setCreateDialog(false)}
                    fOnSubmit={appStore.create}
                />
            )}
            {toUpdateApp != null && (
                <UpdateApplicationDialog
                    fClose={() => setToUpdateApp(undefined)}
                    fOnSubmit={(name, description, defaultPriority) =>
                        appStore.update(toUpdateApp.id, name, description, defaultPriority)
                    }
                    initialDescription={toUpdateApp?.description}
                    initialName={toUpdateApp?.name}
                    initialDefaultPriority={toUpdateApp?.defaultPriority}
                />
            )}
            {toDeleteApp != null && (
                <ConfirmDialog
                    title="Confirm Delete"
                    text={'Delete ' + toDeleteApp.name + '?'}
                    fClose={() => setToDeleteApp(undefined)}
                    fOnSubmit={() => appStore.remove(toDeleteApp.id)}
                />
            )}
            {toDeleteImage != null && (
                <ConfirmDialog
                    title="Confirm Delete Image"
                    text={'Delete image for ' + toDeleteImage.name + '?'}
                    fClose={() => setToDeleteImage(undefined)}
                    fOnSubmit={() => appStore.deleteImage(toDeleteImage.id)}
                />
            )}
        </DefaultPage>
    );
});

interface IRowProps {
    name: string;
    value: string;
    noDelete: boolean;
    description: string;
    defaultPriority: number;
    lastUsed: string | null;
    fUpload: VoidFunction;
    fDeleteImage: VoidFunction;
    image: string;
    fDelete: VoidFunction;
    fEdit: VoidFunction;
}

const Row = ({
    name,
    value,
    noDelete,
    description,
    defaultPriority,
    lastUsed,
    fDelete,
    fUpload,
    fDeleteImage,
    image,
    fEdit,
}: IRowProps) => {
    const {classes} = useStyles();
    return (
        <TableRow>
            <TableCell padding="normal">
                <div style={{display: 'flex'}}>
                    <Tooltip title="Delete image" placement="top" arrow>
                        <ButtonBase className={classes.imageContainer} onClick={fDeleteImage}>
                            <img
                                src={config.get('url') + image}
                                alt="app logo"
                                width="40"
                                height="40"
                            />
                        </ButtonBase>
                    </Tooltip>
                    <IconButton onClick={fUpload} style={{height: 40}}>
                        <CloudUpload />
                    </IconButton>
                </div>
            </TableCell>
            <TableCell>{name}</TableCell>
            <TableCell>
                <CopyableSecret value={value} style={{display: 'flex', alignItems: 'center'}} />
            </TableCell>
            <TableCell>{description}</TableCell>
            <TableCell>{defaultPriority}</TableCell>
            <TableCell>
                <LastUsedCell lastUsed={lastUsed} />
            </TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fEdit} className="edit">
                    <Edit />
                </IconButton>
            </TableCell>
            <TableCell align="right" padding="none">
                <IconButton onClick={fDelete} className="delete" disabled={noDelete}>
                    <Delete />
                </IconButton>
            </TableCell>
        </TableRow>
    );
};

export default Applications;
