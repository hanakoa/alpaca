import { NAME, FORM_NAME, FORM_FIELD_NAMES } from './constants';
import { formValueSelector } from 'redux-form';

const loadingSelector = state => state[NAME].loading;
const successSelector = state => state[NAME].success;

const formSelector = formValueSelector(FORM_NAME);
const firstNameSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.firstName);
const lastNameSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.lastName);
const emailAddressSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.emailAddress);
const citizenshipSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.citizenship);
const reasonSelector = state => formSelector(state, FORM_FIELD_NAMES.reason);

export {
  loadingSelector,
  successSelector,
  firstNameSelector,
  lastNameSelector,
  emailAddressSelector,
  citizenshipSelector,
  reasonSelector,
};
