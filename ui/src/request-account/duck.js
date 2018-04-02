export const REQUEST_ACCOUNT_REQUEST =
  'request-account/REQUEST_ACCOUNT_REQUEST';
export const REQUEST_ACCOUNT_SUCCESS =
  'request-account/REQUEST_ACCOUNT_SUCCESS';
export const REQUEST_ACCOUNT_FAILURE =
  'request-account/REQUEST_ACCOUNT_FAILURE';

const initialState = {
  loading: false,
  success: false,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case REQUEST_ACCOUNT_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case REQUEST_ACCOUNT_SUCCESS:
      return {
        ...state,
        loading: false,
        success: true,
      };
    case REQUEST_ACCOUNT_FAILURE:
      return {
        ...state,
        loading: false,
        success: false,
      };
    default:
      return state;
  }
}

export function requestAccountRequest() {
  return {
    type: REQUEST_ACCOUNT_REQUEST,
  };
}

export function requestAccountSuccess() {
  return {
    type: REQUEST_ACCOUNT_SUCCESS,
  };
}

export function requestAccountFailure() {
  return {
    type: REQUEST_ACCOUNT_FAILURE,
  };
}
