// Copyright 2018 The Chromium Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.6.1
// source: persisted_entity_data.proto

package sync_pb

import (
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

// Sync proto to store entity data similar to what the legacy Directory used
// to store, used to persist data locally and never sent through the wire.
//
// Because it's conceptually similar to SyncEntity (actual protocol) and it's
// unclear how big this'll grow, we've kept compatibility with SyncEntity by
// using the same field numbers.
type PersistedEntityData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// See corresponding fields in SyncEntity for details.
	Name      *string          `protobuf:"bytes,8,opt,name=name" json:"name,omitempty"`
	Specifics *EntitySpecifics `protobuf:"bytes,21,opt,name=specifics" json:"specifics,omitempty"`
}

func (x *PersistedEntityData) Reset() {
	*x = PersistedEntityData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_persisted_entity_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PersistedEntityData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PersistedEntityData) ProtoMessage() {}

func (x *PersistedEntityData) ProtoReflect() protoreflect.Message {
	mi := &file_persisted_entity_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PersistedEntityData.ProtoReflect.Descriptor instead.
func (*PersistedEntityData) Descriptor() ([]byte, []int) {
	return file_persisted_entity_data_proto_rawDescGZIP(), []int{0}
}

func (x *PersistedEntityData) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *PersistedEntityData) GetSpecifics() *EntitySpecifics {
	if x != nil {
		return x.Specifics
	}
	return nil
}

var File_persisted_entity_data_proto protoreflect.FileDescriptor

var file_persisted_entity_data_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x70, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x64, 0x5f, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x1a, 0x16, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x61,
	0x0a, 0x13, 0x50, 0x65, 0x72, 0x73, 0x69, 0x73, 0x74, 0x65, 0x64, 0x45, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x44, 0x61, 0x74, 0x61, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x36, 0x0a, 0x09, 0x73, 0x70, 0x65,
	0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x53, 0x70, 0x65,
	0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x52, 0x09, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63,
	0x73, 0x42, 0x36, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75,
	0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01, 0x5a, 0x09,
	0x2e, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62,
}

var (
	file_persisted_entity_data_proto_rawDescOnce sync.Once
	file_persisted_entity_data_proto_rawDescData = file_persisted_entity_data_proto_rawDesc
)

func file_persisted_entity_data_proto_rawDescGZIP() []byte {
	file_persisted_entity_data_proto_rawDescOnce.Do(func() {
		file_persisted_entity_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_persisted_entity_data_proto_rawDescData)
	})
	return file_persisted_entity_data_proto_rawDescData
}

var file_persisted_entity_data_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_persisted_entity_data_proto_goTypes = []interface{}{
	(*PersistedEntityData)(nil), // 0: sync_pb.PersistedEntityData
	(*EntitySpecifics)(nil),     // 1: sync_pb.EntitySpecifics
}
var file_persisted_entity_data_proto_depIdxs = []int32{
	1, // 0: sync_pb.PersistedEntityData.specifics:type_name -> sync_pb.EntitySpecifics
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_persisted_entity_data_proto_init() }
func file_persisted_entity_data_proto_init() {
	if File_persisted_entity_data_proto != nil {
		return
	}
	file_entity_specifics_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_persisted_entity_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PersistedEntityData); i {
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
			RawDescriptor: file_persisted_entity_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_persisted_entity_data_proto_goTypes,
		DependencyIndexes: file_persisted_entity_data_proto_depIdxs,
		MessageInfos:      file_persisted_entity_data_proto_msgTypes,
	}.Build()
	File_persisted_entity_data_proto = out.File
	file_persisted_entity_data_proto_rawDesc = nil
	file_persisted_entity_data_proto_goTypes = nil
	file_persisted_entity_data_proto_depIdxs = nil
}
