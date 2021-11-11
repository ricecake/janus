import deepmerge from 'deepmerge';

const common = {
	identity: {
		response_type: 'code',
		scope: 'openid profile',
		oidc_path: '/oauth',
		automaticSilentRenew: true,
		validateSubOnSilentRenew: true,
		loadUserInfo: false,
	},
};
const dev = {
	hosts: {
		idp_path: 'https://login.devhost.dev',
	},
	identity: {
		client_id: 'NR9eiBJ6SjO5v02lkx63Jw',
	},
};
const production = {
	hosts: {
		idp_path: 'https://login.greenstuff.io',
	},
	identity: {
		client_id: 'm1qCCDi0QJuEjcgZD2mx9Q',
	},
};

const MergedConfig = deepmerge.all([
	common,
	process.env.production ? production : dev,
]);

export default MergedConfig;
