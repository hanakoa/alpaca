export const SEND_PASSWORD_RESET_REQUEST =
  'password-reset/SEND_PASSWORD_RESET_REQUEST';
export const SEND_PASSWORD_RESET_SUCCESS =
  'password-reset/SEND_PASSWORD_RESET_SUCCESS';
export const SEND_PASSWORD_RESET_FAILURE =
  'password-reset/SEND_PASSWORD_RESET_FAILURE';
export const UPDATE_SEND_OPTIONS = 'password-reset/UPDATE_SEND_OPTIONS';

const initialState = {
  loading: false,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case SEND_PASSWORD_RESET_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case SEND_PASSWORD_RESET_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case SEND_PASSWORD_RESET_FAILURE:
      return {
        ...state,
        loading: false,
      };
    case UPDATE_SEND_OPTIONS:
      return {
        ...state,
        options: action.options,
      };
    default:
      return state;
  }
}

export function updateOptions(options) {
  return {
    type: UPDATE_SEND_OPTIONS,
    options,
  };
}

export function sendPasswordResetRequest() {
  return {
    type: SEND_PASSWORD_RESET_REQUEST,
  };
}

export function sendPasswordResetSuccess() {
  return {
    type: SEND_PASSWORD_RESET_SUCCESS,
  };
}

export function sendPasswordResetFailure() {
  return {
    type: SEND_PASSWORD_RESET_FAILURE,
  };
}
