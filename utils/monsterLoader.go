package utils

import (
	"encoding/json"
	"io"
	"os"

	"ROMProject/data"
	log "github.com/sirupsen/logrus"
)

type MonsterInfo struct {
	LoadShowSize     int     `json:"LoadShowSize"`
	NameZh           string  `json:"NameZh"`
	Guild            string  `json:"Guild"`
	HeadDefaultColor int     `json:"HeadDefaultColor"`
	Type             string  `json:"Type"`
	Atk              int     `json:"Atk"`
	Position         string  `json:"Position"`
	MAtk             int     `json:"MAtk"`
	Icon             string  `json:"Icon"`
	Zone             string  `json:"Zone"`
	Id               int     `json:"id"`
	Flee             int     `json:"Flee"`
	Level            int     `json:"Level"`
	AccessRange      int     `json:"AccessRange"`
	Desc             string  `json:"Desc"`
	Shape            string  `json:"Shape"`
	MDef             int     `json:"MDef"`
	Nature           string  `json:"Nature"`
	LoadShowRotate   int     `json:"LoadShowRotate"`
	Race             string  `json:"Race"`
	Move             int     `json:"move"`
	SpawnSE          string  `json:"SpawnSE"`
	Hp               int     `json:"Hp"`
	LoadShowPose     []int   `json:"LoadShowPose"`
	Features         int     `json:"Features"`
	AtkSpd           float64 `json:"AtkSpd"`
	MoveSpdRate      int     `json:"MoveSpdRate"`
	MoveSpd          int     `json:"MoveSpd"`
	Hit              int     `json:"Hit"`
	Behaviors        int     `json:"Behaviors"`
	Def              int     `json:"Def"`
	Body             int     `json:"Body"`
	DeathEffect      string  `json:"DeathEffect"`
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
