import { createActions, handleActions, combineActions } from 'redux-actions';

const defaultState = {
	email: '',
	password: '',
	email_valid: '',
	submitable: false,
	loading: false,
};

export const { changeEmail, changePassword, loginStart, loginFinish } =
	createActions(
		{
			changeEmail: (email = '') => ({ email }),
			changePassword: (password = '') => ({ password }),
		},
		'LOGIN_START',
		'LOGIN_FINISH',
		{ prefix: 'janus/login' }
	);

export const initiateLogin = (event) => {
	event.preventDefault();
	return (dispatch, getState) => {
		if (getState().login.submitable) {
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
				body: JSON.stringify({
					email: state.login.email,
					password: state.login.password,
				}),
			})
				.then((res) => dispatch(loginFinish()))
				.then(() => {
					let loginForm = document.createElement('form');
					loginForm.setAttribute('action', window.location);
					loginForm.setAttribute('method', 'post');
					loginForm.setAttribute('hidden', 'true');
					document.body.appendChild(loginForm);
					loginForm.submit();
				});
		}
	};
};

const reducer = handleActions(
	{
		[changeEmail]: (state, { payload: email }) => merge(state, email),
		[changePassword]: (state, { payload: password }) =>
			merge(state, password),
		[combineActions(changeEmail, changePassword)]: (state, msg) =>
			merge(state, validate(state, msg)),
	},
	defaultState
);

const validate = (state, { payload }) => {
	let mergeState = Object.assign({}, state, payload);
	let newState = {};

	newState.email_valid = mergeState.email.length > 0;

	newState.submitable = newState.email_valid && !newState.loading;
	return newState;
};

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;
