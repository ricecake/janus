import React, { PureComponent } from "react";
import Button from '@material-ui/core/Button';
import TextField from '@material-ui/core/TextField';
import { Grid } from '@material-ui/core';
import { connect } from "react-redux";
import { changeName, changePassword, changePasswordVerifier, submitForm } from "Include/reducers/activation";
import { bindActionCreators } from 'redux'

class ActivationDetails extends PureComponent {
	constructor(props) {
		super(props);

		console.log(props);
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
						// error={this.state.preferred_name.length <= 0}
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
						// error={this.state.password.length <= 0}
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
						// error={this.props.state.verify_password.length > 0 && this.props.state.verify_password != this.props.state.password }
						onChange={e => this.props.changePasswordVerifier(e.target.value) }
						value={ this.props.verify_password }
					/>

					<Button variant="contained" color="primary" onClick={ this.props.submitForm }>
						Activate User
					</Button>
				</Grid>
			</div>
		);
	}
}

const stateToProps = ({activation}) => activation;
const dispatchToProps = (dispatch) => bindActionCreators({
	changeName, changePassword, changePasswordVerifier, submitForm
}, dispatch);

export default connect(stateToProps, dispatchToProps)(ActivationDetails);