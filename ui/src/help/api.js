import { HELP_API_URL, post } from '../app';

const submitHelpRequest = body =>
  post(`${HELP_API_URL}/help`, body)
    .then(response => ({ response }))
    .catch(error => ({ error }));

export default {
  submitHelpRequest,
};
