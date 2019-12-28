import { createActions, handleActions } from 'redux-actions';

const defaultState = {};

export const {} = createActions({}, { prefix: "janus/login" });

export const initiateLogin = () => (dispatch, getState) => {};

const reducer = handleActions({}, defaultState);

export default reducer;