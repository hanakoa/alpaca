import { call, put, takeEvery, select } from 'redux-saga/effects';
import passwordApi from '../api';
import {
  RESET_PASSWORD_WITH_CODE_REQUEST,
  resetPasswordWithCodeSuccess,
  resetPasswordWithCodeFailure,
} from './duck';
import { codeSelector } from './selectors';
import { showMessage } from '../../snackbar';

export function* resetPasswordWithCode(action) {
  const code = yield select(codeSelector);

  const body = {
    code,
  };

  const { response, errors } = yield call(
    passwordApi.resetPasswordWithCode,
    body
  );
  if (response) {
    yield put(resetPasswordWithCodeSuccess());
  } else {
    yield put(resetPasswordWithCodeFailure());
  }
}

export function* resetPasswordWithCodeSaga() {
  yield takeEvery(RESET_PASSWORD_WITH_CODE_REQUEST, resetPasswordWithCode);
}
