// Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
//
// Sync protocol datatype extension for custom search engines.

// If you change or add any fields in this file, update proto_visitors.h and
// potentially proto_enum_conversions.{h, cc}.

// Fields that are not used anymore should be marked [deprecated = true].

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.2
// source: search_engine_specifics.proto

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

type SearchEngineSpecifics_ActiveStatus int32

const (
	// The default state when a SE is auto-added. Unspecified SE are inactive.
	SearchEngineSpecifics_ACTIVE_STATUS_UNSPECIFIED SearchEngineSpecifics_ActiveStatus = 0
	// The SE is active and can be triggered via the omnibox.
	SearchEngineSpecifics_ACTIVE_STATUS_TRUE SearchEngineSpecifics_ActiveStatus = 1
	// The SE has been manually set to inactive by the user.
	SearchEngineSpecifics_ACTIVE_STATUS_FALSE SearchEngineSpecifics_ActiveStatus = 2
)

// Enum value maps for SearchEngineSpecifics_ActiveStatus.
var (
	SearchEngineSpecifics_ActiveStatus_name = map[int32]string{
		0: "ACTIVE_STATUS_UNSPECIFIED",
		1: "ACTIVE_STATUS_TRUE",
		2: "ACTIVE_STATUS_FALSE",
	}
	SearchEngineSpecifics_ActiveStatus_value = map[string]int32{
		"ACTIVE_STATUS_UNSPECIFIED": 0,
		"ACTIVE_STATUS_TRUE":        1,
		"ACTIVE_STATUS_FALSE":       2,
	}
)

func (x SearchEngineSpecifics_ActiveStatus) Enum() *SearchEngineSpecifics_ActiveStatus {
	p := new(SearchEngineSpecifics_ActiveStatus)
	*p = x
	return p
}

func (x SearchEngineSpecifics_ActiveStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (SearchEngineSpecifics_ActiveStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_search_engine_specifics_proto_enumTypes[0].Descriptor()
}

func (SearchEngineSpecifics_ActiveStatus) Type() protoreflect.EnumType {
	return &file_search_engine_specifics_proto_enumTypes[0]
}

func (x SearchEngineSpecifics_ActiveStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics_ActiveStatus) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = SearchEngineSpecifics_ActiveStatus(num)
	return nil
}

// Deprecated: Use SearchEngineSpecifics_ActiveStatus.Descriptor instead.
func (SearchEngineSpecifics_ActiveStatus) EnumDescriptor() ([]byte, []int) {
	return file_search_engine_specifics_proto_rawDescGZIP(), []int{0, 0}
}

// Properties of custom search engine sync objects.
type SearchEngineSpecifics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The description of the search engine.
	ShortName *string `protobuf:"bytes,1,opt,name=short_name,json=shortName" json:"short_name,omitempty"`
	// The search engine keyword for omnibox access.
	Keyword *string `protobuf:"bytes,2,opt,name=keyword" json:"keyword,omitempty"`
	// A URL to the favicon to show in the search engines options page.
	FaviconUrl *string `protobuf:"bytes,3,opt,name=favicon_url,json=faviconUrl" json:"favicon_url,omitempty"`
	// The actual parameterized search engine query URL.
	Url *string `protobuf:"bytes,4,opt,name=url" json:"url,omitempty"`
	// A flag signifying whether it is safe to automatically modify this search
	// engine entry or not.
	SafeForAutoreplace *bool `protobuf:"varint,5,opt,name=safe_for_autoreplace,json=safeForAutoreplace" json:"safe_for_autoreplace,omitempty"`
	// The URL to the OSD file this search engine entry came from.
	OriginatingUrl *string `protobuf:"bytes,6,opt,name=originating_url,json=originatingUrl" json:"originating_url,omitempty"`
	// The date this search engine entry was created. A UTC timestamp with units
	// in microseconds.
	DateCreated *int64 `protobuf:"varint,7,opt,name=date_created,json=dateCreated" json:"date_created,omitempty"`
	// A list of supported input encodings.
	InputEncodings *string `protobuf:"bytes,8,opt,name=input_encodings,json=inputEncodings" json:"input_encodings,omitempty"`
	// Obsolete field. This used to represent whether or not this entry is shown
	// in the list of default search engines.
	//
	// Deprecated: Do not use.
	DeprecatedShowInDefaultList *bool `protobuf:"varint,9,opt,name=deprecated_show_in_default_list,json=deprecatedShowInDefaultList" json:"deprecated_show_in_default_list,omitempty"`
	// The parameterized URL that provides suggestions as the user types.
	SuggestionsUrl *string `protobuf:"bytes,10,opt,name=suggestions_url,json=suggestionsUrl" json:"suggestions_url,omitempty"`
	// The ID associated with the prepopulate data this search engine comes from.
	// Set to zero if it was not prepopulated.
	PrepopulateId *int32 `protobuf:"varint,11,opt,name=prepopulate_id,json=prepopulateId" json:"prepopulate_id,omitempty"`
	// DEPRECATED: Whether to autogenerate a keyword for the search engine or not.
	// Do not write to this field in the future.  We preserve this for now so we
	// can read the field in order to migrate existing data that sets this bit.
	//
	// Deprecated: Do not use.
	AutogenerateKeyword *bool `protobuf:"varint,12,opt,name=autogenerate_keyword,json=autogenerateKeyword" json:"autogenerate_keyword,omitempty"`
	// ID 13 reserved - previously used by |logo_id|, now deprecated.
	// Obsolete field. This used to represent whether or not this search engine
	// entry was created automatically by an administrator via group policy. This
	// notion no longer exists amongst synced search engines as we do not want to
	// sync managed search engines.
	// optional bool deprecated_created_by_policy = 14;
	//
	// Deprecated: Do not use.
	InstantUrl *string `protobuf:"bytes,15,opt,name=instant_url,json=instantUrl" json:"instant_url,omitempty"`
	// ID 16 reserved - previously used by |id|, now deprecated.
	// The last time this entry was modified by a user. A UTC timestamp with units
	// in microseconds.
	LastModified *int64 `protobuf:"varint,17,opt,name=last_modified,json=lastModified" json:"last_modified,omitempty"`
	// The primary identifier of this search engine entry for Sync.
	SyncGuid *string `protobuf:"bytes,18,opt,name=sync_guid,json=syncGuid" json:"sync_guid,omitempty"`
	// A list of URL patterns that can be used, in addition to |url|, to extract
	// search terms from a URL.
	AlternateUrls []string `protobuf:"bytes,19,rep,name=alternate_urls,json=alternateUrls" json:"alternate_urls,omitempty"`
	// Deprecated: Do not use.
	SearchTermsReplacementKey *string `protobuf:"bytes,20,opt,name=search_terms_replacement_key,json=searchTermsReplacementKey" json:"search_terms_replacement_key,omitempty"`
	// The parameterized URL that provides image results according to the image
	// content or image URL provided by user.
	ImageUrl *string `protobuf:"bytes,21,opt,name=image_url,json=imageUrl" json:"image_url,omitempty"`
	// The following post_params are comma-separated lists used to specify the
	// post parameters for the corresponding search URL.
	SearchUrlPostParams      *string `protobuf:"bytes,22,opt,name=search_url_post_params,json=searchUrlPostParams" json:"search_url_post_params,omitempty"`
	SuggestionsUrlPostParams *string `protobuf:"bytes,23,opt,name=suggestions_url_post_params,json=suggestionsUrlPostParams" json:"suggestions_url_post_params,omitempty"`
	// Deprecated: Do not use.
	InstantUrlPostParams *string `protobuf:"bytes,24,opt,name=instant_url_post_params,json=instantUrlPostParams" json:"instant_url_post_params,omitempty"`
	ImageUrlPostParams   *string `protobuf:"bytes,25,opt,name=image_url_post_params,json=imageUrlPostParams" json:"image_url_post_params,omitempty"`
	// The parameterized URL for a search provider specified new tab page.
	NewTabUrl *string `protobuf:"bytes,26,opt,name=new_tab_url,json=newTabUrl" json:"new_tab_url,omitempty"`
	// Whether a search engine is 'active' and can be triggered via the omnibox by
	// typing in the relevant keyword.
	IsActive *SearchEngineSpecifics_ActiveStatus `protobuf:"varint,27,opt,name=is_active,json=isActive,enum=sync_pb.SearchEngineSpecifics_ActiveStatus" json:"is_active,omitempty"`
}

func (x *SearchEngineSpecifics) Reset() {
	*x = SearchEngineSpecifics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_search_engine_specifics_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchEngineSpecifics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchEngineSpecifics) ProtoMessage() {}

func (x *SearchEngineSpecifics) ProtoReflect() protoreflect.Message {
	mi := &file_search_engine_specifics_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchEngineSpecifics.ProtoReflect.Descriptor instead.
func (*SearchEngineSpecifics) Descriptor() ([]byte, []int) {
	return file_search_engine_specifics_proto_rawDescGZIP(), []int{0}
}

func (x *SearchEngineSpecifics) GetShortName() string {
	if x != nil && x.ShortName != nil {
		return *x.ShortName
	}
	return ""
}

func (x *SearchEngineSpecifics) GetKeyword() string {
	if x != nil && x.Keyword != nil {
		return *x.Keyword
	}
	return ""
}

func (x *SearchEngineSpecifics) GetFaviconUrl() string {
	if x != nil && x.FaviconUrl != nil {
		return *x.FaviconUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetUrl() string {
	if x != nil && x.Url != nil {
		return *x.Url
	}
	return ""
}

func (x *SearchEngineSpecifics) GetSafeForAutoreplace() bool {
	if x != nil && x.SafeForAutoreplace != nil {
		return *x.SafeForAutoreplace
	}
	return false
}

func (x *SearchEngineSpecifics) GetOriginatingUrl() string {
	if x != nil && x.OriginatingUrl != nil {
		return *x.OriginatingUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetDateCreated() int64 {
	if x != nil && x.DateCreated != nil {
		return *x.DateCreated
	}
	return 0
}

func (x *SearchEngineSpecifics) GetInputEncodings() string {
	if x != nil && x.InputEncodings != nil {
		return *x.InputEncodings
	}
	return ""
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics) GetDeprecatedShowInDefaultList() bool {
	if x != nil && x.DeprecatedShowInDefaultList != nil {
		return *x.DeprecatedShowInDefaultList
	}
	return false
}

func (x *SearchEngineSpecifics) GetSuggestionsUrl() string {
	if x != nil && x.SuggestionsUrl != nil {
		return *x.SuggestionsUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetPrepopulateId() int32 {
	if x != nil && x.PrepopulateId != nil {
		return *x.PrepopulateId
	}
	return 0
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics) GetAutogenerateKeyword() bool {
	if x != nil && x.AutogenerateKeyword != nil {
		return *x.AutogenerateKeyword
	}
	return false
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics) GetInstantUrl() string {
	if x != nil && x.InstantUrl != nil {
		return *x.InstantUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetLastModified() int64 {
	if x != nil && x.LastModified != nil {
		return *x.LastModified
	}
	return 0
}

func (x *SearchEngineSpecifics) GetSyncGuid() string {
	if x != nil && x.SyncGuid != nil {
		return *x.SyncGuid
	}
	return ""
}

func (x *SearchEngineSpecifics) GetAlternateUrls() []string {
	if x != nil {
		return x.AlternateUrls
	}
	return nil
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics) GetSearchTermsReplacementKey() string {
	if x != nil && x.SearchTermsReplacementKey != nil {
		return *x.SearchTermsReplacementKey
	}
	return ""
}

func (x *SearchEngineSpecifics) GetImageUrl() string {
	if x != nil && x.ImageUrl != nil {
		return *x.ImageUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetSearchUrlPostParams() string {
	if x != nil && x.SearchUrlPostParams != nil {
		return *x.SearchUrlPostParams
	}
	return ""
}

func (x *SearchEngineSpecifics) GetSuggestionsUrlPostParams() string {
	if x != nil && x.SuggestionsUrlPostParams != nil {
		return *x.SuggestionsUrlPostParams
	}
	return ""
}

// Deprecated: Do not use.
func (x *SearchEngineSpecifics) GetInstantUrlPostParams() string {
	if x != nil && x.InstantUrlPostParams != nil {
		return *x.InstantUrlPostParams
	}
	return ""
}

func (x *SearchEngineSpecifics) GetImageUrlPostParams() string {
	if x != nil && x.ImageUrlPostParams != nil {
		return *x.ImageUrlPostParams
	}
	return ""
}

func (x *SearchEngineSpecifics) GetNewTabUrl() string {
	if x != nil && x.NewTabUrl != nil {
		return *x.NewTabUrl
	}
	return ""
}

func (x *SearchEngineSpecifics) GetIsActive() SearchEngineSpecifics_ActiveStatus {
	if x != nil && x.IsActive != nil {
		return *x.IsActive
	}
	return SearchEngineSpecifics_ACTIVE_STATUS_UNSPECIFIED
}

var File_search_engine_specifics_proto protoreflect.FileDescriptor

var file_search_engine_specifics_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x65, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x5f,
	0x73, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x70, 0x62, 0x22, 0x97, 0x09, 0x0a, 0x15, 0x53, 0x65, 0x61,
	0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69, 0x6e, 0x65, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69,
	0x63, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x68, 0x6f, 0x72, 0x74, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x18, 0x0a, 0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x66,
	0x61, 0x76, 0x69, 0x63, 0x6f, 0x6e, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x66, 0x61, 0x76, 0x69, 0x63, 0x6f, 0x6e, 0x55, 0x72, 0x6c, 0x12, 0x10, 0x0a, 0x03,
	0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x30,
	0x0a, 0x14, 0x73, 0x61, 0x66, 0x65, 0x5f, 0x66, 0x6f, 0x72, 0x5f, 0x61, 0x75, 0x74, 0x6f, 0x72,
	0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x73, 0x61,
	0x66, 0x65, 0x46, 0x6f, 0x72, 0x41, 0x75, 0x74, 0x6f, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65,
	0x12, 0x27, 0x0a, 0x0f, 0x6f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x5f,
	0x75, 0x72, 0x6c, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x6f, 0x72, 0x69, 0x67, 0x69,
	0x6e, 0x61, 0x74, 0x69, 0x6e, 0x67, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x61, 0x74,
	0x65, 0x5f, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0b, 0x64, 0x61, 0x74, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x64, 0x12, 0x27, 0x0a, 0x0f,
	0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x65, 0x6e, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x73, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x45, 0x6e, 0x63, 0x6f,
	0x64, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x48, 0x0a, 0x1f, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x73, 0x68, 0x6f, 0x77, 0x5f, 0x69, 0x6e, 0x5f, 0x64, 0x65, 0x66, 0x61,
	0x75, 0x6c, 0x74, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x18, 0x09, 0x20, 0x01, 0x28, 0x08, 0x42, 0x02,
	0x18, 0x01, 0x52, 0x1b, 0x64, 0x65, 0x70, 0x72, 0x65, 0x63, 0x61, 0x74, 0x65, 0x64, 0x53, 0x68,
	0x6f, 0x77, 0x49, 0x6e, 0x44, 0x65, 0x66, 0x61, 0x75, 0x6c, 0x74, 0x4c, 0x69, 0x73, 0x74, 0x12,
	0x27, 0x0a, 0x0f, 0x73, 0x75, 0x67, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x75,
	0x72, 0x6c, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x75, 0x67, 0x67, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x55, 0x72, 0x6c, 0x12, 0x25, 0x0a, 0x0e, 0x70, 0x72, 0x65, 0x70,
	0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x0d, 0x70, 0x72, 0x65, 0x70, 0x6f, 0x70, 0x75, 0x6c, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12,
	0x35, 0x0a, 0x14, 0x61, 0x75, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x5f,
	0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x08, 0x42, 0x02, 0x18,
	0x01, 0x52, 0x13, 0x61, 0x75, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x4b,
	0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x23, 0x0a, 0x0b, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e,
	0x74, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52,
	0x0a, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x55, 0x72, 0x6c, 0x12, 0x23, 0x0a, 0x0d, 0x6c,
	0x61, 0x73, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x11, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4d, 0x6f, 0x64, 0x69, 0x66, 0x69, 0x65, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x73, 0x79, 0x6e, 0x63, 0x5f, 0x67, 0x75, 0x69, 0x64, 0x18, 0x12, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x79, 0x6e, 0x63, 0x47, 0x75, 0x69, 0x64, 0x12, 0x25, 0x0a,
	0x0e, 0x61, 0x6c, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x74, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x73, 0x18,
	0x13, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0d, 0x61, 0x6c, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x74, 0x65,
	0x55, 0x72, 0x6c, 0x73, 0x12, 0x43, 0x0a, 0x1c, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x5f, 0x74,
	0x65, 0x72, 0x6d, 0x73, 0x5f, 0x72, 0x65, 0x70, 0x6c, 0x61, 0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74,
	0x5f, 0x6b, 0x65, 0x79, 0x18, 0x14, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52, 0x19,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x54, 0x65, 0x72, 0x6d, 0x73, 0x52, 0x65, 0x70, 0x6c, 0x61,
	0x63, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6d, 0x61,
	0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x15, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d,
	0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x12, 0x33, 0x0a, 0x16, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68,
	0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73,
	0x18, 0x16, 0x20, 0x01, 0x28, 0x09, 0x52, 0x13, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x55, 0x72,
	0x6c, 0x50, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x3d, 0x0a, 0x1b, 0x73,
	0x75, 0x67, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x70,
	0x6f, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x17, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x18, 0x73, 0x75, 0x67, 0x67, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x55, 0x72, 0x6c,
	0x50, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x39, 0x0a, 0x17, 0x69, 0x6e,
	0x73, 0x74, 0x61, 0x6e, 0x74, 0x5f, 0x75, 0x72, 0x6c, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x18, 0x20, 0x01, 0x28, 0x09, 0x42, 0x02, 0x18, 0x01, 0x52,
	0x14, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6e, 0x74, 0x55, 0x72, 0x6c, 0x50, 0x6f, 0x73, 0x74, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x31, 0x0a, 0x15, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75,
	0x72, 0x6c, 0x5f, 0x70, 0x6f, 0x73, 0x74, 0x5f, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x18, 0x19,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x50, 0x6f,
	0x73, 0x74, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x1e, 0x0a, 0x0b, 0x6e, 0x65, 0x77, 0x5f,
	0x74, 0x61, 0x62, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6e,
	0x65, 0x77, 0x54, 0x61, 0x62, 0x55, 0x72, 0x6c, 0x12, 0x48, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61,
	0x63, 0x74, 0x69, 0x76, 0x65, 0x18, 0x1b, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x2b, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x5f, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x45, 0x6e, 0x67, 0x69,
	0x6e, 0x65, 0x53, 0x70, 0x65, 0x63, 0x69, 0x66, 0x69, 0x63, 0x73, 0x2e, 0x41, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x22, 0x5e, 0x0a, 0x0c, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x1d, 0x0a, 0x19, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x5f, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x54, 0x52, 0x55, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x41, 0x43, 0x54,
	0x49, 0x56, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x41, 0x4c, 0x53, 0x45,
	0x10, 0x02, 0x42, 0x2b, 0x0a, 0x25, 0x6f, 0x72, 0x67, 0x2e, 0x63, 0x68, 0x72, 0x6f, 0x6d, 0x69,
	0x75, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x2e, 0x73, 0x79,
	0x6e, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x48, 0x03, 0x50, 0x01,
}

var (
	file_search_engine_specifics_proto_rawDescOnce sync.Once
	file_search_engine_specifics_proto_rawDescData = file_search_engine_specifics_proto_rawDesc
)

func file_search_engine_specifics_proto_rawDescGZIP() []byte {
	file_search_engine_specifics_proto_rawDescOnce.Do(func() {
		file_search_engine_specifics_proto_rawDescData = protoimpl.X.CompressGZIP(file_search_engine_specifics_proto_rawDescData)
	})
	return file_search_engine_specifics_proto_rawDescData
}

var file_search_engine_specifics_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_search_engine_specifics_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_search_engine_specifics_proto_goTypes = []interface{}{
	(SearchEngineSpecifics_ActiveStatus)(0), // 0: sync_pb.SearchEngineSpecifics.ActiveStatus
	(*SearchEngineSpecifics)(nil),           // 1: sync_pb.SearchEngineSpecifics
}
var file_search_engine_specifics_proto_depIdxs = []int32{
	0, // 0: sync_pb.SearchEngineSpecifics.is_active:type_name -> sync_pb.SearchEngineSpecifics.ActiveStatus
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_search_engine_specifics_proto_init() }
func file_search_engine_specifics_proto_init() {
	if File_search_engine_specifics_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_search_engine_specifics_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchEngineSpecifics); i {
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
			RawDescriptor: file_search_engine_specifics_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_search_engine_specifics_proto_goTypes,
		DependencyIndexes: file_search_engine_specifics_proto_depIdxs,
		EnumInfos:         file_search_engine_specifics_proto_enumTypes,
		MessageInfos:      file_search_engine_specifics_proto_msgTypes,
	}.Build()
	File_search_engine_specifics_proto = out.File
	file_search_engine_specifics_proto_rawDesc = nil
	file_search_engine_specifics_proto_goTypes = nil
	file_search_engine_specifics_proto_depIdxs = nil
}
