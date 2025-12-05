package config

import (
	"ROMProject/utils"
	"bytes"
	"io"
	"os"

	"ROMProject/data"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type EsConfig struct {
	Urls []string `yaml:"urls"`
}

type EnchantConfig struct {
	AutoSave        bool              `yaml:"autoSave" json:"自动保存"`
	EnchantType     string            `yaml:"enchantType" json:"附魔类型(高级/中级/低级)"`
	EnchantEquipPos string            `yaml:"enchantEquipPos" json:"附魔部位(武器/副手/盔甲/鞋子/披风/头饰/饰品1/饰品2/背部/尾部/脸部/嘴部)"`
	Condition       EnchantCondition  `yaml:"condition" json:"附魔目标条件"`
	AutoBuyCoin     AutoBuyCoinConfig `yaml:"autoBuyCoin" json:"自动购买金币"`
	EnchantCount    uint32            `yaml:"enchantCount" json:"附魔次数"`
}

func (c *EnchantConfig) ParseFromInterface(config map[string]any) EnchantConfig {
	utils.ParseConfigFromInterface(config, c)
	return *c
}

type AutoBuyCoinConfig struct {
	Enable        bool  `yaml:"enable" json:"enable"`
	MinZenyToKeep int64 `yaml:"minZenyToKeep" json:"min_zeny_to_keep"`
	NumCoinsToBuy int   `yaml:"numCoinsToBuy" default:"1000" json:"num_coins_to_buy"`
}

type EnchantCondition struct {
	Attributes []string `yaml:"attributes" json:"属性(必须跟游戏里面的属性描述一样要不然可能无法识别 比如 '暴伤% > 80')"`
	Extras     []string `yaml:"extras" json:"词条(比如 '尖锐4')"`
}

type HuntConfig struct {
	CarryTeam      bool     `yaml:"CarryTeam" json:"组队一起飞"`
	UseDoubleEXP   bool     `yaml:"UseDoubleEXP" json:"使用洋洋"`
	GameDuration   int      `yaml:"GameDuration" json:"狩猎时长(小时) 0为无限制"`
	Mini           []string `yaml:"Mini" json:"Mini狩猎清单"`
	MVP            []string `yaml:"MVP" json:"MVP狩猎清单"`
	HMVP           []string `yaml:"HMVP" json:"HMVP狩猎清单"`
	PrepEliteCD    int      `yaml:"PrepEliteCD" json:"备战精英技能冷却时间(秒)"`
	TimerFly       int      `yaml:"TimerFly" json:"固定时间使用翅膀(0为不使用)"`
	TargetMonsters []string `yaml:"TargetMonsters" json:"狩猎魔物清单"`
	TargetItems    []string `yaml:"TargetItems" json:"狩猎物品清单"`
	Map            string   `yaml:"TargetMap" json:"狩猎地图"`
	NatureType     string   `yaml:"NatureType" json:"使用属性类型"`
}

func (c *HuntConfig) BossInfoParseFromInterface(config map[string]interface{}) HuntConfig {
	if val, ok := config["CarryTeam"].(bool); ok {
		c.CarryTeam = val
	}
	if val, ok := config["GameDuration"].(float64); ok {
		c.GameDuration = int(val)
	}
	if val, ok := config["Mini"].([]interface{}); ok {
		mini := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				mini[i] = str
			}
		}
		c.Mini = mini
	}
	if val, ok := config["MVP"].([]interface{}); ok {
		mvp := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				mvp[i] = str
			}
		}
		c.MVP = mvp
	}
	if val, ok := config["HMVP"].([]interface{}); ok {
		hmvp := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				hmvp[i] = str
			}
		}
		c.HMVP = hmvp
	}
	if val, ok := config["PrepEliteCD"].(float64); ok {
		c.PrepEliteCD = int(val)
	}
	if val, ok := config["UseDoubleEXP"].(bool); ok {
		c.UseDoubleEXP = val
	}
	return *c
}
func (c *HuntConfig) HuntInfoParseFromInterface(config map[string]interface{}) HuntConfig {
	if val, ok := config["PrepEliteCD"].(float64); ok {
		c.PrepEliteCD = int(val)
	}
	if val, ok := config["TargetMonsters"].([]interface{}); ok {
		TargetMonsters := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				TargetMonsters[i] = str
			}
		}
		c.TargetMonsters = TargetMonsters
	}
	if val, ok := config["TargetItems"].([]interface{}); ok {
		TargetItems := make([]string, len(val))
		for i, v := range val {
			if str, ok := v.(string); ok {
				TargetItems[i] = str
			}
		}
		c.TargetItems = TargetItems
	}
	if val, ok := config["Map"].(string); ok {
		c.Map = string(val)
	}
	if val, ok := config["TimerFly"].(float64); ok {
		c.TimerFly = int(val)
	}
	if val, ok := config["UseDoubleEXP"].(bool); ok {
		c.UseDoubleEXP = val
	}
	if val, ok := config["NatureType"].(string); ok {
		c.NatureType = string(val)
	}
	return *c
}

func (c *HuntConfig) GetDefault() HuntConfig {
	return HuntConfig{
		CarryTeam:    false,
		GameDuration: 0,
		Mini: []string{
			"狸猫", "蓝疯兔", "波利之王", "摇滚蝗虫", "幽灵波利", "蛙王", "直升机哥布灵", "龙蝇", "流浪之狼",
			"枯树精", "狮鹫兽", "安毕斯", "妖君", "兽人婴儿", "南瓜先生", "半龙人", "草精",
			"鹗枭首领", "爱丽丝女仆", "艾斯恩魔女", "弑神者", "迷幻之王",
		},
		MVP: []string{
			"天使波利", "黄金虫", "恶魔波利", "海盗之王", "海神", "哥布灵首领", "蜂后", "蚁后",
			"皮里恩", "虎王", "俄塞里斯", "月夜猫", "兽人英雄", "犬妖首领", "死灵", "阿特罗斯",
			"兽人酋长", "鹗枭男爵", "血腥骑士", "巴风特", "黑暗之王",
		},
		HMVP: []string{
			"卡仑", "狼外婆",
		},
		PrepEliteCD: 30,
	}
}

type ServerConfigs struct {
	AuthServer     string `yaml:"authServer"`
	AuthPass       string `yaml:"authPass"`
	ZoneId         uint32 `yaml:"zoneId"`
	Char           uint   `yaml:"char"`
	AutoCreateChar bool   `yaml:"autoCreateChar"`
	// if not set, use random string during auto create char
	CharacterName string        `yaml:"characterName"`
	EnchantConfig EnchantConfig `yaml:"enchantConfig"`
	HuntConfig    HuntConfig    `yaml:"HuntConfig"`
	// Username and password for login if accId is not set
	Username           string             `yaml:"username"`
	Password           string             `yaml:"password"`
	Region             int                `yaml:"region"`
	Version            string             `yaml:"version"`
	ServerId           uint32             `yaml:"serverId"`
	AccId              uint64             `yaml:"accId"`
	Ip                 string             `yaml:"ip"`
	Domain             string             `yaml:"domain"`
	LineGrp            string             `yaml:"lineGrp"`
	Device             string             `yaml:"device"`
	DeviceId           string             `yaml:"deviceId"`
	ClientVer          string             `yaml:"clientVer"`
	LangZone           uint32             `yaml:"langZone"`
	IpPort             string             `yaml:"ipPort"`
	GameServer         string             `yaml:"gameServer"`
	Phone              string             `yaml:"phone"`
	SafeDevice         string             `yaml:"safeDevice"`
	Sha1Str            string             `yaml:"sha1Str"`
	AccessToken        string             `yaml:"accessToken"`
	ResVer             string             `yaml:"resourceVer"`
	PlatformVer        string             `yaml:"platV"`
	AppPreVer          uint32             `yaml:"appPreVer"`
	Lang               uint32             `yaml:"lang"`
	Model              string             `yaml:"model"`
	PhoneVer           string             `yaml:"phoneVer"`
	Authoriz           string             `yaml:"authoriz"`
	AuthParams         map[string]string  `yaml:"authParams"`
	EsConfig           EsConfig           `yaml:"elasticsearch"`
	TeamConfig         TeamConfig         `yaml:"team"`
	Loglines           int                `yaml:"-"`
	ChatMaxSize        int                `yaml:"-"`
	TradeMonitorConfig TradeMonitorConfig `yaml:"tradeMonitorConfig"`
}

func (s *ServerConfigs) GetChatMaxSize() int {
	if s.ChatMaxSize <= 0 {
		return 500
	}
	return s.ChatMaxSize
}

func (s *ServerConfigs) SetTeamLeader(name string) {
	s.TeamConfig.LeaderName = name
}

func (s *ServerConfigs) SetFollowTeamLeader(yes bool) {
	s.TeamConfig.FollowTeamLeader = yes
}

func parseConfigYaml(r io.Reader, sc *ServerConfigs) error {
	decoder := yaml.NewDecoder(r)
	err := decoder.Decode(sc)
	if err != nil {
		return err
	}

	if sc.Region < 1 {
		sc.Region = 1
	} else {
		sc.Region -= 1
	}

	if sc.EnchantConfig.AutoBuyCoin.Enable {
		if sc.EnchantConfig.AutoBuyCoin.NumCoinsToBuy <= 0 {
			sc.EnchantConfig.AutoBuyCoin.NumCoinsToBuy = 1000
		}
		if sc.EnchantConfig.AutoBuyCoin.MinZenyToKeep <= 0 {
			sc.EnchantConfig.AutoBuyCoin.MinZenyToKeep = 100000000
		}
	}

	return nil
}

func NewServerConfigs(configYaml string) *ServerConfigs {
	configPath := configYaml
	configs := &ServerConfigs{}
	err := parseConfigYaml(bytes.NewReader(data.ConfigYml), configs)
	if err != nil {
		log.Fatalf("failed to parse default config yaml: %v", err)
	}
	if configYaml == "" {
		configPath = "config.yml"
	}
	f, err := os.Open(configYaml)
	if err != nil {
		log.Errorf("failed to load %s: %s", configPath, err)
		log.Exit(2)
	}
	defer f.Close()
	err = parseConfigYaml(f, configs)
	if err != nil {
		log.Errorf("parse config yaml failed: %s", err)
		return configs
	}

	return configs
}
