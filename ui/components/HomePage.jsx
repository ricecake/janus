import React, { useEffect } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Avatar from '@material-ui/core/Avatar';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Identicon from 'react-identicons';

import { Link, Show, Hide } from 'Component/Helpers';

import { fetchAllowedClients } from 'Include/reducers/home';

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
			/>
			<CardContent>
				<a href={props.base_uri}>{props.display_name}</a>
			</CardContent>
		</Card>
	</Grid>
);

const ContextDetails = (props) => (
	<Grid item md={12} lg={4}>
		{' '}
		{/** or should this be auto? */}
		<Card variant="outlined">
			<CardHeader
				avatar={
					<Avatar alt={props.display_name}>
						<Identicon string={props.context} size="25" />
					</Avatar>
				}
				title={props.display_name}
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

const HomeAppMenu = (props) => {
	const classes = useStyles();

	useEffect(() => {
		props.fetchAllowedClients();
	}, []);

	return (
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
	);
};

const stateToProps = ({ home }) => ({ ...home });
const dispatchToProps = (dispatch) =>
	bindActionCreators({ fetchAllowedClients }, dispatch);

export default connect(stateToProps, dispatchToProps)(HomeAppMenu);
