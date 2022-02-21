import React from 'react';
import Grid from '@material-ui/core/Grid';

import ProfilePage from './frame';

import { Webauthn } from './Webauthn';
import { Password } from './Password';
import { webauthnCapable } from 'Include/webauthn';
import { Show } from 'Component/Helpers';

const Authentication = () => {
	return (
		<ProfilePage>
			<Show If={webauthnCapable()}>
				<Grid item>
					<Webauthn />
				</Grid>
			</Show>
			<Grid item>
				<Password />
			</Grid>
		</ProfilePage>
	);
};
export default Authentication;
