export const BEGIN_PASSWORD_RESET_REQUEST =
  'password-reset/BEGIN_PASSWORD_RESET_REQUEST';
export const BEGIN_PASSWORD_RESET_SUCCESS =
  'password-reset/BEGIN_PASSWORD_RESET_SUCCESS';
export const BEGIN_PASSWORD_RESET_FAILURE =
  'password-reset/BEGIN_PASSWORD_RESET_FAILURE';

const initialState = {
  loading: false,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case BEGIN_PASSWORD_RESET_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case BEGIN_PASSWORD_RESET_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case BEGIN_PASSWORD_RESET_FAILURE:
      return {
        ...state,
        loading: false,
      };
    default:
      return state;
  }
}

export function beginPasswordReset() {
  return {
    type: BEGIN_PASSWORD_RESET_REQUEST,
  };
}

export function beginPasswordResetSuccess() {
  return {
    type: BEGIN_PASSWORD_RESET_SUCCESS,
  };
}

export function beginPasswordResetFailure() {
  return {
    type: BEGIN_PASSWORD_RESET_FAILURE,
  };
}
