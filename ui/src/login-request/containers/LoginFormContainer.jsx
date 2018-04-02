import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { Link } from 'react-router-dom';
import classNames from 'classnames';
import { loadingSelector } from '../selectors';
import { loginRequest } from '../duck';
import LoginForm from '../components/LoginForm';
import Button from 'material-ui/Button';

const LoginLinks = () => {
  const padding = 'py-1';
  const linkClass = classNames('d-flex', 'justify-content-center', padding);
  return [
    <div key="login-link-password-reset" className={linkClass}>
      <Link to={`/account/begin_password_reset`}>
        <button type="button" className="btn btn-sm btn-link">
          Forgot Password?
        </button>
      </Link>
    </div>,
    <div key="login-link-account-activation" className={linkClass}>
      <Link to={`/account-activation`}>
        <button type="button" className="btn btn-sm btn-link">
          Account Activation Request
        </button>
      </Link>
    </div>,
    <div key="login-link-help" className={classNames(linkClass, 'mb-2')}>
      <Link to={`/help`}>
        <button type="button" className="btn btn-sm btn-link">
          Need Help?
        </button>
      </Link>
    </div>,
  ];
};

class LoginFormContainer extends Component {
  handleSubmit(event) {
    event.preventDefault();

    const { redirectTo } = this.props;
    const { onLoginSubmit } = this.props.actions;
    onLoginSubmit(redirectTo);
  }

  render() {
    const { submitting } = this.props;

    return (
      <div className="px-5">
        <LoginForm {...this.props} />

        <Button
          variant="raised"
          color="primary"
          className="d-flex mx-auto my-4"
          disabled={submitting}
          onClick={event => this.handleSubmit(event)}>
          Sign in
        </Button>

        <LoginLinks />
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
      onLoginSubmit: loginRequest,
    },
    dispatch
  ),
});

export default connect(mapStateToProps, mapDispatchToProps)(LoginFormContainer);
