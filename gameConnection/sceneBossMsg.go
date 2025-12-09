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

func (g *GameConnection) GetMvpInfoList() []Cmd.BossInfoItem {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()

	if len(g.BossInfo.GetBosslist()) > 0 {
		original := g.BossInfo.GetBosslist()
		copied := make([]Cmd.BossInfoItem, len(original))
		for i, v := range original {
			if v != nil {
				copied[i] = *v
			}
		}
		return copied
	}
	return nil
}

func (g *GameConnection) GetMiniInfoList() []Cmd.BossInfoItem {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()

	if len(g.BossInfo.GetMinilist()) > 0 {
		original := g.BossInfo.Minilist
		copied := make([]Cmd.BossInfoItem, len(original))
		for i, v := range original {
			if v != nil {
				copied[i] = *v
			}
		}
		return copied
	}
	return nil
}

func (g *GameConnection) GetBossInfoList() []Cmd.BossInfoItem {

	MvpInfoList := g.GetMvpInfoList()
	MiniInfoList := g.GetMiniInfoList()

	BossInfoList := append(MvpInfoList, MiniInfoList...)

	return BossInfoList
}

func (g *GameConnection) GetMvpInfoByName(name string) Cmd.BossInfoItem {

	MVPList := g.GetMvpInfoList()

	MonsterInfo, ok := g.MonsterItemsByName[name]

	if ok {
		for _, v := range MVPList {
			if int(*v.Id) == MonsterInfo.Id {
				return v
			}
		}
	}

	return Cmd.BossInfoItem{}
}

func (g *GameConnection) GetMiniInfoByName(name string) Cmd.BossInfoItem {

	MiniList := g.GetMiniInfoList()

	MonsterInfo, ok := g.MonsterItemsByName[name]

	if ok {
		for _, v := range MiniList {
			if int(*v.Id) == MonsterInfo.Id {
				return v
			}
		}
	}
	return Cmd.BossInfoItem{}
}

func (g *GameConnection) GetBossInfoByName(name string) Cmd.BossInfoItem {

	BossInfolist := g.GetBossInfoList()

	for _, v := range BossInfolist {
		MonsterInfo := g.GetMonsterItemById(*v.Id)

		if MonsterInfo.NameZh == name {
			return v
		}
	}

	return Cmd.BossInfoItem{}
}
