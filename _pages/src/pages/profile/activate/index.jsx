/* eslint-disable import/no-extraneous-dependencies */
import React from "react";
import ReactDOM from "react-dom";
import BasePage from 'Component/BasePage';
import HelloWorld from 'Component/HelloWorld';
import Oidc from 'oidc-client';

ReactDOM.render((
	<BasePage>
		<HelloWorld />
	</BasePage>
), document.getElementById('main'));

var url = window.location.origin;
var settings = {
	authority: url,
	client_id: 'NR9eiBJ6SjO5v02lkx63Jw',
	response_type: 'code',
	scope: 'openid',
	silent_redirect_uri: url + '/static/oidc.html?mode=silent',
	automaticSilentRenew:true,
	validateSubOnSilentRenew: true,
	client_secret: "Example#1",
	loadUserInfo: false,
};
var mgr = new Oidc.UserManager(settings)
mgr.signinSilent({state:'some data'}).then(function(user) {
	mgr.getUser().then(function(user) {
		console.log("got user", user);
	}).catch(function(err) {
		console.log(err);
	});
}).catch(function(err) {
	console.log(err);
});
