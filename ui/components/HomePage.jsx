import PropTypes from 'prop-types';
import React, { useEffect } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import Avatar from '@material-ui/core/Avatar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import IconButton from '@material-ui/core/IconButton';
import Typography from '@material-ui/core/Typography';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import Identicon from 'react-identicons';
import ExitToAppOutlinedIcon from '@material-ui/icons/ExitToAppOutlined';

import { Link, Show, Hide, NavButton } from 'Component/Helpers';
import { fetchAllowedClients } from 'Include/reducers/home';

import { hasRole } from 'Include/permissions';

const useStyles = makeStyles((theme) => ({
	paper: {
		marginTop: theme.spacing(8),
		display: 'flex',
		flexDirection: 'column',
		alignItems: 'center',
	},
	avatar: {
		margin: theme.spacing(1),
		backgroundColor: theme.palette.secondary.main,
	},
	form: {
		width: '100%', // Fix IE 11 issue.
		marginTop: theme.spacing(1),
	},
	submit: {
		margin: theme.spacing(3, 0, 2),
	},
	root: {
		flexGrow: 1,
	},
	menuButton: {
		marginRight: theme.spacing(2),
	},
	title: {
		flexGrow: 1,
	},
	app_bar: {
		marginBottom: theme.spacing(1),
	},
}));

const ClientDetails = (props) => (
	<Grid item xs="auto">
		<Card variant="outlined">
			<CardHeader
				avatar={
					<Avatar alt={props.display_name}>
						<Identicon string={props.client_id} size="25" />
					</Avatar>
				}
				title={props.display_name}
				subheader={props.description}
			/>
			<CardContent>
				<Grid
					container
					direction="row"
					justifyContent="space-evenly"
					alignItems="center"
				>
					<NavButton
						to={props.base_uri}
						variant="contained"
						color="secondary"
						startIcon={<ExitToAppOutlinedIcon />}
					>
						Login
					</NavButton>
				</Grid>
			</CardContent>
		</Card>
	</Grid>
);

ClientDetails.propTypes = {
	base_uri: PropTypes.any,
	client_id: PropTypes.any,
	description: PropTypes.any,
	display_name: PropTypes.any,
};

const ContextDetails = (props) => (
	<Grid item md={12} lg={4}>
		{/** or should this be auto? */}
		<Card variant="outlined">
			<CardHeader
				avatar={
					<Avatar alt={props.display_name}>
						<Identicon string={props.context} size="25" />
					</Avatar>
				}
				title={props.display_name}
				subheader={props.description}
			/>
			<CardContent>
				<Grid
					container
					spacing={2}
					direction="row"
					justifyContent="space-evenly"
					alignItems="center"
				>
					{props.clients.map((client) => (
						<ClientDetails key={client.client_id} {...client} />
					))}
				</Grid>
			</CardContent>
		</Card>
	</Grid>
);

ContextDetails.propTypes = {
	clients: PropTypes.shape({
		map: PropTypes.func,
	}),
	context: PropTypes.any,
	description: PropTypes.any,
	display_name: PropTypes.any,
};

const ResponsiveAppBar = (props) => {
	const navigate = useNavigate();
	const classes = useStyles();
	const [anchorEl, setAnchorEl] = React.useState(null);
	const open = Boolean(anchorEl);

	const handleMenu = (event) => {
		setAnchorEl(event.currentTarget);
	};

	const handleClose = () => {
		setAnchorEl(null);
	};

	return (
		<AppBar position="static" className={classes.app_bar}>
			<Toolbar>
				<Typography variant="h6" className={classes.title}>
					Applications
				</Typography>
				<Typography variant="subtitle1" className={classes.title}>
					{props.profile.preferred_name}
				</Typography>
				<div>
					<IconButton
						aria-label="account of current user"
						aria-controls="menu-appbar"
						aria-haspopup="true"
						onClick={handleMenu}
						color="inherit"
					>
						<Avatar alt={props.profile.preferred_name}>
							<Identicon string={props.profile.sub} size="25" />
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
				</div>
			</Toolbar>
		</AppBar>
	);
};

ResponsiveAppBar.propTypes = {
	profile: PropTypes.shape({
		preferred_name: PropTypes.any,
		sub: PropTypes.any,
	}),
};

const HomeAppMenu = (props) => {
	useEffect(() => {
		props.fetchAllowedClients();
	}, []);

	return (
		<React.Fragment>
			<ResponsiveAppBar {...props} />
			<Container component="main" maxWidth="xl">
				<CssBaseline />
				<Grid
					container
					spacing={2}
					direction="row"
					justifyContent="space-evenly"
					alignItems="baseline"
				>
					{props.clientDetails.map((context) => (
						<ContextDetails key={context.context} {...context} />
					))}
				</Grid>
			</Container>
		</React.Fragment>
	);
};

HomeAppMenu.propTypes = {
	clientDetails: PropTypes.shape({
		map: PropTypes.func,
	}),
	fetchAllowedClients: PropTypes.func,
};

const stateToProps = ({
	home,
	oidc: {
		user: { profile },
	},
}) => ({
	profile,
	...home,
});
const dispatchToProps = (dispatch) =>
	bindActionCreators({ fetchAllowedClients }, dispatch);

export default connect(stateToProps, dispatchToProps)(HomeAppMenu);
