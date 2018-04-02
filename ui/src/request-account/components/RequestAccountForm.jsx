import React, { Component } from 'react';
import { Field, reduxForm } from 'redux-form';
import { FORM_FIELD_NAMES, FORM_NAME } from '../constants';
import {
  required,
  maxLength,
  renderTextArea,
  renderTextField,
} from '../../redux-form-helpers';
import { LoadingButton } from '../../app';

class RequestAccountForm extends Component {
  handleSubmit = event => {
    event.preventDefault();

    const { onSubmit } = this.props;
    onSubmit();
  };

  render() {
    const { handleSubmit, loading, success } = this.props;

    return (
      <div className="mx-4">
        <form onSubmit={handleSubmit}>
          <Field
            name={FORM_FIELD_NAMES.firstName}
            label="First Name"
            component={renderTextField}
            validate={[required, maxLength(255)]}
            autoFocus={true}
          />
          <Field
            name={FORM_FIELD_NAMES.lastName}
            label="Last Name"
            component={renderTextField}
            validate={[required, maxLength(255)]}
          />
          <Field
            name={FORM_FIELD_NAMES.emailAddress}
            label="Email"
            component={renderTextField}
            validate={[required, maxLength(255)]}
          />
          <Field
            name={FORM_FIELD_NAMES.citizenship}
            label="Citizenship"
            component={renderTextField}
          />
          <Field
            name={FORM_FIELD_NAMES.reason}
            label="Reason to Join"
            validate={[required]}
            component={renderTextArea}
            rows={5}
          />
          <div>
            <LoadingButton
              buttonText="Send Request"
              loading={loading}
              success={success}
            />
          </div>
        </form>
      </div>
    );
  }
}

export default reduxForm({
  form: FORM_NAME,
})(RequestAccountForm);
