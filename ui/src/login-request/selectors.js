import { NAME, FORM_NAME, FORM_FIELD_NAMES } from './constants';
import { formValueSelector } from 'redux-form';

const loadingSelector = state => state[NAME].loading;

const formSelector = formValueSelector(FORM_NAME);
const loginSelector = state => formSelector(state, FORM_FIELD_NAMES.login);
const passwordSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.password);

export { loadingSelector, loginSelector, passwordSelector };
