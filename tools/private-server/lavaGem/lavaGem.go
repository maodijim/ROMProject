package main

import (
	"context"
	"flag"
	"fmt"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

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
	flyMutexWait          = float64(6500)
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

func main() {
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
	flyMutex = &sync.Mutex{}
	start()
}

func start() {
	g.GameServerLogin()
	g.ChangeMap(g.Role.GetMapId())

	if g.Role.GetMapId() != gameTypes.MapId_MagmaDungeon1F.Uint32() {
		g.ExitMapWait(gameTypes.MapId_MagmaDungeon1F.Uint32())
		g.ChangeMap(gameTypes.MapId_MagmaDungeon1F.Uint32())
	}

	spawnNotification := make(chan bool)
	getNotification(spawnNotification)
	ticker := time.NewTicker(time.Second * 10)
	lavaGemCount = getLavaGemCount()
	targetId := uint64(0)
	go func() {
		for {
			if targetId != 0 && g.AtkStat.GetCurrentTargetId() == targetId && time.Since(lastPosUpdate) > time.Second*30 {
				log.Infof("卡住了")
				g.AtkStat.SetCurrentTargetId(0)
				targetId = 0
				useFlyWing()
			} else if targetId == 0 && time.Since(lastPosUpdate) > time.Second*60 {
				log.Infof("没有目标卡住了")
				useFlyWing()
				lastPosUpdate = time.Now()
			} else if g.AtkStat.GetCurrentTargetId() != targetId {
				targetId = g.AtkStat.GetCurrentTargetId()
				lastPosUpdate = time.Now()
			}
			time.Sleep(time.Second * 5)
		}
	}()
	for {
		select {
		case <-ticker.C:
			hasTarget := g.IsMonsterInRange("爆炎小恶魔 ", "爆炎小恶魔★")
			if fightStar {
				continue
			} else if !hasTarget {
				log.Infof("没有找到目标爆炎小恶魔")
				useFlyWing()
			} else if !fighting {
				log.Infof("开始打爆炎小恶魔")
				fightCtx, fightCancel = context.WithCancel(context.Background())
				fightDevil(fightCtx)
			}
			curCount := getLavaGemCount()
			log.Infof("当前熔岩宝石数量 %d, 打了 %d", curCount, curCount-lavaGemCount)
		case <-spawnNotification:
			log.Infof("收到通知爆炎小恶魔★ 出现了")
			pickupCount = 0
			fighting = false
			fightCancel()
		flyLoop:
			for {
				if g.IsMonsterInRange("爆炎小恶魔★") && !fightStar {
					log.Infof("找到爆炎小恶魔★")
					fightDevilStar()
				} else if flyWingUseCount > maxFlyWingUseCount || pickupCount > maxPickupCount {
					if flyWingUseCount > maxFlyWingUseCount {
						log.Infof("%d个翅膀找不到爆炎小恶魔★ 放弃", maxFlyWingUseCount)
					} else if pickupCount > maxPickupCount {
						log.Infof("拾取了%d个物品 放弃", maxPickupCount)
					}
					flyWingUseCount = 0
					pickupCount = 0
					break flyLoop
				} else if !g.IsMonsterInRange("爆炎小恶魔★") {
					log.Infof("找不到爆炎小恶魔★ 使用翅膀")
					useFlyWing()
					flyWingUseCount++
				}
				time.Sleep(time.Millisecond * 1000)
			}
			buyFlyWing()
			fightStar = false
			fightCancel()
		}
	}
}

func fightDevil(ctx context.Context) {
	fighting = true
	g.EnableAutoAttack(ctx, "爆炎小恶魔 ")
}

func fightDevilStar() {
	fightStar = true
	fightCtx, fightCancel = context.WithCancel(context.Background())
	g.EnableAutoAttack(fightCtx, "爆炎小恶魔★")
}

func getLavaGemCount() uint32 {
	iData := g.FindPackItemByName("熔岩宝石", Cmd.EPackType_EPACKTYPE_MAIN)
	if iData == nil {
		return 0
	}
	return iData.GetBase().GetCount()
}

func getNotification(spawnNotification chan bool) {
	g.AddNotifier(gameTypes.NtfType_UserActionDialog)
	g.AddNotifier(gameTypes.NtfType_UserItemPickup)
	dialogNtf := g.Notifier(gameTypes.NtfType_UserActionDialog)
	pickupNtf := g.Notifier(gameTypes.NtfType_UserItemPickup)
	go func() {
		for {
			select {
			case item := <-pickupNtf:
				name := g.FindItemNameById(item.(*Cmd.MapItem).GetId())
				log.Infof("拾取了 %s", name)
				if name == "熔岩宝石" {
					flyMutex.Lock()
					log.Infof("拾取熔岩宝石中停%f秒", flyMutexWait/1000)
					time.Sleep(time.Millisecond * time.Duration(flyMutexWait))
					flyMutex.Unlock()
				}
				pickupCount++
			case newNtf := <-dialogNtf:
				ntf, ok := newNtf.(*Cmd.UserActionNtf)
				if ok && ntf.GetValue() == 65073 {
					spawnNotification <- true
				}
			}
		}
	}()
}

func useFlyWing() {
	flyMutex.Lock()
	defer flyMutex.Unlock()
	g.UseFlyWing()
	item := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil && item.GetBase().GetCount() > 0 {
		log.Infof("使用苍蝇翅膀 还有%d个", item.GetBase().GetCount())
	} else {
		log.Warn("没有找到苍蝇翅膀")
		_ = g.GetMainPackItems()
	}
}

func buyFlyWing() {
	if item := g.FindPackItemByName("苍蝇翅膀", Cmd.EPackType_EPACKTYPE_MAIN); item == nil || item.GetBase().GetCount() > 1000 {
		return
	}
	shopConfig, err := g.QueryShopConfig(gameTypes.ShopType_Item, 1)
	if err != nil {
		log.Errorf("查询商店配置失败 %s", err)
		return
	}
	for _, item := range shopConfig.GetGoods() {
		if item.GetItemid() == 5024 {
			log.Infof("购买50苍蝇翅膀")
			g.BuyShopItem(item, 50)
		}
	}
}

func printNearbyNpcs(stopNpc context.Context) {
	log.Printf("Nearby NPCs:")
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopNpc.Done():
				return
			case <-ticker.C:
				npcList := map[string][]int{}
				for _, npc := range g.GetMapNpcs() {
					if _, ok := npcList[npc.GetName()]; ok {
						npcList[npc.GetName()][0] += 1
					} else {
						npcList[npc.GetName()] = []int{1, int(*npc.Id)}
					}
				}
				output := "\n"
				for k, v := range npcList {
					output += fmt.Sprintf("名字：%s，数量%d\n", k, v[0])
				}
				log.Printf("NPC: %s", output)
			}
		}
	}()
}
