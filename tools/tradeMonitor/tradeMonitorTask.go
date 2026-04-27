package tradeMonitor

import (
	"context"
	"errors"
	"fmt"
	"io"
	"maps"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/esClient"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var (
	ErrNoItemFound       = errors.New("no item found")
	DefaultMaxPossession = uint32(99999999)
	MaxSellItems         = 8
	pointDiscount        = 0.75
)

type Task struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter io.Writer
	logger    *log.Logger
	items     *utils.ItemsLoader
}

func (t *Task) GetContext() context.Context {
	return t.ctx
}

func (t *Task) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(t.GC.LogWriter(), writer)
	t.logWriter = mw
	t.logger.SetOutput(mw)
}

func (t *Task) Start() {
	t.GC.Configs.AutoCreateChar = true
	t.GC.GameServerLogin()

	go func() {
		puConfig := t.GC.Configs.TradeMonitorConfig.GetPurchaseItems()
		for {
			select {
			case <-t.ctx.Done():
				t.logger.Info("停止交易监控任务")
				t.cancel()
				t.GC.Close()
				return
			default:
				if t.GC.Configs.TradeMonitorConfig.GetEnablePurchase() {
					// 处理自动交易
					t.autoTrade(puConfig)
				}

				t.monitorTrades()
				t.logger.Infof("等待 %d 秒后进行下一轮监控", t.GC.Configs.TradeMonitorConfig.GetMonitorInterval())
				time.Sleep(time.Duration(t.GC.Configs.TradeMonitorConfig.GetMonitorInterval()) * time.Second)
			}
		}
	}()
}

func (t *Task) Stop() {
	t.cancel()
}

func (t *Task) monitorTrades() {
	ch := make(chan uint32)
	wg := sync.WaitGroup{}
	tradeResults := []tradeItem{}
	t.logger.Infof("starting %d workers", t.GC.Configs.TradeMonitorConfig.GetNumberWorkers())

	for i := 0; i < t.GC.Configs.TradeMonitorConfig.GetNumberWorkers(); i++ {
		wg.Add(1)
		go queryItems(ch, &tradeResults, &wg, t.GC)
	}

	for _, catName := range t.GC.Configs.TradeMonitorConfig.GetWatchCategories() {
		// 这里添加具体的类别监控逻辑
		t.logger.Infof("监控类别: %s 的交易情况", catName)
		// 示例日志输出
		if _, ok := gameTypes.TradeZhCategories[catName]; !ok {
			t.logger.Warnf("未知的交易类别: %s", catName)
			t.logger.Warnf("支持的类别有: %v", utils.GetMapKeys(gameTypes.TradeZhCategories))
			continue
		}
		// 这里添加具体的类别监控逻辑
		catResults := t.GC.QueryCat(gameTypes.TradeZhCategories[catName])
		t.logger.Infof("%v", catResults)

		for _, itemId := range catResults.GetPubLists() {
			ch <- itemId
		}
		for _, itemId := range catResults.GetLists() {
			ch <- itemId
		}
	}

	for _, itemName := range t.GC.Configs.TradeMonitorConfig.GetWatchItems() {
		// 这里添加具体的交易监控逻辑
		t.logger.Infof("监控物品: %s 的交易情况", itemName)
		// 示例日志输出
	}

	close(ch)
	wg.Wait()

	if t.GC.Configs.TradeMonitorConfig.GetESHostPort() != "" {
		t.logger.Infof("上传交易数据到 Elasticsearch: %s", t.GC.Configs.TradeMonitorConfig.GetESHostPort())
		// 这里添加具体的上传逻辑
		t.uploadTradeRecords(tradeResults)
	}
}

func (t *Task) uploadTradeRecords(detail []tradeItem) {
	logger := t.GC.GetLogger()
	client, err := esClient.NewEsClient([]string{t.GC.Configs.TradeMonitorConfig.GetESHostPort()})
	if err != nil {
		logger.Errorf("failed to create elasticsearch client: %s", err)
		return
	}
	bulk := client.Bulk()
	for _, val := range detail {
		baseInfo := val.TradeBaseInfo
		sellInfo := val.TradeSellInfo
		salePrice := baseInfo.GetPrice()
		if baseInfo.GetDownRate() != 0 {
			salePrice = uint32(float64(salePrice) * float64(baseInfo.GetDownRate()) * 0.001)
		} else if baseInfo.GetUpRate() != 0 {
			salePrice = uint32(float64(salePrice) * (float64(baseInfo.GetUpRate())*0.001 + 1))
		}
		serverIdWithLine, _ := strconv.ParseUint(
			fmt.Sprintf("%d%d", t.GC.Configs.ZoneId, t.GC.Configs.ServerId),
			10,
			32,
		)
		template := esClient.ExchangeTemplate{
			ServerId:     uint32(serverIdWithLine),
			ItemId:       val.TradeBaseInfo.GetItemid(),
			ItemName:     t.GC.FindItemNameById(baseInfo.GetItemid()),
			ItemPrice:    uint64(salePrice),
			ItemCategory: t.GC.GetItemCat(baseInfo.GetItemid()),
			ItemRefineLv: baseInfo.GetRefineLv(),
			Count:        baseInfo.GetCount(),
			ItemEnhance:  baseInfo.GetItemData().GetEnchant(),
			TimeStamp:    time.Now(),
			TradeType:    baseInfo.GetType(),
			Guid:         baseInfo.GetGuid(),
		}
		if baseInfo.GetItemData().GetEquip() != nil {
			template.IsDamage = baseInfo.GetItemData().GetEquip().GetDamage()
		}
		if baseInfo.GetItemData() != nil && baseInfo.GetItemData().GetEnchant() != nil {
			template.ItemEnhance = baseInfo.GetItemData().GetEnchant()
		}
		if baseInfo.GetPublicityId() != 0 {
			template.IsPub = true
			template.ExpireTime = baseInfo.GetEndTime()
			if sellInfo != nil {
				template.BuyerCount = sellInfo.GetBuyerCount()
			}
		}
		req := elastic.NewBulkIndexRequest().Index(template.GetIndexName()).Doc(template)
		bulk.Add(req)
	}
	rsp, err := bulk.Do(context.Background())
	if err != nil {
		if err.Error() == "No bulk actions to commit" {
			retrySec := 60
			logger.Warnf("No information retrieve from exchange retrying in %d seconds", retrySec)
			time.Sleep(time.Duration(retrySec) * time.Second)
			return
		}
		logger.Errorf("failed to send bulk insert: %s", err)
	} else {
		logger.Infof("trying to upload %d trade records", len(detail))
		logger.Infof("response from elasticsearch: %d failed", len(rsp.Failed()))
	}
}

func (t *Task) autoTrade(purchaseConfig []config.PurchaseItem) {
	mails := t.GC.GetMails()
	t.logger.Infof("邮件数量: %d", len(mails))
	// 获取邮件
	for _, mail := range mails {
		if len(mail.GetAttach().GetAttachs()) > 0 {
			attachs := mail.GetAttach().GetAttachs()
			t.logger.Infof("邮件有附件: %s", attachs)
			t.logger.Infof("收取邮件 标题：%s 发送人：%s 内容：%s", mail.GetTitle(), mail.GetSender(), mail.GetMsg())
			// t.GC.GetMailAttachment(mail.GetId())
		}
	}

	// 检查交易记录
	tradeHistory, _ := t.GC.QueryTradeHistoryLog(0)
	t.logger.Infof("购买记录有%d页", tradeHistory.GetTotalPageCount())
	t.handleTradeHistory(tradeHistory)
	if tradeHistory.GetTotalPageCount() > 1 {
		for i := uint32(1); i < tradeHistory.GetTotalPageCount(); i++ {
			time.Sleep(2500 * time.Millisecond)
			history, _ := t.GC.QueryTradeHistoryLog(i)
			t.handleTradeHistory(history)
			newLogList := tradeHistory.GetLogList()
			newLogList = append(newLogList, history.GetLogList()...)
			tradeHistory.LogList = newLogList
		}
	}

	// 处理 买/卖 交易
	for _, pc := range purchaseConfig {
		if pc.ItemName == "" {
			continue
		}
		itemName := pc.ItemName
		itemId := t.items.GetItemIdByName(itemName)
		possessionCount, itemData := t.findPackItemCountById(itemId)
		// 买
		if pc.IsBuyAction() {
			if t.GC.Role.GetSilver() != 0 && t.GC.Role.GetSilver() < pc.MinZenyKeep {
				t.logger.Warnf("角色身上zeny %d 低于 设定最低可交易zeny %d 跳过购买", t.GC.Role.GetSilver(), pc.MinZenyKeep)
				continue
			}
			if possessionCount > pc.MaxPossession {
				t.logger.Infof("身上有%d个%s 大于最大拥有值%d 跳过购买", possessionCount, itemName, pc.MaxPossession)
				continue
			}
			err := t.buyItem(pc, tradeHistory, itemId)
			if err != nil {
				if err == ErrNoItemFound {
					t.logger.Infof("没有在交易所找到 %s 跳过购买...", pc.ItemName)
				}
			}
		} else if pc.IsSellAction() {
			// 卖
			pendingSells := t.GC.QueryPendingSells()
			time.Sleep(2 * time.Second)
			if len(pendingSells.GetLists()) >= MaxSellItems {
				t.logger.Warnf("已达到最大可同时上架数量 %d", MaxSellItems)
				continue
			}
			if possessionCount <= pc.MaxPossession {
				t.logger.Infof("身上有%d个%s 小于最大拥有值%d 跳过出售", possessionCount, itemName, pc.MaxPossession)
				continue
			}
			// Not available in EP 5.0
			// if itemData[0].GetBase().GetIsfavorite() {
			// 	log.Warnf("%s是喜爱物品不能出售", itemName)
			// }
			sellCount := pc.PurchaseCount
			if sellCount > possessionCount {
				sellCount = possessionCount
				pc.PurchaseCount = pc.PurchaseCount - possessionCount
			}
			if pc.MaxPossession != DefaultMaxPossession && possessionCount-sellCount < pc.MaxPossession {
				sellCount = possessionCount - pc.MaxPossession
			}
			err := t.sellItem(itemData, pc, sellCount, possessionCount)
			if err != nil {
				t.logger.Errorf("上架%d个%s失败: %s", pc.PurchaseCount, itemName, err)
				continue
			}
		}
	}
}

func (t *Task) findPackItemCountById(itemId uint32) (itemCount uint32, itemData []*Cmd.ItemData) {
	packItems := t.GC.Role.GetPackItems()
	for _, packItem := range packItems {
		pi := maps.Clone(packItem)
		for _, item := range pi {
			if itemId == item.GetBase().GetId() {
				itemCount += item.GetBase().GetCount()
				itemData = append(itemData, item)
			}
		}
	}
	return itemCount, itemData
}

func (t *Task) handleTradeHistory(tradeHistory *Cmd.MyTradeLogRecordTradeCmd) {
	t.logger.Infof("检查购买记录第%d页", tradeHistory.GetIndex())
	for _, tradeLog := range tradeHistory.GetLogList() {
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			t.logger.Infof("回收抢购失败的金币")
			t.takeFailedMoney(tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuySuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NoramlBuy) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			t.takeTradeLog(tradeLog)
		}
		if (tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell) &&
			tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
			t.takeTradeLog(tradeLog)
		}
	}
}

func (t *Task) takeFailedMoney(tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicityBuyFail &&
		tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		t.logger.Infof("取回抢购失败 %d个%s %d zeny", tradeLog.GetFailcount(), t.items.GetItemName(tradeLog.GetItemid()), tradeLog.GetRetmoney())
		t.GC.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		newSilver := t.GC.Role.GetSilver() + uint64(tradeLog.GetRetmoney())
		t.GC.Role.Silver = &newSilver
		time.Sleep(500 * time.Millisecond)
	}
}

func (t *Task) takeTradeLog(tradeLog *Cmd.LogItemInfo) {
	if tradeLog.GetStatus() == Cmd.ETakeStatus_ETakeStatus_CanTakeGive {
		if tradeLog.GetLogtype() == Cmd.EOperType_EOperType_NormalSell ||
			tradeLog.GetLogtype() == Cmd.EOperType_EOperType_PublicitySellSuccess {
			t.logger.Infof("卖出 %d个%s 赚取 %d zeny",
				tradeLog.GetCount(),
				t.items.GetItemName(tradeLog.GetItemid()),
				tradeLog.GetGetmoney(),
			)
		} else {
			t.logger.Infof("取回从%s购买的物品 %d个%s 花费 %d zeny",
				tradeLog.GetNameInfo().GetName(),
				tradeLog.GetCount(),
				t.items.GetItemName(tradeLog.GetItemid()),
				tradeLog.GetCostmoney(),
			)
		}
		t.GC.TakeLogTrade(tradeLog.GetId(), tradeLog.GetLogtype())
		time.Sleep(500 * time.Millisecond)
	}
}

func (t *Task) buyItem(pItem config.PurchaseItem, tradeHistory *Cmd.MyTradeLogRecordTradeCmd, itemId uint32) (err error) {
	priceList := t.GC.QueryItemPrice(itemId, 0)
	if len(priceList) == 0 {
		return ErrNoItemFound
	}
	t.logger.Infof("购买 %s 物品ID: %d", pItem.ItemName, itemId)
	// 跳过不买摆摊 除非比交易所便宜
	for _, item := range priceList {
		if item.GetType() == Cmd.ETradeType_ETRADETYPE_BOOTH && item.GetUpRate() == 0 {
			t.tradeItem(tradeHistory, pItem, item)
		} else if item.GetType() == Cmd.ETradeType_ETRADETYPE_TRADE {
			t.tradeItem(tradeHistory, pItem, item)
		}
	}
	return err
}

func (t *Task) tradeItem(tradeHistory *Cmd.MyTradeLogRecordTradeCmd, pItem config.PurchaseItem, itemInfo *Cmd.TradeItemBaseInfo) {
	itemName := pItem.ItemName
	purchaseCount := pItem.PurchaseCount
	itemCurPrice := itemInfo.GetPrice()
	itemCounts := itemInfo.GetCount()
	leaveCount := pItem.GetLeaveMinCount()
	if leaveCount > 0 && leaveCount <= itemCounts {
		t.logger.Infof("交易所 %s 最低保有量 %d", itemName, leaveCount)
		purchaseCount -= leaveCount
		itemCounts -= leaveCount
	} else if leaveCount > itemCounts {
		t.logger.Infof("交易所%s最低保有量%d大于出售量%d 跳过购买", itemName, leaveCount, itemCounts)
		return
	}
	time.Sleep(time.Second)
	if itemInfo.GetUpRate() != 0 {
		itemCurPrice = uint32(math.Round(float64(itemCurPrice) * float64(1+itemInfo.GetUpRate()) / 1000 * pointDiscount))
	}
	if itemInfo.GetDownRate() != 0 {
		itemCurPrice = uint32(math.Round(float64(itemCurPrice) * float64(itemInfo.GetDownRate()) / 1000 * pointDiscount))
	}

	// 计算可以买入多少
	buyNum := math.Min(float64(purchaseCount), float64(itemCounts))
	if itemInfo.GetPublicityId() > 0 {
		// Check whether we have pending purchase
		pendingCount := t.hasPendingPurchase(tradeHistory, itemInfo.GetItemid(), uint64(itemCurPrice))
		t.logger.Infof("已抢购 %d个 %s 中", pendingCount, itemName)
		buyNum = math.Min(float64(purchaseCount), float64(itemCounts-pendingCount))
	}

	if mismatches, err := pItem.CompareRefineLv(itemInfo); itemInfo.GetItemData() != nil && (len(mismatches) > 0 || err != nil) {
		if err != nil {
			t.logger.Errorf("交易所 %s 精炼等级 设定购买等级 %s 比较失败: %s 跳过购买", itemName, pItem.RefineLv, err)
		}
		for _, _ = range mismatches {
			t.logger.Infof("交易所 %s 精炼等级 %d 设定购买等级 %s 跳过购买",
				itemName,
				itemInfo.GetRefineLv(),
				pItem.RefineLv,
			)
		}
		return
	}

	if itemInfo.GetItemData().GetEquip().GetRefinelv() > 0 && itemInfo.GetItemData().GetEquip().GetDamage() != pItem.DamageEquip {
		t.logger.Infof("交易所 %s 是破损 %t 设定购买破损 %t 跳过购买",
			itemName,
			itemInfo.GetItemData().GetEquip().GetDamage(),
			pItem.DamageEquip,
		)
		return
	}

	if uint64(itemCurPrice) < pItem.MaxPurchasePrice && buyNum > 0 {
		t.logger.Infof("购买 %d个%s 交易所有%d个 价格: %d", uint32(buyNum), itemName, itemInfo.GetCount(), itemCurPrice)
		result, _ := t.GC.BuyItem(uint32(buyNum), itemInfo)
		t.logger.Infof("购买结果: %v", result)
		if result.Ret != nil && result.GetRet() == Cmd.ETRADE_RET_CODE_ETRADE_RET_CODE_SUCCESS {
			t.logger.Infof("角色剩余 %d zeny", t.GC.Role.GetSilver())
		}
	} else if buyNum == 0 {
		t.logger.Infof("已经申请购入所有交易所 %s", itemName)
	} else {
		t.logger.Infof("%s 价格 %d 比设定最高购买价 %d 高 跳过",
			itemName, itemCurPrice, pItem.MaxPurchasePrice,
		)
	}
}

func (t *Task) sellItem(itemData []*Cmd.ItemData, pItem config.PurchaseItem, sellCount, possessionCount uint32) (err error) {
	if len(itemData) < 1 {
		return nil
	}
	itemId := itemData[0].GetBase().GetId()
	itemName := t.items.GetItemName(itemId)
	price := t.GC.ReqServerPrice(itemData[0])
	time.Sleep(1 * time.Second)

	if price.GetCount() > pItem.MaxExchangeCount && price.GetPrice() >= uint32(pItem.MinSellPrice) {
		t.logger.Infof("交易所有%d个%s 超过最大数量卖出%d", price.GetCount(), itemName, pItem.MaxExchangeCount)
	} else if price.GetCount() > pItem.MaxExchangeCount && price.GetPrice() < uint32(pItem.MinSellPrice) {
		t.logger.Infof("交易所有%d个%s 但价格%d低于设定最低价格%d 跳过出售",
			price.GetCount(),
			itemName,
			price.GetPrice(),
			pItem.MinSellPrice,
		)
		return err
	} else if price.GetPrice() < uint32(pItem.MaxPurchasePrice) {
		t.logger.Warnf("%s价格%d 低于最低上架价%d 交易所数量%d个",
			t.items.GetItemName(price.GetItemData().GetBase().GetId()),
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
			t.logger.Infof("上架出售 %d个%s 价格 %d zeny id:%d", sellCount, t.items.GetItemName(itemId), price.GetPrice(), itemId)
			time.Sleep(2 * time.Second)
			result := t.GC.SellItem(sellCount, price, item)
			t.logger.Infof("上架 %s 结果: %v", itemName, result)
		}
	}
	return err
}

func (t *Task) hasPendingPurchase(tradeHis *Cmd.MyTradeLogRecordTradeCmd, itemId uint32, tradePrice uint64) uint32 {
	for _, tradeLog := range tradeHis.GetLogList() {
		if tradeLog.GetItemid() == itemId && tradeLog.GetPrice() == uint32(tradePrice) && int64(tradeLog.GetEndtime()) > time.Now().Unix() {
			return tradeLog.GetTotalcount()
		}
	}
	return 0
}

type tradeItem struct {
	TradeBaseInfo *Cmd.TradeItemBaseInfo
	TradeSellInfo *Cmd.ItemSellInfoRecordTradeCmd
}

func queryItems(ch chan uint32, details *[]tradeItem, wg *sync.WaitGroup, connection *gameConnection.GameConnection) {
	logger := connection.GetLogger()
	defer wg.Done()
	initStart := true
	for itemId := range ch {
		if initStart {
			time.Sleep(time.Second + time.Duration(rand.Int31n(500))*time.Millisecond)
			initStart = false
		}
		logger.Infof("Requesting information for exchange item id: %d", itemId)
		detail := connection.QueryItemPrice(itemId, 0)
		if len(detail) > 0 {
			// traverse all sub items under same itemId
			for _, v := range detail {
				trade := tradeItem{
					TradeBaseInfo: v,
				}
				// logger.Infof("%d response found", itemId)
				if trade.TradeBaseInfo.GetRefineLv() > 0 {
					logger.Infof("ID: %d; 物品 %s 有精炼等级 %d", itemId, connection.FindItemNameById(itemId), trade.TradeBaseInfo.GetRefineLv())
				}
				logger.Infof("ID: %d, %s; 价格：%d 数量：%d",
					itemId,
					connection.FindItemNameById(itemId),
					v.GetPrice(),
					v.GetCount(),
				)
				if detail[0].GetPublicityId() != 0 {
					logger.Infof("查询公示卖家数量")
					sellInfo := connection.QueryItemSellInfo(itemId, detail[0].GetPublicityId())
					if sellInfo != nil {
						logger.Infof("商品 %s 有 %d 人抢购", connection.FindItemNameById(itemId), sellInfo.GetBuyerCount())
					}
					trade.TradeSellInfo = sellInfo
				}
				*details = append(*details, trade)
			}
		}
	}
}

func NewTradeMonitorTask(ctx context.Context, gc *gameConnection.GameConnection) *Task {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &Task{
		GC:        gc,
		ctx:       newCtx,
		cancel:    cancel,
		logWriter: mw,
		logger:    logger,
		items:     utils.NewItemsLoader("", "", ""),
	}
}
