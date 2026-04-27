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
	ItemCount             []uint32
	lastPosUpdate         = time.Now()
	lastPos               Cmd.ScenePos
	flyMutex              *sync.Mutex
	flyMutexWait          = float64(7000)
)

func i32(v int32) *int32 { return &v }

var Pos_03 = []Cmd.ScenePos{
	{X: i32(21948), Y: i32(-583), Z: i32(43399)},
	{X: i32(-14088), Y: i32(73), Z: i32(-56505)},
}

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

	for i := 0; i < len(g.Configs.HuntConfig.TargetItems); i++ {
		curCount := getItemCount(g.Configs.HuntConfig.TargetItems[i])
		log.Infof("当前%s数量 %d", g.Configs.HuntConfig.TargetItems[i])
		ItemCount = append(ItemCount, curCount)
	}

	MapID := gameTypes.MapNameZh[g.Configs.HuntConfig.Map].Uint32()

	ticker := time.NewTicker(time.Second * 5)
	ticker2 := time.NewTicker(time.Second * 30)

	targetId := uint64(0)
	go func() {
		for {
			select {
			case <-ticker.C:
				if targetId != 0 && g.AtkStat.GetCurrentTargetId() == targetId && time.Since(lastPosUpdate) > time.Second*3 {
					log.Infof("卡住了")
					g.AtkStat.SetCurrentTargetId(0)
					targetId = 0
					useFlyWing()
				} else if g.AtkStat.GetCurrentTargetId() != targetId {
					targetId = g.AtkStat.GetCurrentTargetId()
					lastPosUpdate = time.Now()
				}
			case <-ticker2.C:
				for i := 0; i < len(g.Configs.HuntConfig.TargetItems); i++ {
					curCount := getItemCount(g.Configs.HuntConfig.TargetItems[i])
					log.Infof("当前%s数量 %d, 打了 %d", g.Configs.HuntConfig.TargetItems[i], curCount, curCount-ItemCount[i])
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	for {
		if g.Role.GetMapId() != MapID {
			if MapID == gameTypes.MapId_LhzDun03.Uint32() {

				g.GoToMap(gameTypes.MapId_LhzDun01.Uint32())

				g.ChangeMap(gameTypes.MapId_LhzDun01.Uint32())

				for i := 0; i < len(g.Configs.HuntConfig.TargetItems); i++ {
					g.MoveChartWait(Pos_03[i])
				}
			} else {
				g.GoToMap(MapID)

				g.ChangeMap(MapID)
			}

			time.Sleep(time.Millisecond * 1000)
			useFlyWing()
			time.Sleep(time.Millisecond * 3200)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
		} else {
			hasTarget := g.IsMonsterInRange(g.Configs.HuntConfig.TargetMonsters...)
			if !hasTarget {
				log.Infof("没有找到目标怪物")
				useFlyWing()
			} else if !fighting {
				log.Infof("附近找到目标怪物，开始自动挂机，坐稳了")
				fightCtx, fightCancel = context.WithCancel(context.Background())
				g.EnableAutoAttack(fightCtx, g.Configs.HuntConfig.TargetMonsters...)
				fighting = true
			}
		}
		time.Sleep(time.Millisecond * 1000)
	}
}

func getItemCount(ItemName string) uint32 {
	iData := g.FindPackItemByName(ItemName, Cmd.EPackType_EPACKTYPE_MAIN)
	if iData == nil {
		return 0
	}
	return iData.GetBase().GetCount()
}

func useFlyWing() {
	flyMutex.Lock()
	defer flyMutex.Unlock()
	buyFlyWing()
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
			log.Infof("购买10000苍蝇翅膀")
			g.BuyShopItem(item, 10000)
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

func Useskill() {
	log.Infof("使用装死!")
	num := int32(1)
	dir := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    g.Role.Pos,
		Dir:    &dir,
	}
	g.SkillCmd(10020001, pData, true)
}
