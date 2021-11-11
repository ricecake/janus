import React from 'react';
import { default as MuiLink } from '@material-ui/core/Link';
import { Link as RrLink } from 'react-router-dom';

export const Hide = ({ If: condition, children }) => {
	if (condition) {
		return null;
	}

	return <React.Fragment>{children}</React.Fragment>;
};

export const Show = ({ If: condition, children }) => (
	<Hide If={!condition}> {children} </Hide>
);

export const Link = (props) => <MuiLink component={RrLink} {...props} />;
