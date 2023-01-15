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
	csVer = "v0.0.1"
)

var (
	g                    *gameConnection.GameConnection
	startSockCount       = uint32(0)
	curSockCount         = uint32(0)
	flyWingNotFoundCount = 0
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func main() {
	log.Infof("ROM auto Christmas sock version: %s", csVer)
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
	start()
}

func start() {
	g.GameServerLogin()
	inGameTicker := time.NewTicker(2 * time.Second)
waitForInGame:
	for {
		select {
		case <-inGameTicker.C:
			log.Infof("等待角色进入游戏")
			if g.Role.GetInGame() {
				log.Infof("角色已进入游戏")
				inGameTicker.Stop()
				break waitForInGame
			}
		case <-time.After(15 * time.Second):
			log.Warnf("等待进入游戏超时")
			inGameTicker.Stop()
			break waitForInGame
		}
	}
	_ = g.GetAllPackItems()
	g.ChangeMap(g.Role.GetMapId())

	if g.Role.GetMapId() != gameTypes.MapId_ToyFactory1F.Uint32() {
		g.ExitMapWait(gameTypes.MapId_ToyFactory1F.Uint32())
	}

	// notification for new 玩具士兵★
	g.AddNotifier(gameTypes.NtfType_BossWorldNtf)
	g.AddNotifier(gameTypes.NtfType_UserActionDialog)
	bossNtf := g.Notifier(gameTypes.NtfType_BossWorldNtf)
	dialogNtf := g.Notifier(gameTypes.NtfType_UserActionDialog)
	buyFlyWingTicker := time.NewTicker(60 * time.Second)
	go func() {
		for {
			select {
			case newNtf := <-bossNtf:
				bossNtf, ok := newNtf.(*Cmd.WorldBossNtf)
				if !ok {
					log.Infof("Boss怪物诞生 %s", g.GetMonsterNameById(bossNtf.GetNpcid()))
				}
			case newNtf := <-dialogNtf:
				dialogNtf, ok := newNtf.(*Cmd.UserActionNtf)
				// 53170 玩具士兵★ 诞生通知
				if ok && dialogNtf.GetValue() == 53170 {
					log.Info("玩具士兵★ 诞生了")
				} else if ok {
					log.Warnf("对话框id %d", dialogNtf.GetValue())
				}
			case <-buyFlyWingTicker.C:
				wingCount := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
				if wingCount != nil && wingCount.GetCount() < 1000 {
					shopConfig, err := g.QueryShopConfig(gameTypes.ShopType_Item, 1)
					if err != nil {
						log.Errorf("查询商店配置失败 %s", err)
						continue
					}
					for _, item := range shopConfig.GetGoods() {
						if item.GetItemid() == 5024 {
							log.Infof("购买100苍蝇翅膀")
							g.BuyShopItem(item, 100)
						}
					}
				}
			}
		}
	}()

	startSockCount = getSockCount()
	for {
		findMonster()
		curSockCount = getSockCount()
		log.Infof("挂机获得圣诞袜子数量: %d; 现在有袜子: %d", curSockCount-startSockCount, curSockCount)
		// g.ExitMapWait(gameConnection.MapId_Yuno.Uint32())
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute*3)
		log.Infof("开始休息打怪 3分钟")
		g.EnableAutoAttack(ctx, "玩具士兵", "玩具士兵★", "玩偶熊")
	rest:
		for {
			select {
			case <-ctx.Done():
				log.Infof("停止休息")
				cancel()
				break rest
			}
		}
		// g.ExitMapWait(gameConnection.MapId_ToyFactory1F.Uint32())
	}
}

func useFlyWing() {
	g.UseFlyWing()
	item := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil && item.GetCount() > 0 {

		log.Infof("使用苍蝇翅膀 还有%d个", item.GetCount())

		go func() {
			for {
				select {
				case <-time.After(time.Second * 1):
					newCount := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN).GetCount()
					if newCount == item.GetCount() {
						if flyWingNotFoundCount > 10 {
							log.Warnf("使用苍蝇翅膀失败怀疑卡住了重连")
							flyWingNotFoundCount = 0
							g.Reconnect()
						}
						flyWingNotFoundCount++
					} else {
						flyWingNotFoundCount = 0
					}
				}
			}
		}()
		return
	}
	log.Infof("没有找到苍蝇翅膀")
	_ = g.GetAllPackItems()
}

func findMonster() {
	max := 30
	for count := 0; count < max; count++ {
		if count == max-1 {
			log.Warnf("累计%d次翅膀找不到怪物休息一下", max)
			break
		}
		disableAttack, cancelAttack := context.WithCancel(context.Background())
		defer cancelAttack()
		ticker := time.NewTicker(10 * time.Second)
		useFlyWing()
		monster := []string{"玩具士兵★"}
		stuckCount := time.Now()
	fightLoop:
		for {
			select {
			case <-time.After(5 * time.Second):
				if !g.IsMonsterInRange(monster...) {
					log.Warnf("找不到怪物")
					ticker.Stop()
					cancelAttack()
					break fightLoop
				}
				g.EnableAutoAttack(disableAttack, monster...)
			case <-ticker.C:
				if !g.IsMonsterInRange(monster...) {
					log.Warnf("找不到怪物")
					ticker.Stop()
					cancelAttack()
					break fightLoop
				}
				log.Infof("找到怪物%s", monster)
				if time.Since(stuckCount) > time.Second*45 {
					log.Warnf("卡住了")
					ticker.Stop()
					cancelAttack()
					break fightLoop
				}
			}
		}
	}
}

func getSockCount() uint32 {
	iData := g.FindPackItemByName("圣诞袜子", Cmd.EPackType_EPACKTYPE_MAIN)
	return iData.GetCount()
}
