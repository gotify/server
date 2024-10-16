import React, {useState} from 'react';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid2';
import TextField from '@mui/material/TextField';
import {useNavigate} from 'react-router';
import Container from '../common/Container';
import DefaultPage from '../common/DefaultPage';
import * as config from '../config';
import {useAppDispatch, useAppSelector} from '../store';
import {login} from '../store/auth-actions.ts';
import RegistrationDialog from './Register';

const Login = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const [showRegisterDialog, setShowRegisterDialog] = useState(false);
    const connectionErrorMessage = useAppSelector((state) => state.ui.connectionErrorMessage);

    const handleLogin = (event) => {
        event.preventDefault();
        const fd = new FormData(event.target);
        const loginData = Object.fromEntries(fd.entries());
        dispatch(login(loginData.username, loginData.password))
            .then(() => navigate('/'))
            .catch((error) => {
                console.log(error);
                event.target.reset();
            });
    };

    const registerButton = () => {
        if (config.get('register'))
            return (
                <Button
                    id="register"
                    variant="contained"
                    color="primary"
                    onClick={() => setShowRegisterDialog(true)}>
                    Register
                </Button>
            );
        else return null;
    };

    return (
        <DefaultPage title="Login" rightControl={registerButton()} maxWidth={250}>
            <Grid size={12} style={{textAlign: 'center'}}>
                <Container>
                    <form onSubmit={handleLogin} id="login-form">
                        <TextField
                            name="username"
                            variant="standard"
                            autoFocus
                            className="name"
                            label="Username"
                            margin="dense"
                            autoComplete="username"
                        />
                        <TextField
                            name="password"
                            variant="standard"
                            type="password"
                            className="password"
                            label="Password"
                            margin="normal"
                            autoComplete="current-password"
                        />
                        <Button
                            type="submit"
                            variant="contained"
                            size="large"
                            className="login"
                            color="primary"
                            disabled={!!connectionErrorMessage}
                            style={{marginTop: 15, marginBottom: 5}}>
                            Login
                        </Button>
                    </form>
                </Container>
            </Grid>
            {showRegisterDialog && (
                <RegistrationDialog
                    fClose={() => setShowRegisterDialog(false)}
                    fOnSubmit={this.props.currentUser.register}
                />
            )}
        </DefaultPage>
    );
};

export default Login;
