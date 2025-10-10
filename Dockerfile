# Copyright 2023 SGNL.ai, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG GOLANG_IMAGE=golang:1.24-bookworm

FROM ${GOLANG_IMAGE} as build

RUN apt-get update && apt-get install -y \
    unzip=6.0*

ARG PROTOBUF_VERSION=23.3
RUN curl -fSsL "https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOBUF_VERSION}/protoc-${PROTOBUF_VERSION}-linux-x86_64.zip" > /tmp/protoc.zip \
    && (cd /usr/local && unzip /tmp/protoc.zip 'bin/protoc' 'include/*')  \
    && chmod +x /usr/local/bin/protoc \
    && rm -f /tmp/protoc.zip

ARG PROTOC_GEN_GO_VERSION=1.28.1
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v${PROTOC_GEN_GO_VERSION}

ARG PROTOC_GEN_GO_RPC_VERSION=1.3.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v${PROTOC_GEN_GO_RPC_VERSION}

WORKDIR /build
COPY . ./

RUN protoc \
    --proto_path=/usr/local/include \
    --proto_path=api \
    --go_out=paths=source_relative:api \
    --go-grpc_out=paths=source_relative:api \
    --go_opt=paths=source_relative \
    adapter/v1/adapter.proto