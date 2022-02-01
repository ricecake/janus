import React, { useEffect } from 'react';
import UAParser from 'ua-parser-js';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';

import { makeStyles } from '@material-ui/core/styles';
import IconButton from '@material-ui/core/IconButton';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListSubheader from '@material-ui/core/ListSubheader';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Collapse from '@material-ui/core/Collapse';
import ListItemAvatar from '@material-ui/core/ListItemAvatar';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import DeleteIcon from '@material-ui/icons/Delete';
import ProfilePage from './frame';
import Identicon from 'react-identicons';

import {
	fetchLogins,
	deleteAccessContext,
	deleteSession,
} from 'Include/reducers/profile';

const useStyles = makeStyles((theme) => ({
	root: {
		padding: '2px 4px',
		display: 'flex',
		alignItems: 'center',
		width: 400,
	},
	input: {
		marginLeft: theme.spacing(1),
		flex: 1,
	},
	iconButton: {
		padding: 10,
	},
	divider: {
		height: 28,
		margin: 4,
	},
}));

const LoginContext = ({
	code,
	deleteAccessContext,
	created_at,
	client: { display_name, client_id },
}) => {
	const created = new Date(created_at);

	return (
		<ListItem>
			<ListItemAvatar>
				<Identicon string={client_id} size="25" />
			</ListItemAvatar>
			<ListItemText
				primary={display_name}
				secondary={`${created.toDateString()} ${created.toLocaleTimeString()}`}
			/>
			<ListItemSecondaryAction>
				<IconButton
					edge="end"
					aria-label="delete"
					color="primary"
					onClick={() => deleteAccessContext(code)}
				>
					<DeleteIcon />
				</IconButton>
			</ListItemSecondaryAction>
		</ListItem>
	);
};
const AppContext = ({
	display_name,
	description,
	access_context,
	deleteAccessContext,
}) => {
	const [open, setOpen] = React.useState(true);
	React.useEffect(() => {
		console.log(access_context);
	});
	return (
		<>
			<ListItem button onClick={() => setOpen(!open)}>
				<ListItemText primary={display_name} secondary={description} />
				{open ? <ExpandLess /> : <ExpandMore />}
			</ListItem>
			<Collapse in={open} timeout="auto" unmountOnExit>
				<List>
					<ListSubheader>Logins</ListSubheader>

					{access_context.map((access) => (
						<LoginContext
							key={access.code}
							deleteAccessContext={deleteAccessContext}
							{...access}
						/>
					))}
				</List>
			</Collapse>
		</>
	);
};

const BrowserSession = ({
	code,
	created_at,
	user_agent,
	context,
	deleteSession,
	deleteAccessContext,
}) => {
	const created = new Date(created_at);
	var parser = new UAParser();
	parser.setUA(user_agent);
	const {
		browser: { name: browserName },
		os: { name: osName },
	} = parser.getResult();
	return (
		<Card variant="outlined">
			<CardHeader
				title={`${browserName} on ${osName}`}
				subheader={`${created.toDateString()} ${created.toLocaleTimeString()}`}
				action={
					<IconButton
						edge="end"
						aria-label="delete"
						color="primary"
						onClick={() => deleteSession(code)}
					>
						<DeleteIcon />
					</IconButton>
				}
			/>
			<CardContent>
				<List>
					<ListSubheader>Apps groups</ListSubheader>
					{context.map((con) => (
						<AppContext
							key={con.code}
							deleteAccessContext={deleteAccessContext}
							{...con}
						/>
					))}
				</List>
			</CardContent>
		</Card>
	);
};
const LoginsBase = ({
	fetchLogins,
	deleteAccessContext,
	deleteSession,
	logins,
}) => {
	const classes = useStyles();

	React.useEffect(() => {
		fetchLogins();
	}, []);

	React.useEffect(() => {
		console.log(logins);
	});

	return (
		<ProfilePage>
			{logins.Logins.map((session) => (
				<BrowserSession
					key={session.code}
					deleteAccessContext={deleteAccessContext}
					deleteSession={deleteSession}
					{...session}
				/>
			))}
		</ProfilePage>
	);
};

const stateToProps = ({ profile }) => ({ ...profile });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{
			fetchLogins,
			deleteAccessContext,
			deleteSession,
		},
		dispatch
	);

export const Logins = connect(stateToProps, dispatchToProps)(LoginsBase);
export default Logins;
