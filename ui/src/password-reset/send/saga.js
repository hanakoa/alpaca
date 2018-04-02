import { call, put, takeEvery, select } from 'redux-saga/effects';
import passwordApi from '../api';
import {
  SEND_PASSWORD_RESET_REQUEST,
  sendPasswordResetSuccess,
  sendPasswordResetFailure,
} from './duck';
import { optionsSelector } from './selectors';
import { showMessage } from '../../snackbar';

export function* sendPasswordResetCode(action) {
  const option = yield select(optionsSelector);

  const body = {
    option,
  };

  const { response } = yield call(passwordApi.resetPasswordWithCode, body);
  if (response) {
    yield put(sendPasswordResetSuccess());
  } else {
    yield put(sendPasswordResetFailure());
  }
}

export function* sendPasswordResetCodeSaga() {
  yield takeEvery(SEND_PASSWORD_RESET_REQUEST, sendPasswordResetCode);
}
