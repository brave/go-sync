// Copyright 2019 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: web_app_specifics.proto

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

// This enum should be a subset of the DisplayMode enum in
// chrome/browser/web_applications/proto/web_app.proto and
// third_party/blink/public/mojom/manifest/display_mode.mojom
type WebAppSpecifics_UserDisplayMode int32

const (
	// UNDEFINED is never serialized.
	WebAppSpecifics_BROWSER WebAppSpecifics_UserDisplayMode = 1
	// MINIMAL_UI is never serialized.
	WebAppSpecifics_STANDALONE WebAppSpecifics_UserDisplayMode = 3 // FULLSCREEN is never serialized.
)

// Enum value maps for WebAppSpecifics_UserDisplayMode.
var (
	WebAppSpecifics_UserDisplayMode_name = map[int32]string{
		1: "BROWSER",
		3: "STANDALONE",
	}
	WebAppSpecifics_UserDisplayMode_value = map[string]int32{
		"BROWSER":    1,
		"STANDALONE": 3,
	}
)

func (x WebAppSpecifics_UserDisplayMode) Enum() *WebAppSpecifics_UserDisplayMode {
	p := new(WebAppSpecifics_UserDisplayMode)
	*p = x
	return p
}

func (x WebAppSpecifics_UserDisplayMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (WebAppSpecifics_UserDisplayMode) Descriptor() protoreflect.EnumDescriptor {
	return file_web_app_specifics_proto_enumTypes[0].Descriptor()
}

func (WebAppSpecifics_UserDisplayMode) Type() protoreflect.EnumType {
	return &file_web_app_specifics_proto_enumTypes[0]
}

func (x WebAppSpecifics_UserDisplayMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *WebAppSpecifics_UserDisplayMode) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = WebAppSpecifics_UserDisplayMode(num)
	return nil
}

// Deprecated: Use WebAppSpecifics_UserDisplayMode.Descriptor instead.
func (WebAppSpecifics_UserDisplayMode) EnumDescriptor() ([]byte, []int) {
	return file_web_app_specifics_proto_rawDescGZIP(), []int{0, 0}
}

// WebApp data. This is a synced part of
// chrome/browser/web_applications/proto/web_app.proto data.
type WebAppSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LaunchUrl       *string                          `protobuf:"bytes,1,opt,name=launch_url,json=launchUrl" json:"launch_url,omitempty"`
	Name            *string                          `protobuf:"bytes,2,opt,name=name" json:"name,omitempty"`
	UserDisplayMode *WebAppSpecifics_UserDisplayMode `protobuf:"varint,3,opt,name=user_display_mode,json=userDisplayMode,enum=sync_pb.WebAppSpecifics_UserDisplayMode" json:"user_display_mode,omitempty"`
	ThemeColor      *uint32                          `protobuf:"varint,4,opt,name=theme_color,json=themeColor" json:"theme_color,omitempty"`
}

func (x *WebAppSpecifics) Reset() {
	*x = WebAppSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_web_app_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebAppSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebAppSpecifics) ProtoMessage() {}

func (x *WebAppSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_web_app_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebAppSpecifics.ProtoReflect.Descriptor instead.
func (*WebAppSpecifics) Descriptor() ([]byte, []int) {
	return file_web_app_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *WebAppSpecifics) GetLaunchUrl() string {
	if x != nil && x.LaunchUrl != nil {
		return *x.LaunchUrl
	}
	return ""
}

func (x *WebAppSpecifics) GetName() string {
	if x != nil && x.Name != nil {
		return *x.Name
	}
	return ""
}

func (x *WebAppSpecifics) GetUserDisplayMode() WebAppSpecifics_UserDisplayMode {
	if x != nil && x.UserDisplayMode != nil {
		return *x.UserDisplayMode
	}
	return WebAppSpecifics_BROWSER
}

func (x *WebAppSpecifics) GetThemeColor() uint32 {
	if x != nil && x.ThemeColor != nil {
		return *x.ThemeColor
	}
	return 0
}

var File_web_app_specifics_proto protoreflect.FileDescriptor

var file_web_app_specifics_proto_rawDesc = []byte{
	0x0a, 0x17, 0x77, 0x65, 0x62, 0x5f, 0x61, 0x70, 0x70, 0x5f, 0x73, 0x70, 0x65, 0x63, 0x69, 0x66,
	0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x5f,
	0x70, 0x62, 0x22, 0xeb, 0x01, 0x0a, 0x0f, 0x57, 0x65, 0x62, 0x41, 0x70, 0x70, 0x53, 0x70, 0x65,
	0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x75, 0x6e, 0x63, 0x68,
	0x5f, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x61, 0x75, 0x6e,
	0x63, 0x68, 0x55, 0x72, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x54, 0x0a, 0x11, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x28, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x57,
	0x65, 0x62, 0x41, 0x70, 0x70, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x55,
	0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x0f,
	0x75, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4d, 0x6f, 0x64, 0x65, 0x12,
	0x1f, 0x0a, 0x0b, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x5f, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x74, 0x68, 0x65, 0x6d, 0x65, 0x43, 0x6f, 0x6c, 0x6f, 0x72,
	0x22, 0x2e, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4d,
	0x6f, 0x64, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x42, 0x52, 0x4f, 0x57, 0x53, 0x45, 0x52, 0x10, 0x01,
	0x12, 0x0e, 0x0a, 0x0a, 0x53, 0x54, 0x41, 0x4e, 0x44, 0x41, 0x4c, 0x4f, 0x4e, 0x45, 0x10, 0x03,
	0x42, 0x2b, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d,
	0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01,
}

var (
	file_web_app_specifics_proto_rawDescOnce sync.Once
	file_web_app_specifics_proto_rawDescData = file_web_app_specifics_proto_rawDesc
)

func file_web_app_specifics_proto_rawDescGZIP() []byte {
	file_web_app_specifics_proto_rawDescOnce.Do(func() {
		file_web_app_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_web_app_specifics_proto_rawDescData)
	})
	return file_web_app_specifics_proto_rawDescData
}

var file_web_app_specifics_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_web_app_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_web_app_specifics_proto_goTypes = []interface{}{
	(WebAppSpecifics_UserDisplayMode)(0), // 0: sync_pb.WebAppSpecifics.UserDisplayMode
	(*WebAppSpecifics)(nil),              // 1: sync_pb.WebAppSpecifics
}
var file_web_app_specifics_proto_depIdxs = []int32{
	0, // 0: sync_pb.WebAppSpecifics.user_display_mode:type_name -> sync_pb.WebAppSpecifics.UserDisplayMode
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_web_app_specifics_proto_init() }
func file_web_app_specifics_proto_init() {
	if File_web_app_specifics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_web_app_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebAppSpecifics); i {
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
			RawDescriptor: file_web_app_specifics_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_web_app_specifics_proto_goTypes,
		DependencyIndexes: file_web_app_specifics_proto_depIdxs,
		EnumInfos:         file_web_app_specifics_proto_enumTypes,
		MessageInfos:      file_web_app_specifics_proto_msgTypes,
	}.Build()
	File_web_app_specifics_proto = out.File
	file_web_app_specifics_proto_rawDesc = nil
	file_web_app_specifics_proto_goTypes = nil
	file_web_app_specifics_proto_depIdxs = nil
}