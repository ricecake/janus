import { createActions, handleActions } from 'redux-actions';
import { MakeMerge } from './helpers';
import { doWebauthnRegister } from 'Include/webauthn';
import userManager from 'Include/userManager';

const defaultState = {
	loading: false,
	loaded: false,
	user_details: {},
	error: undefined,
	authenticators: [],
	logins: {
		Logins: [],
	},
};

const handleFetchError = (res) => {
	/*
	TODO: this should be in a shared location.
	it should also handle "bad user", "bad creds", and "no permissions"
	Should almost certainly update the server to distinguish between 401 bad auth and 403 no perms
	*/
	if (!res.ok) {
		if (res.status === 401) {
			userManager.removeUser();
		}
		throw res;
	}
	return res;
};

export const {
	profileError,
	detailStartFetch,
	detailFinishFetch,
	detailStartUpdate,
	detailFinishUpdate,
	finishAuthenticatorFetch,
	finishLoginFetch,
} = createActions(
	'PROFILE_ERROR',

	'DETAIL_START_FETCH',
	'DETAIL_FINISH_FETCH',

	'DETAIL_START_UPDATE',
	'DETAIL_FINISH_UPDATE',

	'PASSWORD_START_CHANGE',
	'PASSWORD_FINISH_CHANGE',

	'FINISH_AUTHENTICATOR_FETCH',

	'FINISH_LOGIN_FETCH',
	{ prefix: 'janus/profile' }
);

export const fetchUserDetails = () => {
	return (dispatch, getState) => {
		dispatch(detailStartFetch());
		let state = getState();
		fetch('/profile/api/detail', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(detailFinishFetch(methods)))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

export const fetchAuthenticators = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		fetch('/profile/api/authenticator', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishAuthenticatorFetch(methods)))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

export const deleteAuthenticator = (name) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		fetch(`/profile/api/authenticator?name=${name}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then(() => dispatch(fetchAuthenticators()))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

export const initiatePasswordChange = (pass, verify) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return (
			fetch('/profile/api/password', {
				method: 'POST',
				headers: {
					Authorization: `Bearer ${state.oidc.user.access_token}`,
					'Content-Type': 'application/json',
				},
				body: JSON.stringify({
					password: pass,
					verify_password: verify,
				}),
			})
				.then(handleFetchError)
				// .then((methods) => dispatch(finishLoginFetch(methods)))
				.catch((err) => {
					console.log(err);
					dispatch(profileError('Something went wrong'));
				})
		);
	};
};

export const initiateWebauthnEnroll = (name) => {
	return (dispatch, getState) => {
		// dispatch(webauthnStart());
		doWebauthnRegister(name)
			.then(handleFetchError)
			.then(() => dispatch(fetchAuthenticators()))
			// .then(() => dispatch(webauthnFinish()))
			.catch((err) => {
				console.log(err);
				dispatch(profileError('Something went wrong'));
			});
	};
};

export const fetchLogins = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		fetch('/profile/api/login', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishLoginFetch(methods)))
			.catch((err) => {
				console.log(err);
				dispatch(profileError('Something went wrong'));
			});
	};
};

// TODO: Need to force a refresh of auth session when you get a 401

export const deleteSession = (code) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		fetch(`/profile/api/login/session?code=${code}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then(() => dispatch(fetchLogins()))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

export const deleteAccessContext = (code) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		fetch(`/profile/api/login?code=${code}`, {
			method: 'DELETE',
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then(() => dispatch(fetchLogins()))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

export const updateUserDetails = ({ PreferredName, GivenName, FamilyName }) => {
	return (dispatch, getState) => {
		dispatch(detailStartUpdate());
		let state = getState();
		return fetch('/profile/api/detail', {
			method: 'POST',
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ PreferredName, GivenName, FamilyName }),
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(detailFinishUpdate(methods)))
			.catch(() => dispatch(profileError('Something went wrong')));
	};
};

const reducer = handleActions(
	{
		[detailStartFetch]: (state) =>
			merge(state, { loading: true, loaded: false }),
		[detailFinishFetch]: (state, { payload: details }) =>
			merge(state, {
				loading: false,
				loaded: true,
				user_details: details,
			}),
		[detailFinishUpdate]: (state, { payload: details }) =>
			merge(state, {
				loading: false,
				loaded: true,
				user_details: details,
			}),
		[finishAuthenticatorFetch]: (state, { payload: details }) =>
			merge(state, {
				authenticators: details,
			}),
		[finishLoginFetch]: (state, { payload: logins }) =>
			merge(state, { logins }),
		[profileError]: (state, { payload: error }) =>
			merge(state, { loading: false, error }),
	},
	defaultState
);

const merge = MakeMerge();

export default reducer;
