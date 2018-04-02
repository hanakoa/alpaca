export const RESET_PASSWORD_WITH_CODE_REQUEST =
  'forgot-password/RESET_PASSWORD_WITH_CODE_REQUEST';
export const RESET_PASSWORD_WITH_CODE_SUCCESS =
  'forgot-password/RESET_PASSWORD_WITH_CODE_SUCCESS';
export const RESET_PASSWORD_WITH_CODE_FAILURE =
  'forgot-password/RESET_PASSWORD_WITH_CODE_FAILURE';

const initialState = {
  loading: false,
  options: null,
};

export default function auth(state = initialState, action) {
  switch (action.type) {
    case RESET_PASSWORD_WITH_CODE_REQUEST:
      return {
        ...state,
        loading: true,
      };
    case RESET_PASSWORD_WITH_CODE_SUCCESS:
      return {
        ...state,
        loading: false,
      };
    case RESET_PASSWORD_WITH_CODE_FAILURE:
      return {
        ...state,
        loading: false,
      };
    default:
      return state;
  }
}

export function resetPasswordWithCodeRequest() {
  return {
    type: RESET_PASSWORD_WITH_CODE_REQUEST,
  };
}

export function resetPasswordWithCodeSuccess() {
  return {
    type: RESET_PASSWORD_WITH_CODE_SUCCESS,
  };
}

export function resetPasswordWithCodeFailure() {
  return {
    type: RESET_PASSWORD_WITH_CODE_FAILURE,
  };
}
