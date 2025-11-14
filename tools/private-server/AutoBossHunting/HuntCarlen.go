package main

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"time"

	log "github.com/sirupsen/logrus"
)

type Work int

const (
	Init Work = iota
	TeleportMap
	MOVE_RABBIDSPOS
	CHECK_RABBIDS
	HUNT_RABBIDS
	MOVE_CARLENPOS
	CHECK_CARLEN
	HUNT_CARLEN
	End
)

var (
	WorkState  Work  = Init
	PosCount   int32 = int32(0)
	HaveCarlen bool  = false
	FindCarlen bool  = false
)

func i32(v int32) *int32 { return &v }

var RabbidsPos = []Cmd.ScenePos{
	{X: i32(10165), Y: i32(27), Z: i32(-57086)},
	{X: i32(20845), Y: i32(27), Z: i32(-85483)},
	{X: i32(-29854), Y: i32(27), Z: i32(-84827)},
	{X: i32(-41499), Y: i32(127), Z: i32(-51326)},
}

func HuntCarlen() {
	for {
		CheckCarlonApear() //确认卡伦是否复活
		switch WorkState {
		//初始化
		case Init:
			log.Infof("开始狩猎卡仑")
			FindCarlen = false
			HaveCarlen = false
			PosCount = 0
			Transition(TeleportMap)
			break
		//传送到姜饼城
		case TeleportMap:
			log.Infof("传送到姜饼城")
			if g.Configs.HuntConfig.CarryTeam {
				g.TeamGoToMap(gameTypes.MapId_GingerbreadCity.Uint32())
			} else {
				g.GoToMap(gameTypes.MapId_GingerbreadCity.Uint32())
			}

			time.Sleep(time.Millisecond * 1000)
			useFlyWing()
			time.Sleep(time.Millisecond * 3200)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
			Useskill()
			time.Sleep(time.Millisecond * 1000)
			Transition(MOVE_RABBIDSPOS)
			break
		//移动到疯兔地点
		case MOVE_RABBIDSPOS:
			if int(PosCount) < len(RabbidsPos) {
				if g.MoveChartWait(RabbidsPos[PosCount]) {
					log.Infof("抵达疯兔出生点%d", PosCount+1)
					Transition(CHECK_RABBIDS)
				}
			} else {
				if HaveCarlen {
					log.Infof("疯兔狩猎完成，开始狩猎卡仑")
					Transition(MOVE_CARLENPOS)
				} else {
					if _, ok := HiddenMVPList["卡仑"]; ok {
						obj := HiddenMVPList["卡仑"]
						obj.RespawnTime = time.Now().Add(10 * time.Minute)
						HiddenMVPList["卡仑"] = obj
					} else {
						log.Error("卡仑不在列表中")
					}
					log.Infof("未搜寻到卡仑，重新查找时间:%s", HiddenMVPList["卡仑"].RespawnTime.Format("2006-01-02 15:04:05"))
					Transition(Init)
					return
				}
			}
			break
		//确认有疯兔
		case CHECK_RABBIDS:
			if g.IsMonsterInRange("疯兔") {
				log.Infof("出生点%d发现疯兔，开始狩猎疯兔", PosCount+1)
				HaveCarlen = true
				fightMonstStar("疯兔")
				Transition(HUNT_RABBIDS)
			} else {
				log.Infof("出生点%d未发现疯兔，到疯兔点%d", PosCount+1, PosCount+2)
				fightCancel()
				PosCount++
				Transition(MOVE_RABBIDSPOS)
			}
			break
		//狩猎疯兔
		case HUNT_RABBIDS:
			if !g.IsMonsterInRange("疯兔") {
				log.Infof("出生点%d疯兔狩猎完成，到疯兔点%d", PosCount+1, PosCount+2)
				PosCount++
				fightCancel()
				Transition(MOVE_RABBIDSPOS)
			}
			break
		//移动到卡伦出生位置
		case MOVE_CARLENPOS:
			if g.MoveChartWait(g.ParsePos(-29854, 62, -87337)) {
				Transition(CHECK_RABBIDS)
			}
			break
		//确认卡伦
		case CHECK_CARLEN:
			if g.IsMonsterInRange("卡仑") {
				fightMonstStar("卡仑")
				Transition(HUNT_RABBIDS)
			} else {
				fightCancel()
				PosCount++
				Transition(MOVE_RABBIDSPOS)
			}
			break
		//狩猎卡伦
		case HUNT_CARLEN:
			if g.IsMonsterInRange("卡仑") {
				TargetID := g.AtkStat.GetCurrentTargetId()
				MapNPC := g.GetMapNpcs()
				if TargetID != 0 && MapNPC[TargetID].Attrs != nil {
					MonsterHP := utils.GetNpcAttrValByType(MapNPC[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
					MHP := utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
					if MonsterHP != tempMHP || MHP != tempUHP {
						tempMHP = MonsterHP
						tempUHP = MHP
						log.Infof("%s 未死亡，剩余血量:%d", "卡仑", MonsterHP)
						log.Infof("我的血量", MHP)
					} else if MonsterHP == 0 {
						Transition(End)
					}
				}
			} else {
				Transition(End)
			}
			break
		//结束
		case End:
			if _, ok := HiddenMVPList["卡仑"]; ok {
				obj := HiddenMVPList["卡仑"]
				obj.RespawnTime = time.Now().Add(30 * time.Minute)
				HiddenMVPList["卡仑"] = obj
			} else {
				log.Error("卡仑不在列表中")
			}
			fightCancel()
			log.Infof("卡仑已死亡，复活时间:%s", HiddenMVPList["卡仑"].RespawnTime.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Millisecond * 3000)
			Transition(Init)
			return
		}

		time.Sleep(time.Millisecond * 100)
	}
}
func CheckCarlonApear() bool {
	if !FindCarlen && g.IsMonsterInRange("卡仑") {
		FindCarlen = true
		log.Infof("发现卡仑，开始狩猎")
		fightCancel()
		fightMonstStar("卡仑")
		Transition(HUNT_CARLEN)
		return true
	} else {
		return false
	}
}

func Transition(SwitchState Work) {
	if WorkState != SwitchState {
		WorkState = SwitchState
	}
}
