import React, {Component} from 'react';
import Button from 'material-ui/Button';
import Grid from 'material-ui/Grid';
import TextField from 'material-ui/TextField';
import Typography from 'material-ui/Typography';
import Container from '../component/Container';
import * as UserAction from '../actions/UserAction';
import DefaultPage from '../component/DefaultPage';
import PropTypes from 'prop-types';

class Login extends Component {
    static propTypes = {
        loginFailed: PropTypes.bool.isRequired,
    };

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
        const {loginFailed} = this.props;
        return (
            <DefaultPage title="Login" maxWidth={250} hideButton={true}>
                <Grid item xs={12} style={{textAlign: 'center'}}>
                    <Container>
                        <form>
                            <TextField id="name" label="Username" margin="dense" value={username}
                                       onChange={this.handleChange.bind(this, 'username')}/>
                            <TextField type="password" id="password" label="Password" margin="normal"
                                       value={password} onChange={this.handleChange.bind(this, 'password')}/>
                            <Button variant="raised" size="large" color="primary"
                                    style={{marginTop: 15, marginBottom: 5}} onClick={this.login}>
                                Login
                            </Button>
                            {loginFailed && <Typography>Login Failed</Typography>}
                        </form>
                    </Container>
                </Grid>
            </DefaultPage>
        );
    }
}

export default Login;
