import { call, put, select } from 'redux-saga/effects';
import { cloneableGenerator } from 'redux-saga/utils';
import { helpSubmission } from '../saga';
import helpApi from '../api';
import {
  nameSelector,
  emailAddressSelector,
  descriptionSelector,
} from '../selectors';
import {
  HELP_SUBMISSION_REQUEST,
  helpSubmissionSuccess,
  helpSubmissionFailure,
} from '../duck';

const sagaFactory = cloneableGenerator(helpSubmission)({
  type: HELP_SUBMISSION_REQUEST,
});

// THIS IS IMPORTANT TO UNDERSTAND FOR NEWCOMERS
// https://github.com/redux-saga/redux-saga/issues/114#issuecomment-184622588
// Anytime you use the 'return value' of yield, e.g
// const name = yield select(getName);
// you have to pass the mocked result to the next next() call

describe('help sagas', () => {
  it('should successfully create a help submission', () => {
    const saga = sagaFactory.clone();

    expect(saga.next().value).toEqual(select(nameSelector));
    expect(saga.next('Kevin Chen').value).toEqual(select(emailAddressSelector));
    expect(saga.next('fake.email@gmail.com').value).toEqual(
      select(descriptionSelector)
    );

    expect(saga.next("I'm having issues!").value).toEqual(
      call(helpApi.submitHelpRequest, {
        name: 'Kevin Chen',
        emailAddress: 'fake.email@gmail.com',
        description: "I'm having issues!",
        issueType: 'Browser Issues',
      })
    );

    expect(saga.next().value).toEqual(put(helpSubmissionSuccess()));

    // The saga should be finished now...
    expect(saga.next().done).toBeTruthy();
  });

  it('should unsuccessfully create a help submission', () => {
    const saga = sagaFactory.clone();

    expect(saga.next().value).toEqual(select(nameSelector));
    expect(saga.next('Kevin Chen').value).toEqual(select(emailAddressSelector));
    expect(saga.next('fake.email@gmail.com').value).toEqual(
      select(descriptionSelector)
    );

    expect(saga.next("I'm having issues!").value).toEqual(
      call(helpApi.submitHelpRequest, {
        name: 'Kevin Chen',
        emailAddress: 'fake.email@gmail.com',
        description: "I'm having issues!",
        issueType: 'Browser Issues',
      })
    );

    const errorResponse = {
      errors: [
        {
          code: 1,
          message: 'Shit went down',
        },
      ],
    };
    expect(saga.throw(errorResponse).value).toEqual(
      put(helpSubmissionFailure(errorResponse.errors))
    );

    // The saga should be finished now...
    expect(saga.next().done).toBeTruthy();
  });
});
