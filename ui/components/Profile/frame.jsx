import PropTypes from 'prop-types';
import React from 'react';
import { Dashboard } from 'Component/Dashboard';
import Grid from '@material-ui/core/Grid';

export const ProfilePage = (props) => (
	<Dashboard
		root="/profile"
		title={'Profile Management'}
		categories={[
			{ id: 'Logins' }, // this should be a list of browser logins, with individual app contexts listed underneath
			{ id: 'Authentication' },
		]}
	>
		<Grid
			container
			spacing={2}
			direction="row"
			justifyContent="space-evenly"
			alignItems="baseline"
		>
			{props.children}
		</Grid>
	</Dashboard>
);

ProfilePage.propTypes = {
	children: PropTypes.any,
};
export default ProfilePage;
