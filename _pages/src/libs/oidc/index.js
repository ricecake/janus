import Oidc from 'oidc-client';


Oidc.Log.logger = console;
Oidc.Log.level = Oidc.Log.INFO;

let manager = new Oidc.UserManager({response_mode:'query'});

let url = new URL(document.location);
let params = url.searchParams;

switch (params.get("mode")) {
	case "silent":
		manager.signinCallback();
		break;
	default:
		console.log("IT BROKEN");
}