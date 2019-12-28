import { createUserManager } from 'redux-oidc';

let serverParamsElm = document.getElementById('server-params');
let serverParams = {};
if (serverParamsElm) {
	serverParams = JSON.parse(serverParamsElm.innerHTML);
}

var url = window.location.origin;
const userManagerConfig = {
	authority: url,
	response_type: 'code',
	scope: 'openid',
	redirect_uri: url + '/static/oidc.html?mode=normal',
	silent_redirect_uri: url + '/static/oidc.html?mode=silent',
	automaticSilentRenew:true,
	validateSubOnSilentRenew: true,
	loadUserInfo: false,
	... serverParams
};

const userManager = createUserManager(userManagerConfig);

export default userManager;