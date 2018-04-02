import { fork } from 'redux-saga/effects';
import { beginPasswordResetSaga } from './begin/saga';
import { resetPasswordWithCodeSaga } from './confirm-pin-reset/saga';
import { sendPasswordResetCodeSaga } from './send/saga';

export const sagas = [
  fork(resetPasswordWithCodeSaga),
  fork(beginPasswordResetSaga),
  fork(sendPasswordResetCodeSaga),
];
