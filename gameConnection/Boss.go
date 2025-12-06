package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

func (g *GameConnection) GetMvpInfoList() []*Cmd.BossInfoItem {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if g.BossInfo != nil && g.BossInfo.Bosslist != nil {
		original := g.BossInfo.Bosslist
		copied := make([]*Cmd.BossInfoItem, len(original))
		copy(copied, original)
		return copied
	}
	return nil
}

func (g *GameConnection) GetMiniInfoList() []*Cmd.BossInfoItem {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	if g.BossInfo != nil && g.BossInfo.Minilist != nil {
		original := g.BossInfo.Minilist
		copied := make([]*Cmd.BossInfoItem, len(original))
		copy(copied, original)
		return copied
	}
	return nil
}

func (g *GameConnection) GetBossInfoList() []*Cmd.BossInfoItem {

	MvpInfoList := g.GetMvpInfoList()
	MiniInfoList := g.GetMiniInfoList()

	BossInfoList := append(MvpInfoList, MiniInfoList...)

	return BossInfoList
}

func (g *GameConnection) GetMvpInfoByNmae(name string) *Cmd.BossInfoItem {

	MVPList := g.GetMvpInfoList()

	MonsterInfo, ok := g.MonsterItemsByName[name]

	if ok {
		for _, v := range MVPList {
			if int(*v.Id) == MonsterInfo.Id {
				return v
			}
		}
	}
	return nil
}

func (g *GameConnection) GetMiniInfoByNmae(name string) *Cmd.BossInfoItem {

	MiniList := g.GetMiniInfoList()

	MonsterInfo, ok := g.MonsterItemsByName[name]

	if ok {
		for _, v := range MiniList {
			if int(*v.Id) == MonsterInfo.Id {
				return v
			}
		}
	}
	return nil
}

func (g *GameConnection) GetBossInfoByNmae(name string) *Cmd.BossInfoItem {

	BossInfolist := g.GetBossInfoList()

	for _, v := range BossInfolist {
		MonsterInfo := g.GetMonsterItemById(*v.Id)

		if MonsterInfo.NameZh == name {
			return v
		}
	}

	return nil
}
