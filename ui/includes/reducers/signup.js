import { createActions, handleActions } from 'redux-actions';
import { doWebauthnRegister } from 'Include/webauthn';

const defaultState = {
	email: '',
	loading: false,
	enrolled: false,
	webauthn: false,
	password: false,
};

export const {
	signupStart,
	signupFinish,
	webauthnStart,
	webauthnFinish,
	passwordStart,
	passwordFinish,
} = createActions(
	'SIGNUP_START',
	'SIGNUP_FINISH',
	'WEBAUTHN_START',
	'WEBAUTHN_FINISH',
	'PASSWORD_START',
	'PASSWORD_FINISH',
	{ prefix: 'janus/signup' }
);

export const initiateSignup = (preferred_name, email) => {
	return (dispatch, getState) => {
		dispatch(signupStart());
		fetch('/signup', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				preferred_name,
				email,
			}),
		}).then((res) => dispatch(signupFinish(email)));
	};
};

export const initiateWebauthnEnroll = () => {
	return (dispatch, getState) => {
		dispatch(webauthnStart());
		doWebauthnRegister().then(() => dispatch(webauthnFinish()));
	};
};

export const initiatePasswordEnroll = (password, verify_password) => {
	return (dispatch, getState) => {
		dispatch(passwordStart());
		fetch('/signup/password', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				password: password,
				verify_password: verify_password,
			}),
		}).then((res) => dispatch(passwordFinish()));
	};
};

const reducer = handleActions(
	{
		[signupStart]: (state) => merge(state, { loading: true }),
		[webauthnStart]: (state) => merge(state, { loading: true }),
		[passwordStart]: (state) => merge(state, { loading: true }),
		[signupFinish]: (state, { payload: email }) =>
			merge(state, { enrolled: true, loading: false, email: email }),
		[webauthnFinish]: (state) =>
			merge(state, { webauthn: true, loading: false }),
		[passwordFinish]: (state) =>
			merge(state, { password: true, loading: false }),
	},
	defaultState
);

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;
