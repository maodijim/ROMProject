package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

const (
	BuyAction               = "buy"
	SellAction              = "sell"
	DefaultMaxPurchasePrice = 2000000000
	DefaultPurchaseCount    = 1
	DefaultBuyInterval      = 300
	DefaultMaxPossession    = 99999999
	DefaultItemAction       = BuyAction
	DefaultMaxExchangeCount = 99999999
)

type PurchaseConfig struct {
	AuthPass    string          `yaml:"authPass"`
	BuyInterval uint32          `yaml:"buyInterval"`
	MinZenyKeep uint64          `yaml:"minZenyKeep"`
	EnterMap    bool            `yaml:"enterMap"`
	Items       []*PurchaseItem `yaml:"items"`
}

type PurchaseItem struct {
	ItemName         string `yaml:"itemName"`
	MaxPurchasePrice uint64 `yaml:"maxPurchasePrice"`
	PurchaseCount    uint32 `yaml:"purchaseCount"`
	MaxPossession    uint32 `yaml:"maxPossession"`
	Action           string `yaml:"action"`
	MaxExchangeCount uint32 `yaml:"maxExchangeCount"`
	MinSellPrice     uint64 `yaml:"minSellPrice"`
}

func (p *PurchaseItem) IsBuyAction() bool {
	return p.Action == "买" || strings.ToLower(p.Action) == BuyAction
}

func (p *PurchaseItem) IsSellAction() bool {
	return p.Action == "卖" || strings.ToLower(p.Action) == SellAction
}

func PurchaseConfigParser(purchaseYml string) *PurchaseConfig {
	ymlPath := purchaseYml
	var pConfig *PurchaseConfig
	if ymlPath == "" {
		ymlPath = "purchaseItems.yaml"
	}
	f, err := os.Open(ymlPath)
	if err != nil {
		log.Errorf("failed to open %s: %s", ymlPath, err)
		log.Exit(2)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&pConfig)
	if err != nil {
		log.Errorf("parse purchase item yaml failed: %s", err)
		log.Exit(3)
	}

	if pConfig.BuyInterval == 0 {
		pConfig.BuyInterval = DefaultBuyInterval
	}

	for _, item := range pConfig.Items {
		if item.Action == "" {
			item.Action = DefaultItemAction
		}
		if item.MaxPurchasePrice == 0 {
			if item.IsBuyAction() {
				item.MaxPurchasePrice = DefaultMaxPurchasePrice
			} else if item.IsSellAction() {
				item.MaxPurchasePrice = 0
			}

		}
		if item.PurchaseCount == 0 {
			item.PurchaseCount = DefaultPurchaseCount
		}
		if item.MaxExchangeCount == 0 {
			item.MaxExchangeCount = DefaultMaxExchangeCount
		}
		if item.MaxPossession == 0 {
			if item.IsBuyAction() {
				item.MaxPossession = DefaultMaxPossession
			} else if item.IsSellAction() {
				item.MaxPossession = 0
			}
		}
	}
	return pConfig
}
