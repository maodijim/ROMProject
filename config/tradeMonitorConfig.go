package config

import (
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
)

type TradeMonitorConfig struct {
	MonitorInterval       int      `yaml:"monitorInterval" json:"监控间隔"` // in seconds
	WatchItems            []string `yaml:"watchItems" json:"监控物品"`
	WatchCategories       []string `yaml:"watchCategories" json:"监控类别"`
	ElasticsearchHostPort string   `yaml:"elasticsearchHostPort" json:"elasticsearchHostPort"`
	NumberWorkers         int      `yaml:"numberWorkers" json:"工作线程数"`
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
