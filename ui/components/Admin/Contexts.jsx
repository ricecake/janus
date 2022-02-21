import React from 'react';
import Grid from '@material-ui/core/Grid';

import AdminPage from './frame';
import { AutoEditForm } from 'Component/Helpers';

import ContextList from './ContextList';

const Contexts = () => {
	return (
		<AdminPage>
			<Grid item>
				<ContextList />
			</Grid>
			{/* <Grid item>
				<AutoEditForm
					onSubmit={console.log}
					title="Application Groups"
					initialValues={{
						code: 'Foobar',
					}}
					fields={['name', 'description']}
					labels={{
						name: 'Group Name',
						description: 'Description',
					}}
				/>
			</Grid>
			<Grid item>Example</Grid> */}
		</AdminPage>
	);
};
export default Contexts;
