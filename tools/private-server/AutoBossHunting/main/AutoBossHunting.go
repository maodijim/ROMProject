package main

import (
	"context"
	"flag"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/AutoBossHunting"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

type HiddenMVP struct {
	Info        utils.MonsterInfo
	RespawnTime time.Time
}

var (
	g                     *gameConnection.GameConnection
	fightCtx, fightCancel = context.WithCancel(context.Background())
	flyWingUseCount       = 0
	maxFlyWingUseCount    = 20
	fightStar             = false
	pickupCount           = uint32(0)
	maxPickupCount        = uint32(100)
	fighting              = false
	lavaGemCount          = uint32(0)
	lastPosUpdate         = time.Now()
	lastPos               Cmd.ScenePos
	flyMutex              *sync.Mutex
	flyMutexWait          = float64(7000)
	TargetMonster         Cmd.BossInfoItem
	HuntingCount          = uint32(0)
	MitionCompelete       = bool(true)
	LastHp                = int32(0)
	tempMHP               = int32(0)
	tempUHP               = int32(0)
	HuntHidMVP            = bool(false)
	HiddenMVPList         map[string]HiddenMVP
	TargetHiddenMVP       HiddenMVP
	StartTime             = time.Now()
)

const (
	ver = "0.1.2"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func main() {
	log.Infof("自动BOSS狩猎版本 %s", ver)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()
	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}

	bossHuntTask := AutoBossHunting.NewBossHuntTask(context.Background(), g)
	bossHuntTask.Start()
	select {}
}
