import { call, put, takeEvery, select, fork } from 'redux-saga/effects';
import { delay } from 'redux-saga';
import requestAccountApi from './api';
import {
  REQUEST_ACCOUNT_REQUEST,
  requestAccountSuccess,
  requestAccountFailure,
} from './duck';
import {
  firstNameSelector,
  lastNameSelector,
  emailAddressSelector,
  citizenshipSelector,
  reasonSelector,
} from './selectors';
import { showMessage } from '../snackbar';

export function* requestAccount(action) {
  const body = {
    firstName: yield select(firstNameSelector),
    lastName: yield select(lastNameSelector),
    emailAddress: yield select(emailAddressSelector),
    citizenship: yield select(citizenshipSelector),
    reason: yield select(reasonSelector),
  };

  const { response, error } = yield call(
    requestAccountApi.requestAccount,
    body
  );

  if (response) {
    yield put(requestAccountSuccess());
    yield put(showMessage('Request Account succeeded!'));

    // We reset success to false after a short bit
    yield call(delay, 2000);
    yield put(requestAccountFailure());
  } else {
    yield put(showMessage('Request Account failed'));
    yield put(requestAccountFailure());
  }
}

export function* requestAccountSaga() {
  yield takeEvery(REQUEST_ACCOUNT_REQUEST, requestAccount);
}

export const sagas = [fork(requestAccountSaga)];
