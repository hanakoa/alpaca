import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { Field, reduxForm } from 'redux-form';
import { Link } from 'react-router-dom';
import Button from 'material-ui/Button';
import { loadingSelector } from './selectors';
import { renderTextField, required } from '../../redux-form-helpers';
import { NAME, FIELD_NAME } from './constants';
import { Loading } from '../../app';
import { beginPasswordReset } from './duck';

class BeginPasswordResetContainer extends Component {
  beginPasswordReset = () => {
    this.props.actions.beginPasswordReset();
  };

  render() {
    const { loading } = this.props;
    return (
      <div>
        {loading && <Loading />}
        <div className="mx-4 px-5">
          <div className="card-title d-flex my-4 password-reset-form-header">
            Find your account
          </div>
          <form onSubmit={this.beginPasswordReset}>
            <Field
              name={FIELD_NAME}
              label="Email, phone number, or username"
              component={renderTextField}
              validate={[required]}
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
              onClick={this.beginPasswordReset}>
              Search
            </Button>
          </div>
        </div>
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  loading: loadingSelector,
});

const mapDispatchToProps = dispatch => ({
  actions: bindActionCreators(
    {
      beginPasswordReset,
    },
    dispatch
  ),
});

BeginPasswordResetContainer = reduxForm({
  form: NAME,
})(BeginPasswordResetContainer);

BeginPasswordResetContainer = connect(mapStateToProps, mapDispatchToProps)(
  BeginPasswordResetContainer
);

export default BeginPasswordResetContainer;
