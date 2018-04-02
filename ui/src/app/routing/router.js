import React from 'react';
import {
  ConnectedRouter,
  routerReducer,
  routerMiddleware,
} from 'react-router-redux';
import { Redirect, Route, Switch } from 'react-router';
import createHistory from 'history/createBrowserHistory';
import { LoginNavContainer } from '../../login-nav';
import { HelpContainer } from '../../help';
import {
  BeginPasswordResetContainer,
  ConfirmPinResetContainer,
  SendPasswordResetContainer,
} from '../../password-reset';

const NotFound = () => <div>Not Found</div>;

const history = createHistory({
  basename: '/',
});

/**
 * The middleware and reducer are used to hook up the router to redux store
 * so that state is kept in sync with router
 */
export const middleware = routerMiddleware(history);
export const reducer = routerReducer;

/**
 * All application routes registered here
 */
export default (
  <ConnectedRouter history={history}>
    <Switch>
      <Route exact path="/" render={() => <Redirect to="/login" />} />
      <Route exact name="login" path="/login" component={LoginNavContainer} />
      <Route exact name="help" path="/help" component={HelpContainer} />
      <Route
        exact
        name="begin_password_reset"
        path="/account/begin_password_reset"
        component={BeginPasswordResetContainer}
      />
      <Route
        exact
        name="send_password_reset"
        path="/account/send_password_reset"
        component={SendPasswordResetContainer}
      />
      <Route
        exact
        name="confirm_pin_reset"
        path="/account/confirm_pin_reset"
        component={ConfirmPinResetContainer}
      />
      <Route component={NotFound} status={404} />
    </Switch>
  </ConnectedRouter>
);
