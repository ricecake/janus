import { createUserManager } from 'redux-oidc';

var serverVars = JSON.parse(document.getElementById('openid-client-params').innerHTML);

var url = window.location.origin;
const userManagerConfig = {
	authority: url,
	response_type: 'code',
	scope: 'openid',
	silent_redirect_uri: url + '/static/oidc.html?mode=silent',
	automaticSilentRenew:true,
	validateSubOnSilentRenew: true,
	loadUserInfo: false,
	... serverVars
};

const userManager = createUserManager(userManagerConfig);

export default userManager;