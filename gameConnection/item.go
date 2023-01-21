package gameConnection

import (
	"errors"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	SceneUserItemCmdId = Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"]
)

func (g *GameConnection) UseItem(itemGuid string, count uint32) {

	cmd := &Cmd.ItemUse{
		Itemguid: &itemGuid,
	}
	if count > 0 {
		cmd.Count = &count
	}
	_ = g.sendProtoCmd(cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_ITEMUSE"])
}

func (g *GameConnection) GetItemCount(itemId uint32, source Cmd.ESource) {
	cmd := Cmd.GetCountItemCmd{
		Itemid: &itemId,
		Source: &source,
	}
	_ = g.sendProtoCmd(&cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_GETCOUNT"])
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

func (g *GameConnection) EquipItem(guid string, pos Cmd.EEquipPos, oper Cmd.EEquipOper) {
	cmd := &Cmd.Equip{
		Oper: &oper,
		Pos:  &pos,
		Guid: &guid,
	}
	_ = g.sendProtoCmd(cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_EQUIP"])
}

func (g *GameConnection) FindItemNameById(itemId uint32) string {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	if val, ok := g.Items[itemId]; ok {
		return val.NameZh
	}
	return ""
}

func (g *GameConnection) FindPackItemByName(name string, packType Cmd.EPackType) (itemData *Cmd.ItemData) {
	var itemId uint32
	g.Mutex.RLock()
	if val, ok := g.ItemsByName[name]; ok {
		for _, item := range val.Items {
			id, _ := item.Id.Int64()
			itemId = uint32(id)
			break
		}
	}
	g.Mutex.RUnlock()
	if itemId == 0 {
		log.Warnf("item name for id %s not found", name)
	}
	items := g.Role.GetPackItems()
	g.Mutex.RLock()
	for _, item := range items[packType] {
		if item.GetBase().GetId() == itemId {
			itemData = item
			break
		}
	}
	g.Mutex.RUnlock()
	return itemData
}

func (g *GameConnection) FindPackItemById(itemId uint32, packType Cmd.EPackType) (itemData *Cmd.ItemData) {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	packItem := g.Role.GetPackItems()
	if packItem == nil {
		return itemData
	}
	for _, item := range packItem[packType] {
		if item.GetBase().GetId() == itemId {
			itemData = item
			return itemData
		}
	}
	log.Warnf("item id %d not found", itemId)
	return itemData
}

func (g *GameConnection) EquipItemByName(name string, pos Cmd.EEquipPos, oper Cmd.EEquipOper) (err error) {
	itemInfo := g.FindPackItemByName(name, Cmd.EPackType_EPACKTYPE_MAIN).GetBase()
	if itemInfo == nil {
		return errors.New("item not found")
	}
	if pos == Cmd.EEquipPos_EEQUIPPOS_MIN {
		pos = g.GetItemEquipPos(itemInfo)
	}
	g.EquipItem(itemInfo.GetGuid(), pos, oper)
	return nil
}

func (g *GameConnection) FindPackItemByGuid(guid string, packType Cmd.EPackType) (itemData *Cmd.ItemData) {
	packItem := g.Role.GetPackItemsByType(packType)
	if packItem == nil {
		return itemData
	}
	for _, item := range packItem {
		if item.GetBase().GetGuid() == guid {
			itemData = item
			return itemData
		}
	}
	return itemData
}

func (g *GameConnection) GetItemEquipPos(item *Cmd.ItemInfo) Cmd.EEquipPos {
	switch item.GetEquipType() {
	case Cmd.EEquipType_EEQUIPTYPE_WEAPON:
		return Cmd.EEquipPos_EEQUIPPOS_WEAPON
	case Cmd.EEquipType_EEQUIPTYPE_SHIELD:
		return Cmd.EEquipPos_EEQUIPPOS_SHIELD
	case Cmd.EEquipType_EEQUIPTYPE_HEAD:
		return Cmd.EEquipPos_EEQUIPPOS_HEAD
	case Cmd.EEquipType_EEQUIPTYPE_ARMOUR:
		return Cmd.EEquipPos_EEQUIPPOS_ARMOUR
	case Cmd.EEquipType_EEQUIPTYPE_ACCESSORY:
		return Cmd.EEquipPos_EEQUIPPOS_ACCESSORY1
	case Cmd.EEquipType_EEQUIPTYPE_ROBE:
		return Cmd.EEquipPos_EEQUIPPOS_ROBE
	case Cmd.EEquipType_EEQUIPTYPE_SHOES:
		return Cmd.EEquipPos_EEQUIPPOS_SHOES
	}
	return Cmd.EEquipPos_EEQUIPPOS_MIN
}

func (g *GameConnection) UseFlyWing() {
	item := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		log.Warnf("fly wing not found")
		return
	}
	g.UseItem(item.GetBase().GetGuid(), 1)
}
