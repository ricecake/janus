import React from 'react';
import { Dashboard } from 'Component/Dashboard';

export const ProfilePage = (props) => (
	<Dashboard
		root="/profile"
		title={'Profile Management'}
		categories={[
			{ id: 'Sessions' },
			{ id: 'Logins' },
			{ id: 'Authentication' },
		]}
	>
		{props.children}
	</Dashboard>
);
export default ProfilePage;
