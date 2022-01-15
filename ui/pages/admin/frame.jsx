import React from 'react';
import { Dashboard } from 'Component/Dashboard';

export const AdminPage = (props) => (
	<Dashboard
		root="/admin"
		title={'System Management'}
		categories={[
			{ id: 'Users' },
			{ id: 'Application Groups', path: 'contexts' },
			{ id: 'Clients' },
			{ id: 'Roles' },
			{ id: 'Actions' },
		]}
	>
		{props.children}
	</Dashboard>
);
export default AdminPage;
