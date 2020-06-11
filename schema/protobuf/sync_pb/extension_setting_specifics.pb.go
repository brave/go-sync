// Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Sync protocol datatype extension for an extension setting.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: extension_setting_specifics.proto

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

// Properties of extension setting sync objects.
type ExtensionSettingSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Id of the extension the setting is for.
	ExtensionId *string `protobuf:"bytes,1,opt,name=extension_id,json=extensionId" json:"extension_id,omitempty"`
	// Setting key.
	Key *string `protobuf:"bytes,2,opt,name=key" json:"key,omitempty"`
	// Setting value serialized as JSON.
	Value *string `protobuf:"bytes,3,opt,name=value" json:"value,omitempty"`
}

func (x *ExtensionSettingSpecifics) Reset() {
	*x = ExtensionSettingSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_extension_setting_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ExtensionSettingSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExtensionSettingSpecifics) ProtoMessage() {}

func (x *ExtensionSettingSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_extension_setting_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExtensionSettingSpecifics.ProtoReflect.Descriptor instead.
func (*ExtensionSettingSpecifics) Descriptor() ([]byte, []int) {
	return file_extension_setting_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *ExtensionSettingSpecifics) GetExtensionId() string {
	if x != nil && x.ExtensionId != nil {
		return *x.ExtensionId
	}
	return ""
}

func (x *ExtensionSettingSpecifics) GetKey() string {
	if x != nil && x.Key != nil {
		return *x.Key
	}
	return ""
}

func (x *ExtensionSettingSpecifics) GetValue() string {
	if x != nil && x.Value != nil {
		return *x.Value
	}
	return ""
}

var File_extension_setting_specifics_proto protoreflect.FileDescriptor

var file_extension_setting_specifics_proto_rawDesc = []byte{
	0x0a, 0x21, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x65, 0x74, 0x74,
	0x69, 0x6e, 0x67, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x22, 0x66, 0x0a, 0x19,
	0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67,
	0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x65, 0x78, 0x74,
	0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0b, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x42, 0x2b, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50,
	0x01,
}

var (
	file_extension_setting_specifics_proto_rawDescOnce sync.Once
	file_extension_setting_specifics_proto_rawDescData = file_extension_setting_specifics_proto_rawDesc
)

func file_extension_setting_specifics_proto_rawDescGZIP() []byte {
	file_extension_setting_specifics_proto_rawDescOnce.Do(func() {
		file_extension_setting_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_extension_setting_specifics_proto_rawDescData)
	})
	return file_extension_setting_specifics_proto_rawDescData
}

var file_extension_setting_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_extension_setting_specifics_proto_goTypes = []interface{}{
	(*ExtensionSettingSpecifics)(nil), // 0: sync_pb.ExtensionSettingSpecifics
}
var file_extension_setting_specifics_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_extension_setting_specifics_proto_init() }
func file_extension_setting_specifics_proto_init() {
	if File_extension_setting_specifics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_extension_setting_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ExtensionSettingSpecifics); i {
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
			RawDescriptor: file_extension_setting_specifics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_extension_setting_specifics_proto_goTypes,
		DependencyIndexes: file_extension_setting_specifics_proto_depIdxs,
		MessageInfos:      file_extension_setting_specifics_proto_msgTypes,
	}.Build()
	File_extension_setting_specifics_proto = out.File
	file_extension_setting_specifics_proto_rawDesc = nil
	file_extension_setting_specifics_proto_goTypes = nil
	file_extension_setting_specifics_proto_depIdxs = nil
}