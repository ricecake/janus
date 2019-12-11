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
    client_id: 'KKw_TXyeSfOTg8E81D42xg',
    response_type: 'code',
    scope: 'openid',
    silent_redirect_uri: url + '/static/code',
    automaticSilentRenew:true,
    validateSubOnSilentRenew: true,
};
var mgr = new Oidc.UserManager(settings)
mgr.signinSilent({state:'some data'}).then(function(user) {
	console.log("signed in", user);
}).catch(function(err) {
	console.log(err);
});
mgr.getUser().then(function(user) {
	console.log("got user", user);
}).catch(function(err) {
	console.log(err);
});

console.log(mgr);