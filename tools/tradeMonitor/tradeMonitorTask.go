package tradeMonitor

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	"ROMProject/Cmds"
	"ROMProject/esClient"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

type Task struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter io.Writer
	logger    *log.Logger
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
		for {
			select {
			case <-t.ctx.Done():
				t.logger.Info("停止交易监控任务")
				t.cancel()
				t.GC.Close()
				return
			default:
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
	log.Infof("starting %d workers", t.GC.Configs.TradeMonitorConfig.GetNumberWorkers())

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
		log.Infof("%v", catResults)

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
	}
}
