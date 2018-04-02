import React from 'react';
import Snackbar from 'material-ui/Snackbar';
import { connect } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import { bindActionCreators } from 'redux';
import { clearMessage } from './duck';
import { messageSelector, openSelector } from './selectors';

class PositionedSnackbar extends React.Component {
  render() {
    const { open, message } = this.props;
    const { onClose } = this.props.actions;
    return (
      <div>
        <Snackbar
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
          open={open}
          autoHideDuration={3000}
          onClose={() => onClose()}
          SnackbarContentProps={{
            'aria-describedby': 'message-id',
          }}
          message={<span id="message-id">{message}</span>}
        />
      </div>
    );
  }
}

const mapStateToProps = createStructuredSelector({
  message: messageSelector,
  open: openSelector,
});

const mapDispatchToProps = dispatch => ({
  actions: bindActionCreators(
    {
      onClose: clearMessage,
    },
    dispatch
  ),
});

export default connect(mapStateToProps, mapDispatchToProps)(PositionedSnackbar);
