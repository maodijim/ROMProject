package main

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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
	BuyInterval int32           `yaml:"buyInterval"`
	MinZenyKeep uint64          `yaml:"minZenyKeep"`
	EnterMap    bool            `yaml:"enterMap"`
	Items       []*PurchaseItem `yaml:"items"`
	// Should Log out after purchase complete
	LogOut bool `yaml:"logOut"`
}

type PurchaseItem struct {
	ItemName         string `yaml:"itemName"`
	MaxPurchasePrice uint64 `yaml:"maxPurchasePrice"`
	PurchaseCount    uint32 `yaml:"purchaseCount"`
	MaxPossession    uint32 `yaml:"maxPossession"`
	Action           string `yaml:"action"`
	MaxExchangeCount uint32 `yaml:"maxExchangeCount"`
	MinSellPrice     uint64 `yaml:"minSellPrice"`
	LeaveMinCount    uint32 `yaml:"leaveMinCount"`
	RefineLv         string `yaml:"refineLv"`
	DamageEquip      bool   `yaml:"damageEquip"`
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
	re, err := regexp.Compile(`([><!=]+)=?\s?(\d)`)
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
