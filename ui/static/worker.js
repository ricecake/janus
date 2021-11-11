/* This is a start of a proof of concept for handling redirects in a code flow.
Since the auth call in the background would return a redirect to the callback url, which would then return data, but couldn't actually reach out and grab the token, it's instead possible to install a webworker that will intercept the redirected call (maybe), and then exchange the code for a token.  Allowing both a code flow, and a spa without a background server to handle that part.

*/
self.addEventListener('fetch', (event) => {
	const url = new URL(event.request.url);
	if (url.pathname == '/static/code') {
		const query = {};
		for (const [key, value] of url.searchParams.entries()) {
			query[key] = value;
		}
		return event.respondWith(new Response(JSON.stringify(query)));
	}
	event.respondWith(
		caches.match(event.request).then((response) => {
			return response || fetch(event.request);
		})
	);
});

//let x =await fetch("/login?response_type=code&state=test1&prompt=any&redirect_uri=https%3A%2F%2Flogin.devhost.dev/profile/code&scope=openid&client_id=KKw_TXyeSfOTg8E81D42xg"); x.url
