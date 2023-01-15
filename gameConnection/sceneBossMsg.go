package gameConnection

import (
	Cmd "ROMProject/Cmds"
	notifier "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
)

func (g *GameConnection) HandleSceneBossMsg(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.BossParam_value["BOSS_WORLD_NTF"]:
		param = &Cmd.WorldBossNtf{}
		err := utils.ParseCmd(rawData, param)
		if err != nil && g.Notifier(notifier.NtfType_BossWorldNtf) != nil {
			g.Notifier(notifier.NtfType_BossWorldNtf) <- param.(*Cmd.WorldBossNtf)
		}
	}
	return param, err
}
