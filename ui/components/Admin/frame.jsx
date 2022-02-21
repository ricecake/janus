import PropTypes from 'prop-types';
import React from 'react';
import { Dashboard } from 'Component/Dashboard';
import Grid from '@material-ui/core/Grid';
import { withLogin } from 'Include/userManager';

export const AdminPage = withLogin((props) => (
	<Dashboard
		root="/admin"
		title={'System Management'}
		categories={[
			{ id: 'Users' },
			{ id: 'Application Groups', path: 'contexts' },
			{ id: 'Clients' },
			{ id: 'Roles' },
			{ id: 'Actions' },
			{ id: 'Groups' },
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
));
export default AdminPage;
