import React from "react";
import { ThemeProvider, Grid } from '@material-ui/core/';
import { createMuiTheme } from '@material-ui/core/styles';
import green from '@material-ui/core/colors/green';
import teal from '@material-ui/core/colors/teal';

const theme = createMuiTheme({
	palette: {
	  primary: green,
	  secondary: teal,
	},
	status: {
	  danger: 'orange',
	},
  });

const BasePage = props => (
	<ThemeProvider theme={theme}>
		<Grid container direction="row" justify="center" alignItems="center">
			{props.children}
		</Grid>
	</ThemeProvider>
);

export default BasePage;