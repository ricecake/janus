import { createActions, handleActions } from 'redux-actions';

const url = new URL(window.location);
let query = {};
for (const [key, value] of url.searchParams.entries()) {
	query[key] = value;
}

const defaultState = { query: query, serverParams: window.__PRELOADED_STATE__ };

delete window.__PRELOADED_STATE__;

export const {} = createActions({}, { prefix: 'janus/context' });

const reducer = handleActions({}, defaultState);

export default reducer;
