package gameConnection

import (
	Cmd "ROMProject/Cmds"
	notifier "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"github.com/golang/protobuf/proto"
)

func (g *GameConnection) HandleSceneBossMsg(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.BossParam_value["BOSS_LIST_USER_CMD"]:
		param = &Cmd.BossListUserCmd{}
		err = utils.ParseCmd(rawData, param)
		if err == nil && g.Notifier(notifier.NtfType_BossListUserCmd) != nil {
			go func() {
				g.Notifier(notifier.NtfType_BossListUserCmd) <- param.(*Cmd.BossListUserCmd)
			}()
		}
		if err == nil {
			g.Mutex.Lock()
			g.BossInfo = param.(*Cmd.BossListUserCmd)
			g.Mutex.Unlock()
		}

	case Cmd.BossParam_value["BOSS_WORLD_NTF"]:
		param = &Cmd.WorldBossNtf{}
		err := utils.ParseCmd(rawData, param)
		if err == nil && g.Notifier(notifier.NtfType_BossWorldNtf) != nil {
			g.Notifier(notifier.NtfType_BossWorldNtf) <- param.(*Cmd.WorldBossNtf)
		}
	}
	return param, err
}

func (g *GameConnection) GetBossInfo() chan interface{} {
	cmd := &Cmd.BossListUserCmd{}
	g.AddNotifier(notifier.NtfType_BossListUserCmd)
	ntf := g.Notifier(notifier.NtfType_BossListUserCmd)
	_ = g.sendProtoCmd(cmd, Cmd.Command_value["SCENE_BOSS_PROTOCMD"], Cmd.BossParam_value["BOSS_LIST_USER_CMD"])
	return ntf
}
