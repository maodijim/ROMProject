package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	_ "net/http/pprof"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	autoLvVer = "v0.0.3"
)

var (
	disable, cancelAutoAtk = context.WithCancel(context.Background())
	lvUpChoice             = -1
	goblinLvUpChoice       = -1
	mountainLvUpChoice     = -1
	items                  = &utils.ItemsLoader{}
	g                      = &gameConnection.GameConnection{}
	npcList                = map[string][]int{}
	// 北区练级点
	northUPPos = [][]Cmd.ScenePos{
		// 右下
		{
			g.ParsePos(41716, 2851, -48749),
		},
		// 左下
		{
			g.ParsePos(-13186, 2835, -45833),
		},
		// 左上
		{
			g.ParsePos(30676, 2851, 100065),
			g.ParsePos(6317, 2851, 101614),
		},
		// 右下2
		{
			g.ParsePos(25073, 2852, -26499),
			g.ParsePos(63123, 2835, -20295),
		},
		// 右上
		{
			g.ParsePos(58811, 2851, -116973),
		},
	}
	goblinPos = [][]Cmd.ScenePos{
		// 中下,
		{
			g.ParsePos(22303, 8698, -18031),
		},
		// 左下
		{
			g.ParsePos(-49275, 12865, -67336),
		},
		// 中下2
		{
			g.ParsePos(22055, 8838, -6449),
		},
		// 左下
		{
			g.ParsePos(-2916, 8697, -976),
		},
		// 左
		{
			g.ParsePos(-27150, 10798, 24148),
		},
	}
	mountainPos = [][]Cmd.ScenePos{
		// 中下
		{
			g.ParsePos(21006, -1524, 3415),
		},
		// 中右
		{
			g.ParsePos(41908, -1223, 27975),
		},
		// 右上
		{
			g.ParsePos(21006, -1524, 3415),
			g.ParsePos(-39263, 12865, -62454),
		},
		// 左下
		// {
		// 	g.ParsePos(-34175, -9845, -4474),
		// },
	}
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
}

type StopAttackCondition struct {
	JobLevel   uint64
	BaseLevel  uint64
	NoMonster  bool
	Standstill bool
}

func southGateScript() {

	// 南门的脚本
	// 拿新手礼包
	if g.Role.GetMapId() != 2 {
		log.Warnf("角色不在南门地图")
		return
	} else if g.Role.GetJobLevel() >= 10 && g.Role.GetProfession() != Cmd.EProfession_EPROFESSION_NOVICE {
		log.Warnf("角色已经达到10级并转职")
		g.ExitMap(gameTypes.MapId_Protera.Uint32())
		return
	}
	g.RunQuestStep(10101, 0, 0, 0)
	time.Sleep(1 * time.Second)
	g.RunQuestStep(80000015, 0, 0, 0)
	time.Sleep(3 * time.Second)
	g.MoveChartWait(g.ParsePos(-30560, 20, 53584))
	time.Sleep(2 * time.Second)
	// g.VisitNpc(2147483836)
	_, _ = g.VisitNpcByName("宝雅")
	time.Sleep(3 * time.Second)

	// 新手礼包任务
	log.Infof("开始新手礼包任务")
	g.RunQuestStep(10100, 0, 0, 0)
	time.Sleep(2 * time.Second)
	g.RunQuestStep(10100, 0, 0, 2)
	time.Sleep(3 * time.Second)

	// 穿新手装备打怪
	item := g.FindPackItemByName("百万集结礼包", Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil {
		log.Warnf("使用百万集结礼包")
		g.UseItem(item.GetBase().GetGuid(), 1)
	} else {
		log.Warnf("没有找到百万集结礼包")
	}
	time.Sleep(3 * time.Second)
	putOnEquip("光明圣徒战靴", Cmd.EEquipPos_EEQUIPPOS_SHOES)
	putOnEquip("光明圣徒之盾", Cmd.EEquipPos_EEQUIPPOS_SHIELD)
	putOnEquip("光明圣徒之铠", Cmd.EEquipPos_EEQUIPPOS_ARMOUR)
	putOnEquip("光明圣徒披风", Cmd.EEquipPos_EEQUIPPOS_ROBE)
	putOnEquip("圣徒项链", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY1)
	putOnEquip("圣徒之戒", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY2)

	g.MoveChartWait(g.ParsePos(10129, 34, 3834))

	// 打怪
	lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 10, JobLevel: 10, NoMonster: false}, "绿棉虫")

	// 接转职任务
	log.Infof("开始接转职任务")
	// 对话宝雅
	g.MoveChartWait(g.ParsePos(-30560, 20, 53584))
	time.Sleep(9 * time.Second)
	_, err := g.VisitNpcByName("宝雅")
	if err != nil {
		log.Errorf("对话宝雅失败，%v", err)
	}
	// g.VisitNpc(2147483836)
	g.RunQuestStep(40001, 0, 0, 0)

	// 去普隆德拉
	g.MoveChartWait(g.ParsePos(-7452, 153, 69612))

	g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
}

func pronteraScript() {
	if g.Role.GetMapId() != 1 {
		log.Warnf("角色不在普隆德拉地图")
		return
	} else if g.Role.GetJobLevel() < 10 || g.Role.GetRoleLevel() < 10 {
		log.Errorf("角色基础/职业等级不足10级, 无法转职")
		return
	} else if g.Role.GetProfession() == Cmd.EProfession_EPROFESSION_ARCHER {
		log.Warnf("角色已经是弓箭手了")
		return
	}
	log.Info("开始执行普成脚本")
	// 第一个转职任务
	g.MoveChartWait(g.ParsePos(8533, -3320, -59734))
	g.RunQuestStep(40100, 0, 0, 0)
	time.Sleep(2 * time.Second)

	// 赛尼亚
	log.Infof("开始赛尼亚任务")
	// g.MoveChartWait(g.ParsePos(13935, -3320, -60135))
	err := g.MoveToNpcWait("赛尼亚")
	if err != nil {
		log.Errorf("移动到赛尼亚失败，%v", err)
	}
	time.Sleep(2 * time.Second)
	// g.VisitNpc(2147483734)
	g.VisitNpcByName("赛尼亚")
	g.RunQuestStep(40002, 0, 0, 0)
	time.Sleep(1 * time.Second)
	g.RunQuestStep(40002, 0, 0, 2)
	time.Sleep(3 * time.Second)

	// 娜莎
	log.Infof("开始娜莎任务")
	// g.MoveChartWait(g.ParsePos(19855, -3320, -58490))
	err = g.MoveToNpcWait("娜莎")
	if err != nil {
		log.Errorf("移动到娜莎失败，%v", err)
	}
	time.Sleep(2 * time.Second)
	// g.VisitNpc(2147483738)
	_, err = g.VisitNpcByName("娜莎")
	if err != nil {
		log.Errorf("对话娜莎失败，%v", err)
	}
	g.RunQuestStep(40002, 0, 0, 4)
	time.Sleep(1 * time.Second)
	_, _ = g.VisitNpcByName("娜莎")
	g.RunQuestStep(40002, 0, 0, 5)
	_, _ = g.VisitNpcByName("娜莎")
	g.RunQuestStep(40002, 0, 0, 6)
	time.Sleep(1 * time.Second)
	_, _ = g.VisitNpcByName("娜莎")
	g.RunQuestStep(40002, 0, 0, 8)

	// 开始泉水拍照
	g.MoveChartWait(g.ParsePos(12351, -3312, -52460))
	time.Sleep(2 * time.Second)
	g.TakePhoto(nil, g.ParsePos(12351, -3312, -52460))
	time.Sleep(1 * time.Second)
	g.RunQuestStep(40002, 0, 0, 10)
	g.SceneryCmd(55)

	// 对话娜莎2
	// g.MoveChartWait(g.ParsePos(19855, -3320, -58490))
	_ = g.MoveToNpcWait("娜莎")
	time.Sleep(2 * time.Second)
	// g.VisitNpc(2147483738)
	_, _ = g.VisitNpcByName("娜莎")
	g.RunQuestStep(40002, 0, 0, 11)
	time.Sleep(3 * time.Second)
	_, _ = g.VisitNpcByName("娜莎")
	g.RunQuestStep(40003, 0, 0, 0)
	time.Sleep(3 * time.Second)

	// 对话卡普拉
	// g.MoveChartWait(g.ParsePos(20521, -3320, -49221))
	err = g.MoveToNpcWait("卡普拉服务人员")
	if err != nil {
		log.Errorf("移动到卡普拉失败，%v", err)
	}
	time.Sleep(3 * time.Second)
	// g.VisitNpc(2147483724)
	_, _ = g.VisitNpcByName("卡普拉服务人员")
	g.RunQuestStep(40003, 0, 0, 1)
	time.Sleep(2 * time.Second)
	_, _ = g.VisitNpcByName("卡普拉服务人员")
	g.RunQuestStep(40004, 0, 0, 0)
	time.Sleep(2 * time.Second)
	_, _ = g.VisitNpcByName("卡普拉服务人员")
	g.RunQuestStep(40004, 0, 0, 12)
	g.RunQuestStep(40004, 0, 0, 14)
	time.Sleep(2 * time.Second)
	_, _ = g.VisitNpcByName("卡普拉服务人员")
	g.RunQuestStep(40004, 0, 0, 15)
	time.Sleep(3 * time.Second)
	g.RunQuestStep(99999100, 0, 0, 0)
	g.QuestRaidCmd(311000001)
	g.RunQuestStep(311000001, 0, 0, 0)
	time.Sleep(2 * time.Second)
	g.RunQuestStep(40007, 0, 0, 1)
	time.Sleep(2 * time.Second)

	// 去转职大厅
	g.ExitMapWait(gameTypes.MapId_RoomAdvanced.Uint32())
	time.Sleep(2 * time.Second)
}

func jobScript() {
	if g.Role.GetMapId() != gameTypes.MapId_RoomAdvanced.Uint32() {
		log.Warnf("角色不在转职大厅地图")
		return
	} else if g.Role.GetJobLevel() < 10 || g.Role.GetRoleLevel() < 10 {
		log.Warnf("角色基础/职业等级不足10级, 无法转职")
		return
	} else if g.Role.GetProfession() == Cmd.EProfession_EPROFESSION_ARCHER {
		log.Warnf("角色已经是弓箭手了")
		putOnEquip("百万击破", Cmd.EEquipPos_EEQUIPPOS_HEAD)
		putOnEquip("圣徒之弓[1]", Cmd.EEquipPos_EEQUIPPOS_WEAPON)
		g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
		return
	} else if g.Role.GetProfession() != Cmd.EProfession_EPROFESSION_NOVICE {
		log.Warnf("角色不是初心者，无法转职")
		return
	}
	log.Info("开始执行转职脚本")

	// 转职会长
	log.Infof("开始转职会长")
	g.RunQuestStep(11000001, 0, 0, 0)
	time.Sleep(1 * time.Second)
	_ = g.MoveToNpcWait("希盖伊兹")
	_, _ = g.VisitNpcByName("希盖伊兹")
	g.RunQuestStep(11000001, 0, 0, 5)
	time.Sleep(2 * time.Second)

	// 转职猎人
	log.Infof("开始转职猎人")
	// g.MoveChartWait(g.ParsePos(-2835, 110, -19238))
	_ = g.MoveToNpcWait("卡巴克")

	_, _ = g.VisitNpcByName("卡巴克")
	g.RunQuestStep(11040011, 0, 9, 0)

	_, _ = g.VisitNpcByName("卡巴克")
	g.RunQuestStep(11040011, 0, 0, 1)

	_, _ = g.VisitNpcByName("卡巴克")
	g.RunQuestStep(11040011, 0, 0, 2)

	g.Answer(0, 401, 2)
	g.Answer(0, 402, 2)
	g.Answer(0, 403, 1)
	time.Sleep(1 * time.Second)

	_, _ = g.VisitNpcByName("卡巴克")
	g.RunQuestStep(11040011, 0, 0, 22)
	time.Sleep(1 * time.Second)

	_, _ = g.VisitNpcByName("卡巴克")
	g.RunQuestStep(11040011, 0, 6, 23)
	time.Sleep(1 * time.Second)
	g.QuestRaidCmd(11040011)
	time.Sleep(1 * time.Second)

	// 转职地图
	g.ChangeMap(10041)
	time.Sleep(2 * time.Second)
	if g.Role.GetMapId() != 10041 {
		log.Warnf("转职训练地图失败")
	} else {
		log.Infof("转职训练地图成功")
		g.MoveChartWait(g.ParsePos(-3377, 5732, 3193))
		// 清理所有怪物
		// 第一轮
		disable, cancelAutoAtk = context.WithCancel(context.Background())
		lvUp(disable, cancelAutoAtk, StopAttackCondition{NoMonster: true}, "魔化树精")
		time.Sleep(5 * time.Second)
		// 第二轮
		disable, cancelAutoAtk = context.WithCancel(context.Background())
		lvUp(disable, cancelAutoAtk, StopAttackCondition{NoMonster: true}, "魔化树精")
		time.Sleep(5 * time.Second)

		// 拉杆回家
		// g.MoveChartWait(g.ParsePos(-10741, 5732, 3311))
		err := g.MoveToNpcWait("拉杆")
		if err != nil {
			log.Warnf("拉杆失败: %v", err)
			g.MoveChartWait(g.ParsePos(-10741, 5732, 3311))
		}
		time.Sleep(3 * time.Second)
		// g.VisitNpc(2174823957)
		_, err = g.VisitNpcByName("拉杆")
		if err != nil {
			log.Warnf("拉杆失败: %v", err)
			g.VisitNpc(2174823957)
		}
		g.RunQuestStep(11570001, 0, 0, 0)
		time.Sleep(1 * time.Second)
		_, _ = g.VisitNpcByName("拉杆")
		g.RunQuestStep(11570001, 0, 0, 2)
		time.Sleep(3 * time.Second)
		g.ExitMapWait(gameTypes.MapId_RoomAdvanced.Uint32())
		g.Reconnect()
		time.Sleep(10 * time.Second)
	}

	// 回到转职大厅转职
	if g.Role.GetMapId() != 1001 {
		log.Warnf("角色不在转职大厅地图")
	} else {
		// g.MoveChartWait(g.ParsePos(-2835, 110, -19238))
		_ = g.MoveToNpcWait("卡巴克")
		time.Sleep(2 * time.Second)
		// g.VisitNpc(2147492054)
		_, _ = g.VisitNpcByName("卡巴克")
		g.RunQuestStep(11040011, 0, 8, 25)
		g.QuestRaidCmd(11040011)
		time.Sleep(1 * time.Second)
		g.ChangeMap(10042)
		time.Sleep(3 * time.Second)
	}

	// 转职礼堂
	if g.Role.GetMapId() != 10042 {
		log.Warnf("转职礼堂失败")
	} else {
		log.Infof("转职礼堂成功")
		g.MoveChart(g.ParsePos(0, 316, 3000))
		time.Sleep(5 * time.Second)
		g.MoveChart(g.ParsePos(0, 309, 1800))
		time.Sleep(5 * time.Second)
		g.MoveChart(g.ParsePos(0, 718, 29500))
		time.Sleep(9 * time.Second)
		g.RunQuestStep(11040012, 0, 0, 0)
		g.RunQuestStep(11040012, 0, 0, 1)
		time.Sleep(5 * time.Second)
		err := g.MoveToNpcWait("卡巴克")
		if err != nil {
			log.Errorf("卡巴克对话失败: %v", err)
		}
		_, _ = g.VisitNpcByName("卡巴克")
		g.RunQuestStep(11140014, 0, 2, 0)
		time.Sleep(3 * time.Second)
		g.ExitMapWait(gameTypes.MapId_RoomAdvanced.Uint32())
		g.Reconnect()
		time.Sleep(10 * time.Second)
	}

	if g.Role.GetMapId() != 1001 {
		log.Warnf("角色不在转职大厅地图")
	} else {
		if g.Role.GetProfession() == Cmd.EProfession_EPROFESSION_ARCHER {
			log.Infof("转职成功")
		} else {
			log.Warnf("转职失败")
		}

		// 学习技能
		// g.MoveChartWait(g.ParsePos(-55, 124, -15889))
		// g.VisitNpc(2147492053)
		err := g.MoveToNpcWait("希盖伊兹")
		if err != nil {
			log.Errorf("希盖伊兹对话失败: %v", err)
		}
		_, _ = g.VisitNpcByName("希盖伊兹")
		g.RunQuestStep(11500006, 0, 0, 0)
		time.Sleep(5 * time.Second)

		// g.MoveChartWait(g.ParsePos(302, 128, -15622))
		err = g.MoveToNpcWait("赛尼亚")
		if err != nil {
			log.Errorf("赛尼亚对话失败: %v", err)
		}
		// g.VisitNpc(2191695762)
		_, _ = g.VisitNpcByName("赛尼亚")
		g.RunQuestStep(11500006, 0, 11, 3)
		_, _ = g.VisitNpcByName("赛尼亚")
		g.RunQuestStep(11500006, 0, 0, 14)
		time.Sleep(2 * time.Second)

		_ = g.MoveToNpcWait("希盖伊兹")
		_, _ = g.VisitNpcByName("希盖伊兹")
		g.RunQuestStep(400040001, 0, 0, 0)
		g.RunQuestStep(400040001, 0, 0, 3)
		g.RunQuestStep(400040001, 0, 3, 3)
	}

	if g.Role.GetProfession() == Cmd.EProfession_EPROFESSION_NOVICE {
		g.ExitMapWait(gameTypes.MapId_ProteraSouth.Uint32())
		log.Fatalf("還是初心者转职失敗, 返回南門")
	}

	putOnEquip("百万击破", Cmd.EEquipPos_EEQUIPPOS_HEAD)
	putOnEquip("圣徒之弓[1]", Cmd.EEquipPos_EEQUIPPOS_WEAPON)

	g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
}

func westGate() {
	if g.Role.GetProfession() != Cmd.EProfession_EPROFESSION_ARCHER {
		log.Errorf("角色不是弓箭手")
		g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
		return
	}
	if g.GetAtk() > 190 {
		log.Infof("攻击力已经大于190，不需要再打西门")
		return
	}
	if g.Role.GetMapId() != 5 {
		g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
		time.Sleep(2 * time.Second)
		g.ExitMapWait(gameTypes.MapId_ProteraWest.Uint32())
		g.Reconnect()
		time.Sleep(10 * time.Second)
	}
	log.Warnf("攻击力不足在西门练练吧")

	g.MoveChartWait(g.ParsePos(-54016, 10133, -17510))
	lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 18}, "溜溜猴")
	g.ExitMapWait(gameTypes.MapId_Protera.Uint32())
	time.Sleep(2 * time.Second)
}

func northGate() {
	if g.Role.GetProfession() != Cmd.EProfession_EPROFESSION_ARCHER {
		log.Errorf("角色不是弓箭手")
		return
	} else if g.Role.GetRoleLevel() >= 37 {
		log.Infof("角色等级已经大于37，不需要再打北门")
		return
	}
	if g.Role.GetMapId() == 1 {
		log.Warnf("角色在主城地图")
		g.ExitMapWait(gameTypes.MapId_ProteraRoom1F.Uint32())
		time.Sleep(2 * time.Second)
		g.ChangeMap(gameTypes.MapId_ProteraRoom1F.Uint32())
		if g.Role.GetMapId() != gameTypes.MapId_ProteraRoom1F.Uint32() {
			log.Warnf("角色不在普隆德拉皇家区1F地图")
			return
		} else {
			log.Infof("角色在普隆德拉皇家区1F地图")
			g.ExitMapWait(gameTypes.MapId_ProteraNorth.Uint32())
		}
	} else if g.Role.GetRoleLevel() >= 37 {
		log.Errorf("角色等级大于37，不需要再打北门")
		return
	} else if g.Role.GetMapId() != gameTypes.MapId_ProteraNorth.Uint32() {
		g.ExitMapWait(gameTypes.MapId_ProteraNorth.Uint32())
		g.ChangeMap(gameTypes.MapId_ProteraNorth.Uint32())
	}

	// 随机选择 1 或者 0 , 2, 3
	choice := lvUpChoice
	if lvUpChoice < 0 || choice > len(northUPPos)-1 {
		rand.Seed(time.Now().UnixNano())
		choice = rand.Intn(4)
	}
	var name string
	switch choice {
	case 0:
		name = "右下"
	case 1:
		name = "左下"
	case 2:
		name = "左上"
	case 3:
		name = "右下2"
	}
	log.Infof("选择练级点 %d, %s %v", choice, name, northUPPos[choice][len(northUPPos[choice])-1])
	for _, pos := range northUPPos[choice] {
		g.MoveChartWait(pos)
	}

	lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 37}, "森灵")
	time.Sleep(2 * time.Second)
}

func goblinForest() {
	if g.Role.GetRoleLevel() < 37 {
		log.Errorf("角色等级不足37级")
		return
	} else if g.Role.GetMapId() != gameTypes.MapId_GoblinForest.Uint32() {
		g.ExitMapWait(gameTypes.MapId_GoblinForest.Uint32())
		time.Sleep(2 * time.Second)
	}

	log.Info("在哥布林森林练级")

	choice := goblinLvUpChoice
	if choice < 0 || choice > len(goblinPos)-1 {
		rand.Seed(time.Now().UnixNano())
		choice = rand.Intn(len(goblinPos))
	}
	switch choice {
	case 0:
		log.Infof("选择练级点 %d, %s %v", choice, "中下", goblinPos[choice][len(goblinPos[choice])-1])
	case 1:
		log.Infof("选择练级点 %d, %s %v", choice, "左下", goblinPos[choice][len(goblinPos[choice])-1])
	case 2:
		log.Infof("选择练级点 %d, %s %v", choice, "中下2", goblinPos[choice][len(goblinPos[choice])-1])
	case 3:
		log.Infof("选择练级点 %d, %s %v", choice, "左下", goblinPos[choice][len(goblinPos[choice])-1])
	}
	for _, pos := range goblinPos[choice] {
		g.MoveChartWait(pos)
	}
	lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 50}, "喷射哥布灵", "弓箭哥布灵")
}

func mjolnirMountains() {
	if g.Role.GetRoleLevel() < 37 {
		log.Errorf("角色等级不足37级")
		return
	} else if g.Role.GetMapId() != gameTypes.MapId_MjolnirMountains.Uint32() {
		g.ExitMapWait(gameTypes.MapId_MjolnirMountains.Uint32())
		g.Reconnect()
		time.Sleep(10 * time.Second)
	}

	log.Info("在妙勒尼山脉练级")

	choice := mountainLvUpChoice
	if choice < 0 || choice > len(mountainPos)-1 {
		rand.Seed(time.Now().UnixNano())
		choice = rand.Intn(len(mountainPos))
	}
	switch choice {
	case 0:
		log.Infof("选择练级点 %d, %s %v", choice, "中下", goblinPos[choice][len(goblinPos[choice])-1])
	case 1:
		log.Infof("选择练级点 %d, %s %v", choice, "中右", goblinPos[choice][len(goblinPos[choice])-1])
	case 2:
		log.Infof("选择练级点 %d, %s %v", choice, "右上", goblinPos[choice][len(goblinPos[choice])-1])
	case 3:
		log.Infof("选择练级点 %d, %s %v", choice, "左下", goblinPos[choice][len(goblinPos[choice])-1])
	}
	for _, pos := range mountainPos[choice] {
		g.MoveChartWait(pos)
	}
	ticker := time.NewTicker(10 * time.Second)
	goBackTicker := time.NewTicker(5 * time.Minute)
	startPos := g.Role.GetPos()
	lastPosUpdate := time.Now()
	for {
		select {
		case <-time.After(1 * time.Second):
			if g.Role.GetRoleLevel() < 50 {
				lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 50}, "蜂兵")
			} else if g.Role.GetRoleLevel() < 60 {
				choice = rand.Intn(len(mountainPos) - 1)
				lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 60}, "噬人花")
			}
		case <-goBackTicker.C:
			log.Infof("5分钟后回到起始点")
			cancelAutoAtk()
			for _, pos := range mountainPos[choice] {
				g.MoveChartWait(pos)
			}
			disable, cancelAutoAtk = context.WithCancel(context.Background())
			if g.Role.GetRoleLevel() < 50 {
				lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 50}, "蜂兵")
			} else if g.Role.GetRoleLevel() < 60 {
				lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 60}, "噬人花")
			}
		case <-ticker.C:
			if g.Role.GetRoleLevel() >= 60 {
				log.Infof("角色等级达到60级")
				ticker.Stop()
				goBackTicker.Stop()
				cancelAutoAtk()
				return
			}
			if !g.Role.IsEqualPos(startPos) {
				lastPosUpdate = time.Now()
				startPos = g.Role.GetPos()
			}
			if time.Since(lastPosUpdate) > 60*time.Second {
				log.Infof("角色卡住了, 重新进入地图")
				cancelAutoAtk()
				for _, pos := range mountainPos[choice] {
					g.MoveChartWait(pos)
				}
				disable, cancelAutoAtk = context.WithCancel(context.Background())
				if g.Role.GetRoleLevel() < 50 {
					lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 50}, "蜂兵")
				} else if g.Role.GetRoleLevel() < 60 {
					lvUp(disable, cancelAutoAtk, StopAttackCondition{BaseLevel: 60}, "噬人花")
				}
			}
		}
	}
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
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-stopNpc:
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
	return stopNpc
}

func expBuff() {
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				item := g.FindPackItemByName("暖身料理", Cmd.EPackType_EPACKTYPE_MAIN)
				if item != nil {
					log.Infof("使用暖身料理 %d 个", item.GetBase().GetCount())
					g.UseItem(item.GetBase().GetGuid(), item.GetBase().GetCount())
				}
				item = g.FindPackItemByName("锁链雷锭", Cmd.EPackType_EPACKTYPE_MAIN)
				if item == nil {
					log.Warnf("未找到锁链雷锭")
				} else {
					log.Infof("背包还有锁链雷锭 %d 个", item.GetBase().GetCount())
					hasBuff := g.GetBuffNameByRegex("锁链雷锭")
					if hasBuff != "" {
						log.Warnf("已经有锁链经验buff了")
						continue
					}
					log.Infof("使用锁链雷锭")
					g.UseItem(item.GetBase().GetGuid(), 1)
					log.Infof("还剩下 %d 个锁链雷锭", item.GetBase().GetCount()-1)
				}
			}
		}
	}()

}

func lvUp(disable context.Context, cancelAutoAtk context.CancelFunc, condition StopAttackCondition, monsters ...string) {
	go func() {
		targetId := g.AtkStat.GetCurrentTargetId()
		lastPosUpdate := time.Now()
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-disable.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if targetId != 0 && g.AtkStat.GetCurrentTargetId() == targetId && time.Since(lastPosUpdate) > time.Second*30 {
					log.Infof("卡住了")
					useFlyWing()
				} else if targetId == 0 && time.Since(lastPosUpdate) > time.Second*60 {
					log.Infof("没有目标卡住了")
					useFlyWing()
				} else if g.AtkStat.GetCurrentTargetId() != targetId {
					targetId = g.AtkStat.GetCurrentTargetId()
					lastPosUpdate = time.Now()
				}
			}
		}
	}()
	stopNpc := printNearbyNpcs(g)
	if condition.Standstill {
		g.AtkStat.SetStandstill(true)
	}
	g.EnableAutoAttack(disable, monsters...)
	jlv := g.Role.GetJobLevel()
	blv := g.Role.GetRoleLevel()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-disable.Done():
			log.Infof("停止练级")
			stopNpc <- true
			cancelAutoAtk()
			return
		case <-ticker.C:
			jlv = g.Role.GetJobLevel()
			blv = g.Role.GetRoleLevel()
			log.Infof("现在角色等级: %d, 职业等级: %d", blv, jlv)
			log.Infof("当前血量: %f %%, 经验: %d", g.GetHpPer()*100, g.Role.GetRoleExp())
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
				g.AtkStat.SetStandstill(false)
				stopNpc <- true
				cancelAutoAtk()
			}
			if g.GetHpPer() == 0 {
				log.Infof("死亡")
				stopNpc <- true
				cancelAutoAtk()
				go func() {
					g.Reconnect()
					time.Sleep(5 * time.Second)

					mjolnirMountains()
				}()
			}
		}
	}
}

func AutoAttrPoint() {
	go func() {
		ticker := time.NewTicker(120 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if g.Role.GetTotalPoint() >= 2 {
					log.Infof("自动分配属性点")
					var attrs []int32
					var err error
					if g.Role.GetRoleLevel() >= 35 {
						attrs, err = g.AddAttrPoint(0, 1, 0, 0, 0, 0)
					} else if g.Role.GetRoleLevel() < 35 || g.Role.GetRoleLevel() >= 50 {
						attrs, err = g.AddAttrPoint(0, 0, 0, 0, 1, 0)
					} else {
						attrs, err = g.AddAttrPoint(0, 0, 0, 0, 1, 2)
					}
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

func AutoSkillLearn() {
	skillsToLearn := []string{
		"苍鹰之眼", "10",
		"元素箭矢", "10",
		"二连矢", "10",
		"鹗枭之眼", "10",
	}
	go func() {
		ticker := time.NewTicker(120 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for i := 0; i < len(skillsToLearn); i += 2 {
					if g.Role.GetSkillPoint() >= 1 {
						skillName := skillsToLearn[i]
						skillLevel, _ := strconv.ParseUint(skillsToLearn[i+1], 10, 32)
						lv := g.GetLearnedSkillLevelByName(skillName)
						if lv < uint32(skillLevel) {
							id := g.GetSkillIdByName(skillName, uint32(lv+1))
							g.LevelUpSkill(
								[]uint32{
									id,
								},
								Cmd.ELevelupType_ELEVELUPTYPE_MT,
							)
							log.Infof("自动学习技能 %s, id %d, lv %d", skillName, id, lv+1)
							time.Sleep(1 * time.Second)
						}
					}
				}
			}
		}
	}()
}

func useFlyWing() {
	g.UseFlyWing()
	item := g.FindPackItemById(5024, Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil && item.GetBase().GetCount() > 0 {
		log.Infof("使用苍蝇翅膀 还有%d个", item.GetBase().GetCount())
		if item.GetBase().GetCount() < 10 {
			buyFlyWing()
		}
	} else {
		log.Warn("没有找到苍蝇翅膀")
		_ = g.GetAllPackItems()
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
			log.Infof("购买10苍蝇翅膀")
			g.BuyShopItem(item, 10)
		}
	}
}

func main() {
	log.Infof("ROM auto level up version: %s", autoLvVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	enablePprof := flag.Bool("pprof", false, "Enable pprof")
	choice := flag.Int("choice", -1,
		"Choice of north gate lv up point: \n"+
			"0 bottom right\n"+
			"1 bottom left\n"+
			"2 top left\n"+
			"3 bottom right2\n"+
			"4 top right\n",
	)
	goblinChoice := flag.Int("goblinChoice", -1,
		"Choice of goblin forest lv up point: \n"+
			"0 bottom mid\n"+
			"1 bottom left\n"+
			"2 bottom mid 2\n"+
			"3 bottom left2\n"+
			"4 top left\n",
	)
	mountainChoice := flag.Int("mountainChoice", -1,
		"Choice of mountain lv up point: \n"+
			"0 bottom mid\n"+
			"1 mid left\n"+
			"2 mid right\n"+
			"3 top right\n",
	)
	flag.Parse()
	lvUpChoice = *choice
	goblinLvUpChoice = *goblinChoice
	mountainLvUpChoice = *mountainChoice
	if *enablePprof {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	items = utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	g.AddNotifier("INTER_QUESTION")
	g.GameServerLogin()
	g.ChangeMap(g.Role.GetMapId())
	time.Sleep(3 * time.Second)

	AutoAttrPoint()
	AutoSkillLearn()
	expBuff()
	southGateScript()
	pronteraScript()
	jobScript()
	putOnEquip("百万击破", Cmd.EEquipPos_EEQUIPPOS_HEAD)
	putOnEquip("圣徒之弓[1]", Cmd.EEquipPos_EEQUIPPOS_WEAPON)
	putOnEquip("光明圣徒战靴", Cmd.EEquipPos_EEQUIPPOS_SHOES)
	putOnEquip("光明圣徒之盾", Cmd.EEquipPos_EEQUIPPOS_SHIELD)
	putOnEquip("光明圣徒之铠", Cmd.EEquipPos_EEQUIPPOS_ARMOUR)
	putOnEquip("光明圣徒披风", Cmd.EEquipPos_EEQUIPPOS_ROBE)
	putOnEquip("圣徒项链", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY1)
	putOnEquip("圣徒之戒", Cmd.EEquipPos_EEQUIPPOS_ACCESSORY2)
	westGate()
	northGate()
	// goblinForest()
	mjolnirMountains()
	<-time.After(15 * time.Second)
}
