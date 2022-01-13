import { createActions, handleActions } from 'redux-actions';
import { doWebauthnRegister } from 'Include/webauthn';

const defaultState = {
	email: '',
	loading: false,
	enrolled: false,
	webauthn: false,
	password: false,
	error: undefined,
};

export const {
	signupStart,
	signupFinish,
	webauthnStart,
	webauthnFinish,
	passwordStart,
	passwordFinish,
	signupError,
} = createActions(
	'SIGNUP_START',
	'SIGNUP_FINISH',
	'WEBAUTHN_START',
	'WEBAUTHN_FINISH',
	'PASSWORD_START',
	'PASSWORD_FINISH',
	'SIGNUP_ERROR',
	{ prefix: 'janus/signup' }
);

const handleFetchError = (res) => {
	if (!res.ok) {
		throw res;
	}
	return res;
};

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
		})
			.then(handleFetchError)
			.then((res) => dispatch(signupFinish(email)))
			.catch(() => dispatch(signupError('Something went wrong')));
	};
};

export const initiateWebauthnEnroll = () => {
	return (dispatch, getState) => {
		dispatch(webauthnStart());
		doWebauthnRegister()
			.then(handleFetchError)
			.then(() => dispatch(webauthnFinish()))
			.catch(() => dispatch(signupError('Something went wrong')));
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
		})
			.then(handleFetchError)
			.then((res) => dispatch(passwordFinish()))
			.catch(() => dispatch(signupError('Something went wrong')));
	};
};

const reducer = handleActions(
	{
		[signupError]: (state, { payload: error }) => merge(state, { error }),
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
