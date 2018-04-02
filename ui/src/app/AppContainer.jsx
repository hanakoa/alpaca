import React, { Component } from 'react';
import './styles/main.css';
import { router } from './routing';
import { DevTools } from './tools';
import { Snackbar } from '../snackbar';

class AppContainer extends Component {
  render() {
    return (
      <div>
        <Snackbar />
        <div className="card mx-auto mt-4">
          {process.env.NODE_ENV === 'development' && <DevTools />}
          <div className="card-block">{router}</div>
        </div>
      </div>
    );
  }
}

export default AppContainer;
