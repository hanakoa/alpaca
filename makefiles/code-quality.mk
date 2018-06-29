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
