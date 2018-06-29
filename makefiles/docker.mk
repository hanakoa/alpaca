.PHONY: docker
docker:
	@$(MAKE) docker-remove
	@$(MAKE) docker-rebuild
	@$(MAKE) docker-start
	@$(MAKE) docker-seed

.PHONY: docker-build
docker-build:
	docker image build -t hanakoa/alpaca-auth-api:v0.0.1 -f auth/Dockerfile .
	docker image build -t hanakoa/alpaca-mfa-api:v0.0.1 -f mfa/Dockerfile .
	docker image build -t hanakoa/alpaca-password-reset-api:v0.0.1 -f password-reset/Dockerfile .
	docker-compose build

.PHONY: docker-rebuild
docker-rebuild:
	docker image build -t hanakoa/alpaca-auth-api:v0.0.1 -f auth/Dockerfile . --no-cache
	docker image build -t hanakoa/alpaca-mfa-api:v0.0.1 -f mfa/Dockerfile . --no-cache
	docker image build -t hanakoa/alpaca-password-reset-api:v0.0.1 -f password-reset/Dockerfile . --no-cache
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
