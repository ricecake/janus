import React, { useEffect } from 'react';
import BasePage from 'Component/BasePage';

import { OidcProvider } from 'redux-oidc';
import store from 'Include/store';
import userManager, { withLogin } from 'Include/userManager';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import {
	changeName,
	changePassword,
	changePasswordVerifier,
	submitForm,
	startSignin,
} from 'Include/reducers/activation';

// Base64 to ArrayBuffer
function bufferDecode(value) {
	return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
}

// ArrayBuffer to URLBase64
function bufferEncode(value) {
	return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
		.replace(/\+/g, '-')
		.replace(/\//g, '_')
		.replace(/=/g, '');
}

export const WebauthnPage = withLogin((props) => {
	console.log(props);

	useEffect(() => {
		fetch('/webauthn/register/start', {
			method: 'POST',
		})
			.then((response) => response.json())
			.then((data) => {
				console.log(data);
				data.publicKey.challenge = bufferDecode(
					data.publicKey.challenge
				);
				data.publicKey.user.id = bufferDecode(data.publicKey.user.id);

				if (data.publicKey.excludeCredentials) {
					data.publicKey.excludeCredentials.forEach(
						(cred, index, excludes) => {
							excludes[index].id = bufferDecode(cred.id);
						}
					);
				}

				return navigator.credentials.create(data);
			})
			.then((data) => {
				let attestationObject = data.response.attestationObject;
				let clientDataJSON = data.response.clientDataJSON;
				let rawId = data.rawId;

				fetch('/webauthn/register/finish', {
					method: 'POST',
					body: JSON.stringify({
						id: data.id,
						rawId: bufferEncode(rawId),
						type: data.type,
						response: {
							attestationObject: bufferEncode(attestationObject),
							clientDataJSON: bufferEncode(clientDataJSON),
						},
					}),
				}).then(() => {
					fetch('/webauthn/login/start/geoffcake@gmail.com', {
						method: 'POST',
					})
						.then((res) => res.json())
						.then((credentialRequestOptions) => {
							credentialRequestOptions.publicKey.challenge =
								bufferDecode(
									credentialRequestOptions.publicKey.challenge
								);
							credentialRequestOptions.publicKey.allowCredentials.forEach(
								function (listItem) {
									listItem.id = bufferDecode(listItem.id);
								}
							);

							return navigator.credentials.get({
								publicKey: credentialRequestOptions.publicKey,
							});
						})
						.then((assertion) => {
							console.log(assertion);
							let authData = assertion.response.authenticatorData;
							let clientDataJSON =
								assertion.response.clientDataJSON;
							let rawId = assertion.rawId;
							let sig = assertion.response.signature;
							let userHandle = assertion.response.userHandle;

							return fetch(
								'/webauthn/login/finish/geoffcake@gmail.com',
								{
									method: 'POST',
									headers: {
										'Content-type': 'application/json',
									},
									body: JSON.stringify({
										id: assertion.id,
										rawId: bufferEncode(rawId),
										type: assertion.type,
										response: {
											authenticatorData:
												bufferEncode(authData),
											clientDataJSON:
												bufferEncode(clientDataJSON),
											signature: bufferEncode(sig),
											userHandle:
												bufferEncode(userHandle),
										},
									}),
								}
							);
						})
						.then((success) => {
							alert('successfully logged in !');
							return;
						});
				});
			});
	});

	return (
		<React.Fragment>
			<OidcProvider store={store} userManager={userManager}>
				<BasePage>Cactus Toaster</BasePage>
			</OidcProvider>
		</React.Fragment>
	);
});

const stateToProps = ({ activation, oidc }) => ({
	...activation,
	user: oidc.user,
});
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{
			changeName,
			changePassword,
			changePasswordVerifier,
			submitForm,
			startSignin,
		},
		dispatch
	);

export default connect(stateToProps, dispatchToProps)(WebauthnPage);
