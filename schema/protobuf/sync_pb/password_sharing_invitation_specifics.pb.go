// Copyright 2023 The Chromium Authors
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.1
// source: password_sharing_invitation_specifics.proto

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

type PasswordSharingInvitationData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PasswordData *PasswordSharingInvitationData_PasswordData `protobuf:"bytes,1,opt,name=password_data,json=passwordData" json:"password_data,omitempty"`
}

func (x *PasswordSharingInvitationData) Reset() {
	*x = PasswordSharingInvitationData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasswordSharingInvitationData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordSharingInvitationData) ProtoMessage() {}

func (x *PasswordSharingInvitationData) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordSharingInvitationData.ProtoReflect.Descriptor instead.
func (*PasswordSharingInvitationData) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *PasswordSharingInvitationData) GetPasswordData() *PasswordSharingInvitationData_PasswordData {
	if x != nil {
		return x.PasswordData
	}
	return nil
}

// Contains user profile information.
type UserDisplayInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Primary email address of the user.
	Email *string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	// The user's full name.
	DisplayName *string `protobuf:"bytes,2,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
	// Portrait photo of the user.
	ProfileImageUrl *string `protobuf:"bytes,3,opt,name=profile_image_url,json=profileImageUrl" json:"profile_image_url,omitempty"`
}

func (x *UserDisplayInfo) Reset() {
	*x = UserDisplayInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDisplayInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDisplayInfo) ProtoMessage() {}

func (x *UserDisplayInfo) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDisplayInfo.ProtoReflect.Descriptor instead.
func (*UserDisplayInfo) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{1}
}

func (x *UserDisplayInfo) GetEmail() string {
	if x != nil && x.Email != nil {
		return *x.Email
	}
	return ""
}

func (x *UserDisplayInfo) GetDisplayName() string {
	if x != nil && x.DisplayName != nil {
		return *x.DisplayName
	}
	return ""
}

func (x *UserDisplayInfo) GetProfileImageUrl() string {
	if x != nil && x.ProfileImageUrl != nil {
		return *x.ProfileImageUrl
	}
	return ""
}

type UserInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Obfuscated Gaia ID.
	UserId          *string          `protobuf:"bytes,1,opt,name=user_id,json=userId" json:"user_id,omitempty"`
	UserDisplayInfo *UserDisplayInfo `protobuf:"bytes,2,opt,name=user_display_info,json=userDisplayInfo" json:"user_display_info,omitempty"`
	// Latest user's public key registered on the server.
	CrossUserSharingPublicKey *CrossUserSharingPublicKey `protobuf:"bytes,3,opt,name=cross_user_sharing_public_key,json=crossUserSharingPublicKey" json:"cross_user_sharing_public_key,omitempty"`
}

func (x *UserInfo) Reset() {
	*x = UserInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserInfo) ProtoMessage() {}

func (x *UserInfo) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserInfo.ProtoReflect.Descriptor instead.
func (*UserInfo) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{2}
}

func (x *UserInfo) GetUserId() string {
	if x != nil && x.UserId != nil {
		return *x.UserId
	}
	return ""
}

func (x *UserInfo) GetUserDisplayInfo() *UserDisplayInfo {
	if x != nil {
		return x.UserDisplayInfo
	}
	return nil
}

func (x *UserInfo) GetCrossUserSharingPublicKey() *CrossUserSharingPublicKey {
	if x != nil {
		return x.CrossUserSharingPublicKey
	}
	return nil
}

// Encryption key used to encrypt PasswordSharingInvitationData.
type SharingSymmetricKey struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyValue []byte `protobuf:"bytes,1,opt,name=key_value,json=keyValue" json:"key_value,omitempty"`
}

func (x *SharingSymmetricKey) Reset() {
	*x = SharingSymmetricKey{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SharingSymmetricKey) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SharingSymmetricKey) ProtoMessage() {}

func (x *SharingSymmetricKey) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SharingSymmetricKey.ProtoReflect.Descriptor instead.
func (*SharingSymmetricKey) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{3}
}

func (x *SharingSymmetricKey) GetKeyValue() []byte {
	if x != nil {
		return x.KeyValue
	}
	return nil
}

// Incoming invitations for password sending.
type IncomingPasswordSharingInvitationSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique client tag for the invitation. This does *not* have to be the same
	// GUID as for the outgoing invitation.
	Guid *string `protobuf:"bytes,1,opt,name=guid" json:"guid,omitempty"`
	// Profile information about the sender of the password. Sender's public key
	// is used to authenticate the sender for `encrypted_key_for_recipient`.
	SenderInfo *UserInfo `protobuf:"bytes,2,opt,name=sender_info,json=senderInfo" json:"sender_info,omitempty"`
	// PasswordSharingInvitationData, encrypted using the encryption key from
	// `encrypted_key_for_recipient`.
	EncryptedPasswordSharingInvitationData []byte `protobuf:"bytes,3,opt,name=encrypted_password_sharing_invitation_data,json=encryptedPasswordSharingInvitationData" json:"encrypted_password_sharing_invitation_data,omitempty"`
	// An unsynced field for use internally on the client. This field should
	// never be set in any network-based communications because it contains
	// unencrypted material.
	ClientOnlyUnencryptedData *PasswordSharingInvitationData `protobuf:"bytes,4,opt,name=client_only_unencrypted_data,json=clientOnlyUnencryptedData" json:"client_only_unencrypted_data,omitempty"`
	// Encrypted SharingSymmetricKey using recipient's public key corresponding to
	// `recipient_key_version` and sender's private key to authenticate the
	// sender, see https://www.rfc-editor.org/rfc/rfc9180.html.
	EncryptedKeyForRecipient []byte  `protobuf:"bytes,5,opt,name=encrypted_key_for_recipient,json=encryptedKeyForRecipient" json:"encrypted_key_for_recipient,omitempty"`
	RecipientKeyVersion      *uint32 `protobuf:"varint,6,opt,name=recipient_key_version,json=recipientKeyVersion" json:"recipient_key_version,omitempty"`
}

func (x *IncomingPasswordSharingInvitationSpecifics) Reset() {
	*x = IncomingPasswordSharingInvitationSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *IncomingPasswordSharingInvitationSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*IncomingPasswordSharingInvitationSpecifics) ProtoMessage() {}

func (x *IncomingPasswordSharingInvitationSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use IncomingPasswordSharingInvitationSpecifics.ProtoReflect.Descriptor instead.
func (*IncomingPasswordSharingInvitationSpecifics) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{4}
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetGuid() string {
	if x != nil && x.Guid != nil {
		return *x.Guid
	}
	return ""
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetSenderInfo() *UserInfo {
	if x != nil {
		return x.SenderInfo
	}
	return nil
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetEncryptedPasswordSharingInvitationData() []byte {
	if x != nil {
		return x.EncryptedPasswordSharingInvitationData
	}
	return nil
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetClientOnlyUnencryptedData() *PasswordSharingInvitationData {
	if x != nil {
		return x.ClientOnlyUnencryptedData
	}
	return nil
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetEncryptedKeyForRecipient() []byte {
	if x != nil {
		return x.EncryptedKeyForRecipient
	}
	return nil
}

func (x *IncomingPasswordSharingInvitationSpecifics) GetRecipientKeyVersion() uint32 {
	if x != nil && x.RecipientKeyVersion != nil {
		return *x.RecipientKeyVersion
	}
	return 0
}

// Outgoing invitations for password sending.
type OutgoingPasswordSharingInvitationSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique client tag for the invitation, generated by the client.
	Guid *string `protobuf:"bytes,1,opt,name=guid" json:"guid,omitempty"`
	// Recipient's user identifier (obfuscated Gaia ID).
	RecipientUserId *string `protobuf:"bytes,2,opt,name=recipient_user_id,json=recipientUserId" json:"recipient_user_id,omitempty"`
	// The actual data, contains an encrypted PasswordSharingInvitationData using
	// an encryption key from `encrypted_key_for_recipient`.
	EncryptedPasswordSharingInvitationData []byte `protobuf:"bytes,3,opt,name=encrypted_password_sharing_invitation_data,json=encryptedPasswordSharingInvitationData" json:"encrypted_password_sharing_invitation_data,omitempty"`
	// An unsynced field for use internally on the client. This field should
	// never be set in any network-based communications because it contains
	// unencrypted material.
	ClientOnlyUnencryptedData *PasswordSharingInvitationData `protobuf:"bytes,4,opt,name=client_only_unencrypted_data,json=clientOnlyUnencryptedData" json:"client_only_unencrypted_data,omitempty"`
	// Encrypted SharingSymmetricKey using recipient's public key corresponding to
	// `recipient_key_version`.
	EncryptedKeyForRecipient []byte  `protobuf:"bytes,5,opt,name=encrypted_key_for_recipient,json=encryptedKeyForRecipient" json:"encrypted_key_for_recipient,omitempty"`
	RecipientKeyVersion      *uint32 `protobuf:"varint,6,opt,name=recipient_key_version,json=recipientKeyVersion" json:"recipient_key_version,omitempty"`
	// Version of Public key of the sender which is used to authenticate the
	// sender of the password. Must be equal to the latest committed version.
	SenderKeyVersion *uint32 `protobuf:"varint,7,opt,name=sender_key_version,json=senderKeyVersion" json:"sender_key_version,omitempty"`
}

func (x *OutgoingPasswordSharingInvitationSpecifics) Reset() {
	*x = OutgoingPasswordSharingInvitationSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutgoingPasswordSharingInvitationSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutgoingPasswordSharingInvitationSpecifics) ProtoMessage() {}

func (x *OutgoingPasswordSharingInvitationSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutgoingPasswordSharingInvitationSpecifics.ProtoReflect.Descriptor instead.
func (*OutgoingPasswordSharingInvitationSpecifics) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{5}
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetGuid() string {
	if x != nil && x.Guid != nil {
		return *x.Guid
	}
	return ""
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetRecipientUserId() string {
	if x != nil && x.RecipientUserId != nil {
		return *x.RecipientUserId
	}
	return ""
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetEncryptedPasswordSharingInvitationData() []byte {
	if x != nil {
		return x.EncryptedPasswordSharingInvitationData
	}
	return nil
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetClientOnlyUnencryptedData() *PasswordSharingInvitationData {
	if x != nil {
		return x.ClientOnlyUnencryptedData
	}
	return nil
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetEncryptedKeyForRecipient() []byte {
	if x != nil {
		return x.EncryptedKeyForRecipient
	}
	return nil
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetRecipientKeyVersion() uint32 {
	if x != nil && x.RecipientKeyVersion != nil {
		return *x.RecipientKeyVersion
	}
	return 0
}

func (x *OutgoingPasswordSharingInvitationSpecifics) GetSenderKeyVersion() uint32 {
	if x != nil && x.SenderKeyVersion != nil {
		return *x.SenderKeyVersion
	}
	return 0
}

// Contains password fields required for sending. See PasswordSpecificsData
// for field descriptions.
type PasswordSharingInvitationData_PasswordData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PasswordValue *string `protobuf:"bytes,1,opt,name=password_value,json=passwordValue" json:"password_value,omitempty"`
	// See PasswordSpecificsData::Scheme for values.
	Scheme          *int32  `protobuf:"varint,2,opt,name=scheme" json:"scheme,omitempty"`
	SignonRealm     *string `protobuf:"bytes,3,opt,name=signon_realm,json=signonRealm" json:"signon_realm,omitempty"`
	Origin          *string `protobuf:"bytes,4,opt,name=origin" json:"origin,omitempty"`
	UsernameElement *string `protobuf:"bytes,5,opt,name=username_element,json=usernameElement" json:"username_element,omitempty"`
	UsernameValue   *string `protobuf:"bytes,6,opt,name=username_value,json=usernameValue" json:"username_value,omitempty"`
	PasswordElement *string `protobuf:"bytes,7,opt,name=password_element,json=passwordElement" json:"password_element,omitempty"`
	DisplayName     *string `protobuf:"bytes,8,opt,name=display_name,json=displayName" json:"display_name,omitempty"`
	AvatarUrl       *string `protobuf:"bytes,9,opt,name=avatar_url,json=avatarUrl" json:"avatar_url,omitempty"`
}

func (x *PasswordSharingInvitationData_PasswordData) Reset() {
	*x = PasswordSharingInvitationData_PasswordData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_password_sharing_invitation_specifics_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PasswordSharingInvitationData_PasswordData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PasswordSharingInvitationData_PasswordData) ProtoMessage() {}

func (x *PasswordSharingInvitationData_PasswordData) ProtoReflect() protoreflect.Message {
	mi := &file_password_sharing_invitation_specifics_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PasswordSharingInvitationData_PasswordData.ProtoReflect.Descriptor instead.
func (*PasswordSharingInvitationData_PasswordData) Descriptor() ([]byte, []int) {
	return file_password_sharing_invitation_specifics_proto_rawDescGZIP(), []int{0, 0}
}

func (x *PasswordSharingInvitationData_PasswordData) GetPasswordValue() string {
	if x != nil && x.PasswordValue != nil {
		return *x.PasswordValue
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetScheme() int32 {
	if x != nil && x.Scheme != nil {
		return *x.Scheme
	}
	return 0
}

func (x *PasswordSharingInvitationData_PasswordData) GetSignonRealm() string {
	if x != nil && x.SignonRealm != nil {
		return *x.SignonRealm
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetOrigin() string {
	if x != nil && x.Origin != nil {
		return *x.Origin
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetUsernameElement() string {
	if x != nil && x.UsernameElement != nil {
		return *x.UsernameElement
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetUsernameValue() string {
	if x != nil && x.UsernameValue != nil {
		return *x.UsernameValue
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetPasswordElement() string {
	if x != nil && x.PasswordElement != nil {
		return *x.PasswordElement
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetDisplayName() string {
	if x != nil && x.DisplayName != nil {
		return *x.DisplayName
	}
	return ""
}

func (x *PasswordSharingInvitationData_PasswordData) GetAvatarUrl() string {
	if x != nil && x.AvatarUrl != nil {
		return *x.AvatarUrl
	}
	return ""
}

var File_password_sharing_invitation_specifics_proto protoreflect.FileDescriptor

var file_password_sharing_invitation_specifics_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x69,
	0x6e, 0x67, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x70,
	0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73,
	0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x1a, 0x16, 0x6e, 0x69, 0x67, 0x6f, 0x72, 0x69, 0x5f, 0x73,
	0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xc3,
	0x03, 0x0a, 0x1d, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69,
	0x6e, 0x67, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61,
	0x12, 0x58, 0x0a, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x33, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70,
	0x62, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e,
	0x67, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x2e,
	0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x52, 0x0c, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x1a, 0xc7, 0x02, 0x0a, 0x0c, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x25, 0x0a, 0x0e, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x06, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x69,
	0x67, 0x6e, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x61, 0x6c, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x73, 0x69, 0x67, 0x6e, 0x6f, 0x6e, 0x52, 0x65, 0x61, 0x6c, 0x6d, 0x12, 0x16, 0x0a,
	0x06, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x6f,
	0x72, 0x69, 0x67, 0x69, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x5f, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x45, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x12, 0x25, 0x0a, 0x0e, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x70, 0x61, 0x73, 0x73, 0x77,
	0x6f, 0x72, 0x64, 0x5f, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x45, 0x6c, 0x65, 0x6d, 0x65,
	0x6e, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x76, 0x61, 0x74, 0x61,
	0x72, 0x55, 0x72, 0x6c, 0x22, 0x76, 0x0a, 0x0f, 0x55, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x21, 0x0a,
	0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x2a, 0x0a, 0x11, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6d, 0x61, 0x67,
	0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x49, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x22, 0xcf, 0x01, 0x0a,
	0x08, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65,
	0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x12, 0x44, 0x0a, 0x11, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x64, 0x69, 0x73, 0x70, 0x6c,
	0x61, 0x79, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x44, 0x69, 0x73,
	0x70, 0x6c, 0x61, 0x79, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x64, 0x0a, 0x1d, 0x63, 0x72, 0x6f, 0x73,
	0x73, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x70,
	0x75, 0x62, 0x6c, 0x69, 0x63, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x6f, 0x73, 0x73, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63,
	0x4b, 0x65, 0x79, 0x52, 0x19, 0x63, 0x72, 0x6f, 0x73, 0x73, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68,
	0x61, 0x72, 0x69, 0x6e, 0x67, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79, 0x22, 0x32,
	0x0a, 0x13, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x53, 0x79, 0x6d, 0x6d, 0x65, 0x74, 0x72,
	0x69, 0x63, 0x4b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6b, 0x65, 0x79, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0xac, 0x03, 0x0a, 0x2a, 0x49, 0x6e, 0x63, 0x6f, 0x6d, 0x69, 0x6e, 0x67, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63,
	0x73, 0x12, 0x12, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x32, 0x0a, 0x0b, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x5f,
	0x69, 0x6e, 0x66, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x73, 0x79, 0x6e,
	0x63, 0x5f, 0x70, 0x62, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x0a, 0x73,
	0x65, 0x6e, 0x64, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x5a, 0x0a, 0x2a, 0x65, 0x6e, 0x63,
	0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f,
	0x73, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x26, 0x65,
	0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64,
	0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x67, 0x0a, 0x1c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x75, 0x6e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x50, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68,
	0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x19, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x4f, 0x6e, 0x6c, 0x79, 0x55,
	0x6e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3d,
	0x0a, 0x1b, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x5f,
	0x66, 0x6f, 0x72, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x0c, 0x52, 0x18, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x4b, 0x65,
	0x79, 0x46, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a,
	0x15, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x13, 0x72, 0x65,
	0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f,
	0x6e, 0x22, 0xd2, 0x03, 0x0a, 0x2a, 0x4f, 0x75, 0x74, 0x67, 0x6f, 0x69, 0x6e, 0x67, 0x50, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x76,
	0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73,
	0x12, 0x12, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x67, 0x75, 0x69, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x5a, 0x0a, 0x2a, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x5f, 0x73, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x69,
	0x6e, 0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x26, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x50,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e,
	0x76, 0x69, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x67, 0x0a, 0x1c,
	0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x6f, 0x6e, 0x6c, 0x79, 0x5f, 0x75, 0x6e, 0x65, 0x6e,
	0x63, 0x72, 0x79, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x26, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x50, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x53, 0x68, 0x61, 0x72, 0x69, 0x6e, 0x67, 0x49, 0x6e, 0x76, 0x69,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x19, 0x63, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x4f, 0x6e, 0x6c, 0x79, 0x55, 0x6e, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74, 0x65,
	0x64, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3d, 0x0a, 0x1b, 0x65, 0x6e, 0x63, 0x72, 0x79, 0x70, 0x74,
	0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x72, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x18, 0x65, 0x6e, 0x63, 0x72,
	0x79, 0x70, 0x74, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x46, 0x6f, 0x72, 0x52, 0x65, 0x63, 0x69, 0x70,
	0x69, 0x65, 0x6e, 0x74, 0x12, 0x32, 0x0a, 0x15, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x13, 0x72, 0x65, 0x63, 0x69, 0x70, 0x69, 0x65, 0x6e, 0x74, 0x4b, 0x65,
	0x79, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x2c, 0x0a, 0x12, 0x73, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x5f, 0x6b, 0x65, 0x79, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x73, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x4b, 0x65, 0x79, 0x56,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x36, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68,
	0x72, 0x6f, 0x6d, 0x69, 0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x73, 0x2e, 0x73, 0x79, 0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48,
	0x03, 0x50, 0x01, 0x5a, 0x09, 0x2e, 0x2f, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62,
}

var (
	file_password_sharing_invitation_specifics_proto_rawDescOnce sync.Once
	file_password_sharing_invitation_specifics_proto_rawDescData = file_password_sharing_invitation_specifics_proto_rawDesc
)

func file_password_sharing_invitation_specifics_proto_rawDescGZIP() []byte {
	file_password_sharing_invitation_specifics_proto_rawDescOnce.Do(func() {
		file_password_sharing_invitation_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_password_sharing_invitation_specifics_proto_rawDescData)
	})
	return file_password_sharing_invitation_specifics_proto_rawDescData
}

var file_password_sharing_invitation_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_password_sharing_invitation_specifics_proto_goTypes = []interface{}{
	(*PasswordSharingInvitationData)(nil),              // 0: sync_pb.PasswordSharingInvitationData
	(*UserDisplayInfo)(nil),                            // 1: sync_pb.UserDisplayInfo
	(*UserInfo)(nil),                                   // 2: sync_pb.UserInfo
	(*SharingSymmetricKey)(nil),                        // 3: sync_pb.SharingSymmetricKey
	(*IncomingPasswordSharingInvitationSpecifics)(nil), // 4: sync_pb.IncomingPasswordSharingInvitationSpecifics
	(*OutgoingPasswordSharingInvitationSpecifics)(nil), // 5: sync_pb.OutgoingPasswordSharingInvitationSpecifics
	(*PasswordSharingInvitationData_PasswordData)(nil), // 6: sync_pb.PasswordSharingInvitationData.PasswordData
	(*CrossUserSharingPublicKey)(nil),                  // 7: sync_pb.CrossUserSharingPublicKey
}
var file_password_sharing_invitation_specifics_proto_depIdxs = []int32{
	6, // 0: sync_pb.PasswordSharingInvitationData.password_data:type_name -> sync_pb.PasswordSharingInvitationData.PasswordData
	1, // 1: sync_pb.UserInfo.user_display_info:type_name -> sync_pb.UserDisplayInfo
	7, // 2: sync_pb.UserInfo.cross_user_sharing_public_key:type_name -> sync_pb.CrossUserSharingPublicKey
	2, // 3: sync_pb.IncomingPasswordSharingInvitationSpecifics.sender_info:type_name -> sync_pb.UserInfo
	0, // 4: sync_pb.IncomingPasswordSharingInvitationSpecifics.client_only_unencrypted_data:type_name -> sync_pb.PasswordSharingInvitationData
	0, // 5: sync_pb.OutgoingPasswordSharingInvitationSpecifics.client_only_unencrypted_data:type_name -> sync_pb.PasswordSharingInvitationData
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_password_sharing_invitation_specifics_proto_init() }
func file_password_sharing_invitation_specifics_proto_init() {
	if File_password_sharing_invitation_specifics_proto != nil {
		return
	}
	file_nigori_specifics_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_password_sharing_invitation_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasswordSharingInvitationData); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDisplayInfo); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserInfo); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SharingSymmetricKey); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*IncomingPasswordSharingInvitationSpecifics); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutgoingPasswordSharingInvitationSpecifics); i {
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
		file_password_sharing_invitation_specifics_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PasswordSharingInvitationData_PasswordData); i {
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
			RawDescriptor: file_password_sharing_invitation_specifics_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_password_sharing_invitation_specifics_proto_goTypes,
		DependencyIndexes: file_password_sharing_invitation_specifics_proto_depIdxs,
		MessageInfos:      file_password_sharing_invitation_specifics_proto_msgTypes,
	}.Build()
	File_password_sharing_invitation_specifics_proto = out.File
	file_password_sharing_invitation_specifics_proto_rawDesc = nil
	file_password_sharing_invitation_specifics_proto_goTypes = nil
	file_password_sharing_invitation_specifics_proto_depIdxs = nil
}
