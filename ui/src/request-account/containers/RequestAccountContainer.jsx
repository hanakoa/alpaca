import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import React, { Component } from 'react';
import { requestAccountRequest } from '../duck';
import RequestAccountForm from '../components/RequestAccountForm';
import { loadingSelector, successSelector } from '../selectors';

class RequestAccountContainer extends Component {
  render() {
    const { loading, success } = this.props;
    const { onSubmit } = this.props.actions;
    return (
      <div>
        <RequestAccountForm
          loading={loading}
          success={success}
          onSubmit={onSubmit}
        />
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  loading: loadingSelector,
  success: successSelector,
});

const mapDispatchToProps = dispatch => ({
  actions: bindActionCreators(
    {
      onSubmit: requestAccountRequest,
    },
    dispatch
  ),
});

export default connect(mapStateToProps, mapDispatchToProps)(
  RequestAccountContainer
);
