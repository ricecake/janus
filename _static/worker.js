self.addEventListener('fetch', (event) => {
	const url = new URL(event.request.url);
	if (url.pathname == '/static/code') {
		const query = {};
		for (const [key, value] of url.searchParams.entries()) {
			query[key] = value;
		}
		return event.respondWith(
			new Response(JSON.stringify(query))
		);
	}
	event.respondWith(
		caches.match(event.request).then((response) => {
			return response || fetch(event.request);
		})
	);
});

//let x =await fetch("/login?response_type=code&state=test1&prompt=any&redirect_uri=https%3A%2F%2Flogin.devhost.dev/profile/code&scope=openid&client_id=KKw_TXyeSfOTg8E81D42xg"); x.url