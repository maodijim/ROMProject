package gameConnection

import (
	"errors"
	"time"

	Cmd "ROMProject/Cmds"
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

func (g *GameConnection) GetItemCount(itemId uint32, source Cmd.ESource) (item *Cmd.GetCountItemCmd, err error) {
	cmd := Cmd.GetCountItemCmd{
		Itemid: &itemId,
		Source: &source,
	}
	g.AddNotifier("ITEMPARAM_GETCOUNT")
	_ = g.sendProtoCmd(&cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_GETCOUNT"])
	res, err := g.waitForResponse("ITEMPARAM_GETCOUNT")
	if err != nil {
		return item, err
	}
	item = res.(*Cmd.GetCountItemCmd)
	return item, nil
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
		g.logger.Infof("%s quick selling %d items", g.Role.GetRoleName(), len(sellItems))
		g.sendProtoCmd(cmd, SceneUserItemCmdId, Cmd.ItemParam_value["ITEMPARAM_QUICK_SELLITEM"])
	} else {
		g.logger.Infof("%s no quick sell items found", g.Role.GetRoleName())
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
		g.logger.Warnf("item name for id %s not found", name)
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

// FindPackItemById 根据物品ID查找背包中的物品, 返回第一个找到的物品
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
	g.logger.Warnf("item id %d not found", itemId)
	return itemData
}

// FindPackItemByIdAll FindPackItemById 根据物品ID查找背包中的物品, 返回所有找到的物品
func (g *GameConnection) FindPackItemByIdAll(itemId uint32, packType Cmd.EPackType) (itemDatas []*Cmd.ItemData) {
	g.Mutex.RLock()
	defer g.Mutex.RUnlock()
	packItem := g.Role.GetPackItems()
	if packItem == nil {
		return itemDatas
	}
	for _, item := range packItem[packType] {
		if item.GetBase().GetId() == itemId {
			itemDatas = append(itemDatas, item)
		}
	}
	if len(itemDatas) == 0 {
		g.logger.Warnf("item id %d not found", itemId)
		return itemDatas
	}
	return itemDatas
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
		g.logger.Warnf("fly wing not found")
		return
	}
	g.UseItem(item.GetBase().GetGuid(), 1)
	g.Role.DelaySkillTime = time.Now().Add(3 * time.Second)
}

func (g *GameConnection) UseYggdrasilBerry() {
	item := g.FindPackItemByName("天地树果实", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		g.logger.Warnf("Yggdrasil Berry not found")
		return
	}
	cdTime := time.UnixMilli(int64(item.GetBase().GetCd() - fixedItemCDSubtract))
	if time.Since(cdTime) < 0 {
		g.logger.Warnf("Yggdrasil Berry is in cooldown next use time: %v", cdTime)
		return
	}
	g.UseItem(item.GetBase().GetGuid(), 1)
	g.logger.Debugf("Used Yggdrasil Berry")
}

func (g *GameConnection) UseHoney() {
	item := g.FindPackItemById(12117, Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		g.logger.Warnf("Honey not found")
		return
	}
	cdTime := time.UnixMilli(int64(item.GetBase().GetCd() - fixedItemCDSubtract))
	if time.Since(cdTime) < 0 {
		g.logger.Warnf("Honey is in cooldown next use time: %v", cdTime)
		return
	}
	g.UseItem(item.GetBase().GetGuid(), 1)
	g.logger.Debugf("Used Honey")
}

func (g *GameConnection) ProduceItem(composeId uint32) {
	cmd := &Cmd.Produce{
		Composeid: &composeId,
	}
	_ = g.sendProtoCmd(cmd,
		Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"],
		Cmd.ItemParam_value["ITEMPARAM_PRODUCE"])
}
