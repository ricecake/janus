import React from 'react';
import { LoginBasePage } from 'Component/BasePage';
import HomeAppMenu from 'Component/HomePage';

export const HomePage = () => (
	<React.Fragment>
		<LoginBasePage>
			<HomeAppMenu />
		</LoginBasePage>
	</React.Fragment>
);
export default HomePage;
