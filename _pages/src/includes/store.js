import { createStore, applyMiddleware, compose } from "redux";
import { loadUser } from "redux-oidc";
import reducer from "Include/reducers";
import userManager from "Include/userManager";

// create the middleware with the userManager
// const oidcMiddleware = createOidcMiddleware(userManager);

const loggerMiddleware = store => next => action => {
  console.log("Action type:", action.type);
  console.log("Action payload:", action.payload);
  console.log("State before:", store.getState());
  next(action);
  console.log("State after:", store.getState());
};

const initialState = {};

const createStoreWithMiddleware = compose(
  applyMiddleware(loggerMiddleware)
)(createStore);

const store = createStoreWithMiddleware(reducer, initialState);
loadUser(store, userManager);

export default store;