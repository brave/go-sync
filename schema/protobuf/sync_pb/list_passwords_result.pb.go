// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.1
// source: list_passwords_result.proto

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

// Response to a request sent to Google Mobile Services to request a list of
// passwords.
// ATTENTION(crbug.com/1330911): This proto is being moved to
// components/password_manager/core/browser/protocol folder. Two files exist
// while the migration is in process, this file will be deleted when the
// migration is over. IF YOU MODIFY THIS FILE, PLEASE ALSO MODIFY THE COPY IN
// components/password_manager.
type ListPasswordsResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The list of password entries and corresponding additional info.
	PasswordData []*PasswordWithLocalData `protobuf:"bytes,1,rep,name=password_data,json=passwordData" json:"password_data,omitempty"`
}

func (x *ListPasswordsResult) Reset() {
	*x = ListPasswordsResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_list_passwords_result_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPasswordsResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPasswordsResult) ProtoMessage() {}

func (x *ListPasswordsResult) ProtoReflect() protoreflect.Message {
	mi := &file_list_passwords_result_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPasswordsResult.ProtoReflect.Descriptor instead.
func (*ListPasswordsResult) Descriptor() ([]byte, []int) {
	return file_list_passwords_result_proto_rawDescGZIP(), []int{0}
}

func (x *ListPasswordsResult) GetPasswordData() []*PasswordWithLocalData {
	if x != nil {
		return x.PasswordData
	}
	return nil
}

var File_list_passwords_result_proto protoreflect.FileDescriptor

var file_list_passwords_result_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73,
	0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x1a, 0x1e, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x5f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x5f, 0x64, 0x61, 0x74, 0x61,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5a, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x43, 0x0a,
	0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x57, 0x69, 0x74, 0x68, 0x4c, 0x6f, 0x63, 0x61, 0x6c,
	0x44, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61,
	0x74, 0x61, 0x42, 0x2b, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69,
	0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01,
}

var (
	file_list_passwords_result_proto_rawDescOnce sync.Once
	file_list_passwords_result_proto_rawDescData = file_list_passwords_result_proto_rawDesc
)

func file_list_passwords_result_proto_rawDescGZIP() []byte {
	file_list_passwords_result_proto_rawDescOnce.Do(func() {
		file_list_passwords_result_proto_rawDescData = protoimpl.X.CompressGZIP(file_list_passwords_result_proto_rawDescData)
	})
	return file_list_passwords_result_proto_rawDescData
}

var file_list_passwords_result_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_list_passwords_result_proto_goTypes = []interface{}{
	(*ListPasswordsResult)(nil),   // 0: sync_pb.ListPasswordsResult
	(*PasswordWithLocalData)(nil), // 1: sync_pb.PasswordWithLocalData
}
var file_list_passwords_result_proto_depIdxs = []int32{
	1, // 0: sync_pb.ListPasswordsResult.password_data:type_name -> sync_pb.PasswordWithLocalData
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_list_passwords_result_proto_init() }
func file_list_passwords_result_proto_init() {
	if File_list_passwords_result_proto != nil {
		return
	}
	file_password_with_local_data_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_list_passwords_result_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPasswordsResult); i {
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
			RawDescriptor: file_list_passwords_result_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_list_passwords_result_proto_goTypes,
		DependencyIndexes: file_list_passwords_result_proto_depIdxs,
		MessageInfos:      file_list_passwords_result_proto_msgTypes,
	}.Build()
	File_list_passwords_result_proto = out.File
	file_list_passwords_result_proto_rawDesc = nil
	file_list_passwords_result_proto_goTypes = nil
	file_list_passwords_result_proto_depIdxs = nil
}