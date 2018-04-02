import { combineReducers } from 'redux';
import { constants as beginConstants, reducer as beginReducer } from './begin';
import {
  constants as confirmPinConstants,
  reducer as confirmPinReducer,
} from './confirm-pin-reset';
import { constants as sendConstants, reducer as sendReducer } from './send';

export default combineReducers({
  [beginConstants.NAME]: beginReducer,
  [confirmPinConstants.NAME]: confirmPinReducer,
  [sendConstants.NAME]: sendReducer,
});
