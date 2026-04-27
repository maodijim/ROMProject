package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

var LotteryNameZh = map[string]Cmd.ELotteryType{
	"幻想创造器·宴": Cmd.ELotteryType_ELotteryType_Magic,
	"卡片扭蛋":    Cmd.ELotteryType_ELotteryType_Card,
	"幻想创造器Ⅲ号": Cmd.ELotteryType_ELotteryType_Card,
	"装备扭蛋":    Cmd.ELotteryType_ELotteryType_Equip,
	"头饰扭蛋":    Cmd.ELotteryType_ELotteryType_Head,
	"坐骑扭蛋":    Cmd.ELotteryType_ELotteryType_Max,
	"宠物蛋扭蛋":   Cmd.ELotteryType_ELotteryType_PetEgg,
}

var LotteryTypeToNameMap = map[Cmd.ELotteryType]string{
	Cmd.ELotteryType_ELotteryType_Magic:  "幻想创造器·宴",
	Cmd.ELotteryType_ELotteryType_Card:   "卡片扭蛋",
	Cmd.ELotteryType_ELotteryType_Equip:  "装备扭蛋",
	Cmd.ELotteryType_ELotteryType_Head:   "头饰扭蛋",
	Cmd.ELotteryType_ELotteryType_Max:    "坐骑扭蛋",
	Cmd.ELotteryType_ELotteryType_PetEgg: "宠物蛋扭蛋",
}

var LotteryTypeToVarCountMap = map[Cmd.ELotteryType]Cmd.EVarType{
	Cmd.ELotteryType_ELotteryType_Magic: Cmd.EVarType_EVARTYPE_DAY_LOTTERY_CNT_MAGIC,
}

var LotteryTypePriceMap = map[Cmd.ELotteryType]uint64{
	Cmd.ELotteryType_ELotteryType_Magic: 500,
	Cmd.ELotteryType_ELotteryType_Card:  500,
}
