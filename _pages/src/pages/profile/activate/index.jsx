/* eslint-disable import/no-extraneous-dependencies */
import React from "react";
import ReactDOM from "react-dom";
import HelloWorld from 'Component/HelloWorld';
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

const root = document.getElementById('main');

ReactDOM.render((
	<ThemeProvider theme={theme}>
		<Grid container direction="row" justify="center" alignItems="center">
			<HelloWorld />
		</Grid>
	</ThemeProvider>
), root);
