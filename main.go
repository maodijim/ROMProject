package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/esClient"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

const (
	dev_version = "v0.0.5"
)

type tradeItem struct {
	TradeBaseInfo *Cmd.TradeItemBaseInfo
	TradeSellInfo *Cmd.ItemSellInfoRecordTradeCmd
}

func init() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func queryItems(ch chan uint32, details *[]tradeItem, wg *sync.WaitGroup, connection *gameConnection.GameConnection, items *utils.ItemsLoader) {
	defer wg.Done()
	initStart := true
	for itemId := range ch {
		if initStart {
			time.Sleep(time.Second + time.Duration(rand.Int31n(500))*time.Millisecond)
			initStart = false
		}
		log.Infof("Requesting information for exchange item id: %d", itemId)
		detail := connection.QueryItemPrice(itemId, 0)
		if len(detail) > 0 {
			// traverse all sub items under same itemId
			for _, v := range detail {
				trade := tradeItem{
					TradeBaseInfo: v,
				}
				log.Infof("%d response found", itemId)
				log.Infof("%s；价格：%d 数量：%d",
					items.GetItemName(itemId),
					v.GetPrice(),
					v.GetCount(),
				)
				if detail[0].GetPublicityId() != 0 {
					log.Infof("查询公示卖家数量")
					sellInfo := connection.QueryItemSellInfo(itemId, detail[0].GetPublicityId())
					if sellInfo != nil {
						log.Infof("商品 %s 有 %d 人抢购", items.GetItemName(itemId), sellInfo.GetBuyerCount())
					}
					trade.TradeSellInfo = sellInfo
				}
				*details = append(*details, trade)
			}
		}
	}
}

func main() {
	// hostname := "gw-m-ro.xd.com"
	// ip, err := net.LookupIP(hostname)
	// if err != nil {
	//	log.Errorf("%s", err)
	// }
	// log.Infof("ip for %s is :%s", hostname, ip)

	confFile := flag.String("configPath", "", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	skillJson := flag.String("skillJson", "skills.yml", "json file of skills")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	workerNum := flag.Int("worker", 3, "Number of work to pull Trade Info")
	sleepInterval := flag.Int("sleepFor", 1000, "Interval between each pull circle")
	pullOnce := flag.Bool("once", false, "Pull Trade Information once")
	flag.Parse()

	log.Infof("Version: %s", dev_version)

	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillJson)

	wg := sync.WaitGroup{}
	ch := make(chan uint32)
	detail := make([]tradeItem, 0)

	cont := true
	if cont {
		gameConnect := gameConnection.NewConnection(conf, skills, items)
		if *enableDebug {
			gameConnect.DebugMsg = true
			log.SetLevel(log.DebugLevel)
		}
		gameConnect.GameServerLogin()
		cats := []uint32{
			12,   // 图纸
			1001, // 药剂/效果
			1002, // 精炼
			1003, // 卷轴/唱片
			1004, // 材料
			1005, // 节日材料
			1007, // 宠物材料
			1008, // 宠物头饰图纸
			1009, // 宠物头饰
			1010, // cards - 武器
			1011, // cards - 副手
			1012, // cards - 盔甲
			1013, // cards - 披风
			1014, // cards - 鞋子
			1015, // cards - 饰品
			1016, // cards - 头部

			1025, // weapon
			1026, // 副手
			1027, // 盔甲
			1028, // 披风
			1029, // 鞋子
			1030, // 饰品

			1045, // 时装
			1052, // 限定特典
		}

		for {
			if gameConnect.Role.GetInGame() {
				ch = make(chan uint32)
				log.Infof("starting %d workers", *workerNum)
				for t := 0; t < *workerNum; t++ {
					wg.Add(1)
					go queryItems(ch, &detail, &wg, gameConnect, items)
				}
				for _, cat := range cats {
					log.Infof("Requesting information for catergory id: %d", cat)
					results := gameConnect.QueryCat(cat)
					log.Infof("%v", results)

					for _, itemId := range results.GetPubLists() {
						ch <- itemId
					}
					for _, itemId := range results.GetLists() {
						ch <- itemId
					}
				}
				close(ch)
				wg.Wait()

				log.Infof("Total %d results", len(detail))
				log.Info("End Query")
				// Insert to elasticsearch
				ctx := context.Background()
				client := esClient.NewEsClient(conf.EsConfig.Urls)
				bulk := client.Bulk()
				now := time.Now()
				for _, val := range detail {
					baseInfo := val.TradeBaseInfo
					sellInfo := val.TradeSellInfo
					salePrice := baseInfo.GetPrice()
					if baseInfo.GetDownRate() != 0 {
						salePrice = uint64(float64(salePrice) * float64(baseInfo.GetDownRate()) * 0.001)
					} else if baseInfo.GetUpRate() != 0 {
						salePrice = uint64(float64(salePrice) * (float64(baseInfo.GetUpRate())*0.001 + 1))
					}
					serverIdWithLine, _ := strconv.ParseUint(
						fmt.Sprintf("%d%d", gameConnect.Configs.ZoneId, gameConnect.Configs.ServerId),
						10,
						32,
					)
					template := esClient.ExchangeTemplate{
						ServerId:     uint32(serverIdWithLine),
						ItemId:       val.TradeBaseInfo.GetItemid(),
						ItemName:     items.GetItemName(baseInfo.GetItemid()),
						ItemPrice:    salePrice,
						ItemCategory: items.GetItemCat(baseInfo.GetItemid()),
						ItemRefineLv: baseInfo.GetRefineLv(),
						Count:        baseInfo.GetCount(),
						ItemEnhance:  baseInfo.GetItemData().GetEnchant(),
						TimeStamp:    now,
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
				rsp, err := bulk.Do(ctx)
				if err != nil {
					if err.Error() == "No bulk actions to commit" {
						retrySec := 60
						log.Warnf("No information retrieve from exchange retrying in %d seconds", retrySec)
						time.Sleep(time.Duration(retrySec) * time.Second)
						continue
					}
					log.Errorf("failed to send bulk insert: %s", err)
				} else {
					log.Infof("response from elasticsearch: %d failed", len(rsp.Failed()))
				}
				detail = []tradeItem{}

				if *pullOnce {
					os.Exit(0)
				}

			}

			if gameConnect.Role.GetInGame() {
				sleepFor := utils.RandomSleepTime(*sleepInterval, 500)
				log.Infof("sleeping for %d seconds", sleepFor)
				time.Sleep(time.Duration(sleepFor) * time.Second)
			} else {
				log.Info("Sleeping for 15 seconds to wait for game client connect")
				time.Sleep(15 * time.Second)
			}
		}
	}

}
