import { createActions, handleActions, combineActions } from 'redux-actions';

const defaultState = {
	preferred_name: '',
	email: '',
	name_valid: false,
	email_valid: false,
	loading: false,
	enrolled: false,
};

export const {changeName, changeEmail, signupStart, signupFinish} = createActions({
	changeName: (preferred_name = "")=>({ preferred_name }),
	changeEmail: (email = "") => ({ email }),
},'SIGNUP_START', 'SIGNUP_FINISH', { prefix: "janus/signup" });

export const initiateSignup = (event) => {
	event.preventDefault();
	return (dispatch, getState) => {
		if (getState().signup.submitable) {
			dispatch(signupStart());
			let state = getState();
			fetch("/signup", {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					preferred_name: state.signup.preferred_name,
					email: state.signup.email,
				}),
			})
			.then(res => dispatch(signupFinish()));
		}
	};
};

const reducer = handleActions({
	[changeName]: (state, { payload: name }) => (merge(state, name)),
	[changeEmail]: (state, { payload: email }) => (merge(state, email)),
	[signupStart]: (state) => merge(state, {loading: true}),
	[signupFinish]: (state) => merge(state, {enrolled: true, loading: false}),
	[combineActions(changeName, changeEmail)]: (state, msg) => merge(state, validate(state, msg)),
}, defaultState);

const validate = (state, { payload }) => {
	let mergeState = Object.assign({}, state, payload);
	let newState = {};

	newState.name_valid  = mergeState.preferred_name.length > 0;
	newState.email_valid = mergeState.email.length > 0;

	newState.submitable = newState.name_valid && newState.email_valid && !newState.loading && !newState.enrolled;
	return newState;
};

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;