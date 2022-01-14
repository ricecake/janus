import { createUserManager } from 'redux-oidc';
import React, { useEffect } from 'react';
import { useSelector } from 'react-redux';

import Config from 'Include/config';
import store from 'Include/store';

var url = Config.hosts.idp_path;
const userManagerConfig = {
	...Config.identity,
	authority: url,
	redirect_uri: url + '/callbacks/oidc/?mode=normal',
	silent_redirect_uri: url + '/callbacks/oidc/?mode=silent',
};

const userManager = createUserManager(userManagerConfig);

export const ensureLoginEffect = () => {
	let state = store.getState();
	console.log(state);
	if (!state.oidc.user || state.oidc.user.expired) {
		sessionStorage.setItem('loc', window.location.href);
		userManager.signinSilent().catch(() => {
			userManager.signinRedirect();
		});
	}
	return;
};

export const withLogin = (WrappedComponent) => (props) => {
	useEffect(ensureLoginEffect);

	let user = useSelector(({ oidc: { user } }) => user);
	if (!user || user.expired) {
		return <div>loading</div>;
	}
	return <WrappedComponent {...props} />;
};

export default userManager;
