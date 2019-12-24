/* eslint-disable import/no-extraneous-dependencies */
import React from "react";
import ReactDOM from "react-dom";
import BasePage from 'Component/BasePage';
import ActivationDetails from 'Component/ActivationDetails';


import { Provider } from 'react-redux';
import { OidcProvider } from 'redux-oidc';
import store from 'Include/store';
import userManager from 'Include/userManager';



ReactDOM.render((
	<Provider store={store}>
		<OidcProvider store={store} userManager={userManager}>
			<BasePage>
				<ActivationDetails />
			</BasePage>
		</OidcProvider>
	</Provider>
), document.getElementById('main'));
