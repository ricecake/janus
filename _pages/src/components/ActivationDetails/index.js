import React from "react";
import { Component } from 'react';
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import { Grid } from '@material-ui/core';
import { connect } from "react-redux";
import { changeName, changePassword, changePasswordVerifier, submitForm } from "Include/reducers/activation";
import { bindActionCreators } from 'redux'

class ActivationDetails extends Component {
	static defaultProps = {
		preferred_name: '',
		password: '',
		verify_password: '',
	};

	constructor(props) {
		super(props);

		console.log(props);

		this.state = {
			... this.defaultProps,
			... props
		};
	}

	render(props) {
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
						label="Name"
						type="text"
						variant="outlined"
						margin="normal"
						error={this.state.preferred_name.length <= 0}
						onChange={ (e) => this.setState({ preferred_name: e.target.value }) }
						value={ this.state.preferred_name }
					/>
					<TextField
						required
						name='password'
						label="Password"
						type="password"
						variant="outlined"
						margin="normal"
						error={this.state.password.length <= 0}
						onChange={ (e) => this.setState({ password: e.target.value }) }
						value={ this.state.password }
					/>
					<TextField
						required
						name='verify_password'
						label="Verify Password"
						type="password"
						variant="outlined"
						margin="normal"
						error={this.state.verify_password.length > 0 && this.state.verify_password != this.state.password }
						onChange={ (e) => this.setState({ verify_password: e.target.value }) }
						value={ this.state.verify_password }
					/>

					<Button variant="contained" color="primary" onClick={()=>{
						fetch("/profile/api/activate", {
							method: 'POST',
							headers: {
								'Content-Type': 'application/json',
								'Authorization': `Bearer ${ this.state.access_token }`,
							},
							body: JSON.stringify({
								password: this.state.password,
								verify_password: this.state.verify_password,
								preferred_name: this.state.preferred_name,
							}),
						})
					}}>
						Activate User
					</Button>
				</Grid>
			</div>
		);
	}
}

const stateToProps = (state) => state;
const dispatchToProps = (dispatch) => bindActionCreators({
	changeName, changePassword, changePasswordVerifier, submitForm
}, dispatch);

export default connect(stateToProps, dispatchToProps)(ActivationDetails);