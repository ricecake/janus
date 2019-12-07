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

setTimeout(()=>(
fetch("/login?response_type=code&state=tyu&prompt=any&redirect_uri=https%3A%2F%2Flogin.devhost.dev/static/code&scope=openid&client_id=KKw_TXyeSfOTg8E81D42xg").then(console.log)
), 2000);

var url = window.location.origin;
var settings = {
    authority: url,
    client_id: 'KKw_TXyeSfOTg8E81D42xg',
    //client_id: 'spa.short',
    redirect_uri: url + '/static/code',
    post_logout_redirect_uri: url + '/static/code',
    response_type: 'code',
    //response_mode: 'fragment',
    scope: 'openid profile api',
    //scope: 'openid profile api offline_access',
    
    popup_redirect_uri: url + '/static/code',
    popup_post_logout_redirect_uri: url + '/static/code',
    
    silent_redirect_uri: url + '/static/code',
    automaticSilentRenew:false,
    validateSubOnSilentRenew: true,
    //silentRequestTimeout:10000,

    monitorAnonymousSession : true,

    filterProtocolClaims: true,
    loadUserInfo: true,
    revokeAccessTokenOnSignout : true
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