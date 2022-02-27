import { handleActions } from 'redux-actions';

const url = new URL(window.location);
let query = {};
for (const [key, value] of url.searchParams.entries()) {
	query[key] = value;
}

const defaultState = {
	query: query,
	roles: {},
	serverParams: window.__PRELOADED_STATE__,
};

delete window.__PRELOADED_STATE__;

const reducer = handleActions({}, defaultState);

export default reducer;
