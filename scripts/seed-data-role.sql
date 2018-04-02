DROP TABLE IF EXISTS role CASCADE;
DROP TABLE IF EXISTS role_membership CASCADE;

-- Create syntax for TABLE 'role'
CREATE TABLE role (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name varchar(50) NOT NULL
) WITH (OIDS=FALSE);

-- Create syntax for TABLE 'role_membership'
CREATE TABLE role_membership (
  id bigint PRIMARY KEY,
  created_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_timestamp timestamp NULL DEFAULT NULL,
  last_modified_timestamp timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  person_id bigint NOT NULL,
  role_id bigint NOT NULL,
  expiration_timestamp timestamp NULL DEFAULT NULL
) WITH (OIDS=FALSE);

ALTER TABLE role_membership
  ADD CONSTRAINT role_membership_role_id_fkey
FOREIGN KEY (role_id)
REFERENCES role(id);

CREATE INDEX role_name_idx ON role (name);
CREATE INDEX role_membership_person_id_idx ON role_membership (person_id);
CREATE INDEX role_membership_role_id_idx ON role_membership (role_id);