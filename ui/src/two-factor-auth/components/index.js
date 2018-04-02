import React from 'react';
import { maxLength, renderField, required } from '../../redux-form-helpers';
import { FORM_NAMES, FORM_FIELD_NAMES } from '../constants';
import { Field, reduxForm } from 'redux-form';
import { OneFieldOneButtonForm } from '../../app';

const SendCodeForm = reduxForm({
  form: FORM_NAMES.sendCodeForm,
})(({ handleSubmit }) => (
  <OneFieldOneButtonForm
    header="2-Factor Authentication"
    buttonText="Send Code"
    handleSubmit={handleSubmit}>
    <Field
      name={FORM_FIELD_NAMES.sendCodeOption}
      fieldName={FORM_FIELD_NAMES.sendCodeOption}
      label="Code transmission option: SMS or email?"
      validate={[required, maxLength(255)]}
      component={renderField}
    />
  </OneFieldOneButtonForm>
));

const VerifyCodeForm = reduxForm({
  form: FORM_NAMES.verifyCodeForm,
})(({ handleSubmit }) => (
  <OneFieldOneButtonForm
    header=""
    buttonText="Verify"
    handleSubmit={handleSubmit}>
    <Field
      name={FORM_FIELD_NAMES.code}
      fieldName={FORM_FIELD_NAMES.code}
      label="Enter 16 Digit Verification Code"
      validate={[required, maxLength(16)]}
      component={renderField}
    />
  </OneFieldOneButtonForm>
));

export { SendCodeForm, VerifyCodeForm };
