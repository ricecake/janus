import React from 'react';
import Avatar from '@material-ui/core/Avatar';
import Button from '@material-ui/core/Button';
import CssBaseline from '@material-ui/core/CssBaseline';
import TextField from '@material-ui/core/TextField';
import Grid from '@material-ui/core/Grid';
import LockOutlinedIcon from '@material-ui/icons/LockOutlined';
import Typography from '@material-ui/core/Typography';
import { makeStyles } from '@material-ui/core/styles';
import Container from '@material-ui/core/Container';

import { connect } from "react-redux";
import { changeName, changePassword, changePasswordVerifier, submitForm, startSignin } from "Include/reducers/activation";
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


const ActivationDetails = (props) => {
	if (!props.user) {
		props.startSignin();
		return null;
	}
	const classes = useStyles();

	return (
		<Container component="main" maxWidth="xs">
		<CssBaseline />
		<div className={classes.paper}>
		  <Avatar className={classes.avatar}>
			<LockOutlinedIcon />
		  </Avatar>
		  <Typography component="h1" variant="h5">
			Activate User
		  </Typography>
		  <form className={classes.form} onSubmit={ props.initiateSignup } noValidate>
			<Grid container spacing={2}>
			  <Grid item xs={12}>
				<TextField
					required
					fullWidth
					autoFocus
					name='preferred_name'
					label="What should we call you?"
					type="text"
					variant="outlined"
					margin="normal"
					autoComplete="preferred_name"
					error={!props.name_valid}
					helperText={props.name_valid?'':"We have to call you something!"}
					onChange={e => props.changeName(e.target.value)}
					value={ props.preferred_name }
				/>
			  </Grid>
			  <Grid item xs={12}>
			  <TextField
					required
					fullWidth
					name='password'
					label="Password"
					type="password"
					variant="outlined"
					margin="normal"
					error={!props.password_valid}
					helperText="Password must be at least eight characters long"
					onChange={e => props.changePassword(e.target.value)}
					value={ props.password }
				/>
			  </Grid>
			  <Grid item xs={12}>
			  <TextField
					required
					fullWidth
					name='verify_password'
					label="Verify Password"
					type="password"
					variant="outlined"
					margin="normal"
					error={!props.password_match}
					helperText={props.password_match? '':"Passwords don't seem to match..."}
					onChange={e => props.changePasswordVerifier(e.target.value) }
					value={ props.verify_password }
				/>
			  </Grid>
			</Grid>
			<Button
				fullWidth
				type="submit"
				variant="contained"
				color="primary"
				onClick={ props.submitForm }
				disabled={!props.submitable}
			>
				Activate User
			</Button>
		  </form>
		</div>
	  </Container>
	);
};

const stateToProps = ({activation, oidc}) => ({...activation, user: oidc.user });
const dispatchToProps = (dispatch) => bindActionCreators({
	changeName, changePassword, changePasswordVerifier, submitForm, startSignin
}, dispatch);

export default connect(stateToProps, dispatchToProps)(ActivationDetails);