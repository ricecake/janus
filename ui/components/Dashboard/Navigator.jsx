import React from 'react';
// import PropTypes from 'prop-types';
import clsx from 'clsx';
import { withStyles } from '@material-ui/core/styles';
import Divider from '@material-ui/core/Divider';
import Drawer from '@material-ui/core/Drawer';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemIcon from '@material-ui/core/ListItemIcon';
import ListItemText from '@material-ui/core/ListItemText';
import HomeIcon from '@material-ui/icons/Home';

import { NavLink } from 'react-router-dom';
import { makeStyles } from '@material-ui/styles';

const useStyles = makeStyles((theme) => ({
	categoryHeader: {
		paddingTop: theme.spacing(2),
		paddingBottom: theme.spacing(2),
		'&:hover,&:focus': {
			backgroundColor: 'rgba(255, 255, 255, 0.08)',
		},
	},
	categoryHeaderPrimary: {
		color: theme.palette.common.white,
	},
	item: {
		paddingTop: 1,
		paddingBottom: 1,
		color: 'rgba(255, 255, 255, 0.7)',
		outline: 'none',
		'text-decoration': 'none',
		'&:hover,&:focus': {
			backgroundColor: 'rgba(255, 255, 255, 0.08)',
		},
	},
	itemCategory: {
		backgroundColor: '#232f3e',
		boxShadow: '0 -1px 0 #404854 inset',
		paddingTop: theme.spacing(2),
		paddingBottom: theme.spacing(2),
	},
	firebase: {
		fontSize: 24,
		color: theme.palette.common.white,
	},
	itemActiveItem: {
		color: '#4fc3f7',
	},
	itemPrimary: {
		fontSize: 'inherit',
	},
	itemIcon: {
		minWidth: 'auto',
		marginRight: theme.spacing(2),
	},
	divider: {
		marginTop: theme.spacing(0),
	},
}));

function Navigator(props) {
	console.log(props);
	const classes = useStyles(props);
	const { title, categories, root = '', ...other } = props;

	return (
		<Drawer variant="permanent" {...other}>
			<List disablePadding>
				<ListItem className={clsx(classes.firebase)}>{title}</ListItem>
				<Divider className={classes.divider} />

				<NavLink
					exact
					to={`${root}/`.toLowerCase()}
					className={({ isActive }) =>
						[
							classes.item,
							isActive ? classes.itemActiveItem : '',
						].join(' ')
					}
				>
					<ListItem className={clsx(classes.categoryHeader)}>
						<ListItemIcon className={classes.itemIcon}>
							<HomeIcon />
						</ListItemIcon>
						<ListItemText
							classes={{
								primary: classes.itemPrimary,
							}}
						>
							Overview
						</ListItemText>
					</ListItem>
				</NavLink>
				<Divider className={classes.divider} />
				{categories.map(({ id, path, children = [] }) => (
					<React.Fragment key={id}>
						<NavLink
							to={`${root}/${path || id}/`.toLowerCase()}
							className={({ isActive }) =>
								[
									classes.item,
									isActive ? classes.itemActiveItem : '',
								].join(' ')
							}
						>
							<ListItem className={classes.categoryHeader}>
								<ListItemText>{id}</ListItemText>
							</ListItem>
						</NavLink>
						{children.map(
							({
								id: childId,
								path: childPath,
								icon,
								active,
							}) => (
								<NavLink
									key={childId}
									to={`${root}/${id}/${
										childPath || childId
									}/`.toLowerCase()}
									className={({ isActive }) =>
										[
											classes.item,
											isActive
												? classes.itemActiveItem
												: '',
										].join(' ')
									}
								>
									<ListItem key={childId} button>
										<ListItemIcon
											className={classes.itemIcon}
										>
											{icon}
										</ListItemIcon>
										<ListItemText
											classes={{
												primary: classes.itemPrimary,
											}}
										>
											{childId}
										</ListItemText>
									</ListItem>
								</NavLink>
							)
						)}
						<Divider className={classes.divider} />
					</React.Fragment>
				))}
			</List>
		</Drawer>
	);
}

// Navigator.propTypes = {
// 	classes: PropTypes.object.isRequired,
// };

export default Navigator;