import React from 'react';
import BasePage from 'Component/BasePage';
import LoginForm from 'Component/LoginForm';

export const LoginPage = (props) => (
	<React.Fragment>
		<BasePage>
			<LoginForm />
		</BasePage>
	</React.Fragment>
);
export default LoginPage;
