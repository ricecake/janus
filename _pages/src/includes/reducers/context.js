import { createActions, handleActions, combineActions } from 'redux-actions';

const url = new URL(window.location);
let query = {};
for (const [key, value] of url.searchParams.entries()) {
	query[key] = value;
}

const defaultState = { query: query };

export const {} = createActions({});

const reducer = handleActions({}, defaultState);

export default reducer;