import { call, put, takeEvery, select, fork } from 'redux-saga/effects';
import helpApi from './api';
import {
  HELP_SUBMISSION_REQUEST,
  helpSubmissionSuccess,
  helpSubmissionFailure,
} from './duck';
import {
  nameSelector,
  emailAddressSelector,
  descriptionSelector,
} from './selectors';

export function* helpSubmission(action) {
  const body = {
    name: yield select(nameSelector),
    emailAddress: yield select(emailAddressSelector),
    description: yield select(descriptionSelector),
    issueType: 'Browser Issues',
  };

  const { response, errors } = yield call(helpApi.submitHelpRequest, body);
  if (response) {
    yield put(helpSubmissionSuccess());
  } else {
    yield put(helpSubmissionFailure());
  }
}

export function* helpSubmissionSaga() {
  yield takeEvery(HELP_SUBMISSION_REQUEST, helpSubmission);
}

export const sagas = [fork(helpSubmissionSaga)];
