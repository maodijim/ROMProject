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

func (g *GameConnection) EnchantEquip(enchantType Cmd.EEnchantType, equipGuid string) {
	cmd := Cmd.EnchantEquip{
		Type: &enchantType,
		Guid: &equipGuid,
	}
	_ = g.sendProtoCmd(
		&cmd,
		SceneUserItemCmdId,
		Cmd.ItemParam_value["ITEMPARAM_ENCHANT"],
	)
}

func (g *GameConnection) EnchantSave(itemGuid string) {
	save := true
	cmd := Cmd.ProcessEnchantItemCmd{
		Itemid: &itemGuid,
		Save:   &save,
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

func (g *GameConnection) EnchantGetPreviewByItemGuid(itemGuid string, packType Cmd.EPackType) *Cmd.EnchantData {
	item := g.FindPackItemByGuid(itemGuid, packType)
	if item == nil {
		return &Cmd.EnchantData{}
	}
	return item.GetPreviewenchant()
}

func (g *GameConnection) EnchantPreviewContains(equipGuid string, preview *EnchantCompare) bool {
	enchantPreview := g.EnchantGetPreviewByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	enchantNow := g.EnchantGetByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	if enchantPreview == nil {
		return false
	}
	extrasPreview := enchantPreview.GetExtras()
	attrsPreview := enchantPreview.GetAttrs()
	attrsNow := enchantNow.GetAttrs()
	for _, extra := range extrasPreview {
		for _, targetExtra := range preview.GetExtras() {
			if extra.GetBuffid() == targetExtra.GetBuffid() {
				return true
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
									return true
								} else {
									return false
								}
							}
							return true
						}
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
						for _, attrNow := range attrsNow {
							if attrNow.GetType() == targetAttr.GetType() {
								if attrNow.GetValue() < targetAttr.GetValue() {
									return true
								} else {
									return false
								}
							} else {
								return true
							}
						}
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
