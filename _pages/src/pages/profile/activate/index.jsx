/* eslint-disable import/no-extraneous-dependencies */
import React from "react";
import ReactDOM from "react-dom";
import BasePage from 'Component/BasePage';
import ActivationDetails from 'Component/ActivationDetails';
import Oidc from 'oidc-client';



var serverVars = JSON.parse(document.getElementById('openid-client-params').innerHTML);

var url = window.location.origin;
var settings = {
	authority: url,
	response_type: 'code',
	scope: 'openid',
	silent_redirect_uri: url + '/static/oidc.html?mode=silent',
	automaticSilentRenew:true,
	validateSubOnSilentRenew: true,
	loadUserInfo: false,
	... serverVars
};

var mgr = new Oidc.UserManager(settings)
ReactDOM.render((
	<BasePage>
		<ActivationDetails userManager={ mgr } />
	</BasePage>
), document.getElementById('main'));
