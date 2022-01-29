import PropTypes from 'prop-types';
import React from 'react';
import { Dashboard } from 'Component/Dashboard';

export const ProfilePage = (props) => (
	<Dashboard
		root="/profile"
		title={'Profile Management'}
		categories={[
			{ id: 'Logins' }, // this should be a list of browser logins, with individual app contexts listed underneath
			{ id: 'Authentication' },
		]}
	>
		{props.children}
	</Dashboard>
);

ProfilePage.propTypes = {
	children: PropTypes.any,
};
export default ProfilePage;
