import { PASSWORD_RESET_URL, post } from '../app';

const endpoint = 'password-reset';

const sendResetPasswordCodeViaEmail = body =>
  post(`${PASSWORD_RESET_URL}/${endpoint}`, body)
    .then(response => ({ response: response.data }))
    .catch(error => ({ error }));

const resetPasswordWithCode = body =>
  post(`${PASSWORD_RESET_URL}/${endpoint}`, body)
    .then(response => ({ response: response.data }))
    .catch(error => ({ error }));

export default {
  sendResetPasswordCodeViaEmail,
  resetPasswordWithCode,
};
