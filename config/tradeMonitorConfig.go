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
	ItemName         string `yaml:"itemName" json:"itemName" label:"物品名称"`
	MaxPurchasePrice uint64 `yaml:"maxPurchasePrice" json:"maxPurchasePrice" label:"最大购买价格"`
	PurchaseCount    uint32 `yaml:"purchaseCount" json:"purchaseCount" label:"购买数量"`
	MaxPossession    uint32 `yaml:"maxPossession" json:"maxPossession" label:"最大拥有数量"`
	TradeAction      string `yaml:"TradeAction" json:"TradeAction" label:"操作(买/卖)"`
	MaxExchangeCount uint32 `yaml:"maxExchangeCount" json:"maxExchangeCount" label:"交易所保留最大数量"`
	MinSellPrice     uint64 `yaml:"minSellPrice" json:"minSellPrice" label:"最小出售价格"`
	LeaveMinCount    uint32 `yaml:"leaveMinCount" json:"leaveMinCount" label:"交易所保留最小数量"`
	RefineLv         string `yaml:"refineLv" json:"refineLv" label:"精炼等级比较(=0, >=5, <10)"`
	DamageEquip      bool   `yaml:"damageEquip" json:"damageEquip" label:"是否购买损坏装备"`
	MinZenyKeep      uint64 `yaml:"minZenyKeep" json:"minZenyKeep" label:"购买后保留最小zeny"`
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
	return p.TradeAction == "买" || strings.ToLower(p.TradeAction) == BuyAction
}

func (p *PurchaseItem) IsSellAction() bool {
	return p.TradeAction == "卖" || strings.ToLower(p.TradeAction) == SellAction
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
	MonitorInterval       int            `yaml:"monitorInterval" json:"monitorInterval" label:"监控间隔"` // in seconds
	WatchItems            []string       `yaml:"watchItems" json:"watchItems" label:"监控物品"`
	WatchCategories       []string       `yaml:"watchCategories" json:"watchCategories" label:"监控类别"`
	ElasticsearchHostPort string         `yaml:"elasticsearchHostPort" json:"elasticsearchHostPort" label:"elasticsearchHostPort"`
	NumberWorkers         int            `yaml:"numberWorkers" json:"numberWorkers" label:"工作线程数"`
	EnablePurchase        bool           `yaml:"enablePurchase" json:"enablePurchase" label:"启用买卖物品"`
	BuyItems              []PurchaseItem `yaml:"buyItem" json:"buyItem" label:"买卖物品"`
}

func (c *TradeMonitorConfig) GetPurchaseItems() []PurchaseItem {
	return c.BuyItems
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
		BuyItems:        make([]PurchaseItem, 1),
	}
}

func (c *TradeMonitorConfig) GetNumberWorkers() int {
	if c.NumberWorkers <= 0 {
		return 3
	}
	return c.NumberWorkers
}
