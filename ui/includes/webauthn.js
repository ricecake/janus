import {
	browserSupportsWebAuthn,
	startAuthentication,
	startRegistration,
	browserSupportsWebAuthnAutofill,
} from '@simplewebauthn/browser';

// This should have webauthn specific helpers.
// These should include a higher order component to wrap a needs-login component
// This also needs some helpers relating to if webauthn is supported.

export const webauthnCapable = browserSupportsWebAuthn;

export const strongWebauthnAvailable = () =>
	window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();

const handleFetchError = (res) => {
	if (!res.ok) {
		throw res;
	}
	return res;
};

export const doWebauthnRegister = (name) =>
	fetch('/webauthn/register/start', {
		method: 'POST',
	})
		.then(handleFetchError)
		.then((response) => response.json())
		.then((data) => {
			return startRegistration(data.publicKey);
		})
		.then((data) => {
			return fetch(`/webauthn/register/finish?name=${name}`, {
				method: 'POST',
				body: JSON.stringify(data),
			}).then(handleFetchError);
		});

export const doWebauthnLogin = (email, client_id) =>
	fetch(`/webauthn/login/start/${client_id}/${email}`, {
		method: 'POST',
	})
		.then(handleFetchError)
		.then((res) => res.json())
		.then((credentialRequestOptions) => {
			return startAuthentication(credentialRequestOptions.publicKey);
		})
		.then((assertion) => {
			return fetch(`/webauthn/login/finish/${client_id}/${email}`, {
				method: 'POST',
				headers: {
					'Content-type': 'application/json',
				},
				body: JSON.stringify(assertion),
			}).then(handleFetchError);
		});

export const doMediatedWebauthn = (client_id, autofil = false) => {
	return browserSupportsWebAuthnAutofill().then((supportAutofill) => {
		if (supportAutofill || !autofil) {
			return fetch(`/webauthn/mediated/start/${client_id}`, {
				method: 'POST',
			})
				.then(handleFetchError)
				.then((res) => res.json())
				.then((credentialRequestOptions) => {
					return startAuthentication(
						credentialRequestOptions.publicKey,
						autofil
					);
				})
				.then((assertion) => {
					return fetch(`/webauthn/mediated/finish/${client_id}`, {
						method: 'POST',
						headers: {
							'Content-type': 'application/json',
						},
						body: JSON.stringify(assertion),
					}).then(handleFetchError);
				});
		} else {
			return Promise.reject();
		}
	});
};
