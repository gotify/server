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

const Applications = observer(() => {
    const {appStore} = useStores();
    const apps = appStore.getItems();
    const [toDeleteApp, setToDeleteApp] = useState<IApplication>();
    const [toUpdateApp, setToUpdateApp] = useState<IApplication>();
    const [createDialog, setCreateDialog] = useState<boolean>(false);

    const fileInputRef = useRef<HTMLInputElement>(null);
    const uploadId = useRef(-1);

    useEffect(() => void appStore.refresh(), []);

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
        if (['image/png', 'image/jpeg', 'image/gif'].indexOf(file.type) !== -1) {
            appStore.uploadImage(uploadId.current, file);
        } else {
            alert('Uploaded file must be of type png, jpeg or gif.');
        }
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
                                <TableCell>Sort Order</TableCell>
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
                                    sortOrder={app.sortOrder}
                                    image={app.image}
                                    name={app.name}
                                    value={app.token}
                                    lastUsed={app.lastUsed}
                                    fUpload={() => handleImageUploadClick(app.id)}
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
                    fOnSubmit={(name, description, defaultPriority, sortOrder) =>
                        appStore.update(toUpdateApp.id, name, description, defaultPriority, sortOrder)
                    }
                    initialDescription={toUpdateApp?.description}
                    initialName={toUpdateApp?.name}
                    initialDefaultPriority={toUpdateApp?.defaultPriority}
                    initialSortOrder={appStore.getByID(toUpdateApp.id).sortOrder}
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
        </DefaultPage>
    );
});

interface IRowProps {
    name: string;
    value: string;
    noDelete: boolean;
    description: string;
    defaultPriority: number;
    sortOrder: number;
    lastUsed: string | null;
    fUpload: VoidFunction;
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
    sortOrder,
    lastUsed,
    fDelete,
    fUpload,
    image,
    fEdit,
}: IRowProps) => {
    return (
        <TableRow>
            <TableCell padding="normal">
                <div style={{display: 'flex'}}>
                    <img src={config.get('url') + image} alt="app logo" width="40" height="40" />
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
            <TableCell>{sortOrder}</TableCell>
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
