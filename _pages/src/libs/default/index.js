'use strict';

if (false){//'serviceWorker' in navigator) {
	navigator.serviceWorker.register('/static/worker.js', {scope: '/static/'})
	.then((reg) => {
		if(reg.installing) {
			console.log('Service worker installing');
		  } else if(reg.waiting) {
			console.log('Service worker installed');
		  } else if(reg.active) {
			console.log('Service worker active');
		  }
		console.log('Registration succeeded. Scope is ' + reg.scope);
		reg.update().then(console.log);
	}).catch((error) => {
		// registration failed
		console.log('Registration failed with ' + error);
	});
}