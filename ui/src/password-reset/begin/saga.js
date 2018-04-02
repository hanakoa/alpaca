import { call, put, takeEvery, select } from 'redux-saga/effects';
import { push } from 'react-router-redux';
import passwordApi from '../api';
import { accountSelector } from './selectors';
import {
  BEGIN_PASSWORD_RESET_REQUEST,
  beginPasswordResetFailure,
  beginPasswordResetSuccess,
} from './duck';
import { updateOptions } from '../send';
import { showMessage } from '../../snackbar';

export function* beginPasswordReset(action) {
  const body = {
    account: yield select(accountSelector),
  };

  const { response } = yield call(
    passwordApi.sendResetPasswordCodeViaEmail,
    body
  );

  if (response) {
    yield put(beginPasswordResetSuccess());
    let hasOptions = !!response.email_addresses;
    if (hasOptions) {
      console.log('UPDATING SEND OPTIONS:', response);
      yield put(updateOptions(transformResponse(response)));
      yield put(push('/account/send_password_reset'));
    } else {
      yield put(showMessage('Request sent to that account (if it exists)'));
      yield put(push('/'));
    }
  } else {
    yield put(showMessage('Request sent to that account (if it exists)'));
    yield put(beginPasswordResetFailure());
  }
}

function transformResponse(response) {
  let arr = [];
  if (response.email_addresses) {
    response.email_addresses.forEach(e =>
      arr.push({ type: 'email_address', id: e.id, value: e.email_address })
    );
  }
  if (response.phone_numbers) {
    response.phone_numbers.forEach(p =>
      arr.push({ type: 'phone_number', id: p.id, value: p.phone_number })
    );
  }
  return arr;
}

export function* beginPasswordResetSaga() {
  yield takeEvery(BEGIN_PASSWORD_RESET_REQUEST, beginPasswordReset);
}
