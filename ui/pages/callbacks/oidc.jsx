import React from 'react';
import Oidc from 'oidc-client';
import userManager from 'Include/userManager';
import { useSearchParams } from 'react-router-dom';

export const OidcCallback = (props) => {
	const [searchParams] = useSearchParams();

	React.useEffect(() => {
		Oidc.Log.logger = console;
		Oidc.Log.level = Oidc.Log.DEBUG;

		console.log(searchParams.get('mode'));
		switch (searchParams.get('mode')) {
			case 'normal':
				userManager.signinRedirectCallback().then(() => {
					let redir = sessionStorage.getItem('loc');
					if (redir) {
						console.log(`redirect: ${redir}`);
						sessionStorage.removeItem('loc');
						window.location = redir;
					}
				});
				break;
			case 'silent':
				userManager.signinSilentCallback();
				break;
			default:
				console.log('IT BROKEN');
		}

		// if we have a redirect in params, bounce to that path, but only inside this domain
	}, []);

	return <span>Auth...</span>;
};
export default OidcCallback;
