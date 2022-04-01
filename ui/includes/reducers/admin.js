import { createActions, handleActions } from 'redux-actions';
import { MakeMerge } from './helpers';
import userManager from 'Include/userManager';

const defaultState = {
	error: undefined,
	loading: false,
	loaded: false,
	contexts: [],
	clients: [],
	users: [],
	roles: [],
	actions: [],
	groups: [],
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
	finishContextFetch,
	finishContextUpdate,
	finishContextCreate,
	finishClientFetch,
	finishClientUpdate,
	finishClientCreate,
	finishUserFetch,
	finishRoleFetch,
	finishActionFetch,
	finishGroupFetch,
	adminError,
} = createActions(
	'FINISH_CONTEXT_FETCH',
	'FINISH_CONTEXT_UPDATE',
	'FINISH_CONTEXT_CREATE',
	'FINISH_CLIENT_FETCH',
	'FINISH_CLIENT_UPDATE',
	'FINISH_CLIENT_CREATE',
	'FINISH_USER_FETCH',
	'FINISH_ROLE_FETCH',
	'FINISH_ACTION_FETCH',
	'FINISH_GROUP_FETCH',
	'ADMIN_ERROR',
	{ prefix: 'janus/admin' }
);

export const fetchContexts = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/context', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishContextFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

export const updateContext = ({ Code, Name, Description }) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/context', {
			method: 'PATCH',
			body: JSON.stringify({ Code, Name, Description }),
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
				'Content-Type': 'application/json',
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishContextUpdate(methods)))
			.then(() => dispatch(fetchContexts()))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

export const createContext = ({ Name, Description = '' }) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/context', {
			method: 'POST',
			body: JSON.stringify({ Name, Description }),
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
				'Content-Type': 'application/json',
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishContextCreate(methods)))
			.then(() => dispatch(fetchContexts()))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

export const fetchClients = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/client', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishClientFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

export const updateClient = (args) => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/client', {
			method: 'PATCH',
			body: JSON.stringify(args),
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
				'Content-Type': 'application/json',
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishClientUpdate(methods)))
			.then(() => dispatch(fetchClients()))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

export const fetchUsers = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/user', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishUserFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};
export const fetchRoles = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/role', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishRoleFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};
export const fetchActions = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/action', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishActionFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};
export const fetchGroups = () => {
	return (dispatch, getState) => {
		// dispatch(detailStartFetch());
		let state = getState();
		return fetch('/admin/api/group', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(finishGroupFetch(methods)))
			.catch(() => dispatch(adminError('Something went wrong')));
	};
};

const reducer = handleActions(
	{
		[adminError]: (state, { payload: error }) =>
			merge(state, { loading: false, error }),
		[finishContextFetch]: (state, { payload: contexts }) =>
			merge(state, { loading: false, contexts }),
		[finishClientFetch]: (state, { payload: clients }) =>
			merge(state, { loading: false, clients }),
		[finishUserFetch]: (state, { payload: users }) =>
			merge(state, { loading: false, users }),
		[finishRoleFetch]: (state, { payload: roles }) =>
			merge(state, { loading: false, roles }),
		[finishActionFetch]: (state, { payload: actions }) =>
			merge(state, { loading: false, actions }),
		[finishGroupFetch]: (state, { payload: groups }) =>
			merge(state, { loading: false, groups }),
	},
	defaultState
);

const merge = MakeMerge();

export default reducer;
