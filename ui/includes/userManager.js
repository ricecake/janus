import { createUserManager } from 'redux-oidc';

import Config from 'Include/config';

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

export default userManager;
