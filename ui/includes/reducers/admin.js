import { createActions, handleActions } from 'redux-actions';
import { MakeMerge } from './helpers';
import userManager from 'Include/userManager';

const defaultState = {
	error: undefined,
	loading: false,
	loaded: false,
	contexts: [],
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
	adminError,
} = createActions(
	'FINISH_CONTEXT_FETCH',
	'FINISH_CONTEXT_UPDATE',
	'FINISH_CONTEXT_CREATE',
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

const reducer = handleActions(
	{
		[adminError]: (state, { payload: error }) =>
			merge(state, { loading: false, error }),
		[finishContextFetch]: (state, { payload: contexts }) =>
			merge(state, { loading: false, contexts }),
	},
	defaultState
);

const merge = MakeMerge();

export default reducer;
