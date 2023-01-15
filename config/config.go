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

type ServerConfigs struct {
	AuthServer     string `yaml:"authServer"`
	ZoneId         uint32 `yaml:"zoneId"`
	Char           uint   `yaml:"char"`
	AutoCreateChar bool   `yaml:"autoCreateChar"`
	// if not set, use random string during auto create char
	CharacterName string `yaml:"characterName"`
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
		log.Exit(3)
	}

	return configs
}
