import { createActions, handleActions, combineActions } from 'redux-actions';

const defaultState = {
	preferred_name: '',
	password: '',
	verify_password: '',
	submitable: false,
	password_valid: false,
	password_match: true,
	name_valid: false,
};

export const { changeName, changePassword, changePasswordVerifier, submitForm } = createActions({
	changeName: (name = "")=>({ name }),
	changePassword: (password = "")=>({ password }),
	changePasswordVerifier: (verifier = "")=>({ verifier }),
	submitForm: ()=>({}),
});

const validate = (state, { payload }) => {
	let mergeState = Object.assign({}, state, payload);
	let newState = {};

	newState.name_valid     = mergeState.preferred_name.length > 0;
	newState.password_valid = mergeState.password.length >= 8;
	newState.password_match = mergeState.password === mergeState.verify_password;

	newState.submitable = newState.password_valid && newState.password_match && newState.name_valid;
	return newState;
};

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

const reducer = handleActions({
	[changeName]: (state, { payload: { name } }) => (merge(state, { preferred_name: name })),
	[changePassword]: (state, { payload: { password } }) => (merge(state, { password: password, password_match: password === state.verify_password })),
	[changePasswordVerifier]: (state, { payload: { verifier } }) => (merge(state, { verify_password: verifier, password_match: state.password === verifier })),
	[combineActions(changeName, changePassword, changePasswordVerifier)]: (state, msg) => merge(state, validate(state, msg)),
	[submitForm]: (state, { payload }) => {
		console.log(state, payload);
		// fetch("/profile/api/activate", {
		// 	method: 'POST',
		// 	headers: {
		// 		'Content-Type': 'application/json',
		// 		'Authorization': `Bearer ${ this.state.access_token }`,
		// 	},
		// 	body: JSON.stringify({
		// 		password: this.state.password,
		// 		verify_password: this.state.verify_password,
		// 		preferred_name: this.state.preferred_name,
		// 	}),
		// });
		return state;
	},
}, defaultState);

export default reducer;