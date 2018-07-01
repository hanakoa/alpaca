```bash
# create person
http POST localhost:8080/person username="kevin_chen" email_address="kevin.chen.bulk.test@gmail.com"

# get people
http GET localhost:8080/person

# set password
http PUT localhost:8080/person/983088048847720448/password password="potato-tomato-cherry-gun"

# delete person
http DELETE localhost:8080/person/983088048847720448

# login
http POST localhost:8080/token login="kevin_chen" password="potato-tomato-cherry-gun"
```
