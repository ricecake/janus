import { createActions, handleActions, combineActions } from 'redux-actions';
import { doWebauthnLogin } from 'Include/webauthn';

const defaultState = {
	loading: false,
	Email: false,
	Password: false,
	Totp: false,
	Webauthn: false,
	methods: false,
};

export const {
	loginStart,
	loginFinish,
	passwordStart,
	passwordFinish,
	webauthnStart,
	webauthnFinish,
	magicStart,
	magicFinish,
	methodsStart,
	methodsFinish,
} = createActions(
	'LOGIN_START',
	'LOGIN_FINISH',
	'PASSWORD_START',
	'PASSWORD_FINISH',
	'WEBAUTHN_START',
	'WEBAUTHN_FINISH',
	'MAGIC_START',
	'MAGIC_FINISH',
	'METHODS_START',
	'METHODS_FINISH',
	{ prefix: 'janus/login' }
);

const resubmitLoginForm = () => {
	let loginForm = document.createElement('form');
	loginForm.setAttribute('action', window.location);
	loginForm.setAttribute('method', 'post');
	loginForm.setAttribute('hidden', 'true');
	document.body.appendChild(loginForm);
	loginForm.submit();
};

export const fetchAuthMethods = (email) => {
	return (dispatch, getState) => {
		dispatch(methodsStart());
		fetch('/check/authenticators', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ email }),
		})
			.then((res) => res.json())
			.then((methods) => dispatch(methodsFinish(methods)));
	};
};

export const doPasswordAuth = (email, password, totp) => {
	return (dispatch, getState) => {
		dispatch(loginStart());

		let state = getState();
		var url = new URL(window.location.origin);

		url.pathname = '/check/auth';
		url.searchParams.append('client_id', state.context.query.client_id);

		fetch(url, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ email, password }),
		})
			.then((res) => dispatch(loginFinish()))
			.then(() => resubmitLoginForm());
	};
};

export const doWebauthn = (email) => {
	return (dispatch, getState) => {
		dispatch(webauthnStart());
		doWebauthnLogin(email)
			.then(() => dispatch(webauthnFinish()))
			.then(() => resubmitLoginForm());
	};
};

const doMagicLoginLink = () => {};

const completeLogin = () => {};

const reducer = handleActions(
	{
		[methodsFinish]: (state, { payload: methods }) =>
			merge(state, { methods: true, ...methods }),
		[combineActions(
			loginStart,
			passwordStart,
			webauthnStart,
			magicStart,
			methodsStart
		)]: (state) => merge(state, { loading: true }),
		[combineActions(
			loginFinish,
			passwordFinish,
			webauthnFinish,
			magicFinish,
			methodsFinish
		)]: (state) => merge(state, { loading: false }),
	},
	defaultState
);

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;
