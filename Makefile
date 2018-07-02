.PHONY: all
all:
	@$(MAKE) build

.PHONY: auth
auth:
	@$(MAKE) build-auth
	@$(MAKE) run-auth

.PHONY: password-reset
password-reset:
	@$(MAKE) build-password-reset
	@$(MAKE) run-password-reset

.PHONY: mfa
mfa:
	@$(MAKE) build-mfa
	@$(MAKE) run-mfa

.PHONY: build
build:
	@$(MAKE) build-auth
	@$(MAKE) build-password-reset
	@$(MAKE) build-mfa

.PHONY: build-auth
build-auth:
	go build -o ./bin/alpaca-auth ./auth

.PHONY: build-password-reset
build-password-reset:
	go build -o ./bin/alpaca-password-reset ./password-reset

.PHONY: build-mfa
build-mfa:
	go build -o ./bin/alpaca-mfa ./mfa

.PHONY: run-auth
run-auth:
	RABBITMQ_ENABLED=false ORIGIN_ALLOWED=http://localhost:3000 DB_PASSWORD=password DB_HOST=localhost ALPACA_SECRET=4FFFA6A10E744158464EB55133A475673264748804882A1B4F8106D545C584EF ./bin/alpaca-auth

.PHONY: run-password-reset
run-password-reset:
	RABBITMQ_ENABLED=false ORIGIN_ALLOWED=http://localhost:3000 DB_PASSWORD=password DB_HOST=localhost ./bin/alpaca-password-reset

.PHONY: run-mfa
run-mfa:
	TWILIO_ACCOUNT_SID=${TWILIO_ACCOUNT_SID} TWILIO_AUTH_TOKEN=${TWILIO_AUTH_TOKEN} TWILIO_PHONE_NUMBER=${TWILIO_PHONE_NUMBER} RABBITMQ_ENABLED=false ORIGIN_ALLOWED=http://localhost:3000 DB_PASSWORD=password DB_HOST=localhost ./bin/alpaca-mfa

.PHONY: test
test:
	go test -v ./auth

.PHONY: circle-test
circle-test:
	go test -cover -coverprofile=/home/ubuntu/coverage.out -v ./auth

.PHONY: convey
convey:
#	go get github.com/smartystreets/goconvey
	${GOPATH}/bin/goconvey

.PHONY: list-users
list-users:
	http localhost:8080/person

.PHONY: create-user
create-user:
	http POST localhost:8080/person username="kevin_chen" emailAddress="kevin.chen.bulk@gmail.com"

include makefiles/*.mk
