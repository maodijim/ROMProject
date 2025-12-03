package AutoBossHunting

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

type HiddenMVP struct {
	Info                 utils.MonsterInfo
	Map                  gameTypes.MapId
	RespawnTime          time.Time
	PrerequisiteMonsters utils.MonsterInfo
	PrerequisitePosList  []Cmd.ScenePos
	BossPosList          []Cmd.ScenePos
}

type BossHuntTask struct {
	GC              *gameConnection.GameConnection
	ctx             context.Context
	cancel          context.CancelFunc
	logWriter       io.Writer
	logger          *log.Logger
	flyWingUseCount int
	fightStar       bool
	pickupCount     uint32
	lastPosUpdate   time.Time
	startTime       time.Time
	lastPos         Cmd.ScenePos
	flyMutex        sync.Mutex
	huntingCount    uint32
	targetMonster   Cmd.BossInfoItem
	mitionCompelete bool
	lastHp          int32
	tempMHP         int32
	tempUHP         int32
	huntHidMVP      bool
	targetHiddenMVP *HiddenMVP
	hiddenMVPList   map[string]*HiddenMVP
	fightCtx        context.Context
	fightCancel     context.CancelFunc
	workState       Work
	posCount        int32
	BossposCount    int32
	haveBoss        bool
	findBoss        bool
}

func (b *BossHuntTask) GetContext() context.Context {
	return b.ctx
}

type Work int

const (
	Init Work = iota
	TeleportMap
	MOVE_PrerequisiteMonstersPOS
	CHECK_PrerequisiteMonsters
	HUNT_PrerequisiteMonsters
	MOVE_BOSSPOS
	CHECK_BOSS
	HUNT_BOSS
	End
)

func i32(v int32) *int32 { return &v }

func (b *BossHuntTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(b.GC.LogWriter(), writer)
	b.logWriter = mw
	b.logger.SetOutput(mw)
}

func (b *BossHuntTask) Start() {
	b.GC.ShouldChangeScene = true
	b.GC.GameServerLogin()

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		// 等待登录完成
		oneTick := time.After(time.Millisecond * 50)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				//b.checkBossLive()
			case <-b.ctx.Done():
				b.logger.Info("停止自动BOSS狩猎任务")
				b.cancel()
				b.GC.Close()
				return
			case <-oneTick:
				go b.startHunt()
			}
		}

	}()
}

func (b *BossHuntTask) Stop() {
	b.cancel()
}

func (b *BossHuntTask) checkBossLive() bool {
	if b.GC.BossInfo != nil {
		targetHidMvpList := b.GC.Configs.HuntConfig.HMVP
		BossInfo := b.GC.BossInfo

		// 查找隐藏BOSS
		for _, v := range targetHidMvpList {
			if _, ok := b.hiddenMVPList[v]; ok {
				if time.Since(b.hiddenMVPList[v].RespawnTime) > time.Second*0 {
					b.huntHidMVP = true
					b.targetHiddenMVP = b.hiddenMVPList[v]
					b.logger.Info("%s已复活，进行狩猎", v)
					return true
				} else {
					remaining := time.Until(b.hiddenMVPList[v].RespawnTime)
					minutes := int(remaining.Minutes())
					b.logger.Infof("%s 目標時間尚未到，還有約 %d 分鐘\n", v, minutes)
				}
			}
		}
		// 查找BOSS清单
		for _, v := range BossInfo.Bosslist {
			if Contains(b.GC.Configs.HuntConfig.MVP, b.GC.MonsterItems[v.GetId()].NameZh) && v.GetMapid() != gameTypes.MapId_LabyrinthForest.Uint32() && v.GetSettime() == 0 {
				if v.RefreshTime == nil {
					b.logger.Infof(" %s 已复活，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
					b.targetMonster = *v
					b.mitionCompelete = false
					return true
				} else {
					past, diffMin := IsPastTime(v.GetRefreshTime())
					if past {
						b.logger.Infof(" %s 已复活，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
						b.targetMonster = *v
						b.mitionCompelete = false
						return true
					} else {
						b.logger.Infof("%s 目標時間尚未到，還有約 %d 分鐘\n", b.GC.MonsterItems[v.GetId()].NameZh, -diffMin)
					}
				}
			}
		}
		// 查找Mini清单
		for _, v := range BossInfo.Minilist {
			if Contains(b.GC.Configs.HuntConfig.Mini, b.GC.MonsterItems[v.GetId()].NameZh) && v.GetMapid() != gameTypes.MapId_LabyrinthForest.Uint32() {
				if v.RefreshTime == nil {
					b.logger.Infof(" %s 已复活，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
					b.targetMonster = *v
					b.mitionCompelete = false
					return true
				} else {
					past, diffMin := IsPastTime(v.GetRefreshTime())
					if past {
						b.logger.Infof(" %s 已复活，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
						b.targetMonster = *v
						b.mitionCompelete = false
						return true
					} else {
						b.logger.Infof("%s 目標時間尚未到，還有約 %d 分鐘\n", b.GC.MonsterItems[v.GetId()].NameZh, -diffMin)
					}
				}
			}
		}
	}

	return false
}

func (b *BossHuntTask) startHunt() {
	ctx, cancel := context.WithCancel(b.ctx)
	defer cancel()
	b.GC.GetBossInfo()
	b.checkBossLive()
	b.lastPos = b.GC.Role.GetPos()

	StartNum := b.GC.Role.GetLottery()

	ticker := time.NewTicker(time.Second * 10)
	ticker2 := time.NewTicker(time.Second * 1)

	targetId := uint64(0)
	// 隐藏MVP清单
	ConfigMVPName := b.GC.Configs.HuntConfig.HMVP
	b.hiddenMVPList = map[string]*HiddenMVP{}
	for _, v := range ConfigMVPName {
		if b.GC.GetMonsterIdByName(v) != 0 {
			HMVP := &HiddenMVP{}
			switch v {
			case "卡仑":
				HMVP.Info = b.GC.MonsterItemsByName[v]
				HMVP.RespawnTime = time.Now()
				HMVP.Map = gameTypes.MapId_GingerbreadCity
				HMVP.PrerequisiteMonsters = b.GC.MonsterItemsByName["疯兔"]
				HMVP.PrerequisitePosList = CrazyRabbitPos
				HMVP.BossPosList = CarlenPos
			case "狼外婆":
				HMVP.Info = b.GC.MonsterItemsByName[v]
				HMVP.RespawnTime = time.Now()
				HMVP.Map = gameTypes.MapId_MistyForest
				HMVP.PrerequisiteMonsters = b.GC.MonsterItemsByName["尖叫魔"]
				HMVP.PrerequisitePosList = ScreamingDemonPos
				HMVP.BossPosList = BigBadWolfPos
			}
			b.hiddenMVPList[v] = HMVP
		}
	}

	go func() {
		b.lastPosUpdate = time.Now()
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("停止自动BOSS狩猎任务监控协程")
				return
			default:
				if targetId != 0 && b.GC.AtkStat.GetCurrentTargetId() == targetId {
					TargetID := b.GC.AtkStat.GetCurrentTargetId()
					npcs := b.GC.GetMapNpcs()
					if _, ok := npcs[TargetID]; ok {
						HP := utils.GetNpcAttrValByType(b.GC.MapNpcs[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
						if b.lastHp == 0 || b.lastHp != HP {
							b.lastHp = HP
							b.lastPosUpdate = time.Now()
						} else if time.Since(b.lastPosUpdate) > time.Second*10 {
							b.logger.Infof("卡住了")
							b.GC.AtkStat.SetCurrentTargetId(0)
							b.useFlyWing()
							targetId = 0
						}
					}
				} else if b.fightStar && targetId == 0 && time.Since(b.lastPosUpdate) > time.Second*10 {
					b.logger.Infof("没有目标卡住了")
					b.fightStar = false
					b.useFlyWing()
					b.lastPosUpdate = time.Now()
				} else if b.GC.AtkStat.GetCurrentTargetId() != targetId {
					targetId = b.GC.AtkStat.GetCurrentTargetId()
					b.lastPosUpdate = time.Now()
				}

				time.Sleep(time.Second * 2)
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("停止自动BOSS狩猎任务监控协程2")
				return
			// 确认BOSS是否被他人狩猎
			case <-ticker.C:
				if !b.huntHidMVP {
					b.GC.GetBossInfo()
					if !b.mitionCompelete && !b.fightStar {
						if !b.checkTargetBossLive() {
							b.logger.Infof("%s 死亡，重新查找", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh)
							b.buyFlyWing()
							b.fightStar = false
							b.fightCancel()
							b.mitionCompelete = true
							b.checkBossLive()
							b.logger.Infof("已获取彩币数量:%d,已狩猎数量:%d", b.GC.Role.GetLottery()-StartNum, b.huntingCount)
						} else {
							b.logger.Infof("%s 未死亡，继续寻找", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh)
						}
					}
				}
			// 查找BOSS清单
			case <-ticker2.C:
				if b.mitionCompelete && !b.huntHidMVP {
					if !b.checkBossLive() {
						b.logger.Infof("没有找到目标，躺平吧!")
						time.Sleep(time.Millisecond * 10000)
					} else {
						b.CheckCloseTime()
						b.logger.Infof("已获取彩币数量:%d,已狩猎数量:%d", b.GC.Role.GetLottery()-StartNum, b.huntingCount)
						if !b.huntHidMVP {
							b.mitionCompelete = false
						}
					}
				}
			}

		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("停止自动BOSS狩猎任务主协程")
				return
			default:
				if !b.mitionCompelete && b.targetMonster.GetMapid() != 0 {
					if b.GC.Role.GetMapId() != b.targetMonster.GetMapid() {
						if b.GC.Configs.HuntConfig.CarryTeam {
							b.GC.TeamGoToMap(b.targetMonster.GetMapid())
						} else {
							b.GC.GoToMap(b.targetMonster.GetMapid())
						}

						b.GC.ChangeMap(b.targetMonster.GetMapid())
						time.Sleep(time.Millisecond * 1000)
						b.useFlyWing()
						time.Sleep(time.Millisecond * 3200)
						b.useSkill()
						time.Sleep(time.Millisecond * 1000)
						b.useSkill()
						time.Sleep(time.Millisecond * 1000)
					} else if b.GC.IsMonsterInRange(b.GC.MonsterItems[b.targetMonster.GetId()].NameZh) && !b.fightStar {
						b.logger.Infof("找到%s", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh)
						b.fightMonstStar(b.GC.MonsterItems[b.targetMonster.GetId()].NameZh, b.targetMonster.GetId())
						b.flyWingUseCount = 0
						b.pickupCount = 0
						b.tempMHP = 0
						b.tempUHP = 0
					} else if !b.GC.IsMonsterInRange(b.GC.MonsterItems[b.targetMonster.GetId()].NameZh) {
						b.logger.Infof("找不到%s 使用翅膀", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh)
						b.useFlyWing()
						b.flyWingUseCount++
						b.fightStar = false
						b.fightCancel()
					} else if b.fightStar {
						if !b.checkTargetBossLive() {
							b.logger.Infof("%s 死亡，重新查找", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh)
							b.huntingCount++
							b.fightStar = false
							b.fightCancel()
							b.mitionCompelete = true
							time.Sleep(time.Second * 3)
						} else {
							TargetID := b.GC.AtkStat.GetCurrentTargetId()
							MapNPC := b.GC.GetMapNpcs()
							if _, ok := MapNPC[TargetID]; ok {
								MonsterHP := utils.GetNpcAttrValByType(MapNPC[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
								MHP := utils.GetNpcAttrValByType(b.GC.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
								if MonsterHP != b.tempMHP || MHP != b.tempUHP {
									b.tempMHP = MonsterHP
									b.tempUHP = MHP
									b.logger.Infof("%s 未死亡，剩余血量:%d", b.GC.MonsterItems[b.targetMonster.GetId()].NameZh, MonsterHP)
									b.logger.Infof("我的血量:%d", MHP)
								} else if MonsterHP == 0 {
									b.GC.GetBossInfo()
								}
							}
						}
					}

					time.Sleep(time.Millisecond * 1000)

				} else if b.huntHidMVP {
					b.huntCarlen()
					b.huntHidMVP = false
				} else {
					time.Sleep(time.Millisecond * 1000)
				}
			}

		}
	}()
	<-ctx.Done()
}

// fightMonstStar 开始自动打怪, MonsterName:怪物名称, monsterId:怪物ID(用于避免同名怪物)
func (b *BossHuntTask) fightMonstStar(MonsterName string, monsterId uint32) {
	var m utils.MonsterInfo
	if monsterId > 0 {
		m = b.GC.GetMonsterItemById(monsterId)
	} else {
		m = b.GC.GetMonsterItemByName(MonsterName)
	}
	nature := m.Nature
	if nature != "" {
		if b.GC.Role.GetProfession() >= Cmd.EProfession_EPROFESSION_ARCHER && b.GC.Role.GetProfession() <= Cmd.EProfession_EPROFESSION_RANGER {
			if nature == gameTypes.NatureType_Fire {
				b.useElementArrow(gameTypes.WaterArrow)
			} else if nature == gameTypes.NatureType_Water {
				b.useElementArrow(gameTypes.WindArrow)
			} else if nature == gameTypes.NatureType_Wind {
				b.useElementArrow(gameTypes.EarthArrow)
			} else if nature == gameTypes.NatureType_Earth {
				b.useElementArrow(gameTypes.FireArrow)
			} else if nature == gameTypes.NatureType_Undead || nature == gameTypes.NatureType_Shawdow {
				b.useElementArrow(gameTypes.SliverArrow)
			} else {
				b.useElementArrow(gameTypes.FireArrow)
			}
		}
	}

	b.fightCancel()
	b.fightStar = true
	b.fightCtx, b.fightCancel = context.WithCancel(context.Background())
	b.GC.EnableAutoAttack(b.fightCtx, MonsterName)
	b.lastPosUpdate = time.Now()
}

func (b *BossHuntTask) useFlyWing() {
	b.flyMutex.Lock()
	defer b.flyMutex.Unlock()
	b.buyFlyWing()
	b.GC.UseFlyWing()
	item := b.GC.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil && item.GetBase().GetCount() > 0 {
		b.logger.Infof("使用苍蝇翅膀 还有%d个", item.GetBase().GetCount())
	} else {
		b.logger.Warn("没有找到苍蝇翅膀")
		_ = b.GC.GetMainPackItems()
	}
}

func (b *BossHuntTask) buyFlyWing() {
	if item := b.GC.FindPackItemByName("苍蝇翅膀", Cmd.EPackType_EPACKTYPE_MAIN); item == nil || item.GetBase().GetCount() > 1000 {
		return
	}
	shopConfig, err := b.GC.QueryShopConfig(gameTypes.ShopType_Item, 1)
	if err != nil {
		b.logger.Errorf("查询商店配置失败 %s", err)
		return
	}
	for _, item := range shopConfig.GetGoods() {
		if item.GetItemid() == 5024 {
			b.logger.Infof("购买1000苍蝇翅膀")
			b.GC.BuyShopItem(item, 1000)
		}
	}
}

func (b *BossHuntTask) checkTargetBossLive() bool {
	if b.targetMonster.GetId() == 0 {
		return false
	}
	BossInfo := *b.GC.BossInfo
	// 确认MVP复活时间
	for _, v := range BossInfo.Bosslist {
		if v.GetId() == b.targetMonster.GetId() && v.GetMapid() == b.targetMonster.GetMapid() {
			if v.RefreshTime == nil {
				b.logger.Debugf(" %s 还活着，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
				return true
			} else {
				past, diffMin := IsPastTime(v.GetRefreshTime())
				if past {
					b.logger.Debugf(" %s 还活着，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
					return true
				} else {
					b.logger.Debugf("%s 已死亡，還有約 %d 分鐘", b.GC.MonsterItems[v.GetId()].NameZh, -diffMin)
					return false
				}
			}
		}
	}
	// 确认Mini复活时间
	for _, v := range BossInfo.Minilist {
		if v.GetId() == b.targetMonster.GetId() && v.GetMapid() == b.targetMonster.GetMapid() {
			if v.RefreshTime == nil {
				b.logger.Debugf(" %s 还活着，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
				return true
			} else {
				past, diffMin := IsPastTime(v.GetRefreshTime())
				if past {
					b.logger.Debugf(" %s 还活着，进行狩猎", b.GC.MonsterItems[v.GetId()].NameZh)
					return true
				} else {
					b.logger.Debugf("%s 已死亡，還有約 %d 分鐘", b.GC.MonsterItems[v.GetId()].NameZh, -diffMin)
					return false
				}
			}
		}
	}
	return false
}

func (b *BossHuntTask) useSkill() {
	b.logger.Infof("使用装死!")
	num := int32(1)
	dir := int32(utils.GetNpcDataValByType(b.GC.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    b.GC.Role.Pos,
		Dir:    &dir,
	}
	b.GC.SkillCmd(10020001, pData, true)
}

func (b *BossHuntTask) CheckCloseTime() {
	GameDuration := b.GC.Configs.HuntConfig.GameDuration
	if GameDuration > 0 {
		remaining := time.Until(b.startTime) + time.Duration(GameDuration)*time.Hour
		hour := int(remaining.Hours())
		minutes := int(remaining.Minutes()) % 60
		seconds := int(remaining.Seconds()) % 3600
		if GameDuration > 0 && hour <= 0 && minutes <= 0 && seconds <= 0 {
			b.logger.Info("游戏时间到，准备关闭游戏")
			b.Stop()
			return
		} else {
			b.logger.Infof("距离关闭还有%d小时，%d分钟", hour, minutes)
		}

	} else {
		remaining := time.Until(b.startTime)
		hour := -int(remaining.Hours())
		minutes := -int(remaining.Minutes()) % 60
		b.logger.Infof("已挂机%d小时，%d分钟", hour, minutes)
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

func NewBossHuntTask(ctx context.Context, gc *gameConnection.GameConnection) *BossHuntTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	fiightCtx, fightCancel := context.WithCancel(newCtx)
	return &BossHuntTask{
		GC:              gc,
		ctx:             newCtx,
		cancel:          cancel,
		logWriter:       mw,
		logger:          logger,
		startTime:       time.Now(),
		lastPosUpdate:   time.Now(),
		hiddenMVPList:   make(map[string]*HiddenMVP),
		flyMutex:        sync.Mutex{},
		lastPos:         Cmd.ScenePos{},
		targetMonster:   Cmd.BossInfoItem{},
		fightCancel:     fightCancel,
		fightCtx:        fiightCtx,
		workState:       Init,
		targetHiddenMVP: &HiddenMVP{},
	}
}

func (b *BossHuntTask) transition(SwitchState Work) {
	if b.workState != SwitchState {
		b.workState = SwitchState
	}
}
