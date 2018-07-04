```bash
# create account
http POST localhost:8080/account username="kevin_chen" email_address="kevin.chen.bulk.test@gmail.com"

# get accounts
http GET localhost:8080/account

# set password
http PUT localhost:8080/account/983088048847720448/password password="potato-tomato-cherry-gun"

# delete account
http DELETE localhost:8080/account/983088048847720448

# login
http POST localhost:8080/token login="kevin_chen" password="potato-tomato-cherry-gun"
```
