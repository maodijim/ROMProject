package main

import (
	"flag"
	"time"

	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	autoScriptVer = "0.0.1"
)

func init() {
	// log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func main() {
	log.Infof("ROM auto script version %s", autoScriptVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "exchangeItems.json", "Exchange Item Json Path")
	scriptYml := flag.String("script", "script.yml", "yaml file of the script actions")
	buffFile := flag.String("buffPath", "buff.json", "Buff Json Path")
	skillJson := flag.String("skillJson", "skills.yml", "json file of skills")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()

	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	scriptActions := config.ScriptParser(*scriptYml)
	skills := utils.NewSkillParser(*skillJson)
	log.Infof("%v", scriptActions)
	conf := config.NewServerConfigs(*confFile)
	gameConnect := gameConnection.NewConnection(conf, skills, items)
	gameConnect.ShouldChangeScene = true

	if *enableDebug {
		gameConnect.DebugMsg = *enableDebug
		log.SetLevel(log.DebugLevel)
	}

	gameConnect.GameServerLogin()
	quit := make(chan bool)
	gameConnect.CheckForFubenInviteInBackground(quit)
	disable := make(chan bool)
	gameConnect.EnableAutoAttack([]string{"all"}, disable)
	// gameConnect.InviteTeamExpFuben()
	gameConnect.AutoSubmitWantedQuest()
	go func() {
		time.Sleep(10 * time.Second)
		if !gameConnect.IsCookWorkSpaceInUse() {
			log.Infof("没有宠物烹饪打工, 开始任命一个")
			gameConnect.BattlePetToWork()
		} else {
			gameConnect.GetCookWorkSpaceReward()
		}
		gameConnect.DailySignIn()
	}()
	for {
		if gameConnect.Role.GetInGame() {

			// log.Infof("附近的NPCS")
			// for _, npc := range gameConnect.MapNpcs {
			//	log.Infof("NPC: %s, 血量: %d",
			//		npc.GetName(),
			//		utils.GetNpcAttrValByType(npc.GetAttrs(), Cmd.EAttrType_EATTRTYPE_HP),
			//	)
			// }
			// log.Infof("有%d只NPC", len(gameConnect.MapNpcs))

			time.Sleep(10 * time.Second)

		} else {
			log.Infof("等待角色进入游戏")
			time.Sleep(10 * time.Second)
		}
	}
}
