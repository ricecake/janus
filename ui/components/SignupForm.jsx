import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import Avatar from '@material-ui/core/Avatar';
import Box from '@material-ui/core/Box';
import Button from '@material-ui/core/Button';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import CardHeader from '@material-ui/core/CardHeader';
import Container from '@material-ui/core/Container';
import CssBaseline from '@material-ui/core/CssBaseline';
import Grid from '@material-ui/core/Grid';
import LinearProgress from '@material-ui/core/LinearProgress';
import Paper from '@material-ui/core/Paper';
import TextField from '@material-ui/core/TextField';
import Typography from '@material-ui/core/Typography';
import Alert from '@material-ui/lab/Alert';

import { makeStyles } from '@material-ui/core/styles';

import FingerprintOutlinedIcon from '@material-ui/icons/FingerprintOutlined';
import EmailOutlinedIcon from '@material-ui/icons/EmailOutlined';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';

import { Link, Show, Hide } from 'Component/Helpers';

import {
	initiateSignup,
	initiateWebauthnEnroll,
	initiatePasswordEnroll,
} from 'Include/reducers/signup';
import { webauthnCapable } from 'Include/webauthn';

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
		marginTop: theme.spacing(3),
	},
	submit: {
		margin: theme.spacing(3, 0, 2),
	},
}));

// Ask for email; advance to Identity
const Identity = ({ changeStep, enrolled, initiateSignup }) => {
	React.useEffect(() => {
		if (enrolled) {
			changeStep('pick');
		}
	}, [enrolled]);

	const [email, setEmail] = React.useState('');
	const [name, setName] = React.useState('');

	return (
		<Box>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					initiateSignup(name, email);
				}}
			>
				<Grid item>
					<TextField
						variant="outlined"
						required
						fullWidth
						id="email"
						label="Email Address"
						name="email"
						margin="normal"
						autoComplete="email"
						autoFocus
						onChange={(e) => setEmail(e.target.value)}
						error={!!email && !/^\S+@\S+\.\S+$/.test(email)}
					/>
				</Grid>
				<Grid item>
					<TextField
						autoComplete="name"
						name="preferred_name"
						variant="outlined"
						required
						fullWidth
						id="preferred_name"
						margin="normal"
						label="What should we call you?"
						onChange={(e) => setName(e.target.value)}
					/>
				</Grid>
				<Button
					fullWidth
					type="submit"
					variant="contained"
					color="primary"
				>
					Choose Authentication Method
				</Button>
			</form>
		</Box>
	);
};

// Over choices for login, after identifying options
// user/pass+totp, webauthn; advance to choice
const PickLogin = ({ changeStep, ...props }) => (
	<Box>
		<ButtonGroup fullWidth orientation="vertical">
			<Show If={webauthnCapable()}>
				<Button
					startIcon={<FingerprintOutlinedIcon />}
					onClick={() => changeStep('webauthn')}
					fullWidth
					variant="contained"
				>
					Configure Platform Authentication
				</Button>
			</Show>
			<Button
				startIcon={<LockOutlinedIcon />}
				onClick={() => changeStep('password')}
				fullWidth
				variant="contained"
			>
				Configure Password Authentication
			</Button>
			<Button
				startIcon={<EmailOutlinedIcon />}
				onClick={() => changeStep('finish')}
				fullWidth
				variant="contained"
			>
				Magic Link Only
			</Button>
		</ButtonGroup>
	</Box>
);

// configure password auth -- probably merge with totp step; advance to finish
const SetupPassword = ({
	changeStep,
	loading,
	password,
	email,
	initiatePasswordEnroll,
}) => {
	React.useEffect(() => {
		if (password) {
			changeStep('finish');
		}
	}, [password]);

	const [pass, setPass] = React.useState('');
	const [verify, setVerify] = React.useState('');

	return (
		<Box>
			<form
				onSubmit={(e) => {
					e.preventDefault();
					if (verify && verify === pass) {
						initiatePasswordEnroll(pass, verify);
					}
				}}
			>
				<input name="username" type="hidden" value={email} />
				<Show If={loading}>
					<LinearProgress />
				</Show>
				<Grid item>
					<TextField
						required
						fullWidth
						autoFocus
						disabled={loading}
						name="password"
						label="Password"
						type="password"
						variant="outlined"
						margin="normal"
						autoComplete="new-password"
						onChange={(e) => setPass(e.target.value)}
						error={!!pass && pass.length < 8}
					/>
				</Grid>
				<Grid item>
					<TextField
						required
						fullWidth
						disabled={loading}
						name="verify_password"
						label="Verify Password"
						type="password"
						variant="outlined"
						margin="normal"
						autoComplete="new-password"
						onChange={(e) => setVerify(e.target.value)}
						error={!!verify && verify !== pass}
						helperText={
							!!verify && !!pass && verify !== pass
								? "Passwords don't seem to match..."
								: ''
						}
					/>
				</Grid>
				<Grid
					container
					direction="row"
					justifyContent="space-around"
					alignItems="center"
				>
					<Button
						disabled={loading || (!!verify && verify !== pass)}
						type="submit"
						variant="contained"
						color="primary"
					>
						Finish
					</Button>
				</Grid>
			</form>
		</Box>
	);
};

// configure webauthn; offer also configure password; advance to password or finish
const SetupWebauthn = ({
	changeStep,
	webauthn,
	initiateWebauthnEnroll,
	loading,
}) => (
	<Box>
		<Grid item>
			<Show If={loading}>
				<LinearProgress />
			</Show>
			<Hide If={loading || webauthn}>
				<Grid
					container
					direction="column"
					justifyContent="space-around"
					alignItems="center"
				>
					<Typography>
						Configure advanced, secure signin using your device.
						<br />
						Your biometric information or PIN will not be shared.
					</Typography>
					<Grid
						container
						direction="row"
						justifyContent="space-around"
						alignItems="center"
					>
						<Button
							onClick={() => initiateWebauthnEnroll()}
							variant="contained"
							color="primary"
							disabled={webauthn || loading}
						>
							Start Setup
						</Button>
					</Grid>
				</Grid>
			</Hide>
		</Grid>
		<Grid item>
			<Grid
				container
				direction="row"
				justifyContent="space-between"
				alignItems="center"
			>
				<Button
					onClick={() => changeStep('password')}
					variant="contained"
					color="primary"
					disabled={!webauthn}
				>
					Configure Password
				</Button>

				<Button
					onClick={() => changeStep('finish')}
					variant="contained"
					color="primary"
					disabled={!webauthn}
				>
					Finish
				</Button>
			</Grid>
		</Grid>
	</Box>
);
// tell user everthing worked.  offer button that moves to login and thence to destination
const Finish = ({ login }) => (
	<Grid
		container
		direction="row"
		justifyContent="space-around"
		alignItems="center"
	>
		<Button
			onClick={() => (window.location = login)}
			variant="contained"
			color="primary"
		>
			Login
		</Button>
	</Grid>
);

const stepTitles = {
	identity: 'User Details',
	pick: 'Choose Authentication Method',
	password: 'Configure Password',
	webauthn: 'Setup Platform Authentication',
	finish: 'Finish Signup',
};

const FormPage = ({
	step,
	changeStep,
	initiateSignup,
	initiateWebauthnEnroll,
	initiatePasswordEnroll,
	enrolled,
	webauthn,
	password,
	loading,
	email,
	context,
}) => {
	const page = () => {
		switch (step) {
			case 'identity':
				return (
					<Identity
						changeStep={changeStep}
						initiateSignup={initiateSignup}
						enrolled={enrolled}
					/>
				);
			case 'pick':
				return <PickLogin changeStep={changeStep} />;
			case 'password':
				return (
					<SetupPassword
						changeStep={changeStep}
						initiatePasswordEnroll={initiatePasswordEnroll}
						password={password}
						loading={loading}
						email={email}
					/>
				);
			case 'webauthn':
				return (
					<SetupWebauthn
						changeStep={changeStep}
						initiateWebauthnEnroll={initiateWebauthnEnroll}
						webauthn={webauthn}
						loading={loading}
					/>
				);
			case 'finish':
				return (
					<Finish
						changeStep={changeStep}
						login={`/login?${context.serverParams.RawQuery}`}
					/>
				);
			default:
				return (
					<Finish
						changeStep={changeStep}
						login={`/login?${context.serverParams.RawQuery}`}
					/>
				);
		}
	};

	return (
		<Container>
			<CardHeader
				title={stepTitles[step] || 'Error'}
				avatar={<LockOutlinedIcon />}
			/>
			{page()}
		</Container>
	);
};

const SignupForm = (props) => {
	const classes = useStyles();

	const [activeStep, setActiveStep] = React.useState('identity');

	return (
		<Container component="main" maxWidth="sm">
			<CssBaseline />
			<Paper className={classes.paper}>
				<Avatar className={classes.avatar}>
					<LockOutlinedIcon />
				</Avatar>
				<Typography component="h1" variant="h5">
					Sign up
				</Typography>
				<Show If={props.error}>
					<Alert severity="error">{props.error}</Alert>
				</Show>
				<FormPage
					step={activeStep}
					changeStep={setActiveStep}
					{...props}
				/>
				<Grid container justify="flex-end">
					<Grid item>
						<Link
							to={`/login?${props.context.serverParams.RawQuery}`}
							variant="body2"
						>
							Already have an account? Sign in
						</Link>
					</Grid>
				</Grid>
			</Paper>
		</Container>
	);
};

const stateToProps = ({ signup, context }) => ({ ...signup, context });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{ initiateSignup, initiateWebauthnEnroll, initiatePasswordEnroll },
		dispatch
	);

export default connect(stateToProps, dispatchToProps)(SignupForm);
