import PropTypes from 'prop-types';
import React, { useEffect } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import { ButtonGroup } from '@material-ui/core';
import Snackbar from '@material-ui/core/Snackbar';
import MuiAlert from '@material-ui/lab/Alert';
import { Show, Hide } from 'Component/Helpers';
import LinearProgress from '@material-ui/core/LinearProgress';
import Grid from '@material-ui/core/Grid';

import {
	fetchUserDetails,
	updateUserDetails,
	initiatePasswordChange,
} from 'Include/reducers/profile';

const PasswordBase = ({ loading }) => {
	const [pass, setPass] = React.useState('');
	const [verify, setVerify] = React.useState('');

	return (
		<Card>
			<CardHeader title="Change Password" />
			<CardContent>
				<form
					onSubmit={(e) => {
						e.preventDefault();
						if (verify && verify === pass) {
							initiatePasswordChange(pass, verify);
						}
					}}
				>
					<Show If={loading}>
						<LinearProgress />
					</Show>
					<Grid item>
						<TextField
							required
							fullWidth
							autoFocus
							disabled={loading}
							name="password"
							label="Password"
							type="password"
							variant="outlined"
							margin="normal"
							autoComplete="new-password"
							onChange={(e) => setPass(e.target.value)}
							error={!!pass && pass.length < 8}
						/>
					</Grid>
					<Grid item>
						<TextField
							required
							fullWidth
							disabled={loading}
							name="verify_password"
							label="Verify Password"
							type="password"
							variant="outlined"
							margin="normal"
							autoComplete="new-password"
							onChange={(e) => setVerify(e.target.value)}
							error={!!verify && verify !== pass}
							helperText={
								!!verify && !!pass && verify !== pass
									? "Passwords don't seem to match..."
									: ''
							}
						/>
					</Grid>
					<Grid
						container
						direction="row"
						justifyContent="space-around"
						alignItems="center"
					>
						<Button
							disabled={loading || (!!verify && verify !== pass)}
							type="submit"
							variant="contained"
							color="primary"
						>
							Finish
						</Button>
					</Grid>
				</form>
			</CardContent>
		</Card>
	);
};

const stateToProps = ({ profile }) => ({ ...profile });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{ fetchUserDetails, updateUserDetails, initiatePasswordChange },
		dispatch
	);

export const Password = connect(stateToProps, dispatchToProps)(PasswordBase);
