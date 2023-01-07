package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
)

func (g *GameConnection) HandleRecordUserTradeProtoCmd(cmdParamId int32, rawData []byte) (err error) {
	var param proto.Message
	switch cmdParamId {
	case Cmd.RecordUserTradeParam_value["MY_PENDING_LIST_RECORDTRADE"]:
		param = &Cmd.MyPendingListRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		pRes := param.(*Cmd.MyPendingListRecordTradeCmd)
		if len(pRes.Lists) > 0 {
			g.Mutex.Lock()
			g.pendingSells = pRes
			g.Mutex.Unlock()
		}

	case Cmd.RecordUserTradeParam_value["SELL_ITEM_RECORDTRADE"]:
		param = &Cmd.SellItemRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		sellRes := param.(*Cmd.SellItemRecordTradeCmd)
		g.Mutex.Lock()
		g.sellItem[sellRes.GetItemInfo().GetItemid()] = sellRes
		g.Mutex.Unlock()

	case Cmd.RecordUserTradeParam_value["REQ_SERVER_PRICE_RECORDTRADE"]:
		param = &Cmd.ReqServerPriceRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		reqServerPrice := param.(*Cmd.ReqServerPriceRecordTradeCmd)
		if reqServerPrice.GetPrice() > 0 {
			g.Mutex.Lock()
			g.reqServerPrice[reqServerPrice.GetItemData().GetBase().GetId()] = reqServerPrice
			g.Mutex.Unlock()
		}

	case Cmd.RecordUserTradeParam_value["BUY_ITEM_RECORDTRADE"]:
		param = &Cmd.BuyItemRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		buyRes := param.(*Cmd.BuyItemRecordTradeCmd)
		g.Mutex.Lock()
		g.buyItem[buyRes.ItemInfo.GetItemid()] = buyRes
		g.Mutex.Unlock()

	case Cmd.RecordUserTradeParam_value["MY_TRADE_LOG_LIST_RECORDTRADE"]:
		param = &Cmd.MyTradeLogRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		history := param.(*Cmd.MyTradeLogRecordTradeCmd)
		if len(history.GetLogList()) > 0 {
			g.Mutex.Lock()
			g.tradeHistory = history
			g.Mutex.Unlock()
		}

	case Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"]:
		param = &Cmd.DetailPendingListRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		detail := param.(*Cmd.DetailPendingListRecordTradeCmd)
		if detail.GetSearchCond() != nil {
			g.Mutex.Lock()
			g.tradeDetail[detail.GetSearchCond().GetItemId()] = detail
			g.Mutex.Unlock()
		}

	case Cmd.RecordUserTradeParam_value["BRIEF_PENDING_LIST_RECORDTRADE"]:
		param = &Cmd.BriefPendingListRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		brief := param.(*Cmd.BriefPendingListRecordTradeCmd)
		if brief.GetCategory() != 0 {
			g.Mutex.Lock()
			g.tradeBrief[brief.GetCategory()] = brief
			g.Mutex.Unlock()
		}

	case Cmd.RecordUserTradeParam_value["ITEM_SELL_INFO_RECORDTRADE"]:
		param = &Cmd.ItemSellInfoRecordTradeCmd{}
		err = utils.ParseCmd(rawData, param)
		sellInfo := param.(*Cmd.ItemSellInfoRecordTradeCmd)
		if sellInfo.GetItemid() != 0 {
			g.Mutex.Lock()
			g.sellInfo[sellInfo.GetItemid()] = sellInfo
			g.Mutex.Unlock()
		}
	}
	return err
}
