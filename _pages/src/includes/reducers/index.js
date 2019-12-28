import { combineReducers } from 'redux';
import { reducer as oidcReducer } from 'redux-oidc';
import activationReducer from 'Include/reducers/activation';
import contextReducer from 'Include/reducers/context';
import loginReducer from 'Include/reducers/login';
import signupReducer from 'Include/reducers/signup';

const reducer = combineReducers({
	oidc: oidcReducer,
	activation: activationReducer,
	context: contextReducer,
	login: loginReducer,
	signup: signupReducer,
});

export default reducer;