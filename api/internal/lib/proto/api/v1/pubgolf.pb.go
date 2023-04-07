// PubGolf defines the app-facing API service for the in-game apps.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: api/v1/pubgolf.proto

package apiv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ClientVersionResponse_VersionStatus int32

const (
	ClientVersionResponse_VERSION_STATUS_UNSPECIFIED  ClientVersionResponse_VersionStatus = 0
	ClientVersionResponse_VERSION_STATUS_OK           ClientVersionResponse_VersionStatus = 1
	ClientVersionResponse_VERSION_STATUS_OUTDATED     ClientVersionResponse_VersionStatus = 2
	ClientVersionResponse_VERSION_STATUS_INCOMPATIBLE ClientVersionResponse_VersionStatus = 3
)

// Enum value maps for ClientVersionResponse_VersionStatus.
var (
	ClientVersionResponse_VersionStatus_name = map[int32]string{
		0: "VERSION_STATUS_UNSPECIFIED",
		1: "VERSION_STATUS_OK",
		2: "VERSION_STATUS_OUTDATED",
		3: "VERSION_STATUS_INCOMPATIBLE",
	}
	ClientVersionResponse_VersionStatus_value = map[string]int32{
		"VERSION_STATUS_UNSPECIFIED":  0,
		"VERSION_STATUS_OK":           1,
		"VERSION_STATUS_OUTDATED":     2,
		"VERSION_STATUS_INCOMPATIBLE": 3,
	}
)

func (x ClientVersionResponse_VersionStatus) Enum() *ClientVersionResponse_VersionStatus {
	p := new(ClientVersionResponse_VersionStatus)
	*p = x
	return p
}

func (x ClientVersionResponse_VersionStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ClientVersionResponse_VersionStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_api_v1_pubgolf_proto_enumTypes[0].Descriptor()
}

func (ClientVersionResponse_VersionStatus) Type() protoreflect.EnumType {
	return &file_api_v1_pubgolf_proto_enumTypes[0]
}

func (x ClientVersionResponse_VersionStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ClientVersionResponse_VersionStatus.Descriptor instead.
func (ClientVersionResponse_VersionStatus) EnumDescriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{1, 0}
}

type ClientVersionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientVersion uint32 `protobuf:"varint,1,opt,name=client_version,json=clientVersion,proto3" json:"client_version,omitempty"`
}

func (x *ClientVersionRequest) Reset() {
	*x = ClientVersionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientVersionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientVersionRequest) ProtoMessage() {}

func (x *ClientVersionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientVersionRequest.ProtoReflect.Descriptor instead.
func (*ClientVersionRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{0}
}

func (x *ClientVersionRequest) GetClientVersion() uint32 {
	if x != nil {
		return x.ClientVersion
	}
	return 0
}

type ClientVersionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VersionStatus ClientVersionResponse_VersionStatus `protobuf:"varint,1,opt,name=version_status,json=versionStatus,proto3,enum=api.v1.ClientVersionResponse_VersionStatus" json:"version_status,omitempty"`
}

func (x *ClientVersionResponse) Reset() {
	*x = ClientVersionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientVersionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientVersionResponse) ProtoMessage() {}

func (x *ClientVersionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientVersionResponse.ProtoReflect.Descriptor instead.
func (*ClientVersionResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{1}
}

func (x *ClientVersionResponse) GetVersionStatus() ClientVersionResponse_VersionStatus {
	if x != nil {
		return x.VersionStatus
	}
	return ClientVersionResponse_VERSION_STATUS_UNSPECIFIED
}

type GetScheduleRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventKey          string  `protobuf:"bytes,1,opt,name=event_key,json=eventKey,proto3" json:"event_key,omitempty"`
	CachedDataVersion *uint32 `protobuf:"varint,2,opt,name=cached_data_version,json=cachedDataVersion,proto3,oneof" json:"cached_data_version,omitempty"`
}

func (x *GetScheduleRequest) Reset() {
	*x = GetScheduleRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetScheduleRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleRequest) ProtoMessage() {}

func (x *GetScheduleRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScheduleRequest.ProtoReflect.Descriptor instead.
func (*GetScheduleRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{2}
}

func (x *GetScheduleRequest) GetEventKey() string {
	if x != nil {
		return x.EventKey
	}
	return ""
}

func (x *GetScheduleRequest) GetCachedDataVersion() uint32 {
	if x != nil && x.CachedDataVersion != nil {
		return *x.CachedDataVersion
	}
	return 0
}

type GetScheduleResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LatestDataVersion uint32                        `protobuf:"varint,1,opt,name=latest_data_version,json=latestDataVersion,proto3" json:"latest_data_version,omitempty"`
	Schedule          *GetScheduleResponse_Schedule `protobuf:"bytes,2,opt,name=schedule,proto3,oneof" json:"schedule,omitempty"`
}

func (x *GetScheduleResponse) Reset() {
	*x = GetScheduleResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetScheduleResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleResponse) ProtoMessage() {}

func (x *GetScheduleResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScheduleResponse.ProtoReflect.Descriptor instead.
func (*GetScheduleResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{3}
}

func (x *GetScheduleResponse) GetLatestDataVersion() uint32 {
	if x != nil {
		return x.LatestDataVersion
	}
	return 0
}

func (x *GetScheduleResponse) GetSchedule() *GetScheduleResponse_Schedule {
	if x != nil {
		return x.Schedule
	}
	return nil
}

type GetVenueRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventKey  string   `protobuf:"bytes,1,opt,name=event_key,json=eventKey,proto3" json:"event_key,omitempty"`
	VenueKeys []uint32 `protobuf:"varint,2,rep,packed,name=venue_keys,json=venueKeys,proto3" json:"venue_keys,omitempty"`
}

func (x *GetVenueRequest) Reset() {
	*x = GetVenueRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVenueRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVenueRequest) ProtoMessage() {}

func (x *GetVenueRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVenueRequest.ProtoReflect.Descriptor instead.
func (*GetVenueRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{4}
}

func (x *GetVenueRequest) GetEventKey() string {
	if x != nil {
		return x.EventKey
	}
	return ""
}

func (x *GetVenueRequest) GetVenueKeys() []uint32 {
	if x != nil {
		return x.VenueKeys
	}
	return nil
}

type GetVenueResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Map of requested venue keys to Venue objects.
	Venues map[uint32]*GetVenueResponse_VenueWrapper `protobuf:"bytes,1,rep,name=venues,proto3" json:"venues,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *GetVenueResponse) Reset() {
	*x = GetVenueResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVenueResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVenueResponse) ProtoMessage() {}

func (x *GetVenueResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVenueResponse.ProtoReflect.Descriptor instead.
func (*GetVenueResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{5}
}

func (x *GetVenueResponse) GetVenues() map[uint32]*GetVenueResponse_VenueWrapper {
	if x != nil {
		return x.Venues
	}
	return nil
}

type GetScheduleResponse_Schedule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// List of past venues. Does not include the current venue.
	VisitedVenueKeys []uint32 `protobuf:"varint,1,rep,packed,name=visited_venue_keys,json=visitedVenueKeys,proto3" json:"visited_venue_keys,omitempty"`
	// Optional in the case that the event hasn't started yet.
	CurrentVenueKey *uint32 `protobuf:"varint,2,opt,name=current_venue_key,json=currentVenueKey,proto3,oneof" json:"current_venue_key,omitempty"`
	// Optional in the case that the next venue isn't yet visible to players, or after the second to last venue. The next venue key only becomes visible X mins before the next venue's start time.
	NextVenueKey   *uint32                `protobuf:"varint,3,opt,name=next_venue_key,json=nextVenueKey,proto3,oneof" json:"next_venue_key,omitempty"`
	NextVenueStart *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=next_venue_start,json=nextVenueStart,proto3,oneof" json:"next_venue_start,omitempty"`
	EventEnd       *timestamppb.Timestamp `protobuf:"bytes,5,opt,name=event_end,json=eventEnd,proto3" json:"event_end,omitempty"`
}

func (x *GetScheduleResponse_Schedule) Reset() {
	*x = GetScheduleResponse_Schedule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetScheduleResponse_Schedule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetScheduleResponse_Schedule) ProtoMessage() {}

func (x *GetScheduleResponse_Schedule) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetScheduleResponse_Schedule.ProtoReflect.Descriptor instead.
func (*GetScheduleResponse_Schedule) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{3, 0}
}

func (x *GetScheduleResponse_Schedule) GetVisitedVenueKeys() []uint32 {
	if x != nil {
		return x.VisitedVenueKeys
	}
	return nil
}

func (x *GetScheduleResponse_Schedule) GetCurrentVenueKey() uint32 {
	if x != nil && x.CurrentVenueKey != nil {
		return *x.CurrentVenueKey
	}
	return 0
}

func (x *GetScheduleResponse_Schedule) GetNextVenueKey() uint32 {
	if x != nil && x.NextVenueKey != nil {
		return *x.NextVenueKey
	}
	return 0
}

func (x *GetScheduleResponse_Schedule) GetNextVenueStart() *timestamppb.Timestamp {
	if x != nil {
		return x.NextVenueStart
	}
	return nil
}

func (x *GetScheduleResponse_Schedule) GetEventEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.EventEnd
	}
	return nil
}

type GetVenueResponse_Venue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Global ID for the venue in ULID format (26 characters, base32), not to be confused with the venue key.
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Address string suitable for display or using for a mapping query.
	Address  string `protobuf:"bytes,3,opt,name=address,proto3" json:"address,omitempty"`
	ImageUrl string `protobuf:"bytes,4,opt,name=image_url,json=imageUrl,proto3" json:"image_url,omitempty"`
}

func (x *GetVenueResponse_Venue) Reset() {
	*x = GetVenueResponse_Venue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVenueResponse_Venue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVenueResponse_Venue) ProtoMessage() {}

func (x *GetVenueResponse_Venue) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVenueResponse_Venue.ProtoReflect.Descriptor instead.
func (*GetVenueResponse_Venue) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{5, 0}
}

func (x *GetVenueResponse_Venue) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *GetVenueResponse_Venue) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *GetVenueResponse_Venue) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *GetVenueResponse_Venue) GetImageUrl() string {
	if x != nil {
		return x.ImageUrl
	}
	return ""
}

// VenueWrapper allows us to return an empty wrapper in the case of an invalid or unauthorized venue ID.
type GetVenueResponse_VenueWrapper struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Venue *GetVenueResponse_Venue `protobuf:"bytes,1,opt,name=venue,proto3,oneof" json:"venue,omitempty"`
}

func (x *GetVenueResponse_VenueWrapper) Reset() {
	*x = GetVenueResponse_VenueWrapper{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_pubgolf_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetVenueResponse_VenueWrapper) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetVenueResponse_VenueWrapper) ProtoMessage() {}

func (x *GetVenueResponse_VenueWrapper) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_pubgolf_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetVenueResponse_VenueWrapper.ProtoReflect.Descriptor instead.
func (*GetVenueResponse_VenueWrapper) Descriptor() ([]byte, []int) {
	return file_api_v1_pubgolf_proto_rawDescGZIP(), []int{5, 1}
}

func (x *GetVenueResponse_VenueWrapper) GetVenue() *GetVenueResponse_Venue {
	if x != nil {
		return x.Venue
	}
	return nil
}

var File_api_v1_pubgolf_proto protoreflect.FileDescriptor

var file_api_v1_pubgolf_proto_rawDesc = []byte{
	0x0a, 0x14, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x6c, 0x66,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x3d, 0x0a, 0x14, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x63, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x0d, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x22, 0xf2,
	0x01, 0x0a, 0x15, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x52, 0x0a, 0x0e, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x2b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e,
	0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x0d, 0x76,
	0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x84, 0x01, 0x0a,
	0x0d, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e,
	0x0a, 0x1a, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x15,
	0x0a, 0x11, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53,
	0x5f, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x1b, 0x0a, 0x17, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x4f, 0x55, 0x54, 0x44, 0x41, 0x54, 0x45, 0x44,
	0x10, 0x02, 0x12, 0x1f, 0x0a, 0x1b, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x43, 0x4f, 0x4d, 0x50, 0x41, 0x54, 0x49, 0x42, 0x4c,
	0x45, 0x10, 0x03, 0x22, 0x7e, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x33, 0x0a, 0x13, 0x63, 0x61, 0x63, 0x68, 0x65, 0x64,
	0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x11, 0x63, 0x61, 0x63, 0x68, 0x65, 0x64, 0x44, 0x61, 0x74,
	0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x88, 0x01, 0x01, 0x42, 0x16, 0x0a, 0x14, 0x5f,
	0x63, 0x61, 0x63, 0x68, 0x65, 0x64, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x76, 0x65, 0x72, 0x73,
	0x69, 0x6f, 0x6e, 0x22, 0xf2, 0x03, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x6c,
	0x61, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69,
	0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x11, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x45, 0x0a, 0x08, 0x73,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75,
	0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x53, 0x63, 0x68, 0x65, 0x64,
	0x75, 0x6c, 0x65, 0x48, 0x00, 0x52, 0x08, 0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x88,
	0x01, 0x01, 0x1a, 0xd6, 0x02, 0x0a, 0x08, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12,
	0x2c, 0x0a, 0x12, 0x76, 0x69, 0x73, 0x69, 0x74, 0x65, 0x64, 0x5f, 0x76, 0x65, 0x6e, 0x75, 0x65,
	0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x10, 0x76, 0x69, 0x73,
	0x69, 0x74, 0x65, 0x64, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x2f, 0x0a,
	0x11, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x5f, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x0f, 0x63, 0x75, 0x72, 0x72,
	0x65, 0x6e, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12, 0x29,
	0x0a, 0x0e, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x5f, 0x6b, 0x65, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x0c, 0x6e, 0x65, 0x78, 0x74, 0x56, 0x65,
	0x6e, 0x75, 0x65, 0x4b, 0x65, 0x79, 0x88, 0x01, 0x01, 0x12, 0x49, 0x0a, 0x10, 0x6e, 0x65, 0x78,
	0x74, 0x5f, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48,
	0x02, 0x52, 0x0e, 0x6e, 0x65, 0x78, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x88, 0x01, 0x01, 0x12, 0x37, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x65, 0x6e,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x52, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x45, 0x6e, 0x64, 0x42, 0x14, 0x0a,
	0x12, 0x5f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x5f,
	0x6b, 0x65, 0x79, 0x42, 0x11, 0x0a, 0x0f, 0x5f, 0x6e, 0x65, 0x78, 0x74, 0x5f, 0x76, 0x65, 0x6e,
	0x75, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x6e, 0x65, 0x78, 0x74, 0x5f,
	0x76, 0x65, 0x6e, 0x75, 0x65, 0x5f, 0x73, 0x74, 0x61, 0x72, 0x74, 0x42, 0x0b, 0x0a, 0x09, 0x5f,
	0x73, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x22, 0x4d, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x6e, 0x75, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65,
	0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x1d, 0x0a, 0x0a, 0x76, 0x65, 0x6e, 0x75,
	0x65, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x09, 0x76, 0x65,
	0x6e, 0x75, 0x65, 0x4b, 0x65, 0x79, 0x73, 0x22, 0xeb, 0x02, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x56,
	0x65, 0x6e, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x06,
	0x76, 0x65, 0x6e, 0x75, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x06, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x73, 0x1a, 0x62, 0x0a, 0x05, 0x56, 0x65,
	0x6e, 0x75, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x55, 0x72, 0x6c, 0x1a, 0x53,
	0x0a, 0x0c, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x12, 0x39,
	0x0a, 0x05, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x48, 0x00, 0x52,
	0x05, 0x76, 0x65, 0x6e, 0x75, 0x65, 0x88, 0x01, 0x01, 0x42, 0x08, 0x0a, 0x06, 0x5f, 0x76, 0x65,
	0x6e, 0x75, 0x65, 0x1a, 0x60, 0x0a, 0x0b, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x3b, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x56, 0x65, 0x6e, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x65,
	0x6e, 0x75, 0x65, 0x57, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0xeb, 0x01, 0x0a, 0x0e, 0x50, 0x75, 0x62, 0x47, 0x6f, 0x6c,
	0x66, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4e, 0x0a, 0x0d, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74,
	0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x12, 0x3f, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x12, 0x17,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x18, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x47, 0x65, 0x74, 0x56, 0x65, 0x6e, 0x75, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x6c, 0x66, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x6c,
	0x66, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c,
	0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b,
	0x61, 0x70, 0x69, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_pubgolf_proto_rawDescOnce sync.Once
	file_api_v1_pubgolf_proto_rawDescData = file_api_v1_pubgolf_proto_rawDesc
)

func file_api_v1_pubgolf_proto_rawDescGZIP() []byte {
	file_api_v1_pubgolf_proto_rawDescOnce.Do(func() {
		file_api_v1_pubgolf_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_pubgolf_proto_rawDescData)
	})
	return file_api_v1_pubgolf_proto_rawDescData
}

var file_api_v1_pubgolf_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_api_v1_pubgolf_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_api_v1_pubgolf_proto_goTypes = []interface{}{
	(ClientVersionResponse_VersionStatus)(0), // 0: api.v1.ClientVersionResponse.VersionStatus
	(*ClientVersionRequest)(nil),             // 1: api.v1.ClientVersionRequest
	(*ClientVersionResponse)(nil),            // 2: api.v1.ClientVersionResponse
	(*GetScheduleRequest)(nil),               // 3: api.v1.GetScheduleRequest
	(*GetScheduleResponse)(nil),              // 4: api.v1.GetScheduleResponse
	(*GetVenueRequest)(nil),                  // 5: api.v1.GetVenueRequest
	(*GetVenueResponse)(nil),                 // 6: api.v1.GetVenueResponse
	(*GetScheduleResponse_Schedule)(nil),     // 7: api.v1.GetScheduleResponse.Schedule
	(*GetVenueResponse_Venue)(nil),           // 8: api.v1.GetVenueResponse.Venue
	(*GetVenueResponse_VenueWrapper)(nil),    // 9: api.v1.GetVenueResponse.VenueWrapper
	nil,                                      // 10: api.v1.GetVenueResponse.VenuesEntry
	(*timestamppb.Timestamp)(nil),            // 11: google.protobuf.Timestamp
}
var file_api_v1_pubgolf_proto_depIdxs = []int32{
	0,  // 0: api.v1.ClientVersionResponse.version_status:type_name -> api.v1.ClientVersionResponse.VersionStatus
	7,  // 1: api.v1.GetScheduleResponse.schedule:type_name -> api.v1.GetScheduleResponse.Schedule
	10, // 2: api.v1.GetVenueResponse.venues:type_name -> api.v1.GetVenueResponse.VenuesEntry
	11, // 3: api.v1.GetScheduleResponse.Schedule.next_venue_start:type_name -> google.protobuf.Timestamp
	11, // 4: api.v1.GetScheduleResponse.Schedule.event_end:type_name -> google.protobuf.Timestamp
	8,  // 5: api.v1.GetVenueResponse.VenueWrapper.venue:type_name -> api.v1.GetVenueResponse.Venue
	9,  // 6: api.v1.GetVenueResponse.VenuesEntry.value:type_name -> api.v1.GetVenueResponse.VenueWrapper
	1,  // 7: api.v1.PubGolfService.ClientVersion:input_type -> api.v1.ClientVersionRequest
	3,  // 8: api.v1.PubGolfService.GetSchedule:input_type -> api.v1.GetScheduleRequest
	5,  // 9: api.v1.PubGolfService.GetVenue:input_type -> api.v1.GetVenueRequest
	2,  // 10: api.v1.PubGolfService.ClientVersion:output_type -> api.v1.ClientVersionResponse
	4,  // 11: api.v1.PubGolfService.GetSchedule:output_type -> api.v1.GetScheduleResponse
	6,  // 12: api.v1.PubGolfService.GetVenue:output_type -> api.v1.GetVenueResponse
	10, // [10:13] is the sub-list for method output_type
	7,  // [7:10] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_api_v1_pubgolf_proto_init() }
func file_api_v1_pubgolf_proto_init() {
	if File_api_v1_pubgolf_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_pubgolf_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientVersionRequest); i {
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
		file_api_v1_pubgolf_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientVersionResponse); i {
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
		file_api_v1_pubgolf_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetScheduleRequest); i {
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
		file_api_v1_pubgolf_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetScheduleResponse); i {
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
		file_api_v1_pubgolf_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVenueRequest); i {
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
		file_api_v1_pubgolf_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVenueResponse); i {
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
		file_api_v1_pubgolf_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetScheduleResponse_Schedule); i {
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
		file_api_v1_pubgolf_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVenueResponse_Venue); i {
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
		file_api_v1_pubgolf_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetVenueResponse_VenueWrapper); i {
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
	file_api_v1_pubgolf_proto_msgTypes[2].OneofWrappers = []interface{}{}
	file_api_v1_pubgolf_proto_msgTypes[3].OneofWrappers = []interface{}{}
	file_api_v1_pubgolf_proto_msgTypes[6].OneofWrappers = []interface{}{}
	file_api_v1_pubgolf_proto_msgTypes[8].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_pubgolf_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_pubgolf_proto_goTypes,
		DependencyIndexes: file_api_v1_pubgolf_proto_depIdxs,
		EnumInfos:         file_api_v1_pubgolf_proto_enumTypes,
		MessageInfos:      file_api_v1_pubgolf_proto_msgTypes,
	}.Build()
	File_api_v1_pubgolf_proto = out.File
	file_api_v1_pubgolf_proto_rawDesc = nil
	file_api_v1_pubgolf_proto_goTypes = nil
	file_api_v1_pubgolf_proto_depIdxs = nil
}