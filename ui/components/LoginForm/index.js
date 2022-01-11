import React, { useEffect } from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Grid from '@material-ui/core/Grid';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';
import ButtonGroup from '@material-ui/core/ButtonGroup';
import FingerprintOutlinedIcon from '@material-ui/icons/FingerprintOutlined';
import EmailOutlinedIcon from '@material-ui/icons/EmailOutlined';
import LinearProgress from '@material-ui/core/LinearProgress';
import Alert from '@material-ui/lab/Alert';
import Paper from '@material-ui/core/Paper';

import { Link, Show, Hide } from 'Component/Helpers';

import { connect } from 'react-redux';
import {
	fetchAuthMethods,
	doWebauthn,
	doPasswordAuth,
	doMagicLoginLink,
} from 'Include/reducers/login';
import { bindActionCreators } from 'redux';
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
		marginTop: theme.spacing(1),
	},
	submit: {
		margin: theme.spacing(3, 0, 2),
	},
}));

const LoginForm = (props) => {
	const classes = useStyles();
	const [email, setEmail] = React.useState('');
	const [password, setPassword] = React.useState('');
	const [picked, setPicked] = React.useState(false);

	useEffect(() => {
		if (props.methods) {
			if (props.Webauthn && !props.Password) {
				setPicked('webauthn');
				props.doWebauthn(email);
			} else if (!props.Webauthn && props.Password) {
				setPicked('password');
			}
		}
	}, [props.methods]);

	useEffect(() => {
		if (props.error) {
			setPicked('');
		}
	}, [props.error]);

	return (
		<Container component="main" maxWidth="sm">
			<CssBaseline />
			<Paper className={classes.paper}>
				<Avatar className={classes.avatar}>
					<LockOutlinedIcon />
				</Avatar>
				<Typography component="h1" variant="h5">
					Sign in
				</Typography>
				<Show If={props.error}>
					<Alert severity="error">{props.error}</Alert>
				</Show>
				<Show If={props.emailSent}>
					<Alert severity="success">Login email sent!</Alert>
				</Show>
				<Container>
					<form
						className={classes.form}
						onSubmit={(e) => {
							e.preventDefault();
							props.fetchAuthMethods(email);
						}}
					>
						<TextField
							variant="outlined"
							margin="normal"
							required
							fullWidth
							id="email"
							label="Email Address"
							name="email"
							autoComplete="email"
							autoFocus
							disabled={props.methods}
							onChange={(e) => setEmail(e.target.value)}
							error={!!email && !/^\S+@\S+\.\S+$/.test(email)}
						/>
						<Hide If={props.methods}>
							<Grid container justify="flex-end">
								<Button
									variant="contained"
									color="primary"
									type="submit"
									disabled={
										!email || !/^\S+@\S+\.\S+$/.test(email)
									}
								>
									Next
								</Button>
							</Grid>
						</Hide>
					</form>
				</Container>
				<Show If={props.loading}>
					<LinearProgress />
				</Show>

				<Show If={props.Password && picked === 'password'}>
					<Container>
						<form
							className={classes.form}
							onSubmit={(e) => {
								e.preventDefault();
								props.doPasswordAuth(email, password);
							}}
						>
							<TextField
								variant="outlined"
								margin="normal"
								required
								fullWidth
								name="password"
								label="Password"
								type="password"
								id="password"
								autoComplete="current-password"
								onChange={(e) => setPassword(e.target.value)}
							/>
							<Button
								fullWidth
								type="submit"
								variant="contained"
								color="primary"
								className={classes.submit}
								disabled={!password}
							>
								Sign In
							</Button>
						</form>
					</Container>
				</Show>

				<Show If={props.methods}>
					<Container>
						<ButtonGroup fullWidth orientation="vertical">
							<Show
								If={
									webauthnCapable() &&
									props.Webauthn &&
									!picked
								}
							>
								<Button
									startIcon={<FingerprintOutlinedIcon />}
									onClick={() => {
										setPicked('webauthn');
										props.doWebauthn(email);
									}}
									fullWidth
									variant="contained"
								>
									Platform Authentication
								</Button>
							</Show>
							<Show If={props.Password && !picked}>
								<Button
									startIcon={<LockOutlinedIcon />}
									onClick={() => {
										setPicked('password');
									}}
									fullWidth
									variant="contained"
								>
									Password Authentication
								</Button>
							</Show>
							<Show If={props.Email && !props.emailSent}>
								<Button
									startIcon={<EmailOutlinedIcon />}
									onClick={() => {
										props.doMagicLoginLink(email);
									}}
									fullWidth
									variant="contained"
								>
									Magic Link Email
								</Button>
							</Show>
						</ButtonGroup>
					</Container>
				</Show>

				<Grid container justify="flex-end">
					{/* <Grid item xs>
			  <Link to="#" variant="body2">
				Forgot password?
			  </Link>
			</Grid> */}
					<Grid item>
						<Link
							to={`/signup?${props.context.serverParams.RawQuery}`}
							variant="body2"
						>
							{"Don't have an account? Sign Up"}
						</Link>
					</Grid>
				</Grid>
			</Paper>
		</Container>
	);
};

const stateToProps = ({ login, context }) => ({ ...login, context });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{
			fetchAuthMethods,
			doWebauthn,
			doPasswordAuth,
			doMagicLoginLink,
		},
		dispatch
	);

export default connect(stateToProps, dispatchToProps)(LoginForm);
