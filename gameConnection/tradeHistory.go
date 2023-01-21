package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

const (
	tradeHis = "tradeHist"
)

func (g *GameConnection) QueryTradeHistoryLog(pageIndex uint32) (result *Cmd.MyTradeLogRecordTradeCmd, err error) {
	cmd := &Cmd.MyTradeLogRecordTradeCmd{
		Charid: g.Role.RoleId,
	}
	if pageIndex != 0 {
		cmd.Index = &pageIndex
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["MY_TRADE_LOG_LIST_RECORDTRADE"],
	)
	err = g.waitForTradeHistory()

	// Retries
	if err == ErrQueryTimeout {
		log.Warnf("trade history timed out retrying %d", g.retries[tradeHis])
		if g.retries[tradeHis] < maxRetry {
			g.retries[tradeHis] += 1
			result, err = g.QueryTradeHistoryLog(pageIndex)
			if err == nil {
				g.retries[tradeHis] = 0
				return result, err
			}
		}
	} else {
		g.retries[tradeHis] = 0
	}

	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for trade history return error: %s", err)
	} else if g.tradeHistory != nil {
		result = g.tradeHistory
	}
	g.tradeHistory = &Cmd.MyTradeLogRecordTradeCmd{}
	g.Mutex.Unlock()
	return result, err
}

func (g *GameConnection) waitForTradeHistory() (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for trade history response")
		g.Mutex.Lock()
		if g.tradeHistory.GetLogList() == nil {
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

func (g *GameConnection) TakeLogTrade(logId uint64, logType Cmd.EOperType) {
	cmd := &Cmd.TakeLogCmd{
		Log: &Cmd.LogItemInfo{
			Id:      &logId,
			Logtype: &logType,
		},
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["TAKE_LOG_TRADE_PARAM"],
	)
}

func (g *GameConnection) QueryPendingSells() (result *Cmd.MyPendingListRecordTradeCmd) {
	result = &Cmd.MyPendingListRecordTradeCmd{}
	cmd := &Cmd.MyPendingListRecordTradeCmd{
		Charid: g.Role.RoleId,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"],
		Cmd.RecordUserTradeParam_value["MY_PENDING_LIST_RECORDTRADE"],
	)
	err := g.waitForMyPendingSellResp()
	g.Mutex.Lock()
	if err != nil {
		log.Errorf("query for trade history return error: %s", err)
	} else if g.pendingSells != nil {
		result = g.pendingSells
		g.pendingSells = &Cmd.MyPendingListRecordTradeCmd{}
	}
	g.Mutex.Unlock()
	return result
}

func (g *GameConnection) waitForMyPendingSellResp() (err error) {
	startQueryTime := time.Now()
	for start := startQueryTime; time.Since(start) < queryTimeout; {
		log.Debugf("Checking for trade history response")
		g.Mutex.Lock()
		if len(g.pendingSells.GetLists()) == 0 {
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

func (g *GameConnection) takeTradeLog(tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess {
			log.Infof("卖出 %d个%s 赚取 %d zeny",
				tradeLog.GetCount(),
				g.Items[tradeLog.GetItemid()].NameZh,
				tradeLog.GetGetmoney(),
			)
		} else {
			log.Infof("取回从%s购买的物品 %d个%s 花费 %d zeny",
				tradeLog.GetNameInfo().GetName(),
				tradeLog.GetCount(),
				g.Items[tradeLog.GetItemid()].NameZh,
				tradeLog.GetCostmoney(),
			)
		}
		g.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		time.Sleep(500 * time.Millisecond)
	}
}

func (g *GameConnection) takeFailedMoney(tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
		tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		log.Infof("取回抢购失败 %d个%s %d zeny", tradeLog.GetFailcount(), g.Items[tradeLog.GetItemid()].NameZh, tradeLog.GetRetmoney())
		g.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		newSilver := g.Role.GetSilver() + uint64(tradeLog.GetRetmoney())
		g.Role.Silver = &newSilver
		time.Sleep(500 * time.Millisecond)
	}
}

func (g *GameConnection) HandleTradeHistory(tradeHistory *Cmd.MyTradeLogRecordTradeCmd) {
	log.Infof("检查购买记录第%d页", tradeHistory.GetIndex())
	for _, tradeLog := range tradeHistory.GetLogList() {
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			log.Infof("回收抢购失败的金币")
			g.takeFailedMoney(tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuySuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NoramlBuy) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			g.takeTradeLog(tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			g.takeTradeLog(tradeLog)
		}
		// if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess ||
		//	tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell && itemName == "哈比的羽毛") {
		//	hasNewRecord = true
		//	excel.AddRecord(tradeLog, itemName)
		// }
	}
}
