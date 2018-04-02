import { NAME, FORM_NAMES, FORM_FIELD_NAMES } from './constants';
import { formValueSelector } from 'redux-form';

const loadingSelector = state => state[NAME].loading;

const sendCodeFormSelector = formValueSelector(FORM_NAMES.sendCodeForm);
const verifyCodeFormSelector = formValueSelector(FORM_NAMES.verifyCodeForm);

const sendCodeOptionSelector = state =>
  sendCodeFormSelector(state, FORM_FIELD_NAMES.sendCodeOption);

const codeSelector = state =>
  verifyCodeFormSelector(state, FORM_FIELD_NAMES.code);

export { loadingSelector, sendCodeOptionSelector, codeSelector };
