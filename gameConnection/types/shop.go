package gameConnection

type ShopType uint32

const (
	ShopType_Item    ShopType = 600
	ShopType_Lottery ShopType = 650
)

func (s ShopType) Uint32() uint32 {
	return uint32(s)
}
