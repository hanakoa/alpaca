export const SEND_CODE_REQUEST = 'mfa/SEND_CODE_REQUEST';
export const SEND_CODE_SUCCESS = 'mfa/SEND_CODE_SUCCESS';
export const SEND_CODE_FAILURE = 'mfa/SEND_CODE_FAILURE';

export const VERIFY_CODE_REQUEST = 'mfa/VERIFY_CODE_REQUEST';
export const VERIFY_CODE_SUCCESS = 'mfa/VERIFY_CODE_SUCCESS';
export const VERIFY_CODE_FAILURE = 'mfa/VERIFY_CODE_FAILURE';

const initialState = {
  loading: false,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case VERIFY_CODE_REQUEST:
    case SEND_CODE_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case VERIFY_CODE_SUCCESS:
    case SEND_CODE_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case VERIFY_CODE_FAILURE:
    case SEND_CODE_FAILURE:
      return {
        ...state,
        loading: false,
      };
    default:
      return state;
  }
}

export function sendCodeRequest() {
  return {
    type: SEND_CODE_REQUEST,
  };
}

export function sendCodeSuccess() {
  return {
    type: SEND_CODE_SUCCESS,
  };
}

export function sendCodeFailure() {
  return {
    type: SEND_CODE_FAILURE,
  };
}

export function verifyCodeRequest() {
  return {
    type: VERIFY_CODE_REQUEST,
  };
}

export function verifyCodeSuccess() {
  return {
    type: VERIFY_CODE_SUCCESS,
  };
}

export function verifyCodeFailure() {
  return {
    type: VERIFY_CODE_FAILURE,
  };
}
