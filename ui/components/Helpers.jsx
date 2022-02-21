import React from 'react';
import { default as MuiLink } from '@material-ui/core/Link';
import { Link as RouterLink } from 'react-router-dom';
import Button from '@material-ui/core/Button';

import Fab from '@material-ui/core/Fab';
import AddIcon from '@material-ui/icons/Add';
import SendIcon from '@material-ui/icons/Send';
import CloseIcon from '@material-ui/icons/Close';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import CardContent from '@material-ui/core/CardContent';
import TextField from '@material-ui/core/TextField';
import { ButtonGroup } from '@material-ui/core';
import LinearProgress from '@material-ui/core/LinearProgress';

export const Hide = ({ If: condition, children }) => {
	if (condition) {
		return null;
	}

	return <React.Fragment>{children}</React.Fragment>;
};

export const Show = ({ If: condition, children }) => (
	<Hide If={!condition}> {children} </Hide>
);

const RrLink = (props) => {
	return /^https?:\/\//.test(props.to) ? (
		<a href={props.to} {...props} />
	) : (
		<RouterLink {...props} />
	);
};

export const Link = (props) => <MuiLink component={RrLink} {...props} />;

export const NavButton = (props) => <Button component={RrLink} {...props} />;

/*
Something like this but for both "editable" and "auto multi field new".
Basically want a generic profile edit component, and a generic, multiinput NewAuthButtin thing.
const DynamicForm = ({onSubmit, fields, placeholders, order}) => {
	for (const [key, value] of Object.entries(fields)) {
		<Input name={key} defaultValue={value} placeholder={placeholders[key]} />
	}

}

*/

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

export const NewButton = ({
	onSubmit,
	onCancel = () => {},
	title,
	placeholder,
}) => {
	const classes = useStyles();
	const [open, setOpen] = React.useState(false);
	const [name, setName] = React.useState('');

	const fab = (
		<Fab color="primary" variant="extended" onClick={() => setOpen(true)}>
			<AddIcon />
			{title}
		</Fab>
	);

	const input = (
		<Paper
			component="form"
			className={classes.root}
			onSubmit={(e) => {
				e.preventDefault();
				setOpen(false);
				onSubmit(name);
			}}
		>
			<InputBase
				required
				autoFocus
				className={classes.input}
				placeholder={placeholder}
				onChange={(e) => setName(e.target.value)}
			/>
			<IconButton
				color="primary"
				className={classes.iconButton}
				onClick={() => {
					setOpen(false);
					onCancel();
				}}
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

export const AutoEditForm = ({
	onSubmit,
	title,
	fields = [],
	initialValues = {},
	labels = {},
	autocomplete = {},
}) => {
	const classes = useStyles();

	const [profileView, setProfileView] = React.useState(true);
	const [loading, setLoading] = React.useState(false);

	const [values, setValues] = React.useState(initialValues);
	const setValue = (name, value) => setValues({ ...values, [name]: value });

	return (
		<Card>
			<CardHeader title={title} />
			<CardContent>
				<Hide If={loading}>
					<fieldset disabled={profileView}>
						{fields.map((field) => (
							<TextField
								key={field}
								defaultValue={initialValues[field]}
								onChange={(e) =>
									setValue(field, e.target.value)
								}
								margin="normal"
								variant="outlined"
								label={labels[field] || field}
								name={field}
								autoComplete={autocomplete[field]}
							/>
						))}
					</fieldset>
				</Hide>
				<Show If={profileView && !loading}>
					<Button onClick={() => setProfileView(false)}>Edit</Button>
				</Show>
				<Hide If={profileView || loading}>
					<ButtonGroup>
						<Button onClick={() => setProfileView(true)}>
							Cancel
						</Button>
						<Button
							onClick={() => {
								setProfileView(true);
								onSubmit(values);
								//TODO: take whatever is returned from onsubmit, and make it into a promise.
								// then set loading to be false when it resolved.
							}}
						>
							Save
						</Button>
					</ButtonGroup>
				</Hide>
				<Show If={loading}>
					<LinearProgress />
				</Show>
			</CardContent>
		</Card>
	);
};
