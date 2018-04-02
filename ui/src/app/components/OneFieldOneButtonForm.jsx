import React from 'react';
import Button from 'material-ui/Button';

const OneFieldOneButtonForm = ({
  header,
  buttonText,
  handleSubmit,
  children,
}) => (
  <div className="mx-4 px-5">
    <div className="card-title d-flex my-4 password-reset-form-header">
      {header}
    </div>
    <form onSubmit={handleSubmit}>{children}</form>
    <div className="d-flex justify-content-center">
      <Button
        variant="raised"
        color="primary"
        className="d-flex mx-auto my-4"
        onClick={handleSubmit}>
        {buttonText}
      </Button>
    </div>
  </div>
);

export default OneFieldOneButtonForm;
