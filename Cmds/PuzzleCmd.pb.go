// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: PuzzleCmd.proto

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

type PuzzleParam int32

const (
	PuzzleParam_PUZZLEPARAM_QUERYACTLIST  PuzzleParam = 1
	PuzzleParam_PUZZLEPARAM_ITEMNTF       PuzzleParam = 3
	PuzzleParam_PUZZLEPARAM_ACTIVIEPUZZLE PuzzleParam = 4
)

// Enum value maps for PuzzleParam.
var (
	PuzzleParam_name = map[int32]string{
		1: "PUZZLEPARAM_QUERYACTLIST",
		3: "PUZZLEPARAM_ITEMNTF",
		4: "PUZZLEPARAM_ACTIVIEPUZZLE",
	}
	PuzzleParam_value = map[string]int32{
		"PUZZLEPARAM_QUERYACTLIST":  1,
		"PUZZLEPARAM_ITEMNTF":       3,
		"PUZZLEPARAM_ACTIVIEPUZZLE": 4,
	}
)

func (x PuzzleParam) Enum() *PuzzleParam {
	p := new(PuzzleParam)
	*p = x
	return p
}

func (x PuzzleParam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PuzzleParam) Descriptor() protoreflect.EnumDescriptor {
	return file_PuzzleCmd_proto_enumTypes[0].Descriptor()
}

func (PuzzleParam) Type() protoreflect.EnumType {
	return &file_PuzzleCmd_proto_enumTypes[0]
}

func (x PuzzleParam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *PuzzleParam) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = PuzzleParam(num)
	return nil
}

// Deprecated: Use PuzzleParam.Descriptor instead.
func (PuzzleParam) EnumDescriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{0}
}

type EPuzzleState int32

const (
	EPuzzleState_EPUZZLESTATE_MIN       EPuzzleState = 0
	EPuzzleState_EPUZZLESTATE_UNACTIVE  EPuzzleState = 1
	EPuzzleState_EPUZZLESTATE_CANACTIVE EPuzzleState = 2
	EPuzzleState_EPUZZLESTATE_ACTIVE    EPuzzleState = 3
	EPuzzleState_EPUZZLESTATE_MAX       EPuzzleState = 4
)

// Enum value maps for EPuzzleState.
var (
	EPuzzleState_name = map[int32]string{
		0: "EPUZZLESTATE_MIN",
		1: "EPUZZLESTATE_UNACTIVE",
		2: "EPUZZLESTATE_CANACTIVE",
		3: "EPUZZLESTATE_ACTIVE",
		4: "EPUZZLESTATE_MAX",
	}
	EPuzzleState_value = map[string]int32{
		"EPUZZLESTATE_MIN":       0,
		"EPUZZLESTATE_UNACTIVE":  1,
		"EPUZZLESTATE_CANACTIVE": 2,
		"EPUZZLESTATE_ACTIVE":    3,
		"EPUZZLESTATE_MAX":       4,
	}
)

func (x EPuzzleState) Enum() *EPuzzleState {
	p := new(EPuzzleState)
	*p = x
	return p
}

func (x EPuzzleState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EPuzzleState) Descriptor() protoreflect.EnumDescriptor {
	return file_PuzzleCmd_proto_enumTypes[1].Descriptor()
}

func (EPuzzleState) Type() protoreflect.EnumType {
	return &file_PuzzleCmd_proto_enumTypes[1]
}

func (x EPuzzleState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *EPuzzleState) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = EPuzzleState(num)
	return nil
}

// Deprecated: Use EPuzzleState.Descriptor instead.
func (EPuzzleState) EnumDescriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{1}
}

type PuzzleItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Actid   *uint32       `protobuf:"varint,1,opt,name=actid,def=0" json:"actid,omitempty"`
	Puzzled *uint32       `protobuf:"varint,2,opt,name=puzzled,def=0" json:"puzzled,omitempty"`
	Process *uint32       `protobuf:"varint,3,opt,name=process,def=0" json:"process,omitempty"`
	State   *EPuzzleState `protobuf:"varint,4,opt,name=state,enum=Cmd.EPuzzleState,def=0" json:"state,omitempty"`
}

// Default values for PuzzleItem fields.
const (
	Default_PuzzleItem_Actid   = uint32(0)
	Default_PuzzleItem_Puzzled = uint32(0)
	Default_PuzzleItem_Process = uint32(0)
	Default_PuzzleItem_State   = EPuzzleState_EPUZZLESTATE_MIN
)

func (x *PuzzleItem) Reset() {
	*x = PuzzleItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_PuzzleCmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PuzzleItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PuzzleItem) ProtoMessage() {}

func (x *PuzzleItem) ProtoReflect() protoreflect.Message {
	mi := &file_PuzzleCmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PuzzleItem.ProtoReflect.Descriptor instead.
func (*PuzzleItem) Descriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{0}
}

func (x *PuzzleItem) GetActid() uint32 {
	if x != nil && x.Actid != nil {
		return *x.Actid
	}
	return Default_PuzzleItem_Actid
}

func (x *PuzzleItem) GetPuzzled() uint32 {
	if x != nil && x.Puzzled != nil {
		return *x.Puzzled
	}
	return Default_PuzzleItem_Puzzled
}

func (x *PuzzleItem) GetProcess() uint32 {
	if x != nil && x.Process != nil {
		return *x.Process
	}
	return Default_PuzzleItem_Process
}

func (x *PuzzleItem) GetState() EPuzzleState {
	if x != nil && x.State != nil {
		return *x.State
	}
	return Default_PuzzleItem_State
}

type ActItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Actid *uint32       `protobuf:"varint,1,opt,name=actid,def=0" json:"actid,omitempty"`
	Items []*PuzzleItem `protobuf:"bytes,2,rep,name=items" json:"items,omitempty"`
}

// Default values for ActItem fields.
const (
	Default_ActItem_Actid = uint32(0)
)

func (x *ActItem) Reset() {
	*x = ActItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_PuzzleCmd_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActItem) ProtoMessage() {}

func (x *ActItem) ProtoReflect() protoreflect.Message {
	mi := &file_PuzzleCmd_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActItem.ProtoReflect.Descriptor instead.
func (*ActItem) Descriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{1}
}

func (x *ActItem) GetActid() uint32 {
	if x != nil && x.Actid != nil {
		return *x.Actid
	}
	return Default_ActItem_Actid
}

func (x *ActItem) GetItems() []*PuzzleItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type QueryActPuzzleCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd     *Command     `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=68" json:"cmd,omitempty"`
	Param   *PuzzleParam `protobuf:"varint,2,opt,name=param,enum=Cmd.PuzzleParam,def=1" json:"param,omitempty"`
	Actitem []*ActItem   `protobuf:"bytes,3,rep,name=actitem" json:"actitem,omitempty"`
}

// Default values for QueryActPuzzleCmd fields.
const (
	Default_QueryActPuzzleCmd_Cmd   = Command_PUZZLE_PROTOCMD
	Default_QueryActPuzzleCmd_Param = PuzzleParam_PUZZLEPARAM_QUERYACTLIST
)

func (x *QueryActPuzzleCmd) Reset() {
	*x = QueryActPuzzleCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_PuzzleCmd_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryActPuzzleCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryActPuzzleCmd) ProtoMessage() {}

func (x *QueryActPuzzleCmd) ProtoReflect() protoreflect.Message {
	mi := &file_PuzzleCmd_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryActPuzzleCmd.ProtoReflect.Descriptor instead.
func (*QueryActPuzzleCmd) Descriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{2}
}

func (x *QueryActPuzzleCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_QueryActPuzzleCmd_Cmd
}

func (x *QueryActPuzzleCmd) GetParam() PuzzleParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_QueryActPuzzleCmd_Param
}

func (x *QueryActPuzzleCmd) GetActitem() []*ActItem {
	if x != nil {
		return x.Actitem
	}
	return nil
}

type PuzzleItemNtf struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command      `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=68" json:"cmd,omitempty"`
	Param *PuzzleParam  `protobuf:"varint,2,opt,name=param,enum=Cmd.PuzzleParam,def=3" json:"param,omitempty"`
	Items []*PuzzleItem `protobuf:"bytes,3,rep,name=items" json:"items,omitempty"`
}

// Default values for PuzzleItemNtf fields.
const (
	Default_PuzzleItemNtf_Cmd   = Command_PUZZLE_PROTOCMD
	Default_PuzzleItemNtf_Param = PuzzleParam_PUZZLEPARAM_ITEMNTF
)

func (x *PuzzleItemNtf) Reset() {
	*x = PuzzleItemNtf{}
	if protoimpl.UnsafeEnabled {
		mi := &file_PuzzleCmd_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PuzzleItemNtf) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PuzzleItemNtf) ProtoMessage() {}

func (x *PuzzleItemNtf) ProtoReflect() protoreflect.Message {
	mi := &file_PuzzleCmd_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PuzzleItemNtf.ProtoReflect.Descriptor instead.
func (*PuzzleItemNtf) Descriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{3}
}

func (x *PuzzleItemNtf) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_PuzzleItemNtf_Cmd
}

func (x *PuzzleItemNtf) GetParam() PuzzleParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_PuzzleItemNtf_Param
}

func (x *PuzzleItemNtf) GetItems() []*PuzzleItem {
	if x != nil {
		return x.Items
	}
	return nil
}

type ActivePuzzleCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd      *Command     `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=68" json:"cmd,omitempty"`
	Param    *PuzzleParam `protobuf:"varint,2,opt,name=param,enum=Cmd.PuzzleParam,def=4" json:"param,omitempty"`
	Actid    *uint32      `protobuf:"varint,3,opt,name=actid,def=0" json:"actid,omitempty"`
	Puzzleid *uint32      `protobuf:"varint,4,opt,name=puzzleid,def=0" json:"puzzleid,omitempty"`
}

// Default values for ActivePuzzleCmd fields.
const (
	Default_ActivePuzzleCmd_Cmd      = Command_PUZZLE_PROTOCMD
	Default_ActivePuzzleCmd_Param    = PuzzleParam_PUZZLEPARAM_ACTIVIEPUZZLE
	Default_ActivePuzzleCmd_Actid    = uint32(0)
	Default_ActivePuzzleCmd_Puzzleid = uint32(0)
)

func (x *ActivePuzzleCmd) Reset() {
	*x = ActivePuzzleCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_PuzzleCmd_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivePuzzleCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivePuzzleCmd) ProtoMessage() {}

func (x *ActivePuzzleCmd) ProtoReflect() protoreflect.Message {
	mi := &file_PuzzleCmd_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivePuzzleCmd.ProtoReflect.Descriptor instead.
func (*ActivePuzzleCmd) Descriptor() ([]byte, []int) {
	return file_PuzzleCmd_proto_rawDescGZIP(), []int{4}
}

func (x *ActivePuzzleCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_ActivePuzzleCmd_Cmd
}

func (x *ActivePuzzleCmd) GetParam() PuzzleParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_ActivePuzzleCmd_Param
}

func (x *ActivePuzzleCmd) GetActid() uint32 {
	if x != nil && x.Actid != nil {
		return *x.Actid
	}
	return Default_ActivePuzzleCmd_Actid
}

func (x *ActivePuzzleCmd) GetPuzzleid() uint32 {
	if x != nil && x.Puzzleid != nil {
		return *x.Puzzleid
	}
	return Default_ActivePuzzleCmd_Puzzleid
}

var File_PuzzleCmd_proto protoreflect.FileDescriptor

var file_PuzzleCmd_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x03, 0x43, 0x6d, 0x64, 0x1a, 0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0x9a, 0x01, 0x0a, 0x0a, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x17, 0x0a, 0x05, 0x61, 0x63, 0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d,
	0x3a, 0x01, 0x30, 0x52, 0x05, 0x61, 0x63, 0x74, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x07, 0x70, 0x75,
	0x7a, 0x7a, 0x6c, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07,
	0x70, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x64, 0x12, 0x1b, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x63, 0x65,
	0x73, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07, 0x70, 0x72, 0x6f,
	0x63, 0x65, 0x73, 0x73, 0x12, 0x39, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x11, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x50, 0x75, 0x7a, 0x7a, 0x6c,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x3a, 0x10, 0x45, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22,
	0x49, 0x0a, 0x07, 0x41, 0x63, 0x74, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x17, 0x0a, 0x05, 0x61, 0x63,
	0x74, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x05, 0x61, 0x63,
	0x74, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0xae, 0x01, 0x0a, 0x11, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x41, 0x63, 0x74, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x43, 0x6d, 0x64,
	0x12, 0x2f, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e,
	0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x0f, 0x50, 0x55, 0x5a,
	0x5a, 0x4c, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d,
	0x64, 0x12, 0x40, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x10, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x3a, 0x18, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f,
	0x51, 0x55, 0x45, 0x52, 0x59, 0x41, 0x43, 0x54, 0x4c, 0x49, 0x53, 0x54, 0x52, 0x05, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x12, 0x26, 0x0a, 0x07, 0x61, 0x63, 0x74, 0x69, 0x74, 0x65, 0x6d, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x41, 0x63, 0x74, 0x49, 0x74,
	0x65, 0x6d, 0x52, 0x07, 0x61, 0x63, 0x74, 0x69, 0x74, 0x65, 0x6d, 0x22, 0xa4, 0x01, 0x0a, 0x0d,
	0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x4e, 0x74, 0x66, 0x12, 0x2f, 0x0a,
	0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x0f, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3b,
	0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e,
	0x43, 0x6d, 0x64, 0x2e, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a,
	0x13, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x49, 0x54, 0x45,
	0x4d, 0x4e, 0x54, 0x46, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x25, 0x0a, 0x05, 0x69,
	0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x22, 0xbd, 0x01, 0x0a, 0x0f, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65, 0x50, 0x75, 0x7a,
	0x7a, 0x6c, 0x65, 0x43, 0x6d, 0x64, 0x12, 0x2f, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x3a, 0x0f, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x41, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x50, 0x75, 0x7a,
	0x7a, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x19, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45,
	0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x49, 0x45, 0x50, 0x55, 0x5a,
	0x5a, 0x4c, 0x45, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x17, 0x0a, 0x05, 0x61, 0x63,
	0x74, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x05, 0x61, 0x63,
	0x74, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x08, 0x70, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x69, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x08, 0x70, 0x75, 0x7a, 0x7a, 0x6c, 0x65,
	0x69, 0x64, 0x2a, 0x63, 0x0a, 0x0b, 0x50, 0x75, 0x7a, 0x7a, 0x6c, 0x65, 0x50, 0x61, 0x72, 0x61,
	0x6d, 0x12, 0x1c, 0x0a, 0x18, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d,
	0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x41, 0x43, 0x54, 0x4c, 0x49, 0x53, 0x54, 0x10, 0x01, 0x12,
	0x17, 0x0a, 0x13, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x49,
	0x54, 0x45, 0x4d, 0x4e, 0x54, 0x46, 0x10, 0x03, 0x12, 0x1d, 0x0a, 0x19, 0x50, 0x55, 0x5a, 0x5a,
	0x4c, 0x45, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x49, 0x45, 0x50,
	0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x10, 0x04, 0x2a, 0x8a, 0x01, 0x0a, 0x0c, 0x45, 0x50, 0x75, 0x7a,
	0x7a, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x50, 0x55, 0x5a,
	0x5a, 0x4c, 0x45, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x10, 0x00, 0x12, 0x19,
	0x0a, 0x15, 0x45, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55,
	0x4e, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x45, 0x50, 0x55,
	0x5a, 0x5a, 0x4c, 0x45, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x43, 0x41, 0x4e, 0x41, 0x43, 0x54,
	0x49, 0x56, 0x45, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45,
	0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x56, 0x45, 0x10, 0x03, 0x12, 0x14,
	0x0a, 0x10, 0x45, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4d,
	0x41, 0x58, 0x10, 0x04,
}

var (
	file_PuzzleCmd_proto_rawDescOnce sync.Once
	file_PuzzleCmd_proto_rawDescData = file_PuzzleCmd_proto_rawDesc
)

func file_PuzzleCmd_proto_rawDescGZIP() []byte {
	file_PuzzleCmd_proto_rawDescOnce.Do(func() {
		file_PuzzleCmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_PuzzleCmd_proto_rawDescData)
	})
	return file_PuzzleCmd_proto_rawDescData
}

var file_PuzzleCmd_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_PuzzleCmd_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_PuzzleCmd_proto_goTypes = []interface{}{
	(PuzzleParam)(0),          // 0: Cmd.PuzzleParam
	(EPuzzleState)(0),         // 1: Cmd.EPuzzleState
	(*PuzzleItem)(nil),        // 2: Cmd.PuzzleItem
	(*ActItem)(nil),           // 3: Cmd.ActItem
	(*QueryActPuzzleCmd)(nil), // 4: Cmd.QueryActPuzzleCmd
	(*PuzzleItemNtf)(nil),     // 5: Cmd.PuzzleItemNtf
	(*ActivePuzzleCmd)(nil),   // 6: Cmd.ActivePuzzleCmd
	(Command)(0),              // 7: Cmd.Command
}
var file_PuzzleCmd_proto_depIdxs = []int32{
	1,  // 0: Cmd.PuzzleItem.state:type_name -> Cmd.EPuzzleState
	2,  // 1: Cmd.ActItem.items:type_name -> Cmd.PuzzleItem
	7,  // 2: Cmd.QueryActPuzzleCmd.cmd:type_name -> Cmd.Command
	0,  // 3: Cmd.QueryActPuzzleCmd.param:type_name -> Cmd.PuzzleParam
	3,  // 4: Cmd.QueryActPuzzleCmd.actitem:type_name -> Cmd.ActItem
	7,  // 5: Cmd.PuzzleItemNtf.cmd:type_name -> Cmd.Command
	0,  // 6: Cmd.PuzzleItemNtf.param:type_name -> Cmd.PuzzleParam
	2,  // 7: Cmd.PuzzleItemNtf.items:type_name -> Cmd.PuzzleItem
	7,  // 8: Cmd.ActivePuzzleCmd.cmd:type_name -> Cmd.Command
	0,  // 9: Cmd.ActivePuzzleCmd.param:type_name -> Cmd.PuzzleParam
	10, // [10:10] is the sub-list for method output_type
	10, // [10:10] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_PuzzleCmd_proto_init() }
func file_PuzzleCmd_proto_init() {
	if File_PuzzleCmd_proto != nil {
		return
	}
	file_xCmd_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_PuzzleCmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PuzzleItem); i {
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
		file_PuzzleCmd_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActItem); i {
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
		file_PuzzleCmd_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryActPuzzleCmd); i {
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
		file_PuzzleCmd_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PuzzleItemNtf); i {
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
		file_PuzzleCmd_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivePuzzleCmd); i {
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
			RawDescriptor: file_PuzzleCmd_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_PuzzleCmd_proto_goTypes,
		DependencyIndexes: file_PuzzleCmd_proto_depIdxs,
		EnumInfos:         file_PuzzleCmd_proto_enumTypes,
		MessageInfos:      file_PuzzleCmd_proto_msgTypes,
	}.Build()
	File_PuzzleCmd_proto = out.File
	file_PuzzleCmd_proto_rawDesc = nil
	file_PuzzleCmd_proto_goTypes = nil
	file_PuzzleCmd_proto_depIdxs = nil
}