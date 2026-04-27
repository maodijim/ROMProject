package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
)

type EnchantCompare struct {
	Cmd.EnchantData
	Attrs [][]*EnchantAttrCompare
}

func (e *EnchantCompare) GetExtras() []*Cmd.EnchantExtra {
	return e.Extras
}

func (e *EnchantCompare) GetAttrs() [][]*EnchantAttrCompare {
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

func (g *GameConnection) EnchantPreviewContains(equipGuid string, preview *EnchantCompare, bothCondition bool, allAttrMustMatch bool) (success bool, targetNum uint32) {
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
		extraMatch := false
		attrMatch := false
		attrMatchCount := []map[any]bool{}
		for _, extra := range extrasPreview {
			for _, targetExtra := range preview.GetExtras() {
				if extra.GetBuffid() == targetExtra.GetBuffid() {
					extraMatch = true
				}
			}
		}
		// preset the target attr match count

		for attrIndex, targetAttrs := range preview.GetAttrs() {
			attrMatchCount = append(attrMatchCount, make(map[any]bool))
			for _, targetAttr := range targetAttrs {
				attrMatchCount[attrIndex][targetAttr.GetType()] = false
			}
		}
		for _, attr := range attrsPreview {
			for attrIndex, targetAttrs := range preview.GetAttrs() {
				for _, targetAttr := range targetAttrs {
					if attr.GetType() == targetAttr.GetType() {
						switch targetAttr.Condition {
						case ">":
							if attr.GetValue() > targetAttr.GetValue() {
								for _, attrNow := range attrsNow {
									if attrNow.GetType() == targetAttr.GetType() {
										if attrNow.GetValue() < targetAttr.GetValue() {
											attrMatch = true
											attrMatchCount[attrIndex][targetAttr.GetType()] = true
										} else if attr.GetValue() > attrNow.GetValue() {
											attrMatch = true
											attrMatchCount[attrIndex][targetAttr.GetType()] = true
										}
									}
									attrMatch = true
									attrMatchCount[attrIndex][targetAttr.GetType()] = true
								}
							}
						case "<":
							if attr.GetValue() < targetAttr.GetValue() {
								attrMatch = true
								attrMatchCount[attrIndex][targetAttr.GetType()] = true
							}
						case "=":
							if attr.GetValue() == targetAttr.GetValue() {
								attrMatch = true
								attrMatchCount[attrIndex][targetAttr.GetType()] = true
							}
						case ">=":
							if attr.GetValue() >= targetAttr.GetValue() {
								for _, attrNow := range attrsNow {
									if attrNow.GetType() == targetAttr.GetType() {
										if attrNow.GetValue() <= targetAttr.GetValue() {
											attrMatch = true
											attrMatchCount[attrIndex][targetAttr.GetType()] = true
										}
									} else {
										attrMatch = true
										attrMatchCount[attrIndex][targetAttr.GetType()] = true
									}
								}
							}
						case "<=":
							if attr.GetValue() <= targetAttr.GetValue() {
								attrMatch = true
								attrMatchCount[attrIndex][targetAttr.GetType()] = true
							}
						case "!=":
							if attr.GetValue() != targetAttr.GetValue() {
								attrMatch = true
								attrMatchCount[attrIndex][targetAttr.GetType()] = true
							}
						}
					}
				}
			}
		}
		if bothCondition {
			if extraMatch && attrMatch {
				if allAttrMustMatch {
					for _, countMap := range attrMatchCount {
						if utils.AllValuesTrue(countMap) {
							return true, targetNum
						}
					}
				} else {
					return true, targetNum
				}
			}
		} else if extraMatch || attrMatch {
			if allAttrMustMatch && attrMatch {
				for _, countMap := range attrMatchCount {
					if utils.AllValuesTrue(countMap) {
						return true, targetNum
					}
				}
			} else {
				return true, targetNum
			}
		}
	}
	return false, 0
}

func (g *GameConnection) EnchantContains(equipGuid string, preview *EnchantCompare, bothCondition bool, allAttrMustMatch bool) bool {
	enchant := g.EnchantGetByItemGuid(equipGuid, Cmd.EPackType_EPACKTYPE_EQUIP)
	if enchant == nil {
		return false
	}
	extras := enchant.GetExtras()
	attrs := enchant.GetAttrs()
	extraMatch := false
	attrMatch := false
	attrMatchCount := []map[any]bool{}
	for _, extra := range extras {
		for _, targetExtra := range preview.GetExtras() {
			if extra.GetBuffid() == targetExtra.GetBuffid() {
				extraMatch = true
			}
		}
	}
	// preset the target attr match count
	for attrIndex, targetAttrs := range preview.GetAttrs() {
		attrMatchCount = append(attrMatchCount, make(map[any]bool))
		for _, targetAttr := range targetAttrs {
			attrMatchCount[attrIndex][targetAttr.GetType()] = false
		}
	}
	for _, attr := range attrs {
		for attrIndex, targetAttrs := range preview.GetAttrs() {
			for _, targetAttr := range targetAttrs {
				if attr.GetType() == targetAttr.GetType() {
					switch targetAttr.Condition {
					case ">":
						if attr.GetValue() > targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					case "<":
						if attr.GetValue() < targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					case "=":
						if attr.GetValue() == targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					case ">=":
						if attr.GetValue() >= targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					case "<=":
						if attr.GetValue() <= targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					case "!=":
						if attr.GetValue() != targetAttr.GetValue() {
							attrMatch = true
							attrMatchCount[attrIndex][targetAttr.GetType()] = true
						}
					}
				}
			}
		}
	}
	if bothCondition {
		if extraMatch && attrMatch {
			if allAttrMustMatch {
				for _, countMap := range attrMatchCount {
					if utils.AllValuesTrue(countMap) {
						return true
					}
				}
			} else {
				return true
			}
		}
	} else if extraMatch || attrMatch {
		if allAttrMustMatch && attrMatch {
			for _, countMap := range attrMatchCount {
				if utils.AllValuesTrue(countMap) {
					return true
				}
			}
		} else {
			return true
		}
	}
	return false
}
