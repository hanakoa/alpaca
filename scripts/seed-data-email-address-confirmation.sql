DROP TABLE IF EXISTS email_address_confirmation_code CASCADE;

-- Create syntax for TABLE 'email_address_confirmation_code'
CREATE TABLE email_address_confirmation_code (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed boolean NOT NULL DEFAULT FALSE,
  email_address_id bigint NOT NULL
) WITH (OIDS=FALSE);

CREATE INDEX email_address_confirmation_code_email_address_id_idx ON email_address_confirmation_code (email_address_id);