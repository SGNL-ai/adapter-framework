# Copyright 2025 SGNL.ai, Inc.

GOLANG_IMAGE               := golang:1.26-bookworm
PROTOBUF_VERSION           := 23.3
PROTOC_GEN_GO_VERSION      := v1.36.6
PROTOC_GEN_GO_GRPC_VERSION := v1.6.1
GRPC_GATEWAY_MODULE        := github.com/grpc-ecosystem/grpc-gateway/v2
GRPC_GATEWAY_VERSION       := v2.29.0
MODULE                     := github.com/sgnl-ai/adapter-framework

.PHONY: proto
proto:
	docker run --rm \
	  -v "$(CURDIR):/workspace" \
	  -w /workspace \
	  $(GOLANG_IMAGE) \
	  bash -c '\
	    apt-get update -q && apt-get install -y -q unzip curl 2>&1 | tail -3 && \
	    ARCH=$$(uname -m) && \
	    if [ "$$ARCH" = "aarch64" ]; then PROTOC_ARCH=aarch_64; else PROTOC_ARCH=x86_64; fi && \
	    curl -fSsL "https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOBUF_VERSION)/protoc-$(PROTOBUF_VERSION)-linux-$${PROTOC_ARCH}.zip" > /tmp/protoc.zip && \
	    cd /usr/local && unzip -q /tmp/protoc.zip bin/protoc "include/*" && \
	    chmod +x /usr/local/bin/protoc && rm /tmp/protoc.zip && \
	    mkdir -p /tmp/googleapis/google/api && \
	    curl -fSsL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > /tmp/googleapis/google/api/annotations.proto && \
	    curl -fSsL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > /tmp/googleapis/google/api/http.proto && \
	    go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GEN_GO_VERSION) && \
	    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GEN_GO_GRPC_VERSION) && \
	    go install $(GRPC_GATEWAY_MODULE)/protoc-gen-grpc-gateway@$(GRPC_GATEWAY_VERSION) && \
	    cd /workspace && \
	    protoc \
	      --proto_path=. \
	      --proto_path=/tmp/googleapis \
	      --go_out=. --go_opt=module=$(MODULE) \
	      --go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
	      --grpc-gateway_out=. --grpc-gateway_opt=module=$(MODULE) \
	      proto/grpc_proxy/v1/http.proto \
	      proto/grpc_proxy/v1/sql.proto \
	      proto/grpc_proxy/v1/ldap.proto \
	      proto/grpc_proxy/v1/grpc_proxy.proto \
	  '
