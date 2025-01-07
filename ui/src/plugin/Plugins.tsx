import React, {useEffect} from 'react';
import {Link} from 'react-router-dom';
import Grid from '@mui/material/Grid2';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Settings from '@mui/icons-material/Settings';
import {Switch, Button} from '@mui/material';
import DefaultPage from '../common/DefaultPage';
import CopyableSecret from '../common/CopyableSecret';
import {useAppDispatch, useAppSelector} from '../store';
import {changePluginEnableState, fetchPlugins} from '../plugin/plugin-actions.ts';
import {IPlugin} from '../types';

const Plugins = () => {
    const dispatch = useAppDispatch();
    const plugins = useAppSelector((state) => state.plugin.items);

    useEffect(() => {
        dispatch(fetchPlugins());
    }, [dispatch]);

    const handleChangePluginStatus = async (plugin: IPlugin) => {
        await dispatch(
            changePluginEnableState(plugin.id, !plugin.enabled)
        )
    }

    return (
        <DefaultPage title="Plugins" maxWidth={1000}>
            <Grid size={12}>
                    <Paper elevation={6} style={{overflowX: 'auto'}}>
                    <Table id="plugin-table">
                        <TableHead>
                            <TableRow>
                                <TableCell>ID</TableCell>
                                <TableCell>Enabled</TableCell>
                                <TableCell>Name</TableCell>
                                <TableCell>Token</TableCell>
                                <TableCell>Details</TableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {plugins.map((plugin: IPlugin) => (
                                <Row
                                    key={plugin.token}
                                    id={plugin.id}
                                    token={plugin.token}
                                    name={plugin.name}
                                    enabled={plugin.enabled}
                                    fToggleStatus={() => handleChangePluginStatus(plugin)}
                                />
                            ))}
                        </TableBody>
                    </Table>
                </Paper>
            </Grid>
        </DefaultPage>
    );
};

interface IRowProps {
    id: number;
    name: string;
    token: string;
    enabled: boolean;
    fToggleStatus: VoidFunction;
}

const Row = ({name, id, token, enabled, fToggleStatus}: IRowProps) => (
    <TableRow>
        <TableCell>{id}</TableCell>
        <TableCell>
            <Switch
                checked={enabled}
                onClick={fToggleStatus}
                className="switch"
                data-enabled={enabled}
            />
        </TableCell>
        <TableCell>{name}</TableCell>
        <TableCell>
            <CopyableSecret value={token} style={{display: 'flex', alignItems: 'center'}} />
        </TableCell>
        <TableCell align="right" padding="none">
            <Link to={'/plugins/' + id}>
                <Button>
                    <Settings />
                </Button>
            </Link>
        </TableCell>
    </TableRow>
);

export default Plugins;
