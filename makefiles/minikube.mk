.PHONY: mk-start
mk-start:
	minikube start --memory 2048 --cpus 2 --vm-driver=hyperkit

.PHONY: mk-stop
mk-stop:
	minikube stop || true

.PHONY: mk-upgrade
mk-upgrade:
	@$(MAKE) mk-stop
	minikube delete || true
	rm -rf ~/.minikube
	brew cask reinstall minikube

.PHONY: mk-build
mk-build:
	@eval $$(minikube docker-env) ;\
	docker image build -t hanakoa/alpaca-auth-api:v0.0.1 -f auth/Dockerfile .
#	docker image build -t hanakoa/alpaca-password-reset-api:v0.0.1 -f password-reset/Dockerfile .
#	docker image build -t hanakoa/alpaca-ui -f ui/Dockerfile .

.PHONY: mk-rebuild
mk-rebuild:
	@eval $$(minikube docker-env) ;\
	docker image build -t hanakoa/alpaca-auth-api:v0.0.1 -f auth/Dockerfile . --no-cache
#	docker image build -t hanakoa/alpaca-password-reset-api:v0.0.1 -f password-reset/Dockerfile . --no-cache
#	docker image build -t hanakoa/alpaca-ui -f ui/Dockerfile . --no-cache