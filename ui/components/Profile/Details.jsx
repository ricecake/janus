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

import { fetchUserDetails, updateUserDetails } from 'Include/reducers/profile';

const ProfileDetailsBase = ({
	fetchUserDetails,
	updateUserDetails,
	loading,
	user_details,
}) => {
	const [PreferredName, setPreferredName] = React.useState(
		user_details.PreferredName
	);
	const [GivenName, setGivenName] = React.useState(user_details.GivenName);
	const [FamilyName, setFamilyName] = React.useState(user_details.FamilyName);

	useEffect(() => {
		fetchUserDetails();
	}, []);

	useEffect(()=>{
		setPreferredName(user_details.PreferredName);
		setGivenName(user_details.GivenName);
		setFamilyName(user_details.FamilyName);
	}, [user_details]);

	const [profileView, setProfileView] = React.useState(true);
	const [open, setOpen] = React.useState(false);

	const handleClick = () => {
		setOpen(true);
	};

	const handleClose = (event, reason) => {
		if (reason === 'clickaway') {
			return;
		}

		setOpen(false);
	};
	return (
		<Card>
			<CardHeader title="Profile Details" />
			<CardContent>
				<Hide If={loading}>
					<fieldset disabled={profileView}>
						<TextField
							defaultValue={user_details.PreferredName}
							onChange={(e) => setPreferredName(e.target.value)}
							margin="normal"
							variant="outlined"
							label="Preferred Name"
							name="preferred_name"
							autoComplete="nickname"
						/>
						<TextField
							defaultValue={user_details.GivenName}
							onChange={(e) => setGivenName(e.target.value)}
							margin="normal"
							variant="outlined"
							label="Given Name"
							name="given_name"
							autoComplete="given-name"
						/>
						<TextField
							defaultValue={user_details.FamilyName}
							onChange={(e) => setFamilyName(e.target.value)}
							margin="normal"
							variant="outlined"
							label="Family Name"
							name="family_name"
							autoComplete="family-name"
						/>
					</fieldset>
				</Hide>
				<Show If={profileView && !loading}>
					<Button onClick={() => setProfileView(false)}>Edit</Button>
				</Show>
				<Hide If={profileView || loading}>
					<ButtonGroup>
						<Button onClick={() => setProfileView(true)}>
							Cancel
						</Button>
						<Button
							onClick={() => {
								setProfileView(true);
								updateUserDetails({
									PreferredName,
									GivenName,
									FamilyName,
								}).then(() => {
									handleClick();
								});
							}}
						>
							Save
						</Button>
					</ButtonGroup>
				</Hide>
				<Show If={loading}>
					<LinearProgress />
				</Show>
				<Snackbar
					open={open}
					autoHideDuration={5000}
					onClose={handleClose}
				>
					<MuiAlert onClose={handleClose} severity="success">
						Profile updated
					</MuiAlert>
				</Snackbar>
			</CardContent>
		</Card>
	);
};

ProfileDetailsBase.propTypes = {
	fetchUserDetails: PropTypes.func,
	loading: PropTypes.any,
	user_details: PropTypes.shape({
		FamilyName: PropTypes.any,
		GivenName: PropTypes.any,
		PreferredName: PropTypes.any,
	}),
};

const stateToProps = ({ profile }) => ({ ...profile });
const dispatchToProps = (dispatch) =>
	bindActionCreators({ fetchUserDetails, updateUserDetails }, dispatch);

export const ProfileDetails = connect(
	stateToProps,
	dispatchToProps
)(ProfileDetailsBase);
