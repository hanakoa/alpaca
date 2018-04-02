import React from 'react';
import { required, maxLength } from './validators';
import TextField from 'material-ui/TextField';
import Checkbox from 'material-ui/Checkbox';
import { FormControlLabel } from 'material-ui/Form';
import PasswordInput from './password-input';

const renderTextArea = props => (
  <TextField
    label={props.label}
    placeholder={props.label}
    multiline
    margin="normal"
    fullWidth={props.fullWidth || true}
    {...props.input}
  />
);

const renderTextField = props => (
  <TextField
    hintText={props.label}
    floatingLabelText={props.label}
    errorText={props.touched && props.error}
    fullWidth={props.fullWidth || true}
    {...props.input}
    {...props}
  />
);

const renderPassword = PasswordInput;

const renderCheckbox = props => (
  <FormControlLabel
    control={
      <Checkbox
        checked={!!props.value}
        onChange={props.onChange}
        value={props.name}
      />
    }
    label={props.label}
  />
);

export {
  required,
  maxLength,
  renderCheckbox,
  renderPassword,
  renderTextField,
  renderTextArea,
};
