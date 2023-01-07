package utils

import (
	"encoding/json"
	"io"
	"os"

	"ROMProject/data"
	log "github.com/sirupsen/logrus"
)

type SkillItem struct {
	AttackStatus          string `json:"AttackStatus"`
	AttackEp              string `json:"Attack_EP"`
	CCT                   string `json:"CCT"`
	CD                    string `json:"CD"`
	CastAct               string `json:"CastAct"`
	Camps                 string `json:"Camps"`
	Cost                  string `json:"Cost"`
	DelayCd               string `json:"DelayCD"`
	Duration              string `json:"duration"`
	Elementparam          string `json:"elementparam"`
	Effect                string `json:"effect"`
	FCT                   string `json:"FCT"`
	FieldareaCannotImmune string `json:"fieldarea_cannot_immune"`
	FireEp                string `json:"Fire_EP"`
	ForbiUse              string `json:"ForbiUse"`
	Icon                  string `json:"Icon"`
	ItemId                string `json:"itemID"`
	IncludeSelf           string `json:"include_self"`
	LaunchRange           string `json:"Launch_Range"`
	LaunchType            string `json:"Launch_Type"`
	Level                 string `json:"Level"`
	Logic                 string `json:"Logic"`
	NameZh                string `json:"NameZh"`
	NextId                string `json:"NextID"`
	NextNewId             string `json:"NextNewID"`
	RollType              string `json:"RollType"`
	SeAttack              string `json:"SE_attack"`
	SeHit                 string `json:"SE_hit"`
	SelectTarget          string `json:"select_target"`
	SelectHide            string `json:"select_hide"`
	SeMiss                string `json:"SE_miss"`
	SeCast                string `json:"SE_cast"`
	SkillId               string `json:"skillid"`
	SkillHit              string `json:"skillHit"`
	SkillType             string `json:"SkillType"`
	TargetEP              string `json:"Target_EP"`
	DamChangePer          string `json:"damChangePer"`
	Id                    string `json:"id"`
	Interval              string `json:"interval"`
	IsCountTrap           string `json:"isCountTrap"`
	ItemType              string `json:"itemtype"`
	NoAction              string `json:"no_action"`
	NoSelect              string `json:"no_select"`
	NoTarget              string `json:"no_target"`
	MaxCount              string `json:"max_count"`
	Menu                  string `json:"menu"`
	Type                  string `json:"type"`
	TeamRange             string `json:"team_range"`
	TrapEffect            string `json:"trap_effect"`
	Range                 string `json:"range"`
	RangeNum              string `json:"range_num"`
	RiskId                string `json:"riskId"`
	ProType               string `json:"ProType"`
	SP                    string `json:"sp"`
	Speed                 string `json:"speed"`
	Spotter               string `json:"spotter"`
	Time                  string `json:"time"`
	Value                 string `json:"value"`
}

func readSkillItem(skillJson string) map[uint32]SkillItem {
	var b []byte
	var err error
	if skillJson != "" {
		fName := skillJson
		jsonFile, err := os.Open(skillJson)
		if err != nil {
			log.Errorf("failed to open %s: %s", fName, err)
			b = data.SkillsJson
		} else {
			b, _ = io.ReadAll(jsonFile)
		}
	} else {
		b = data.SkillsJson
	}

	var jLoader map[uint32]SkillItem
	err = json.Unmarshal(b, &jLoader)
	if err != nil {
		log.Errorf("loading exchange items with error: %s", err)
	}
	log.Infof("loaded %d skill items", len(jLoader))
	return jLoader
}

func NewSkillParser(skillJson string) map[uint32]SkillItem {
	skills := readSkillItem(skillJson)
	return skills
}
