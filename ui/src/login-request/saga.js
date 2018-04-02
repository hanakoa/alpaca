import { call, put, takeEvery, select, fork } from 'redux-saga/effects';
import loginApi from './api';
import { LOGIN_REQUEST, loginSuccess, loginFailure } from './duck';
import { loginSelector, passwordSelector } from './selectors';
import { showMessage } from '../snackbar';

export function* login(action) {
  const { payload: { redirectTo } } = action;

  const body = {
    login: yield select(loginSelector),
    password: yield select(passwordSelector),
  };

  const { response, errors } = yield call(loginApi.login, body);

  if (response) {
    yield put(loginSuccess(response));
    yield put(showMessage(JSON.stringify(response)));
  } else {
    yield put(loginFailure());
    yield put(showMessage('Login failed'));
  }
}

export function* loginSaga() {
  yield takeEvery(LOGIN_REQUEST, login);
}

export const sagas = [fork(loginSaga)];
