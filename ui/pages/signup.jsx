import React from 'react';
import BasePage from 'Component/BasePage';
import SignupForm from 'Component/SignupForm';

export const SignupPage = (props) => (
	<React.Fragment>
		<BasePage>
			<SignupForm />
		</BasePage>
	</React.Fragment>
);
export default SignupPage;
