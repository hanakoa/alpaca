import { formValueSelector } from 'redux-form';
import { NAME as PASSWORD_RESET } from '../constants';
import { NAME, FIELD_NAMES } from './constants';

const loadingSelector = state => state[PASSWORD_RESET][NAME].loading;
const optionsSelector = state => state[PASSWORD_RESET][NAME].options;
const formSelector = formValueSelector(NAME);
const optionSelector = state => formSelector(state, FIELD_NAMES.option);

export { optionSelector, optionsSelector, loadingSelector };
