package AutoHunting

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

type HuntTask struct {
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
	fightCtx        context.Context
	fightCancel     context.CancelFunc
}

var (
	ItemCount     []uint32
	lastPosUpdate time.Time
	lastFlyTime   time.Time
)

func i32(v int32) *int32 { return &v }

var Pos_03 = []Cmd.ScenePos{
	{X: i32(21948), Y: i32(-583), Z: i32(43399)},
	{X: i32(-14088), Y: i32(357), Z: i32(-56505)},
	// 2f东区 -> 2f 西区 入口
	{X: i32(46754), Y: i32(357), Z: i32(857)},
}

func (b *HuntTask) Start() {
	b.GC.ShouldChangeScene = true
	b.GC.GameServerLogin()

	go func() {
		ticker := time.NewTicker(time.Second * 10)
		// 等待登录完成
		oneTick := time.After(time.Millisecond * 50)
		defer ticker.Stop()
		for {
			select {
			case <-b.ctx.Done():
				b.logger.Info("停止自动BOSS狩猎任务")
				b.cancel()
				b.GC.Close()
				return
			case <-oneTick:
				go b.StartHunt()
			}
		}

	}()
}

func (b *HuntTask) StartHunt() {
	ctx, cancel := context.WithCancel(b.ctx)
	defer cancel()

	for i := 0; i < len(b.GC.Configs.HuntConfig.TargetItems); i++ {
		curCount := b.getItemCount(b.GC.Configs.HuntConfig.TargetItems[i])
		b.logger.Infof("当前%s数量 %d", b.GC.Configs.HuntConfig.TargetItems[i])
		ItemCount = append(ItemCount, curCount)
	}

	b.GC.GetAllPackItems()

	MapID := gameTypes.MapNameZh[b.GC.Configs.HuntConfig.Map].Uint32()

	if b.GC.Role.GetMapId() == MapID {
		b.EnableGodMode()
	}

	ticker := time.NewTicker(time.Second * 1)
	ticker2 := time.NewTicker(time.Second * 30)

	targetId := uint64(0)
	go func() {
		for {
			select {
			case <-ctx.Done():
				b.logger.Info("停止自动狩猎任务监控协程1")
				return
			case <-ticker.C:
				if targetId != 0 && b.GC.AtkStat.GetCurrentTargetId() == targetId && time.Since(lastPosUpdate) > time.Second*10 {
					b.logger.Infof("卡住了")
					b.GC.AtkStat.SetCurrentTargetId(0)
					targetId = 0
					b.useFlyWing()
				} else if b.fightStar && targetId == 0 && time.Since(lastPosUpdate) > time.Second*10 {
					b.logger.Infof("没有目标卡住了")
					b.useFlyWing()
					lastPosUpdate = time.Now()
				} else if b.GC.AtkStat.GetCurrentTargetId() != targetId {
					targetId = b.GC.AtkStat.GetCurrentTargetId()
					lastPosUpdate = time.Now()
				}
			case <-ticker2.C:
				for i := 0; i < len(b.GC.Configs.HuntConfig.TargetItems); i++ {
					curCount := b.getItemCount(b.GC.Configs.HuntConfig.TargetItems[i])
					b.logger.Infof("当前%s数量 %d, 打了 %d", b.GC.Configs.HuntConfig.TargetItems[i], curCount, curCount-ItemCount[i])
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			b.logger.Infof("主循环结束")
			return
		default:
		}

		if b.GC.Role.GetMapId() != MapID {
			if MapID == gameTypes.MapId_LhzDun03.Uint32() {
				b.GC.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
				time.Sleep(time.Millisecond * 500)
				b.GC.MoveChartWait(Pos_03[0])
				time.Sleep(time.Millisecond * 500)
				b.GC.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, b.GC.Role.GetPos())
				time.Sleep(time.Millisecond * 500)
				b.GC.MoveChartWait(Pos_03[1])
				time.Sleep(time.Millisecond * 500)
				b.GC.ExitMapPos(gameTypes.MapId_LhzDun02.Uint32(), 3, b.GC.Role.GetPos())
				time.Sleep(time.Millisecond * 500)
			} else if MapID == gameTypes.MapId_LhzDun02.Uint32() {
				b.GC.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
				time.Sleep(time.Millisecond * 500)
				b.GC.MoveChartWait(Pos_03[1])
				time.Sleep(time.Millisecond * 500)
				b.GC.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, b.GC.Role.GetPos())
				time.Sleep(time.Millisecond * 500)
			} else if MapID == gameTypes.MapId_LhzDun02West.Uint32() {
				b.GC.GoToMap(gameTypes.MapId_LhzDun01.Uint32())
				time.Sleep(time.Millisecond * 500)
				b.GC.MoveChartWait(Pos_03[0])
				time.Sleep(time.Millisecond * 500)
				b.GC.ExitMapPos(gameTypes.MapId_LhzDun01.Uint32(), 2, b.GC.Role.GetPos())
				time.Sleep(time.Millisecond * 500)
				b.GC.MoveChartWait(Pos_03[2])
				time.Sleep(time.Millisecond * 500)
				b.GC.ExitMapPos(gameTypes.MapId_LhzDun02.Uint32(), 2, b.GC.Role.GetPos())
			}

			b.EnableGodMode()
			b.GC.CheckDraculaBuff()

		} else {
			if b.GC.Configs.HuntConfig.UseDoubleEXP && b.GC.Role.GetBuffById(6062) == nil {
				b.UsesEXP()
			} else if b.GC.Configs.HuntConfig.TimerFly > 0 && time.Since(lastFlyTime) > time.Second*time.Duration(b.GC.Configs.HuntConfig.TimerFly) {
				b.useFlyWing()
				lastFlyTime = time.Now()
			} else if !b.GC.IsMonsterInRange(b.GC.Configs.HuntConfig.TargetMonsters...) {
				b.logger.Infof("没有找到目标怪物")
				b.useFlyWing()
			} else if !b.fightStar {
				b.logger.Infof("附近找到目标怪物，开始自动挂机，坐稳了")
				b.fightCtx, b.fightCancel = context.WithCancel(context.Background())
				b.GC.EnableAutoAttack(b.fightCtx, b.GC.Configs.HuntConfig.TargetMonsters...)
				b.GC.Role.SetSkillCd(50057001, time.Now().Add(time.Second*4))
				b.fightStar = true
				lastPosUpdate = time.Now()
			}
		}

		time.Sleep(time.Second) // ✅ 控制主循环节奏
	}
}

func (b *HuntTask) Stop() {
	b.cancel()
}

func (b *HuntTask) getItemCount(ItemName string) uint32 {
	iData := b.GC.FindPackItemByName(ItemName, Cmd.EPackType_EPACKTYPE_MAIN)
	if iData == nil {
		return 0
	}
	return iData.GetBase().GetCount()
}

func (b *HuntTask) useFlyWing() {
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
	b.fightCancel()
	b.fightStar = false
	lastFlyTime = time.Now()
	time.Sleep(time.Second * 2)
}

func (b *HuntTask) buyFlyWing() {
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
			b.logger.Infof("购买10000苍蝇翅膀")
			b.GC.BuyShopItem(item, 10000)
		}
	}
}

func (b *HuntTask) Useskill() {
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

func (b *HuntTask) UsesEXP() {
	b.logger.Infof("使用羊羊助力!")
	num := int32(1)
	dir := int32(utils.GetNpcDataValByType(b.GC.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    b.GC.Role.Pos,
		Dir:    &dir,
	}
	b.GC.SkillCmd(1556001, pData, true)
}

func NewHuntTask(ctx context.Context, gc *gameConnection.GameConnection) *HuntTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	fiightCtx, fightCancel := context.WithCancel(newCtx)
	return &HuntTask{
		GC:            gc,
		ctx:           newCtx,
		cancel:        cancel,
		logWriter:     mw,
		logger:        logger,
		startTime:     time.Now(),
		lastPosUpdate: time.Now(),
		flyMutex:      sync.Mutex{},
		lastPos:       Cmd.ScenePos{},
		fightCancel:   fightCancel,
		fightCtx:      fiightCtx,
	}
}

func (b *HuntTask) EnableGodMode() {
	time.Sleep(time.Millisecond * 1000)
	b.useFlyWing()
	time.Sleep(time.Millisecond * 3200)
	b.Useskill()
	time.Sleep(time.Millisecond * 1000)
	b.Useskill()
	time.Sleep(time.Millisecond * 1000)
}
