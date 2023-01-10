package utils

import (
	"encoding/json"
	"io"
	"os"

	"ROMProject/data"
	log "github.com/sirupsen/logrus"
)

type MonsterInfo struct {
	LoadShowSize     json.Number   `json:"LoadShowSize"`
	NameZh           string        `json:"NameZh"`
	Guild            string        `json:"Guild"`
	HeadDefaultColor json.Number   `json:"HeadDefaultColor"`
	Type             string        `json:"Type"`
	Atk              json.Number   `json:"Atk"`
	Position         string        `json:"Position"`
	MAtk             json.Number   `json:"MAtk"`
	Icon             string        `json:"Icon"`
	Zone             string        `json:"Zone"`
	Id               int           `json:"id"`
	Flee             json.Number   `json:"Flee"`
	Level            json.Number   `json:"Level"`
	AccessRange      json.Number   `json:"AccessRange"`
	Desc             string        `json:"Desc"`
	Shape            string        `json:"Shape"`
	MDef             json.Number   `json:"MDef"`
	Nature           string        `json:"Nature"`
	LoadShowRotate   json.Number   `json:"LoadShowRotate"`
	Race             string        `json:"Race"`
	Move             json.Number   `json:"move"`
	SpawnSE          string        `json:"SpawnSE"`
	Hp               json.Number   `json:"Hp"`
	LoadShowPose     []json.Number `json:"LoadShowPose"`
	Features         json.Number   `json:"Features"`
	AtkSpd           float64       `json:"AtkSpd"`
	MoveSpdRate      json.Number   `json:"MoveSpdRate"`
	MoveSpd          json.Number   `json:"MoveSpd"`
	Hit              json.Number   `json:"Hit"`
	Behaviors        json.Number   `json:"Behaviors"`
	Def              json.Number   `json:"Def"`
	Body             json.Number   `json:"Body"`
	DeathEffect      string        `json:"DeathEffect"`
}

func MonsterParser(monsterJson string) map[uint32]MonsterInfo {
	var b []byte
	var jLoader map[uint32]MonsterInfo
	if monsterJson != "" {
		monsterFile, err := os.Open(monsterJson)
		if err != nil {
			log.Error("reading monster json file failed using default:%s", err)
			b = data.MonstersJson
		} else {
			defer monsterFile.Close()
			b, _ = io.ReadAll(monsterFile)
		}
	} else {
		b = data.MonstersJson
	}

	err := json.Unmarshal(b, &jLoader)
	if err != nil {
		log.Errorf("failed to unmarshal monster json: %s", err)
	}
	return jLoader
}
