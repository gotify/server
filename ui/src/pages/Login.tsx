import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import React, {Component, FormEvent} from 'react';
import Container from '../component/Container';
import DefaultPage from '../component/DefaultPage';
import {currentUser} from '../stores/CurrentUser';
import {observable} from 'mobx';
import {observer} from 'mobx-react';

@observer
class Login extends Component {
    @observable
    private username = '';
    @observable
    private password = '';

    public render() {
        const {username, password} = this;
        return (
            <DefaultPage title="Login" maxWidth={250} hideButton={true}>
                <Grid item xs={12} style={{textAlign: 'center'}}>
                    <Container>
                        <form onSubmit={this.preventDefault} id="login-form">
                            <TextField
                                autoFocus
                                className="name"
                                label="Username"
                                margin="dense"
                                value={username}
                                onChange={(e) => (this.username = e.target.value)}
                            />
                            <TextField
                                type="password"
                                className="password"
                                label="Password"
                                margin="normal"
                                value={password}
                                onChange={(e) => (this.password = e.target.value)}
                            />
                            <Button
                                type="submit"
                                variant="raised"
                                size="large"
                                className="login"
                                color="primary"
                                style={{marginTop: 15, marginBottom: 5}}
                                onClick={this.login}>
                                Login
                            </Button>
                        </form>
                    </Container>
                </Grid>
            </DefaultPage>
        );
    }

    private login = (e: React.MouseEvent<HTMLInputElement>) => {
        e.preventDefault();
        currentUser.login(this.username, this.password);
    };

    private preventDefault = (e: FormEvent<HTMLFormElement>) => e.preventDefault();
}

export default Login;
