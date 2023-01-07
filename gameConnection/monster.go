package gameConnection

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
