package autoDailyTask

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
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
		d.performItemCombineTask()
	}
	if d.GC.Configs.DailyTaskConfig.EnableKanBan {
		d.logger.Infof("开始执行看板任务...")
		d.performKanBanTask()
	}
	if d.GC.Configs.DailyTaskConfig.EnableWasteLandWeed {
		d.logger.Infof("开始执行荒地除草任务...")
		d.performWasteLandWeedTask()
	}
	if d.GC.Configs.DailyTaskConfig.EnableYuno && d.GC.Role.GetRoleLevel() >= 100 {
		d.logger.Infof("开始执行朱诺任务...")
		d.performYunoTask()
	} else if d.GC.Configs.DailyTaskConfig.EnableCrack {
		d.logger.Infof("开始执行裂隙/朱诺任务...")
		if d.GC.Role.GetRoleLevel() >= 100 {
			d.logger.Warnf("当前角色等级已达100级及以上，推荐执行朱诺任务...")
		}
		d.performCrackTask()
	}
	d.logger.Infof("每日日常任务执行完毕。")
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
					d.GC.ChangeMap(gameTypes.MapId_Protera.Uint32())
					time.Sleep(time.Second * 2)
					continue
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
					time.Sleep(time.Millisecond * 2500)
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
		i := d.GC.FindPackItemByName("荒境除草卡片礼包", Cmd.EPackType_EPACKTYPE_MAIN)
		var startItemCount uint32 = 0
		if i != nil {
			startItemCount = i.GetBase().GetCount()
		}
		d.logger.Infof("当前拥有荒境除草卡片礼包数量: %d", startItemCount)
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
		d.GC.MoveChartWait(d.GC.ParsePos(129957, -21393, -177177))
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
			i = d.GC.FindPackItemByName("荒境除草卡片礼包", Cmd.EPackType_EPACKTYPE_MAIN)
			startItemCount = uint32(0)
			if i != nil {
				startItemCount = i.GetBase().GetCount()
			}
			d.logger.Infof("当前拥有荒境除草卡片礼包数量: %d", startItemCount)
			d.logger.Infof("荒地除草任务已停止。")
			return
		default:
			reward = d.checkRewardCount()
			if reward >= 20 {
				d.logger.Infof("当前荒境除草卡片礼包数量已达%d个及以上，无需继续完成除草任务。", reward)
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

func (d *DailyTask) performItemCombineTask() {
	d.logger.Infof("执行物品合成任务中...")
	items := map[string]uint32{
		"精装卡册的残页": 15,
		"卡册残页":    15,
	}
	for itemName, reqCount := range items {
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
			if i.GetBase().GetCount() < reqCount {
				d.logger.Infof("背包中物品 %s 数量不足%d个，合成完成。", itemName, reqCount)
				break
			}
			d.logger.Infof("背包中还有物品 %s %d个，继续合成...", itemName, i.GetBase().GetCount())
			time.Sleep(time.Millisecond * 500)
		}
		time.Sleep(time.Second * 1)
	}
}

func (d *DailyTask) performCrackTask() {
	// 裂隙任务
	teamCfg := config.TeamConfig{
		LeaderName: d.GC.Role.GetRoleName(),
	}
	stage := 0

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Infof("裂隙任务已停止。")
			return
		default:
			switch stage {
			case 0:
				// 0. 检查是否完成任务
				sealQuest, _ := d.GC.QuerySealQuest()
				if sealQuest.GetDonetimes() >= sealQuest.GetMaxtimes() {
					d.logger.Infof("今日裂隙任务已完成，任务结束。")
					return
				}
				d.logger.Infof("今日剩余裂隙任务次数: %d/%d", sealQuest.GetMaxtimes()-sealQuest.GetDonetimes(), sealQuest.GetMaxtimes())
				stage = 1
			case 1:
				// 1. 检查是否组队
				if d.GC.GetCurrentTeamName() == "" {
					d.logger.Infof("当前没有队伍，开始组队...")
					d.GC.AutoCreateJoinTeam(teamCfg)
					time.Sleep(time.Second * 5)
				} else {
					stage = 2
				}
			case 2:
				// 2. 检查是否队长
				if d.GC.Role.AcceptSeal.GetSeal() == uint32(gameTypes.SealQuestType_WestGate) {
					d.logger.Infof("当前已经接受西门裂隙任务，前往裂隙位置...")
					stage = 5
					continue
				}
				// 如果不是队长，退队重新组队
				if d.GC.GetTeamLeaderName(true) != d.GC.Role.GetRoleName() {
					d.logger.Infof("当前不是队长，退队重新组队...")
					d.GC.ExitTeam()
					time.Sleep(time.Second * 3)
					d.GC.AutoCreateJoinTeam(teamCfg)
					time.Sleep(time.Second * 5)
				} else {
					stage = 3
				}
			case 3:
				// 3. 去普隆德拉接取西门裂隙任务
				if d.GC.Role.GetMapId() != gameTypes.MapId_Protera.Uint32() {
					d.logger.Infof("当前地图不是普隆德拉，飞去普隆德拉中...")
					d.GC.GoToMap(gameTypes.MapId_Protera.Uint32())
					time.Sleep(time.Second * 5)
				} else {
					stage = 4
				}
			case 4:
				// 4. 接受裂隙任务
				d.logger.Infof("完成西门裂隙任务...")
				d.GC.MoveChart(d.GC.ParsePos(-23437, 16, 500))
				time.Sleep(time.Millisecond * 1500)
				_ = d.GC.MoveToNpcWait("裂隙监视者")
				time.Sleep(time.Second)
				_, err := d.GC.VisitNpcByName("裂隙监视者")
				if err != nil {
					d.logger.Errorf("访问裂隙监视者失败: %v", err)
					continue
				}
				sealQuest, _ := d.GC.QuerySealQuest()
				hasWestGateQuest := false
				for _, q := range sealQuest.GetConfigid() {
					if q == uint32(gameTypes.SealQuestType_WestGate) {
						hasWestGateQuest = true
						break
					}
				}
				if !hasWestGateQuest {
					d.logger.Infof("当前没有西门裂隙任务，取消任务...")
					return
				}

				quest, _ := d.GC.AcceptSealQuest(gameTypes.SealQuestType_WestGate)
				d.logger.Infof("已接受%s任务: %s", gameTypes.SealQuestType_WestGate.String(), quest.GetPos())
				stage = 5
			case 5:
				sealQuest, _ := d.GC.QuerySealQuest()
				if sealQuest.GetMaxtimes() != 0 && sealQuest.GetDonetimes() >= sealQuest.GetMaxtimes() {
					d.logger.Infof("今日裂隙任务已完成%d次，任务结束。", sealQuest.GetDonetimes())
					return
				}
				// 5. 前往裂隙位置
				if d.GC.Role.GetMapId() != gameTypes.MapId_ProteraWest.Uint32() {
					d.logger.Infof("当前地图不是普隆德拉西门，飞去普隆德拉西门中...")
					d.GC.GoToMap(gameTypes.MapId_ProteraWest.Uint32())
					time.Sleep(time.Second * 5)
				}
				curSealPos := d.GC.Role.AcceptSeal.GetPos()
				d.logger.Infof("前往西门裂隙位置... 坐标: X=%d, Y=%d, Z=%d", curSealPos.X, curSealPos.Y, curSealPos.Z)
				d.GC.MoveChartWait(*curSealPos)
				time.Sleep(time.Second * 2)
				d.logger.Infof("到达西门裂隙位置，开始完成任务...")
				time.Sleep(time.Second * 5)
				npc, _ := d.GC.VisitNpcByName("时空裂隙")
				time.Sleep(time.Second * 2)
				d.GC.BeginSealQuest(npc.GetId())
				d.logger.Infof("等待裂隙任务完成")
				time.Sleep(time.Second * 10)
				for {
					if _, ok := gameTypes.MapNameZh[d.GC.Role.GetMapName()]; !ok {
						d.logger.Infof("等待裂隙消失...")
						time.Sleep(time.Second * 5)
					} else {
						d.logger.Infof("西门裂隙任务已完成，继续下一次...")
						break
					}
				}
				time.Sleep(time.Second * 5)
			}
		}
	}
}

func (d *DailyTask) performYunoTask() {
	// 朱诺任务
	lName := d.GC.Configs.DailyTaskConfig.YunoTeamLeader
	if lName == "" {
		lName = d.GC.Role.GetRoleName()
	}
	teamCfg := config.TeamConfig{
		LeaderName: lName,
	}
	stage := 0

	for {
		select {
		case <-d.ctx.Done():
			d.logger.Infof("朱诺任务已停止。")
			return
		default:
			switch stage {
			case 0:
				// 0. 检查是否完成任务
				sealQuest, _ := d.GC.QuerySealQuest()
				if !d.GC.Configs.DailyTaskConfig.YunForceContinue && sealQuest.GetDonetimes() >= sealQuest.GetMaxtimes() {
					d.logger.Infof("今日朱诺任务已完成，任务结束。")
					return
				}
				d.logger.Infof("今日剩余朱诺任务次数: %d/%d", sealQuest.GetMaxtimes()-sealQuest.GetDonetimes(), sealQuest.GetMaxtimes())
				stage = 1
			case 1:
				// 1. 检查是否队长
				// 如果不是队长，退队重新组队
				if d.GC.GetTeamLeaderName(true) != lName {
					d.logger.Infof("当前队长是%s, 不是%s，退队重新组队...", d.GC.GetTeamLeaderName(true), d.GC.Configs.DailyTaskConfig.YunoTeamLeader)
					d.GC.ExitTeam()
					time.Sleep(time.Second * 3)
					d.GC.AutoCreateJoinTeam(teamCfg)
					time.Sleep(time.Second * 5)
				} else {
					stage = 2
				}
			case 2:
				// 3. 如果是队长去朱诺接取朱诺任务
				// 如果不是队长跟随队长等待
				if d.GC.IsTeamLeader(d.GC.Role.GetRoleId(), true) && d.GC.Role.GetMapId() != gameTypes.MapId_Yuno.Uint32() {
					d.logger.Infof("当前地图不是朱诺，飞去朱诺中...")
					d.GC.GoToMap(gameTypes.MapId_Yuno.Uint32())
					time.Sleep(time.Second * 5)
				} else if d.GC.IsTeamLeader(d.GC.Role.GetRoleId(), true) && d.GC.Role.GetMapId() == gameTypes.MapId_Yuno.Uint32() {
					stage = 3
				} else {
					d.logger.Infof("跟随队长 %s 中...", d.GC.GetTeamLeaderName(true))
					d.GC.FollowUser(d.GC.GetTeamLeader(true))
					stage = 4
				}
			case 3:
				// 3. 接受朱诺任务
				d.logger.Infof("开始朱诺任务...")
				time.Sleep(time.Millisecond * 1500)
				_ = d.GC.MoveToNpcWait("卡莱克·乌迪")
				time.Sleep(time.Second)
				_, err := d.GC.VisitNpcByName("卡莱克·乌迪")
				if err != nil {
					d.logger.Errorf("访问卡莱克·乌迪失败: %v", err)
					continue
				}
				d.GC.RunQuestStep(390990001, 0, 0, 0)

				time.Sleep(time.Second * 2)
				if d.GC.Role.GetMapId() == 60123 {
					stage = 5
				} else {
					d.logger.Infof("等待传送到火焰之地地图...")
					time.Sleep(time.Second * 3)
				}
			case 5:
				sealQuest, _ := d.GC.QuerySealQuest()
				if !d.GC.Configs.DailyTaskConfig.YunForceContinue && sealQuest.GetMaxtimes() != 0 && sealQuest.GetDonetimes() >= sealQuest.GetMaxtimes() {
					d.logger.Infof("今日朱诺任务已完成%d次，任务结束。", sealQuest.GetDonetimes())
					return
				}

				d.logger.Infof("今日朱诺任务完成次数: %d/%d", sealQuest.GetDonetimes(), sealQuest.GetMaxtimes())

				if d.GC.Role.GetMapId() != 60123 {
					d.logger.Infof("当前地图不是火焰之地, 等待传送中...")
					time.Sleep(time.Second * 5)
					continue
				}
				if d.GC.Role.GetRoleName() != d.GC.GetTeamLeaderName(true) {
					d.logger.Infof("等待队长完全任务完成中...")
					time.Sleep(time.Second * 10)
					continue
				}
				// 5. 完成朱诺副本
				// 5.1 对话NPC开始任务
				d.logger.Infof("执行朱诺副本...")
				d.GC.UpdateQueryTimeout(time.Second * 5)
				d.GC.UseElementArrow(gameTypes.SliverArrow)
				_ = d.GC.MoveToNpcWait("卡莱克·乌迪")
				time.Sleep(time.Second)
				_, err := d.GC.VisitNpcByName("卡莱克·乌迪")
				if err != nil {
					d.logger.Errorf("访问卡莱克·乌迪失败: %v", err)
					continue
				}

				lastStepId := uint32(20839)
				stepSync := &Cmd.FubenStepSyncCmd{
					Id: &lastStepId,
				}

				stepCtx, stepCancel := context.WithCancel(d.ctx)
				defer stepCancel()

				go func() {
					// 接收副本步骤通知
					for {
						select {
						case <-stepCtx.Done():
							d.logger.Infof("朱诺副本步骤同步监听已停止。")
							return
						case <-stepCtx.Done():
							d.logger.Infof("朱诺副本步骤同步监听已停止。")
							return
						default:
							d.GC.AddNotifier(gameTypes.NtfType_FubenStepSync)
							res, err := d.GC.WaitForResponse(gameTypes.NtfType_FubenStepSync)
							if err != nil {
								continue
							}
							if res != nil {
								stepSync = res.(*Cmd.FubenStepSyncCmd)
								lastStepId = stepSync.GetId()
							}
						}
					}
				}()

				// 20839 朱诺的危机
				d.yunoStep(stepSync)
				ticker := time.NewTicker(time.Second * 10)
			mainLoop:
				for {
					select {
					case <-d.ctx.Done():
						d.logger.Infof("朱诺副本步骤同步任务已停止。")
						return
					case <-ticker.C:
						if d.GC.Configs.DailyTaskConfig.YunForceContinue {
							d.logger.Infof("强制继续模式开启，忽略任务完成检查，继续执行朱诺任务...")
							ticker.Stop()
							break
						}
						sealQuest, _ = d.GC.QuerySealQuest()
						if sealQuest.GetMaxtimes() != 0 && sealQuest.GetDonetimes() >= sealQuest.GetMaxtimes() {
							d.logger.Infof("今日朱诺任务已完成%d次，任务结束。", sealQuest.GetDonetimes())
							ticker.Stop()
							return
						}
					default:
						if d.GC.Role.GetMapId() == gameTypes.MapId_Yuno.Uint32() {
							d.logger.Infof("朱诺副本未完成，在朱诺城...")
							stepCancel()
							stage = 3
							curHpPer := d.GC.GetHpPer()
							for curHpPer < 0.95 {
								d.logger.Infof("当前血量低于95%%，等待恢复血量...")
								d.logger.Infof("当前血量%d, %f%%", d.GC.GetCurrentHp(), d.GC.GetHpPer()*100)
								time.Sleep(time.Second * 5)
								curHpPer = d.GC.GetHpPer()
							}
							break mainLoop
						}
						if lastStepId == 21076 {
							d.yunoStep(stepSync)
							stepCancel()
							d.logger.Infof("朱诺副本已完成，返回朱诺城...")
							d.GC.GoToMap(gameTypes.MapId_Yuno.Uint32())
							return
						} else if stepSync.GetId() == lastStepId {
							d.yunoStep(stepSync)
						}
						time.Sleep(time.Millisecond * 1500)
					}
				}

				// 5.2 第一关
				// 20844
				// 20860
				// 20862
				// 20882
				// 20883
				// 20893
				// 20894
				// 20899

				// 5.3 第二关
				// 20903
				// 20911
				// 20913
				// 20922
				// 20941
				// 20943 挑战火焰教主教
				// 20958
				// 20963
				// 20967
				// 20975

				// 5.4 第三关
				// 20976 追寻左右护法的踪迹
				// 20977
				// 20978 追寻左右护法的踪迹
				// 20979
				// 20980 继续前进追寻左右护法
				// 20981
				// 20983
				// 21001

				// 5.5 第四关
				// 21005 前往阻止仪式
				// 21013
				// 21017
				// 21022
				// 21031 击败火焰教左右护法
				// 21038
				// 21042
				// 21046
				// 21054
				// 21060

				// 5.6 第五关
				// 21066 击败火焰领主-莫特奈尔
				// 21071
				// 21076
			}
		}
	}
}

func (d *DailyTask) yunoStep(stepSync *Cmd.FubenStepSyncCmd) {
	d.GC.FubenStepSync(stepSync.GetId())
	msg := fmt.Sprintf("%s - %s",
		stepSync.GetConfig().GetDescInfo(),
		stepSync.GetConfig().GetTraceInfo())
	d.logger.Infof("下个朱诺副本步骤: %d, %s", stepSync.GetId(), msg)
	if stepSync.GetConfig().GetDescInfo() == "走" {
		nextPos := d.GC.ExtractFubenStepSyncPos(stepSync)
		d.logger.Infof("走到指定位置... %s", nextPos.String())
		d.GC.MoveChartWait(nextPos)
	}
	if stepSync.GetConfig().GetDescInfo() == "说" && stepSync.GetConfig().GetContent() == "visit" {
		npcPar := d.GC.GetFubenStepSyncParam(stepSync, "npc")
		npcId, _ := strconv.ParseUint(npcPar.GetValue(), 10, 32)
		pos := d.GC.ExtractFubenStepSyncPos(stepSync)
		d.logger.Infof("走到NPC位置... %s", pos.String())
		d.GC.MoveChartWait(pos)
		time.Sleep(time.Second * 2)
		d.logger.Infof("访问NPC... %d", npcId)
		_ = d.GC.MoveToNpcIdWait(uint32(npcId))
	}
	if stepSync.GetConfig().GetDescInfo() == "说" && stepSync.GetConfig().GetContent() == "use" {
		pos := d.GC.ExtractFubenStepSyncPos(stepSync)
		d.logger.Infof("走到物品位置... %s", pos.String())
		d.GC.MoveChartWait(pos)
	}
	if stepSync.GetConfig().GetDescInfo() == "等击杀" {
		time.Sleep(time.Second * 3)
		d.logger.Infof("开始击杀怪物")
		if stepSync.GetId() == 20894 {
			d.killYunoTarget("火焰教大祭祀")
		} else if stepSync.GetId() == 20958 {
			d.killYunoTarget("火焰教主教")
		} else {
			d.killYunoTarget("火焰教左护法", "火焰教右护法", "火焰教信徒", "狂暴火精灵", "火焰领主莫特奈尔")
		}
	}
}

func (d *DailyTask) killYunoTarget(monsterNames ...string) {
	atkCtx, atkCancel := context.WithCancel(context.Background())
	d.GC.EnableAutoAttack(atkCtx, monsterNames...)
	time.Sleep(time.Second * 5)
	for {
		select {
		case <-d.ctx.Done():
			atkCancel()
			d.logger.Infof("朱诺副本击杀%s任务已停止。", monsterNames)
			return
		default:
			if d.GC.IsMonsterInRange(monsterNames...) {
				time.Sleep(time.Second * 5)
			} else {
				d.logger.Infof("%s已击杀完成。", monsterNames)
				atkCancel()
				return
			}
		}
	}
}

func NewDailyTask(ctx context.Context, gc *gameConnection.GameConnection) *DailyTask {
	taskCtx, cancel := context.WithCancel(ctx)
	logger := log.New()
	mw := io.MultiWriter(gc.LogWriter(), os.Stdout)
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
