import { NAME, FORM_NAME, FORM_FIELD_NAMES } from './constants';
import { formValueSelector } from 'redux-form';

const loadingSelector = state => state[NAME].loading;

const formSelector = formValueSelector(FORM_NAME);
const nameSelector = state => formSelector(state, FORM_FIELD_NAMES.name);
const emailAddressSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.emailAddress);
const descriptionSelector = state =>
  formSelector(state, FORM_FIELD_NAMES.description);

export {
  loadingSelector,
  nameSelector,
  emailAddressSelector,
  descriptionSelector,
};
