import PropTypes from 'prop-types';
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

const withSuspense =
	(Element, fallback = <div>Loading...</div>) =>
	(props) =>
		(
			<Suspense fallback={fallback}>
				<Element {...props} />
			</Suspense>
		);

const Home = withSuspense(lazy(() => import('Page/home')));
const Login = withSuspense(lazy(() => import('Page/login')));
const Logout = withSuspense(lazy(() => import('Page/logout')));
const Signup = withSuspense(lazy(() => import('Page/signup')));
const OidcCallback = withSuspense(lazy(() => import('Page/callbacks/oidc')));

const ProfileIndex = withSuspense(lazy(() => import('Page/profile')));
const ProfileLogin = withSuspense(lazy(() => import('Page/profile/logins')));
const ProfileAuthentication = withSuspense(
	lazy(() => import('Page/profile/authentication'))
);

const AdminIndex = withSuspense(lazy(() => import('Page/admin')));
const AdminAction = withSuspense(lazy(() => import('Page/admin/actions')));
const AdminClient = withSuspense(lazy(() => import('Page/admin/clients')));
const AdminContext = withSuspense(lazy(() => import('Page/admin/contexts')));
const AdminRole = withSuspense(lazy(() => import('Page/admin/roles')));
const AdminUser = withSuspense(lazy(() => import('Page/admin/users')));

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
								<Route path="login" element={<Login />} />
								<Route path="logout" element={<Logout />} />
								<Route path="signup" element={<Signup />} />
								<Route path="profile">
									<Route index element={<ProfileIndex />} />
									<Route
										path="authentication"
										element={<ProfileAuthentication />}
									/>
									<Route
										path="logins"
										element={<ProfileLogin />}
									/>
								</Route>
								<Route path="admin">
									<Route index element={<AdminIndex />} />
									<Route
										path="actions"
										element={<AdminAction />}
									/>
									<Route
										path="clients"
										element={<AdminClient />}
									/>
									<Route
										path="contexts"
										element={<AdminContext />}
									/>
									<Route
										path="roles"
										element={<AdminRole />}
									/>
									<Route
										path="users"
										element={<AdminUser />}
									/>
								</Route>
								<Route path="callbacks">
									<Route
										path="oidc"
										element={<OidcCallback />}
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

App.propTypes = {
	classes: PropTypes.shape({
		app: PropTypes.any,
		root: PropTypes.any,
	}),
};

const stateToProps = () => ({});
const dispatchToProps = (dispatch) => bindActionCreators({}, dispatch);

export default connect(stateToProps, dispatchToProps)(withStyles(styles)(App));
