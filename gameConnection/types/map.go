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

func (m MapId) Uint32() uint32 {
	return uint32(m)
}
