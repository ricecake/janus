import { createActions, handleActions } from 'redux-actions';

const defaultState = {
	loading: false,
	loaded: false,
	clientDetails: [],
	error: undefined,
};

const handleFetchError = (res) => {
	if (!res.ok) {
		throw res;
	}
	return res;
};

export const { applistStart, applistFinish, applistError } = createActions(
	'APPLIST_START',
	'APPLIST_FINISH',
	'APPLIST_ERROR',
	{ prefix: 'janus/home' }
);

export const fetchAllowedClients = (email) => {
	return (dispatch, getState) => {
		dispatch(applistStart());
		let state = getState();
		fetch('/profile/api/applist', {
			headers: {
				Authorization: `Bearer ${state.oidc.user.access_token}`,
			},
		})
			.then(handleFetchError)
			.then((res) => res.json())
			.then((methods) => dispatch(applistFinish(methods)))
			.catch(() => dispatch(applistError('Something went wrong')));
	};
};

const reducer = handleActions(
	{
		[applistStart]: (state) =>
			merge(state, { loading: true, loaded: false }),
		[applistFinish]: (state, { payload: details }) =>
			merge(state, {
				loading: false,
				loaded: true,
				clientDetails: details,
			}),
		[applistError]: (state) => merge(state, { loading: false }),
	},
	defaultState
);

const merge = (oldState, newState) => Object.assign({}, oldState, newState);

export default reducer;
