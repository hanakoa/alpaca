import axios from 'axios';

const isDevOrTest =
  process.env.NODE_ENV === 'development' || process.env.NODE_ENV === 'test';

const AUTH_API_URL = isDevOrTest
  ? 'http://localhost:8080'
  : window.alpacaConfig.auth.baseApiUrl;

const PASSWORD_RESET_URL = isDevOrTest
  ? 'http://localhost:8081'
  : window.alpacaConfig.passwordReset.baseApiUrl;

const HELP_API_URL = isDevOrTest
  ? 'http://localhost:8084/help'
  : window.alpacaConfig.help.baseApiUrl;

const REQUEST_ACCOUNT_API_URL = isDevOrTest
  ? 'http://localhost:8085/requestAccount'
  : window.alpacaConfig.accountRequests.baseApiUrl;

const post = (endpoint, body) => {
  return axios.post(endpoint, body);
};

export {
  post,
  AUTH_API_URL,
  PASSWORD_RESET_URL,
  HELP_API_URL,
  REQUEST_ACCOUNT_API_URL,
};
