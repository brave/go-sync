// Copyright 2012 The Chromium Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Sync protocol for debug info clients can send to the sync server.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.1
// source: client_debug_info.proto

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

// Per-type hint information.
type TypeHint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The data type this hint applied to.
	DataTypeId *int32 `protobuf:"varint,1,opt,name=data_type_id,json=dataTypeId" json:"data_type_id,omitempty"`
	// Whether or not a valid hint is provided.
	HasValidHint *bool `protobuf:"varint,2,opt,name=has_valid_hint,json=hasValidHint" json:"has_valid_hint,omitempty"`
}

func (x *TypeHint) Reset() {
	*x = TypeHint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_debug_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TypeHint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TypeHint) ProtoMessage() {}

func (x *TypeHint) ProtoReflect() protoreflect.Message {
	mi := &file_client_debug_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TypeHint.ProtoReflect.Descriptor instead.
func (*TypeHint) Descriptor() ([]byte, []int) {
	return file_client_debug_info_proto_rawDescGZIP(), []int{0}
}

func (x *TypeHint) GetDataTypeId() int32 {
	if x != nil && x.DataTypeId != nil {
		return *x.DataTypeId
	}
	return 0
}

func (x *TypeHint) GetHasValidHint() bool {
	if x != nil && x.HasValidHint != nil {
		return *x.HasValidHint
	}
	return false
}

// The additional info here is from the StatusController. They get sent when
// the event SYNC_CYCLE_COMPLETED  is sent.
type SyncCycleCompletedEventInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// These new conflict counters replace the ones above.
	// TODO(crbug.com/1315573): Deprecated in M103.
	//
	// Deprecated: Marked as deprecated in client_debug_info.proto.
	NumEncryptionConflicts *int32 `protobuf:"varint,4,opt,name=num_encryption_conflicts,json=numEncryptionConflicts" json:"num_encryption_conflicts,omitempty"`
	// Deprecated: Marked as deprecated in client_debug_info.proto.
	NumHierarchyConflicts *int32 `protobuf:"varint,5,opt,name=num_hierarchy_conflicts,json=numHierarchyConflicts" json:"num_hierarchy_conflicts,omitempty"`
	NumSimpleConflicts    *int32 `protobuf:"varint,6,opt,name=num_simple_conflicts,json=numSimpleConflicts" json:"num_simple_conflicts,omitempty"` // No longer sent since M24.
	NumServerConflicts    *int32 `protobuf:"varint,7,opt,name=num_server_conflicts,json=numServerConflicts" json:"num_server_conflicts,omitempty"`
	// Counts to track the effective usefulness of our GetUpdate requests.
	NumUpdatesDownloaded *int32 `protobuf:"varint,8,opt,name=num_updates_downloaded,json=numUpdatesDownloaded" json:"num_updates_downloaded,omitempty"`
	// TODO(crbug.com/1315573): Deprecated in M103.
	//
	// Deprecated: Marked as deprecated in client_debug_info.proto.
	NumReflectedUpdatesDownloaded *int32 `protobuf:"varint,9,opt,name=num_reflected_updates_downloaded,json=numReflectedUpdatesDownloaded" json:"num_reflected_updates_downloaded,omitempty"`
	// |caller_info| was mostly replaced by |get_updates_origin|; now it only
	// contains the |notifications_enabled| flag.
	CallerInfo       *GetUpdatesCallerInfo       `protobuf:"bytes,10,opt,name=caller_info,json=callerInfo" json:"caller_info,omitempty"`
	GetUpdatesOrigin *SyncEnums_GetUpdatesOrigin `protobuf:"varint,12,opt,name=get_updates_origin,json=getUpdatesOrigin,enum=sync_pb.SyncEnums_GetUpdatesOrigin" json:"get_updates_origin,omitempty"`
}

func (x *SyncCycleCompletedEventInfo) Reset() {
	*x = SyncCycleCompletedEventInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_debug_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncCycleCompletedEventInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncCycleCompletedEventInfo) ProtoMessage() {}

func (x *SyncCycleCompletedEventInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_debug_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncCycleCompletedEventInfo.ProtoReflect.Descriptor instead.
func (*SyncCycleCompletedEventInfo) Descriptor() ([]byte, []int) {
	return file_client_debug_info_proto_rawDescGZIP(), []int{1}
}

// Deprecated: Marked as deprecated in client_debug_info.proto.
func (x *SyncCycleCompletedEventInfo) GetNumEncryptionConflicts() int32 {
	if x != nil && x.NumEncryptionConflicts != nil {
		return *x.NumEncryptionConflicts
	}
	return 0
}

// Deprecated: Marked as deprecated in client_debug_info.proto.
func (x *SyncCycleCompletedEventInfo) GetNumHierarchyConflicts() int32 {
	if x != nil && x.NumHierarchyConflicts != nil {
		return *x.NumHierarchyConflicts
	}
	return 0
}

func (x *SyncCycleCompletedEventInfo) GetNumSimpleConflicts() int32 {
	if x != nil && x.NumSimpleConflicts != nil {
		return *x.NumSimpleConflicts
	}
	return 0
}

func (x *SyncCycleCompletedEventInfo) GetNumServerConflicts() int32 {
	if x != nil && x.NumServerConflicts != nil {
		return *x.NumServerConflicts
	}
	return 0
}

func (x *SyncCycleCompletedEventInfo) GetNumUpdatesDownloaded() int32 {
	if x != nil && x.NumUpdatesDownloaded != nil {
		return *x.NumUpdatesDownloaded
	}
	return 0
}

// Deprecated: Marked as deprecated in client_debug_info.proto.
func (x *SyncCycleCompletedEventInfo) GetNumReflectedUpdatesDownloaded() int32 {
	if x != nil && x.NumReflectedUpdatesDownloaded != nil {
		return *x.NumReflectedUpdatesDownloaded
	}
	return 0
}

func (x *SyncCycleCompletedEventInfo) GetCallerInfo() *GetUpdatesCallerInfo {
	if x != nil {
		return x.CallerInfo
	}
	return nil
}

func (x *SyncCycleCompletedEventInfo) GetGetUpdatesOrigin() SyncEnums_GetUpdatesOrigin {
	if x != nil && x.GetUpdatesOrigin != nil {
		return *x.GetUpdatesOrigin
	}
	return SyncEnums_UNKNOWN_ORIGIN
}

type DebugEventInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Each of the following fields correspond to different kinds of events. as
	// a result, only one is set during any single DebugEventInfo.
	// A singleton event. See enum definition.
	SingletonEvent *SyncEnums_SingletonDebugEventType `protobuf:"varint,1,opt,name=singleton_event,json=singletonEvent,enum=sync_pb.SyncEnums_SingletonDebugEventType" json:"singleton_event,omitempty"`
	// A sync cycle completed.
	SyncCycleCompletedEventInfo *SyncCycleCompletedEventInfo `protobuf:"bytes,2,opt,name=sync_cycle_completed_event_info,json=syncCycleCompletedEventInfo" json:"sync_cycle_completed_event_info,omitempty"`
	// A datatype triggered a nudge.
	NudgingDatatype *int32 `protobuf:"varint,3,opt,name=nudging_datatype,json=nudgingDatatype" json:"nudging_datatype,omitempty"`
	// A notification triggered a nudge.
	DatatypesNotifiedFromServer []int32 `protobuf:"varint,4,rep,name=datatypes_notified_from_server,json=datatypesNotifiedFromServer" json:"datatypes_notified_from_server,omitempty"`
}

func (x *DebugEventInfo) Reset() {
	*x = DebugEventInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_debug_info_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugEventInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugEventInfo) ProtoMessage() {}

func (x *DebugEventInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_debug_info_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugEventInfo.ProtoReflect.Descriptor instead.
func (*DebugEventInfo) Descriptor() ([]byte, []int) {
	return file_client_debug_info_proto_rawDescGZIP(), []int{2}
}

func (x *DebugEventInfo) GetSingletonEvent() SyncEnums_SingletonDebugEventType {
	if x != nil && x.SingletonEvent != nil {
		return *x.SingletonEvent
	}
	return SyncEnums_CONNECTION_STATUS_CHANGE
}

func (x *DebugEventInfo) GetSyncCycleCompletedEventInfo() *SyncCycleCompletedEventInfo {
	if x != nil {
		return x.SyncCycleCompletedEventInfo
	}
	return nil
}

func (x *DebugEventInfo) GetNudgingDatatype() int32 {
	if x != nil && x.NudgingDatatype != nil {
		return *x.NudgingDatatype
	}
	return 0
}

func (x *DebugEventInfo) GetDatatypesNotifiedFromServer() []int32 {
	if x != nil {
		return x.DatatypesNotifiedFromServer
	}
	return nil
}

type DebugInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Events []*DebugEventInfo `protobuf:"bytes,1,rep,name=events" json:"events,omitempty"`
	// Whether cryptographer is ready to encrypt and decrypt data.
	CryptographerReady *bool `protobuf:"varint,2,opt,name=cryptographer_ready,json=cryptographerReady" json:"cryptographer_ready,omitempty"`
	// Cryptographer has pending keys which indicates the correct passphrase
	// has not been provided yet.
	CryptographerHasPendingKeys *bool `protobuf:"varint,3,opt,name=cryptographer_has_pending_keys,json=cryptographerHasPendingKeys" json:"cryptographer_has_pending_keys,omitempty"`
	// Indicates client has dropped some events to save bandwidth.
	EventsDropped *bool `protobuf:"varint,4,opt,name=events_dropped,json=eventsDropped" json:"events_dropped,omitempty"`
}

func (x *DebugInfo) Reset() {
	*x = DebugInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_client_debug_info_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DebugInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DebugInfo) ProtoMessage() {}

func (x *DebugInfo) ProtoReflect() protoreflect.Message {
	mi := &file_client_debug_info_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DebugInfo.ProtoReflect.Descriptor instead.
func (*DebugInfo) Descriptor() ([]byte, []int) {
	return file_client_debug_info_proto_rawDescGZIP(), []int{3}
}

func (x *DebugInfo) GetEvents() []*DebugEventInfo {
	if x != nil {
		return x.Events
	}
	return nil
}

func (x *DebugInfo) GetCryptographerReady() bool {
	if x != nil && x.CryptographerReady != nil {
		return *x.CryptographerReady
	}
	return false
}

func (x *DebugInfo) GetCryptographerHasPendingKeys() bool {
	if x != nil && x.CryptographerHasPendingKeys != nil {
		return *x.CryptographerHasPendingKeys
	}
	return false
}

func (x *DebugInfo) GetEventsDropped() bool {
	if x != nil && x.EventsDropped != nil {
		return *x.EventsDropped
	}
	return false
}

var File_client_debug_info_proto protoreflect.FileDescriptor

var file_client_debug_info_proto_rawDesc = []byte{
	0x0a, 0x17, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x65, 0x62, 0x75, 0x67, 0x5f, 0x69,
	0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x79, 0x6e, 0x63, 0x5f,
	0x70, 0x62, 0x1a, 0x1d, 0x67, 0x65, 0x74, 0x5f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f,
	0x63, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x10, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x65, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x52, 0x0a, 0x08, 0x54, 0x79, 0x70, 0x65, 0x48, 0x69, 0x6e, 0x74, 0x12,
	0x20, 0x0a, 0x0c, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x64, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x49,
	0x64, 0x12, 0x24, 0x0a, 0x0e, 0x68, 0x61, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x5f, 0x68,
	0x69, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x68, 0x61, 0x73, 0x56, 0x61,
	0x6c, 0x69, 0x64, 0x48, 0x69, 0x6e, 0x74, 0x22, 0xf8, 0x04, 0x0a, 0x1b, 0x53, 0x79, 0x6e, 0x63,
	0x43, 0x79, 0x63, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x3c, 0x0a, 0x18, 0x6e, 0x75, 0x6d, 0x5f, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69,
	0x63, 0x74, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x16, 0x6e,
	0x75, 0x6d, 0x45, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x6e, 0x66,
	0x6c, 0x69, 0x63, 0x74, 0x73, 0x12, 0x3a, 0x0a, 0x17, 0x6e, 0x75, 0x6d, 0x5f, 0x68, 0x69, 0x65,
	0x72, 0x61, 0x72, 0x63, 0x68, 0x79, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x73,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x15, 0x6e, 0x75, 0x6d, 0x48,
	0x69, 0x65, 0x72, 0x61, 0x72, 0x63, 0x68, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74,
	0x73, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x12, 0x6e, 0x75, 0x6d, 0x53, 0x69, 0x6d, 0x70, 0x6c, 0x65, 0x43, 0x6f, 0x6e, 0x66, 0x6c, 0x69,
	0x63, 0x74, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x6e, 0x75, 0x6d, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x5f, 0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x12, 0x6e, 0x75, 0x6d, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66,
	0x6c, 0x69, 0x63, 0x74, 0x73, 0x12, 0x34, 0x0a, 0x16, 0x6e, 0x75, 0x6d, 0x5f, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x73, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x05, 0x52, 0x14, 0x6e, 0x75, 0x6d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x73, 0x44, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x12, 0x4b, 0x0a, 0x20, 0x6e,
	0x75, 0x6d, 0x5f, 0x72, 0x65, 0x66, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x73, 0x5f, 0x64, 0x6f, 0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x18,
	0x09, 0x20, 0x01, 0x28, 0x05, 0x42, 0x02, 0x18, 0x01, 0x52, 0x1d, 0x6e, 0x75, 0x6d, 0x52, 0x65,
	0x66, 0x6c, 0x65, 0x63, 0x74, 0x65, 0x64, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x44, 0x6f,
	0x77, 0x6e, 0x6c, 0x6f, 0x61, 0x64, 0x65, 0x64, 0x12, 0x3e, 0x0a, 0x0b, 0x63, 0x61, 0x6c, 0x6c,
	0x65, 0x72, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x73, 0x43, 0x61, 0x6c, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x63, 0x61,
	0x6c, 0x6c, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x51, 0x0a, 0x12, 0x67, 0x65, 0x74, 0x5f,
	0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x73, 0x5f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x23, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x53,
	0x79, 0x6e, 0x63, 0x45, 0x6e, 0x75, 0x6d, 0x73, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x73, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x52, 0x10, 0x67, 0x65, 0x74, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x73, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x4a, 0x04, 0x08, 0x01, 0x10,
	0x02, 0x4a, 0x04, 0x08, 0x02, 0x10, 0x03, 0x4a, 0x04, 0x08, 0x03, 0x10, 0x04, 0x4a, 0x04, 0x08,
	0x0b, 0x10, 0x0c, 0x52, 0x0c, 0x73, 0x79, 0x6e, 0x63, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x75, 0x63,
	0x6b, 0x52, 0x16, 0x6e, 0x75, 0x6d, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f,
	0x63, 0x6f, 0x6e, 0x66, 0x6c, 0x69, 0x63, 0x74, 0x73, 0x52, 0x1a, 0x6e, 0x75, 0x6d, 0x5f, 0x6e,
	0x6f, 0x6e, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x6f, 0x6e, 0x66,
	0x6c, 0x69, 0x63, 0x74, 0x73, 0x52, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x69, 0x6e,
	0x66, 0x6f, 0x22, 0xe3, 0x02, 0x0a, 0x0e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x45, 0x76, 0x65, 0x6e,
	0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x53, 0x0a, 0x0f, 0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74,
	0x6f, 0x6e, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2a,
	0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x53, 0x79, 0x6e, 0x63, 0x45, 0x6e, 0x75,
	0x6d, 0x73, 0x2e, 0x53, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x44, 0x65, 0x62, 0x75,
	0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0e, 0x73, 0x69, 0x6e, 0x67,
	0x6c, 0x65, 0x74, 0x6f, 0x6e, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x6a, 0x0a, 0x1f, 0x73, 0x79,
	0x6e, 0x63, 0x5f, 0x63, 0x79, 0x63, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74,
	0x65, 0x64, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x53, 0x79,
	0x6e, 0x63, 0x43, 0x79, 0x63, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x1b, 0x73, 0x79, 0x6e, 0x63, 0x43,
	0x79, 0x63, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x64, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x29, 0x0a, 0x10, 0x6e, 0x75, 0x64, 0x67, 0x69, 0x6e,
	0x67, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0f, 0x6e, 0x75, 0x64, 0x67, 0x69, 0x6e, 0x67, 0x44, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x43, 0x0a, 0x1e, 0x64, 0x61, 0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x73, 0x5f, 0x6e,
	0x6f, 0x74, 0x69, 0x66, 0x69, 0x65, 0x64, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x73, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x18, 0x04, 0x20, 0x03, 0x28, 0x05, 0x52, 0x1b, 0x64, 0x61, 0x74, 0x61, 0x74,
	0x79, 0x70, 0x65, 0x73, 0x4e, 0x6f, 0x74, 0x69, 0x66, 0x69, 0x65, 0x64, 0x46, 0x72, 0x6f, 0x6d,
	0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4a, 0x04, 0x08, 0x05, 0x10, 0x06, 0x52, 0x1a, 0x64, 0x61,
	0x74, 0x61, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x61, 0x73, 0x73, 0x6f, 0x63, 0x69, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x73, 0x22, 0xd9, 0x01, 0x0a, 0x09, 0x44, 0x65, 0x62,
	0x75, 0x67, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x2f, 0x0a, 0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62,
	0x2e, 0x44, 0x65, 0x62, 0x75, 0x67, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x06, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x61, 0x64, 0x79, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x70,
	0x68, 0x65, 0x72, 0x52, 0x65, 0x61, 0x64, 0x79, 0x12, 0x43, 0x0a, 0x1e, 0x63, 0x72, 0x79, 0x70,
	0x74, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x65, 0x72, 0x5f, 0x68, 0x61, 0x73, 0x5f, 0x70, 0x65,
	0x6e, 0x64, 0x69, 0x6e, 0x67, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x1b, 0x63, 0x72, 0x79, 0x70, 0x74, 0x6f, 0x67, 0x72, 0x61, 0x70, 0x68, 0x65, 0x72, 0x48,
	0x61, 0x73, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x25, 0x0a,
	0x0e, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x5f, 0x64, 0x72, 0x6f, 0x70, 0x70, 0x65, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0d, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x73, 0x44, 0x72, 0x6f,
	0x70, 0x70, 0x65, 0x64, 0x42, 0x36, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f,
	0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50,
	0x01, 0x5a, 0x09, 0x2e, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62,
}

var (
	file_client_debug_info_proto_rawDescOnce sync.Once
	file_client_debug_info_proto_rawDescData = file_client_debug_info_proto_rawDesc
)

func file_client_debug_info_proto_rawDescGZIP() []byte {
	file_client_debug_info_proto_rawDescOnce.Do(func() {
		file_client_debug_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_client_debug_info_proto_rawDescData)
	})
	return file_client_debug_info_proto_rawDescData
}

var file_client_debug_info_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_client_debug_info_proto_goTypes = []interface{}{
	(*TypeHint)(nil),                       // 0: sync_pb.TypeHint
	(*SyncCycleCompletedEventInfo)(nil),    // 1: sync_pb.SyncCycleCompletedEventInfo
	(*DebugEventInfo)(nil),                 // 2: sync_pb.DebugEventInfo
	(*DebugInfo)(nil),                      // 3: sync_pb.DebugInfo
	(*GetUpdatesCallerInfo)(nil),           // 4: sync_pb.GetUpdatesCallerInfo
	(SyncEnums_GetUpdatesOrigin)(0),        // 5: sync_pb.SyncEnums.GetUpdatesOrigin
	(SyncEnums_SingletonDebugEventType)(0), // 6: sync_pb.SyncEnums.SingletonDebugEventType
}
var file_client_debug_info_proto_depIdxs = []int32{
	4, // 0: sync_pb.SyncCycleCompletedEventInfo.caller_info:type_name -> sync_pb.GetUpdatesCallerInfo
	5, // 1: sync_pb.SyncCycleCompletedEventInfo.get_updates_origin:type_name -> sync_pb.SyncEnums.GetUpdatesOrigin
	6, // 2: sync_pb.DebugEventInfo.singleton_event:type_name -> sync_pb.SyncEnums.SingletonDebugEventType
	1, // 3: sync_pb.DebugEventInfo.sync_cycle_completed_event_info:type_name -> sync_pb.SyncCycleCompletedEventInfo
	2, // 4: sync_pb.DebugInfo.events:type_name -> sync_pb.DebugEventInfo
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_client_debug_info_proto_init() }
func file_client_debug_info_proto_init() {
	if File_client_debug_info_proto != nil {
		return
	}
	file_get_updates_caller_info_proto_init()
	file_sync_enums_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_client_debug_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TypeHint); i {
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
		file_client_debug_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncCycleCompletedEventInfo); i {
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
		file_client_debug_info_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugEventInfo); i {
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
		file_client_debug_info_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DebugInfo); i {
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
			RawDescriptor: file_client_debug_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_client_debug_info_proto_goTypes,
		DependencyIndexes: file_client_debug_info_proto_depIdxs,
		MessageInfos:      file_client_debug_info_proto_msgTypes,
	}.Build()
	File_client_debug_info_proto = out.File
	file_client_debug_info_proto_rawDesc = nil
	file_client_debug_info_proto_goTypes = nil
	file_client_debug_info_proto_depIdxs = nil
}
