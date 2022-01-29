import React from 'react';
import Grid from '@material-ui/core/Grid';

import ProfilePage from './frame';
import { ProfileDetails } from './Details';

const ProfileOverview = () => {
	return (
		<ProfilePage>
			<Grid item>
				<ProfileDetails />
			</Grid>
		</ProfilePage>
	);
};
export default ProfileOverview;
