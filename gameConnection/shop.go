package gameConnection

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"

	log "github.com/sirupsen/logrus"
)

var (
	sessionUserShopCmdId = Cmd.Command_value["SESSION_USER_SHOP_PROTOCMD"]
)

func (g *GameConnection) QueryShopConfig(shopType gameTypes.ShopType, shopId uint32) (result *Cmd.QueryShopConfigCmd, err error) {
	cmd := Cmd.QueryShopConfigCmd{
		Type:   (*uint32)(&shopType),
		Shopid: &shopId,
	}
	g.AddNotifier(gameTypes.NtfType_ShopQueryShopConfig)
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserShopCmdId,
		Cmd.ShopParam_value["SHOPPARAM_QUERY_SHOP_CONFIG"],
	)
	res, err := g.waitForResponse(gameTypes.NtfType_ShopQueryShopConfig)
	if err != nil {
		return nil, err
	}
	return res.(*Cmd.QueryShopConfigCmd), nil
}

func (g *GameConnection) BuyShopItem(shopItem *Cmd.ShopItem, count uint32) {
	price := shopItem.GetMoneycount()
	id := shopItem.GetId()
	g.AddNotifier(gameTypes.NtfType_ShopBuyItem)
	cmd := Cmd.BuyShopItem{
		Price: &price,
		Count: &count,
		Id:    &id,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserShopCmdId,
		Cmd.ShopParam_value["SHOPPARAM_BUYITEM"],
	)
	res, err := g.waitForResponse(gameTypes.NtfType_ShopBuyItem)
	if err != nil {
		return
	}
	log.Infof("buy shop item %v", res)
}

func (g *GameConnection) QueryPringleShopConfig() (*Cmd.QueryShopConfigCmd, error) {
	return g.QueryShopConfig(gameTypes.ShopType_Pringle, 1)
}

func (g *GameConnection) QueryQuickBuyShopConfig(itemId uint32, itemid ...uint32) (res *Cmd.QueryQuickBuyConfigCmd, err error) {
	cmd := Cmd.QueryQuickBuyConfigCmd{
		Itemids: append([]uint32{itemId}, itemid...),
	}
	g.AddNotifier(gameTypes.NtfType_ShopQueryQuickBuyConfig)
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserShopCmdId,
		Cmd.ShopParam_value["SHOPPARAM_QUICKBUY_SHOP_CONFIG]"],
	)
	r, err := g.waitForResponse(gameTypes.NtfType_ShopQueryQuickBuyConfig)
	if err != nil {
		return nil, err
	}
	return r.(*Cmd.QueryQuickBuyConfigCmd), nil
}
