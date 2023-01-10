package main

import (
	"errors"
	"flag"
	"math"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	traderVer     = "0.1.6"
	pointDiscount = 0.75
	MaxSellItems  = 8
)

var (
	items          = &utils.ItemsLoader{}
	ErrNoItemFound = errors.New("no item found")
)

func init() {
	// log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	log.Infof("ROM auto trader Version: %s", traderVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	purchaseItemYml := flag.String("purchaseItems", "purchaseItems.yml", "yaml file of the list of items to purchase")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()

	items = utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	purchaseConfig := PurchaseConfigParser(*purchaseItemYml)
	excelFile := NewTradeExcel("", "Sheet1")
	gameConnect := gameConnection.NewConnection(conf, nil, items)
	if *enableDebug {
		gameConnect.DebugMsg = *enableDebug
		log.SetLevel(log.DebugLevel)
	}

	gameConnect.GameServerLogin()
	shouldReconnect := false
	for {
		purchaseConfig = PurchaseConfigParser(*purchaseItemYml)
		if gameConnect.Role.GetInGame() {
			if shouldReconnect && gameConnect.IsEnteringGameMap() {
				gameConnect.Reconnect()
				shouldReconnect = false
				continue
			}
			gameConnect.MoveChart(gameConnect.Role.GetPos())

			if !gameConnect.Role.GetAuthConfirm() {
				gameConnect.AuthConfirm(purchaseConfig.AuthPass)
			}

			// gameConnect.ChangeMap(uint32(1))
			// 获取书包和个人仓库的道具
			_ = gameConnect.GetAllPackItems()
			// 检查交易记录
			tradeHistory, _ := gameConnect.QueryTradeHistoryLog(0)
			log.Infof("购买记录有%d页", tradeHistory.GetTotalPageCount())
			handleTradeHistory(gameConnect, tradeHistory, excelFile)
			if tradeHistory.GetTotalPageCount() > 1 {
				for i := uint32(1); i < tradeHistory.GetTotalPageCount(); i++ {
					time.Sleep(2500 * time.Millisecond)
					history, _ := gameConnect.QueryTradeHistoryLog(i)
					handleTradeHistory(gameConnect, history, excelFile)
					newLogList := tradeHistory.GetLogList()
					newLogList = append(newLogList, history.GetLogList()...)
					tradeHistory.LogList = newLogList
				}
			}

			// 处理 买/卖 交易
			for _, pItem := range purchaseConfig.Items {
				itemName := pItem.ItemName
				itemId := items.GetItemIdByName(itemName)
				possessionCount, itemData := findPackItemCountById(gameConnect, itemId)
				// 买
				if pItem.IsBuyAction() {
					if gameConnect.Role.GetSilver() != 0 && gameConnect.Role.GetSilver() < purchaseConfig.MinZenyKeep {
						log.Warnf("角色身上zeny %d 低于 设定最低可交易zeny %d 跳过购买", gameConnect.Role.GetSilver(), purchaseConfig.MinZenyKeep)
						continue
					}
					if possessionCount > pItem.MaxPossession {
						log.Infof("身上有%d个%s 大于最大拥有值%d 跳过购买", possessionCount, itemName, pItem.MaxPossession)
						continue
					}
					err := buyItem(pItem, tradeHistory, itemId, gameConnect)
					if err != nil {
						if err == ErrNoItemFound {
							log.Infof("没有在交易所找到 %s 跳过购买...", pItem.ItemName)
							continue
						}
					}
				} else if pItem.IsSellAction() {
					// 卖
					pendingSells := gameConnect.QueryPendingSells()
					time.Sleep(2 * time.Second)
					if len(pendingSells.GetLists()) >= MaxSellItems {
						log.Warnf("已达到最大可同时上架数量 %d", MaxSellItems)
					}
					if possessionCount <= pItem.MaxPossession {
						log.Infof("身上有%d个%s 小于最大拥有值%d 跳过出售", possessionCount, itemName, pItem.MaxPossession)
						continue
					}
					if itemData[0].GetBase().GetIsfavorite() {
						log.Warnf("%s是喜爱物品不能出售", itemName)
					}
					sellCount := pItem.PurchaseCount
					if sellCount > possessionCount {
						sellCount = possessionCount
						pItem.PurchaseCount = pItem.PurchaseCount - possessionCount
					}
					if pItem.MaxPossession != DefaultMaxPossession && possessionCount-sellCount < pItem.MaxPossession {
						sellCount = possessionCount - pItem.MaxPossession
					}
					err := sellItem(itemData, pItem, sellCount, possessionCount, gameConnect)
					if err != nil {
						log.Errorf("上架%d个%s失败: %s", pItem.PurchaseCount, itemName, err)
					}
				}
			}

			// 游戏角色进入游戏地图
			// 每次购买交易会重新登录游戏
			if purchaseConfig.EnterMap {
				gameConnect.ChangeMap(gameConnect.Role.GetMapId())
				shouldReconnect = true
			}

			waitTime := time.Duration(purchaseConfig.BuyInterval) * time.Second
			log.Infof("等待 %f 秒", waitTime.Seconds())
			time.Sleep(waitTime)
		} else {
			log.Infof("等待角色进入游戏")
			time.Sleep(15 * time.Second)
		}
	}
}

// buy item from exchange
func buyItem(pItem *PurchaseItem, tradeHistory *Cmd.MyTradeLogRecordTradeCmd, itemId uint32, gameConnect *gameConnection.GameConnection) (err error) {
	priceList := gameConnect.QueryItemPrice(itemId, 0)
	if len(priceList) == 0 {
		return ErrNoItemFound
	}
	log.Infof("购买 %s 物品ID: %d", pItem.ItemName, itemId)
	// 跳过不买摆摊 除非比交易所便宜
	for _, item := range priceList {
		if item.GetType() == Cmd.ETradeType_ETRADETYPE_BOOTH && item.GetUpRate() == 0 {
			tradeItem(gameConnect, tradeHistory, pItem, item)
		} else if item.GetType() == Cmd.ETradeType_ETRADETYPE_TRADE {
			tradeItem(gameConnect, tradeHistory, pItem, item)
		}
	}
	return err
}

func sellItem(itemData []*Cmd.ItemData, pItem *PurchaseItem, sellCount, possessionCount uint32, gameConnect *gameConnection.GameConnection) (err error) {
	if len(itemData) < 1 {
		return nil
	}
	itemId := itemData[0].GetBase().GetId()
	itemName := items.GetItemName(itemId)
	price := gameConnect.ReqServerPrice(itemData[0])
	time.Sleep(1 * time.Second)

	if price.GetCount() > pItem.MaxExchangeCount && price.GetPrice() >= pItem.MinSellPrice {
		log.Infof("交易所有%d个%s 超过最大数量卖出%d", price.GetCount(), itemName, pItem.MaxExchangeCount)
	} else if price.GetCount() > pItem.MaxExchangeCount && price.GetPrice() < pItem.MinSellPrice {
		log.Infof("交易所有%d个%s 但价格%d低于设定最低价格%d 跳过出售",
			price.GetCount(),
			itemName,
			price.GetPrice(),
			pItem.MinSellPrice,
		)
		return err
	} else if price.GetPrice() < pItem.MaxPurchasePrice {
		log.Warnf("%s价格%d 低于最低上架价%d 交易所数量%d个",
			items.GetItemName(price.GetItemData().GetBase().GetId()),
			price.GetPrice(),
			pItem.MaxPurchasePrice,
			price.GetCount(),
		)
		return err
	}
	if price.GetPrice() == 0 {
		return ErrNoItemFound
	}
	if sellCount > 0 && possessionCount > 0 {
		if sellCount > possessionCount {
			sellCount = possessionCount
		}
		for _, item := range itemData {
			log.Infof("上架出售 %d个%s 价格 %d zeny id:%d", sellCount, items.GetItemName(itemId), price.GetPrice(), itemId)
			time.Sleep(2 * time.Second)
			result := gameConnect.SellItem(sellCount, price, item)
			log.Infof("上架 %s 结果: %v", itemName, result)
		}
	}
	return err
}

func hasPendingPurchase(tradeHis *Cmd.MyTradeLogRecordTradeCmd, itemId uint32, tradePrice uint64) uint32 {
	for _, tradeLog := range tradeHis.GetLogList() {
		if tradeLog.GetItemid() == itemId && tradeLog.GetPrice() == tradePrice && int64(tradeLog.GetEndtime()) > time.Now().Unix() {
			return tradeLog.GetTotalcount()
		}
	}
	return 0
}

func tradeItem(gameConnect *gameConnection.GameConnection, tradeHistory *Cmd.MyTradeLogRecordTradeCmd, pItem *PurchaseItem, itemInfo *Cmd.TradeItemBaseInfo) {
	itemName := pItem.ItemName
	purchaseCount := pItem.PurchaseCount
	itemCurPrice := itemInfo.GetPrice()
	itemCounts := itemInfo.GetCount()
	leaveCount := pItem.GetLeaveMinCount()
	if leaveCount > 0 && leaveCount <= itemCounts {
		log.Infof("交易所 %s 最低保有量 %d", itemName, leaveCount)
		purchaseCount -= leaveCount
		itemCounts -= leaveCount
	} else if leaveCount > itemCounts {
		log.Infof("交易所%s最低保有量%d大于出售量%d 跳过购买", itemName, leaveCount, itemCounts)
		return
	}
	time.Sleep(2 * time.Second)
	if itemInfo.GetUpRate() != 0 {
		itemCurPrice = uint64(math.Round(float64(itemCurPrice) * float64(1+itemInfo.GetUpRate()) / 1000 * pointDiscount))
	}
	if itemInfo.GetDownRate() != 0 {
		itemCurPrice = uint64(math.Round(float64(itemCurPrice) * float64(itemInfo.GetDownRate()) / 1000 * pointDiscount))
	}

	// 计算可以买入多少
	buyNum := math.Min(float64(purchaseCount), float64(itemCounts))
	if itemInfo.GetPublicityId() > 0 {
		// Check whether we have pending purchase
		pendingCount := hasPendingPurchase(tradeHistory, itemInfo.GetItemid(), itemCurPrice)
		log.Infof("已抢购 %d个 %s 中", pendingCount, itemName)
		buyNum = math.Min(float64(purchaseCount), float64(itemCounts-pendingCount))
	}

	if itemCurPrice < pItem.MaxPurchasePrice && buyNum > 0 {
		log.Infof("购买 %d个%s 交易所有%d个", uint32(buyNum), itemName, itemInfo.GetCount())
		result, _ := gameConnect.BuyItem(uint32(buyNum), itemInfo)
		log.Infof("购买结果: %v", result)
		if result.Ret != nil && result.GetRet() == Cmd.ETRADE_RET_CODE_ETRADE_RET_CODE_SUCCESS {
			silverInBag := gameConnect.Role.GetSilver() - uint64(buyNum)*itemCurPrice
			gameConnect.Role.Silver = &silverInBag
			log.Infof("角色剩余 %d zeny", silverInBag)
		}
	} else if buyNum == 0 {
		log.Infof("已经申请购入所有交易所 %s", itemName)
	} else {
		log.Infof("%s 价格 %d 比设定最高购买价 %d 高 跳过",
			itemName, itemCurPrice, pItem.MaxPurchasePrice,
		)
	}
}

func handleTradeHistory(connection *gameConnection.GameConnection, tradeHistory *Cmd.MyTradeLogRecordTradeCmd, excel *TradeExcel) {
	hasNewRecord := false
	log.Infof("检查购买记录第%d页", tradeHistory.GetIndex())
	for _, tradeLog := range tradeHistory.GetLogList() {
		itemName := items.GetItemName(tradeLog.GetItemid())
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			log.Infof("回收抢购失败的金币")
			takeFailedMoney(connection, tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuySuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NoramlBuy) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			hasNewRecord = true
			excel.AddRecord(tradeLog, items.GetItemName(tradeLog.GetItemid()))
			takeTradeLog(connection, tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			hasNewRecord = true
			excel.AddRecord(tradeLog, itemName)
			takeTradeLog(connection, tradeLog)
		}
		// if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess ||
		//	tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell && itemName == "哈比的羽毛") {
		//	hasNewRecord = true
		//	excel.AddRecord(tradeLog, itemName)
		// }
	}
	if hasNewRecord {
		excel.WriteExcel()
	}
}

func takeTradeLog(connection *gameConnection.GameConnection, tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess {
			log.Infof("卖出 %d个%s 赚取 %d zeny",
				tradeLog.GetCount(),
				items.GetItemName(tradeLog.GetItemid()),
				tradeLog.GetGetmoney(),
			)
		} else {
			log.Infof("取回从%s购买的物品 %d个%s 花费 %d zeny",
				tradeLog.GetNameInfo().GetName(),
				tradeLog.GetCount(),
				items.GetItemName(tradeLog.GetItemid()),
				tradeLog.GetCostmoney(),
			)
		}
		connection.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		time.Sleep(500 * time.Millisecond)
	}
}

func takeFailedMoney(gameConnect *gameConnection.GameConnection, tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
		tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		log.Infof("取回抢购失败 %d个%s %d zeny", tradeLog.GetFailcount(), items.GetItemName(tradeLog.GetItemid()), tradeLog.GetRetmoney())
		gameConnect.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		newSilver := gameConnect.Role.GetSilver() + tradeLog.GetRetmoney()
		gameConnect.Role.Silver = &newSilver
		time.Sleep(500 * time.Millisecond)
	}
}

func findPackItemCountById(gameConnect *gameConnection.GameConnection, itemId uint32) (itemCount uint32, itemData []*Cmd.ItemData) {
	packItems := gameConnect.Role.GetPackItems()
	for _, packItem := range packItems {
		for _, item := range packItem {
			if itemId == item.GetBase().GetId() {
				itemCount += item.GetBase().GetCount()
				itemData = append(itemData, item)
			}
		}
	}
	return itemCount, itemData
}
