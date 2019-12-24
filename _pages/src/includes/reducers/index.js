import { combineReducers } from 'redux';
import { reducer as oidcReducer } from 'redux-oidc';
import activationReducer from 'Include/reducers/activation';
import contextReducer from 'Include/reducers/context';

const reducer = combineReducers({
	oidc: oidcReducer,
	activation: activationReducer,
	context: contextReducer,
});

export default reducer;