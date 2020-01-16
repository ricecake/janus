import Oidc from 'oidc-client';


Oidc.Log.logger = console;
Oidc.Log.level = Oidc.Log.INFO;

var path = window.location.origin;
const userManagerConfig = {
	authority: path,
	response_type: 'code',
	scope: 'openid profile',
	redirect_uri: path + '/static/oidc.html?mode=normal',
	silent_redirect_uri: path + '/static/oidc.html?mode=silent',
	automaticSilentRenew:true,
	validateSubOnSilentRenew: true,
	loadUserInfo: false,
	response_mode:'query',
};

let manager = new Oidc.UserManager(userManagerConfig);

let url = new URL(document.location);
let params = url.searchParams;

switch (params.get("mode")) {
	case "normal":
	case "silent":
		manager.signinCallback();
		break;
	default:
		console.log("IT BROKEN");
}