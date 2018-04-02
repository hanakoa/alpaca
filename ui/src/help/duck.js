export const HELP_SUBMISSION_REQUEST = 'help/HELP_SUBMISSION_REQUEST';
export const HELP_SUBMISSION_SUCCESS = 'help/HELP_SUBMISSION_SUCCESS';
export const HELP_SUBMISSION_FAILURE = 'help/HELP_SUBMISSION_FAILURE';

const initialState = {
  loading: false,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case HELP_SUBMISSION_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case HELP_SUBMISSION_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case HELP_SUBMISSION_FAILURE:
      return {
        ...state,
        loading: false,
      };
    default:
      return state;
  }
}

export function helpSubmissionRequest() {
  return {
    type: HELP_SUBMISSION_REQUEST,
  };
}

export function helpSubmissionSuccess() {
  return {
    type: HELP_SUBMISSION_SUCCESS,
  };
}

export function helpSubmissionFailure() {
  return {
    type: HELP_SUBMISSION_FAILURE,
  };
}
