# Password Resets

This µService's sole purpose is to perform password resets.

It consists of 3 endpoints:

* `POST /password-reset` creates a reset code and fires an email to the intended recipient.
  * The client won't know if an email does not exist. We don't want to leak that info.
* `GET /password-reset/{code}` confirms that a usuable and unexpired code exists for the given UUID.
  * We return a 404 if nothing is found, otherwise a 200.
* `PUT /password-reset` performs the reset and renders all prior codes as unusable.

## grpc client

This library has a grpc dependency on the main auth µService for two reasons:

* to get the `accountId` for a given email address
  * this is because this database persists accountId, not email addresses, which can change with time
* to perform the password reset

## Using the API

### Sending the reset code

```bash
http POST localhost:8081/password-reset account="kevin.chen.bulk@gmail.com"
```

The `account` field can be a username, email address, or phone number.

If nothing fatal happens, this endpoint returns a 200 OK, even if an account does not exist, so we do not leak that info.
If a account has additional means of login, besides just a primary email address, then we return a list of options:

```json
{
  "account_id": 1,
  "code": "B390A341-A88A-4F5A-5F8E-12656368D328",
  "options": [
    {
      "type": "email_address",
      "value": "ke**************@g****.com"
    },
    {
      "type": "phone_number",
      "value": "87"
    }
  ]
}
```

### Confirming the code exists

```bash
http GET localhost:8081/password-reset/B390A341-A88A-4F5A-5F8E-12656368D328
```

If no valid (usable, unused, unexpired) code exists, we return a 404 Not Found:

```json
{
  "error": "No password reset code for: 29538B59-204F-450E-9253-8A2762D07C1f"
}
```

Otherwise we return a 200 OK.

### Triggering password reset

```bash
http PUT localhost:8081/password-reset email_address="kevin.chen.bulk@gmail.com" code="38E45CD6-D09C-4F41-AB24-FCB4795B7335" password="Standing-On-Shoulders-Of-Giants93"
```
