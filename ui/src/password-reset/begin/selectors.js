import { formValueSelector } from 'redux-form';
import { NAME as PASSWORD_RESET } from '../constants';
import { NAME, FIELD_NAME } from './constants';

const loadingSelector = state => state[PASSWORD_RESET][NAME].loading;
const formSelector = formValueSelector(NAME);
const accountSelector = state => formSelector(state, FIELD_NAME);

export { accountSelector, loadingSelector };
