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
					b.GC.CheckuseFlyWing()
				} else if b.fightStar && targetId == 0 && time.Since(lastPosUpdate) > time.Second*10 {
					b.logger.Infof("没有目标卡住了")
					b.GC.CheckuseFlyWing()
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

		b.GC.InMap(gameTypes.MapNameZh[b.GC.Configs.HuntConfig.Map].Uint32())

		b.UseNature(gameTypes.GetNatureTypeFromZhFast(b.GC.Configs.HuntConfig.NatureType))

		if b.GC.Configs.HuntConfig.UseDoubleEXP && b.GC.Role.GetBuffById(6062) == nil {
			b.UsesEXP()
		} else if b.GC.Configs.HuntConfig.TimerFly > 0 && time.Since(lastFlyTime) > time.Second*time.Duration(b.GC.Configs.HuntConfig.TimerFly) {
			b.GC.CheckuseFlyWing()
			lastFlyTime = time.Now()
		} else if !b.GC.IsMonsterInRange(b.GC.Configs.HuntConfig.TargetMonsters...) {
			b.logger.Infof("没有找到目标怪物")
			b.GC.CheckuseFlyWing()
		} else if !b.fightStar {
			b.logger.Infof("附近找到目标怪物，开始自动挂机，坐稳了")
			b.fightCtx, b.fightCancel = context.WithCancel(context.Background())
			b.GC.EnableAutoAttack(b.fightCtx, b.GC.Configs.HuntConfig.TargetMonsters...)
			b.GC.Role.SetSkillCd(50057001, time.Now().Add(time.Second*4))
			b.fightStar = true
			lastPosUpdate = time.Now()
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

func (b *HuntTask) UseNature(Nature gameTypes.NatureType) {
	if Nature != "" {
		if b.GC.Role.GetProfession() >= Cmd.EProfession_EPROFESSION_ARCHER && b.GC.Role.GetProfession() <= Cmd.EProfession_EPROFESSION_RANGER {
			if Nature == gameTypes.NatureType_Fire {
				b.GC.UseElementArrow(gameTypes.FireArrow)
			} else if Nature == gameTypes.NatureType_Water {
				b.GC.UseElementArrow(gameTypes.WaterArrow)
			} else if Nature == gameTypes.NatureType_Wind {
				b.GC.UseElementArrow(gameTypes.WindArrow)
			} else if Nature == gameTypes.NatureType_Earth {
				b.GC.UseElementArrow(gameTypes.EarthArrow)
			} else if Nature == gameTypes.NatureType_Holy {
				b.GC.UseElementArrow(gameTypes.SliverArrow)
			}
		} else {
			if Nature == gameTypes.NatureType_Fire {
				b.GC.UseElementStone(gameTypes.FireStone)
			} else if Nature == gameTypes.NatureType_Water {
				b.GC.UseElementStone(gameTypes.WaterStone)
			} else if Nature == gameTypes.NatureType_Wind {
				b.GC.UseElementStone(gameTypes.WindStone)
			} else if Nature == gameTypes.NatureType_Earth {
				b.GC.UseElementStone(gameTypes.EarthStone)
			}
		}

	}
}
