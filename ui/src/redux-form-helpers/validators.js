const required = (value, allValues, { pristine }, name) => {
  return value || pristine ? undefined : 'Required';
};

const maxLength = max => value =>
  value && value.length > max ? `Must be ${max} characters or less` : undefined;

export { required, maxLength };
