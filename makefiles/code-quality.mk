.PHONY: lint
lint:
	golint $(SERVICES_DIR)/auth
	golint $(SERVICES_DIR)/password-reset
	golint $(SERVICES_DIR)/mfa

.PHONY: fmt
fmt:
	go fmt $(SERVICES_DIR)/auth
	go fmt $(SERVICES_DIR)/password-reset
	go fmt $(SERVICES_DIR)/mfa

.PHONY: vet
vet:
	go tool vet auth
	go tool vet password-reset
	go tool vet mfa
