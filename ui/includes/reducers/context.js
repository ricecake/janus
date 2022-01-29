import { createActions, handleActions } from 'redux-actions';
import { USER_FOUND } from 'redux-oidc';
import { MakeMerge } from './helpers';
import store from 'Include/store';

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

// export const {} = createActions({}, { prefix: 'janus/context' });

const reducer = handleActions(
	{
		// [USER_FOUND]: (state, { payload: { profile: { ctx, roles } }}) => {
		// 	// TODO: make this extract user roles, and save into roles part of state
		// 	// console.log(profile);
		// 	merge(state, { roles: { roles: roles[ctx].reduce((acc, item) => acc[item]= true, {}) } });
		// },
	},
	defaultState
);

// const merge = MakeMerge();

export default reducer;
