package autoLottery

import (
	"context"
	"io"
	"os"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"

	log "github.com/sirupsen/logrus"
)

type LotteryTask struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter io.Writer
	logger    *log.Logger
}

func (l *LotteryTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(l.GC.LogWriter(), writer)
	l.logWriter = mw
	l.logger.SetOutput(mw)
}

func (l *LotteryTask) Start() {
	l.GC.GameServerLogin()

	if l.GC.Role.GetMapId() != gameTypes.MapId_Protera.Uint32() {
		l.logger.Warnf("当前地图不是普隆德拉，飞去普隆德拉中...")
		time.Sleep(time.Second * 5)
		_ = l.GC.GoToGear(gameTypes.MapId_Protera.Uint32())
		l.GC.ChangeMap(gameTypes.MapId_Protera.Uint32())
	}

	// 移动到抽奖NPC位置
	l.logger.Infof("移动到抽奖NPC位置...")
	time.Sleep(time.Second * 2)
	_ = l.GC.MoveChartWait(l.GC.ParsePos(-25099, -513, -46410))
	time.Sleep(time.Second * 2)
	lotteryName := l.GC.Configs.LotteryConfig.LotteryType
	l.GC.MoveToNpcWait(lotteryName)
	time.Sleep(time.Second)
	npc, err := l.GC.VisitObjectByName(lotteryName)
	if err != nil {
		l.logger.Errorf("访问NPC失败: %v", err)
		return
	}
	time.Sleep(time.Second * 2)

	// 开始抽奖任务
	l.lotteryTask(npc)
}

func (l *LotteryTask) Stop() {
	l.cancel()
}

func (l *LotteryTask) GetContext() context.Context {
	return l.ctx
}

func (l *LotteryTask) lotteryTask(npc Cmd.MapNpc) {
	l.logger.Infof("开始自动抽奖任务...")
	l.GC.SetQueryTimeout(time.Millisecond * 1000)
	lotteryName := l.GC.Configs.LotteryConfig.LotteryType
	lotteryType, ok := gameTypes.LotteryNameZh[l.GC.Configs.LotteryConfig.LotteryType]
	if !ok {
		l.logger.Errorf("未知的抽奖类型: %s", l.GC.Configs.LotteryConfig.LotteryType)
		return
	}
	maxCount := uint32(0)
	dailyCount := uint32(0)
	lotteryPrice := gameTypes.LotteryTypePriceMap[lotteryType]
	var lotteryInfo *Cmd.QueryLotteryInfo

	for {
		select {
		case <-l.ctx.Done():
			l.logger.Infof("自动抽奖任务已停止.")
			return
		default:
			// 检查抽奖机会
			if dailyCount == 0 && lotteryInfo == nil {
				l.logger.Infof("检查抽奖机会...")
				lotteryInfo = l.GC.QueryLotteryInfo(lotteryType)
				dailyCount = lotteryInfo.GetTodayCnt()
				maxCount = lotteryInfo.GetMaxCnt()
			}

			if dailyCount >= maxCount {
				l.logger.Infof("今日%s抽奖机会%d次已用完", lotteryName, maxCount)
				l.Stop()
				continue
			} else {
				l.logger.Infof("今日%s抽奖机会剩余: %d次", lotteryName, maxCount-dailyCount)
			}
			// 在这里添加具体的抽奖逻辑
			drawCount := min(l.GC.Configs.LotteryConfig.DrawCount, maxCount-dailyCount)
			l.logger.Infof("开始抽奖: 抽取 %d 次 %s", drawCount, lotteryName)
			l.logger.Infof("剩余票券: %d", l.GetTicketCount(lotteryInfo))
			ticketId := uint32(0)
			if l.GC.Configs.LotteryConfig.UseTickets {
				ticketId = lotteryInfo.GetInfos()[0].GetSubInfo()[0].GetRecoverItemid()
				drawCount = min(drawCount, l.GetTicketCount(lotteryInfo)/30)
				if drawCount == 0 {
					l.logger.Infof("票券不足，无法继续抽奖。")
					l.Stop()
					continue
				}
			} else {
				drawCount = min(drawCount, uint32(l.GC.Role.GetLottery()/lotteryPrice))
				if drawCount == 0 {
					l.logger.Infof("猫币不足，无法继续抽奖。")
					l.Stop()
					continue
				}
			}

			lotteryCmd := l.GC.LotteryDraw(lotteryType, drawCount, uint32(lotteryPrice)*drawCount, ticketId, npc.GetId())
			l.logger.Infof("抽奖结果: 获得 %d 个物品", len(lotteryCmd.GetItems()))
			if lotteryCmd != nil && lotteryCmd.GetTodayCnt() > 0 {
				dailyCount = lotteryCmd.GetTodayCnt()
			}
		}
	}
}

func (l *LotteryTask) GetTicketCount(lotteryInfo *Cmd.QueryLotteryInfo) (totalCount uint32) {
	if lotteryInfo.GetInfos() == nil {
		return 0
	}
	if lotteryInfo.GetInfos()[0].GetSubInfo() == nil {
		return 0
	}
	ticketId := lotteryInfo.GetInfos()[0].GetSubInfo()[0].GetRecoverItemid()
	item := l.GC.FindPackItemByIdAll(ticketId, Cmd.EPackType_EPACKTYPE_MAIN)
	for _, it := range item {
		totalCount += it.GetBase().GetCount()
	}
	return totalCount
}

func NewLotteryTask(ctx context.Context, gc *gameConnection.GameConnection) *LotteryTask {
	newCtx, cancel := context.WithCancel(ctx)
	mw := io.MultiWriter(os.Stdout, gc.LogWriter())
	logger := log.New()
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &LotteryTask{
		GC:        gc,
		ctx:       newCtx,
		cancel:    cancel,
		logWriter: mw,
		logger:    logger,
	}
}
