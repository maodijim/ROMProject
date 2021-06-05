// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: TechTreeCmd.proto

package Cmd

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

type TechTreeParam int32

const (
	TechTreeParam_TECHTREEPARAM_UNLOCK        TechTreeParam = 1
	TechTreeParam_TECHTREEPARAM_SYNCINFO      TechTreeParam = 2
	TechTreeParam_TECHTREEPARAM_TOY_UNLOCK    TechTreeParam = 3
	TechTreeParam_TECHTREEPARAM_TOY_SYNCINFO  TechTreeParam = 4
	TechTreeParam_TECHTREEPARAM_MAKE_TOY      TechTreeParam = 5
	TechTreeParam_TECHTREEPARAM_TOY_TRANS_POS TechTreeParam = 6
)

// Enum value maps for TechTreeParam.
var (
	TechTreeParam_name = map[int32]string{
		1: "TECHTREEPARAM_UNLOCK",
		2: "TECHTREEPARAM_SYNCINFO",
		3: "TECHTREEPARAM_TOY_UNLOCK",
		4: "TECHTREEPARAM_TOY_SYNCINFO",
		5: "TECHTREEPARAM_MAKE_TOY",
		6: "TECHTREEPARAM_TOY_TRANS_POS",
	}
	TechTreeParam_value = map[string]int32{
		"TECHTREEPARAM_UNLOCK":        1,
		"TECHTREEPARAM_SYNCINFO":      2,
		"TECHTREEPARAM_TOY_UNLOCK":    3,
		"TECHTREEPARAM_TOY_SYNCINFO":  4,
		"TECHTREEPARAM_MAKE_TOY":      5,
		"TECHTREEPARAM_TOY_TRANS_POS": 6,
	}
)

func (x TechTreeParam) Enum() *TechTreeParam {
	p := new(TechTreeParam)
	*p = x
	return p
}

func (x TechTreeParam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TechTreeParam) Descriptor() protoreflect.EnumDescriptor {
	return file_TechTreeCmd_proto_enumTypes[0].Descriptor()
}

func (TechTreeParam) Type() protoreflect.EnumType {
	return &file_TechTreeCmd_proto_enumTypes[0]
}

func (x TechTreeParam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *TechTreeParam) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = TechTreeParam(num)
	return nil
}

// Deprecated: Use TechTreeParam.Descriptor instead.
func (TechTreeParam) EnumDescriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{0}
}

type TechTreeLeafInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Leafid *uint32 `protobuf:"varint,1,opt,name=leafid,def=0" json:"leafid,omitempty"`
	Level  *uint32 `protobuf:"varint,2,opt,name=level,def=0" json:"level,omitempty"`
}

// Default values for TechTreeLeafInfo fields.
const (
	Default_TechTreeLeafInfo_Leafid = uint32(0)
	Default_TechTreeLeafInfo_Level  = uint32(0)
)

func (x *TechTreeLeafInfo) Reset() {
	*x = TechTreeLeafInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TechTreeLeafInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TechTreeLeafInfo) ProtoMessage() {}

func (x *TechTreeLeafInfo) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TechTreeLeafInfo.ProtoReflect.Descriptor instead.
func (*TechTreeLeafInfo) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{0}
}

func (x *TechTreeLeafInfo) GetLeafid() uint32 {
	if x != nil && x.Leafid != nil {
		return *x.Leafid
	}
	return Default_TechTreeLeafInfo_Leafid
}

func (x *TechTreeLeafInfo) GetLevel() uint32 {
	if x != nil && x.Level != nil {
		return *x.Level
	}
	return Default_TechTreeLeafInfo_Level
}

type TechTreeUnlockLeafCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd    *Command          `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param  *TechTreeParam    `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=1" json:"param,omitempty"`
	Leaf   *TechTreeLeafInfo `protobuf:"bytes,4,opt,name=leaf" json:"leaf,omitempty"`
	Treeid *uint32           `protobuf:"varint,5,opt,name=treeid" json:"treeid,omitempty"`
}

// Default values for TechTreeUnlockLeafCmd fields.
const (
	Default_TechTreeUnlockLeafCmd_Cmd   = Command_TECHTREE_PROTOCMD
	Default_TechTreeUnlockLeafCmd_Param = TechTreeParam_TECHTREEPARAM_UNLOCK
)

func (x *TechTreeUnlockLeafCmd) Reset() {
	*x = TechTreeUnlockLeafCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TechTreeUnlockLeafCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TechTreeUnlockLeafCmd) ProtoMessage() {}

func (x *TechTreeUnlockLeafCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TechTreeUnlockLeafCmd.ProtoReflect.Descriptor instead.
func (*TechTreeUnlockLeafCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{1}
}

func (x *TechTreeUnlockLeafCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_TechTreeUnlockLeafCmd_Cmd
}

func (x *TechTreeUnlockLeafCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_TechTreeUnlockLeafCmd_Param
}

func (x *TechTreeUnlockLeafCmd) GetLeaf() *TechTreeLeafInfo {
	if x != nil {
		return x.Leaf
	}
	return nil
}

func (x *TechTreeUnlockLeafCmd) GetTreeid() uint32 {
	if x != nil && x.Treeid != nil {
		return *x.Treeid
	}
	return 0
}

type TechTreeSyncLeafCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd    *Command            `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param  *TechTreeParam      `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=2" json:"param,omitempty"`
	Leaves []*TechTreeLeafInfo `protobuf:"bytes,3,rep,name=leaves" json:"leaves,omitempty"`
}

// Default values for TechTreeSyncLeafCmd fields.
const (
	Default_TechTreeSyncLeafCmd_Cmd   = Command_TECHTREE_PROTOCMD
	Default_TechTreeSyncLeafCmd_Param = TechTreeParam_TECHTREEPARAM_SYNCINFO
)

func (x *TechTreeSyncLeafCmd) Reset() {
	*x = TechTreeSyncLeafCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TechTreeSyncLeafCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TechTreeSyncLeafCmd) ProtoMessage() {}

func (x *TechTreeSyncLeafCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TechTreeSyncLeafCmd.ProtoReflect.Descriptor instead.
func (*TechTreeSyncLeafCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{2}
}

func (x *TechTreeSyncLeafCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_TechTreeSyncLeafCmd_Cmd
}

func (x *TechTreeSyncLeafCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_TechTreeSyncLeafCmd_Param
}

func (x *TechTreeSyncLeafCmd) GetLeaves() []*TechTreeLeafInfo {
	if x != nil {
		return x.Leaves
	}
	return nil
}

type AddToyDrawingCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd       *Command       `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param     *TechTreeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=3" json:"param,omitempty"`
	Drawingid *uint32        `protobuf:"varint,3,opt,name=drawingid,def=0" json:"drawingid,omitempty"`
}

// Default values for AddToyDrawingCmd fields.
const (
	Default_AddToyDrawingCmd_Cmd       = Command_TECHTREE_PROTOCMD
	Default_AddToyDrawingCmd_Param     = TechTreeParam_TECHTREEPARAM_TOY_UNLOCK
	Default_AddToyDrawingCmd_Drawingid = uint32(0)
)

func (x *AddToyDrawingCmd) Reset() {
	*x = AddToyDrawingCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddToyDrawingCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddToyDrawingCmd) ProtoMessage() {}

func (x *AddToyDrawingCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddToyDrawingCmd.ProtoReflect.Descriptor instead.
func (*AddToyDrawingCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{3}
}

func (x *AddToyDrawingCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_AddToyDrawingCmd_Cmd
}

func (x *AddToyDrawingCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_AddToyDrawingCmd_Param
}

func (x *AddToyDrawingCmd) GetDrawingid() uint32 {
	if x != nil && x.Drawingid != nil {
		return *x.Drawingid
	}
	return Default_AddToyDrawingCmd_Drawingid
}

type SyncToyDrawingCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd      *Command       `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param    *TechTreeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=4" json:"param,omitempty"`
	Drawings []uint32       `protobuf:"varint,4,rep,name=drawings" json:"drawings,omitempty"`
}

// Default values for SyncToyDrawingCmd fields.
const (
	Default_SyncToyDrawingCmd_Cmd   = Command_TECHTREE_PROTOCMD
	Default_SyncToyDrawingCmd_Param = TechTreeParam_TECHTREEPARAM_TOY_SYNCINFO
)

func (x *SyncToyDrawingCmd) Reset() {
	*x = SyncToyDrawingCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SyncToyDrawingCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SyncToyDrawingCmd) ProtoMessage() {}

func (x *SyncToyDrawingCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SyncToyDrawingCmd.ProtoReflect.Descriptor instead.
func (*SyncToyDrawingCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{4}
}

func (x *SyncToyDrawingCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_SyncToyDrawingCmd_Cmd
}

func (x *SyncToyDrawingCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_SyncToyDrawingCmd_Param
}

func (x *SyncToyDrawingCmd) GetDrawings() []uint32 {
	if x != nil {
		return x.Drawings
	}
	return nil
}

type TechTreeMakeToyCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd       *Command       `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param     *TechTreeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=5" json:"param,omitempty"`
	Drawingid *uint32        `protobuf:"varint,3,opt,name=drawingid,def=0" json:"drawingid,omitempty"`
	Count     *uint32        `protobuf:"varint,4,opt,name=count,def=0" json:"count,omitempty"`
}

// Default values for TechTreeMakeToyCmd fields.
const (
	Default_TechTreeMakeToyCmd_Cmd       = Command_TECHTREE_PROTOCMD
	Default_TechTreeMakeToyCmd_Param     = TechTreeParam_TECHTREEPARAM_MAKE_TOY
	Default_TechTreeMakeToyCmd_Drawingid = uint32(0)
	Default_TechTreeMakeToyCmd_Count     = uint32(0)
)

func (x *TechTreeMakeToyCmd) Reset() {
	*x = TechTreeMakeToyCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TechTreeMakeToyCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TechTreeMakeToyCmd) ProtoMessage() {}

func (x *TechTreeMakeToyCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TechTreeMakeToyCmd.ProtoReflect.Descriptor instead.
func (*TechTreeMakeToyCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{5}
}

func (x *TechTreeMakeToyCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_TechTreeMakeToyCmd_Cmd
}

func (x *TechTreeMakeToyCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_TechTreeMakeToyCmd_Param
}

func (x *TechTreeMakeToyCmd) GetDrawingid() uint32 {
	if x != nil && x.Drawingid != nil {
		return *x.Drawingid
	}
	return Default_TechTreeMakeToyCmd_Drawingid
}

func (x *TechTreeMakeToyCmd) GetCount() uint32 {
	if x != nil && x.Count != nil {
		return *x.Count
	}
	return Default_TechTreeMakeToyCmd_Count
}

type ToyTransSetPosCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command       `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=73" json:"cmd,omitempty"`
	Param *TechTreeParam `protobuf:"varint,2,opt,name=param,enum=Cmd.TechTreeParam,def=6" json:"param,omitempty"`
	Pos   *ScenePos      `protobuf:"bytes,3,opt,name=pos" json:"pos,omitempty"`
}

// Default values for ToyTransSetPosCmd fields.
const (
	Default_ToyTransSetPosCmd_Cmd   = Command_TECHTREE_PROTOCMD
	Default_ToyTransSetPosCmd_Param = TechTreeParam_TECHTREEPARAM_TOY_TRANS_POS
)

func (x *ToyTransSetPosCmd) Reset() {
	*x = ToyTransSetPosCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_TechTreeCmd_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ToyTransSetPosCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ToyTransSetPosCmd) ProtoMessage() {}

func (x *ToyTransSetPosCmd) ProtoReflect() protoreflect.Message {
	mi := &file_TechTreeCmd_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ToyTransSetPosCmd.ProtoReflect.Descriptor instead.
func (*ToyTransSetPosCmd) Descriptor() ([]byte, []int) {
	return file_TechTreeCmd_proto_rawDescGZIP(), []int{6}
}

func (x *ToyTransSetPosCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_ToyTransSetPosCmd_Cmd
}

func (x *ToyTransSetPosCmd) GetParam() TechTreeParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_ToyTransSetPosCmd_Param
}

func (x *ToyTransSetPosCmd) GetPos() *ScenePos {
	if x != nil {
		return x.Pos
	}
	return nil
}

var File_TechTreeCmd_proto protoreflect.FileDescriptor

var file_TechTreeCmd_proto_rawDesc = []byte{
	0x0a, 0x11, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x43, 0x6d, 0x64, 0x1a, 0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x10, 0x54, 0x65, 0x63, 0x68, 0x54,
	0x72, 0x65, 0x65, 0x4c, 0x65, 0x61, 0x66, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x19, 0x0a, 0x06, 0x6c,
	0x65, 0x61, 0x66, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x06,
	0x6c, 0x65, 0x61, 0x66, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x05, 0x6c, 0x65, 0x76, 0x65, 0x6c, 0x22,
	0xcd, 0x01, 0x0a, 0x15, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x55, 0x6e, 0x6c, 0x6f,
	0x63, 0x6b, 0x4c, 0x65, 0x61, 0x66, 0x43, 0x6d, 0x64, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3e, 0x0a, 0x05,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a,
	0x14, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x55,
	0x4e, 0x4c, 0x4f, 0x43, 0x4b, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x29, 0x0a, 0x04,
	0x6c, 0x65, 0x61, 0x66, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x4c, 0x65, 0x61, 0x66, 0x49, 0x6e, 0x66,
	0x6f, 0x52, 0x04, 0x6c, 0x65, 0x61, 0x66, 0x12, 0x16, 0x0a, 0x06, 0x74, 0x72, 0x65, 0x65, 0x69,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x74, 0x72, 0x65, 0x65, 0x69, 0x64, 0x22,
	0xb9, 0x01, 0x0a, 0x13, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x53, 0x79, 0x6e, 0x63,
	0x4c, 0x65, 0x61, 0x66, 0x43, 0x6d, 0x64, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x40, 0x0a, 0x05, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x43, 0x6d, 0x64, 0x2e,
	0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x16, 0x54,
	0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x53, 0x59, 0x4e,
	0x43, 0x49, 0x4e, 0x46, 0x4f, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x2d, 0x0a, 0x06,
	0x6c, 0x65, 0x61, 0x76, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x43,
	0x6d, 0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x4c, 0x65, 0x61, 0x66, 0x49,
	0x6e, 0x66, 0x6f, 0x52, 0x06, 0x6c, 0x65, 0x61, 0x76, 0x65, 0x73, 0x22, 0xaa, 0x01, 0x0a, 0x10,
	0x41, 0x64, 0x64, 0x54, 0x6f, 0x79, 0x44, 0x72, 0x61, 0x77, 0x69, 0x6e, 0x67, 0x43, 0x6d, 0x64,
	0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e,
	0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43,
	0x48, 0x54, 0x52, 0x45, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03,
	0x63, 0x6d, 0x64, 0x12, 0x42, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x12, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65,
	0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x18, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45,
	0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x4f, 0x59, 0x5f, 0x55, 0x4e, 0x4c, 0x4f, 0x43, 0x4b,
	0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x09, 0x64, 0x72, 0x61, 0x77, 0x69,
	0x6e, 0x67, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x09, 0x64,
	0x72, 0x61, 0x77, 0x69, 0x6e, 0x67, 0x69, 0x64, 0x22, 0xa8, 0x01, 0x0a, 0x11, 0x53, 0x79, 0x6e,
	0x63, 0x54, 0x6f, 0x79, 0x44, 0x72, 0x61, 0x77, 0x69, 0x6e, 0x67, 0x43, 0x6d, 0x64, 0x12, 0x31,
	0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43, 0x48, 0x54,
	0x52, 0x45, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d,
	0x64, 0x12, 0x44, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x12, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x50,
	0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1a, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41,
	0x52, 0x41, 0x4d, 0x5f, 0x54, 0x4f, 0x59, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x49, 0x4e, 0x46, 0x4f,
	0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x72, 0x61, 0x77, 0x69,
	0x6e, 0x67, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x64, 0x72, 0x61, 0x77, 0x69,
	0x6e, 0x67, 0x73, 0x22, 0xc3, 0x01, 0x0a, 0x12, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65,
	0x4d, 0x61, 0x6b, 0x65, 0x54, 0x6f, 0x79, 0x43, 0x6d, 0x64, 0x12, 0x31, 0x0a, 0x03, 0x63, 0x6d,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x40, 0x0a,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x43,
	0x6d, 0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x3a, 0x16, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f,
	0x4d, 0x41, 0x4b, 0x45, 0x5f, 0x54, 0x4f, 0x59, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12,
	0x1f, 0x0a, 0x09, 0x64, 0x72, 0x61, 0x77, 0x69, 0x6e, 0x67, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x09, 0x64, 0x72, 0x61, 0x77, 0x69, 0x6e, 0x67, 0x69, 0x64,
	0x12, 0x17, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x3a,
	0x01, 0x30, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xae, 0x01, 0x0a, 0x11, 0x54, 0x6f,
	0x79, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x53, 0x65, 0x74, 0x50, 0x6f, 0x73, 0x43, 0x6d, 0x64, 0x12,
	0x31, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43,
	0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x11, 0x54, 0x45, 0x43, 0x48,
	0x54, 0x52, 0x45, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63,
	0x6d, 0x64, 0x12, 0x45, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x12, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x54, 0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1b, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x4f, 0x59, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x5f, 0x50,
	0x4f, 0x53, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x03, 0x70, 0x6f, 0x73,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x53, 0x63, 0x65,
	0x6e, 0x65, 0x50, 0x6f, 0x73, 0x52, 0x03, 0x70, 0x6f, 0x73, 0x2a, 0xc0, 0x01, 0x0a, 0x0d, 0x54,
	0x65, 0x63, 0x68, 0x54, 0x72, 0x65, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x18, 0x0a, 0x14,
	0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x55, 0x4e,
	0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52,
	0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x49, 0x4e, 0x46, 0x4f,
	0x10, 0x02, 0x12, 0x1c, 0x0a, 0x18, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41,
	0x52, 0x41, 0x4d, 0x5f, 0x54, 0x4f, 0x59, 0x5f, 0x55, 0x4e, 0x4c, 0x4f, 0x43, 0x4b, 0x10, 0x03,
	0x12, 0x1e, 0x0a, 0x1a, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41,
	0x4d, 0x5f, 0x54, 0x4f, 0x59, 0x5f, 0x53, 0x59, 0x4e, 0x43, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x04,
	0x12, 0x1a, 0x0a, 0x16, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41,
	0x4d, 0x5f, 0x4d, 0x41, 0x4b, 0x45, 0x5f, 0x54, 0x4f, 0x59, 0x10, 0x05, 0x12, 0x1f, 0x0a, 0x1b,
	0x54, 0x45, 0x43, 0x48, 0x54, 0x52, 0x45, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x54, 0x4f,
	0x59, 0x5f, 0x54, 0x52, 0x41, 0x4e, 0x53, 0x5f, 0x50, 0x4f, 0x53, 0x10, 0x06,
}

var (
	file_TechTreeCmd_proto_rawDescOnce sync.Once
	file_TechTreeCmd_proto_rawDescData = file_TechTreeCmd_proto_rawDesc
)

func file_TechTreeCmd_proto_rawDescGZIP() []byte {
	file_TechTreeCmd_proto_rawDescOnce.Do(func() {
		file_TechTreeCmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_TechTreeCmd_proto_rawDescData)
	})
	return file_TechTreeCmd_proto_rawDescData
}

var file_TechTreeCmd_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_TechTreeCmd_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_TechTreeCmd_proto_goTypes = []interface{}{
	(TechTreeParam)(0),            // 0: Cmd.TechTreeParam
	(*TechTreeLeafInfo)(nil),      // 1: Cmd.TechTreeLeafInfo
	(*TechTreeUnlockLeafCmd)(nil), // 2: Cmd.TechTreeUnlockLeafCmd
	(*TechTreeSyncLeafCmd)(nil),   // 3: Cmd.TechTreeSyncLeafCmd
	(*AddToyDrawingCmd)(nil),      // 4: Cmd.AddToyDrawingCmd
	(*SyncToyDrawingCmd)(nil),     // 5: Cmd.SyncToyDrawingCmd
	(*TechTreeMakeToyCmd)(nil),    // 6: Cmd.TechTreeMakeToyCmd
	(*ToyTransSetPosCmd)(nil),     // 7: Cmd.ToyTransSetPosCmd
	(Command)(0),                  // 8: Cmd.Command
	(*ScenePos)(nil),              // 9: Cmd.ScenePos
}
var file_TechTreeCmd_proto_depIdxs = []int32{
	8,  // 0: Cmd.TechTreeUnlockLeafCmd.cmd:type_name -> Cmd.Command
	0,  // 1: Cmd.TechTreeUnlockLeafCmd.param:type_name -> Cmd.TechTreeParam
	1,  // 2: Cmd.TechTreeUnlockLeafCmd.leaf:type_name -> Cmd.TechTreeLeafInfo
	8,  // 3: Cmd.TechTreeSyncLeafCmd.cmd:type_name -> Cmd.Command
	0,  // 4: Cmd.TechTreeSyncLeafCmd.param:type_name -> Cmd.TechTreeParam
	1,  // 5: Cmd.TechTreeSyncLeafCmd.leaves:type_name -> Cmd.TechTreeLeafInfo
	8,  // 6: Cmd.AddToyDrawingCmd.cmd:type_name -> Cmd.Command
	0,  // 7: Cmd.AddToyDrawingCmd.param:type_name -> Cmd.TechTreeParam
	8,  // 8: Cmd.SyncToyDrawingCmd.cmd:type_name -> Cmd.Command
	0,  // 9: Cmd.SyncToyDrawingCmd.param:type_name -> Cmd.TechTreeParam
	8,  // 10: Cmd.TechTreeMakeToyCmd.cmd:type_name -> Cmd.Command
	0,  // 11: Cmd.TechTreeMakeToyCmd.param:type_name -> Cmd.TechTreeParam
	8,  // 12: Cmd.ToyTransSetPosCmd.cmd:type_name -> Cmd.Command
	0,  // 13: Cmd.ToyTransSetPosCmd.param:type_name -> Cmd.TechTreeParam
	9,  // 14: Cmd.ToyTransSetPosCmd.pos:type_name -> Cmd.ScenePos
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_TechTreeCmd_proto_init() }
func file_TechTreeCmd_proto_init() {
	if File_TechTreeCmd_proto != nil {
		return
	}
	file_xCmd_proto_init()
	file_ProtoCommon_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_TechTreeCmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TechTreeLeafInfo); i {
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
		file_TechTreeCmd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TechTreeUnlockLeafCmd); i {
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
		file_TechTreeCmd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TechTreeSyncLeafCmd); i {
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
		file_TechTreeCmd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddToyDrawingCmd); i {
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
		file_TechTreeCmd_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SyncToyDrawingCmd); i {
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
		file_TechTreeCmd_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TechTreeMakeToyCmd); i {
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
		file_TechTreeCmd_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ToyTransSetPosCmd); i {
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
			RawDescriptor: file_TechTreeCmd_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_TechTreeCmd_proto_goTypes,
		DependencyIndexes: file_TechTreeCmd_proto_depIdxs,
		EnumInfos:         file_TechTreeCmd_proto_enumTypes,
		MessageInfos:      file_TechTreeCmd_proto_msgTypes,
	}.Build()
	File_TechTreeCmd_proto = out.File
	file_TechTreeCmd_proto_rawDesc = nil
	file_TechTreeCmd_proto_goTypes = nil
	file_TechTreeCmd_proto_depIdxs = nil
}
