import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import SendIcon from '@material-ui/icons/Send';
import CloseIcon from '@material-ui/icons/Close';

import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemSecondaryAction from '@material-ui/core/ListItemSecondaryAction';
import ListItemText from '@material-ui/core/ListItemText';
import DeleteIcon from '@material-ui/icons/Delete';
import { NewButton } from 'Component/Helpers';

import {
	initiateWebauthnEnroll,
	fetchAuthenticators,
	deleteAuthenticator,
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

const NewAuthButton = ({ initiateWebauthnEnroll }) => {
	const classes = useStyles();
	const [open, setOpen] = React.useState(false);
	const [name, setName] = React.useState('');

	const fab = (
		<Fab color="primary" variant="extended" onClick={() => setOpen(true)}>
			<AddIcon />
			New Authenticator
		</Fab>
	);

	const input = (
		<Paper
			component="form"
			className={classes.root}
			onSubmit={(e) => {
				e.preventDefault();
				setOpen(false);
				initiateWebauthnEnroll(name);
			}}
		>
			<InputBase
				required
				autoFocus
				className={classes.input}
				placeholder="Create new authenticator"
				onChange={(e) => setName(e.target.value)}
			/>
			<IconButton
				color="primary"
				className={classes.iconButton}
				onClick={() => setOpen(false)}
			>
				<CloseIcon />
			</IconButton>
			<Divider className={classes.divider} orientation="vertical" />
			<IconButton
				disabled={!name}
				color="primary"
				className={classes.iconButton}
				type="submit"
			>
				<SendIcon />
			</IconButton>
		</Paper>
	);

	return open ? input : fab;
};

const Authenticator = ({ Name, deleteAuthenticator, CreatedAt }) => {
	const created = new Date(CreatedAt);
	return (
		<ListItem key={Name}>
			<ListItemText primary={Name} secondary={created.toDateString()} />
			<ListItemSecondaryAction>
				<IconButton
					edge="end"
					aria-label="delete"
					color="primary"
					onClick={() => deleteAuthenticator(Name)}
				>
					<DeleteIcon />
				</IconButton>
			</ListItemSecondaryAction>
		</ListItem>
	);
};

const WebauthnBase = (props) => {
	const classes = useStyles();

	React.useEffect(() => {
		props.fetchAuthenticators();
	}, []);
	React.useEffect(() => {
		console.log(props);
	});
	return (
		<Card>
			<CardHeader title="Platform Authenticators" />
			<CardContent>
				<List>
					{/* TODO: this should be an mui list, with created date as secondary text.
                also, move the signup component to use the widget from here.  Need a shared location. */}
					{props.authenticators.map((authenticator) => (
						<>
							<Authenticator
								key={authenticator.Name}
								deleteAuthenticator={props.deleteAuthenticator}
								{...authenticator}
							/>
						</>
					))}
				</List>
				<NewButton
					onSubmit={props.initiateWebauthnEnroll}
					title="Setup new authenticator"
					placeholder="Authenticator Name"
				/>
			</CardContent>
		</Card>
	);
};

const stateToProps = ({ profile }) => ({ ...profile });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{ initiateWebauthnEnroll, fetchAuthenticators, deleteAuthenticator },
		dispatch
	);

export const Webauthn = connect(stateToProps, dispatchToProps)(WebauthnBase);
