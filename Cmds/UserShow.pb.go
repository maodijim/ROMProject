// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.4.0
// source: UserShow.proto

package Cmd

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

type EUserShowParam int32

const (
	EUserShowParam_EUSERSHOW_NEW_PHOTO_FRAME       EUserShowParam = 1
	EUserShowParam_EUSERSHOW_SYNC_PHOTO_FRAME      EUserShowParam = 2
	EUserShowParam_EUSERSHOW_NEW_BACKGROUND_FRAME  EUserShowParam = 3
	EUserShowParam_EUSERSHOW_SYNC_BACKGROUND_FRAME EUserShowParam = 4
	EUserShowParam_EUSERSHOW_USE_PHOTO_FRAME       EUserShowParam = 5
	EUserShowParam_EUSERSHOW_USE_BACKGROUND_FRAME  EUserShowParam = 6
)

// Enum value maps for EUserShowParam.
var (
	EUserShowParam_name = map[int32]string{
		1: "EUSERSHOW_NEW_PHOTO_FRAME",
		2: "EUSERSHOW_SYNC_PHOTO_FRAME",
		3: "EUSERSHOW_NEW_BACKGROUND_FRAME",
		4: "EUSERSHOW_SYNC_BACKGROUND_FRAME",
		5: "EUSERSHOW_USE_PHOTO_FRAME",
		6: "EUSERSHOW_USE_BACKGROUND_FRAME",
	}
	EUserShowParam_value = map[string]int32{
		"EUSERSHOW_NEW_PHOTO_FRAME":       1,
		"EUSERSHOW_SYNC_PHOTO_FRAME":      2,
		"EUSERSHOW_NEW_BACKGROUND_FRAME":  3,
		"EUSERSHOW_SYNC_BACKGROUND_FRAME": 4,
		"EUSERSHOW_USE_PHOTO_FRAME":       5,
		"EUSERSHOW_USE_BACKGROUND_FRAME":  6,
	}
)

func (x EUserShowParam) Enum() *EUserShowParam {
	p := new(EUserShowParam)
	*p = x
	return p
}

func (x EUserShowParam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EUserShowParam) Descriptor() protoreflect.EnumDescriptor {
	return file_UserShow_proto_enumTypes[0].Descriptor()
}

func (EUserShowParam) Type() protoreflect.EnumType {
	return &file_UserShow_proto_enumTypes[0]
}

func (x EUserShowParam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *EUserShowParam) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = EUserShowParam(num)
	return nil
}

// Deprecated: Use EUserShowParam.Descriptor instead.
func (EUserShowParam) EnumDescriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{0}
}

type UnlockPhotoFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=1" json:"param,omitempty"`
	Id    *uint32         `protobuf:"varint,3,opt,name=id,def=0" json:"id,omitempty"`
	Del   *bool           `protobuf:"varint,4,opt,name=del,def=0" json:"del,omitempty"`
}

// Default values for UnlockPhotoFrame fields.
const (
	Default_UnlockPhotoFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_UnlockPhotoFrame_Param = EUserShowParam_EUSERSHOW_NEW_PHOTO_FRAME
	Default_UnlockPhotoFrame_Id    = uint32(0)
	Default_UnlockPhotoFrame_Del   = bool(false)
)

func (x *UnlockPhotoFrame) Reset() {
	*x = UnlockPhotoFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockPhotoFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockPhotoFrame) ProtoMessage() {}

func (x *UnlockPhotoFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockPhotoFrame.ProtoReflect.Descriptor instead.
func (*UnlockPhotoFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{0}
}

func (x *UnlockPhotoFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_UnlockPhotoFrame_Cmd
}

func (x *UnlockPhotoFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_UnlockPhotoFrame_Param
}

func (x *UnlockPhotoFrame) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_UnlockPhotoFrame_Id
}

func (x *UnlockPhotoFrame) GetDel() bool {
	if x != nil && x.Del != nil {
		return *x.Del
	}
	return Default_UnlockPhotoFrame_Del
}

type SyncAllPhotoFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=2" json:"param,omitempty"`
	Ids   []uint32        `protobuf:"varint,3,rep,name=ids" json:"ids,omitempty"`
}

// Default values for SyncAllPhotoFrame fields.
const (
	Default_SyncAllPhotoFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_SyncAllPhotoFrame_Param = EUserShowParam_EUSERSHOW_SYNC_PHOTO_FRAME
)

func (x *SyncAllPhotoFrame) Reset() {
	*x = SyncAllPhotoFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllPhotoFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllPhotoFrame) ProtoMessage() {}

func (x *SyncAllPhotoFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllPhotoFrame.ProtoReflect.Descriptor instead.
func (*SyncAllPhotoFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{1}
}

func (x *SyncAllPhotoFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_SyncAllPhotoFrame_Cmd
}

func (x *SyncAllPhotoFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_SyncAllPhotoFrame_Param
}

func (x *SyncAllPhotoFrame) GetIds() []uint32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type SelectPhotoFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=5" json:"param,omitempty"`
	Id    *uint32         `protobuf:"varint,3,opt,name=id,def=0" json:"id,omitempty"`
}

// Default values for SelectPhotoFrame fields.
const (
	Default_SelectPhotoFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_SelectPhotoFrame_Param = EUserShowParam_EUSERSHOW_USE_PHOTO_FRAME
	Default_SelectPhotoFrame_Id    = uint32(0)
)

func (x *SelectPhotoFrame) Reset() {
	*x = SelectPhotoFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectPhotoFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectPhotoFrame) ProtoMessage() {}

func (x *SelectPhotoFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectPhotoFrame.ProtoReflect.Descriptor instead.
func (*SelectPhotoFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{2}
}

func (x *SelectPhotoFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_SelectPhotoFrame_Cmd
}

func (x *SelectPhotoFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_SelectPhotoFrame_Param
}

func (x *SelectPhotoFrame) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_SelectPhotoFrame_Id
}

type UnlockBackgroundFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=3" json:"param,omitempty"`
	Id    *uint32         `protobuf:"varint,3,opt,name=id,def=0" json:"id,omitempty"`
	Del   *bool           `protobuf:"varint,4,opt,name=del,def=0" json:"del,omitempty"`
}

// Default values for UnlockBackgroundFrame fields.
const (
	Default_UnlockBackgroundFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_UnlockBackgroundFrame_Param = EUserShowParam_EUSERSHOW_NEW_BACKGROUND_FRAME
	Default_UnlockBackgroundFrame_Id    = uint32(0)
	Default_UnlockBackgroundFrame_Del   = bool(false)
)

func (x *UnlockBackgroundFrame) Reset() {
	*x = UnlockBackgroundFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UnlockBackgroundFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UnlockBackgroundFrame) ProtoMessage() {}

func (x *UnlockBackgroundFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UnlockBackgroundFrame.ProtoReflect.Descriptor instead.
func (*UnlockBackgroundFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{3}
}

func (x *UnlockBackgroundFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_UnlockBackgroundFrame_Cmd
}

func (x *UnlockBackgroundFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_UnlockBackgroundFrame_Param
}

func (x *UnlockBackgroundFrame) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_UnlockBackgroundFrame_Id
}

func (x *UnlockBackgroundFrame) GetDel() bool {
	if x != nil && x.Del != nil {
		return *x.Del
	}
	return Default_UnlockBackgroundFrame_Del
}

type SyncAllBackgroundFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=4" json:"param,omitempty"`
	Ids   []uint32        `protobuf:"varint,3,rep,name=ids" json:"ids,omitempty"`
}

// Default values for SyncAllBackgroundFrame fields.
const (
	Default_SyncAllBackgroundFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_SyncAllBackgroundFrame_Param = EUserShowParam_EUSERSHOW_SYNC_BACKGROUND_FRAME
)

func (x *SyncAllBackgroundFrame) Reset() {
	*x = SyncAllBackgroundFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncAllBackgroundFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncAllBackgroundFrame) ProtoMessage() {}

func (x *SyncAllBackgroundFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncAllBackgroundFrame.ProtoReflect.Descriptor instead.
func (*SyncAllBackgroundFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{4}
}

func (x *SyncAllBackgroundFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_SyncAllBackgroundFrame_Cmd
}

func (x *SyncAllBackgroundFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_SyncAllBackgroundFrame_Param
}

func (x *SyncAllBackgroundFrame) GetIds() []uint32 {
	if x != nil {
		return x.Ids
	}
	return nil
}

type SelectBackgroundFrame struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=225" json:"cmd,omitempty"`
	Param *EUserShowParam `protobuf:"varint,2,opt,name=param,enum=Cmd.EUserShowParam,def=6" json:"param,omitempty"`
	Id    *uint32         `protobuf:"varint,3,opt,name=id,def=0" json:"id,omitempty"`
}

// Default values for SelectBackgroundFrame fields.
const (
	Default_SelectBackgroundFrame_Cmd   = Command_USERSHOW_PROTOCMD
	Default_SelectBackgroundFrame_Param = EUserShowParam_EUSERSHOW_USE_BACKGROUND_FRAME
	Default_SelectBackgroundFrame_Id    = uint32(0)
)

func (x *SelectBackgroundFrame) Reset() {
	*x = SelectBackgroundFrame{}
	if protoimpl.UnsafeEnabled {
		mi := &file_UserShow_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SelectBackgroundFrame) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SelectBackgroundFrame) ProtoMessage() {}

func (x *SelectBackgroundFrame) ProtoReflect() protoreflect.Message {
	mi := &file_UserShow_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SelectBackgroundFrame.ProtoReflect.Descriptor instead.
func (*SelectBackgroundFrame) Descriptor() ([]byte, []int) {
	return file_UserShow_proto_rawDescGZIP(), []int{5}
}

func (x *SelectBackgroundFrame) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_SelectBackgroundFrame_Cmd
}

func (x *SelectBackgroundFrame) GetParam() EUserShowParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_SelectBackgroundFrame_Param
}

func (x *SelectBackgroundFrame) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_SelectBackgroundFrame_Id
}

var File_UserShow_proto protoreflect.FileDescriptor

var file_UserShow_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x03, 0x43, 0x6d, 0x64, 0x1a, 0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xb7, 0x01, 0x0a, 0x10, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x68, 0x6f, 0x74,
	0x6f, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x3a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x44, 0x0a, 0x05, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45,
	0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x19, 0x45,
	0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x4e, 0x45, 0x57, 0x5f, 0x50, 0x48, 0x4f,
	0x54, 0x4f, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12,
	0x11, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x17, 0x0a, 0x03, 0x64, 0x65, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x3a,
	0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x03, 0x64, 0x65, 0x6c, 0x22, 0x9f, 0x01, 0x0a, 0x11,
	0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x46, 0x72, 0x61, 0x6d,
	0x65, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x55, 0x53,
	0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52,
	0x03, 0x63, 0x6d, 0x64, 0x12, 0x45, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x55, 0x73, 0x65, 0x72, 0x53,
	0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1a, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53,
	0x48, 0x4f, 0x57, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f, 0x50, 0x48, 0x4f, 0x54, 0x4f, 0x5f, 0x46,
	0x52, 0x41, 0x4d, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x10, 0x0a, 0x03, 0x69,
	0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0x9e, 0x01,
	0x0a, 0x10, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x46, 0x72, 0x61,
	0x6d, 0x65, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x55,
	0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44,
	0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x44, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x55, 0x73, 0x65, 0x72,
	0x53, 0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x19, 0x45, 0x55, 0x53, 0x45, 0x52,
	0x53, 0x48, 0x4f, 0x57, 0x5f, 0x55, 0x53, 0x45, 0x5f, 0x50, 0x48, 0x4f, 0x54, 0x4f, 0x5f, 0x46,
	0x52, 0x41, 0x4d, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x11, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x02, 0x69, 0x64, 0x22, 0xc1,
	0x01, 0x0a, 0x15, 0x55, 0x6e, 0x6c, 0x6f, 0x63, 0x6b, 0x42, 0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x3a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x49, 0x0a, 0x05, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x45, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a,
	0x1e, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x4e, 0x45, 0x57, 0x5f, 0x42,
	0x41, 0x43, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x52,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x11, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x03, 0x64, 0x65, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x03, 0x64,
	0x65, 0x6c, 0x22, 0xa9, 0x01, 0x0a, 0x16, 0x53, 0x79, 0x6e, 0x63, 0x41, 0x6c, 0x6c, 0x42, 0x61,
	0x63, 0x6b, 0x67, 0x72, 0x6f, 0x75, 0x6e, 0x64, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a,
	0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48,
	0x4f, 0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64,
	0x12, 0x4a, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1f, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f,
	0x53, 0x59, 0x4e, 0x43, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f,
	0x46, 0x52, 0x41, 0x4d, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x10, 0x0a, 0x03,
	0x69, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x03, 0x69, 0x64, 0x73, 0x22, 0xa8,
	0x01, 0x0a, 0x15, 0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x42, 0x61, 0x63, 0x6b, 0x67, 0x72, 0x6f,
	0x75, 0x6e, 0x64, 0x46, 0x72, 0x61, 0x6d, 0x65, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x3a, 0x11, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x49, 0x0a, 0x05, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x45, 0x55, 0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a,
	0x1e, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x55, 0x53, 0x45, 0x5f, 0x42,
	0x41, 0x43, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x52,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x11, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x02, 0x69, 0x64, 0x2a, 0xdb, 0x01, 0x0a, 0x0e, 0x45, 0x55,
	0x73, 0x65, 0x72, 0x53, 0x68, 0x6f, 0x77, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1d, 0x0a, 0x19,
	0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x4e, 0x45, 0x57, 0x5f, 0x50, 0x48,
	0x4f, 0x54, 0x4f, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x45,
	0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x5f, 0x50, 0x48,
	0x4f, 0x54, 0x4f, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x10, 0x02, 0x12, 0x22, 0x0a, 0x1e, 0x45,
	0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x4e, 0x45, 0x57, 0x5f, 0x42, 0x41, 0x43,
	0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x10, 0x03, 0x12,
	0x23, 0x0a, 0x1f, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x53, 0x59, 0x4e,
	0x43, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f, 0x46, 0x52, 0x41,
	0x4d, 0x45, 0x10, 0x04, 0x12, 0x1d, 0x0a, 0x19, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f,
	0x57, 0x5f, 0x55, 0x53, 0x45, 0x5f, 0x50, 0x48, 0x4f, 0x54, 0x4f, 0x5f, 0x46, 0x52, 0x41, 0x4d,
	0x45, 0x10, 0x05, 0x12, 0x22, 0x0a, 0x1e, 0x45, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57,
	0x5f, 0x55, 0x53, 0x45, 0x5f, 0x42, 0x41, 0x43, 0x4b, 0x47, 0x52, 0x4f, 0x55, 0x4e, 0x44, 0x5f,
	0x46, 0x52, 0x41, 0x4d, 0x45, 0x10, 0x06,
}

var (
	file_UserShow_proto_rawDescOnce sync.Once
	file_UserShow_proto_rawDescData = file_UserShow_proto_rawDesc
)

func file_UserShow_proto_rawDescGZIP() []byte {
	file_UserShow_proto_rawDescOnce.Do(func() {
		file_UserShow_proto_rawDescData = protoimpl.X.CompressGZIP(file_UserShow_proto_rawDescData)
	})
	return file_UserShow_proto_rawDescData
}

var file_UserShow_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_UserShow_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_UserShow_proto_goTypes = []interface{}{
	(EUserShowParam)(0),            // 0: Cmd.EUserShowParam
	(*UnlockPhotoFrame)(nil),       // 1: Cmd.UnlockPhotoFrame
	(*SyncAllPhotoFrame)(nil),      // 2: Cmd.SyncAllPhotoFrame
	(*SelectPhotoFrame)(nil),       // 3: Cmd.SelectPhotoFrame
	(*UnlockBackgroundFrame)(nil),  // 4: Cmd.UnlockBackgroundFrame
	(*SyncAllBackgroundFrame)(nil), // 5: Cmd.SyncAllBackgroundFrame
	(*SelectBackgroundFrame)(nil),  // 6: Cmd.SelectBackgroundFrame
	(Command)(0),                   // 7: Cmd.Command
}
var file_UserShow_proto_depIdxs = []int32{
	7,  // 0: Cmd.UnlockPhotoFrame.cmd:type_name -> Cmd.Command
	0,  // 1: Cmd.UnlockPhotoFrame.param:type_name -> Cmd.EUserShowParam
	7,  // 2: Cmd.SyncAllPhotoFrame.cmd:type_name -> Cmd.Command
	0,  // 3: Cmd.SyncAllPhotoFrame.param:type_name -> Cmd.EUserShowParam
	7,  // 4: Cmd.SelectPhotoFrame.cmd:type_name -> Cmd.Command
	0,  // 5: Cmd.SelectPhotoFrame.param:type_name -> Cmd.EUserShowParam
	7,  // 6: Cmd.UnlockBackgroundFrame.cmd:type_name -> Cmd.Command
	0,  // 7: Cmd.UnlockBackgroundFrame.param:type_name -> Cmd.EUserShowParam
	7,  // 8: Cmd.SyncAllBackgroundFrame.cmd:type_name -> Cmd.Command
	0,  // 9: Cmd.SyncAllBackgroundFrame.param:type_name -> Cmd.EUserShowParam
	7,  // 10: Cmd.SelectBackgroundFrame.cmd:type_name -> Cmd.Command
	0,  // 11: Cmd.SelectBackgroundFrame.param:type_name -> Cmd.EUserShowParam
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_UserShow_proto_init() }
func file_UserShow_proto_init() {
	if File_UserShow_proto != nil {
		return
	}
	file_xCmd_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_UserShow_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockPhotoFrame); i {
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
		file_UserShow_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncAllPhotoFrame); i {
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
		file_UserShow_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectPhotoFrame); i {
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
		file_UserShow_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UnlockBackgroundFrame); i {
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
		file_UserShow_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncAllBackgroundFrame); i {
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
		file_UserShow_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SelectBackgroundFrame); i {
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
			RawDescriptor: file_UserShow_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_UserShow_proto_goTypes,
		DependencyIndexes: file_UserShow_proto_depIdxs,
		EnumInfos:         file_UserShow_proto_enumTypes,
		MessageInfos:      file_UserShow_proto_msgTypes,
	}.Build()
	File_UserShow_proto = out.File
	file_UserShow_proto_rawDesc = nil
	file_UserShow_proto_goTypes = nil
	file_UserShow_proto_depIdxs = nil
}
