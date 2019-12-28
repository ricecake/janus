import { createActions, handleActions } from 'redux-actions';

const url = new URL(window.location);
let query = {};
for (const [key, value] of url.searchParams.entries()) {
	query[key] = value;
}


let serverParamsElm = document.getElementById('server-params');
let serverParams = {};
if (serverParamsElm) {
	serverParams = JSON.parse(serverParamsElm.innerHTML);
}

const defaultState = { query: query, serverParams: serverParams };

export const {} = createActions({});

const reducer = handleActions({}, defaultState);

export default reducer;