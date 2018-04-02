import reducer, {
  helpSubmissionRequest,
  helpSubmissionSuccess,
  helpSubmissionFailure,
} from '../duck';

describe('help reducer', () => {
  it('should return the initial state', () => {
    expect(reducer(undefined, {})).toEqual({
      loading: false,
    });
  });

  it('should handle HELP_SUBMISSION_REQUEST', () => {
    expect(
      reducer(
        {
          loading: false,
        },
        helpSubmissionRequest('/projects')
      )
    ).toEqual({
      loading: true,
    });
  });

  it('should handle LOGIN_SUCCESS', () => {
    expect(
      reducer(
        {
          loading: true,
        },
        helpSubmissionSuccess()
      )
    ).toEqual({
      loading: false,
    });
  });

  it('should handle LOGIN_FAILURE', () => {
    expect(
      reducer(
        {
          loading: true,
        },
        helpSubmissionFailure()
      )
    ).toEqual({
      loading: false,
    });
  });
});
