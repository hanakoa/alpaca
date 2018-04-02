import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { Field, reduxForm } from 'redux-form';
import { Link } from 'react-router-dom';
import Button from 'material-ui/Button';
import { loadingSelector } from './selectors';
import { maxLength, renderTextField, required } from '../../redux-form-helpers';
import { NAME, FIELD_NAME } from './constants';
import { Loading } from '../../app';

class ConfirmPinResetContainer extends Component {
  render() {
    const { loading } = this.props;

    return (
      <div>
        {loading && <Loading />}
        <div className="mx-4 px-5">
          <div className="card-title d-flex my-4 password-reset-form-header">
            Check your phone
          </div>
          <div>
            We've texted a code to the phone number ending in '87'. Once you
            receive the code, enter it below to reset your password.
          </div>
          <form onSubmit={() => console.log('submit')}>
            <Field
              name={FIELD_NAME}
              label="Enter Code"
              component={renderTextField}
              validate={[required, maxLength(16)]}
            />
          </form>
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
              onClick={() => console.log('submit')}>
              Reset
            </Button>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  loading: loadingSelector,
  // lastTwoDigits: lastTwoDigitsSelector,
});

const mapDispatchToProps = dispatch => ({
  actions: bindActionCreators({}, dispatch),
});

ConfirmPinResetContainer = reduxForm({
  form: NAME,
})(ConfirmPinResetContainer);

ConfirmPinResetContainer = connect(mapStateToProps, mapDispatchToProps)(
  ConfirmPinResetContainer
);

export default ConfirmPinResetContainer;
