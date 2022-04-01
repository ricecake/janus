import React from 'react';
import Grid from '@material-ui/core/Grid';

import AdminPage from './frame';
import { ClientList } from './ClientList';

const Clients = () => {
	return (
		<AdminPage>
			<Grid item xs={12}>
				<ClientList />
			</Grid>
		</AdminPage>
	);
};
export default Clients;

// a list of Contexts, with netsted datagrid of clients with Edit, move context, delete and create
// maybe datagrid is overkill?  but also super easy...

/*
Can edit things like "moving" with singleSelect and valueOptions on the datagrid, and can use renderEditCell for more complicated notions, like adding groups to a user.
*/
