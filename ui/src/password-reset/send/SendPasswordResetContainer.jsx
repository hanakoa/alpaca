import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { Field, reduxForm } from 'redux-form';
import { Link } from 'react-router-dom';
import Button from 'material-ui/Button';
import { loadingSelector, optionsSelector } from './selectors';
import { NAME, FIELD_NAMES } from './constants';
import { Loading } from '../../app';

class SendPasswordResetContainer extends Component {
  renderOption = option => {
    return (
      <label>
        <Field
          name={FIELD_NAMES.option}
          component="input"
          type="radio"
          value={`${option.type};${option.id}`}
        />
        {option.type === 'email_address' && (
          <span>
            Email a link to <strong>{option.value}</strong>
          </span>
        )}
        {option.type === 'phone_number' && (
          <span>
            Text a code to my phone ending in <strong>{option.value}</strong>
          </span>
        )}
      </label>
    );
  };

  submit = () => {
    console.log('submit');
  };

  render() {
    const { loading, options } = this.props;

    return (
      <div>
        {loading && <Loading />}
        <div className="mx-4 px-5">
          <div className="card-title d-flex my-4 password-reset-form-header">
            How do you want to reset your password?
          </div>
          <form onSubmit={this.submit}>{options.map(this.renderOption)}</form>
          <div className="d-flex justify-content-center">
            <Button
              component={Link}
              to="/"
              variant="raised"
              color="secondary"
              className="d-flex mx-auto my-4">
              Back
            </Button>
            <Button
              variant="raised"
              color="primary"
              className="d-flex mx-auto my-4"
              onClick={this.submit}>
              Continue
            </Button>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  loading: loadingSelector,
  options: optionsSelector,
});

const mapDispatchToProps = dispatch => ({
  actions: bindActionCreators({}, dispatch),
});

SendPasswordResetContainer = reduxForm({
  form: NAME,
})(SendPasswordResetContainer);

SendPasswordResetContainer = connect(mapStateToProps, mapDispatchToProps)(
  SendPasswordResetContainer
);

export default SendPasswordResetContainer;
