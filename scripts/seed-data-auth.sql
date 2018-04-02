DROP TABLE IF EXISTS email_address CASCADE;
DROP TABLE IF EXISTS phone_number CASCADE;
DROP TABLE IF EXISTS login_attempt CASCADE;
DROP TABLE IF EXISTS password CASCADE;
DROP TABLE IF EXISTS person CASCADE;

-- Create syntax for TABLE 'email_address'
CREATE TABLE email_address (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  confirmed boolean NOT NULL DEFAULT FALSE,
  is_primary boolean NOT NULL,
  email_address varchar(255) NOT NULL,
  person_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'phone_number'
CREATE TABLE phone_number (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  phone_number varchar(50) NOT NULL,
  person_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'login_attempt'
CREATE TABLE login_attempt (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  success boolean NOT NULL,
  person_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'password'
CREATE TABLE password (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  iteration_count int DEFAULT NULL NOT NULL,
  salt bytea NOT NULL,
  password_hash bytea DEFAULT NULL,
  person_id bigint NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'person'
CREATE TABLE person (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  disabled boolean NOT NULL DEFAULT FALSE,
  multi_factor_required boolean NOT NULL DEFAULT FALSE,
  username varchar(25) NOT NULL,
  current_password_id bigint DEFAULT NULL,
  primary_email_address_id bigint DEFAULT NULL,
  UNIQUE (username)
) WITH (OIDS=FALSE);

-- Create foreign key constraints
ALTER TABLE email_address
  ADD CONSTRAINT email_address_person_id_fkey
FOREIGN KEY (person_id)
REFERENCES person(id);

ALTER TABLE phone_number
  ADD CONSTRAINT phone_number_person_id_fkey
FOREIGN KEY (person_id)
REFERENCES person(id);

ALTER TABLE login_attempt
  ADD CONSTRAINT login_attempt_person_id_fkey
FOREIGN KEY (person_id)
REFERENCES person(id);

ALTER TABLE password
  ADD CONSTRAINT password_person_id_fkey
FOREIGN KEY (person_id)
REFERENCES person(id);

ALTER TABLE person
  ADD CONSTRAINT person_current_password_id_fkey
FOREIGN KEY (current_password_id)
REFERENCES password(id);

ALTER TABLE person
  ADD CONSTRAINT person_primary_email_address_id_fkey
FOREIGN KEY (primary_email_address_id)
REFERENCES email_address(id);

-- Create indexes
CREATE INDEX email_address_email_address_idx ON email_address (email_address);
CREATE INDEX email_address_person_id_idx ON email_address (person_id);
CREATE INDEX login_attempt_person_id_idx ON login_attempt (person_id);
CREATE INDEX login_attempt_created_timestamp_idx ON login_attempt (created_timestamp);
CREATE INDEX password_person_id_idx ON password (person_id);
CREATE INDEX person_current_password_id_idx ON person (current_password_id);
CREATE INDEX person_primary_email_address_id_idx ON person (primary_email_address_id);