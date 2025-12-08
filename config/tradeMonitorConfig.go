package config

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

const (
	BuyAction  = "buy"
	SellAction = "sell"
)

type PurchaseItem struct {
	ItemName         string `yaml:"itemName" json:"物品名称"`
	MaxPurchasePrice uint64 `yaml:"maxPurchasePrice" json:"最大购买价格"`
	PurchaseCount    uint32 `yaml:"purchaseCount" json:"购买数量"`
	MaxPossession    uint32 `yaml:"maxPossession" json:"最大拥有数量"`
	Action           string `yaml:"action" json:"操作(买/卖)"`
	MaxExchangeCount uint32 `yaml:"maxExchangeCount" json:"交易所保留最大数量"`
	MinSellPrice     uint64 `yaml:"minSellPrice" json:"最小出售价格"`
	LeaveMinCount    uint32 `yaml:"leaveMinCount" json:"交易所保留最小数量"`
	RefineLv         string `yaml:"refineLv" json:"精炼等级比较=0 >=5 <10"`
	DamageEquip      bool   `yaml:"damageEquip" json:"是否购买损坏装备"`
	MinZenyKeep      uint64 `yaml:"minZenyKeep" json:"购买后保留最小zeny"`
}

func (p *PurchaseItem) ParseConfigFromInterface(config map[string]any) PurchaseItem {
	utils.ParseConfigFromInterface(config, p)
	if p.RefineLv == "" || p.RefineLv == "0" {
		p.RefineLv = "=0"
	}
	return *p
}

func (p *PurchaseItem) GetLeaveMinCount() uint32 {
	return p.LeaveMinCount
}

func (p *PurchaseItem) IsBuyAction() bool {
	return p.Action == "买" || strings.ToLower(p.Action) == BuyAction
}

func (p *PurchaseItem) IsSellAction() bool {
	return p.Action == "卖" || strings.ToLower(p.Action) == SellAction
}

func (p *PurchaseItem) CompareRefineLv(info *Cmd.TradeItemBaseInfo) (mismatches []string, err error) {
	if p.RefineLv == "" {
		return mismatches, nil
	}
	compare, lv, err := p.ParseRefineLv()
	if err != nil {
		log.Errorf("compile refine lv regex failed: %s", err)
		return mismatches, err
	}
	if len(compare) != len(lv) {
		log.Errorf("invalid refine lv compare string")
		return mismatches, errors.New("lv and compare string not the same length")
	}
	for i, comparison := range compare {
		switch comparison {
		case ">":
			if info.GetRefineLv() <= lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case ">=":
			if info.GetRefineLv() < lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case "<":
			if info.GetRefineLv() >= lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case "<=":
			if info.GetRefineLv() > lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case "!=":
			if info.GetRefineLv() == lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case "!":
			if info.GetRefineLv() == lv[i] {
				mismatches = append(mismatches, comparison)
			}
		case "=":
			if info.GetRefineLv() != lv[i] {
				mismatches = append(mismatches, comparison)
			}
		default:
			log.Errorf("unknown compare string")
			return mismatches, errors.New("invalid refine lv compare string")
		}
	}
	return mismatches, nil
}

func (p *PurchaseItem) ParseRefineLv() (compare []string, lv []uint32, err error) {
	re, err := regexp.Compile(`([><!=]+)=?\s?(\d+)`)
	if err != nil {
		return compare, lv, err
	}
	matches := re.FindAllStringSubmatch(p.RefineLv, -1)
	if len(matches) == 0 {
		return compare, lv, errors.New("invalid refine lv compare string")
	}
	for _, match := range matches {
		if len(match) != 3 {
			return compare, lv, errors.New("invalid refine lv compare string")
		}
		comparison := match[1]
		matchLv := match[2]
		targetLv, _ := strconv.ParseUint(matchLv, 10, 32)
		compare = append(compare, comparison)
		lv = append(lv, uint32(targetLv))
	}
	return compare, lv, nil
}

type TradeMonitorConfig struct {
	MonitorInterval       int          `yaml:"monitorInterval" json:"监控间隔"` // in seconds
	WatchItems            []string     `yaml:"watchItems" json:"监控物品"`
	WatchCategories       []string     `yaml:"watchCategories" json:"监控类别"`
	ElasticsearchHostPort string       `yaml:"elasticsearchHostPort" json:"elasticsearchHostPort"`
	NumberWorkers         int          `yaml:"numberWorkers" json:"工作线程数"`
	EnablePurchase        bool         `yaml:"enablePurchase" json:"启用买卖物品"`
	BuyItems1             PurchaseItem `yaml:"buyItem1" json:"购买物品1"`
	BuyItems2             PurchaseItem `yaml:"buyItem2" json:"购买物品2"`
	BuyItems3             PurchaseItem `yaml:"buyItem3" json:"购买物品3"`
	BuyItems4             PurchaseItem `yaml:"buyItem4" json:"购买物品4"`
	BuyItems5             PurchaseItem `yaml:"buyItem5" json:"购买物品5"`
}

func (c *TradeMonitorConfig) GetPurchaseItems() []PurchaseItem {
	return []PurchaseItem{c.BuyItems1, c.BuyItems2, c.BuyItems3, c.BuyItems4, c.BuyItems5}
}

func (c *TradeMonitorConfig) GetEnablePurchase() bool {
	return c.EnablePurchase
}

func (c *TradeMonitorConfig) GetESHostPort() string {
	return c.ElasticsearchHostPort
}

func (c *TradeMonitorConfig) GetMonitorInterval() int {
	if c.MonitorInterval <= 0 {
		return 30
	}
	return c.MonitorInterval
}

func (c *TradeMonitorConfig) GetWatchItems() []string {
	return c.WatchItems
}

func (c *TradeMonitorConfig) GetWatchCategories() []string {
	return c.WatchCategories
}

func (c *TradeMonitorConfig) ParseFromInterface(config map[string]interface{}) TradeMonitorConfig {
	utils.ParseConfigFromInterface(config, c)
	return *c
}

func (c *TradeMonitorConfig) GetDefault() TradeMonitorConfig {
	return TradeMonitorConfig{
		MonitorInterval: 30,
		WatchItems:      []string{},
		WatchCategories: utils.GetMapKeys(gameTypes.TradeZhCategories),
		NumberWorkers:   3,
	}
}

func (c *TradeMonitorConfig) GetNumberWorkers() int {
	if c.NumberWorkers <= 0 {
		return 3
	}
	return c.NumberWorkers
}
