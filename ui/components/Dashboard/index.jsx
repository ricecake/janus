import React from 'react';
import { makeStyles } from '@material-ui/core/styles';
import CssBaseline from '@material-ui/core/CssBaseline';
import Hidden from '@material-ui/core/Hidden';
import BasePage from 'Component/BasePage';
import userManager, { withLogin } from 'Include/userManager';
import { OidcProvider } from 'redux-oidc';
import store from 'Include/store';

import Navigator from './Navigator';
import Header from './Header';

const drawerWidth = 256;
const useStyles = makeStyles((theme) => ({
	root: {
		display: 'flex',
		minHeight: '100vh',
	},
	drawer: {
		[theme.breakpoints.up('sm')]: {
			width: drawerWidth,
			flexShrink: 0,
		},
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
}));

export const Dashboard = withLogin(
	({ title = '', root = '', categories = [], ...props }) => {
		const [mobileOpen, setMobileOpen] = React.useState(false);
		const classes = useStyles(props);

		const handleDrawerToggle = () => {
			setMobileOpen(!mobileOpen);
		};

		return (
			<OidcProvider store={store} userManager={userManager}>
				<div className={classes.root}>
					<nav className={classes.drawer}>
						<Hidden smUp implementation="js">
							<Navigator
								PaperProps={{
									style: { width: drawerWidth },
								}}
								variant="temporary"
								open={mobileOpen}
								onClose={handleDrawerToggle}
								root={root}
								title={title}
								categories={categories}
							/>
						</Hidden>
						<Hidden xsDown implementation="css">
							<Navigator
								PaperProps={{
									style: { width: drawerWidth },
								}}
								root={root}
								title={title}
								categories={categories}
							/>
						</Hidden>
					</nav>
					<div className={classes.app}>
						<Header onDrawerToggle={handleDrawerToggle} />
						<BasePage>{props.children}</BasePage>
					</div>
				</div>
			</OidcProvider>
		);
	}
);
