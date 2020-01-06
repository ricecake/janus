import { createActions, handleActions, combineActions } from 'redux-actions';
import userManager from 'Include/userManager';
import { USER_FOUND } from "redux-oidc";

const defaultState = {
	preferred_name: '',
	password: '',
	verify_password: '',
	submitable: false,
	password_valid: false,
	password_match: true,
	name_valid: false,
	loading: false,
	activated: false,
};


export const submitForm = (event) =>{
	event.preventDefault();
	return (dispatch, getState) => {
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
					code: state.context.query.code,
				}),
			})
			.then(res =>{
				res.json().then(body => dispatch(submitFormFinish(body))).then(() => {
					let headers = res.headers;
					if (headers.has("X-Redirect-Location")) {
						window.location = headers.get("X-Redirect-Location");
					}
				});
			});
		}
	};
};

export const startSignin = () => (dispatch, getState) => {
	userManager.signinSilent();
	return;
};

export const { changeName, changePassword, changePasswordVerifier, submitFormStart, submitFormFinish } = createActions({
	changeName: (name = "")=>({ name }),
	changePassword: (password = "")=>({ password }),
	changePasswordVerifier: (verifier = "")=>({ verifier }),
	submitFormStart: ()=>({}),
	submitFormFinish: (data)=>({ data }),
}, { prefix: "janus/activation" });

const reducer = handleActions({
	[changeName]: (state, { payload: { name } }) => (merge(state, { preferred_name: name })),
	[changePassword]: (state, { payload: { password } }) => (merge(state, { password: password, password_match: password === state.verify_password })),
	[changePasswordVerifier]: (state, { payload: { verifier } }) => (merge(state, { verify_password: verifier, password_match: state.password === verifier })),

	[submitFormStart]: (state)=> merge(state, { loading: true }),
	[submitFormFinish]: (state, { payload })=> merge(state, { loading: false, activated: payload.Active }),
	[USER_FOUND]: (state, { payload }) => merge(state, { preferred_name: payload.profile.preferred_name || state.preferred_name }),

	[combineActions(changeName, changePassword, changePasswordVerifier, USER_FOUND)]: (state, msg) => merge(state, validate(state, msg)),
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