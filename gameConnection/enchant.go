package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

type EnchantCompare struct {
	Cmd.EnchantData
	Attrs []*EnchantAttrCompare
}

func (e *EnchantCompare) GetExtras() []*Cmd.EnchantExtra {
	return e.Extras
}

func (e *EnchantCompare) GetAttrs() []*EnchantAttrCompare {
	return e.Attrs
}

type EnchantAttrCompare struct {
	Cmd.EnchantAttr
	Condition string
}

func (ea *EnchantAttrCompare) GetType() Cmd.EAttrType {
	return ea.EnchantAttr.GetType()
}

func (ea *EnchantAttrCompare) GetValue() uint32 {
	return ea.EnchantAttr.GetValue()
}

func (ea *EnchantAttrCompare) GetCondition() string {
	return ea.Condition
}

func (g *GameConnection) EnchantEquip(enchantType Cmd.EEnchantType, equipGuid string, enchantCount uint32) {
	defaultEnchantCount := uint32(1)
	cmd := Cmd.EnchantEquip{
		Type:         &enchantType,
		Guid:         &equipGuid,
		EnchantCount: &defaultEnchantCount,
	}
	if enchantCount >= 1 {
		cmd.EnchantCount = &enchantCount
	}
	_ = g.sendProtoCmd(
		&cmd,
		SceneUserItemCmdId,
		Cmd.ItemParam_value["ITEMPARAM_ENCHANT"],
	)
}

func (g *GameConnection) EnchantSave(itemGuid string, enchantNum uint32) {
	save := true
	cmd := Cmd.ProcessEnchantItemCmd{
		Itemid: &itemGuid,
		Save:   &save,
	}
	if enchantNum > 0 {
		cmd.EnchantNum = &enchantNum
	}
	_ = g.sendProtoCmd(
		&cmd,
		SceneUserItemCmdId,
		Cmd.ItemParam_value["ITEMPARAM_PROCESSENCHANT"],
	)
}

func (g *GameConnection) EnchantGetByItemGuid(itemGuid string, packType Cmd.EPackType) *Cmd.EnchantData {
	item := g.FindPackItemByGuid(itemGuid, packType)
	if item == nil {
		return &Cmd.EnchantData{}
	}
	return item.GetEnchant()
}

func (g *GameConnection) EnchantGetPreviewByItemGuid(itemGuid string, packType Cmd.EPackType) []*Cmd.EnchantData {
	item := g.FindPackItemByGuid(itemGuid, packType)
	if item == nil {
		return []*Cmd.EnchantData{}
	}
	return item.GetPreviewenchant()
}

func (g *GameConnection) EnchantPreviewContains(equipGuid string, preview *EnchantCompare) (success bool, targetNum uint32) {
	enchantPreview := g.EnchantGetPreviewByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	enchantNow := g.EnchantGetByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	if enchantPreview == nil {
		return false, targetNum
	}
	for index, newPreview := range enchantPreview {
		targetNum = uint32(index)
		extrasPreview := newPreview.GetExtras()
		attrsPreview := newPreview.GetAttrs()
		attrsNow := enchantNow.GetAttrs()
		for _, extra := range extrasPreview {
			for _, targetExtra := range preview.GetExtras() {
				if extra.GetBuffid() == targetExtra.GetBuffid() {
					return true, targetNum
				}
			}
		}
		for _, attr := range attrsPreview {
			for _, targetAttr := range preview.GetAttrs() {
				if attr.GetType() == targetAttr.GetType() {
					switch targetAttr.Condition {
					case ">":
						if attr.GetValue() > targetAttr.GetValue() {
							for _, attrNow := range attrsNow {
								if attrNow.GetType() == targetAttr.GetType() {
									if attrNow.GetValue() < targetAttr.GetValue() {
										return true, targetNum
									} else if attr.GetValue() > attrNow.GetValue() {
										return true, targetNum
									} else {
										return false, 0
									}
								}
								return true, targetNum
							}
						}
					case "<":
						if attr.GetValue() < targetAttr.GetValue() {
							return true, targetNum
						}
					case "=":
						if attr.GetValue() == targetAttr.GetValue() {
							return true, targetNum
						}
					case ">=":
						if attr.GetValue() >= targetAttr.GetValue() {
							for _, attrNow := range attrsNow {
								if attrNow.GetType() == targetAttr.GetType() {
									if attrNow.GetValue() < targetAttr.GetValue() {
										return true, targetNum
									} else if attr.GetValue() > attrNow.GetValue() {
										return true, targetNum
									} else {
										return false, 0
									}
								} else {
									return true, targetNum
								}
							}
						}
					case "<=":
						if attr.GetValue() <= targetAttr.GetValue() {
							return true, targetNum
						}
					case "!=":
						if attr.GetValue() != targetAttr.GetValue() {
							return true, targetNum
						}
					}
				}
			}
		}
	}
	return false, 0
}

func (g *GameConnection) EnchantContains(equipGuid string, preview *EnchantCompare) bool {
	enchant := g.EnchantGetByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	if enchant == nil {
		return false
	}
	extras := enchant.GetExtras()
	attrs := enchant.GetAttrs()
	for _, extra := range extras {
		for _, targetExtra := range preview.GetExtras() {
			if extra.GetBuffid() == targetExtra.GetBuffid() {
				return true
			}
		}
	}
	for _, attr := range attrs {
		for _, targetAttr := range preview.GetAttrs() {
			if attr.GetType() == targetAttr.GetType() {
				switch targetAttr.Condition {
				case ">":
					if attr.GetValue() > targetAttr.GetValue() {
						return true
					}
				case "<":
					if attr.GetValue() < targetAttr.GetValue() {
						return true
					}
				case "=":
					if attr.GetValue() == targetAttr.GetValue() {
						return true
					}
				case ">=":
					if attr.GetValue() >= targetAttr.GetValue() {
						return true
					}
				case "<=":
					if attr.GetValue() <= targetAttr.GetValue() {
						return true
					}
				case "!=":
					if attr.GetValue() != targetAttr.GetValue() {
						return true
					}
				}
			}
		}
	}
	return false
}
