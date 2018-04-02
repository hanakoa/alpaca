import { all } from 'redux-saga/effects';
import { sagas as loginSagas } from '../../login-request';
import { sagas as helpSagas } from '../../help';
import { sagas as requestAccountSagas } from '../../request-account';
import { sagas as passwordResetSagas } from '../../password-reset';

export default function* sagas() {
  yield all([
    ...loginSagas,
    ...helpSagas,
    ...requestAccountSagas,
    ...passwordResetSagas,
  ]);
}
