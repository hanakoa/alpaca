.PHONY: protoc
protoc:
	protoc -I $(SERVICES_DIR)/auth/pb $(SERVICES_DIR)/auth/pb/*.proto --go_out=plugins=grpc:./services/auth/pb
	protoc -I $(SERVICES_DIR)/mfa/pb $(SERVICES_DIR)/mfa/pb/*.proto --go_out=plugins=grpc:./services/mfa/pb

.PHONY: install-proto
install-proto:
	brew install protobuf
	go get -u github.com/golang/protobuf/protoc-gen-go
