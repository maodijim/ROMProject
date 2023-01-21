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
	MapId_ToyFactory1F     MapId = 49
	MapId_PoringIsland     MapId = 52
	MapId_Yuno             MapId = 63
	MapId_EinbrochField    MapId = 65
	MapId_MagmaDungeon1F   MapId = 66
	MapId_MagmaDungeon2F   MapId = 67
	MapId_MagmaDungeon3F   MapId = 68
	MapId_RoomAdvanced     MapId = 1001
	MapId_RoyalCooking     MapId = 1061
)

func (m MapId) Uint32() uint32 {
	return uint32(m)
}
