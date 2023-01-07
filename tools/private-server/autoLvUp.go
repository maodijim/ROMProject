package main

import (
	"flag"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	autoLvVer = "v0.0.1"
)

var (
	items   = &utils.ItemsLoader{}
	g       = &gameConnection.GameConnection{}
	npcList = map[string][]int{}
	// 北区练级点
	lvUPPos = []*Cmd.ScenePos{
		// 右下
		g.ParsePos(41716, 2851, -48749),
		// 左下
		g.ParsePos(-13186, 2835, -45833),
		// 左上
		g.ParsePos(6317, 2851, -103412),
		// 右上
		g.ParsePos(58811, 2851, -116973),
	}
)

type StopAttackCondition struct {
	JobLevel  uint64
	BaseLevel uint64
	NoMonster bool
}

func southGateScript() {

	// 南门的脚本
	// 拿新手礼包
	if g.Role.GetMapId() != 2 {
		log.Warnf("角色不在南门地图")
		return
	} else if g.Role.GetJobLevel() >= 10 {
		log.Warnf("角色已经达到10级")
		g.ExitMap(1)
		return
	}
	g.RunQuestStep(10101, 0, 0, 0)
	g.RunQuestStep(80000015, 0, 0, 0)
	time.Sleep(2 * time.Second)
	g.MoveChartWait(g.ParsePos(-30560, 20, 53584))
	g.VisitNpc(2147483836)
	time.Sleep(3 * time.Second)

	// 新手礼包任务
	log.Infof("开始新手礼包任务")
	g.RunQuestStep(10100, 0, 0, 0)
	time.Sleep(1 * time.Second)
	g.RunQuestStep(10100, 0, 0, 2)
	time.Sleep(1 * time.Second)

	// 穿新手装备打怪
	item := g.FindItemByName("百万集结礼包", Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil {
		log.Warnf("使用百万集结礼包")
		g.UseItem(item.GetGuid(), 1)
	} else {
		log.Warnf("没有找到百万集结礼包")
	}
	time.Sleep(1 * time.Second)
	putOnEquip("光明圣徒战靴", Cmd.EEquipPos_EEQUIPPOS_SHOES)
	putOnEquip("光明圣徒之盾", Cmd.EEquipPos_EEQUIPPOS_SHIELD)
	putOnEquip("光明圣徒之铠", Cmd.EEquipPos_EEQUIPPOS_ARMOUR)
	putOnEquip("光明圣徒披风", Cmd.EEquipPos_EEQUIPPOS_ROBE)
	putOnEquip("圣徒项链", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY1)
	putOnEquip("圣徒之戒", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY2)

	g.MoveChartWait(g.ParsePos(10129, 34, 3834))

	// 打怪
	lvUp(StopAttackCondition{BaseLevel: 10, JobLevel: 10, NoMonster: false}, "绿棉虫")

	// 接转职任务
	log.Infof("开始接转职任务")
	g.MoveChartWait(g.ParsePos(-30560, 20, 53584))
	g.VisitNpc(2147483836)
	time.Sleep(3 * time.Second)
	g.RunQuestStep(40001, 0, 0, 0)

	// 去普隆德拉
	g.MoveChartWait(g.ParsePos(-7452, 153, 69612))

	g.ExitMapWait(1)
}

func pronteraScript() {
	if g.Role.GetMapId() != 1 {
		log.Warnf("角色不在普隆德拉地图")
		return
	} else if g.Role.GetJobLevel() < 10 || g.Role.GetRoleLevel() < 10 {
		log.Warnf("角色基础/职业等级不足10级, 无法转职")
		return
	} else if g.Role.GetProfession() == Cmd.EProfession_name[int32(Cmd.EProfession_EPROFESSION_ARCHER)] {
		log.Warnf("角色已经是弓箭手了")
		return
	}
	log.Info("开始执行普成脚本")
	// 第一个转职任务
	g.MoveChartWait(g.ParsePos(8533, -3320, -59734))
	g.RunQuestStep(40100, 0, 0, 0)

	// 塞尼亚
	log.Infof("开始塞尼亚任务")
	g.MoveChartWait(g.ParsePos(13935, -3320, -60135))
	time.Sleep(2 * time.Second)
	g.VisitNpc(2147483734)
	g.RunQuestStep(40002, 0, 0, 0)
	time.Sleep(1 * time.Second)
	g.RunQuestStep(40002, 0, 0, 2)
	time.Sleep(1 * time.Second)

	// 娜莎
	log.Infof("开始娜莎任务")
	g.MoveChartWait(g.ParsePos(19855, -3320, -58490))
	time.Sleep(2 * time.Second)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40002, 0, 0, 4)
	time.Sleep(1 * time.Second)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40002, 0, 0, 5)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40002, 0, 0, 6)
	time.Sleep(1 * time.Second)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40002, 0, 0, 8)

	// 开始泉水拍照
	g.MoveChartWait(g.ParsePos(12351, -3312, -52460))
	time.Sleep(2 * time.Second)
	g.TakePhoto(nil, g.ParsePos(12351, -3312, -52460))
	time.Sleep(1 * time.Second)
	g.RunQuestStep(40002, 0, 0, 10)
	g.SceneryCmd(55)

	// 对话娜莎2
	g.MoveChartWait(g.ParsePos(19855, -3320, -58490))
	time.Sleep(1 * time.Second)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40002, 0, 0, 11)
	time.Sleep(1 * time.Second)
	g.VisitNpc(2147483738)
	g.RunQuestStep(40003, 0, 0, 0)
	time.Sleep(1 * time.Second)

	// 对话卡普拉
	g.MoveChartWait(g.ParsePos(20521, -3320, -49221))
	g.VisitNpc(2147483724)
	g.RunQuestStep(40003, 0, 0, 1)
	g.VisitNpc(2147483724)
	g.RunQuestStep(40004, 0, 0, 0)
	time.Sleep(1 * time.Second)
	g.VisitNpc(2147483724)
	g.RunQuestStep(40004, 0, 0, 12)
	g.RunQuestStep(40004, 0, 0, 14)
	g.VisitNpc(2147483724)
	g.RunQuestStep(40004, 0, 0, 15)
	g.RunQuestStep(40007, 0, 0, 1)
	time.Sleep(1 * time.Second)

	// 去转职大厅
	g.ExitMapWait(1001)
	time.Sleep(1 * time.Second)
}

func jobScript() {
	if g.Role.GetMapId() != 1001 {
		log.Warnf("角色不在转职大厅地图")
		return
	} else if g.Role.GetJobLevel() < 10 || g.Role.GetRoleLevel() < 10 {
		log.Warnf("角色基础/职业等级不足10级, 无法转职")
		return
	} else if g.Role.GetProfession() == Cmd.EProfession_name[int32(Cmd.EProfession_EPROFESSION_ARCHER)] {
		log.Warnf("角色已经是弓箭手了")
		g.ExitMap(1)
		return
	}
	log.Info("开始执行转职脚本")

	// 转职会长
	log.Infof("开始转职会长")
	g.RunQuestStep(11000001, 0, 0, 0)
	g.MoveChartWait(g.ParsePos(-55, 124, -15832))
	g.VisitNpc(2147492053)
	g.RunQuestStep(11000001, 0, 0, 5)
	time.Sleep(2 * time.Second)

	// 转职猎人
	log.Infof("开始转职猎人")
	g.MoveChartWait(g.ParsePos(-2835, 110, -19238))

	g.VisitNpc(2147492054)
	g.RunQuestStep(11040011, 0, 9, 0)

	g.VisitNpc(2147492054)
	g.RunQuestStep(11040011, 0, 0, 1)

	g.VisitNpc(2147492054)
	g.RunQuestStep(11040011, 0, 0, 2)

	g.Answer(0, 401, 2)
	g.Answer(0, 402, 2)
	g.Answer(0, 403, 1)
	time.Sleep(1 * time.Second)

	g.VisitNpc(2147492054)
	g.RunQuestStep(11040011, 0, 0, 22)
	time.Sleep(1 * time.Second)

	g.VisitNpc(2147492054)
	g.RunQuestStep(11040011, 0, 6, 23)
	time.Sleep(1 * time.Second)
	g.QuestRaidCmd(11040011)
	time.Sleep(1 * time.Second)

	// 转职地图
	g.ChangeMap(10041)
	time.Sleep(2 * time.Second)
	if g.Role.GetMapId() != 10041 {
		log.Warnf("转职地图失败")
	} else {
		log.Infof("转职地图成功")
		g.MoveChartWait(g.ParsePos(-3377, 5732, 3193))
		// 清理所有怪物
		lvUp(StopAttackCondition{NoMonster: true}, "魔化树精")

		// 拉杆回家
		g.MoveChartWait(g.ParsePos(-10741, 5732, 3311))
		g.VisitNpc(2174823957)
		g.RunQuestStep(11570001, 0, 0, 0)
		g.VisitNpc(2174823957)
		g.RunQuestStep(11570001, 0, 0, 2)
		time.Sleep(1 * time.Second)
		g.ExitMap(1001)
		g.ChangeMap(1001)
		time.Sleep(1 * time.Second)
	}

	// 回到转职大厅转职
	if g.Role.GetMapId() != 1001 {
		log.Warnf("角色不在转职大厅地图")
	} else {
		g.MoveChartWait(g.ParsePos(-2835, 110, -19238))
		g.VisitNpc(2147492054)
		g.RunQuestStep(11040011, 0, 8, 25)
		g.QuestRaidCmd(11040011)
		time.Sleep(1 * time.Second)
		g.ChangeMap(10042)
		time.Sleep(2 * time.Second)
	}

	// 转职礼堂
	if g.Role.GetMapId() != 10042 {
		log.Warnf("转职礼堂失败")
	} else {
		log.Infof("转职礼堂成功")
		g.MoveChartWait(g.ParsePos(0, 316, 3000))
		time.Sleep(2 * time.Second)
		g.MoveChartWait(g.ParsePos(0, 309, 1800))
		time.Sleep(1 * time.Second)
		g.MoveChartWait(g.ParsePos(0, 718, 29500))
		time.Sleep(2 * time.Second)
		g.RunQuestStep(11040012, 0, 0, 0)
		time.Sleep(1 * time.Second)
		g.RunQuestStep(11040012, 0, 0, 1)
		time.Sleep(1 * time.Second)
		g.MoveChartWait(g.ParsePos(46, 718, 31114))
		g.VisitNpc(2176323690)
		g.RunQuestStep(11140014, 0, 2, 0)
		time.Sleep(1 * time.Second)
		g.ChangeMap(1001)
	}

	if g.Role.GetMapId() != 1001 {
		log.Warnf("角色不在转职大厅地图")
	} else {
		log.Infof("转职成功")
		g.MoveChartWait(g.ParsePos(-55, 124, -15889))
		g.VisitNpc(2147492053)
		g.RunQuestStep(11500006, 0, 0, 0)
		time.Sleep(3 * time.Second)
		g.MoveChartWait(g.ParsePos(302, 128, -15622))
		g.VisitNpc(2191695762)
		g.RunQuestStep(11500006, 0, 11, 3)
		g.VisitNpc(2191695762)
		g.RunQuestStep(11500006, 0, 0, 14)
	}

	putOnEquip("百万击破", Cmd.EEquipPos_EEQUIPPOS_HEAD)
	putOnEquip("圣徒之弓[1]", Cmd.EEquipPos_EEQUIPPOS_WEAPON)

	g.ExitMap(1)
}

func westGate() {
	if g.GetAtk() > 200 {
		log.Infof("攻击力已经大于200，不需要再打西门")
		return
	}
	if g.Role.GetMapId() != 5 {
		g.ExitMap(5)
	}
	log.Warnf("攻击力不足在西门练练吧")
	g.MoveChartWait(g.ParsePos(-54016, 10133, -17510))
	lvUp(StopAttackCondition{BaseLevel: 18}, "溜溜猴")
	g.ExitMap(1)
	time.Sleep(2 * time.Second)
}

func northGate() {
	if g.Role.GetMapId() == 1 {
		log.Warnf("角色在主城地图")
		g.ExitMap(47)
		time.Sleep(2 * time.Second)
		g.ChangeMap(47)
		if g.Role.GetMapId() != 47 {
			log.Warnf("角色不在普隆德拉皇家区1F地图")
			return
		} else {
			log.Infof("角色在普隆德拉皇家区1F地图")
			g.ExitMap(42)
		}
	} else if g.Role.GetMapId() != 42 {
		g.ExitMap(42)
		g.ChangeMap(42)
	}
	g.MoveChartWait(lvUPPos[0])
	lvUp(StopAttackCondition{BaseLevel: 50}, "森灵")
	time.Sleep(2 * time.Second)
}
func putOnEquip(name string, pos Cmd.EEquipPos) {
	err := g.EquipItemByName(name, pos, Cmd.EEquipOper_EEQUIPOPER_ON)
	if err != nil {
		log.Warnf("穿戴%s失败: %s", name, err.Error())
	} else {
		log.Warnf("穿戴%s成功", name)
	}
}

func printNearbyNpcs(g *gameConnection.GameConnection) (stopNpc chan bool) {
	log.Printf("Nearby NPCs:")
	stopNpc = make(chan bool)
	go func() {
		for {
			select {
			case <-stopNpc:
				return
			case <-time.Tick(6 * time.Second):
				npcList := map[string][]int{}
				for _, npc := range g.MapNpcs {
					if _, ok := npcList[npc.GetName()]; ok {
						npcList[npc.GetName()][0] += 1
					} else {
						npcList[npc.GetName()] = []int{1, int(*npc.Id)}
					}
				}
				for k, v := range npcList {
					log.Printf("NPC: %s, 数量: %d", k, v)
				}
			}
		}
	}()
	return stopNpc
}

func expBuff() {
	go func() {
		for {
			select {
			case <-time.Tick(60 * time.Second):
				item := g.FindItemByName("暖身料理", Cmd.EPackType_EPACKTYPE_MAIN)
				if item != nil {
					log.Infof("使用暖身料理 %d 个", item.GetCount())
					g.UseItem(item.GetGuid(), item.GetCount())
				}
				item = g.FindItemByName("锁链雷锭", Cmd.EPackType_EPACKTYPE_MAIN)
				if item == nil {
					log.Warnf("未找到锁链雷锭")
				} else {
					log.Infof("背包还有锁链雷锭 %d 个", item.GetCount())
					hasBuff := g.GetBuffNameByRegex("锁链雷锭")
					if hasBuff != "" {
						log.Warnf("已经有锁链经验buff了")
						continue
					}
					log.Infof("使用锁链雷锭")
					g.UseItem(item.GetGuid(), 1)
					log.Infof("还剩下 %d 个锁链雷锭", item.GetCount()-1)
				}
			}
		}
	}()

}

func lvUp(condition StopAttackCondition, monsters ...string) {
	stopNpc := printNearbyNpcs(g)
	disable := make(chan bool)
	g.EnableAutoAttack(monsters, disable)
	jlv := g.Role.GetJobLevel()
	blv := g.Role.GetRoleLevel()
	for {
		select {
		case <-time.Tick(10 * time.Second):
			jlv = g.Role.GetJobLevel()
			blv = g.Role.GetRoleLevel()
			log.Infof("现在角色等级: %d, 职业等级: %d", blv, jlv)
			wanted := condition.NoMonster
			if condition.NoMonster {
				counter := 0
				for _, monster := range npcList {
					if monster[1] >= 10000 {
						counter += monster[0]
					}
				}
				if counter == 0 {
					log.Infof("怪物已经全部死亡")
				} else {
					wanted = !wanted
				}
			}
			if jlv >= condition.JobLevel && blv >= condition.BaseLevel && condition.NoMonster == wanted {
				log.Infof("练级完成")
				stopNpc <- true
				disable <- true
				return
			}
		}
	}
}

func AutoAttrPoint() {
	go func() {
		for {
			select {
			case <-time.Tick(120 * time.Second):
				if g.Role.GetTotalPoint() >= 2 {
					log.Infof("自动分配属性点")
					attrs, err := g.AddAttrPoint(0, 0, 0, 0, 1, 0)
					if err != nil {
						log.Warnf("自动分配属性点失败: %s", err.Error())
					} else {
						log.Infof("自动分配属性点成功: %v", attrs)
					}
				}
			}
		}
	}()
}

func main() {
	log.Infof("ROM auto level up version: %s", autoLvVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()
	items = utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	g.GameServerLogin()
	g.AddNotifier("INTER_QUESTION")

waitForInGame:
	for {
		select {
		case <-time.Tick(3 * time.Second):
			log.Infof("等待角色进入游戏")
			if g.Role.GetInGame() {
				log.Infof("角色已进入游戏")
				break waitForInGame
			}
		case <-time.After(15 * time.Second):
			log.Warnf("等待进入游戏超时")
			break waitForInGame
		}
	}
	_ = g.GetAllPackItems()
	g.ChangeMap(g.Role.GetMapId())
	AutoAttrPoint()

	expBuff()
	southGateScript()
	pronteraScript()
	jobScript()
	westGate()
	northGate()

	<-time.After(20 * time.Second)
}
