import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { makeStyles } from '@material-ui/core/styles';

import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListSubheader from '@material-ui/core/ListSubheader';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Collapse from '@material-ui/core/Collapse';
import ListItemText from '@material-ui/core/ListItemText';

import { DataGrid } from '@material-ui/data-grid';

import { fetchContexts, fetchClients } from 'Include/reducers/admin';

const useStyles = makeStyles((theme) => ({
	root: {
		padding: '2px 4px',
		display: 'flex',
		alignItems: 'center',
		width: 400,
	},
	input: {
		marginLeft: theme.spacing(1),
		flex: 1,
	},
	iconButton: {
		padding: 10,
	},
	divider: {
		height: 28,
		margin: 4,
	},
}));

const clientToRow = ({
	Context,
	ClientId,
	BaseUri,
	DisplayName,
	Description,
}) => ({
	id: ClientId,
	Context: Context,
	ClientId: ClientId,
	BaseUri: BaseUri,
	DisplayName: DisplayName,
	Description: Description,
});

const AppGroup = ({ Name, Description, clients, contexts }) => {
	const [open, setOpen] = React.useState(false);
	const columns = [
		{
			field: 'Context',
			flex: 1,
			headerName: 'Application Group',
			editable: true,
			type: 'singleSelect',
			valueOptions: contexts.map((ctx) => ({
				value: ctx.Code,
				label: ctx.Name,
			})),
			valueGetter: ({ value, ...rest }) => {
				console.log(rest);
				return contexts.find((ctx) => ctx.Code === value).Name;
			},
		},
		{ field: 'ClientId', flex: 1, headerName: 'ID', editable: false },
		{
			field: 'BaseUri',
			flex: 1,
			headerName: 'Base URI',
			editable: false,
		},
		{ field: 'DisplayName', flex: 1, headerName: 'Name', editable: false },
		{
			field: 'Description',
			flex: 1,
			headerName: 'Description',
			editable: false,
		},
	];

	return (
		<>
			<ListItem button onClick={() => setOpen(!open)}>
				<ListItemText primary={Name} secondary={Description} />
				{open ? <ExpandLess /> : <ExpandMore />}
			</ListItem>
			<Collapse in={open} timeout="auto" unmountOnExit>
				<List>
					<ListSubheader>Clients</ListSubheader>
					<DataGrid
						autoHeight={true}
						rows={clients.map((client) => clientToRow(client))}
						columns={columns}
					/>
				</List>
			</Collapse>
		</>
	);
};

const LoginsBase = ({ fetchContexts, fetchClients, contexts, clients }) => {
	const classes = useStyles();

	React.useEffect(() => {
		fetchContexts();
		fetchClients();
	}, []);

	return (
		<List>
			<ListSubheader>Application Groups</ListSubheader>
			{contexts.map((ctx) => (
				<AppGroup
					key={ctx['Code']}
					contexts={contexts}
					clients={clients.filter(
						(client) => client.Context === ctx['Code']
					)}
					{...ctx}
				/>
			))}
		</List>
	);
};

const stateToProps = ({ admin }) => ({ ...admin });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{
			fetchClients,
			fetchContexts,
		},
		dispatch
	);

export const ClientList = connect(stateToProps, dispatchToProps)(LoginsBase);
export default ClientList;
