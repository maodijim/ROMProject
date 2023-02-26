package gameConnection

type MapId uint32

const (
	MapId_Protera          MapId = 1
	MapId_ProteraSouth     MapId = 2
	MapId_ProteraWest      MapId = 5
	MapId_IzludeIsland     MapId = 7
	MapId_Geffen           MapId = 13
	MapId_MjolnirMountains MapId = 14
	MapId_GoblinForest     MapId = 33
	MapId_ProteraNorth     MapId = 42
	MapId_ProteraRoom1F    MapId = 47
	MapId_GingerbreadCity  MapId = 48
	MapId_ToyFactory1F     MapId = 49
	MapId_PoringIsland     MapId = 52
	MapId_Yuno             MapId = 63
	MapId_EinbrochField    MapId = 65
	MapId_MagmaDungeon1F   MapId = 66
	MapId_MagmaDungeon2F   MapId = 67
	MapId_MagmaDungeon3F   MapId = 68
	MapId_Niflheim         MapId = 70
	MapId_MistyForest      MapId = 71
	MapId_Skellington      MapId = 72
	MapId_Hamelin          MapId = 73
	MapId_RoomAdvanced     MapId = 1001
	MapId_RoyalCooking     MapId = 1061
	MapId_Guild            MapId = 10001
)

var (
	MapIdMap = map[uint32]MapId{
		1:     MapId_Protera,
		2:     MapId_ProteraSouth,
		5:     MapId_ProteraWest,
		7:     MapId_IzludeIsland,
		13:    MapId_Geffen,
		14:    MapId_MjolnirMountains,
		33:    MapId_GoblinForest,
		42:    MapId_ProteraNorth,
		47:    MapId_ProteraRoom1F,
		48:    MapId_GingerbreadCity,
		49:    MapId_ToyFactory1F,
		52:    MapId_PoringIsland,
		63:    MapId_Yuno,
		65:    MapId_EinbrochField,
		66:    MapId_MagmaDungeon1F,
		67:    MapId_MagmaDungeon2F,
		68:    MapId_MagmaDungeon3F,
		1001:  MapId_RoomAdvanced,
		1061:  MapId_RoyalCooking,
		10001: MapId_Guild,
	}
)

func (m MapId) Uint32() uint32 {
	return uint32(m)
}
