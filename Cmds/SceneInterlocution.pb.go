// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: SceneInterlocution.proto

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

type InterParam int32

const (
	InterParam_INTERPARAM_NEWINTERLOCUTION InterParam = 1
	InterParam_INTERPARAM_ANSWERINTER      InterParam = 2
	InterParam_INTERPARAM_QUERYINTER       InterParam = 3
	InterParam_INTERPARAM_QUERYPAPER       InterParam = 4
	InterParam_INTERPARAM_PAPERQUESTION    InterParam = 5
	InterParam_INTERPARAM_PAPERRESULT      InterParam = 6
)

// Enum value maps for InterParam.
var (
	InterParam_name = map[int32]string{
		1: "INTERPARAM_NEWINTERLOCUTION",
		2: "INTERPARAM_ANSWERINTER",
		3: "INTERPARAM_QUERYINTER",
		4: "INTERPARAM_QUERYPAPER",
		5: "INTERPARAM_PAPERQUESTION",
		6: "INTERPARAM_PAPERRESULT",
	}
	InterParam_value = map[string]int32{
		"INTERPARAM_NEWINTERLOCUTION": 1,
		"INTERPARAM_ANSWERINTER":      2,
		"INTERPARAM_QUERYINTER":       3,
		"INTERPARAM_QUERYPAPER":       4,
		"INTERPARAM_PAPERQUESTION":    5,
		"INTERPARAM_PAPERRESULT":      6,
	}
)

func (x InterParam) Enum() *InterParam {
	p := new(InterParam)
	*p = x
	return p
}

func (x InterParam) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (InterParam) Descriptor() protoreflect.EnumDescriptor {
	return file_SceneInterlocution_proto_enumTypes[0].Descriptor()
}

func (InterParam) Type() protoreflect.EnumType {
	return &file_SceneInterlocution_proto_enumTypes[0]
}

func (x InterParam) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *InterParam) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = InterParam(num)
	return nil
}

// Deprecated: Use InterParam.Descriptor instead.
func (InterParam) EnumDescriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{0}
}

type EQueryState int32

const (
	EQueryState_EQUERYSTATE_OK             EQueryState = 1
	EQueryState_EQUERYSTATE_ANSWERED_RIGHT EQueryState = 2
	EQueryState_EQUERYSTATE_ANSWERED_WRONG EQueryState = 3
	EQueryState_EQUERYSTATE_FAIL           EQueryState = 4
)

// Enum value maps for EQueryState.
var (
	EQueryState_name = map[int32]string{
		1: "EQUERYSTATE_OK",
		2: "EQUERYSTATE_ANSWERED_RIGHT",
		3: "EQUERYSTATE_ANSWERED_WRONG",
		4: "EQUERYSTATE_FAIL",
	}
	EQueryState_value = map[string]int32{
		"EQUERYSTATE_OK":             1,
		"EQUERYSTATE_ANSWERED_RIGHT": 2,
		"EQUERYSTATE_ANSWERED_WRONG": 3,
		"EQUERYSTATE_FAIL":           4,
	}
)

func (x EQueryState) Enum() *EQueryState {
	p := new(EQueryState)
	*p = x
	return p
}

func (x EQueryState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EQueryState) Descriptor() protoreflect.EnumDescriptor {
	return file_SceneInterlocution_proto_enumTypes[1].Descriptor()
}

func (EQueryState) Type() protoreflect.EnumType {
	return &file_SceneInterlocution_proto_enumTypes[1]
}

func (x EQueryState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *EQueryState) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = EQueryState(num)
	return nil
}

// Deprecated: Use EQueryState.Descriptor instead.
func (EQueryState) EnumDescriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{1}
}

type InterData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Guid    *uint32  `protobuf:"varint,1,opt,name=guid,def=0" json:"guid,omitempty"`
	Interid *uint32  `protobuf:"varint,2,opt,name=interid,def=0" json:"interid,omitempty"`
	Paramid *uint32  `protobuf:"varint,4,opt,name=paramid,def=0" json:"paramid,omitempty"`
	Source  *ESource `protobuf:"varint,3,opt,name=source,enum=Cmd.ESource,def=0" json:"source,omitempty"`
}

// Default values for InterData fields.
const (
	Default_InterData_Guid    = uint32(0)
	Default_InterData_Interid = uint32(0)
	Default_InterData_Paramid = uint32(0)
	Default_InterData_Source  = ESource_ESOURCE_MIN
)

func (x *InterData) Reset() {
	*x = InterData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InterData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InterData) ProtoMessage() {}

func (x *InterData) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InterData.ProtoReflect.Descriptor instead.
func (*InterData) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{0}
}

func (x *InterData) GetGuid() uint32 {
	if x != nil && x.Guid != nil {
		return *x.Guid
	}
	return Default_InterData_Guid
}

func (x *InterData) GetInterid() uint32 {
	if x != nil && x.Interid != nil {
		return *x.Interid
	}
	return Default_InterData_Interid
}

func (x *InterData) GetParamid() uint32 {
	if x != nil && x.Paramid != nil {
		return *x.Paramid
	}
	return Default_InterData_Paramid
}

func (x *InterData) GetSource() ESource {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return Default_InterData_Source
}

type NewInter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd      *Command    `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param    *InterParam `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=1" json:"param,omitempty"`
	Inter    *InterData  `protobuf:"bytes,3,opt,name=inter" json:"inter,omitempty"`
	Npcid    *uint64     `protobuf:"varint,4,opt,name=npcid" json:"npcid,omitempty"`
	Answerid *uint64     `protobuf:"varint,5,opt,name=answerid,def=0" json:"answerid,omitempty"`
}

// Default values for NewInter fields.
const (
	Default_NewInter_Cmd      = Command_SCENE_USER_INTER_PROTOCMD
	Default_NewInter_Param    = InterParam_INTERPARAM_NEWINTERLOCUTION
	Default_NewInter_Answerid = uint64(0)
)

func (x *NewInter) Reset() {
	*x = NewInter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NewInter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NewInter) ProtoMessage() {}

func (x *NewInter) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NewInter.ProtoReflect.Descriptor instead.
func (*NewInter) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{1}
}

func (x *NewInter) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_NewInter_Cmd
}

func (x *NewInter) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_NewInter_Param
}

func (x *NewInter) GetInter() *InterData {
	if x != nil {
		return x.Inter
	}
	return nil
}

func (x *NewInter) GetNpcid() uint64 {
	if x != nil && x.Npcid != nil {
		return *x.Npcid
	}
	return 0
}

func (x *NewInter) GetAnswerid() uint64 {
	if x != nil && x.Answerid != nil {
		return *x.Answerid
	}
	return Default_NewInter_Answerid
}

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd     *Command    `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param   *InterParam `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=2" json:"param,omitempty"`
	Npcid   *uint64     `protobuf:"varint,3,opt,name=npcid" json:"npcid,omitempty"`
	Guid    *uint32     `protobuf:"varint,4,opt,name=guid,def=0" json:"guid,omitempty"`
	Interid *uint32     `protobuf:"varint,5,opt,name=interid,def=0" json:"interid,omitempty"`
	Source  *ESource    `protobuf:"varint,6,opt,name=source,enum=Cmd.ESource,def=0" json:"source,omitempty"`
	Answer  *uint32     `protobuf:"varint,7,opt,name=answer,def=0" json:"answer,omitempty"`
	Correct *bool       `protobuf:"varint,8,opt,name=correct,def=0" json:"correct,omitempty"`
	Paramid *uint32     `protobuf:"varint,9,opt,name=paramid,def=0" json:"paramid,omitempty"`
}

// Default values for Answer fields.
const (
	Default_Answer_Cmd     = Command_SCENE_USER_INTER_PROTOCMD
	Default_Answer_Param   = InterParam_INTERPARAM_ANSWERINTER
	Default_Answer_Guid    = uint32(0)
	Default_Answer_Interid = uint32(0)
	Default_Answer_Source  = ESource_ESOURCE_MIN
	Default_Answer_Answer  = uint32(0)
	Default_Answer_Correct = bool(false)
	Default_Answer_Paramid = uint32(0)
)

func (x *Answer) Reset() {
	*x = Answer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{2}
}

func (x *Answer) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_Answer_Cmd
}

func (x *Answer) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_Answer_Param
}

func (x *Answer) GetNpcid() uint64 {
	if x != nil && x.Npcid != nil {
		return *x.Npcid
	}
	return 0
}

func (x *Answer) GetGuid() uint32 {
	if x != nil && x.Guid != nil {
		return *x.Guid
	}
	return Default_Answer_Guid
}

func (x *Answer) GetInterid() uint32 {
	if x != nil && x.Interid != nil {
		return *x.Interid
	}
	return Default_Answer_Interid
}

func (x *Answer) GetSource() ESource {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return Default_Answer_Source
}

func (x *Answer) GetAnswer() uint32 {
	if x != nil && x.Answer != nil {
		return *x.Answer
	}
	return Default_Answer_Answer
}

func (x *Answer) GetCorrect() bool {
	if x != nil && x.Correct != nil {
		return *x.Correct
	}
	return Default_Answer_Correct
}

func (x *Answer) GetParamid() uint32 {
	if x != nil && x.Paramid != nil {
		return *x.Paramid
	}
	return Default_Answer_Paramid
}

type Query struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command     `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param *InterParam  `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=3" json:"param,omitempty"`
	Npcid *uint64      `protobuf:"varint,3,opt,name=npcid,def=0" json:"npcid,omitempty"`
	Ret   *EQueryState `protobuf:"varint,4,opt,name=ret,enum=Cmd.EQueryState" json:"ret,omitempty"`
}

// Default values for Query fields.
const (
	Default_Query_Cmd   = Command_SCENE_USER_INTER_PROTOCMD
	Default_Query_Param = InterParam_INTERPARAM_QUERYINTER
	Default_Query_Npcid = uint64(0)
)

func (x *Query) Reset() {
	*x = Query{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Query) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Query) ProtoMessage() {}

func (x *Query) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Query.ProtoReflect.Descriptor instead.
func (*Query) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{3}
}

func (x *Query) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_Query_Cmd
}

func (x *Query) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_Query_Param
}

func (x *Query) GetNpcid() uint64 {
	if x != nil && x.Npcid != nil {
		return *x.Npcid
	}
	return Default_Query_Npcid
}

func (x *Query) GetRet() EQueryState {
	if x != nil && x.Ret != nil {
		return *x.Ret
	}
	return EQueryState_EQUERYSTATE_OK
}

type PaperData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     *uint32  `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Result *uint32  `protobuf:"varint,2,opt,name=result" json:"result,omitempty"`
	Source *ESource `protobuf:"varint,3,opt,name=source,enum=Cmd.ESource" json:"source,omitempty"`
}

func (x *PaperData) Reset() {
	*x = PaperData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaperData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaperData) ProtoMessage() {}

func (x *PaperData) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaperData.ProtoReflect.Descriptor instead.
func (*PaperData) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{4}
}

func (x *PaperData) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

func (x *PaperData) GetResult() uint32 {
	if x != nil && x.Result != nil {
		return *x.Result
	}
	return 0
}

func (x *PaperData) GetSource() ESource {
	if x != nil && x.Source != nil {
		return *x.Source
	}
	return ESource_ESOURCE_MIN
}

type QueryPaperResultInterCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command     `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param *InterParam  `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=4" json:"param,omitempty"`
	Datas []*PaperData `protobuf:"bytes,3,rep,name=datas" json:"datas,omitempty"`
}

// Default values for QueryPaperResultInterCmd fields.
const (
	Default_QueryPaperResultInterCmd_Cmd   = Command_SCENE_USER_INTER_PROTOCMD
	Default_QueryPaperResultInterCmd_Param = InterParam_INTERPARAM_QUERYPAPER
)

func (x *QueryPaperResultInterCmd) Reset() {
	*x = QueryPaperResultInterCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QueryPaperResultInterCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueryPaperResultInterCmd) ProtoMessage() {}

func (x *QueryPaperResultInterCmd) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueryPaperResultInterCmd.ProtoReflect.Descriptor instead.
func (*QueryPaperResultInterCmd) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{5}
}

func (x *QueryPaperResultInterCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_QueryPaperResultInterCmd_Cmd
}

func (x *QueryPaperResultInterCmd) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_QueryPaperResultInterCmd_Param
}

func (x *QueryPaperResultInterCmd) GetDatas() []*PaperData {
	if x != nil {
		return x.Datas
	}
	return nil
}

type PaperQuestionInterCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd   *Command    `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param *InterParam `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=5" json:"param,omitempty"`
	Id    *uint32     `protobuf:"varint,3,opt,name=id" json:"id,omitempty"`
}

// Default values for PaperQuestionInterCmd fields.
const (
	Default_PaperQuestionInterCmd_Cmd   = Command_SCENE_USER_INTER_PROTOCMD
	Default_PaperQuestionInterCmd_Param = InterParam_INTERPARAM_PAPERQUESTION
)

func (x *PaperQuestionInterCmd) Reset() {
	*x = PaperQuestionInterCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaperQuestionInterCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaperQuestionInterCmd) ProtoMessage() {}

func (x *PaperQuestionInterCmd) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaperQuestionInterCmd.ProtoReflect.Descriptor instead.
func (*PaperQuestionInterCmd) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{6}
}

func (x *PaperQuestionInterCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_PaperQuestionInterCmd_Cmd
}

func (x *PaperQuestionInterCmd) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_PaperQuestionInterCmd_Param
}

func (x *PaperQuestionInterCmd) GetId() uint32 {
	if x != nil && x.Id != nil {
		return *x.Id
	}
	return 0
}

type PaperResultInterCmd struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Cmd    *Command    `protobuf:"varint,1,opt,name=cmd,enum=Cmd.Command,def=22" json:"cmd,omitempty"`
	Param  *InterParam `protobuf:"varint,2,opt,name=param,enum=Cmd.InterParam,def=6" json:"param,omitempty"`
	Result *PaperData  `protobuf:"bytes,3,opt,name=result" json:"result,omitempty"`
}

// Default values for PaperResultInterCmd fields.
const (
	Default_PaperResultInterCmd_Cmd   = Command_SCENE_USER_INTER_PROTOCMD
	Default_PaperResultInterCmd_Param = InterParam_INTERPARAM_PAPERRESULT
)

func (x *PaperResultInterCmd) Reset() {
	*x = PaperResultInterCmd{}
	if protoimpl.UnsafeEnabled {
		mi := &file_SceneInterlocution_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PaperResultInterCmd) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PaperResultInterCmd) ProtoMessage() {}

func (x *PaperResultInterCmd) ProtoReflect() protoreflect.Message {
	mi := &file_SceneInterlocution_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PaperResultInterCmd.ProtoReflect.Descriptor instead.
func (*PaperResultInterCmd) Descriptor() ([]byte, []int) {
	return file_SceneInterlocution_proto_rawDescGZIP(), []int{7}
}

func (x *PaperResultInterCmd) GetCmd() Command {
	if x != nil && x.Cmd != nil {
		return *x.Cmd
	}
	return Default_PaperResultInterCmd_Cmd
}

func (x *PaperResultInterCmd) GetParam() InterParam {
	if x != nil && x.Param != nil {
		return *x.Param
	}
	return Default_PaperResultInterCmd_Param
}

func (x *PaperResultInterCmd) GetResult() *PaperData {
	if x != nil {
		return x.Result
	}
	return nil
}

var File_SceneInterlocution_proto protoreflect.FileDescriptor

var file_SceneInterlocution_proto_rawDesc = []byte{
	0x0a, 0x18, 0x53, 0x63, 0x65, 0x6e, 0x65, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x6c, 0x6f, 0x63, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x43, 0x6d, 0x64, 0x1a,
	0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8f,
	0x01, 0x0a, 0x09, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x15, 0x0a, 0x04,
	0x67, 0x75, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x04, 0x67,
	0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x69, 0x64,
	0x12, 0x1b, 0x0a, 0x07, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x69, 0x64, 0x12, 0x31, 0x0a,
	0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e,
	0x43, 0x6d, 0x64, 0x2e, 0x45, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x3a, 0x0b, 0x45, 0x53, 0x4f,
	0x55, 0x52, 0x43, 0x45, 0x5f, 0x4d, 0x49, 0x4e, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x22, 0xe4, 0x01, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x39, 0x0a,
	0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x42, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x1b, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50,
	0x41, 0x52, 0x41, 0x4d, 0x5f, 0x4e, 0x45, 0x57, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4c, 0x4f, 0x43,
	0x55, 0x54, 0x49, 0x4f, 0x4e, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x24, 0x0a, 0x05,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x05, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x6e, 0x70, 0x63, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x6e, 0x70, 0x63, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x08, 0x61, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04, 0x3a, 0x01, 0x30, 0x52, 0x08, 0x61,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x69, 0x64, 0x22, 0xd8, 0x02, 0x0a, 0x06, 0x41, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x12, 0x39, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x19, 0x53,
	0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3d, 0x0a,
	0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x43,
	0x6d, 0x64, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x16, 0x49,
	0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x14, 0x0a, 0x05,
	0x6e, 0x70, 0x63, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x6e, 0x70, 0x63,
	0x69, 0x64, 0x12, 0x15, 0x0a, 0x04, 0x67, 0x75, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d,
	0x3a, 0x01, 0x30, 0x52, 0x04, 0x67, 0x75, 0x69, 0x64, 0x12, 0x1b, 0x0a, 0x07, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07, 0x69,
	0x6e, 0x74, 0x65, 0x72, 0x69, 0x64, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x53, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x3a, 0x0b, 0x45, 0x53, 0x4f, 0x55, 0x52, 0x43, 0x45, 0x5f, 0x4d, 0x49,
	0x4e, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x19, 0x0a, 0x06, 0x61, 0x6e, 0x73,
	0x77, 0x65, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x06, 0x61, 0x6e,
	0x73, 0x77, 0x65, 0x72, 0x12, 0x1f, 0x0a, 0x07, 0x63, 0x6f, 0x72, 0x72, 0x65, 0x63, 0x74, 0x18,
	0x08, 0x20, 0x01, 0x28, 0x08, 0x3a, 0x05, 0x66, 0x61, 0x6c, 0x73, 0x65, 0x52, 0x07, 0x63, 0x6f,
	0x72, 0x72, 0x65, 0x63, 0x74, 0x12, 0x1b, 0x0a, 0x07, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x69, 0x64,
	0x18, 0x09, 0x20, 0x01, 0x28, 0x0d, 0x3a, 0x01, 0x30, 0x52, 0x07, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x69, 0x64, 0x22, 0xbd, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x39, 0x0a, 0x03,
	0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3c, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x49, 0x6e, 0x74,
	0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x15, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41,
	0x52, 0x41, 0x4d, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x52, 0x05,
	0x70, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x17, 0x0a, 0x05, 0x6e, 0x70, 0x63, 0x69, 0x64, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x04, 0x3a, 0x01, 0x30, 0x52, 0x05, 0x6e, 0x70, 0x63, 0x69, 0x64, 0x12, 0x22,
	0x0a, 0x03, 0x72, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x10, 0x2e, 0x43, 0x6d,
	0x64, 0x2e, 0x45, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x03, 0x72,
	0x65, 0x74, 0x22, 0x59, 0x0a, 0x09, 0x50, 0x61, 0x70, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x24, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x45, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22, 0xb9, 0x01,
	0x0a, 0x18, 0x51, 0x75, 0x65, 0x72, 0x79, 0x50, 0x61, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x43, 0x6d, 0x64, 0x12, 0x39, 0x0a, 0x03, 0x63, 0x6d,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45,
	0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44,
	0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3c, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x15, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41,
	0x4d, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x50, 0x41, 0x50, 0x45, 0x52, 0x52, 0x05, 0x70, 0x61,
	0x72, 0x61, 0x6d, 0x12, 0x24, 0x0a, 0x05, 0x64, 0x61, 0x74, 0x61, 0x73, 0x18, 0x03, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x50, 0x61, 0x70, 0x65, 0x72, 0x44, 0x61,
	0x74, 0x61, 0x52, 0x05, 0x64, 0x61, 0x74, 0x61, 0x73, 0x22, 0xa3, 0x01, 0x0a, 0x15, 0x50, 0x61,
	0x70, 0x65, 0x72, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x72,
	0x43, 0x6d, 0x64, 0x12, 0x39, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e,
	0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x3a, 0x19,
	0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x4e, 0x54, 0x45, 0x52,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63, 0x6d, 0x64, 0x12, 0x3f,
	0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0f, 0x2e,
	0x43, 0x6d, 0x64, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x3a, 0x18,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x50, 0x41, 0x50, 0x45, 0x52,
	0x51, 0x55, 0x45, 0x53, 0x54, 0x49, 0x4f, 0x4e, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x22,
	0xb7, 0x01, 0x0a, 0x13, 0x50, 0x61, 0x70, 0x65, 0x72, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x49,
	0x6e, 0x74, 0x65, 0x72, 0x43, 0x6d, 0x64, 0x12, 0x39, 0x0a, 0x03, 0x63, 0x6d, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x0c, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x3a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49,
	0x4e, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x52, 0x03, 0x63,
	0x6d, 0x64, 0x12, 0x3d, 0x0a, 0x05, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x0f, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x3a, 0x16, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x50,
	0x41, 0x50, 0x45, 0x52, 0x52, 0x45, 0x53, 0x55, 0x4c, 0x54, 0x52, 0x05, 0x70, 0x61, 0x72, 0x61,
	0x6d, 0x12, 0x26, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0e, 0x2e, 0x43, 0x6d, 0x64, 0x2e, 0x50, 0x61, 0x70, 0x65, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2a, 0xb9, 0x01, 0x0a, 0x0a, 0x49, 0x6e,
	0x74, 0x65, 0x72, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x12, 0x1f, 0x0a, 0x1b, 0x49, 0x4e, 0x54, 0x45,
	0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x4e, 0x45, 0x57, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4c,
	0x4f, 0x43, 0x55, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x01, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x49, 0x4e,
	0x54, 0x45, 0x52, 0x10, 0x02, 0x12, 0x19, 0x0a, 0x15, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41,
	0x52, 0x41, 0x4d, 0x5f, 0x51, 0x55, 0x45, 0x52, 0x59, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x10, 0x03,
	0x12, 0x19, 0x0a, 0x15, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x51,
	0x55, 0x45, 0x52, 0x59, 0x50, 0x41, 0x50, 0x45, 0x52, 0x10, 0x04, 0x12, 0x1c, 0x0a, 0x18, 0x49,
	0x4e, 0x54, 0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x50, 0x41, 0x50, 0x45, 0x52, 0x51,
	0x55, 0x45, 0x53, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x05, 0x12, 0x1a, 0x0a, 0x16, 0x49, 0x4e, 0x54,
	0x45, 0x52, 0x50, 0x41, 0x52, 0x41, 0x4d, 0x5f, 0x50, 0x41, 0x50, 0x45, 0x52, 0x52, 0x45, 0x53,
	0x55, 0x4c, 0x54, 0x10, 0x06, 0x2a, 0x77, 0x0a, 0x0b, 0x45, 0x51, 0x75, 0x65, 0x72, 0x79, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x45, 0x51, 0x55, 0x45, 0x52, 0x59, 0x53, 0x54,
	0x41, 0x54, 0x45, 0x5f, 0x4f, 0x4b, 0x10, 0x01, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x51, 0x55, 0x45,
	0x52, 0x59, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x45, 0x44,
	0x5f, 0x52, 0x49, 0x47, 0x48, 0x54, 0x10, 0x02, 0x12, 0x1e, 0x0a, 0x1a, 0x45, 0x51, 0x55, 0x45,
	0x52, 0x59, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x41, 0x4e, 0x53, 0x57, 0x45, 0x52, 0x45, 0x44,
	0x5f, 0x57, 0x52, 0x4f, 0x4e, 0x47, 0x10, 0x03, 0x12, 0x14, 0x0a, 0x10, 0x45, 0x51, 0x55, 0x45,
	0x52, 0x59, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x46, 0x41, 0x49, 0x4c, 0x10, 0x04,
}

var (
	file_SceneInterlocution_proto_rawDescOnce sync.Once
	file_SceneInterlocution_proto_rawDescData = file_SceneInterlocution_proto_rawDesc
)

func file_SceneInterlocution_proto_rawDescGZIP() []byte {
	file_SceneInterlocution_proto_rawDescOnce.Do(func() {
		file_SceneInterlocution_proto_rawDescData = protoimpl.X.CompressGZIP(file_SceneInterlocution_proto_rawDescData)
	})
	return file_SceneInterlocution_proto_rawDescData
}

var file_SceneInterlocution_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_SceneInterlocution_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_SceneInterlocution_proto_goTypes = []interface{}{
	(InterParam)(0),                  // 0: Cmd.InterParam
	(EQueryState)(0),                 // 1: Cmd.EQueryState
	(*InterData)(nil),                // 2: Cmd.InterData
	(*NewInter)(nil),                 // 3: Cmd.NewInter
	(*Answer)(nil),                   // 4: Cmd.Answer
	(*Query)(nil),                    // 5: Cmd.Query
	(*PaperData)(nil),                // 6: Cmd.PaperData
	(*QueryPaperResultInterCmd)(nil), // 7: Cmd.QueryPaperResultInterCmd
	(*PaperQuestionInterCmd)(nil),    // 8: Cmd.PaperQuestionInterCmd
	(*PaperResultInterCmd)(nil),      // 9: Cmd.PaperResultInterCmd
	(ESource)(0),                     // 10: Cmd.ESource
	(Command)(0),                     // 11: Cmd.Command
}
var file_SceneInterlocution_proto_depIdxs = []int32{
	10, // 0: Cmd.InterData.source:type_name -> Cmd.ESource
	11, // 1: Cmd.NewInter.cmd:type_name -> Cmd.Command
	0,  // 2: Cmd.NewInter.param:type_name -> Cmd.InterParam
	2,  // 3: Cmd.NewInter.inter:type_name -> Cmd.InterData
	11, // 4: Cmd.Answer.cmd:type_name -> Cmd.Command
	0,  // 5: Cmd.Answer.param:type_name -> Cmd.InterParam
	10, // 6: Cmd.Answer.source:type_name -> Cmd.ESource
	11, // 7: Cmd.Query.cmd:type_name -> Cmd.Command
	0,  // 8: Cmd.Query.param:type_name -> Cmd.InterParam
	1,  // 9: Cmd.Query.ret:type_name -> Cmd.EQueryState
	10, // 10: Cmd.PaperData.source:type_name -> Cmd.ESource
	11, // 11: Cmd.QueryPaperResultInterCmd.cmd:type_name -> Cmd.Command
	0,  // 12: Cmd.QueryPaperResultInterCmd.param:type_name -> Cmd.InterParam
	6,  // 13: Cmd.QueryPaperResultInterCmd.datas:type_name -> Cmd.PaperData
	11, // 14: Cmd.PaperQuestionInterCmd.cmd:type_name -> Cmd.Command
	0,  // 15: Cmd.PaperQuestionInterCmd.param:type_name -> Cmd.InterParam
	11, // 16: Cmd.PaperResultInterCmd.cmd:type_name -> Cmd.Command
	0,  // 17: Cmd.PaperResultInterCmd.param:type_name -> Cmd.InterParam
	6,  // 18: Cmd.PaperResultInterCmd.result:type_name -> Cmd.PaperData
	19, // [19:19] is the sub-list for method output_type
	19, // [19:19] is the sub-list for method input_type
	19, // [19:19] is the sub-list for extension type_name
	19, // [19:19] is the sub-list for extension extendee
	0,  // [0:19] is the sub-list for field type_name
}

func init() { file_SceneInterlocution_proto_init() }
func file_SceneInterlocution_proto_init() {
	if File_SceneInterlocution_proto != nil {
		return
	}
	file_xCmd_proto_init()
	file_ProtoCommon_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_SceneInterlocution_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InterData); i {
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
		file_SceneInterlocution_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NewInter); i {
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
		file_SceneInterlocution_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Answer); i {
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
		file_SceneInterlocution_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Query); i {
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
		file_SceneInterlocution_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaperData); i {
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
		file_SceneInterlocution_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QueryPaperResultInterCmd); i {
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
		file_SceneInterlocution_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaperQuestionInterCmd); i {
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
		file_SceneInterlocution_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PaperResultInterCmd); i {
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
			RawDescriptor: file_SceneInterlocution_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_SceneInterlocution_proto_goTypes,
		DependencyIndexes: file_SceneInterlocution_proto_depIdxs,
		EnumInfos:         file_SceneInterlocution_proto_enumTypes,
		MessageInfos:      file_SceneInterlocution_proto_msgTypes,
	}.Build()
	File_SceneInterlocution_proto = out.File
	file_SceneInterlocution_proto_rawDesc = nil
	file_SceneInterlocution_proto_goTypes = nil
	file_SceneInterlocution_proto_depIdxs = nil
}
