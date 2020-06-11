// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Protocol messages used to record the state of the model type store for USS.
// At the time of writing, the model type store uses leveldb, a schemaless
// key-value store. This means that the database's schema is mostly implicit.
// This descriptor isn't intended to fully describe the schema, just keep track
// of which major changes have been applied.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: model_type_store_schema_descriptor.proto

package sync_pb

import (
	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type ModelTypeStoreSchemaDescriptor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VersionNumber *int64 `protobuf:"varint,1,opt,name=version_number,json=versionNumber" json:"version_number,omitempty"`
}

func (x *ModelTypeStoreSchemaDescriptor) Reset() {
	*x = ModelTypeStoreSchemaDescriptor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_model_type_store_schema_descriptor_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ModelTypeStoreSchemaDescriptor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ModelTypeStoreSchemaDescriptor) ProtoMessage() {}

func (x *ModelTypeStoreSchemaDescriptor) ProtoReflect() protoreflect.Message {
	mi := &file_model_type_store_schema_descriptor_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ModelTypeStoreSchemaDescriptor.ProtoReflect.Descriptor instead.
func (*ModelTypeStoreSchemaDescriptor) Descriptor() ([]byte, []int) {
	return file_model_type_store_schema_descriptor_proto_rawDescGZIP(), []int{0}
}

func (x *ModelTypeStoreSchemaDescriptor) GetVersionNumber() int64 {
	if x != nil && x.VersionNumber != nil {
		return *x.VersionNumber
	}
	return 0
}

var File_model_type_store_schema_descriptor_proto protoreflect.FileDescriptor

var file_model_type_store_schema_descriptor_proto_rawDesc = []byte{
	0x0a, 0x28, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x73, 0x74, 0x6f,
	0x72, 0x65, 0x5f, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69,
	0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e, 0x63,
	0x5f, 0x70, 0x62, 0x22, 0x47, 0x0a, 0x1e, 0x4d, 0x6f, 0x64, 0x65, 0x6c, 0x54, 0x79, 0x70, 0x65,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x44, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x5f, 0x6e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x4e, 0x75, 0x6d, 0x62, 0x65, 0x72, 0x42, 0x2b, 0x0a, 0x25,
	0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01,
}

var (
	file_model_type_store_schema_descriptor_proto_rawDescOnce sync.Once
	file_model_type_store_schema_descriptor_proto_rawDescData = file_model_type_store_schema_descriptor_proto_rawDesc
)

func file_model_type_store_schema_descriptor_proto_rawDescGZIP() []byte {
	file_model_type_store_schema_descriptor_proto_rawDescOnce.Do(func() {
		file_model_type_store_schema_descriptor_proto_rawDescData = protoimpl.X.CompressGZIP(file_model_type_store_schema_descriptor_proto_rawDescData)
	})
	return file_model_type_store_schema_descriptor_proto_rawDescData
}

var file_model_type_store_schema_descriptor_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_model_type_store_schema_descriptor_proto_goTypes = []interface{}{
	(*ModelTypeStoreSchemaDescriptor)(nil), // 0: sync_pb.ModelTypeStoreSchemaDescriptor
}
var file_model_type_store_schema_descriptor_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_model_type_store_schema_descriptor_proto_init() }
func file_model_type_store_schema_descriptor_proto_init() {
	if File_model_type_store_schema_descriptor_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_model_type_store_schema_descriptor_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ModelTypeStoreSchemaDescriptor); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_model_type_store_schema_descriptor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_model_type_store_schema_descriptor_proto_goTypes,
		DependencyIndexes: file_model_type_store_schema_descriptor_proto_depIdxs,
		MessageInfos:      file_model_type_store_schema_descriptor_proto_msgTypes,
	}.Build()
	File_model_type_store_schema_descriptor_proto = out.File
	file_model_type_store_schema_descriptor_proto_rawDesc = nil
	file_model_type_store_schema_descriptor_proto_goTypes = nil
	file_model_type_store_schema_descriptor_proto_depIdxs = nil
}