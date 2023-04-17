// Admin defines the admin API service for the game management UI.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: api/v1/admin.proto

package apiv1

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

type CreatePlayerRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventKey string                          `protobuf:"bytes,1,opt,name=event_key,json=eventKey,proto3" json:"event_key,omitempty"`
	Player   *CreatePlayerRequest_PlayerInfo `protobuf:"bytes,2,opt,name=player,proto3" json:"player,omitempty"`
}

func (x *CreatePlayerRequest) Reset() {
	*x = CreatePlayerRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePlayerRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePlayerRequest) ProtoMessage() {}

func (x *CreatePlayerRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePlayerRequest.ProtoReflect.Descriptor instead.
func (*CreatePlayerRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePlayerRequest) GetEventKey() string {
	if x != nil {
		return x.EventKey
	}
	return ""
}

func (x *CreatePlayerRequest) GetPlayer() *CreatePlayerRequest_PlayerInfo {
	if x != nil {
		return x.Player
	}
	return nil
}

type CreatePlayerResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId string `protobuf:"bytes,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"`
}

func (x *CreatePlayerResponse) Reset() {
	*x = CreatePlayerResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePlayerResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePlayerResponse) ProtoMessage() {}

func (x *CreatePlayerResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePlayerResponse.ProtoReflect.Descriptor instead.
func (*CreatePlayerResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePlayerResponse) GetPlayerId() string {
	if x != nil {
		return x.PlayerId
	}
	return ""
}

type ListPlayersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	EventKey string `protobuf:"bytes,1,opt,name=event_key,json=eventKey,proto3" json:"event_key,omitempty"`
}

func (x *ListPlayersRequest) Reset() {
	*x = ListPlayersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPlayersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPlayersRequest) ProtoMessage() {}

func (x *ListPlayersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPlayersRequest.ProtoReflect.Descriptor instead.
func (*ListPlayersRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{2}
}

func (x *ListPlayersRequest) GetEventKey() string {
	if x != nil {
		return x.EventKey
	}
	return ""
}

type ListPlayersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Players []*ListPlayersResponse_PlayerInfo `protobuf:"bytes,1,rep,name=players,proto3" json:"players,omitempty"`
}

func (x *ListPlayersResponse) Reset() {
	*x = ListPlayersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPlayersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPlayersResponse) ProtoMessage() {}

func (x *ListPlayersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPlayersResponse.ProtoReflect.Descriptor instead.
func (*ListPlayersResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{3}
}

func (x *ListPlayersResponse) GetPlayers() []*ListPlayersResponse_PlayerInfo {
	if x != nil {
		return x.Players
	}
	return nil
}

type CreatePlayerRequest_PlayerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name            string           `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	ScoringCategory *ScoringCategory `protobuf:"varint,2,opt,name=scoring_category,json=scoringCategory,proto3,enum=api.v1.ScoringCategory,oneof" json:"scoring_category,omitempty"`
}

func (x *CreatePlayerRequest_PlayerInfo) Reset() {
	*x = CreatePlayerRequest_PlayerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePlayerRequest_PlayerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePlayerRequest_PlayerInfo) ProtoMessage() {}

func (x *CreatePlayerRequest_PlayerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePlayerRequest_PlayerInfo.ProtoReflect.Descriptor instead.
func (*CreatePlayerRequest_PlayerInfo) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{0, 0}
}

func (x *CreatePlayerRequest_PlayerInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CreatePlayerRequest_PlayerInfo) GetScoringCategory() ScoringCategory {
	if x != nil && x.ScoringCategory != nil {
		return *x.ScoringCategory
	}
	return ScoringCategory_SCORING_CATEGORY_UNSPECIFIED
}

type ListPlayersResponse_PlayerInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string           `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ScoringCategory *ScoringCategory `protobuf:"varint,3,opt,name=scoring_category,json=scoringCategory,proto3,enum=api.v1.ScoringCategory,oneof" json:"scoring_category,omitempty"`
}

func (x *ListPlayersResponse_PlayerInfo) Reset() {
	*x = ListPlayersResponse_PlayerInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_admin_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListPlayersResponse_PlayerInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPlayersResponse_PlayerInfo) ProtoMessage() {}

func (x *ListPlayersResponse_PlayerInfo) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_admin_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPlayersResponse_PlayerInfo.ProtoReflect.Descriptor instead.
func (*ListPlayersResponse_PlayerInfo) Descriptor() ([]byte, []int) {
	return file_api_v1_admin_proto_rawDescGZIP(), []int{3, 0}
}

func (x *ListPlayersResponse_PlayerInfo) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ListPlayersResponse_PlayerInfo) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ListPlayersResponse_PlayerInfo) GetScoringCategory() ScoringCategory {
	if x != nil && x.ScoringCategory != nil {
		return *x.ScoringCategory
	}
	return ScoringCategory_SCORING_CATEGORY_UNSPECIFIED
}

var File_api_v1_admin_proto protoreflect.FileDescriptor

var file_api_v1_admin_proto_rawDesc = []byte{
	0x0a, 0x12, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x64, 0x6d, 0x69, 0x6e, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xf2, 0x01, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79,
	0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x76, 0x65,
	0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x3e, 0x0a, 0x06, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x06,
	0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x1a, 0x7e, 0x0a, 0x0a, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x49, 0x6e, 0x66, 0x6f, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x47, 0x0a, 0x10, 0x73, 0x63, 0x6f, 0x72,
	0x69, 0x6e, 0x67, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x6f, 0x72,
	0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x48, 0x00, 0x52, 0x0f, 0x73,
	0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x88, 0x01,
	0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x61,
	0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x22, 0x33, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1b,
	0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64, 0x22, 0x31, 0x0a, 0x12, 0x4c,
	0x69, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x4b, 0x65, 0x79, 0x22, 0xe8,
	0x01, 0x0a, 0x13, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x26, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x4c, 0x69, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x52,
	0x07, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x1a, 0x8e, 0x01, 0x0a, 0x0a, 0x50, 0x6c, 0x61,
	0x79, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x47, 0x0a, 0x10, 0x73,
	0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53,
	0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x48, 0x00,
	0x52, 0x0f, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67, 0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72,
	0x79, 0x88, 0x01, 0x01, 0x42, 0x13, 0x0a, 0x11, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x69, 0x6e, 0x67,
	0x5f, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x32, 0xa5, 0x01, 0x0a, 0x0c, 0x41, 0x64,
	0x6d, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4b, 0x0a, 0x0c, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x12, 0x1b, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x50,
	0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x12, 0x1a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x4c, 0x69, 0x73, 0x74, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x42, 0x40, 0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x70, 0x75, 0x62, 0x67, 0x6f, 0x6c, 0x66, 0x2f, 0x70, 0x75, 0x62, 0x67, 0x6f, 0x6c, 0x66, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c, 0x69, 0x62,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61, 0x70,
	0x69, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_admin_proto_rawDescOnce sync.Once
	file_api_v1_admin_proto_rawDescData = file_api_v1_admin_proto_rawDesc
)

func file_api_v1_admin_proto_rawDescGZIP() []byte {
	file_api_v1_admin_proto_rawDescOnce.Do(func() {
		file_api_v1_admin_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_admin_proto_rawDescData)
	})
	return file_api_v1_admin_proto_rawDescData
}

var file_api_v1_admin_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_api_v1_admin_proto_goTypes = []interface{}{
	(*CreatePlayerRequest)(nil),            // 0: api.v1.CreatePlayerRequest
	(*CreatePlayerResponse)(nil),           // 1: api.v1.CreatePlayerResponse
	(*ListPlayersRequest)(nil),             // 2: api.v1.ListPlayersRequest
	(*ListPlayersResponse)(nil),            // 3: api.v1.ListPlayersResponse
	(*CreatePlayerRequest_PlayerInfo)(nil), // 4: api.v1.CreatePlayerRequest.PlayerInfo
	(*ListPlayersResponse_PlayerInfo)(nil), // 5: api.v1.ListPlayersResponse.PlayerInfo
	(ScoringCategory)(0),                   // 6: api.v1.ScoringCategory
}
var file_api_v1_admin_proto_depIdxs = []int32{
	4, // 0: api.v1.CreatePlayerRequest.player:type_name -> api.v1.CreatePlayerRequest.PlayerInfo
	5, // 1: api.v1.ListPlayersResponse.players:type_name -> api.v1.ListPlayersResponse.PlayerInfo
	6, // 2: api.v1.CreatePlayerRequest.PlayerInfo.scoring_category:type_name -> api.v1.ScoringCategory
	6, // 3: api.v1.ListPlayersResponse.PlayerInfo.scoring_category:type_name -> api.v1.ScoringCategory
	0, // 4: api.v1.AdminService.CreatePlayer:input_type -> api.v1.CreatePlayerRequest
	2, // 5: api.v1.AdminService.ListPlayers:input_type -> api.v1.ListPlayersRequest
	1, // 6: api.v1.AdminService.CreatePlayer:output_type -> api.v1.CreatePlayerResponse
	3, // 7: api.v1.AdminService.ListPlayers:output_type -> api.v1.ListPlayersResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1_admin_proto_init() }
func file_api_v1_admin_proto_init() {
	if File_api_v1_admin_proto != nil {
		return
	}
	file_api_v1_shared_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_admin_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePlayerRequest); i {
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
		file_api_v1_admin_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePlayerResponse); i {
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
		file_api_v1_admin_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPlayersRequest); i {
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
		file_api_v1_admin_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPlayersResponse); i {
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
		file_api_v1_admin_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePlayerRequest_PlayerInfo); i {
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
		file_api_v1_admin_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListPlayersResponse_PlayerInfo); i {
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
	file_api_v1_admin_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_api_v1_admin_proto_msgTypes[5].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_api_v1_admin_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_admin_proto_goTypes,
		DependencyIndexes: file_api_v1_admin_proto_depIdxs,
		MessageInfos:      file_api_v1_admin_proto_msgTypes,
	}.Build()
	File_api_v1_admin_proto = out.File
	file_api_v1_admin_proto_rawDesc = nil
	file_api_v1_admin_proto_goTypes = nil
	file_api_v1_admin_proto_depIdxs = nil
}
