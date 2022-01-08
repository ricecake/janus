// This should have webauthn specific helpers.
// These should include a higher order component to wrap a needs-login component
// This also needs some helpers relating to if webauthn is supported.

export const webauthnCapable = () =>
	typeof window['PublicKeyCredential'] !== 'undefined';

export const strongWebauthnAvailable = () =>
	window.PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();

// Base64 to ArrayBuffer
const bufferDecode = (value) => {
	return Uint8Array.from(atob(value), (c) => c.charCodeAt(0));
};

// ArrayBuffer to URLBase64
const bufferEncode = (value) => {
	return btoa(String.fromCharCode.apply(null, new Uint8Array(value)))
		.replace(/\+/g, '-')
		.replace(/\//g, '_')
		.replace(/=/g, '');
};

export const doWebauthnRegister = () =>
	fetch('/webauthn/register/start', {
		method: 'POST',
	})
		.then((response) => response.json())
		.then((data) => {
			console.log(data);
			data.publicKey.challenge = bufferDecode(data.publicKey.challenge);
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

			return fetch('/webauthn/register/finish', {
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
			});
		});

export const doWebauthnLogin = (email) =>
	fetch(`/webauthn/login/start/${email}`, {
		method: 'POST',
	})
		.then((res) => res.json())
		.then((credentialRequestOptions) => {
			credentialRequestOptions.publicKey.challenge = bufferDecode(
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
			let authData = assertion.response.authenticatorData;
			let clientDataJSON = assertion.response.clientDataJSON;
			let rawId = assertion.rawId;
			let sig = assertion.response.signature;
			let userHandle = assertion.response.userHandle;

			return fetch(`/webauthn/login/finish/${email}`, {
				method: 'POST',
				headers: {
					'Content-type': 'application/json',
				},
				body: JSON.stringify({
					id: assertion.id,
					rawId: bufferEncode(rawId),
					type: assertion.type,
					response: {
						authenticatorData: bufferEncode(authData),
						clientDataJSON: bufferEncode(clientDataJSON),
						signature: bufferEncode(sig),
						userHandle: bufferEncode(userHandle),
					},
				}),
			});
		});
