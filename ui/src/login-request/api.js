import { AUTH_API_URL, post } from '../app';

const login = body =>
  post(`${AUTH_API_URL}/token`, body)
    .then(response => ({ response: response.data }))
    .catch(error => ({ error }));

export default {
  login,
};
