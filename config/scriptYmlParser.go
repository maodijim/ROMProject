package config

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	ActionList = map[string]interface{}{
		"use_skill":   []Cmd.SkillBroadcastUserCmd{},
		"move":        Cmd.ScenePos{},
		"talk_to_npc": Cmd.MapNpc{},
		"use_item":    "",
		"auto_attack": "",
	}
)

type ScriptAction struct {
	ActionName  string      `yaml:"action_name"`
	ActionValue interface{} `yaml:"action_value"`
}

type ScriptActions struct {
	AutoSubmitWantedQuest *bool          `yaml:"auto_submit_wanted_quest"`
	AutoTeamExpFuben      *bool          `yaml:"auto_team_exp_fuben"`
	MapName               *string        `yaml:"mapName"`
	MonsterName           []string       `yaml:"monsterName"`
	UseFlyWing            *bool          `yaml:"useFlyWing"`
	Actions               []ScriptAction `yaml:"actions"`
}

func (s *ScriptActions) GetMapName() string {
	if s.MapName == nil {
		return ""
	}
	return *s.MapName
}

func (s *ScriptActions) ShouldUseFlyWing() bool {
	if s.UseFlyWing == nil {
		return false
	}
	return *s.UseFlyWing
}

func ScriptParser(scriptPath string) *ScriptActions {
	if scriptPath == "" {
		scriptPath = "script.yml"
	}
	var scripActions *ScriptActions
	f, err := os.Open(scriptPath)
	if err != nil {
		log.Errorf("failed to open %s: %v", scriptPath, err)
		log.Exit(2)
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)

	err = decoder.Decode(&scripActions)
	if err != nil {
		log.Errorf("parse script actions failed: %v", err)
		log.Exit(3)
	}

	return scripActions
}
