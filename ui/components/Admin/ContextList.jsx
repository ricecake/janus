/* eslint-disable @typescript-eslint/no-use-before-define */
import React, { useEffect } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';

import { DataGrid } from '@material-ui/data-grid';
import { NewButton } from 'Component/Helpers';

import {
	fetchContexts,
	updateContext,
	createContext,
} from 'Include/reducers/admin';

const ContextListBasic = ({
	fetchContexts,
	updateContext,
	createContext,
	contexts,
}) => {
	useEffect(() => {
		fetchContexts();
	}, []);

	return (
		<>
			<DataGrid
				autoHeight={true}
				rows={contexts.map((ctx) => contextToRow(ctx))}
				columns={columns}
				onCellEditCommit={({ id, field, value }) => {
					let ctx = contexts.find(({ Code }) => Code === id);
					updateContext({ ...ctx, [field]: value });
				}}
			/>
			<NewButton
				title="New Context"
				placeholder="Context Name"
				onSubmit={(name) => {
					createContext({ Name: name });
				}}
			/>
		</>
	);
};

const contextToRow = ({ Code, Name, Description }) => ({
	id: Code,
	Code: Code,
	Name: Name,
	Description: Description,
});

const columns = [
	{ field: 'Code', flex: 1, headerName: 'Code', editable: false },
	{ field: 'Name', flex: 1, headerName: 'Name', editable: true },
	{
		field: 'Description',
		flex: 3,
		headerName: 'Description',
		editable: true,
	},
];

const stateToProps = ({ admin }) => ({ ...admin });
const dispatchToProps = (dispatch) =>
	bindActionCreators(
		{ fetchContexts, updateContext, createContext },
		dispatch
	);

export const ContextList = connect(
	stateToProps,
	dispatchToProps
)(ContextListBasic);
