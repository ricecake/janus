import React from 'react';
import { default as MuiLink } from '@material-ui/core/Link';
import { Link as RouterLink } from 'react-router-dom';
import Button from '@material-ui/core/Button';

export const Hide = ({ If: condition, children }) => {
	if (condition) {
		return null;
	}

	return <React.Fragment>{children}</React.Fragment>;
};

export const Show = ({ If: condition, children }) => (
	<Hide If={!condition}> {children} </Hide>
);

const RrLink = (props) => {
	return /^https?:\/\//.test(props.to) ? (
		<a href={props.to} {...props} />
	) : (
		<RouterLink {...props} />
	);
};

export const Link = (props) => <MuiLink component={RrLink} {...props} />;

export const NavButton = (props) => <Button component={RrLink} {...props} />;
