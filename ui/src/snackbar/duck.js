export const SHOW_MESSAGE = 'snackbar/SHOW_MESSAGE';
export const CLEAR_MESSAGE = 'snackbar/CLEAR_MESSAGE';

const initialState = {
  message: '',
  open: false,
};

export default function snackbar(state = initialState, action) {
  switch (action.type) {
    case SHOW_MESSAGE:
      return {
        ...state,
        message: action.message,
        open: true,
      };
    case CLEAR_MESSAGE:
      return {
        ...state,
        message: '',
        open: false,
      };
    default:
      return state;
  }
}

export function showMessage(message) {
  return {
    type: SHOW_MESSAGE,
    message,
  };
}

export function clearMessage() {
  return {
    type: CLEAR_MESSAGE,
  };
}
