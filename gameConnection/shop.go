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
	g.AddNotifier("SHOPPARAM_QUERY_SHOP_CONFIG")
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserShopCmdId,
		Cmd.ShopParam_value["SHOPPARAM_QUERY_SHOP_CONFIG"],
	)
	res, err := g.waitForResponse("SHOPPARAM_QUERY_SHOP_CONFIG")
	if err != nil {
		return nil, err
	}
	return res.(*Cmd.QueryShopConfigCmd), nil
}

func (g *GameConnection) BuyShopItem(shopItem *Cmd.ShopItem, count uint32) {
	price := shopItem.GetMoneycount()
	id := shopItem.GetId()
	g.AddNotifier("SHOPPARAM_BUYITEM")
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
	res, err := g.waitForResponse("SHOPPARAM_BUYITEM")
	if err != nil {
		return
	}
	log.Infof("buy shop item %v", res)
}
