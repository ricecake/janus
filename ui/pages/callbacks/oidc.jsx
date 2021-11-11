import React from 'react';
import Oidc from 'oidc-client';
import userManager from 'Include/userManager';
import { useSearchParams } from 'react-router-dom';

export const OidcCallback = (props) => {
	const [searchParams] = useSearchParams();

	React.useEffect(() => {
		Oidc.Log.logger = console;
		Oidc.Log.level = Oidc.Log.INFO;

		switch (searchParams.get('mode')) {
			case 'normal':
			case 'silent':
				userManager.signinCallback();
				break;
			default:
				console.log('IT BROKEN');
		}
	}, []);

	return <span>Auth...</span>;
};
export default OidcCallback;
