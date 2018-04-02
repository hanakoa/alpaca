import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { sendCodeRequest, verifyCodeRequest } from '../duck';
import { loadingSelector } from '../selectors';
import { SendCodeForm, VerifyCodeForm } from '../components';
import { Loading } from '../../app';

class TwoFactorAuthContainer extends Component {
  render() {
    const { loading } = this.props;
    const { sendCodeRequest, verifyCodeRequest } = this.props.actions;
    return (
      <div>
        {loading && <Loading />}
        <SendCodeForm handleSubmit={sendCodeRequest} />
        <VerifyCodeForm handleSubmit={verifyCodeRequest} />
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
      sendCodeRequest,
      verifyCodeRequest,
    },
    dispatch
  ),
});

export default connect(mapStateToProps, mapDispatchToProps)(
  TwoFactorAuthContainer
);
