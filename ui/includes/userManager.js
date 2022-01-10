import { createUserManager } from 'redux-oidc';
import React, { useEffect } from 'react';

import Config from 'Include/config';
import store from 'Include/store';

var url = Config.hosts.idp_path;
const userManagerConfig = {
	authority: url,
	response_type: 'code',
	scope: 'openid profile',
	redirect_uri: url + '/callbacks/oidc/?mode=normal',
	silent_redirect_uri: url + '/callbacks/oidc/?mode=silent',
	automaticSilentRenew: true,
	validateSubOnSilentRenew: true,
	loadUserInfo: false,
	client_id: Config.identity.client_id,
};

const userManager = createUserManager(userManagerConfig);

export const ensureLoginEffect = () => {
	let state = store.getState();
	if (!state.oidc.user) {
		sessionStorage.setItem('loc', window.location.href);
		userManager.signinSilent().catch(() => {
			userManager.signinRedirect();
		});
	}
	return;
};

export const withLogin = (WrappedComponent) => (props) => {
	useEffect(ensureLoginEffect);
	return <WrappedComponent {...props} />;
};

export default userManager;
