SERVICES_DIR = ./services

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
	go build -o ./bin/alpaca-auth $(SERVICES_DIR)/auth

.PHONY: build-password-reset
build-password-reset:
	go build -o ./bin/alpaca-password-reset $(SERVICES_DIR)/password-reset

.PHONY: build-mfa
build-mfa:
	go build -o ./bin/alpaca-mfa $(SERVICES_DIR)/mfa

.PHONY: run-auth
run-auth:
	echo "Starting Auth µService" && \
  RABBITMQ_ENABLED=false \
  ORIGIN_ALLOWED=http://localhost:3000 \
  DB_PASSWORD=password \
  DB_HOST=localhost \
  ALPACA_SECRET=4FFFA6A10E744158464EB55133A475673264748804882A1B4F8106D545C584EF \
  ITERATION_COUNT=10000 \
  ./bin/alpaca-auth

.PHONY: run-password-reset
run-password-reset:
	echo "Starting Password Reset µService" && \
  RABBITMQ_ENABLED=false \
  ORIGIN_ALLOWED=http://localhost:3000 \
  DB_PASSWORD=password \
  DB_HOST=localhost \
  ./bin/alpaca-password-reset

.PHONY: run-mfa
run-mfa:
	TWILIO_ACCOUNT_SID=${TWILIO_ACCOUNT_SID} TWILIO_AUTH_TOKEN=${TWILIO_AUTH_TOKEN} TWILIO_PHONE_NUMBER=${TWILIO_PHONE_NUMBER} RABBITMQ_ENABLED=false ORIGIN_ALLOWED=http://localhost:3000 DB_PASSWORD=password DB_HOST=localhost ./bin/alpaca-mfa

.PHONY: test
test:
	cd ${GOPATH}/src/github.com/hanakoa/alpaca/$(SERVICES_DIR)/auth && vgo test -v .
	cd ${GOPATH}/src/github.com/hanakoa/alpaca/$(SERVICES_DIR)/password-reset && pwd

.PHONY: convey
convey:
#	go get github.com/smartystreets/goconvey
	${GOPATH}/bin/goconvey

.PHONY: list-users
list-users:
	http localhost:8080/account

.PHONY: create-user
create-user:
	http POST localhost:8080/account username="kevin_chen" email_address="kevin.chen.bulk.test@gmail.com"

include makefiles/*.mk
