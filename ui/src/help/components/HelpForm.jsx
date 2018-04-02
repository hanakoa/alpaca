import React, { Component } from 'react';
import { Field, reduxForm } from 'redux-form';
import { Link } from 'react-router-dom';
import { FORM_NAME, FORM_FIELD_NAMES } from '../constants';
import {
  required,
  maxLength,
  renderTextField,
  renderTextArea,
} from '../../redux-form-helpers';
import Button from 'material-ui/Button';

const CardTitle = () => <h4 className="card-title d-flex my-4">Help</h4>;

class HelpForm extends Component {
  handleSubmit = event => {
    event.preventDefault();

    const { onSubmit } = this.props;
    onSubmit();
  };

  render() {
    const { submitting, handleSubmit } = this.props;

    return (
      <div className="mx-4 px-5">
        <CardTitle />

        <form onSubmit={handleSubmit}>
          <Field
            name={FORM_FIELD_NAMES.name}
            label="Name"
            component={renderTextField}
            validate={[required, maxLength(255)]}
            autoFocus={true}
          />
          <Field
            name={FORM_FIELD_NAMES.emailAddress}
            label="Email"
            component={renderTextField}
            validate={[required, maxLength(255)]}
          />
          <Field
            name={FORM_FIELD_NAMES.description}
            label="How can we help?"
            component={renderTextArea}
            validate={[required]}
            rows={5}
          />
        </form>

        <div>
          <Button
            variant="raised"
            color="primary"
            className="d-flex mx-auto my-4"
            disabled={submitting}
            onClick={event => this.handleSubmit(event)}>
            Submit
          </Button>
        </div>
        <div className="d-flex justify-content-center mb-4">
          <Button component={Link} to="/" variant="raised" color="secondary">
            Back
          </Button>
        </div>
      </div>
    );
  }
}

export default reduxForm({
  form: FORM_NAME,
})(HelpForm);
