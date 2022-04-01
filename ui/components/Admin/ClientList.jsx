import React from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { DataGrid } from '@material-ui/data-grid';

import SendIcon from '@material-ui/icons/Send';
import CloseIcon from '@material-ui/icons/Close';
import { makeStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import { Snackbar } from '@material-ui/core';
import { Alert } from '@material-ui/lab';

import {
	fetchContexts,
	fetchClients,
	updateClient,
} from 'Include/reducers/admin';

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

const LoginsBase = ({
	fetchContexts,
	fetchClients,
	updateClient,
	contexts,
	clients,
}) => {
	React.useEffect(() => {
		fetchContexts();
		fetchClients();
	}, []);

	const columns = [
		{
			field: 'Context',
			flex: 0.5,
			headerName: 'Application Group',
			editable: true,
			type: 'singleSelect',
			valueOptions: contexts.map((ctx) => ({
				value: ctx.Code,
				label: ctx.Name,
			})),
			valueFormatter: ({ value }) =>
				contexts.find((ctx) => ctx.Code === value)['Name'] || 'Unknown',
		},
		{ field: 'ClientId', flex: 1, headerName: 'ID', editable: false },
		{
			field: 'BaseUri',
			flex: 1.25,
			headerName: 'Base URI',
			editable: true,
		},
		{
			field: 'Secret',
			flex: 1,
			headerName: 'Client Secret',
			editable: true,
			valueFormatter: () => '********',
			renderEditCell: ({ id, value, api, field }) => {
				return (
					<EditSecret
						onSubmit={(newValue, event) => {
							api.setEditCellValue(
								{ id, field, value: newValue },
								event
							);
							api.commitCellChange({ id, field });
							api.setCellMode(id, field, 'view');
						}}
						onCancel={() => {
							api.setCellMode(id, field, 'view');
						}}
					/>
				);
			},
		},
		{ field: 'DisplayName', flex: 1, headerName: 'Name', editable: true },
		{
			field: 'Description',
			flex: 1,
			headerName: 'Description',
			editable: true,
		},
	];

	return (
		<DataGrid
			autoHeight={true}
			rows={clients.map((client) => clientToRow(client))}
			columns={columns}
			onCellEditCommit={({ id, field, value }) => {
				console.log({ id, field, value });
				let ctx = clients.find(({ ClientId }) => ClientId === id);
				updateClient({ ...ctx, [field]: value });
			}}
		/>
	);
};

export const EditSecret = ({ onSubmit, onCancel = () => {} }) => {
	const classes = useStyles();
	const [name, setName] = React.useState('');

	return (
		<Paper
			component="form"
			className={classes.root}
			onSubmit={(e) => {
				e.preventDefault();
				onSubmit(name, e);
			}}
		>
			<InputBase
				className={classes.input}
				placeholder={'New Secret'}
				onChange={(e) => setName(e.target.value)}
			/>
			<IconButton
				color="primary"
				className={classes.iconButton}
				onClick={() => {
					setName('');
					onCancel();
				}}
			>
				<CloseIcon />
			</IconButton>
			<Divider className={classes.divider} orientation="vertical" />
			<IconButton
				color="primary"
				className={classes.iconButton}
				type="submit"
			>
				<SendIcon />
			</IconButton>
		</Paper>
	);
};

const stateToProps = ({ admin }) => ({ ...admin });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{
			fetchClients,
			fetchContexts,
			updateClient,
		},
		dispatch
	);

export const ClientList = connect(stateToProps, dispatchToProps)(LoginsBase);
export default ClientList;
