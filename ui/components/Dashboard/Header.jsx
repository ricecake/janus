import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import Identicon from 'react-identicons';
import { useLocation } from 'react-router-dom';
import PropTypes from 'prop-types';

import AppBar from '@material-ui/core/AppBar';
import Avatar from '@material-ui/core/Avatar';
import Grid from '@material-ui/core/Grid';
import Hidden from '@material-ui/core/Hidden';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import { withStyles } from '@material-ui/core/styles';
import NavigateNextIcon from '@material-ui/icons/NavigateNext';
import Breadcrumbs from '@material-ui/core/Breadcrumbs';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';

import { Show, Link } from 'Component/Helpers';
import { hasRole } from 'Include/permissions';

const urlNameMap = {
	'/': 'Home',
};

const ucfirst = (string = '') => string[0].toUpperCase() + string.slice(1);

const lightColor = 'rgba(255, 255, 255, 0.7)';

const styles = (theme) => ({
	secondaryBar: {
		zIndex: 0,
	},
	menuButton: {
		marginLeft: -theme.spacing(1),
	},
	iconButtonAvatar: {
		padding: 4,
	},
	link: {
		textDecoration: 'none',
		color: lightColor,
		'&:hover': {
			color: theme.palette.common.white,
		},
	},
	button: {
		borderColor: lightColor,
	},
	root: {
		display: 'flex',
		flexDirection: 'column',
		width: 360,
	},
	lists: {
		backgroundColor: theme.palette.background.paper,
		marginTop: theme.spacing(1),
	},
	nested: {
		paddingLeft: theme.spacing(4),
	},
});

const RouterBreadcrumbs = withStyles(styles)((props) => {
	let location = useLocation();
	const { classes } = props;
	const pathnames = ['', ...location.pathname.split('/').filter((x) => x)];

	return (
		<div className={classes.root}>
			<Breadcrumbs
				aria-label="breadcrumb"
				separator={<NavigateNextIcon fontSize="small" />}
			>
				{pathnames.map((value, index) => {
					const last = index === pathnames.length - 1;
					const to = `/${pathnames.slice(1, index + 1).join('/')}`;

					return last ? (
						<Typography
							color="inherit"
							key={to}
							variant="h5"
							component="h1"
						>
							{urlNameMap[to] || ucfirst(value)}
						</Typography>
					) : (
						<Link color="inherit" to={to} key={to} variant="h5">
							{urlNameMap[to] || ucfirst(value)}
						</Link>
					);
				})}
			</Breadcrumbs>
		</div>
	);
});

function Header(props) {
	const { classes, onDrawerToggle } = props;
	const navigate = useNavigate();
	const [anchorEl, setAnchorEl] = React.useState(null);
	const open = Boolean(anchorEl);

	const handleMenu = (event) => {
		setAnchorEl(event.currentTarget);
	};

	const handleClose = () => {
		setAnchorEl(null);
	};

	return (
		<React.Fragment>
			<AppBar color="primary" position="sticky" elevation={0}>
				<Toolbar>
					<Grid container spacing={1} alignItems="center">
						<Hidden smUp>
							<Grid item>
								<IconButton
									color="inherit"
									aria-label="open drawer"
									onClick={onDrawerToggle}
									className={classes.menuButton}
								>
									<MenuIcon />
								</IconButton>
							</Grid>
						</Hidden>
						<Grid item xs>
							<RouterBreadcrumbs />
						</Grid>
						<Grid item>
							<IconButton
								aria-label="account of current user"
								aria-controls="menu-appbar"
								aria-haspopup="true"
								onClick={handleMenu}
								color="inherit"
							>
								<Avatar alt={props.profile.preferred_name}>
									<Identicon
										string={props.profile.sub}
										size="25"
									/>
								</Avatar>
							</IconButton>
							<Menu
								id="menu-appbar"
								anchorEl={anchorEl}
								anchorOrigin={{
									vertical: 'top',
									horizontal: 'right',
								}}
								keepMounted
								transformOrigin={{
									vertical: 'top',
									horizontal: 'right',
								}}
								open={open}
								onClose={handleClose}
							>
								<MenuItem
									onClick={() => {
										handleClose();
										navigate('/profile');
									}}
								>
									Profile
								</MenuItem>
								<Show If={hasRole('Admin')}>
									<MenuItem
										onClick={() => {
											handleClose();
											navigate('/admin');
										}}
									>
										Admin
									</MenuItem>
								</Show>
								<MenuItem
									onClick={() => {
										handleClose();
										navigate('/logout');
									}}
								>
									Logout
								</MenuItem>
							</Menu>
						</Grid>
					</Grid>
				</Toolbar>
			</AppBar>
		</React.Fragment>
	);
}

Header.propTypes = {
	classes: PropTypes.object.isRequired,
	onDrawerToggle: PropTypes.func.isRequired,
};

const stateToProps = ({
	oidc: {
		user: { profile },
	},
}) => ({
	profile,
});
const dispatchToProps = (dispatch) => bindActionCreators({}, dispatch);

export default connect(
	stateToProps,
	dispatchToProps
)(withStyles(styles)(Header));
