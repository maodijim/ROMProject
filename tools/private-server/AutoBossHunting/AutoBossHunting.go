package main

import (
	"context"
	"flag"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
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
	flyMutex = &sync.Mutex{}

	// 隐藏MVP清单
	ConfigMVPName := g.Configs.HuntConfig.HMVP
	HiddenMVPList = map[string]HiddenMVP{}
	for _, v := range ConfigMVPName {
		if g.GetMonsterIdByName(v) != 0 {
			HMVP := HiddenMVP{
				Info:        g.MonsterItemsByName[v],
				RespawnTime: time.Now(),
			}
			HiddenMVPList[v] = HMVP
		}
	}
	start()
}

func start() {

	g.ShouldChangeScene = true
	g.GameServerLogin()

	g.GetBossInfo()
	checkBossLive()

	StartNum := g.Role.GetLottery()

	ticker := time.NewTicker(time.Second * 10)
	ticker2 := time.NewTicker(time.Second * 1)

	targetId := uint64(0)
	go func() {
		lastPosUpdate = time.Now()
		for {
			if !HuntHidMVP {
				if targetId != 0 && g.AtkStat.GetCurrentTargetId() == targetId && g.IsMonsterInRange(g.MonsterItems[*TargetMonster.Id].NameZh) {
					TargetID := g.AtkStat.GetCurrentTargetId()
					if g.MapNpcs[TargetID].Attrs != nil {
						HP := utils.GetNpcAttrValByType(g.MapNpcs[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
						if LastHp == 0 || LastHp > HP {
							LastHp = HP
							lastPosUpdate = time.Now()
						} else if time.Since(lastPosUpdate) > time.Second*10 {
							log.Infof("卡住了")
							g.AtkStat.SetCurrentTargetId(0)
							useFlyWing()
							targetId = 0
						}
					}
				} else if fightStar && targetId == 0 && time.Since(lastPosUpdate) > time.Second*10 {
					log.Infof("没有目标卡住了")
					fightStar = false
					useFlyWing()
					lastPosUpdate = time.Now()
				} else if g.AtkStat.GetCurrentTargetId() != targetId {
					targetId = g.AtkStat.GetCurrentTargetId()
					lastPosUpdate = time.Now()
				}
			}
			time.Sleep(time.Second * 2)
		}
	}()

	go func() {
		for {
			select {
			// 确认BOSS是否被他人狩猎
			case <-ticker.C:
				if !HuntHidMVP {
					g.GetBossInfo()
					if !MitionCompelete && !fightStar {
						if !checkTargetBossLive() {
							log.Infof("%s 死亡，重新查找", g.MonsterItems[*TargetMonster.Id].NameZh)
							buyFlyWing()
							fightStar = false
							fightCancel()
							MitionCompelete = true
							checkBossLive()
							log.Infof("已获取彩币数量:%d,已狩猎数量:%d", g.Role.GetLottery()-StartNum, HuntingCount)
						} else {
							log.Infof("%s 未死亡，继续寻找", g.MonsterItems[*TargetMonster.Id].NameZh)
						}
					}
				}
			// 查找BOSS清单
			case <-ticker2.C:
				if MitionCompelete && !HuntHidMVP {
					if !checkBossLive() {
						log.Infof("没有找到目标，躺平吧!")
						time.Sleep(time.Millisecond * 10000)
					} else {
						log.Infof("已获取彩币数量:%d,已狩猎数量:%d", g.Role.GetLottery()-StartNum, HuntingCount)
						MitionCompelete = false
					}
				}
			}
		}
	}()

	for {
		if !MitionCompelete {
			if g.Role.GetMapId() != *TargetMonster.Mapid {
				if g.Configs.HuntConfig.CarryTeam {
					g.TeamGoToMap(*TargetMonster.Mapid)
				} else {
					g.GoToMap(*TargetMonster.Mapid)
				}

				g.ChangeMap(*TargetMonster.Mapid)
				time.Sleep(time.Millisecond * 1000)
				useFlyWing()
				time.Sleep(time.Millisecond * 3200)
				Useskill()
				time.Sleep(time.Millisecond * 1000)
				Useskill()
				time.Sleep(time.Millisecond * 1000)
			} else if g.IsMonsterInRange(g.MonsterItems[*TargetMonster.Id].NameZh) && !fightStar {
				log.Infof("找到%s", g.MonsterItems[*TargetMonster.Id].NameZh)
				fightMonstStar(g.MonsterItems[*TargetMonster.Id].NameZh, *TargetMonster.Id)
				flyWingUseCount = 0
				pickupCount = 0
				tempMHP = 0
				tempUHP = 0
				/*} else if flyWingUseCount > maxFlyWingUseCount || pickupCount > maxPickupCount {
				if flyWingUseCount > maxFlyWingUseCount {
					log.Infof("%d个翅膀找不到%s 放弃", maxFlyWingUseCount, g.MonsterItems[*TargetMonster.Id].NameZh)
				} else if pickupCount > maxPickupCount {
					log.Infof("拾取了%d个物品 放弃", maxPickupCount)
				}
				flyWingUseCount = 0
				pickupCount = 0*/
			} else if !g.IsMonsterInRange(g.MonsterItems[*TargetMonster.Id].NameZh) {
				log.Infof("找不到%s 使用翅膀", g.MonsterItems[*TargetMonster.Id].NameZh)
				useFlyWing()
				flyWingUseCount++
				fightStar = false
				fightCancel()
			} else if fightStar {
				if !checkTargetBossLive() {
					log.Infof("%s 死亡，重新查找", g.MonsterItems[*TargetMonster.Id].NameZh)
					HuntingCount++
					buyFlyWing()
					fightStar = false
					fightCancel()
					MitionCompelete = true
					time.Sleep(time.Second * 3)
				} else {
					TargetID := g.AtkStat.GetCurrentTargetId()
					MapNPC := g.GetMapNpcs()
					if g.IsMonsterInRange(g.MonsterItems[*TargetMonster.Id].NameZh) && MapNPC[TargetID].Attrs != nil {
						MonsterHP := utils.GetNpcAttrValByType(MapNPC[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
						MHP := utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
						if MonsterHP != tempMHP || MHP != tempUHP {
							tempMHP = MonsterHP
							tempUHP = MHP
							log.Infof("%s 未死亡，剩余血量:%d", g.MonsterItems[*TargetMonster.Id].NameZh, MonsterHP)
							log.Infof("我的血量", MHP)
						} else if MonsterHP == 0 {
							g.GetBossInfo()
						}
					}
				}
			}

			time.Sleep(time.Millisecond * 1000)

		} else if HuntHidMVP {
			switch TargetHiddenMVP.Info.NameZh {
			case "卡仑":
				HuntCarlen()
				HuntHidMVP = false
				break
			}
		} else {
			time.Sleep(time.Millisecond * 1000)
		}
	}
}

// fightMonstStar 开始自动打怪, MonsterName:怪物名称, monsterId:怪物ID(用于避免同名怪物)
func fightMonstStar(MonsterName string, monsterId uint32) {
	var m utils.MonsterInfo
	if monsterId > 0 {
		m = g.GetMonsterItemById(monsterId)
	} else {
		m = g.GetMonsterItemByName(MonsterName)
	}
	nature := m.Nature
	if nature != "" {
		if g.Role.GetProfession() >= 41 && g.Role.GetProfession() <= 44 {
			if nature == gameTypes.NatureType_Fire {
				useElementArrow(gameTypes.WaterArrow)
			} else if nature == gameTypes.NatureType_Water {
				useElementArrow(gameTypes.WindArrow)
			} else if nature == gameTypes.NatureType_Wind {
				useElementArrow(gameTypes.EarthArrow)
			} else if nature == gameTypes.NatureType_Earth {
				useElementArrow(gameTypes.FireArrow)
			} else if nature == gameTypes.NatureType_Undead || nature == gameTypes.NatureType_Shawdow {
				useElementArrow(gameTypes.SliverArrow)
			} else {
				useElementArrow(gameTypes.FireArrow)
			}
		}
	}

	fightCancel()
	fightStar = true
	fightCtx, fightCancel = context.WithCancel(context.Background())
	g.EnableAutoAttack(fightCtx, MonsterName)
	lastPosUpdate = time.Now()
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
			log.Infof("购买1000苍蝇翅膀")
			g.BuyShopItem(item, 1000)
		}
	}
}

func Contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

func checkBossLive() bool {
	if g.BossInfo != nil {
		HiddenMVP := g.Configs.HuntConfig.HMVP
		BossInfo := *g.BossInfo

		// 查找隐藏BOSS
		for _, v := range HiddenMVP {
			if _, ok := HiddenMVPList[v]; ok {
				if time.Since(HiddenMVPList[v].RespawnTime) > time.Second*0 {
					HuntHidMVP = true
					TargetHiddenMVP = HiddenMVPList[v]
					log.Printf("卡仑已复活，进行狩猎")
					return false
				} else {
					remaining := time.Until(HiddenMVPList[v].RespawnTime)
					minutes := int(remaining.Minutes())
					log.Printf("%s 目標時間尚未到，還有約 %d 分鐘\n", v, minutes)
				}
			}
		}
		// 查找BOSS清单
		for _, v := range BossInfo.Bosslist {
			if Contains(g.Configs.HuntConfig.MVP, g.MonsterItems[*v.Id].NameZh) && *v.Mapid != gameTypes.MapId_LabyrinthForest.Uint32() && *v.Settime == 0 {
				if v.RefreshTime == nil {
					log.Printf(" %s 已复活，进行狩猎", g.MonsterItems[*v.Id].NameZh)
					TargetMonster = *v
					MitionCompelete = false
					return true
				} else {
					past, diffMin := IsPastTime(*v.RefreshTime)
					if past {
						log.Printf(" %s 已复活，进行狩猎", g.MonsterItems[*v.Id].NameZh)
						TargetMonster = *v
						MitionCompelete = false
						return true
					} else {
						log.Printf("%s 目標時間尚未到，還有約 %d 分鐘\n", g.MonsterItems[*v.Id].NameZh, -diffMin)
					}
				}
			}
		}
		// 查找Mini清单
		for _, v := range BossInfo.Minilist {
			if Contains(g.Configs.HuntConfig.Mini, g.MonsterItems[*v.Id].NameZh) && *v.Mapid != gameTypes.MapId_LabyrinthForest.Uint32() {
				if v.RefreshTime == nil {
					log.Printf(" %s 已复活，进行狩猎", g.MonsterItems[*v.Id].NameZh)
					TargetMonster = *v
					MitionCompelete = false
					return true
				} else {
					past, diffMin := IsPastTime(*v.RefreshTime)
					if past {
						log.Printf(" %s 已复活，进行狩猎", g.MonsterItems[*v.Id].NameZh)
						TargetMonster = *v
						MitionCompelete = false
						return true
					} else {
						log.Printf("%s 目標時間尚未到，還有約 %d 分鐘\n", g.MonsterItems[*v.Id].NameZh, -diffMin)
					}
				}
			}
		}
	}

	return false
}

func checkTargetBossLive() bool {
	if *TargetMonster.Id == 0 {
		return false
	}
	BossInfo := *g.BossInfo
	// 确认MVP复活时间
	for _, v := range BossInfo.Bosslist {
		if *v.Id == *TargetMonster.Id && *v.Mapid == *TargetMonster.Mapid {
			if v.RefreshTime == nil {
				log.Debugf(" %s 还活着，进行狩猎", g.MonsterItems[*v.Id].NameZh)
				return true
			} else {
				past, diffMin := IsPastTime(*v.RefreshTime)
				if past {
					log.Debugf(" %s 还活着，进行狩猎", g.MonsterItems[*v.Id].NameZh)
					return true
				} else {
					log.Debugf("%s 已死亡，還有約 %d 分鐘", g.MonsterItems[*v.Id].NameZh, -diffMin)
					return false
				}
			}
		}
	}
	// 确认Mini复活时间
	for _, v := range BossInfo.Minilist {
		if *v.Id == *TargetMonster.Id && *v.Mapid == *TargetMonster.Mapid {
			if v.RefreshTime == nil {
				log.Debugf(" %s 还活着，进行狩猎", g.MonsterItems[*v.Id].NameZh)
				return true
			} else {
				past, diffMin := IsPastTime(*v.RefreshTime)
				if past {
					log.Debugf(" %s 还活着，进行狩猎", g.MonsterItems[*v.Id].NameZh)
					return true
				} else {
					log.Debugf("%s 已死亡，還有約 %d 分鐘", g.MonsterItems[*v.Id].NameZh, -diffMin)
					return false
				}
			}
		}
	}
	return false
}

func IsPastTime(timestamp uint32) (bool, int64) {
	// 台灣時區 (UTC+8)
	taiwan := time.FixedZone("CST", 8*3600)

	// 將 timestamp 轉成台灣時間
	targetTime := time.Unix(int64(timestamp), 0).In(taiwan)

	// 取得目前台灣時間
	now := time.Now().In(taiwan)

	// 計算時間差（分鐘）
	diff := now.Sub(targetTime).Minutes()

	// 若目標時間早於現在，回傳 true 以及差值
	return targetTime.Before(now), int64(diff)
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
