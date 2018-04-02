import { call, put, takeEvery, select, fork } from 'redux-saga/effects';
import api from './api';
import {
  SEND_CODE_REQUEST,
  sendCodeSuccess,
  sendCodeFailure,
  VERIFY_CODE_REQUEST,
  verifyCodeSuccess,
  verifyCodeFailure,
} from './duck';
import { sendCodeOptionSelector, codeSelector } from './selectors';

export function* sendCode(action) {
  const body = {
    option: yield select(sendCodeOptionSelector),
  };

  const { response, errors } = yield call(api.sendCode, body);
  if (response) {
    yield put(sendCodeSuccess());
  } else {
    yield put(sendCodeFailure());
  }
}

export function* sendCodeSaga() {
  yield takeEvery(SEND_CODE_REQUEST, sendCode);
}

export function* verifyCode(action) {
  const code = yield select(codeSelector);

  const body = {
    code,
  };

  const { response, errors } = yield call(api.verifyCode, body);

  if (response) {
    yield put(verifyCodeSuccess());
  } else {
    yield put(verifyCodeFailure());
  }
}

export function* verifyCodeSaga() {
  yield takeEvery(VERIFY_CODE_REQUEST, verifyCode);
}

export const sagas = [fork(verifyCodeSaga), fork(sendCodeSaga)];
