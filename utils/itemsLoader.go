package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"ROMProject/data"
	log "github.com/sirupsen/logrus"
)

type ExchangeItem struct {
	Trade         string      `json:"Trade,omitempty"`
	NameZh        string      `json:"NameZh"`
	Id            json.Number `json:"id"`
	Category      string      `json:"Category"`
	ShowTime      string      `json:"ShowTime,omitempty"`
	Overlap       string      `json:"Overlap,omitempty"`
	UnTradeTime   string      `json:"UnTradeTime,omitempty"`
	TFTradeTime   string      `json:"TFTradeTime,omitempty"`
	TFUnTradeTime string      `json:"TRUnTradeTime,,omitempty"`
}

type BuffItem struct {
	Id          json.Number `json:"id"`
	BuffDesc    string      `json:"BuffDesc"`
	BuffIcon    string      `json:"BuffIcon,omitempty"`
	BuffName    string      `json:"BuffName"`
	IconType    string      `json:"IconType,omitempty"`
	NoAttack    string      `json:"NoAttack,omitempty"`
	Odds        string      `json:"Odds,,omitempty"`
	TransformID string      `json:"TransformID,,omitempty"`
	IsDisperse  string      `json:"isdisperse,omitempty"`
	IsGain      string      `json:"isgain,,omitempty"`
	Rate        json.Number `json:"rate,omitempty"`
	Shape       string      `json:"shape,omitempty"`
	Type        string      `json:"type,,omitempty"`
}

type BuffItemByName struct {
	Items []BuffItem
}

// Items All Available items
type Items struct {
	AdventureSort  string      `json:"AdventureSort"`
	AuctionPrice   json.Number `json:"AuctionPrice"`
	AdventureValue string      `json:"AdventureValue"`
	Condition      string      `json:"Condition"`
	Feature        string      `json:"Feature"`
	Desc           string      `json:"Desc"`
	Icon           string      `json:"Icon"`
	Level          json.Number `json:"Level"`
	LoadShowSize   string      `json:"LoadShowSize"`
	MaxNum         uint32      `json:"MaxNum"`
	NameZh         string      `json:"NameZh"`
	NoStorage      string      `json:"NoStorage"`
	Quality        json.Number `json:"Quality"`
	SellPrice      string      `json:"SellPrice"`
	Type           json.Number `json:"Type"`
	Id             json.Number `json:"id"`
}

func (i Items) GetLevel() uint64 {
	if i.Level != "" {
		result, _ := strconv.ParseUint(i.Level.String(), 10, 32)
		return result
	}
	return 0
}

type ItemsByName struct {
	Items []Items
}

type ItemsLoader struct {
	ExchangeItems       map[uint32]ExchangeItem
	BuffItems           map[uint32]BuffItem
	BuffItemsByName     map[string]BuffItemByName
	ExchangeItemsByName map[string]ExchangeItem
	Items               map[uint32]Items
	ItemsByName         map[string]ItemsByName
	Monsters            map[uint32]MonsterInfo
}

func (i *ItemsLoader) GetItemName(itemId uint32) string {
	if i.ExchangeItems != nil {
		if _, ok := i.ExchangeItems[itemId]; ok {
			return i.ExchangeItems[itemId].NameZh
		}
		log.Warnf("item id %d not found", itemId)
	}
	return ""
}

func (i *ItemsLoader) GetItemIdByName(itemName string) uint32 {
	if i.ExchangeItemsByName != nil {
		if _, ok := i.ExchangeItemsByName[itemName]; ok {
			id, _ := strconv.ParseUint(i.ExchangeItemsByName[itemName].Id.String(), 10, 32)
			return uint32(id)
		}
		log.Warnf("item name %s not found", itemName)
	}
	return 0
}

func (i *ItemsLoader) GetItemCat(itemId uint32) uint32 {
	if i.ExchangeItems != nil {
		if _, ok := i.ExchangeItems[itemId]; ok {
			catInt, _ := strconv.ParseUint(i.ExchangeItems[itemId].Category, 10, 32)
			return uint32(catInt)
		}
		log.Warnf("item id %d not found", itemId)
	}
	return 0
}

func (i *ItemsLoader) GetBuffNameById(buffId uint32) string {
	if i.BuffItems != nil {
		if _, ok := i.BuffItems[buffId]; ok {
			return i.BuffItems[buffId].BuffName
		}
		log.Warnf("item id %d not found", buffId)
	}
	return ""
}

func loadBuff(buffJsonPath string) (map[uint32]BuffItem, map[string]BuffItemByName) {
	var b []byte
	var err error
	if buffJsonPath != "" {
		fName := buffJsonPath
		jFile, err := os.Open(fName)
		if err != nil {
			log.Errorf("failed to open %s: %s", fName, err)
		}
		b, _ = ioutil.ReadAll(jFile)
	} else {
		b = data.BuffJson
	}

	var jLoader map[uint32]BuffItem
	jNameLoader := map[string]BuffItemByName{}
	err = json.Unmarshal(b, &jLoader)
	if err != nil {
		log.Errorf("failed load buff json: %s", err)
	}
	log.Infof("loaded %d buff items", len(jLoader))
	for _, val := range jLoader {
		if jNameLoader[val.BuffName].Items != nil {
			item := jNameLoader[val.BuffName].Items
			item = append(item, val)
			jNameLoader[val.BuffName] = BuffItemByName{
				Items: item,
			}
		} else {
			jNameLoader[val.BuffName] = BuffItemByName{Items: []BuffItem{val}}
		}
	}
	return jLoader, jNameLoader
}

func loadExchangeItems(itemJsonPath string) map[uint32]ExchangeItem {
	var b []byte
	var err error
	if itemJsonPath != "" {
		fName := itemJsonPath
		jsonFile, err := os.Open(fName)
		if err != nil {
			log.Errorf("failed to open %s: %s", fName, err)
		}
		b, _ = ioutil.ReadAll(jsonFile)
	} else {
		b = data.ExchangeJson
	}

	var jLoader map[uint32]ExchangeItem
	err = json.Unmarshal(b, &jLoader)
	if err != nil {
		log.Errorf("loading exchange items with error: %s", err)
	}
	log.Infof("loaded %d exchange items", len(jLoader))
	return jLoader
}

func loadItems(itemJsonPath string) (map[uint32]Items, map[string]ItemsByName) {
	var b []byte
	var err error
	if itemJsonPath != "" {
		fName := "items.json"
		fName = itemJsonPath
		jFile, err := os.Open(fName)
		if err != nil {
			log.Errorf("failed to open %s: %s", fName, err)
		}
		b, _ = ioutil.ReadAll(jFile)
	} else {
		b = data.ItemsJson
	}

	var jLoader map[uint32]Items
	jNameLoader := map[string]ItemsByName{}
	err = json.Unmarshal(b, &jLoader)
	if err != nil {
		log.Errorf("failed to load items json: %s", err)
	}
	log.Infof("loaded %d items", len(jLoader))
	for _, val := range jLoader {
		if jNameLoader[val.NameZh].Items != nil {
			item := jNameLoader[val.NameZh].Items
			item = append(item, val)
			jNameLoader[val.NameZh] = ItemsByName{
				Items: item,
			}
		} else {
			jNameLoader[val.NameZh] = ItemsByName{Items: []Items{val}}
		}
	}
	return jLoader, jNameLoader
}

func NewItemsLoader(exchangeItemJsonPath, buffJsonPath, itemsJsonPath string) *ItemsLoader {
	buffs, buffNames := loadBuff(buffJsonPath)
	items, itemNames := loadItems(itemsJsonPath)
	loader := &ItemsLoader{
		ExchangeItems:   loadExchangeItems(exchangeItemJsonPath),
		BuffItems:       buffs,
		BuffItemsByName: buffNames,
		Items:           items,
		ItemsByName:     itemNames,
	}
	itemByName := map[string]ExchangeItem{}
	for _, item := range loader.ExchangeItems {
		itemByName[item.NameZh] = item
	}
	loader.ExchangeItemsByName = itemByName
	return loader
}
