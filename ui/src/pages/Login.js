import React, {Component} from 'react';
import Button from 'material-ui/Button';
import Grid from 'material-ui/Grid';
import TextField from 'material-ui/TextField';
import Container from '../component/Container';
import * as UserAction from '../actions/UserAction';
import DefaultPage from '../component/DefaultPage';

class Login extends Component {
    constructor() {
        super();
        this.state = {username: '', password: ''};
    }

    handleChange(propertyName, event) {
        const state = this.state;
        state[propertyName] = event.target.value;
        this.setState(state);
    }

    login = () => UserAction.login(this.state.username, this.state.password);

    render() {
        const {username, password} = this.state;
        return (
            <DefaultPage title="Login" maxWidth={250} hideButton={true}>
                <Grid item xs={12} style={{textAlign: 'center'}}>
                    <Container>
                        <form onSubmit={(e) => e.preventDefault()}>
                            <TextField id="name" label="Username" margin="dense" value={username}
                                       onChange={this.handleChange.bind(this, 'username')}/>
                            <TextField type="password" id="password" label="Password" margin="normal"
                                       value={password} onChange={this.handleChange.bind(this, 'password')}/>
                            <Button type="submit" variant="raised" size="large" color="primary"
                                    style={{marginTop: 15, marginBottom: 5}} onClick={this.login}>
                                Login
                            </Button>
                        </form>
                    </Container>
                </Grid>
            </DefaultPage>
        );
    }
}

export default Login;
