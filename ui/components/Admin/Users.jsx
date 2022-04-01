import React from 'react';
import Grid from '@material-ui/core/Grid';

import AdminPage from './frame';

const Users = () => {
	return (
		<AdminPage>
			<Grid item>Example</Grid>
			<Grid item>Example</Grid>
		</AdminPage>
	);
};
export default Users;

// A datagrid of Users.
// need a way to have the ability to assign roles, view roles, and likewise groups
// mui multiselect would be good for those.
// should be able to have a multi select with checkboxes and chips, and apply groupings to keep the context roles straight.
//  Since it will get an api reference, should be able to make it a bit easier to control how we pass the args so things dont get weird on update with duplicate role names in different contexts.
