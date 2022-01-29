import _ from 'lodash';

import store from 'Include/store';

export const hasRole = (role) => {
	let state = store.getState();
	const context = _.get(state, ['oidc', 'user', 'profile', 'ctx']);
	const roles = _.get(
		state,
		['oidc', 'user', 'profile', 'roles', context],
		[]
	);
	return roles.includes(role);
};
