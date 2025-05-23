// Copyright 2025 SGNL.ai, Inc.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v3.21.12
// source: proto/grpc_proxy/v1/sql.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// SQLQueryRequest is a wrapper around a marshalled SQL Adapter request to an
// on-premises connector.
type SQLQueryRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Request       string                 `protobuf:"bytes,1,opt,name=request,proto3" json:"request,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SQLQueryRequest) Reset() {
	*x = SQLQueryRequest{}
	mi := &file_proto_grpc_proxy_v1_sql_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SQLQueryRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SQLQueryRequest) ProtoMessage() {}

func (x *SQLQueryRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_grpc_proxy_v1_sql_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SQLQueryRequest.ProtoReflect.Descriptor instead.
func (*SQLQueryRequest) Descriptor() ([]byte, []int) {
	return file_proto_grpc_proxy_v1_sql_proto_rawDescGZIP(), []int{0}
}

func (x *SQLQueryRequest) GetRequest() string {
	if x != nil {
		return x.Request
	}
	return ""
}

// SQLQueryResponse is a wrapper around a marshalled SQL processed response or
// any error (marshalled framework.Error) while processing the request from an on-premises connector.
type SQLQueryResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Response      string                 `protobuf:"bytes,1,opt,name=response,proto3" json:"response,omitempty"`
	Error         string                 `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *SQLQueryResponse) Reset() {
	*x = SQLQueryResponse{}
	mi := &file_proto_grpc_proxy_v1_sql_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SQLQueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SQLQueryResponse) ProtoMessage() {}

func (x *SQLQueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_grpc_proxy_v1_sql_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SQLQueryResponse.ProtoReflect.Descriptor instead.
func (*SQLQueryResponse) Descriptor() ([]byte, []int) {
	return file_proto_grpc_proxy_v1_sql_proto_rawDescGZIP(), []int{1}
}

func (x *SQLQueryResponse) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

func (x *SQLQueryResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_proto_grpc_proxy_v1_sql_proto protoreflect.FileDescriptor

const file_proto_grpc_proxy_v1_sql_proto_rawDesc = "" +
	"\n" +
	"\x1dproto/grpc_proxy/v1/sql.proto\x12\x12sgnl.grpc_proxy.v1\"+\n" +
	"\x0fSQLQueryRequest\x12\x18\n" +
	"\arequest\x18\x01 \x01(\tR\arequest\"D\n" +
	"\x10SQLQueryResponse\x12\x1a\n" +
	"\bresponse\x18\x01 \x01(\tR\bresponse\x12\x14\n" +
	"\x05error\x18\x02 \x01(\tR\x05errorB8Z6github.com/sgnl-ai/adapter-framework/pkg/grpc_proxy/v1b\x06proto3"

var (
	file_proto_grpc_proxy_v1_sql_proto_rawDescOnce sync.Once
	file_proto_grpc_proxy_v1_sql_proto_rawDescData []byte
)

func file_proto_grpc_proxy_v1_sql_proto_rawDescGZIP() []byte {
	file_proto_grpc_proxy_v1_sql_proto_rawDescOnce.Do(func() {
		file_proto_grpc_proxy_v1_sql_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_grpc_proxy_v1_sql_proto_rawDesc), len(file_proto_grpc_proxy_v1_sql_proto_rawDesc)))
	})
	return file_proto_grpc_proxy_v1_sql_proto_rawDescData
}

var file_proto_grpc_proxy_v1_sql_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_grpc_proxy_v1_sql_proto_goTypes = []any{
	(*SQLQueryRequest)(nil),  // 0: sgnl.grpc_proxy.v1.SQLQueryRequest
	(*SQLQueryResponse)(nil), // 1: sgnl.grpc_proxy.v1.SQLQueryResponse
}
var file_proto_grpc_proxy_v1_sql_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_grpc_proxy_v1_sql_proto_init() }
func file_proto_grpc_proxy_v1_sql_proto_init() {
	if File_proto_grpc_proxy_v1_sql_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_grpc_proxy_v1_sql_proto_rawDesc), len(file_proto_grpc_proxy_v1_sql_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_grpc_proxy_v1_sql_proto_goTypes,
		DependencyIndexes: file_proto_grpc_proxy_v1_sql_proto_depIdxs,
		MessageInfos:      file_proto_grpc_proxy_v1_sql_proto_msgTypes,
	}.Build()
	File_proto_grpc_proxy_v1_sql_proto = out.File
	file_proto_grpc_proxy_v1_sql_proto_goTypes = nil
	file_proto_grpc_proxy_v1_sql_proto_depIdxs = nil
}
