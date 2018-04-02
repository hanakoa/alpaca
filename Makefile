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

.PHONY: docker
docker:
	@$(MAKE) docker-remove
	@$(MAKE) docker-rebuild
	@$(MAKE) docker-start
	@$(MAKE) docker-seed

.PHONY: docker-build
docker-build:
	docker image build -t hanakoa/alpaca-auth-api:latest -f auth/Dockerfile .
	docker image build -t hanakoa/alpaca-mfa-api:latest -f mfa/Dockerfile .
	docker image build -t hanakoa/alpaca-password-reset-api:latest -f password-reset/Dockerfile .
	docker-compose build

.PHONY: docker-rebuild
docker-rebuild:
	docker image build -t hanakoa/alpaca-auth-api:latest -f auth/Dockerfile . --no-cache
	docker image build -t hanakoa/alpaca-mfa-api:latest -f mfa/Dockerfile . --no-cache
	docker image build -t hanakoa/alpaca-password-reset-api:latest -f password-reset/Dockerfile . --no-cache
	docker-compose build

.PHONY: docker-remove
docker-remove:
	docker rm --force alpaca-auth-api || true
	docker rm --force alpaca-auth-db || true
	docker rm --force alpaca-password-reset-api || true
	docker rm --force alpaca-password-reset-db || true
	docker rm --force alpaca-rabbitmq || true

.PHONY: docker-start
docker-start:
	docker-compose up -d

.PHONY: docker-stop
docker-stop:
	docker-compose stop



.PHONY: docker-seed
docker-seed:
	./scripts/seed-data.sh "docker"

.PHONY: seed
seed:
	./scripts/seed-data.sh "local"

.PHONY: test-seed
test-seed:
	./scripts/seed-data.sh "test"



minikube-build:
	@eval $$(minikube docker-env) ;\
	docker image build -t hanakoa/alpaca-auth-api:latest -f auth/Dockerfile .
	docker image build -t hanakoa/alpaca-password-reset-api:latest -f password-reset/Dockerfile .
	docker image build -t hanakoa/alpaca-ui -f ui/Dockerfile .

.PHONY: protoc
protoc:
	protoc -I auth auth/pb/auth.proto --go_out=plugins=grpc:auth
	protoc -I mfa mfa/pb/mfa.proto --go_out=plugins=grpc:mfa

.PHONY: lint
lint:
	golint ./auth
	golint ./password-reset
	golint ./mfa

.PHONY: fmt
fmt:
	go fmt ./auth
	go fmt ./password-reset
	go fmt ./mfa

.PHONY: vet
vet:
	go tool vet auth
	go tool vet password-reset
	go tool vet mfa

.PHONY: test
test:
	@$(MAKE) test-seed
	go test -v ./auth



.PHONY: list-users
list-users:
	http localhost:8080/person

.PHONY: create-user
create-user:
	http POST localhost:8080/person username="kevin_chen" emailAddress="kevin.chen.bulk@gmail.com"

.PHONY: install-proto
install-proto:
	go get -u github.com/golang/protobuf/protoc-gen-go