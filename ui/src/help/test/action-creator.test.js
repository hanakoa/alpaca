import * as types from '../duck';
import {
  helpSubmissionRequest,
  helpSubmissionSuccess,
  helpSubmissionFailure,
} from '../duck';

describe('help actions', () => {
  it('should create an action to request a help submission', () => {
    expect(helpSubmissionRequest('/projects')).toEqual({
      type: types.HELP_SUBMISSION_REQUEST,
    });
  });

  it('should create an action to signal help submission success', () => {
    expect(helpSubmissionSuccess()).toEqual({
      type: types.HELP_SUBMISSION_SUCCESS,
    });
  });

  it('should create an action to signal help submission failure', () => {
    expect(
      helpSubmissionFailure([
        {
          code: 1,
          message: 'An error occurred.',
        },
        {
          code: 2,
          message: 'A different error occurred.',
        },
      ])
    ).toEqual({
      type: types.HELP_SUBMISSION_FAILURE,
      errors: [
        {
          code: 1,
          message: 'An error occurred.',
        },
        {
          code: 2,
          message: 'A different error occurred.',
        },
      ],
    });
  });
});
