.PHONY: protoc
protoc:
	protoc -I auth auth/pb/auth.proto --go_out=plugins=grpc:auth
	protoc -I mfa mfa/pb/mfa.proto --go_out=plugins=grpc:mfa

.PHONY: install-proto
install-proto:
	go get -u github.com/golang/protobuf/protoc-gen-go
