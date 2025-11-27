package config

import (
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
	AutoSave        bool              `yaml:"autoSave" json:"autoSave"`
	EnchantType     string            `yaml:"enchantType" json:"enchantType"`
	EnchantEquipPos string            `yaml:"enchantEquipPos" json:"enchantEquipPos"`
	Condition       EnchantCondition  `yaml:"condition" json:"condition"`
	AutoBuyCoin     AutoBuyCoinConfig `yaml:"autoBuyCoin" json:"autoBuyCoin"`
	EnchantCount    uint32            `yaml:"enchantCount" json:"enchantCount"`
}

func (c *EnchantConfig) ParseFromInterface(config map[string]interface{}) EnchantConfig {
	if val, ok := config["autoSave"].(bool); ok {
		c.AutoSave = val
	}
	if val, ok := config["enchantType"].(string); ok {
		c.EnchantType = val
	}
	if val, ok := config["enchantEquipPos"].(string); ok {
		c.EnchantEquipPos = val
	}
	if val, ok := config["enchantCount"].(uint32); ok {
		c.EnchantCount = val
	}
	if val, ok := config["condition"].(map[string]interface{}); ok {
		if attrs, ok := val["attributes"].([]interface{}); ok {
			attr := make([]string, len(attrs))
			for i, v := range attrs {
				if str, ok := v.(string); ok {
					attr[i] = str
				}
			}
			c.Condition.Attributes = attr
		}
		if extras, ok := val["extras"].([]interface{}); ok {
			extr := make([]string, len(extras))
			for i, v := range extras {
				if str, ok := v.(string); ok {
					extr[i] = str
				}
			}
			c.Condition.Extras = extr
		}
	}
	return *c
}

type AutoBuyCoinConfig struct {
	Enable        bool  `yaml:"enable" json:"enable"`
	MinZenyToKeep int64 `yaml:"minZenyToKeep" json:"min_zeny_to_keep"`
	NumCoinsToBuy int   `yaml:"numCoinsToBuy" default:"1000" json:"num_coins_to_buy"`
}

type EnchantCondition struct {
	Attributes []string `yaml:"attributes" json:"attributes"`
	Extras     []string `yaml:"extras" json:"extras"`
}

type HuntConfig struct {
	CarryTeam    bool     `yaml:"CarryTeam" json:"CarryTeam"`
	GameDuration int      `yaml:"GameDuration" json:"GameDuration"`
	Mini         []string `yaml:"Mini" json:"Mini"`
	MVP          []string `yaml:"MVP" json:"MVP"`
	HMVP         []string `yaml:"HMVP" json:"HMVP"`
	PrepEliteCD  int      `yaml:"PrepEliteCD" json:"PrepEliteCD"`
}

func (c *HuntConfig) ParseFromInterface(config map[string]interface{}) HuntConfig {
	if val, ok := config["CarryTeam"].(bool); ok {
		c.CarryTeam = val
	}
	if val, ok := config["GameDuration"].(int); ok {
		c.GameDuration = val
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
	}
	if val, ok := config["PrepEliteCD"].(int); ok {
		c.PrepEliteCD = val
	}
	return *c
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
	Username    string            `yaml:"username"`
	Password    string            `yaml:"password"`
	Region      int               `yaml:"region"`
	Version     string            `yaml:"version"`
	ServerId    uint32            `yaml:"serverId"`
	AccId       uint64            `yaml:"accId"`
	Ip          string            `yaml:"ip"`
	Domain      string            `yaml:"domain"`
	LineGrp     string            `yaml:"lineGrp"`
	Device      string            `yaml:"device"`
	DeviceId    string            `yaml:"deviceId"`
	ClientVer   string            `yaml:"clientVer"`
	LangZone    uint32            `yaml:"langZone"`
	IpPort      string            `yaml:"ipPort"`
	GameServer  string            `yaml:"gameServer"`
	Phone       string            `yaml:"phone"`
	SafeDevice  string            `yaml:"safeDevice"`
	Sha1Str     string            `yaml:"sha1Str"`
	AccessToken string            `yaml:"accessToken"`
	ResVer      string            `yaml:"resourceVer"`
	PlatformVer string            `yaml:"platV"`
	AppPreVer   uint32            `yaml:"appPreVer"`
	Lang        uint32            `yaml:"lang"`
	Model       string            `yaml:"model"`
	PhoneVer    string            `yaml:"phoneVer"`
	Authoriz    string            `yaml:"authoriz"`
	AuthParams  map[string]string `yaml:"authParams"`
	EsConfig    EsConfig          `yaml:"elasticsearch"`
	TeamConfig  TeamConfig        `yaml:"team"`
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
