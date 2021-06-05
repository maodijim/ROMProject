package gameConnection

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	SceneUserItemCmdId = Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"]
)

func (g *GameConnection) UseItem() {

}

func (g *GameConnection) IsQuickSellItem(itemId uint32) bool {
	if _, ok := g.ExchangeItems[itemId]; !ok && g.Items[itemId].GetLevel() != 0 {
		return true
	}
	return false
}

func (g *GameConnection) QuickSellItems() {
	var sellItems []*Cmd.SItem
	for _, pack := range g.Role.PackItems {
		for _, item := range pack {
			if g.IsQuickSellItem(item.GetBase().GetId()) {
				sellItems = append(sellItems, &Cmd.SItem{
					Guid:  item.GetBase().Guid,
					Count: item.GetBase().Count,
				})
			}
		}
	}
	cmd := &Cmd.QuickSellItemCmd{
		Items: sellItems,
	}
	if len(sellItems) > 0 {
		log.Infof("%s quick selling %d items", g.Role.GetRoleName(), len(sellItems))
		g.sendProtoCmd(cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_QUICK_SELLITEM"])
	} else {
		log.Infof("%s no quick sell items found", g.Role.GetRoleName())
	}
}

func (g *GameConnection) GetTempItems() {
	op := Cmd.EEquipOper_EEQUIPOPER_OFFTEMP
	cmd := &Cmd.Equip{
		Oper: &op,
	}
	g.sendProtoCmd(cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_EQUIP"])
}
