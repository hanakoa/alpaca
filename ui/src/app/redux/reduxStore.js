import { applyMiddleware, compose, createStore } from 'redux';
import createSagaMiddleware from 'redux-saga';
import { routerMiddleware } from '../routing';
import { ReduxLogger } from '../tools';
import rootReducer from './rootReducer';
import sagas from './rootSaga';
import DevTools from '../tools/DevTools';

let store;
const initialState = {};
const sagaMiddleware = createSagaMiddleware();

function getStoreEnhancers() {
  switch (process.env.NODE_ENV) {
    case 'production':
      return compose(applyMiddleware(sagaMiddleware, routerMiddleware));
    default:
      //development, testing, etc...
      return compose(
        applyMiddleware(sagaMiddleware, routerMiddleware, ReduxLogger),
        DevTools.instrument()
      );
  }
}

export default function getStore() {
  if (!store) {
    store = createStore(rootReducer, initialState, getStoreEnhancers());
    sagaMiddleware.run(sagas);
  }
  return store;
}
