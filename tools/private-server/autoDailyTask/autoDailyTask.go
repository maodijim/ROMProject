package autoDailyTask

import (
	"context"
	"io"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

type DailyTask struct {
	GC        *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	logWriter io.Writer
	logger    *log.Logger
}

func (d *DailyTask) SetLogger(writer io.Writer) {
	mw := io.MultiWriter(d.GC.LogWriter(), writer)
	d.logWriter = mw
	d.logger.SetOutput(mw)
}

func (d *DailyTask) Start() {
	d.GC.ShouldChangeScene = true
	d.GC.GameServerLogin()

	d.startDailyTasks()
}

func (d *DailyTask) Stop() {
	d.cancel()
}

func (d *DailyTask) GetContext() context.Context {
	return d.ctx
}

func (d *DailyTask) startDailyTasks() {
	d.logger.Infof("开始执行每日日常任务...")
	if d.GC.Configs.DailyTaskConfig.EnableKanBan {
		d.logger.Infof("开始执行看板任务...")
		d.performKanBanTask()
	}
	if d.GC.Configs.DailyTaskConfig.EnableWasteLandWeed {
		d.logger.Infof("开始执行荒地除草任务...")
		d.performWasteLandWeedTask()
	}
}

func (d *DailyTask) performKanBanTask() {
	d.logger.Infof("执行看板任务中...")
	type kanbanStateType string
	const (
		kanbanStateStarted kanbanStateType = "started"
		kanbanStateMoving  kanbanStateType = "moving"
	)
	kanbanState := kanbanStateStarted
	for {
		select {
		case <-d.ctx.Done():
			log.Infof("看板任务已停止。")
			return
		default:
			d.logger.Infof("看板任务状态: %s", kanbanStateStarted)
			switch kanbanState {
			case kanbanStateStarted:
				if d.GC.Role.GetMapId() != gameTypes.MapId_Protera.Uint32() {
					d.logger.Infof("当前地图不是普隆德拉，飞去普隆德拉中...")
					d.GC.GoToMap(gameTypes.MapId_Protera.Uint32())
					time.Sleep(time.Second * 5)
				}
				kanbanState = kanbanStateMoving
			case kanbanStateMoving:
				if d.GC.Role.GetMapId() != gameTypes.MapId_Protera.Uint32() {
					d.logger.Infof("当前地图不是普隆德拉, 重新开始看板任务...")
					kanbanState = kanbanStateStarted
					continue
				}
				// 移动到看板NPC位置
				d.logger.Infof("移动到看板NPC位置...")
				time.Sleep(time.Second * 2)
				// d.GC.MoveChart((d.GC.ParsePos(-25099, -513, -46410)))
				time.Sleep(time.Second * 2)
			}
		}
	}
}

func (d *DailyTask) performWasteLandWeedTask() {
	if d.GC.Role.GetMapId() != gameTypes.MapId_Wasteland.Uint32() {
		if !utils.Contains(d.GC.GotoList.GetMapid(), gameTypes.MapId_Wasteland.Uint32()) {
			d.logger.Warnf("荒境不在可去地图列表中，无法前往荒境完成除草任务，请检查角色已经开通荒境地图传送！")
			return
		}
		d.logger.Infof("当前地图不是荒境，飞去荒境中...")
		d.GC.GoToMap(gameTypes.MapId_Wasteland.Uint32())
		time.Sleep(time.Second * 5)
	}
	// 开始除草任务
	d.logger.Infof("开始自动除草任务...")
	i := d.GC.FindPackItemByName("荒境除草卡片礼包", Cmd.EPackType_EPACKTYPE_MAIN)
	startItemCount := uint32(0)
	if i != nil {
		startItemCount = i.GetBase().GetCount()
	}
	d.logger.Infof("当前拥有荒境除草卡片礼包数量: %d", startItemCount)
	atkCtx, atkCancel := context.WithCancel(context.Background())
	d.GC.EnableAutoAttack(atkCtx, "洛阳荒草")
	for {
		select {
		case <-d.ctx.Done():
			atkCancel()
			d.logger.Infof("荒地除草任务已停止。")
			return
		default:
			i := d.GC.FindPackItemByName("荒境除草卡片礼包", Cmd.EPackType_EPACKTYPE_MAIN)
			newItemCount := uint32(0)
			if i != nil {
				newItemCount = i.GetBase().GetCount()
			}
			d.logger.Infof("获得荒境除草卡片礼包，当前数量: %d", newItemCount)
			if newItemCount >= startItemCount+20 {
				d.logger.Infof("本次除草任务已完成，获得荒境除草卡片礼包数量达到20个，任务结束。")
				atkCancel()
				return
			}
		}
	}
}

func NewDailyTask(ctx context.Context, gc *gameConnection.GameConnection) *DailyTask {
	taskCtx, cancel := context.WithCancel(ctx)
	logger := log.New()
	mw := io.MultiWriter(gc.LogWriter())
	logger.SetOutput(mw)
	logger.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	return &DailyTask{
		GC:        gc,
		ctx:       taskCtx,
		cancel:    cancel,
		logWriter: mw,
		logger:    logger,
	}
}
