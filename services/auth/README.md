# Authentication
This is the primary µService for Alpaca.
It manages user credentials (username, email, phone number)

## Using the API
### Creating a user
```bash
http POST localhost:8080/account username="kevin_chen" email_address="kevin.chen.bulk@gmail.com"
```

### Set password
```bash
http PUT localhost:8080/account/970480225798328320/password password="MyPassword123Is-Super-Secure"
```

### Logging in
```bash
http POST localhost:8080/token login="kevin_chen" password="MyPassword123Is-Super-Secure"
```