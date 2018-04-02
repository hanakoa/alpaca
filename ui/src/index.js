import React from 'react';
import ReactDOM from 'react-dom';
import { getStore } from './app';
import { Provider } from 'react-redux';
import AppContainer from './app/AppContainer';

ReactDOM.render(
  <Provider store={getStore()}>
    <AppContainer />
  </Provider>,
  document.getElementById('root')
);
