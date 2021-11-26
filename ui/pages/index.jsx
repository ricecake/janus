import React, { Suspense, lazy } from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';

import {
	createTheme,
	ThemeProvider,
	withStyles,
} from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

let theme = createTheme({
	palette: {
		primary: {
			light: '#63ccff',
			main: '#009be5',
			dark: '#006db3',
			background: '#18202c',
		},
	},
	typography: {
		h5: {
			fontWeight: 500,
			fontSize: 26,
			letterSpacing: 0.5,
		},
	},
	shape: {
		borderRadius: 8,
	},
	props: {
		MuiTab: {
			disableRipple: true,
		},
	},
	mixins: {
		toolbar: {
			minHeight: 48,
		},
	},
});

theme = {
	...theme,
	overrides: {
		MuiDrawer: {
			paper: {
				backgroundColor: '#18202c',
			},
		},
		MuiButton: {
			label: {
				textTransform: 'none',
			},
			contained: {
				boxShadow: 'none',
				'&:active': {
					boxShadow: 'none',
				},
			},
		},
		MuiTabs: {
			root: {
				marginLeft: theme.spacing(1),
			},
			indicator: {
				height: 3,
				borderTopLeftRadius: 3,
				borderTopRightRadius: 3,
				backgroundColor: theme.palette.common.white,
			},
		},
		MuiTab: {
			root: {
				textTransform: 'none',
				margin: '0 16px',
				minWidth: 0,
				padding: 0,
				[theme.breakpoints.up('md')]: {
					padding: 0,
					minWidth: 0,
				},
			},
		},
		MuiIconButton: {
			root: {
				padding: theme.spacing(1),
			},
		},
		MuiTooltip: {
			tooltip: {
				borderRadius: 4,
			},
		},
		MuiDivider: {
			root: {
				backgroundColor: '#404854',
			},
		},
		MuiListItemText: {
			primary: {
				fontWeight: theme.typography.fontWeightMedium,
			},
		},
		MuiListItemIcon: {
			root: {
				color: 'inherit',
				marginRight: 0,
				'& svg': {
					fontSize: 20,
				},
			},
		},
		MuiAvatar: {
			root: {
				width: 32,
				height: 32,
			},
		},
	},
};

const styles = {
	root: {
		display: 'flex',
		minHeight: '100vh',
	},
	app: {
		flex: 1,
		display: 'flex',
		flexDirection: 'column',
	},
	main: {
		flex: 1,
		padding: theme.spacing(6, 4),
		background: '#eaeff1',
	},
	footer: {
		padding: theme.spacing(2),
		background: '#eaeff1',
	},
};

const Home = (props) => (
	<React.Fragment>This is the default page!</React.Fragment>
);

const withSuspense = (element, fallback = <div>Loading...</div>) => (
	<Suspense fallback={fallback}>{element}</Suspense>
);

const Login = lazy(() => import('Page/login'));
const Signup = lazy(() => import('Page/signup'));
const Activate = lazy(() => import('Page/profile/activate'));
const OidcCallback = lazy(() => import('Page/callbacks/oidc'));

const Webauthn = lazy(() => import('Page/profile/webauthn'));

export const App = (props) => {
	const { classes } = props;

	return (
		<ThemeProvider theme={theme}>
			<CssBaseline />
			<div className={classes.root}>
				<div className={classes.app}>
					<Router>
						<Routes>
							<Route path="/">
								<Route index element={<Home />} />
								<Route
									path="login"
									element={withSuspense(<Login />)}
								/>
								<Route
									path="signup"
									element={withSuspense(<Signup />)}
								/>
								<Route path="profile">
									<Route
										path="activate"
										element={withSuspense(<Activate />)}
									/>
									<Route
										path="webauthn"
										element={withSuspense(<Webauthn />)}
									/>
								</Route>
								<Route path="callbacks">
									<Route
										path="oidc"
										element={withSuspense(<OidcCallback />)}
									/>
								</Route>
							</Route>
						</Routes>
					</Router>
				</div>
			</div>
		</ThemeProvider>
	);
};

const stateToProps = ({}) => ({});
const dispatchToProps = (dispatch) => bindActionCreators({}, dispatch);

export default connect(stateToProps, dispatchToProps)(withStyles(styles)(App));
