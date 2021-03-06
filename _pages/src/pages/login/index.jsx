/* eslint-disable import/no-extraneous-dependencies */
import React from "react";
import ReactDOM from "react-dom";
import BasePage from 'Component/BasePage';
import store from 'Include/store';
import LoginForm from 'Component/LoginForm';
import { Provider } from 'react-redux';

ReactDOM.render((
	<Provider store={store}>
    <BasePage>
      <LoginForm />
    </BasePage>
	</Provider>
), document.getElementById('main'));
