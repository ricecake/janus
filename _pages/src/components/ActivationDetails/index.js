import React, { PureComponent } from "react";
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import { Grid } from '@material-ui/core';
import { connect } from "react-redux";
import { changeName, changePassword, changePasswordVerifier, submitForm, startSignin } from "Include/reducers/activation";
import { bindActionCreators } from 'redux'

class ActivationDetails extends PureComponent {
	constructor(props) {
		super(props);

		console.log(props);
	}

	render(props) {
		if (!this.props.user) {
			this.props.startSignin();
			return null;
		}
		return (
			<div>
				<h1>Activate User</h1>
				<Grid
					container
					direction="column"
					justify="center"
					alignItems="center"
				>
					<TextField
						required
						name='preferred_name'
						label="What should we call you?"
						type="text"
						variant="outlined"
						margin="normal"
						error={!this.props.name_valid}
						helperText={this.props.name_valid?'':"We have to call you something!"}
						onChange={e => this.props.changeName(e.target.value)}
						value={ this.props.preferred_name }
					/>
					<TextField
						required
						name='password'
						label="Password"
						type="password"
						variant="outlined"
						margin="normal"
						error={!this.props.password_valid}
						helperText="Password must be at least eight characters long"
						onChange={e => this.props.changePassword(e.target.value)}
						value={ this.props.password }
					/>
					<TextField
						required
						name='verify_password'
						label="Verify Password"
						type="password"
						variant="outlined"
						margin="normal"
						error={!this.props.password_match}
						helperText={this.props.password_match? '':"Passwords don't seem to match..."}
						onChange={e => this.props.changePasswordVerifier(e.target.value) }
						value={ this.props.verify_password }
					/>

					<Button variant="contained" color="primary" onClick={ this.props.submitForm } disabled={!this.props.submitable}>
						Activate User
					</Button>
				</Grid>
			</div>
		);
	}
}

const stateToProps = ({activation, oidc}) => ({...activation, user: oidc.user });
const dispatchToProps = (dispatch) => bindActionCreators({
	changeName, changePassword, changePasswordVerifier, submitForm, startSignin
}, dispatch);

export default connect(stateToProps, dispatchToProps)(ActivationDetails);