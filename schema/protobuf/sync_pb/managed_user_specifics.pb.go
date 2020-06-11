// Copyright 2013 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Sync protocol datatype extension for managed user settings.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: managed_user_specifics.proto

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

// Properties of managed user sync objects.
type ManagedUserSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A randomly-generated identifier for the managed user.
	Id *string `protobuf:"bytes,1,opt,name=id" json:"id,omitempty"`
	// The human-visible name of the managed user
	Name *string `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	// This flag is set by the server to acknowledge that it has committed a
	// newly created managed user.
	Acknowledged *bool `protobuf:"varint,3,opt,name=acknowledged,def=0" json:"acknowledged,omitempty"`
	// Master key for managed user cryptohome.
	MasterKey *string `protobuf:"bytes,4,opt,name=master_key,json=masterKey" json:"master_key,omitempty"`
	// A string representing the index of the supervised user avatar on Chrome.
	// It has the following format:
	// "chrome-avatar-index:INDEX" where INDEX is an integer.
	ChromeAvatar *string `protobuf:"bytes,5,opt,name=chrome_avatar,json=chromeAvatar" json:"chrome_avatar,omitempty"`
	// A string representing the index of the supervised user avatar on Chrome OS.
	// It has the following format:
	// "chromeos-avatar-index:INDEX" where INDEX is an integer.
	ChromeosAvatar *string `protobuf:"bytes,6,opt,name=chromeos_avatar,json=chromeosAvatar" json:"chromeos_avatar,omitempty"`
	// Key for signing supervised user's password.
	PasswordSignatureKey *string `protobuf:"bytes,7,opt,name=password_signature_key,json=passwordSignatureKey" json:"password_signature_key,omitempty"`
	// Key for encrypting supervised user's password.
	PasswordEncryptionKey *string `protobuf:"bytes,8,opt,name=password_encryption_key,json=passwordEncryptionKey" json:"password_encryption_key,omitempty"`
}

// Default values for ManagedUserSpecifics fields.
const (
	Default_ManagedUserSpecifics_Acknowledged = bool(false)
)

func (x *ManagedUserSpecifics) Reset() {
	*x = ManagedUserSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_managed_user_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ManagedUserSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ManagedUserSpecifics) ProtoMessage() {}

func (x *ManagedUserSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_managed_user_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ManagedUserSpecifics.ProtoReflect.Descriptor instead.
func (*ManagedUserSpecifics) Descriptor() ([]byte, []int) {
	return file_managed_user_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *ManagedUserSpecifics) GetId() string {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return ""
}

func (x *ManagedUserSpecifics) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *ManagedUserSpecifics) GetAcknowledged() bool {
	if x != nil && x.Acknowledged != nil {
		return *x.Acknowledged
	}
	return Default_ManagedUserSpecifics_Acknowledged
}

func (x *ManagedUserSpecifics) GetMasterKey() string {
	if x != nil && x.MasterKey != nil {
		return *x.MasterKey
	}
	return ""
}

func (x *ManagedUserSpecifics) GetChromeAvatar() string {
	if x != nil && x.ChromeAvatar != nil {
		return *x.ChromeAvatar
	}
	return ""
}

func (x *ManagedUserSpecifics) GetChromeosAvatar() string {
	if x != nil && x.ChromeosAvatar != nil {
		return *x.ChromeosAvatar
	}
	return ""
}

func (x *ManagedUserSpecifics) GetPasswordSignatureKey() string {
	if x != nil && x.PasswordSignatureKey != nil {
		return *x.PasswordSignatureKey
	}
	return ""
}

func (x *ManagedUserSpecifics) GetPasswordEncryptionKey() string {
	if x != nil && x.PasswordEncryptionKey != nil {
		return *x.PasswordEncryptionKey
	}
	return ""
}

var File_managed_user_specifics_proto protoreflect.FileDescriptor

var file_managed_user_specifics_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x64, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x22, 0xc0, 0x02, 0x0a, 0x14, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x64, 0x55, 0x73, 0x65, 0x72, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x0c, 0x61, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65,
	0x64, 0x67, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73,
	0x65, 0x52, 0x0c, 0x61, 0x63, 0x6b, 0x6e, 0x6f, 0x77, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x12, 0x23,
	0x0a, 0x0d, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x5f, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x41, 0x76, 0x61,
	0x74, 0x61, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x5f,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x65, 0x6f, 0x73, 0x41, 0x76, 0x61, 0x74, 0x61, 0x72, 0x12, 0x34, 0x0a, 0x16,
	0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x14, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x4b,
	0x65, 0x79, 0x12, 0x36, 0x0a, 0x17, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x15, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x45, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x42, 0x2b, 0x0a, 0x25, 0x6f, 0x72,
	0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01,
}

var (
	file_managed_user_specifics_proto_rawDescOnce sync.Once
	file_managed_user_specifics_proto_rawDescData = file_managed_user_specifics_proto_rawDesc
)

func file_managed_user_specifics_proto_rawDescGZIP() []byte {
	file_managed_user_specifics_proto_rawDescOnce.Do(func() {
		file_managed_user_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_managed_user_specifics_proto_rawDescData)
	})
	return file_managed_user_specifics_proto_rawDescData
}

var file_managed_user_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_managed_user_specifics_proto_goTypes = []interface{}{
	(*ManagedUserSpecifics)(nil), // 0: sync_pb.ManagedUserSpecifics
}
var file_managed_user_specifics_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_managed_user_specifics_proto_init() }
func file_managed_user_specifics_proto_init() {
	if File_managed_user_specifics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_managed_user_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ManagedUserSpecifics); i {
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
			RawDescriptor: file_managed_user_specifics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_managed_user_specifics_proto_goTypes,
		DependencyIndexes: file_managed_user_specifics_proto_depIdxs,
		MessageInfos:      file_managed_user_specifics_proto_msgTypes,
	}.Build()
	File_managed_user_specifics_proto = out.File
	file_managed_user_specifics_proto_rawDesc = nil
	file_managed_user_specifics_proto_goTypes = nil
	file_managed_user_specifics_proto_depIdxs = nil
}