- postToken
  - accepts either an email or username
- getCurrentToken
  - insecure. we no longer offer this.
  - instead, an endpoint will return user info, not jwt

- putPassword

- getPeople
- getPerson

- getPersonByPrimaryEmailAddress
  - only support get person by id
- getPersonByEmailAddress
  - only support get person by id
- getPersonByUsername
  - only support get person by id

- createPerson
- fullyUpdatePerson
- partiallyUpdatePerson
- deletePerson
- getPersonEmailAddresses (*)
- createPersonEmailAddress (*)
- getPersonRoleMemberships (*)
- createPersonRoleMembership (*)
- getUpdatedPeople
- getRecentlyExpiredPasswords (*)

- getEmailAddresses
- getEmailAddress
- postEmailAddress
- putEmailAddress
- deleteEmailAddress

- getRoles
- getRole
- getRoleByName
- createRole
- updateRole
- deleteRole
- getPersonsForRole
- getUpdatedRoles

- getMemberships
- getMembership
- postMembership
- fullyUpdateRoleMembership
- partiallyUpdateRoleMembership
- deleteMembership
- getUpdatedMemberships
- getRecentlyExpiredMemberships
