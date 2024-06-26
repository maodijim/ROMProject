// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.4.0
// source: AstrolabeCmd.proto

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

type AstrolabeParam int32

const (
	AstrolabeParam_ASTROLABEPARAM_QUERY         AstrolabeParam = 1
	AstrolabeParam_ASTROLABEPARAM_ACTIVATE_STAR AstrolabeParam = 2
	AstrolabeParam_ASTROLABEPARAM_QUERY_RESET   AstrolabeParam = 3
	AstrolabeParam_ASTROLABEPARAM_RESET         AstrolabeParam = 4
	AstrolabeParam_ASTROLABEPARAM_PLAN_SAVE     AstrolabeParam = 5
)

// Enum value maps for AstrolabeParam.
var (
	AstrolabeParam_name = map[int32]string{
		1: "ASTROLABEPARAM_QUERY",
		2: "ASTROLABEPARAM_ACTIVATE_STAR",
		3: "ASTROLABEPARAM_QUERY_RESET",
		4: "ASTROLABEPARAM_RESET",
		5: "ASTROLABEPARAM_PLAN_SAVE",
	}
	AstrolabeParam_value = map[string]int32{
		"ASTROLABEPARAM_QUERY":         1,
		"ASTROLABEPARAM_ACTIVATE_STAR": 2,
		"ASTROLABEPARAM_QUERY_RESET":   3,
		"ASTROLABEPARAM_RESET":         4,
		"ASTROLABEPARAM_PLAN_SAVE":     5,
	}
)

func (x AstrolabeParam) Enum() *AstrolabeParam {
	p := new(AstrolabeParam)
	*p = x
	return p
}

func (x AstrolabeParam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AstrolabeParam) Descriptor() protoreflect.EnumDescriptor {
	return file_AstrolabeCmd_proto_enumTypes[0].Descriptor()
}

func (AstrolabeParam) Type() protoreflect.EnumType {
	return &file_AstrolabeCmd_proto_enumTypes[0]
}

func (x AstrolabeParam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *AstrolabeParam) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = AstrolabeParam(num)
	return nil
}

// Deprecated: Use AstrolabeParam.Descriptor instead.
func (AstrolabeParam) EnumDescriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{0}
}

type EAstrolabeType int32

const (
	EAstrolabeType_EASTROLABETYPE_MIN        EAstrolabeType = 0
	EAstrolabeType_EASTROLABETYPE_PROFESSION EAstrolabeType = 1
	EAstrolabeType_EASTROLABETYPE_PLAN       EAstrolabeType = 100
	EAstrolabeType_EASTROLABETYPE_MAX        EAstrolabeType = 101
)

// Enum value maps for EAstrolabeType.
var (
	EAstrolabeType_name = map[int32]string{
		0:   "EASTROLABETYPE_MIN",
		1:   "EASTROLABETYPE_PROFESSION",
		100: "EASTROLABETYPE_PLAN",
		101: "EASTROLABETYPE_MAX",
	}
	EAstrolabeType_value = map[string]int32{
		"EASTROLABETYPE_MIN":        0,
		"EASTROLABETYPE_PROFESSION": 1,
		"EASTROLABETYPE_PLAN":       100,
		"EASTROLABETYPE_MAX":        101,
	}
)

func (x EAstrolabeType) Enum() *EAstrolabeType {
	p := new(EAstrolabeType)
	*p = x
	return p
}

func (x EAstrolabeType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EAstrolabeType) Descriptor() protoreflect.EnumDescriptor {
	return file_AstrolabeCmd_proto_enumTypes[1].Descriptor()
}

func (EAstrolabeType) Type() protoreflect.EnumType {
	return &file_AstrolabeCmd_proto_enumTypes[1]
}

func (x EAstrolabeType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *EAstrolabeType) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = EAstrolabeType(num)
	return nil
}

// Deprecated: Use EAstrolabeType.Descriptor instead.
func (EAstrolabeType) EnumDescriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{1}
}

type AstrolabeCostData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    *uint32 `protobuf:"varint,1,opt,name=id,def=0" json:"id,omitempty"`
	Count *uint32 `protobuf:"varint,2,opt,name=count,def=0" json:"count,omitempty"`
}

// Default values for AstrolabeCostData fields.
const (
	Default_AstrolabeCostData_Id    = uint32(0)
	Default_AstrolabeCostData_Count = uint32(0)
)

func (x *AstrolabeCostData) Reset() {
	*x = AstrolabeCostData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabeCostData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabeCostData) ProtoMessage() {}

func (x *AstrolabeCostData) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabeCostData.ProtoReflect.Descriptor instead.
func (*AstrolabeCostData) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{0}
}

func (x *AstrolabeCostData) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return Default_AstrolabeCostData_Id
}

func (x *AstrolabeCostData) GetCount() uint32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return Default_AstrolabeCostData_Count
}

type AstrolabeQueryCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd           *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=28" json:"cmd,omitempty"`
	Param         *AstrolabeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.AstrolabeParam,def=1" json:"param,omitempty"`
	Stars         []uint32        `protobuf:"varint,3,rep,name=stars" json:"stars,omitempty"`
	Astrolabetype *EAstrolabeType `protobuf:"varint,4,opt,name=astrolabetype,enum=Cmd.EAstrolabeType,def=0" json:"astrolabetype,omitempty"`
}

// Default values for AstrolabeQueryCmd fields.
const (
	Default_AstrolabeQueryCmd_Cmd           = Command_SCENE_USER_ASTROLABE_PROTOCMD
	Default_AstrolabeQueryCmd_Param         = AstrolabeParam_ASTROLABEPARAM_QUERY
	Default_AstrolabeQueryCmd_Astrolabetype = EAstrolabeType_EASTROLABETYPE_MIN
)

func (x *AstrolabeQueryCmd) Reset() {
	*x = AstrolabeQueryCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabeQueryCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabeQueryCmd) ProtoMessage() {}

func (x *AstrolabeQueryCmd) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabeQueryCmd.ProtoReflect.Descriptor instead.
func (*AstrolabeQueryCmd) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{1}
}

func (x *AstrolabeQueryCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AstrolabeQueryCmd_Cmd
}

func (x *AstrolabeQueryCmd) GetParam() AstrolabeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AstrolabeQueryCmd_Param
}

func (x *AstrolabeQueryCmd) GetStars() []uint32 {
	if x != nil {
		return x.Stars
	}
	return nil
}

func (x *AstrolabeQueryCmd) GetAstrolabetype() EAstrolabeType {
	if x != nil && x.Astrolabetype != nil {
		return *x.Astrolabetype
	}
	return Default_AstrolabeQueryCmd_Astrolabetype
}

type AstrolabeActivateStarCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd     *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=28" json:"cmd,omitempty"`
	Param   *AstrolabeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.AstrolabeParam,def=2" json:"param,omitempty"`
	Stars   []uint32        `protobuf:"varint,3,rep,name=stars" json:"stars,omitempty"`
	Success *bool           `protobuf:"varint,5,opt,name=success" json:"success,omitempty"`
}

// Default values for AstrolabeActivateStarCmd fields.
const (
	Default_AstrolabeActivateStarCmd_Cmd   = Command_SCENE_USER_ASTROLABE_PROTOCMD
	Default_AstrolabeActivateStarCmd_Param = AstrolabeParam_ASTROLABEPARAM_ACTIVATE_STAR
)

func (x *AstrolabeActivateStarCmd) Reset() {
	*x = AstrolabeActivateStarCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabeActivateStarCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabeActivateStarCmd) ProtoMessage() {}

func (x *AstrolabeActivateStarCmd) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabeActivateStarCmd.ProtoReflect.Descriptor instead.
func (*AstrolabeActivateStarCmd) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{2}
}

func (x *AstrolabeActivateStarCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AstrolabeActivateStarCmd_Cmd
}

func (x *AstrolabeActivateStarCmd) GetParam() AstrolabeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AstrolabeActivateStarCmd_Param
}

func (x *AstrolabeActivateStarCmd) GetStars() []uint32 {
	if x != nil {
		return x.Stars
	}
	return nil
}

func (x *AstrolabeActivateStarCmd) GetSuccess() bool {
	if x != nil && x.Success != nil {
		return *x.Success
	}
	return false
}

type AstrolabeQueryResetCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command             `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=28" json:"cmd,omitempty"`
	Param *AstrolabeParam      `protobuf:"varint,2,opt,name=param,enum=Cmd.AstrolabeParam,def=3" json:"param,omitempty"`
	Type  *EAstrolabeType      `protobuf:"varint,3,opt,name=type,enum=Cmd.EAstrolabeType" json:"type,omitempty"`
	Items []*AstrolabeCostData `protobuf:"bytes,4,rep,name=items" json:"items,omitempty"`
}

// Default values for AstrolabeQueryResetCmd fields.
const (
	Default_AstrolabeQueryResetCmd_Cmd   = Command_SCENE_USER_ASTROLABE_PROTOCMD
	Default_AstrolabeQueryResetCmd_Param = AstrolabeParam_ASTROLABEPARAM_QUERY_RESET
)

func (x *AstrolabeQueryResetCmd) Reset() {
	*x = AstrolabeQueryResetCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabeQueryResetCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabeQueryResetCmd) ProtoMessage() {}

func (x *AstrolabeQueryResetCmd) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabeQueryResetCmd.ProtoReflect.Descriptor instead.
func (*AstrolabeQueryResetCmd) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{3}
}

func (x *AstrolabeQueryResetCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AstrolabeQueryResetCmd_Cmd
}

func (x *AstrolabeQueryResetCmd) GetParam() AstrolabeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AstrolabeQueryResetCmd_Param
}

func (x *AstrolabeQueryResetCmd) GetType() EAstrolabeType {
	if x != nil && x.Type != nil {
		return *x.Type
	}
	return EAstrolabeType_EASTROLABETYPE_MIN
}

func (x *AstrolabeQueryResetCmd) GetItems() []*AstrolabeCostData {
	if x != nil {
		return x.Items
	}
	return nil
}

type AstrolabeResetCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd     *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=28" json:"cmd,omitempty"`
	Param   *AstrolabeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.AstrolabeParam,def=4" json:"param,omitempty"`
	Stars   []uint32        `protobuf:"varint,3,rep,name=stars" json:"stars,omitempty"`
	Success *bool           `protobuf:"varint,4,opt,name=success" json:"success,omitempty"`
}

// Default values for AstrolabeResetCmd fields.
const (
	Default_AstrolabeResetCmd_Cmd   = Command_SCENE_USER_ASTROLABE_PROTOCMD
	Default_AstrolabeResetCmd_Param = AstrolabeParam_ASTROLABEPARAM_RESET
)

func (x *AstrolabeResetCmd) Reset() {
	*x = AstrolabeResetCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabeResetCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabeResetCmd) ProtoMessage() {}

func (x *AstrolabeResetCmd) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabeResetCmd.ProtoReflect.Descriptor instead.
func (*AstrolabeResetCmd) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{4}
}

func (x *AstrolabeResetCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AstrolabeResetCmd_Cmd
}

func (x *AstrolabeResetCmd) GetParam() AstrolabeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AstrolabeResetCmd_Param
}

func (x *AstrolabeResetCmd) GetStars() []uint32 {
	if x != nil {
		return x.Stars
	}
	return nil
}

func (x *AstrolabeResetCmd) GetSuccess() bool {
	if x != nil && x.Success != nil {
		return *x.Success
	}
	return false
}

type AstrolabePlanSaveCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command        `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=28" json:"cmd,omitempty"`
	Param *AstrolabeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.AstrolabeParam,def=5" json:"param,omitempty"`
	Stars []uint32        `protobuf:"varint,3,rep,name=stars" json:"stars,omitempty"`
}

// Default values for AstrolabePlanSaveCmd fields.
const (
	Default_AstrolabePlanSaveCmd_Cmd   = Command_SCENE_USER_ASTROLABE_PROTOCMD
	Default_AstrolabePlanSaveCmd_Param = AstrolabeParam_ASTROLABEPARAM_PLAN_SAVE
)

func (x *AstrolabePlanSaveCmd) Reset() {
	*x = AstrolabePlanSaveCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_AstrolabeCmd_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AstrolabePlanSaveCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AstrolabePlanSaveCmd) ProtoMessage() {}

func (x *AstrolabePlanSaveCmd) ProtoReflect() protoreflect.Message {
	mi := &file_AstrolabeCmd_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AstrolabePlanSaveCmd.ProtoReflect.Descriptor instead.
func (*AstrolabePlanSaveCmd) Descriptor() ([]byte, []int) {
	return file_AstrolabeCmd_proto_rawDescGZIP(), []int{5}
}

func (x *AstrolabePlanSaveCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AstrolabePlanSaveCmd_Cmd
}

func (x *AstrolabePlanSaveCmd) GetParam() AstrolabeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AstrolabePlanSaveCmd_Param
}

func (x *AstrolabePlanSaveCmd) GetStars() []uint32 {
	if x != nil {
		return x.Stars
	}
	return nil
}

var File_AstrolabeCmd_proto protoreflect.FileDescriptor

var file_AstrolabeCmd_proto_rawDesc = []byte{
	0x0a, 0x12, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x43, 0x6d, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x43, 0x6d, 0x64, 0x1a, 0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x11, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61,
	0x62, 0x65, 0x43, 0x6f, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x11, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52,
	0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xf8, 0x01, 0x0a, 0x11, 0x41, 0x73, 0x74, 0x72, 0x6f,
	0x6c, 0x61, 0x62, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x43, 0x6d, 0x64, 0x12, 0x3d, 0x0a, 0x03,
	0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x1d, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3f, 0x0a, 0x05, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a,
	0x14, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f,
	0x51, 0x55, 0x45, 0x52, 0x59, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x14, 0x0a, 0x05,
	0x73, 0x74, 0x61, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x73, 0x12, 0x4d, 0x0a, 0x0d, 0x61, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x74,
	0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e,
	0x45, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x54, 0x79, 0x70, 0x65, 0x3a, 0x12,
	0x45, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x4d,
	0x49, 0x4e, 0x52, 0x0d, 0x61, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x74, 0x79, 0x70,
	0x65, 0x22, 0xd2, 0x01, 0x0a, 0x18, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x41,
	0x63, 0x74, 0x69, 0x76, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x72, 0x43, 0x6d, 0x64, 0x12, 0x3d,
	0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x1d, 0x53, 0x43, 0x45, 0x4e, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x47, 0x0a,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43,
	0x6d, 0x64, 0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x3a, 0x1c, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41,
	0x4d, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x52, 0x52,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x18,
	0x03, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x22, 0xf5, 0x01, 0x0a, 0x16, 0x41, 0x73, 0x74, 0x72, 0x6f,
	0x6c, 0x61, 0x62, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x65, 0x74, 0x43, 0x6d,
	0x64, 0x12, 0x3d, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c,
	0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x1d, 0x53, 0x43,
	0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41,
	0x42, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64,
	0x12, 0x45, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1a, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x54,
	0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x27, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x41, 0x73, 0x74,
	0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x2c, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x16, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x43,
	0x6f, 0x73, 0x74, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xc3,
	0x01, 0x0a, 0x11, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x52, 0x65, 0x73, 0x65,
	0x74, 0x43, 0x6d, 0x64, 0x12, 0x3d, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a,
	0x1d, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52,
	0x4f, 0x4c, 0x41, 0x42, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03,
	0x63, 0x6d, 0x64, 0x12, 0x3f, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61,
	0x62, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x14, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41,
	0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x54, 0x52, 0x05, 0x70,
	0x61, 0x72, 0x61, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x18, 0x03, 0x20,
	0x03, 0x28, 0x0d, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x22, 0xb0, 0x01, 0x0a, 0x14, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61,
	0x62, 0x65, 0x50, 0x6c, 0x61, 0x6e, 0x53, 0x61, 0x76, 0x65, 0x43, 0x6d, 0x64, 0x12, 0x3d, 0x0a,
	0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x1d, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x43, 0x0a, 0x05,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x3a, 0x18, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d,
	0x5f, 0x50, 0x4c, 0x41, 0x4e, 0x5f, 0x53, 0x41, 0x56, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0d,
	0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x73, 0x2a, 0xa4, 0x01, 0x0a, 0x0e, 0x41, 0x73, 0x74, 0x72,
	0x6f, 0x6c, 0x61, 0x62, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x53,
	0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x51, 0x55, 0x45,
	0x52, 0x59, 0x10, 0x01, 0x12, 0x20, 0x0a, 0x1c, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42,
	0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x41, 0x54, 0x45, 0x5f,
	0x53, 0x54, 0x41, 0x52, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c,
	0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x5f, 0x52,
	0x45, 0x53, 0x45, 0x54, 0x10, 0x03, 0x12, 0x18, 0x0a, 0x14, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c,
	0x41, 0x42, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x52, 0x45, 0x53, 0x45, 0x54, 0x10, 0x04,
	0x12, 0x1c, 0x0a, 0x18, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x50, 0x41, 0x52,
	0x41, 0x4d, 0x5f, 0x50, 0x4c, 0x41, 0x4e, 0x5f, 0x53, 0x41, 0x56, 0x45, 0x10, 0x05, 0x2a, 0x78,
	0x0a, 0x0e, 0x45, 0x41, 0x73, 0x74, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x65, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x16, 0x0a, 0x12, 0x45, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x1d, 0x0a, 0x19, 0x45, 0x41, 0x53, 0x54,
	0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x46, 0x45,
	0x53, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x41, 0x53, 0x54, 0x52,
	0x4f, 0x4c, 0x41, 0x42, 0x45, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x50, 0x4c, 0x41, 0x4e, 0x10, 0x64,
	0x12, 0x16, 0x0a, 0x12, 0x45, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x54, 0x59,
	0x50, 0x45, 0x5f, 0x4d, 0x41, 0x58, 0x10, 0x65,
}

var (
	file_AstrolabeCmd_proto_rawDescOnce sync.Once
	file_AstrolabeCmd_proto_rawDescData = file_AstrolabeCmd_proto_rawDesc
)

func file_AstrolabeCmd_proto_rawDescGZIP() []byte {
	file_AstrolabeCmd_proto_rawDescOnce.Do(func() {
		file_AstrolabeCmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_AstrolabeCmd_proto_rawDescData)
	})
	return file_AstrolabeCmd_proto_rawDescData
}

var file_AstrolabeCmd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_AstrolabeCmd_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_AstrolabeCmd_proto_goTypes = []interface{}{
	(AstrolabeParam)(0),              // 0: Cmd.AstrolabeParam
	(EAstrolabeType)(0),              // 1: Cmd.EAstrolabeType
	(*AstrolabeCostData)(nil),        // 2: Cmd.AstrolabeCostData
	(*AstrolabeQueryCmd)(nil),        // 3: Cmd.AstrolabeQueryCmd
	(*AstrolabeActivateStarCmd)(nil), // 4: Cmd.AstrolabeActivateStarCmd
	(*AstrolabeQueryResetCmd)(nil),   // 5: Cmd.AstrolabeQueryResetCmd
	(*AstrolabeResetCmd)(nil),        // 6: Cmd.AstrolabeResetCmd
	(*AstrolabePlanSaveCmd)(nil),     // 7: Cmd.AstrolabePlanSaveCmd
	(Command)(0),                     // 8: Cmd.Command
}
var file_AstrolabeCmd_proto_depIdxs = []int32{
	8,  // 0: Cmd.AstrolabeQueryCmd.cmd:type_name -> Cmd.Command
	0,  // 1: Cmd.AstrolabeQueryCmd.param:type_name -> Cmd.AstrolabeParam
	1,  // 2: Cmd.AstrolabeQueryCmd.astrolabetype:type_name -> Cmd.EAstrolabeType
	8,  // 3: Cmd.AstrolabeActivateStarCmd.cmd:type_name -> Cmd.Command
	0,  // 4: Cmd.AstrolabeActivateStarCmd.param:type_name -> Cmd.AstrolabeParam
	8,  // 5: Cmd.AstrolabeQueryResetCmd.cmd:type_name -> Cmd.Command
	0,  // 6: Cmd.AstrolabeQueryResetCmd.param:type_name -> Cmd.AstrolabeParam
	1,  // 7: Cmd.AstrolabeQueryResetCmd.type:type_name -> Cmd.EAstrolabeType
	2,  // 8: Cmd.AstrolabeQueryResetCmd.items:type_name -> Cmd.AstrolabeCostData
	8,  // 9: Cmd.AstrolabeResetCmd.cmd:type_name -> Cmd.Command
	0,  // 10: Cmd.AstrolabeResetCmd.param:type_name -> Cmd.AstrolabeParam
	8,  // 11: Cmd.AstrolabePlanSaveCmd.cmd:type_name -> Cmd.Command
	0,  // 12: Cmd.AstrolabePlanSaveCmd.param:type_name -> Cmd.AstrolabeParam
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_AstrolabeCmd_proto_init() }
func file_AstrolabeCmd_proto_init() {
	if File_AstrolabeCmd_proto != nil {
		return
	}
	file_xCmd_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_AstrolabeCmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabeCostData); i {
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
		file_AstrolabeCmd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabeQueryCmd); i {
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
		file_AstrolabeCmd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabeActivateStarCmd); i {
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
		file_AstrolabeCmd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabeQueryResetCmd); i {
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
		file_AstrolabeCmd_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabeResetCmd); i {
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
		file_AstrolabeCmd_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AstrolabePlanSaveCmd); i {
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
			RawDescriptor: file_AstrolabeCmd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_AstrolabeCmd_proto_goTypes,
		DependencyIndexes: file_AstrolabeCmd_proto_depIdxs,
		EnumInfos:         file_AstrolabeCmd_proto_enumTypes,
		MessageInfos:      file_AstrolabeCmd_proto_msgTypes,
	}.Build()
	File_AstrolabeCmd_proto = out.File
	file_AstrolabeCmd_proto_rawDesc = nil
	file_AstrolabeCmd_proto_goTypes = nil
	file_AstrolabeCmd_proto_depIdxs = nil
}
