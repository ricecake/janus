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
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
import DialogContentText from '@material-ui/core/DialogContentText';
import DialogTitle from '@material-ui/core/DialogTitle';

import { connect } from "react-redux";
import { initiateSignup, changeName, changeEmail } from "Include/reducers/signup";
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
	marginTop: theme.spacing(3),
  },
  submit: {
	margin: theme.spacing(3, 0, 2),
  },
}));

const SignupForm = (props) => {
  const classes = useStyles();

  return (
	<Container component="main" maxWidth="xs">
	  <CssBaseline />
	  <div className={classes.paper}>
		<Avatar className={classes.avatar}>
		  <LockOutlinedIcon />
		</Avatar>
		<Typography component="h1" variant="h5">
		  Sign up
		</Typography>
		<form className={classes.form} onClick={ props.initiateSignup } noValidate>
		  <Grid container spacing={2}>
			<Grid item xs={12}>
			  <TextField
				autoComplete="preferred_name"
				name="preferred_name"
				variant="outlined"
				required
				fullWidth
				id="preferred_name"
				label="What should we call you?"
				onChange={e => props.changeName(e.target.value)}
				value={ props.preferred_name }
				autoFocus
			  />
			</Grid>
			<Grid item xs={12}>
			  <TextField
				variant="outlined"
				required
				fullWidth
				id="email"
				label="Email Address"
				name="email"
				autoComplete="email"
				onChange={e => props.changeEmail(e.target.value)}
				value={ props.email }
			  />
			</Grid>
		  </Grid>
		  <Button
			fullWidth
			type="submit"
			variant="contained"
			color="primary"
			className={classes.submit}
			disabled={!props.submitable}
		  >
			Sign Up
		  </Button>
		  <Grid container justify="flex-end">
			<Grid item>
			<Link href={`/login?${ props.context.serverParams.RawQuery }`} variant="body2">
				Already have an account? Sign in
			  </Link>
			</Grid>
		  </Grid>
		</form>
	  </div>
      <Dialog
        open={props.enrolled}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Signup Complete"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
			Signup sucessfully submitted.  Please check your email to finish activation.
          </DialogContentText>
        </DialogContent>
      </Dialog>
	</Container>
  );
}

const stateToProps = ({signup, context}) => ({...signup, context });
const dispatchToProps = (dispatch) => bindActionCreators({ initiateSignup, changeName, changeEmail }, dispatch);

export default connect(stateToProps, dispatchToProps)(SignupForm);