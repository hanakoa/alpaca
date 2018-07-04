DROP TABLE IF EXISTS password_reset_code CASCADE;

-- Create syntax for TABLE 'password_reset_code'
CREATE TABLE password_reset_code (
  code uuid PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiration_timestamp timestamp NOT NULL,
  usable boolean NOT NULL,
  used boolean NOT NULL,
  account_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'email_address'
CREATE TABLE email_address (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed boolean NOT NULL DEFAULT FALSE,
  is_primary boolean NOT NULL,
  email_address varchar(255) NOT NULL,
  account_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'phone_number'
CREATE TABLE phone_number (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed boolean NOT NULL DEFAULT FALSE,
  phone_number varchar(255) NOT NULL,
  account_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- TODO we should store current pass hash to enforce that user actually changes their pass
-- Create syntax for TABLE 'account'
CREATE TABLE account (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  disabled boolean NOT NULL DEFAULT FALSE,
  username varchar(25) NOT NULL,
  current_password_id bigint DEFAULT NULL,
  primary_email_address_id bigint DEFAULT NULL,
  UNIQUE (username)
) WITH (OIDS=FALSE);

CREATE INDEX password_reset_code_account_id_idx ON password_reset_code (account_id);
CREATE INDEX email_address_account_id_idx ON email_address (account_id);
CREATE INDEX phone_number_account_id_idx ON phone_number (account_id);