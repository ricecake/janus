import { combineReducers } from 'redux';
import { reducer as oidcReducer } from 'redux-oidc';
import activationReducer from 'Include/reducers/activation';

const reducer = combineReducers(
  {
    oidc: oidcReducer,
    activation: activationReducer
  }
);

export default reducer;