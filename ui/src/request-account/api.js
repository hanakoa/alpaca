import { REQUEST_ACCOUNT_API_URL, post } from '../app';

const requestAccount = body =>
  post(`${REQUEST_ACCOUNT_API_URL}/requestAccount`, body)
    .then(response => ({ response }))
    .catch(error => ({ error }));

export default {
  requestAccount,
};
