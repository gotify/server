import React, {useRef, useState} from 'react';
import Button from '@mui/material/Button';
import Grid from '@mui/material/Grid2';
import TextField from '@mui/material/TextField';
import {useNavigate} from 'react-router';
import Container from '../common/Container';
import DefaultPage from '../common/DefaultPage';
import * as config from '../config';
import {useAppDispatch, useAppSelector} from '../store';
import {login, register} from './auth-actions.ts';
import RegistrationDialog from './Register';

const Login = () => {
    const dispatch = useAppDispatch();
    const navigate = useNavigate();
    const formRef = useRef<HTMLFormElement>(null);
    const [showRegisterDialog, setShowRegisterDialog] = useState(false);
    const connectionErrorMessage = useAppSelector((state) => state.ui.connectionErrorMessage);

    const handleLogin = (formData: FormData) => {
        const username = formData.get('username') as string;
        const password = formData.get('password') as string;
        dispatch(login(username, password))
            .then(() => navigate('/'))
            .catch((error) => {
                console.log(error);
                formRef.current!.reset();
            });
    };

    const handleRegister = async (username: string, password: string) => {
        return dispatch(register(username, password))
            .then(() => {
                navigate('/');
                return true;
            })
            .catch((error) => {
                console.error(error);
                return false;
            });
    }

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
                    <form action={handleLogin} ref={formRef} id="login-form">
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
                    fOnSubmit={handleRegister}
                />
            )}
        </DefaultPage>
    );
};

export default Login;
