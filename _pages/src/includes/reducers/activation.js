import { createActions, handleActions, combineActions } from 'redux-actions';

const defaultState = {
	preferred_name: '',
	password: '',
	verify_password: '',
	submitable: false,
	password_valid: false,
	password_match: true,
	name_valid: false,
	loading: false,
};


export const submitForm = () => (dispatch, getState) => {
	console.log(getState());
	if (getState().activation.submitable) {
		dispatch(submitFormStart());
		let state = getState();
		fetch("/profile/api/activate", {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${ state.oidc.user.access_token }`,
			},
			body: JSON.stringify({
				password: state.activation.password,
				verify_password: state.activation.verify_password,
				preferred_name: state.activation.preferred_name,
			}),
		})
		.then(res => res.json())
		.then(res => dispatch(submitFormFinish(res)));
	}
};

export const { changeName, changePassword, changePasswordVerifier, submitFormStart, submitFormFinish } = createActions({
	changeName: (name = "")=>({ name }),
	changePassword: (password = "")=>({ password }),
	changePasswordVerifier: (verifier = "")=>({ verifier }),
	submitFormStart: ()=>({}),
	submitFormFinish: (data)=>({ data }),
});

const reducer = handleActions({
	[changeName]: (state, { payload: { name } }) => (merge(state, { preferred_name: name })),
	[changePassword]: (state, { payload: { password } }) => (merge(state, { password: password, password_match: password === state.verify_password })),
	[changePasswordVerifier]: (state, { payload: { verifier } }) => (merge(state, { verify_password: verifier, password_match: state.password === verifier })),
	[combineActions(changeName, changePassword, changePasswordVerifier)]: (state, msg) => merge(state, validate(state, msg)),
	[submitFormStart]: (state)=> merge(state, { loading: true }),
	[submitFormFinish]: (state)=> merge(state, { loading: false }),
}, defaultState);

const validate = (state, { payload }) => {
	let mergeState = Object.assign({}, state, payload);
	let newState = {};

	newState.name_valid     = mergeState.preferred_name.length > 0;
	newState.password_valid = mergeState.password.length >= 8;
	newState.password_match = mergeState.password === mergeState.verify_password;

	newState.submitable = newState.password_valid && newState.password_match && newState.name_valid && !newState.loading;
	return newState;
};

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;