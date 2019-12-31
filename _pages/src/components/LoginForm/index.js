import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Link from '@material-ui/core/Link';
import Grid from '@material-ui/core/Grid';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';

import { connect } from "react-redux";
import { initiateLogin, changeEmail, changePassword } from "Include/reducers/login";
import { bindActionCreators } from 'redux'


const useStyles = makeStyles(theme => ({
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

  return (
	<Container component="main" maxWidth="xs">
	  <CssBaseline />
	  <div className={classes.paper}>
		<Avatar className={classes.avatar}>
		  <LockOutlinedIcon />
		</Avatar>
		<Typography component="h1" variant="h5">
		  Sign in
		</Typography>
		<form className={classes.form} onClick={ props.initiateLogin } noValidate>
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
			onChange={e => props.changeEmail(e.target.value)}
			value={ props.email }
		  />
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
			onChange={e => props.changePassword(e.target.value)}
			value={ props.password }
		  />
		  <Button
			fullWidth
			type="submit"
			variant="contained"
			color="primary"
			className={classes.submit}
			disabled={!props.submitable}
		  >
			Sign In
		  </Button>
		  <Grid container justify="flex-end">
			{/* <Grid item xs>
			  <Link href="#" variant="body2">
				Forgot password?
			  </Link>
			</Grid> */}
			<Grid item>
			  <Link href={`/signup?${ props.context.serverParams.RawQuery }`} variant="body2">
				{"Don't have an account? Sign Up"}
			  </Link>
			</Grid>
		  </Grid>
		</form>
	  </div>
	</Container>
  );
}

const stateToProps = ({login, context}) => ({...login, context });
const dispatchToProps = (dispatch) => bindActionCreators({
	initiateLogin,
	changeEmail,
	changePassword,
}, dispatch);

export default connect(stateToProps, dispatchToProps)(LoginForm);