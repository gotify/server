import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import React, {Component, FormEvent} from 'react';
import Container from '../common/Container';
import DefaultPage from '../common/DefaultPage';
import {observable} from 'mobx';
import {observer} from 'mobx-react';
import {inject, Stores} from '../inject';
import * as config from '../config';
import RegistrationDialog from './Register';

@observer
class Login extends Component<Stores<'currentUser'>> {
    @observable
    private username = '';
    @observable
    private password = '';
    @observable
    private registerDialog = false;

    public render() {
        const {username, password, registerDialog} = this;
        return (
            <DefaultPage title="Login" rightControl={this.registerButton()} maxWidth={250}>
                <Grid item xs={12} style={{textAlign: 'center'}}>
                    <Container>
                        <form onSubmit={this.preventDefault} id="login-form">
                            <TextField
                                autoFocus
                                className="name"
                                label="Username"
                                margin="dense"
                                autoComplete="username"
                                value={username}
                                onChange={(e) => (this.username = e.target.value)}
                            />
                            <TextField
                                type="password"
                                className="password"
                                label="Password"
                                margin="normal"
                                autoComplete="current-password"
                                value={password}
                                onChange={(e) => (this.password = e.target.value)}
                            />
                            <Button
                                type="submit"
                                variant="contained"
                                size="large"
                                className="login"
                                color="primary"
                                disabled={!!this.props.currentUser.connectionErrorMessage}
                                style={{marginTop: 15, marginBottom: 5}}
                                onClick={this.login}>
                                Login
                            </Button>
                        </form>
                    </Container>
                </Grid>
                {registerDialog && (
                    <RegistrationDialog
                        fClose={() => (this.registerDialog = false)}
                        fOnSubmit={this.props.currentUser.register}
                    />
                )}
            </DefaultPage>
        );
    }

    private login = (e: React.MouseEvent<HTMLButtonElement>) => {
        e.preventDefault();
        this.props.currentUser.login(this.username, this.password);
    };

    private registerButton = () => {
        if (config.get('register'))
            return (
                <Button
                    id="register"
                    variant="contained"
                    color="primary"
                    onClick={() => (this.registerDialog = true)}>
                    Register
                </Button>
            );
        else return null;
    };

    private preventDefault = (e: FormEvent<HTMLFormElement>) => e.preventDefault();
}

export default inject('currentUser')(Login);
