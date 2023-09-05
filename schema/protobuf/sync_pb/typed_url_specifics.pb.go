// Copyright 2012 The Chromium Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Sync protocol datatype extension for typed urls.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.1
// source: typed_url_specifics.proto

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

// Properties of typed_url sync objects - fields correspond to similarly named
// fields in history::URLRow.
type TypedUrlSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Actual URL that was typed.
	Url *string `protobuf:"bytes,1,opt,name=url" json:"url,omitempty"`
	// Title of the page corresponding to this URL.
	Title *string `protobuf:"bytes,2,opt,name=title" json:"title,omitempty"`
	// True if the URL should NOT be used for auto-complete.
	Hidden *bool `protobuf:"varint,4,opt,name=hidden" json:"hidden,omitempty"`
	// Timestamps for all visits to this URL.
	Visits []int64 `protobuf:"varint,7,rep,packed,name=visits" json:"visits,omitempty"`
	// The PageTransition::Type for each of the visits in the |visit| array. Both
	// arrays must be the same length.
	VisitTransitions []int32 `protobuf:"varint,8,rep,packed,name=visit_transitions,json=visitTransitions" json:"visit_transitions,omitempty"`
}

func (x *TypedUrlSpecifics) Reset() {
	*x = TypedUrlSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_typed_url_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypedUrlSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypedUrlSpecifics) ProtoMessage() {}

func (x *TypedUrlSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_typed_url_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypedUrlSpecifics.ProtoReflect.Descriptor instead.
func (*TypedUrlSpecifics) Descriptor() ([]byte, []int) {
	return file_typed_url_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *TypedUrlSpecifics) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *TypedUrlSpecifics) GetTitle() string {
	if x != nil && x.Title != nil {
		return *x.Title
	}
	return ""
}

func (x *TypedUrlSpecifics) GetHidden() bool {
	if x != nil && x.Hidden != nil {
		return *x.Hidden
	}
	return false
}

func (x *TypedUrlSpecifics) GetVisits() []int64 {
	if x != nil {
		return x.Visits
	}
	return nil
}

func (x *TypedUrlSpecifics) GetVisitTransitions() []int32 {
	if x != nil {
		return x.VisitTransitions
	}
	return nil
}

var File_typed_url_specifics_proto protoreflect.FileDescriptor

var file_typed_url_specifics_proto_rawDesc = []byte{
	0x0a, 0x19, 0x74, 0x79, 0x70, 0x65, 0x64, 0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x73, 0x70, 0x65, 0x63,
	0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e,
	0x63, 0x5f, 0x70, 0x62, 0x22, 0xd5, 0x01, 0x0a, 0x11, 0x54, 0x79, 0x70, 0x65, 0x64, 0x55, 0x72,
	0x6c, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72,
	0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74,
	0x6c, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x06, 0x68, 0x69, 0x64, 0x64, 0x65, 0x6e, 0x12, 0x1a, 0x0a, 0x06, 0x76, 0x69,
	0x73, 0x69, 0x74, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x03, 0x42, 0x02, 0x10, 0x01, 0x52, 0x06,
	0x76, 0x69, 0x73, 0x69, 0x74, 0x73, 0x12, 0x2f, 0x0a, 0x11, 0x76, 0x69, 0x73, 0x69, 0x74, 0x5f,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28,
	0x05, 0x42, 0x02, 0x10, 0x01, 0x52, 0x10, 0x76, 0x69, 0x73, 0x69, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x4a, 0x04, 0x08,
	0x05, 0x10, 0x06, 0x4a, 0x04, 0x08, 0x06, 0x10, 0x07, 0x52, 0x0b, 0x74, 0x79, 0x70, 0x65, 0x64,
	0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x05, 0x76, 0x69, 0x73, 0x69, 0x74, 0x52, 0x0d, 0x76,
	0x69, 0x73, 0x69, 0x74, 0x65, 0x64, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x36, 0x0a, 0x25,
	0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d,
	0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01, 0x5a, 0x09, 0x2e, 0x2f, 0x73, 0x79, 0x6e,
	0x63, 0x5f, 0x70, 0x62,
}

var (
	file_typed_url_specifics_proto_rawDescOnce sync.Once
	file_typed_url_specifics_proto_rawDescData = file_typed_url_specifics_proto_rawDesc
)

func file_typed_url_specifics_proto_rawDescGZIP() []byte {
	file_typed_url_specifics_proto_rawDescOnce.Do(func() {
		file_typed_url_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_typed_url_specifics_proto_rawDescData)
	})
	return file_typed_url_specifics_proto_rawDescData
}

var file_typed_url_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_typed_url_specifics_proto_goTypes = []interface{}{
	(*TypedUrlSpecifics)(nil), // 0: sync_pb.TypedUrlSpecifics
}
var file_typed_url_specifics_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_typed_url_specifics_proto_init() }
func file_typed_url_specifics_proto_init() {
	if File_typed_url_specifics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_typed_url_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypedUrlSpecifics); i {
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
			RawDescriptor: file_typed_url_specifics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_typed_url_specifics_proto_goTypes,
		DependencyIndexes: file_typed_url_specifics_proto_depIdxs,
		MessageInfos:      file_typed_url_specifics_proto_msgTypes,
	}.Build()
	File_typed_url_specifics_proto = out.File
	file_typed_url_specifics_proto_rawDesc = nil
	file_typed_url_specifics_proto_goTypes = nil
	file_typed_url_specifics_proto_depIdxs = nil
}
