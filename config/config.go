package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

type EsConfig struct {
	Urls []string `yaml:"urls"`
}

type ServerConfigs struct {
	AuthServer  string            `yaml:"authServer"`
	ZoneId      uint32            `yaml:"zoneId"`
	Char        uint              `yaml:"char"`
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

func NewServerConfigs(configYaml string) *ServerConfigs {
	configPath := configYaml
	config := &ServerConfigs{}
	if configYaml == "" {
		configPath = "config.yml"
	}
	f, err := os.Open(configYaml)
	if err != nil {
		log.Errorf("failed to load %s: %s", configPath, err)
		log.Exit(2)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(config)
	if err != nil {
		log.Errorf("parse config yaml failed: %s", err)
		log.Exit(3)
	}

	if config.Region < 1 {
		config.Region = 1
	} else {
		config.Region -= 1
	}

	return config
}
