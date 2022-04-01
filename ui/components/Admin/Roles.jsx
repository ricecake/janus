import React from 'react';
import Grid from '@material-ui/core/Grid';

import AdminPage from './frame';

const Roles = () => {
	return (
		<AdminPage>
			<Grid item>Example</Grid>
			<Grid item>Example</Grid>
		</AdminPage>
	);
};
export default Roles;

// contexts with nested Roles, and create/update/delete, so datagrid
// need a way to have a mapping of actions to roles that's editable
// Need a multiselect for role actions -- mui multi select
