import { HELP_API_URL, post } from '../app';

// TODO fix URL and endpoints
const sendCode = body =>
  post(`${HELP_API_URL}/send`, body)
    .then(response => ({ response }))
    .catch(error => ({ error }));

const verifyCode = body =>
  post(`${HELP_API_URL}/verify`, body)
    .then(response => ({ response }))
    .catch(error => ({ error }));

export default {
  sendCode,
  verifyCode,
};
