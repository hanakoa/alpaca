import { combineReducers } from 'redux';
import { routerReducer } from 'react-router-redux';
import { reducer as formReducer } from 'redux-form';
import {
  constants as loginRequestConstants,
  reducer as loginRequestReducer,
} from '../../login-request';
import { constants as helpConstants, reducer as helpReducer } from '../../help';
import {
  constants as requestAccountConstants,
  reducer as requestAccountReducer,
} from '../../request-account';
import {
  constants as passwordResetConstants,
  reducer as passwordResetReducer,
} from '../../password-reset';
import {
  constants as snackbarConstants,
  reducer as snackbarReducer,
} from '../../snackbar';

export default combineReducers({
  router: routerReducer,
  form: formReducer,
  [snackbarConstants.NAME]: snackbarReducer,
  [loginRequestConstants.NAME]: loginRequestReducer,
  [helpConstants.NAME]: helpReducer,
  [requestAccountConstants.NAME]: requestAccountReducer,
  [passwordResetConstants.NAME]: passwordResetReducer,
});
