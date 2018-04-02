import React, { Component } from 'react';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import HelpForm from '../components/HelpForm';
import { helpSubmissionRequest } from '../duck';
import { loadingSelector } from '../selectors';
import { Loading } from '../../app';

class HelpContainer extends Component {
  render() {
    const { loading } = this.props;
    const { onSubmit } = this.props.actions;
    return (
      <div>
        {loading && <Loading />}
        <HelpForm onSubmit={onSubmit} />
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  loading: loadingSelector,
});

const mapDispatchToProps = dispatch => {
  return {
    actions: bindActionCreators(
      {
        onSubmit: helpSubmissionRequest,
      },
      dispatch
    ),
  };
};

export default connect(mapStateToProps, mapDispatchToProps)(HelpContainer);
