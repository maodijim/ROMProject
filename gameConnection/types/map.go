package gameConnection

type MapId uint32

const (
	MapId_Protera               MapId = 1
	MapId_ProteraSouth          MapId = 2
	MapId_ProteraWest           MapId = 5
	MapId_LabyrinthForest       MapId = 6
	MapId_IzludeIsland          MapId = 7
	MapId_SunkenShip            MapId = 8
	MapId_GhostShip             MapId = 9
	MapID_ByalanIsland          MapId = 10
	MapID_UnderseaTemple        MapId = 12
	MapId_Geffen                MapId = 13
	MapId_MjolnirMountains      MapId = 14
	MapId_Morroc                MapId = 16
	MapId_Pyramid1F             MapId = 17
	MapId_Payon                 MapId = 18
	MapId_PayonSouth            MapId = 19
	MapId_OrcVillage            MapId = 23
	MapId_GlastHeimOutskirt     MapId = 26
	MapId_GlastHeim             MapId = 27
	MapId_PayonForest           MapId = 32
	MapId_GoblinForest          MapId = 33
	MapId_KordtForest           MapId = 34
	MapId_SogratDesert          MapId = 35
	MapId_ProteraNorth          MapId = 42
	MapId_AlDeBaran             MapId = 43
	MapId_ProteraRoom1F         MapId = 47
	MapId_GingerbreadCity       MapId = 48
	MapId_ToyFactory1F          MapId = 49
	MapId_PoringIsland          MapId = 52
	MapId_PayonForestSouth      MapId = 54
	MapId_OrcVillageSouth       MapId = 56
	MapId_GlastHeimOutskirts    MapId = 59
	MapId_Amatsu                MapId = 62
	MapId_Yuno                  MapId = 63
	MapId_BorderCheckpoint      MapId = 64
	MapId_EinbrochField         MapId = 65
	MapId_MagmaDungeon1F        MapId = 66
	MapId_MagmaDungeon2F        MapId = 67
	MapId_MagmaDungeon3F        MapId = 68
	MapId_LesterLighthouse      MapId = 69
	MapId_Niflheim              MapId = 70
	MapId_MistyForest           MapId = 71
	MapId_Skellington           MapId = 72
	MapId_Hamelin               MapId = 73
	MapId_Umbala                MapId = 75
	MapId_Lighthalzen           MapId = 76
	MapId_LighthalzenPlain      MapId = 77
	MapId_LhzDun01              MapId = 78
	MapId_LhzDun02              MapId = 79
	MapId_LhzDun03              MapId = 81
	MapId_ThePlainofIda         MapId = 83
	MapId_Rachel                MapId = 85
	MapId_Lasagna               MapId = 91
	MapId_DoradoIsland          MapId = 92
	MapId_RavioliForest         MapId = 93
	MapId_Luoyang               MapId = 99
	MapId_SunsetBeach           MapId = 101
	MapId_Wasteland             MapId = 102
	MapId_MoonLake              MapId = 106
	MapId_Eclage                MapId = 109
	MapId_TimeGarden            MapId = 110
	MapId_CrypturaAcademy       MapId = 113
	MapId_StarTearsForest       MapId = 114
	MapId_BloomingLand          MapId = 123
	MapId_WindBreath            MapId = 125
	MapId_Comodo                MapId = 126
	MapId_KokomoBeach           MapId = 127
	MapId_MeteorForest          MapId = 140
	MapId_Alberta               MapId = 145
	MapId_TurtleIsland          MapId = 146
	MapId_TearsoftheAncientCity MapId = 149
	MapId_ScJfzc001             MapId = 151
	MapId_AbyssalLake           MapId = 154
	MapId_RoomAdvanced          MapId = 1001
	MapId_RoyalCooking          MapId = 1061
	MapId_Guild                 MapId = 10001
)

var (
	MapIdMap = map[uint32]MapId{
		1:     MapId_Protera,
		2:     MapId_ProteraSouth,
		5:     MapId_ProteraWest,
		6:     MapId_LabyrinthForest,
		7:     MapId_IzludeIsland,
		8:     MapId_SunkenShip,
		9:     MapId_GhostShip,
		10:    MapID_ByalanIsland,
		12:    MapID_UnderseaTemple,
		13:    MapId_Geffen,
		14:    MapId_MjolnirMountains,
		16:    MapId_Morroc,
		17:    MapId_Pyramid1F,
		18:    MapId_Payon,
		19:    MapId_PayonSouth,
		23:    MapId_OrcVillage,
		26:    MapId_GlastHeimOutskirt,
		27:    MapId_GlastHeim,
		32:    MapId_PayonForest,
		33:    MapId_GoblinForest,
		34:    MapId_KordtForest,
		35:    MapId_SogratDesert,
		42:    MapId_ProteraNorth,
		43:    MapId_AlDeBaran,
		47:    MapId_ProteraRoom1F,
		48:    MapId_GingerbreadCity,
		49:    MapId_ToyFactory1F,
		52:    MapId_PoringIsland,
		54:    MapId_PayonForestSouth,
		56:    MapId_OrcVillageSouth,
		59:    MapId_GlastHeimOutskirts,
		62:    MapId_Amatsu,
		63:    MapId_Yuno,
		64:    MapId_BorderCheckpoint,
		65:    MapId_EinbrochField,
		66:    MapId_MagmaDungeon1F,
		67:    MapId_MagmaDungeon2F,
		68:    MapId_MagmaDungeon3F,
		69:    MapId_LesterLighthouse,
		70:    MapId_Niflheim,
		71:    MapId_MistyForest,
		72:    MapId_Skellington,
		73:    MapId_Hamelin,
		75:    MapId_Umbala,
		76:    MapId_Lighthalzen,
		77:    MapId_LighthalzenPlain,
		78:    MapId_LhzDun01,
		79:    MapId_LhzDun02,
		81:    MapId_LhzDun03,
		83:    MapId_ThePlainofIda,
		85:    MapId_Rachel,
		91:    MapId_Lasagna,
		92:    MapId_DoradoIsland,
		93:    MapId_RavioliForest,
		99:    MapId_Luoyang,
		101:   MapId_SunsetBeach,
		102:   MapId_Wasteland,
		106:   MapId_MoonLake,
		109:   MapId_Eclage,
		110:   MapId_TimeGarden,
		113:   MapId_CrypturaAcademy,
		114:   MapId_StarTearsForest,
		123:   MapId_BloomingLand,
		125:   MapId_WindBreath,
		126:   MapId_Comodo,
		127:   MapId_KokomoBeach,
		140:   MapId_MeteorForest,
		145:   MapId_Alberta,
		146:   MapId_TurtleIsland,
		149:   MapId_TearsoftheAncientCity,
		151:   MapId_ScJfzc001,
		154:   MapId_AbyssalLake,
		1001:  MapId_RoomAdvanced,
		1061:  MapId_RoyalCooking,
		10001: MapId_Guild,
	}
)
var MapNameZh = map[string]MapId{
	"普隆德拉":      MapId_Protera,
	"普隆德拉南门":    MapId_ProteraSouth,
	"普隆德拉西门":    MapId_ProteraWest,
	"迷藏森林":      MapId_LabyrinthForest,
	"伊斯鲁得岛":     MapId_IzludeIsland,
	"沉船":        MapId_SunkenShip,
	"幽灵船":       MapId_GhostShip,
	"海底洞窟岛":     MapID_ByalanIsland,
	"海底神殿":      MapID_UnderseaTemple,
	"吉芬":        MapId_Geffen,
	"妙勒尼山脉":     MapId_MjolnirMountains,
	"摩洛克":       MapId_Morroc,
	"金字塔 1F":    MapId_Pyramid1F,
	"斐扬":        MapId_Payon,
	"斐扬南门":      MapId_PayonSouth,
	"兽人村落":      MapId_OrcVillage,
	"古城郊外":      MapId_GlastHeimOutskirt,
	"古城":        MapId_GlastHeim,
	"斐扬森林":      MapId_PayonForest,
	"哥布林森林":     MapId_GoblinForest,
	"科德森林":      MapId_KordtForest,
	"苏克拉特沙漠":    MapId_SogratDesert,
	"普隆德拉北门":    MapId_ProteraNorth,
	"阿尔德巴朗":     MapId_AlDeBaran,
	"普隆德拉大厅 1F": MapId_ProteraRoom1F,
	"姜饼城":       MapId_GingerbreadCity,
	"玩具工厂 1F":   MapId_ToyFactory1F,
	"波利岛":       MapId_PoringIsland,
	"斐扬森林南部":    MapId_PayonForestSouth,
	"兽人村落南部":    MapId_OrcVillageSouth,
	"古城外围":      MapId_GlastHeimOutskirts,
	"天水之国·安塔修":  MapId_Amatsu,
	"尤诺":        MapId_Yuno,
	"边境检查站":     MapId_BorderCheckpoint,
	"艾因布洛克原野":   MapId_EinbrochField,
	"熔岩洞窟 1F":   MapId_MagmaDungeon1F,
	"熔岩洞窟 2F":   MapId_MagmaDungeon2F,
	"熔岩洞窟 3F":   MapId_MagmaDungeon3F,
	"莱斯特灯塔":     MapId_LesterLighthouse,
	"尼芙海姆":      MapId_Niflheim,
	"迷雾森林":      MapId_MistyForest,
	"骷髅洞穴":      MapId_Skellington,
	"哈姆林":       MapId_Hamelin,
	"乌帕拉":       MapId_Umbala,
	"里希塔尔岑":     MapId_Lighthalzen,
	"里希塔尔岑平原":   MapId_LighthalzenPlain,
	"生体地下1F":    MapId_LhzDun01,
	"生体地下2F":    MapId_LhzDun02,
	"生体地下3F":    MapId_LhzDun03,
	"伊达平原":      MapId_ThePlainofIda,
	"瑞秋":        MapId_Rachel,
	"拉萨尼亚":      MapId_Lasagna,
	"多拉多岛":      MapId_DoradoIsland,
	"意大利饺森林":    MapId_RavioliForest,
	"洛阳":        MapId_Luoyang,
	"夕阳海岸":      MapId_SunsetBeach,
	"荒境":        MapId_Wasteland,
	"月之湖":       MapId_MoonLake,
	"伊克莱基":      MapId_Eclage,
	"时间花园":      MapId_TimeGarden,
	"克雷普特学院":    MapId_CrypturaAcademy,
	"星泪森林":      MapId_StarTearsForest,
	"绽放之地":      MapId_BloomingLand,
	"风之森":       MapId_WindBreath,
	"科摩多":       MapId_Comodo,
	"可可蒙海滩":     MapId_KokomoBeach,
	"流星森林":      MapId_MeteorForest,
	"阿尔贝塔":      MapId_Alberta,
	"海龟岛":       MapId_TurtleIsland,
	"古城之泪":      MapId_TearsoftheAncientCity,
	"副本·极限挑战":   MapId_ScJfzc001,
	"深渊之湖":      MapId_AbyssalLake,

	// 特殊地图
	"高级房间":  MapId_RoomAdvanced,
	"皇家料理间": MapId_RoyalCooking,
	"公会领地":  MapId_Guild,
}

func (m MapId) Uint32() uint32 {
	return uint32(m)
}
