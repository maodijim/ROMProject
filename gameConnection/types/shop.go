package gameConnection

type ShopType uint32

const (
	ShopType_Item    ShopType = 600
	ShopType_Lottery ShopType = 650
	ShopType_Pringle ShopType = 3005
)

func (s ShopType) Uint32() uint32 {
	return uint32(s)
}

var TradeZhCategories = map[string]uint32{
	"图纸":      12,
	"药剂/效果":   1001,
	"精炼":      1002,
	"卷轴/唱片":   1003,
	"材料":      1004,
	"节日材料":    1005,
	"宠物材料":    1007,
	"宠物头饰图纸":  1008,
	"宠物头饰":    1009,
	"卡片 - 武器": 1010,
	"卡片 - 副手": 1011,
	"卡片 - 盔甲": 1012,
	"卡片 - 披风": 1013,
	"卡片 - 鞋子": 1014,
	"卡片 - 饰品": 1015,
	"卡片 - 头部": 1016,
	"武器":      1025,
	"副手":      1026,
	"盔甲":      1027,
	"披风":      1028,
	"鞋子":      1029,
	"饰品":      1030,
	"时装":      1045,
	"限定特典":    1052,
	"坐騎":      1020,
}
