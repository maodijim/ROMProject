// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0-devel
// 	protoc        v3.14.0
// source: xCmd.proto

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

type Command int32

const (
	Command_LOGIN_USER_PROTOCMD             Command = 1
	Command_ERROR_USER_PROTOCMD             Command = 2
	Command_SCENE_USER_PROTOCMD             Command = 5
	Command_SCENE_USER_ITEM_PROTOCMD        Command = 6
	Command_SCENE_USER_SKILL_PROTOCMD       Command = 7
	Command_SCENE_USER_QUEST_PROTOCMD       Command = 8
	Command_SCENE_USER2_PROTOCMD            Command = 9
	Command_SCENE_USER_PET_PROTOCMD         Command = 10
	Command_FUBEN_PROTOCMD                  Command = 11
	Command_SCENE_USER_MAP_PROTOCMD         Command = 12
	Command_SCENE_USER_MOUNT_PROTOCMD       Command = 13
	Command_SCENE_BOSS_PROTOCMD             Command = 15
	Command_SCENE_USER_CARRIER_PROTOCMD     Command = 16
	Command_SCENE_USER_ACHIEVE_PROTOCMD     Command = 17
	Command_SCENE_USER_TIP_PROTOCMD         Command = 18
	Command_SCENE_USER_CHATROOM_PROTOCMD    Command = 19
	Command_INFINITE_TOWER_PROTOCMD         Command = 20
	Command_SCENE_USER_SEAL_PROTOCMD        Command = 21
	Command_SCENE_USER_INTER_PROTOCMD       Command = 22
	Command_SCENE_USER_MANUAL_PROTOCMD      Command = 23
	Command_SCENE_USER_CHAT_PROTOCMD        Command = 24
	Command_USER_EVENT_PROTOCMD             Command = 25
	Command_SCENE_USER_TRADE_PROTOCMD       Command = 26
	Command_SCENE_USER_AUGURY_PROTOCMD      Command = 27
	Command_SCENE_USER_ASTROLABE_PROTOCMD   Command = 28
	Command_SCENE_USER_FOOD_PROTOCMD        Command = 29
	Command_SCENE_USER_PHOTO_PROTOCMD       Command = 30
	Command_SCENE_USER_TUTOR_PROTOCMD       Command = 31
	Command_SCENE_USER_BEING_PROTOCMD       Command = 32
	Command_SESSION_USER_GUILD_PROTOCMD     Command = 50
	Command_SESSION_USER_TEAM_PROTOCMD      Command = 51
	Command_SESSION_USER_SHOP_PROTOCMD      Command = 52
	Command_SESSION_USER_WEATHER_PROTOCMD   Command = 53
	Command_SESSION_USER_MAIL_PROTOCMD      Command = 55
	Command_SESSION_USER_SOCIALITY_PROTOCMD Command = 56
	Command_RECORD_USER_TRADE_PROTOCMD      Command = 57
	Command_DOJO_PROTOCMD                   Command = 58
	Command_CHAT_PROTOCMD                   Command = 59
	Command_ACTIVITY_PROTOCMD               Command = 60
	Command_MATCHC_PROTOCMD                 Command = 61
	Command_SESSION_USER_AUTHORIZE_PROTOCMD Command = 62
	Command_AUCTIONC_PROTOCMD               Command = 63
	Command_ACTIVITY_EVENT_PROTOCMD         Command = 64
	Command_WEDDINGC_PROTOCMD               Command = 65
	Command_PVE_CARD_PROTOCMD               Command = 66
	Command_TEAM_RAID_PROTOCMD              Command = 67
	Command_PUZZLE_PROTOCMD                 Command = 68
	Command_TEAM_GROUP_RAID_PROTOCMD        Command = 69
	Command_HOMEC_PROTOCMD                  Command = 70
	Command_ROGUELIKE_PROTOCMD              Command = 71
	Command_ROGUELIKE_PROTOSCMD             Command = 72
	Command_TECHTREE_PROTOCMD               Command = 73
	Command_USER_AFK_PROTOCMD               Command = 74
	Command_GOAL_PROTOCMD                   Command = 75
	Command_RAID_PROTOCMD                   Command = 76
	Command_SESSION_OVERSEAS_TW_PROTOCMD    Command = 80
	Command_SCENE_OVERSEAS_PROTOCMD         Command = 81
	Command_CLIENT_CMD                      Command = 99
	Command_MAX_USER_CMD                    Command = 100
	Command_RECORD_DATA_PROTOCMD            Command = 200
	Command_TRADE_PROTOCMD                  Command = 201
	Command_SESSION_PROTOCMD                Command = 202
	Command_GMTOOLS_PROTOCMD                Command = 203
	Command_LOG_PROTOCMD                    Command = 204
	Command_GATE_SUPER_PROTOCMD             Command = 205
	Command_REGION_PROTOCMD                 Command = 206
	Command_STAT_PROTOCMD                   Command = 207
	Command_SOCIAL_PROTOCMD                 Command = 208
	Command_TEAM_PROTOCMD                   Command = 209
	Command_GUILD_PROTOCMD                  Command = 210
	Command_GZONE_PROTOCMD                  Command = 211
	Command_MATCHS_PROTOCMD                 Command = 212
	Command_AUCTIONS_PROTOCMD               Command = 213
	Command_WEDDINGS_PROTOCMD               Command = 214
	Command_GTEAM_PROTOCMD                  Command = 215
	Command_BOSSS_PROTOCMD                  Command = 216
	Command_INTERACT_PROTOCMD               Command = 217
	Command_CARRIERS_PROTOCMD               Command = 218
	Command_HOMES_PROTOCMD                  Command = 219
	Command_MONITOR_PROTOCMD                Command = 220
	Command_CHATS_PROTOCMD                  Command = 221
	Command_BATTLEPASS_PROTOCMD             Command = 222
	Command_MINIGAME_PROTOCMD               Command = 223
	Command_REWARD_PROTOCMD                 Command = 224
	Command_USERSHOW_PROTOCMD               Command = 225
	Command_ACTHITPOLLY_PROTOCMD            Command = 226
	Command_QUESTS_PROTOCMD                 Command = 227
	Command_MINIGAMES_PROTOCMD              Command = 228
	Command_ACTMINIRO_PROTOCMD              Command = 229
	Command_RAIDS_PROTOCMD                  Command = 230
	Command_NOVICE_NOTEBOOK                 Command = 231
	Command_DISNEY_ACTIVITY_PROTOCMD        Command = 232
	Command_SCENE_USER_MANOR_PROTOCMD       Command = 233
	Command_REG_CMD                         Command = 253
	Command_GATEWAY_CMD                     Command = 250
	Command_SYSTEM_PROTOCMD                 Command = 255
)

// Enum value maps for Command.
var (
	Command_name = map[int32]string{
		1:   "LOGIN_USER_PROTOCMD",
		2:   "ERROR_USER_PROTOCMD",
		5:   "SCENE_USER_PROTOCMD",
		6:   "SCENE_USER_ITEM_PROTOCMD",
		7:   "SCENE_USER_SKILL_PROTOCMD",
		8:   "SCENE_USER_QUEST_PROTOCMD",
		9:   "SCENE_USER2_PROTOCMD",
		10:  "SCENE_USER_PET_PROTOCMD",
		11:  "FUBEN_PROTOCMD",
		12:  "SCENE_USER_MAP_PROTOCMD",
		13:  "SCENE_USER_MOUNT_PROTOCMD",
		15:  "SCENE_BOSS_PROTOCMD",
		16:  "SCENE_USER_CARRIER_PROTOCMD",
		17:  "SCENE_USER_ACHIEVE_PROTOCMD",
		18:  "SCENE_USER_TIP_PROTOCMD",
		19:  "SCENE_USER_CHATROOM_PROTOCMD",
		20:  "INFINITE_TOWER_PROTOCMD",
		21:  "SCENE_USER_SEAL_PROTOCMD",
		22:  "SCENE_USER_INTER_PROTOCMD",
		23:  "SCENE_USER_MANUAL_PROTOCMD",
		24:  "SCENE_USER_CHAT_PROTOCMD",
		25:  "USER_EVENT_PROTOCMD",
		26:  "SCENE_USER_TRADE_PROTOCMD",
		27:  "SCENE_USER_AUGURY_PROTOCMD",
		28:  "SCENE_USER_ASTROLABE_PROTOCMD",
		29:  "SCENE_USER_FOOD_PROTOCMD",
		30:  "SCENE_USER_PHOTO_PROTOCMD",
		31:  "SCENE_USER_TUTOR_PROTOCMD",
		32:  "SCENE_USER_BEING_PROTOCMD",
		50:  "SESSION_USER_GUILD_PROTOCMD",
		51:  "SESSION_USER_TEAM_PROTOCMD",
		52:  "SESSION_USER_SHOP_PROTOCMD",
		53:  "SESSION_USER_WEATHER_PROTOCMD",
		55:  "SESSION_USER_MAIL_PROTOCMD",
		56:  "SESSION_USER_SOCIALITY_PROTOCMD",
		57:  "RECORD_USER_TRADE_PROTOCMD",
		58:  "DOJO_PROTOCMD",
		59:  "CHAT_PROTOCMD",
		60:  "ACTIVITY_PROTOCMD",
		61:  "MATCHC_PROTOCMD",
		62:  "SESSION_USER_AUTHORIZE_PROTOCMD",
		63:  "AUCTIONC_PROTOCMD",
		64:  "ACTIVITY_EVENT_PROTOCMD",
		65:  "WEDDINGC_PROTOCMD",
		66:  "PVE_CARD_PROTOCMD",
		67:  "TEAM_RAID_PROTOCMD",
		68:  "PUZZLE_PROTOCMD",
		69:  "TEAM_GROUP_RAID_PROTOCMD",
		70:  "HOMEC_PROTOCMD",
		71:  "ROGUELIKE_PROTOCMD",
		72:  "ROGUELIKE_PROTOSCMD",
		73:  "TECHTREE_PROTOCMD",
		74:  "USER_AFK_PROTOCMD",
		75:  "GOAL_PROTOCMD",
		76:  "RAID_PROTOCMD",
		80:  "SESSION_OVERSEAS_TW_PROTOCMD",
		81:  "SCENE_OVERSEAS_PROTOCMD",
		99:  "CLIENT_CMD",
		100: "MAX_USER_CMD",
		200: "RECORD_DATA_PROTOCMD",
		201: "TRADE_PROTOCMD",
		202: "SESSION_PROTOCMD",
		203: "GMTOOLS_PROTOCMD",
		204: "LOG_PROTOCMD",
		205: "GATE_SUPER_PROTOCMD",
		206: "REGION_PROTOCMD",
		207: "STAT_PROTOCMD",
		208: "SOCIAL_PROTOCMD",
		209: "TEAM_PROTOCMD",
		210: "GUILD_PROTOCMD",
		211: "GZONE_PROTOCMD",
		212: "MATCHS_PROTOCMD",
		213: "AUCTIONS_PROTOCMD",
		214: "WEDDINGS_PROTOCMD",
		215: "GTEAM_PROTOCMD",
		216: "BOSSS_PROTOCMD",
		217: "INTERACT_PROTOCMD",
		218: "CARRIERS_PROTOCMD",
		219: "HOMES_PROTOCMD",
		220: "MONITOR_PROTOCMD",
		221: "CHATS_PROTOCMD",
		222: "BATTLEPASS_PROTOCMD",
		223: "MINIGAME_PROTOCMD",
		224: "REWARD_PROTOCMD",
		225: "USERSHOW_PROTOCMD",
		226: "ACTHITPOLLY_PROTOCMD",
		227: "QUESTS_PROTOCMD",
		228: "MINIGAMES_PROTOCMD",
		229: "ACTMINIRO_PROTOCMD",
		230: "RAIDS_PROTOCMD",
		231: "NOVICE_NOTEBOOK",
		232: "DISNEY_ACTIVITY_PROTOCMD",
		233: "SCENE_USER_MANOR_PROTOCMD",
		253: "REG_CMD",
		250: "GATEWAY_CMD",
		255: "SYSTEM_PROTOCMD",
	}
	Command_value = map[string]int32{
		"LOGIN_USER_PROTOCMD":             1,
		"ERROR_USER_PROTOCMD":             2,
		"SCENE_USER_PROTOCMD":             5,
		"SCENE_USER_ITEM_PROTOCMD":        6,
		"SCENE_USER_SKILL_PROTOCMD":       7,
		"SCENE_USER_QUEST_PROTOCMD":       8,
		"SCENE_USER2_PROTOCMD":            9,
		"SCENE_USER_PET_PROTOCMD":         10,
		"FUBEN_PROTOCMD":                  11,
		"SCENE_USER_MAP_PROTOCMD":         12,
		"SCENE_USER_MOUNT_PROTOCMD":       13,
		"SCENE_BOSS_PROTOCMD":             15,
		"SCENE_USER_CARRIER_PROTOCMD":     16,
		"SCENE_USER_ACHIEVE_PROTOCMD":     17,
		"SCENE_USER_TIP_PROTOCMD":         18,
		"SCENE_USER_CHATROOM_PROTOCMD":    19,
		"INFINITE_TOWER_PROTOCMD":         20,
		"SCENE_USER_SEAL_PROTOCMD":        21,
		"SCENE_USER_INTER_PROTOCMD":       22,
		"SCENE_USER_MANUAL_PROTOCMD":      23,
		"SCENE_USER_CHAT_PROTOCMD":        24,
		"USER_EVENT_PROTOCMD":             25,
		"SCENE_USER_TRADE_PROTOCMD":       26,
		"SCENE_USER_AUGURY_PROTOCMD":      27,
		"SCENE_USER_ASTROLABE_PROTOCMD":   28,
		"SCENE_USER_FOOD_PROTOCMD":        29,
		"SCENE_USER_PHOTO_PROTOCMD":       30,
		"SCENE_USER_TUTOR_PROTOCMD":       31,
		"SCENE_USER_BEING_PROTOCMD":       32,
		"SESSION_USER_GUILD_PROTOCMD":     50,
		"SESSION_USER_TEAM_PROTOCMD":      51,
		"SESSION_USER_SHOP_PROTOCMD":      52,
		"SESSION_USER_WEATHER_PROTOCMD":   53,
		"SESSION_USER_MAIL_PROTOCMD":      55,
		"SESSION_USER_SOCIALITY_PROTOCMD": 56,
		"RECORD_USER_TRADE_PROTOCMD":      57,
		"DOJO_PROTOCMD":                   58,
		"CHAT_PROTOCMD":                   59,
		"ACTIVITY_PROTOCMD":               60,
		"MATCHC_PROTOCMD":                 61,
		"SESSION_USER_AUTHORIZE_PROTOCMD": 62,
		"AUCTIONC_PROTOCMD":               63,
		"ACTIVITY_EVENT_PROTOCMD":         64,
		"WEDDINGC_PROTOCMD":               65,
		"PVE_CARD_PROTOCMD":               66,
		"TEAM_RAID_PROTOCMD":              67,
		"PUZZLE_PROTOCMD":                 68,
		"TEAM_GROUP_RAID_PROTOCMD":        69,
		"HOMEC_PROTOCMD":                  70,
		"ROGUELIKE_PROTOCMD":              71,
		"ROGUELIKE_PROTOSCMD":             72,
		"TECHTREE_PROTOCMD":               73,
		"USER_AFK_PROTOCMD":               74,
		"GOAL_PROTOCMD":                   75,
		"RAID_PROTOCMD":                   76,
		"SESSION_OVERSEAS_TW_PROTOCMD":    80,
		"SCENE_OVERSEAS_PROTOCMD":         81,
		"CLIENT_CMD":                      99,
		"MAX_USER_CMD":                    100,
		"RECORD_DATA_PROTOCMD":            200,
		"TRADE_PROTOCMD":                  201,
		"SESSION_PROTOCMD":                202,
		"GMTOOLS_PROTOCMD":                203,
		"LOG_PROTOCMD":                    204,
		"GATE_SUPER_PROTOCMD":             205,
		"REGION_PROTOCMD":                 206,
		"STAT_PROTOCMD":                   207,
		"SOCIAL_PROTOCMD":                 208,
		"TEAM_PROTOCMD":                   209,
		"GUILD_PROTOCMD":                  210,
		"GZONE_PROTOCMD":                  211,
		"MATCHS_PROTOCMD":                 212,
		"AUCTIONS_PROTOCMD":               213,
		"WEDDINGS_PROTOCMD":               214,
		"GTEAM_PROTOCMD":                  215,
		"BOSSS_PROTOCMD":                  216,
		"INTERACT_PROTOCMD":               217,
		"CARRIERS_PROTOCMD":               218,
		"HOMES_PROTOCMD":                  219,
		"MONITOR_PROTOCMD":                220,
		"CHATS_PROTOCMD":                  221,
		"BATTLEPASS_PROTOCMD":             222,
		"MINIGAME_PROTOCMD":               223,
		"REWARD_PROTOCMD":                 224,
		"USERSHOW_PROTOCMD":               225,
		"ACTHITPOLLY_PROTOCMD":            226,
		"QUESTS_PROTOCMD":                 227,
		"MINIGAMES_PROTOCMD":              228,
		"ACTMINIRO_PROTOCMD":              229,
		"RAIDS_PROTOCMD":                  230,
		"NOVICE_NOTEBOOK":                 231,
		"DISNEY_ACTIVITY_PROTOCMD":        232,
		"SCENE_USER_MANOR_PROTOCMD":       233,
		"REG_CMD":                         253,
		"GATEWAY_CMD":                     250,
		"SYSTEM_PROTOCMD":                 255,
	}
)

func (x Command) Enum() *Command {
	p := new(Command)
	*p = x
	return p
}

func (x Command) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Command) Descriptor() protoreflect.EnumDescriptor {
	return file_xCmd_proto_enumTypes[0].Descriptor()
}

func (Command) Type() protoreflect.EnumType {
	return &file_xCmd_proto_enumTypes[0]
}

func (x Command) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *Command) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = Command(num)
	return nil
}

// Deprecated: Use Command.Descriptor instead.
func (Command) EnumDescriptor() ([]byte, []int) {
	return file_xCmd_proto_rawDescGZIP(), []int{0}
}

type Nonce struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Timestamp *uint32 `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	Index     *uint32 `protobuf:"varint,2,opt,name=index" json:"index,omitempty"`
	Sign      *string `protobuf:"bytes,3,opt,name=sign" json:"sign,omitempty"`
}

func (x *Nonce) Reset() {
	*x = Nonce{}
	if protoimpl.UnsafeEnabled {
		mi := &file_xCmd_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Nonce) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Nonce) ProtoMessage() {}

func (x *Nonce) ProtoReflect() protoreflect.Message {
	mi := &file_xCmd_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Nonce.ProtoReflect.Descriptor instead.
func (*Nonce) Descriptor() ([]byte, []int) {
	return file_xCmd_proto_rawDescGZIP(), []int{0}
}

func (x *Nonce) GetTimestamp() uint32 {
	if x != nil && x.Timestamp != nil {
		return *x.Timestamp
	}
	return 0
}

func (x *Nonce) GetIndex() uint32 {
	if x != nil && x.Index != nil {
		return *x.Index
	}
	return 0
}

func (x *Nonce) GetSign() string {
	if x != nil && x.Sign != nil {
		return *x.Sign
	}
	return ""
}

var File_xCmd_proto protoreflect.FileDescriptor

var file_xCmd_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x78, 0x43, 0x6d, 0x64, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x03, 0x43, 0x6d,
	0x64, 0x22, 0x4f, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x09, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65,
	0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12,
	0x0a, 0x04, 0x73, 0x69, 0x67, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x73, 0x69,
	0x67, 0x6e, 0x2a, 0xb3, 0x13, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x17,
	0x0a, 0x13, 0x4c, 0x4f, 0x47, 0x49, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x02,
	0x12, 0x17, 0x0a, 0x13, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x05, 0x12, 0x1c, 0x0a, 0x18, 0x53, 0x43, 0x45,
	0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x49, 0x54, 0x45, 0x4d, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x06, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x4b, 0x49, 0x4c, 0x4c, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x07, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x51, 0x55, 0x45, 0x53, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x08, 0x12, 0x18, 0x0a, 0x14, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x32, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x09, 0x12,
	0x1b, 0x0a, 0x17, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x45,
	0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x0a, 0x12, 0x12, 0x0a, 0x0e,
	0x46, 0x55, 0x42, 0x45, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x0b,
	0x12, 0x1b, 0x0a, 0x17, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4d,
	0x41, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x0c, 0x12, 0x1d, 0x0a,
	0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4d, 0x4f, 0x55, 0x4e,
	0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x0d, 0x12, 0x17, 0x0a, 0x13,
	0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x42, 0x4f, 0x53, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x0f, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x43, 0x41, 0x52, 0x52, 0x49, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x10, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x43, 0x48, 0x49, 0x45, 0x56, 0x45, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x11, 0x12, 0x1b, 0x0a, 0x17, 0x53, 0x43, 0x45, 0x4e, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x49, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0x12, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53,
	0x45, 0x52, 0x5f, 0x43, 0x48, 0x41, 0x54, 0x52, 0x4f, 0x4f, 0x4d, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x13, 0x12, 0x1b, 0x0a, 0x17, 0x49, 0x4e, 0x46, 0x49, 0x4e, 0x49,
	0x54, 0x45, 0x5f, 0x54, 0x4f, 0x57, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0x14, 0x12, 0x1c, 0x0a, 0x18, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45,
	0x52, 0x5f, 0x53, 0x45, 0x41, 0x4c, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0x15, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f,
	0x49, 0x4e, 0x54, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x16,
	0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4d,
	0x41, 0x4e, 0x55, 0x41, 0x4c, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x17,
	0x12, 0x1c, 0x0a, 0x18, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x43,
	0x48, 0x41, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x18, 0x12, 0x17,
	0x0a, 0x13, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x19, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x52, 0x41, 0x44, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x1a, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x55, 0x47, 0x55, 0x52, 0x59, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x1b, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x53, 0x54, 0x52, 0x4f, 0x4c, 0x41, 0x42, 0x45, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x1c, 0x12, 0x1c, 0x0a, 0x18, 0x53, 0x43, 0x45,
	0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x46, 0x4f, 0x4f, 0x44, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x1d, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x50, 0x48, 0x4f, 0x54, 0x4f, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x1e, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x55, 0x54, 0x4f, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x1f, 0x12, 0x1d, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55,
	0x53, 0x45, 0x52, 0x5f, 0x42, 0x45, 0x49, 0x4e, 0x47, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0x20, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f,
	0x55, 0x53, 0x45, 0x52, 0x5f, 0x47, 0x55, 0x49, 0x4c, 0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x32, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x45, 0x41, 0x4d, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x33, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x48, 0x4f, 0x50, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x34, 0x12, 0x21, 0x0a, 0x1d, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x57, 0x45, 0x41, 0x54, 0x48, 0x45, 0x52, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x35, 0x12, 0x1e, 0x0a, 0x1a, 0x53, 0x45, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x4d, 0x41, 0x49, 0x4c, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x37, 0x12, 0x23, 0x0a, 0x1f, 0x53, 0x45, 0x53, 0x53,
	0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x53, 0x4f, 0x43, 0x49, 0x41, 0x4c, 0x49,
	0x54, 0x59, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x38, 0x12, 0x1e, 0x0a,
	0x1a, 0x52, 0x45, 0x43, 0x4f, 0x52, 0x44, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x54, 0x52, 0x41,
	0x44, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x39, 0x12, 0x11, 0x0a,
	0x0d, 0x44, 0x4f, 0x4a, 0x4f, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x3a,
	0x12, 0x11, 0x0a, 0x0d, 0x43, 0x48, 0x41, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0x3b, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x43, 0x54, 0x49, 0x56, 0x49, 0x54, 0x59, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x3c, 0x12, 0x13, 0x0a, 0x0f, 0x4d, 0x41,
	0x54, 0x43, 0x48, 0x43, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x3d, 0x12,
	0x23, 0x0a, 0x1f, 0x53, 0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f,
	0x41, 0x55, 0x54, 0x48, 0x4f, 0x52, 0x49, 0x5a, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0x3e, 0x12, 0x15, 0x0a, 0x11, 0x41, 0x55, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x43,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x3f, 0x12, 0x1b, 0x0a, 0x17, 0x41,
	0x43, 0x54, 0x49, 0x56, 0x49, 0x54, 0x59, 0x5f, 0x45, 0x56, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x40, 0x12, 0x15, 0x0a, 0x11, 0x57, 0x45, 0x44, 0x44,
	0x49, 0x4e, 0x47, 0x43, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x41, 0x12,
	0x15, 0x0a, 0x11, 0x50, 0x56, 0x45, 0x5f, 0x43, 0x41, 0x52, 0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0x42, 0x12, 0x16, 0x0a, 0x12, 0x54, 0x45, 0x41, 0x4d, 0x5f, 0x52,
	0x41, 0x49, 0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x43, 0x12, 0x13,
	0x0a, 0x0f, 0x50, 0x55, 0x5a, 0x5a, 0x4c, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0x44, 0x12, 0x1c, 0x0a, 0x18, 0x54, 0x45, 0x41, 0x4d, 0x5f, 0x47, 0x52, 0x4f, 0x55,
	0x50, 0x5f, 0x52, 0x41, 0x49, 0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0x45, 0x12, 0x12, 0x0a, 0x0e, 0x48, 0x4f, 0x4d, 0x45, 0x43, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0x46, 0x12, 0x16, 0x0a, 0x12, 0x52, 0x4f, 0x47, 0x55, 0x45, 0x4c, 0x49,
	0x4b, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x47, 0x12, 0x17, 0x0a,
	0x13, 0x52, 0x4f, 0x47, 0x55, 0x45, 0x4c, 0x49, 0x4b, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x53, 0x43, 0x4d, 0x44, 0x10, 0x48, 0x12, 0x15, 0x0a, 0x11, 0x54, 0x45, 0x43, 0x48, 0x54, 0x52,
	0x45, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x49, 0x12, 0x15, 0x0a,
	0x11, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x41, 0x46, 0x4b, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0x4a, 0x12, 0x11, 0x0a, 0x0d, 0x47, 0x4f, 0x41, 0x4c, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x4b, 0x12, 0x11, 0x0a, 0x0d, 0x52, 0x41, 0x49, 0x44, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x4c, 0x12, 0x20, 0x0a, 0x1c, 0x53, 0x45,
	0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x4f, 0x56, 0x45, 0x52, 0x53, 0x45, 0x41, 0x53, 0x5f, 0x54,
	0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x50, 0x12, 0x1b, 0x0a, 0x17,
	0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x4f, 0x56, 0x45, 0x52, 0x53, 0x45, 0x41, 0x53, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0x51, 0x12, 0x0e, 0x0a, 0x0a, 0x43, 0x4c, 0x49,
	0x45, 0x4e, 0x54, 0x5f, 0x43, 0x4d, 0x44, 0x10, 0x63, 0x12, 0x10, 0x0a, 0x0c, 0x4d, 0x41, 0x58,
	0x5f, 0x55, 0x53, 0x45, 0x52, 0x5f, 0x43, 0x4d, 0x44, 0x10, 0x64, 0x12, 0x19, 0x0a, 0x14, 0x52,
	0x45, 0x43, 0x4f, 0x52, 0x44, 0x5f, 0x44, 0x41, 0x54, 0x41, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0xc8, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x54, 0x52, 0x41, 0x44, 0x45, 0x5f,
	0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xc9, 0x01, 0x12, 0x15, 0x0a, 0x10, 0x53,
	0x45, 0x53, 0x53, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0xca, 0x01, 0x12, 0x15, 0x0a, 0x10, 0x47, 0x4d, 0x54, 0x4f, 0x4f, 0x4c, 0x53, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xcb, 0x01, 0x12, 0x11, 0x0a, 0x0c, 0x4c, 0x4f, 0x47,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xcc, 0x01, 0x12, 0x18, 0x0a, 0x13,
	0x47, 0x41, 0x54, 0x45, 0x5f, 0x53, 0x55, 0x50, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0xcd, 0x01, 0x12, 0x14, 0x0a, 0x0f, 0x52, 0x45, 0x47, 0x49, 0x4f, 0x4e,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xce, 0x01, 0x12, 0x12, 0x0a, 0x0d,
	0x53, 0x54, 0x41, 0x54, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xcf, 0x01,
	0x12, 0x14, 0x0a, 0x0f, 0x53, 0x4f, 0x43, 0x49, 0x41, 0x4c, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0xd0, 0x01, 0x12, 0x12, 0x0a, 0x0d, 0x54, 0x45, 0x41, 0x4d, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd1, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x47, 0x55,
	0x49, 0x4c, 0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd2, 0x01, 0x12,
	0x13, 0x0a, 0x0e, 0x47, 0x5a, 0x4f, 0x4e, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0xd3, 0x01, 0x12, 0x14, 0x0a, 0x0f, 0x4d, 0x41, 0x54, 0x43, 0x48, 0x53, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd4, 0x01, 0x12, 0x16, 0x0a, 0x11, 0x41, 0x55,
	0x43, 0x54, 0x49, 0x4f, 0x4e, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0xd5, 0x01, 0x12, 0x16, 0x0a, 0x11, 0x57, 0x45, 0x44, 0x44, 0x49, 0x4e, 0x47, 0x53, 0x5f, 0x50,
	0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd6, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x47, 0x54,
	0x45, 0x41, 0x4d, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd7, 0x01, 0x12,
	0x13, 0x0a, 0x0e, 0x42, 0x4f, 0x53, 0x53, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0xd8, 0x01, 0x12, 0x16, 0x0a, 0x11, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x41, 0x43, 0x54,
	0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xd9, 0x01, 0x12, 0x16, 0x0a, 0x11,
	0x43, 0x41, 0x52, 0x52, 0x49, 0x45, 0x52, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d,
	0x44, 0x10, 0xda, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x48, 0x4f, 0x4d, 0x45, 0x53, 0x5f, 0x50, 0x52,
	0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xdb, 0x01, 0x12, 0x15, 0x0a, 0x10, 0x4d, 0x4f, 0x4e,
	0x49, 0x54, 0x4f, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xdc, 0x01,
	0x12, 0x13, 0x0a, 0x0e, 0x43, 0x48, 0x41, 0x54, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0xdd, 0x01, 0x12, 0x18, 0x0a, 0x13, 0x42, 0x41, 0x54, 0x54, 0x4c, 0x45, 0x50,
	0x41, 0x53, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xde, 0x01, 0x12,
	0x16, 0x0a, 0x11, 0x4d, 0x49, 0x4e, 0x49, 0x47, 0x41, 0x4d, 0x45, 0x5f, 0x50, 0x52, 0x4f, 0x54,
	0x4f, 0x43, 0x4d, 0x44, 0x10, 0xdf, 0x01, 0x12, 0x14, 0x0a, 0x0f, 0x52, 0x45, 0x57, 0x41, 0x52,
	0x44, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xe0, 0x01, 0x12, 0x16, 0x0a,
	0x11, 0x55, 0x53, 0x45, 0x52, 0x53, 0x48, 0x4f, 0x57, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43,
	0x4d, 0x44, 0x10, 0xe1, 0x01, 0x12, 0x19, 0x0a, 0x14, 0x41, 0x43, 0x54, 0x48, 0x49, 0x54, 0x50,
	0x4f, 0x4c, 0x4c, 0x59, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xe2, 0x01,
	0x12, 0x14, 0x0a, 0x0f, 0x51, 0x55, 0x45, 0x53, 0x54, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f,
	0x43, 0x4d, 0x44, 0x10, 0xe3, 0x01, 0x12, 0x17, 0x0a, 0x12, 0x4d, 0x49, 0x4e, 0x49, 0x47, 0x41,
	0x4d, 0x45, 0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xe4, 0x01, 0x12,
	0x17, 0x0a, 0x12, 0x41, 0x43, 0x54, 0x4d, 0x49, 0x4e, 0x49, 0x52, 0x4f, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xe5, 0x01, 0x12, 0x13, 0x0a, 0x0e, 0x52, 0x41, 0x49, 0x44,
	0x53, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xe6, 0x01, 0x12, 0x14, 0x0a,
	0x0f, 0x4e, 0x4f, 0x56, 0x49, 0x43, 0x45, 0x5f, 0x4e, 0x4f, 0x54, 0x45, 0x42, 0x4f, 0x4f, 0x4b,
	0x10, 0xe7, 0x01, 0x12, 0x1d, 0x0a, 0x18, 0x44, 0x49, 0x53, 0x4e, 0x45, 0x59, 0x5f, 0x41, 0x43,
	0x54, 0x49, 0x56, 0x49, 0x54, 0x59, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0xe8, 0x01, 0x12, 0x1e, 0x0a, 0x19, 0x53, 0x43, 0x45, 0x4e, 0x45, 0x5f, 0x55, 0x53, 0x45, 0x52,
	0x5f, 0x4d, 0x41, 0x4e, 0x4f, 0x52, 0x5f, 0x50, 0x52, 0x4f, 0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10,
	0xe9, 0x01, 0x12, 0x0c, 0x0a, 0x07, 0x52, 0x45, 0x47, 0x5f, 0x43, 0x4d, 0x44, 0x10, 0xfd, 0x01,
	0x12, 0x10, 0x0a, 0x0b, 0x47, 0x41, 0x54, 0x45, 0x57, 0x41, 0x59, 0x5f, 0x43, 0x4d, 0x44, 0x10,
	0xfa, 0x01, 0x12, 0x14, 0x0a, 0x0f, 0x53, 0x59, 0x53, 0x54, 0x45, 0x4d, 0x5f, 0x50, 0x52, 0x4f,
	0x54, 0x4f, 0x43, 0x4d, 0x44, 0x10, 0xff, 0x01,
}

var (
	file_xCmd_proto_rawDescOnce sync.Once
	file_xCmd_proto_rawDescData = file_xCmd_proto_rawDesc
)

func file_xCmd_proto_rawDescGZIP() []byte {
	file_xCmd_proto_rawDescOnce.Do(func() {
		file_xCmd_proto_rawDescData = protoimpl.X.CompressGZIP(file_xCmd_proto_rawDescData)
	})
	return file_xCmd_proto_rawDescData
}

var file_xCmd_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_xCmd_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_xCmd_proto_goTypes = []interface{}{
	(Command)(0),  // 0: Cmd.Command
	(*Nonce)(nil), // 1: Cmd.Nonce
}
var file_xCmd_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_xCmd_proto_init() }
func file_xCmd_proto_init() {
	if File_xCmd_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_xCmd_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Nonce); i {
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
			RawDescriptor: file_xCmd_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_xCmd_proto_goTypes,
		DependencyIndexes: file_xCmd_proto_depIdxs,
		EnumInfos:         file_xCmd_proto_enumTypes,
		MessageInfos:      file_xCmd_proto_msgTypes,
	}.Build()
	File_xCmd_proto = out.File
	file_xCmd_proto_rawDesc = nil
	file_xCmd_proto_goTypes = nil
	file_xCmd_proto_depIdxs = nil
}
