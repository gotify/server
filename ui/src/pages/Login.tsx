import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import TextField from '@material-ui/core/TextField';
import React, {ChangeEvent, Component, FormEvent} from 'react';
import * as UserAction from '../actions/UserAction';
import Container from '../component/Container';
import DefaultPage from '../component/DefaultPage';

interface IState {
    username: string;
    password: string;
}

class Login extends Component<{}, IState> {
    public state = {username: '', password: ''};

    public render() {
        const {username, password} = this.state;
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
                                onChange={this.handleChange.bind(this, 'username')}
                            />
                            <TextField
                                type="password"
                                className="password"
                                label="Password"
                                margin="normal"
                                value={password}
                                onChange={this.handleChange.bind(this, 'password')}
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

    private handleChange(propertyName: string, event: ChangeEvent<HTMLInputElement>) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    private login = (e: React.MouseEvent<HTMLInputElement>) => {
        e.preventDefault();
        UserAction.login(this.state.username, this.state.password);
    };

    private preventDefault = (e: FormEvent<HTMLFormElement>) => e.preventDefault();
}

export default Login;
