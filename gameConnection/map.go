package gameConnection

type MapId uint32

const (
	MapId_Protera       MapId = 1
	MapId_ProteraSouth  MapId = 2
	MapId_ProteraWest   MapId = 5
	MapId_IzludeIsland  MapId = 7
	MapId_Geffen        MapId = 13
	MapId_GoblinForest  MapId = 33
	MapId_ProteraNorth  MapId = 42
	MapId_ProteraRoom1F MapId = 47
	MapId_PoringIsland  MapId = 52
	MapId_RoomAdvanced  MapId = 1001
)

func (m MapId) Uint32() uint32 {
	return uint32(m)
}
