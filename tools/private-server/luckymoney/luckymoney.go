package main

import (
	"context"
	"flag"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	scriptVer = "1.0.0"
)

var (
	items   = &utils.ItemsLoader{}
	g       = &gameConnection.GameConnection{}
	npcList = map[string][]int{}
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func main() {
	log.Infof("ROM auto money up version: %s", scriptVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()

	items = utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	g.GameServerLogin()
	_ = g.GetAllPackItems()
	g.ChangeMap(g.Role.GetMapId())
	if g.Role.GetMapId() != gameTypes.MapId_MagmaDungeon2F.Uint32() {
		g.ExitMapWait(gameTypes.MapId_MagmaDungeon2F.Uint32())
	}

	// 溶洞进门右边 target:{x:47199 y:926 z:29998}
	g.MoveChartWait(g.ParsePos(47199, 926, 29998))

	cancelAttackCtx, cancelAttack := context.WithCancel(context.Background())
	g.EnableAutoAttack(cancelAttackCtx, "烈鬃马  ")
	ticket := time.NewTicker(10 * time.Second)
	chainTick := time.NewTicker(60 * time.Second)
	goBackTick := time.NewTicker(5 * time.Minute)
	lucky := g.FindPackItemByName("B格猫皇家红包1.0", Cmd.EPackType_EPACKTYPE_MAIN)
	startLuckyMoneyCount := lucky.GetBase().GetCount()
	startTime := time.Now()
	for {
		select {
		case <-goBackTick.C:
			log.Infof("归位")
			cancelAttack()
			g.MoveChartWait(g.ParsePos(47199, 926, 29998))
			cancelAttackCtx, cancelAttack = context.WithCancel(context.Background())
			g.EnableAutoAttack(cancelAttackCtx, "烈鬃马  ")
		case <-ticket.C:
			lucky = g.FindPackItemByName("B格猫皇家红包1.0", Cmd.EPackType_EPACKTYPE_MAIN)
			log.Infof("当前红包数量: %d, 打了 %d个红包 耗时: %f 分钟",
				lucky.GetBase().GetCount(),
				lucky.GetBase().GetCount()-startLuckyMoneyCount,
				time.Since(startTime).Minutes(),
			)
		case <-chainTick.C:
			if g.GetCurrentHp() == 0 {
				log.Errorf("死了")
				cancelAttack()
				g.Reconnect()
				time.Sleep(5 * time.Second)
				g.MoveChartWait(g.ParsePos(47199, 926, 29998))
				cancelAttackCtx, cancelAttack = context.WithCancel(context.Background())
				g.EnableAutoAttack(cancelAttackCtx, "烈鬃马  ", "火焰精灵  ")
			}
			useChain()
			useWater()
			useFood()
		}
	}
}

func useChain() {
	useItem("锁链雷锭", "锁链雷锭")
}

func useFire() {
	useItem("火灵原石", "火灵原石")
}

func useWater() {
	useItem("水灵原石", "火灵原石")
}

func useFood() {
	useItem("敏捷料理A", "敏捷料理A")
	useItem("灵巧料理A", "灵巧料理A")
	useItem("幸运料理A", "幸运料理A")
}

func useItem(name, buffName string) {
	item := g.FindPackItemByName(name, Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		log.Warnf("未找到%s", name)
	} else {
		log.Infof("背包还有%s %d 个", name, item.GetBase().GetCount())
		hasBuff := g.GetBuffNameByRegex(buffName)
		if hasBuff != "" {
			log.Warnf("已经有%sbuff了", buffName)
			return
		}
		log.Infof("使用%s", name)
		g.UseItem(item.GetBase().GetGuid(), 1)
	}
}
