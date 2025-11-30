package config

import (
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
)

type TradeMonitorConfig struct {
	MonitorInterval       int      `yaml:"monitorInterval" json:"monitorInterval"` // in seconds
	WatchItems            []string `yaml:"watchItems" json:"watchItems"`
	WatchCategories       []string `yaml:"watchCategories" json:"watchCategories"`
	ElasticsearchHostPort string   `yaml:"elasticsearchHostPort" json:"elasticsearchHostPort"`
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
	if val, ok := config["monitorInterval"].(float64); ok {
		c.MonitorInterval = int(val)
	}
	if val, ok := config["watchItems"].([]interface{}); ok {
		items := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				items[i] = str
			}
		}
		c.WatchItems = items
	}
	if val, ok := config["watchCategories"].([]interface{}); ok {
		categories := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				categories[i] = str
			}
		}
		c.WatchCategories = categories
	}
	if val, ok := config["elasticsearchHostPort"].(string); ok {
		c.ElasticsearchHostPort = val
	}
	return *c
}

func (c *TradeMonitorConfig) GetDefault() TradeMonitorConfig {
	return TradeMonitorConfig{
		MonitorInterval: 30,
		WatchItems:      []string{},
		WatchCategories: utils.GetMapKeys(gameTypes.TradeZhCategories),
	}
}
