// Copyright 2023 The Chromium Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.1
// source: password_sharing_recipients.proto

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

type PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult int32

const (
	PasswordSharingRecipientsResponse_UNKNOWN PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult = 0
	// The user is a member of a family and the request succeeded.
	PasswordSharingRecipientsResponse_SUCCESS PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult = 1
	// Not a family member, used to distinguish from a family with
	// only one member.
	PasswordSharingRecipientsResponse_NOT_FAMILY_MEMBER PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult = 2
)

// Enum value maps for PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult.
var (
	PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult_name = map[int32]string{
		0: "UNKNOWN",
		1: "SUCCESS",
		2: "NOT_FAMILY_MEMBER",
	}
	PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult_value = map[string]int32{
		"UNKNOWN":           0,
		"SUCCESS":           1,
		"NOT_FAMILY_MEMBER": 2,
	}
)

func (x PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) Enum() *PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult {
	p := new(PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult)
	*p = x
	return p
}

func (x PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) Descriptor() protoreflect.EnumDescriptor {
	return file_password_sharing_recipients_proto_enumTypes[0].Descriptor()
}

func (PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) Type() protoreflect.EnumType {
	return &file_password_sharing_recipients_proto_enumTypes[0]
}

func (x PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult(num)
	return nil
}

// Deprecated: Use PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult.Descriptor instead.
func (PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult) EnumDescriptor() ([]byte, []int) {
	return file_password_sharing_recipients_proto_rawDescGZIP(), []int{1, 0}
}

// A message to obtain a list of recipients for sending a password.
type PasswordSharingRecipientsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PasswordSharingRecipientsRequest) Reset() {
	*x = PasswordSharingRecipientsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_recipients_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasswordSharingRecipientsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordSharingRecipientsRequest) ProtoMessage() {}

func (x *PasswordSharingRecipientsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_recipients_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordSharingRecipientsRequest.ProtoReflect.Descriptor instead.
func (*PasswordSharingRecipientsRequest) Descriptor() ([]byte, []int) {
	return file_password_sharing_recipients_proto_rawDescGZIP(), []int{0}
}

type PasswordSharingRecipientsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result *PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult `protobuf:"varint,1,opt,name=result,enum=sync_pb.PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult" json:"result,omitempty"`
	// List of possible recipients for sending a password. Note that public key
	// may be absent if a recipient can’t receive a password (e.g. due to an older
	// Chrome version).
	Recipients []*UserInfo `protobuf:"bytes,2,rep,name=recipients" json:"recipients,omitempty"`
}

func (x *PasswordSharingRecipientsResponse) Reset() {
	*x = PasswordSharingRecipientsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_recipients_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasswordSharingRecipientsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordSharingRecipientsResponse) ProtoMessage() {}

func (x *PasswordSharingRecipientsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_recipients_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordSharingRecipientsResponse.ProtoReflect.Descriptor instead.
func (*PasswordSharingRecipientsResponse) Descriptor() ([]byte, []int) {
	return file_password_sharing_recipients_proto_rawDescGZIP(), []int{1}
}

func (x *PasswordSharingRecipientsResponse) GetResult() PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult {
	if x != nil && x.Result != nil {
		return *x.Result
	}
	return PasswordSharingRecipientsResponse_UNKNOWN
}

func (x *PasswordSharingRecipientsResponse) GetRecipients() []*UserInfo {
	if x != nil {
		return x.Recipients
	}
	return nil
}

var File_password_sharing_recipients_proto protoreflect.FileDescriptor

var file_password_sharing_recipients_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x69,
	0x6e, 0x67, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x1a, 0x2b, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69,
	0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66,
	0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x22, 0x0a, 0x20, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63, 0x69,
	0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x8e, 0x02,
	0x0a, 0x21, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e,
	0x67, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x62, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x4a, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x52,
	0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x31, 0x0a, 0x0a, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a,
	0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x52, 0x0a, 0x1f, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x52, 0x65, 0x63,
	0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x0b, 0x0a,
	0x07, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x53, 0x55,
	0x43, 0x43, 0x45, 0x53, 0x53, 0x10, 0x01, 0x12, 0x15, 0x0a, 0x11, 0x4e, 0x4f, 0x54, 0x5f, 0x46,
	0x41, 0x4d, 0x49, 0x4c, 0x59, 0x5f, 0x4d, 0x45, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x02, 0x42, 0x36,
	0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01, 0x5a, 0x09, 0x2e, 0x2f, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62,
}

var (
	file_password_sharing_recipients_proto_rawDescOnce sync.Once
	file_password_sharing_recipients_proto_rawDescData = file_password_sharing_recipients_proto_rawDesc
)

func file_password_sharing_recipients_proto_rawDescGZIP() []byte {
	file_password_sharing_recipients_proto_rawDescOnce.Do(func() {
		file_password_sharing_recipients_proto_rawDescData = protoimpl.X.CompressGZIP(file_password_sharing_recipients_proto_rawDescData)
	})
	return file_password_sharing_recipients_proto_rawDescData
}

var file_password_sharing_recipients_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_password_sharing_recipients_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_password_sharing_recipients_proto_goTypes = []interface{}{
	(PasswordSharingRecipientsResponse_PasswordSharingRecipientsResult)(0), // 0: sync_pb.PasswordSharingRecipientsResponse.PasswordSharingRecipientsResult
	(*PasswordSharingRecipientsRequest)(nil),                               // 1: sync_pb.PasswordSharingRecipientsRequest
	(*PasswordSharingRecipientsResponse)(nil),                              // 2: sync_pb.PasswordSharingRecipientsResponse
	(*UserInfo)(nil), // 3: sync_pb.UserInfo
}
var file_password_sharing_recipients_proto_depIdxs = []int32{
	0, // 0: sync_pb.PasswordSharingRecipientsResponse.result:type_name -> sync_pb.PasswordSharingRecipientsResponse.PasswordSharingRecipientsResult
	3, // 1: sync_pb.PasswordSharingRecipientsResponse.recipients:type_name -> sync_pb.UserInfo
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_password_sharing_recipients_proto_init() }
func file_password_sharing_recipients_proto_init() {
	if File_password_sharing_recipients_proto != nil {
		return
	}
	file_password_sharing_invitation_specifics_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_password_sharing_recipients_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasswordSharingRecipientsRequest); i {
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
		file_password_sharing_recipients_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasswordSharingRecipientsResponse); i {
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
			RawDescriptor: file_password_sharing_recipients_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_password_sharing_recipients_proto_goTypes,
		DependencyIndexes: file_password_sharing_recipients_proto_depIdxs,
		EnumInfos:         file_password_sharing_recipients_proto_enumTypes,
		MessageInfos:      file_password_sharing_recipients_proto_msgTypes,
	}.Build()
	File_password_sharing_recipients_proto = out.File
	file_password_sharing_recipients_proto_rawDesc = nil
	file_password_sharing_recipients_proto_goTypes = nil
	file_password_sharing_recipients_proto_depIdxs = nil
}
