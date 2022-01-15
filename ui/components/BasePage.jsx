import React from 'react';
import { OidcProvider } from 'redux-oidc';

import { ThemeProvider, Grid } from '@material-ui/core/';
import { createTheme } from '@material-ui/core/styles';
import green from '@material-ui/core/colors/green';
import teal from '@material-ui/core/colors/teal';

import store from 'Include/store';
import userManager, { withLogin } from 'Include/userManager';

const theme = createTheme({
	palette: {
		primary: green,
		secondary: teal,
		background: {
			default:
				'linear-gradient(0deg, rgba(0,203,0,1) 0%, rgba(0,128,128,1) 100%)',
		},
	},
	status: {
		danger: 'orange',
	},
});

const BasePage = (props) => (
	<ThemeProvider theme={theme}>
		<Grid container direction="row" justify="center" alignItems="center">
			{props.children}
		</Grid>
	</ThemeProvider>
);

export const LoginBasePage = (props) => {
	let WrappedBase = withLogin(BasePage);

	return (
		<OidcProvider store={store} userManager={userManager}>
			<WrappedBase {...props} />
		</OidcProvider>
	);
};

export default BasePage;
