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
	if d.GC.Configs.DailyTaskConfig.EnableItemCombine {
		d.logger.Infof("开始执行物品合成任务...")
		d.PerformItemCombineTask()
	}
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
		kanbanStateStarted   kanbanStateType = "started"
		kanbanStateMoving    kanbanStateType = "moving"
		kanbanStateCompleted kanbanStateType = "completed"
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
				d.GC.MoveChart((d.GC.ParsePos(-23437, 16, 500)))
				time.Sleep(time.Millisecond * 1500)
				_ = d.GC.MoveToNpcWait("委托看板")
				time.Sleep(time.Second)
				_, err := d.GC.VisitNpcByName("委托看板")
				if err != nil {
					d.logger.Errorf("访问委托看板失败: %v", err)
					return
				}

				// 开始看板任务
				d.logger.Infof("开始自动看板任务...")
				quest, err := d.GC.GetQuestList(Cmd.EQuestList_EQUESTLIST_CANACCEPT, 101)
				if err != nil {
					d.logger.Errorf("获取看板任务列表失败: %v", err)
					return
				}
				if len(quest.GetList()) == 0 {
					d.logger.Infof("当前没有可接受的看板任务，任务完成。")
					return
				}
				for _, q := range quest.GetList() {
					d.logger.Infof("完成看板任务: %s", q.GetSteps()[0].GetConfig().GetName())
					d.GC.QuickSubmitWantedQuest(q.GetId())
					time.Sleep(time.Second * 2)
				}
				kanbanState = kanbanStateCompleted
			case kanbanStateCompleted:
				d.logger.Infof("看板任务已完成。")
				return
			}
		}
	}
}

func (d *DailyTask) performWasteLandWeedTask() {
	reward := d.checkRewardCount()
	if reward >= 20 {
		d.logger.Infof("当前荒境除草卡片礼包数量已达20个及以上，无需继续完成除草任务。")
		return
	}
	// 飞去荒境地图
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
			reward = d.checkRewardCount()
			if reward >= 20 {
				d.logger.Infof("当前荒境除草卡片礼包数量已达20个及以上，无需继续完成除草任务。")
				atkCancel()
				return
			}
			d.logger.Infof("当前荒境除草卡片礼包数量: %d，继续完成除草任务...", reward)
			time.Sleep(time.Second * 5)
		}
	}
}

func (d *DailyTask) checkRewardCount() uint32 {
	// 检查荒境除草卡片礼包数量
	i, err := d.GC.GetItemCount(80030004, Cmd.ESource_ESOURCE_REWARD)
	if err != nil {
		d.logger.Errorf("获取荒境除草卡片礼包数量失败: %v", err)
		return 0
	}
	return i.GetCount()
}

func (d *DailyTask) PerformItemCombineTask() {
	d.logger.Infof("执行物品合成任务中...")
	items := []string{
		"精装卡册的残页",
		"卡册残页",
	}
	for _, itemName := range items {
		d.logger.Infof("开始合成物品: %s", itemName)
		i := d.GC.FindPackItemByName(itemName, Cmd.EPackType_EPACKTYPE_MAIN)
		if i == nil {
			d.logger.Infof("背包中没有找到物品 %s，跳过合成。", itemName)
			continue
		}
		d.logger.Infof("背包中找到物品 %s %d个，开始合成...", itemName, i.GetBase().GetCount())
		iName, ok := d.GC.ItemsByName[itemName]
		if !ok {
			d.logger.Errorf("物品 %s 未在物品配置表中找到，无法合成。", itemName)
			continue
		}
		composeId, _ := iName.Items[0].ComposeId.Int64()

		for {
			d.GC.ProduceItem(uint32(composeId))
			i = d.GC.FindPackItemByName(itemName, Cmd.EPackType_EPACKTYPE_MAIN)
			if i == nil {
				d.logger.Infof("背包中没有找到物品 %s，合成完成。", itemName)
				break
			}
			if i.GetBase().GetCount() < 15 {
				d.logger.Infof("背包中物品 %s 数量不足15个，合成完成。", itemName)
				break
			}
			d.logger.Infof("背包中还有物品 %s %d个，继续合成...", itemName, i.GetBase().GetCount())
			time.Sleep(time.Millisecond * 500)
		}
		time.Sleep(time.Second * 1)
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
