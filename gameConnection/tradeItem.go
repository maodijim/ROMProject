package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

const (
	MaxSellItems = 8
	buyRetries   = "buyRetries"
)

func (g *GameConnection) BuyItem(buyCount uint32, itemInfo *Cmd.TradeItemBaseInfo) (result *Cmd.BuyItemRecordTradeCmd, err error) {
	result = &Cmd.BuyItemRecordTradeCmd{}
	cmd := &Cmd.BuyItemRecordTradeCmd{
		Charid: g.Role.RoleId,
		Type:   itemInfo.Type,
		ItemInfo: &Cmd.TradeItemBaseInfo{
			Charid:      itemInfo.Charid,
			Itemid:      itemInfo.Itemid,
			Price:       itemInfo.Price,
			Count:       &buyCount,
			OrderId:     itemInfo.OrderId,
			PublicityId: itemInfo.PublicityId,
		},
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["BUY_ITEM_RECORDTRADE"],
	)
	err = g.waitForBuyItemRes(itemInfo.GetItemid())

	// Retries
	if err == ErrQueryTimeout {
		log.Warnf("buy item timed out retrying %d", g.retries[buyRetries])
		if g.retries[buyRetries] < maxRetry {
			g.retries[buyRetries] += 1
			result, err = g.BuyItem(buyCount, itemInfo)
			if err == nil {
				g.retries[buyRetries] = 0
				return result, err
			}
		}
	} else {
		g.retries[buyRetries] = 0
	}

	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for buy item return error: %s", err)
	} else if g.buyItem != nil && g.buyItem[itemInfo.GetItemid()] != nil {
		result = g.buyItem[itemInfo.GetItemid()]
		delete(g.buyItem, itemInfo.GetItemid())
	}
	g.Mutex.Unlock()
	return result, err
}

func (g *GameConnection) waitForBuyItemRes(itemId uint32) (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for buy item response")
		g.Mutex.Lock()
		if g.buyItem == nil || g.buyItem[itemId] == nil {
			g.Mutex.Unlock()
			time.Sleep(2 * time.Second)
			continue
		} else {
			g.Mutex.Unlock()
			break
		}
	}
	if time.Since(startQueryTime) > queryTimeout {
		err = ErrQueryTimeout
	}
	return err
}

func (g *GameConnection) ReqServerPrice(itemData *Cmd.ItemData) (result *Cmd.ReqServerPriceRecordTradeCmd) {
	cmd := &Cmd.ReqServerPriceRecordTradeCmd{
		Charid:   g.Role.RoleId,
		ItemData: itemData,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["REQ_SERVER_PRICE_RECORDTRADE"],
	)
	err := g.waitForReqServerPriceRes(itemData.GetBase().GetId())
	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for request server price return error: %s", err)
	} else if g.reqServerPrice != nil && g.reqServerPrice[itemData.GetBase().GetId()] != nil {
		result = g.reqServerPrice[itemData.GetBase().GetId()]
		delete(g.reqServerPrice, itemData.GetBase().GetId())
	}
	g.Mutex.Unlock()
	return result
}

func (g *GameConnection) waitForReqServerPriceRes(itemId uint32) (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("checking for request server price response")
		g.Mutex.Lock()
		if g.reqServerPrice == nil || g.reqServerPrice[itemId] == nil {
			g.Mutex.Unlock()
			time.Sleep(2 * time.Second)
			continue
		} else {
			g.Mutex.Unlock()
			break
		}
	}
	if time.Since(startQueryTime) > queryTimeout {
		err = ErrQueryTimeout
	}
	return err
}

func (g *GameConnection) SellItem(sellCount uint32, price *Cmd.ReqServerPriceRecordTradeCmd, itemData *Cmd.ItemData) (result *Cmd.SellItemRecordTradeCmd) {
	result = &Cmd.SellItemRecordTradeCmd{}
	itemId := itemData.GetBase().GetId()
	empty := uint32(0)
	source := Cmd.ESource_ESOURCE_NORMAL
	sellPrice := uint32(price.GetPrice())
	cmd := &Cmd.SellItemRecordTradeCmd{
		Charid: g.Role.RoleId,
		ItemInfo: &Cmd.TradeItemBaseInfo{
			Itemid:      &itemId,
			Price:       &sellPrice,
			Count:       &sellCount,
			Guid:        itemData.GetBase().Guid,
			PublicityId: &empty,
			ItemData: &Cmd.ItemData{
				Base: &Cmd.ItemInfo{
					Id:     itemData.GetBase().Id,
					Count:  &sellCount,
					Source: &source,
				},
				Equip:   &Cmd.EquipData{},
				Refine:  &Cmd.RefineData{},
				Enchant: &Cmd.EnchantData{},
				Egg:     &Cmd.EggData{},
				// Not available in EP 5.0
				// Attr:    &Cmd.GemAttrData{},
				Wedding: &Cmd.WeddingData{},
				Sender:  &Cmd.SenderData{},
				// Not available in EP 5.0
				// Skill:   &Cmd.GemSkillData{},
			},
		},
	}

	if price.GetStatetype() == Cmd.StateType_St_WillPublicity || price.GetStatetype() == Cmd.StateType_St_InPublicity {
		pId := uint32(1)
		cmd.ItemInfo.PublicityId = &pId
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["SELL_ITEM_RECORDTRADE"],
	)
	err := g.waitForSellItemRes(itemId)
	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for sell item return error: %s", err)
	} else if g.sellItem != nil && g.sellItem[itemId] != nil {
		result = g.sellItem[itemId]
		delete(g.sellItem, itemId)
	}
	g.Mutex.Unlock()
	return result
}

func (g *GameConnection) waitForSellItemRes(itemId uint32) (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for sell item response")
		g.Mutex.Lock()
		if g.sellItem == nil || g.sellItem[itemId] == nil {
			g.Mutex.Unlock()
			time.Sleep(2 * time.Second)
			continue
		} else {
			g.Mutex.Unlock()
			break
		}
	}
	if time.Since(startQueryTime) > queryTimeout {
		err = ErrQueryTimeout
	}
	return err
}
