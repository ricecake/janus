import React, { useEffect } from 'react';
import BasePage from 'Component/BasePage';
import Grid from '@material-ui/core/Grid';
import userManager from 'Include/userManager';

export const LoginPage = (props) => {
	useEffect(() => {
		userManager.removeUser();
		fetch('/sessions', {
			method: 'DELETE',
		});
	});
	return (
		<React.Fragment>
			<BasePage>
				<Grid
					container
					direction="row"
					justifyContent="space-evenly"
					alignItems="center"
				>
					You are now logged out.
				</Grid>
			</BasePage>
		</React.Fragment>
	);
};
export default LoginPage;
