DOCKER_ORG := hanakoa
DOCKER_IMAGE_VERSION = v0.0.1
SERVICES = auth mfa password-reset

.PHONY: docker
docker:
	@$(MAKE) docker-stop
	@$(MAKE) docker-remove
	@$(MAKE) docker-rebuild
	@$(MAKE) docker-start
	@$(MAKE) docker-seed

.PHONY: docker-build
docker-build:
	for svc in $(SERVICES); do \
		echo Building image for $$svc; \
		docker image build -t $(DOCKER_ORG)/alpaca-$$svc-api:$(DOCKER_IMAGE_VERSION) -t $(DOCKER_ORG)/alpaca-$$svc-api:latest -f $$svc/Dockerfile . ; \
	done
	docker-compose build

.PHONY: docker-rebuild
docker-rebuild:
	for svc in $(SERVICES); do \
		echo Rebuilding image for $$svc; \
		docker image build -t $(DOCKER_ORG)/alpaca-$$svc-api:$(DOCKER_IMAGE_VERSION) -t $(DOCKER_ORG)/alpaca-$$svc-api:latest -f $$svc/Dockerfile . --no-cache ; \
	done
	docker-compose build

.PHONY: docker-remove
docker-remove:
	for svc in $(SERVICES); do \
		echo Removing image for $$svc; \
		docker rm --force alpaca-$$svc-api || true
	done
	@docker rm --force alpaca-rabbitmq || true

.PHONY: docker-start
docker-start:
	docker-compose up -d

.PHONY: docker-stop
docker-stop:
	@docker-compose stop || true
