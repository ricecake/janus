import { combineReducers } from 'redux';
import { reducer as oidcReducer } from 'redux-oidc';
import contextReducer from 'Include/reducers/context';
import loginReducer from 'Include/reducers/login';
import signupReducer from 'Include/reducers/signup';
import homeReducer from 'Include/reducers/home';

const reducer = combineReducers({
	oidc: oidcReducer,
	context: contextReducer,
	login: loginReducer,
	signup: signupReducer,
	home: homeReducer,
});

export default reducer;
