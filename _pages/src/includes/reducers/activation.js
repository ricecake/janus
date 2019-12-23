import { createActions, handleActions } from 'redux-actions';

const defaultState = {
	preferred_name: '',
	password: '',
	verify_password: '',
	submitable: false,
	password_match: true,
};

export const { changeName, changePassword, changePasswordVerifier, submitForm } = createActions({
	NAME: (name = "")=>({ name }),
	PASSWORD: (password = "")=>({ password }),
	PASSWORD_VERIFIER: (verifier = "")=>({ verifier }),
	SUBMIT: ()=>({}),
});

const reducer = handleActions({
	[changeName]: (state, { payload: { name } }) => ({ ... state, preferred_name: name }),
	[changePassword]: (state, { payload: { password } }) => ({ ... state, password: password }),
	[changePasswordVerifier]: (state, { payload: { verifier } }) => ({ ... state, verify_password: verifier }),
	[submitForm]: (state, { payload }) => {
		console.log(state, payload);
		return state;
	},
}, defaultState);

export default reducer;