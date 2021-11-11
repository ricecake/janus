import React from 'react';
import BasePage from 'Component/BasePage';
import ActivationDetails from 'Component/ActivationDetails';

import { Provider } from 'react-redux';
import { OidcProvider } from 'redux-oidc';
import store from 'Include/store';
import userManager from 'Include/userManager';

export const ActivationPage = (props) => (
	<React.Fragment>
		<OidcProvider store={store} userManager={userManager}>
			<BasePage>
				<ActivationDetails />
			</BasePage>
		</OidcProvider>
	</React.Fragment>
);
export default ActivationPage;
