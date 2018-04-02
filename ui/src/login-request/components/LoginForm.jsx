import React, { Component } from 'react';
import { FORM_FIELD_NAMES, FORM_NAME } from '../constants';
import {
  maxLength,
  renderPassword,
  renderTextField,
  required,
} from '../../redux-form-helpers';
import { Field, reduxForm } from 'redux-form';

class LoginForm extends Component {
  render() {
    const { handleSubmit } = this.props;
    return (
      <form onSubmit={handleSubmit}>
        <div>
          <Field
            name={FORM_FIELD_NAMES.login}
            label="Email or username"
            component={renderTextField}
            validate={[required, maxLength(255)]}
            autoFocus={true}
          />
        </div>
        <div>
          <Field
            name={FORM_FIELD_NAMES.password}
            label="Password"
            component={renderPassword}
            validate={[required]}
          />
        </div>
      </form>
    );
  }
}

export default reduxForm({
  form: FORM_NAME,
})(LoginForm);
