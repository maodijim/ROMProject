package gameConnection

import (
	"ROMProject/utils"
)

func (g *GameConnection) GetMonsterIdByName(name string) (monsterId uint32) {
	if val, ok := g.MonsterItemsByName[name]; ok {
		monsterId = uint32(val.Id)
	}
	return monsterId
}

func (g *GameConnection) GetMonsterNameById(id uint32) (monsterName string) {
	if val, ok := g.MonsterItems[id]; ok {
		monsterName = val.NameZh
	}
	return monsterName
}

func (g *GameConnection) GetMonsterItemByName(name string) utils.MonsterInfo {
	if val, ok := g.MonsterItemsByName[name]; ok {
		return val
	}
	return utils.MonsterInfo{}
}

func (g *GameConnection) GetMonsterItemById(id uint32) utils.MonsterInfo {
	if val, ok := g.MonsterItems[id]; ok {
		return val
	}
	return utils.MonsterInfo{}
}
