DROP TABLE IF EXISTS authentication_code CASCADE;

-- Create syntax for TABLE 'authentication_code'
CREATE TABLE authentication_code (
  id uuid PRIMARY KEY,
  code varchar(6) NOT NULL,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiration_timestamp timestamp NOT NULL,
  usable boolean NOT NULL,
  used boolean NOT NULL,
  person_id bigint NOT NULL
) WITH (OIDS=FALSE);

CREATE INDEX authentication_code_person_id_code_expiration_idx ON authentication_code (person_id, code, expiration_timestamp);