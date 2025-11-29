package AutoBossHunting

import (
	"time"

	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

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

func i32(v int32) *int32 { return &v }

var RabbidsPos = []Cmd.ScenePos{
	{X: i32(10165), Y: i32(27), Z: i32(-57086)},
	{X: i32(20845), Y: i32(27), Z: i32(-85483)},
	{X: i32(-29854), Y: i32(27), Z: i32(-84827)},
	{X: i32(-41499), Y: i32(127), Z: i32(-51326)},
}

func (b *BossHuntTask) huntCarlen() {
	for {
		b.checkCarlonApear() // 确认卡伦是否复活
		switch b.workState {
		// 初始化
		case Init:
			log.Infof("开始狩猎卡仑")
			b.findCarlen = false
			b.haveCarlen = false
			b.posCount = 0
			b.transition(TeleportMap)
			break
		// 传送到姜饼城
		case TeleportMap:
			b.logger.Infof("传送到姜饼城")
			if b.GC.Configs.HuntConfig.CarryTeam {
				b.GC.TeamGoToMap(gameTypes.MapId_GingerbreadCity.Uint32())
			} else {
				b.GC.GoToMap(gameTypes.MapId_GingerbreadCity.Uint32())
			}

			time.Sleep(time.Millisecond * 1000)
			b.useFlyWing()
			time.Sleep(time.Millisecond * 3200)
			b.useSkill()
			time.Sleep(time.Millisecond * 1000)
			b.useSkill()
			time.Sleep(time.Millisecond * 1000)
			b.transition(MOVE_RABBIDSPOS)
			break
		// 移动到疯兔地点
		case MOVE_RABBIDSPOS:
			if int(b.posCount) < len(RabbidsPos) {
				if b.GC.MoveChartWait(RabbidsPos[b.posCount]) {
					log.Infof("抵达疯兔出生点%d", b.posCount+1)
					b.transition(CHECK_RABBIDS)
				}
			} else {
				if b.haveCarlen {
					log.Infof("疯兔狩猎完成，开始狩猎卡仑")
					b.transition(MOVE_CARLENPOS)
				} else {
					if _, ok := b.hiddenMVPList["卡仑"]; ok {
						obj := b.hiddenMVPList["卡仑"]
						obj.RespawnTime = time.Now().Add(10 * time.Minute)
						b.hiddenMVPList["卡仑"] = obj
					} else {
						log.Error("卡仑不在列表中")
					}
					log.Infof("未搜寻到卡仑，重新查找时间:%s", b.hiddenMVPList["卡仑"].RespawnTime.Format("2006-01-02 15:04:05"))
					b.transition(Init)
					return
				}
			}
			break
		// 确认有疯兔
		case CHECK_RABBIDS:
			if b.GC.IsMonsterInRange("疯兔") {
				b.logger.Infof("出生点%d发现疯兔，开始狩猎疯兔", b.posCount+1)
				b.haveCarlen = true
				b.fightMonstStar("疯兔", 0)
				b.transition(HUNT_RABBIDS)
			} else {
				b.logger.Infof("出生点%d未发现疯兔，到疯兔点%d", b.posCount+1, b.posCount+2)
				b.fightCancel()
				b.posCount++
				b.transition(MOVE_RABBIDSPOS)
			}
			break
		// 狩猎疯兔
		case HUNT_RABBIDS:
			if !b.GC.IsMonsterInRange("疯兔") {
				b.logger.Infof("出生点%d疯兔狩猎完成，到疯兔点%d", b.posCount+1, b.posCount+2)
				b.posCount++
				b.fightCancel()
				b.transition(MOVE_RABBIDSPOS)
			}
			break
		// 移动到卡伦出生位置
		case MOVE_CARLENPOS:
			if b.GC.MoveChartWait(b.GC.ParsePos(-29854, 62, -87337)) {
				b.transition(CHECK_RABBIDS)
			}
			break
		// 确认卡伦
		case CHECK_CARLEN:
			if b.GC.IsMonsterInRange("卡仑") {
				b.fightMonstStar("卡仑", 0)
				b.transition(HUNT_CARLEN)
			} else {
				b.fightCancel()
				b.posCount++
				b.transition(MOVE_RABBIDSPOS)
			}
			break
		// 狩猎卡伦
		case HUNT_CARLEN:
			if b.GC.IsMonsterInRange("卡仑") {
				TargetID := b.GC.AtkStat.GetCurrentTargetId()
				MapNPC := b.GC.GetMapNpcs()
				if TargetID != 0 && MapNPC[TargetID].Attrs != nil {
					MonsterHP := utils.GetNpcAttrValByType(MapNPC[TargetID].Attrs, Cmd.EAttrType_EATTRTYPE_HP)
					MHP := utils.GetNpcAttrValByType(b.GC.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
					if MonsterHP != b.tempMHP || MHP != b.tempUHP {
						b.tempMHP = MonsterHP
						b.tempUHP = MHP
						b.logger.Infof("%s 未死亡，剩余血量:%d", "卡仑", MonsterHP)
						b.logger.Infof("我的血量:%d", MHP)
					} else if MonsterHP == 0 {
						b.transition(End)
					}
				}
			} else {
				b.transition(End)
			}
			break
		// 结束
		case End:
			if _, ok := b.hiddenMVPList["卡仑"]; ok {
				obj := b.hiddenMVPList["卡仑"]
				obj.RespawnTime = time.Now().Add(30 * time.Minute)
				b.hiddenMVPList["卡仑"] = obj
			} else {
				b.logger.Error("卡仑不在列表中")
			}
			b.fightCancel()
			b.logger.Infof("卡仑已死亡，复活时间:%s", b.hiddenMVPList["卡仑"].RespawnTime.Format("2006-01-02 15:04:05"))
			time.Sleep(time.Millisecond * 3000)
			b.transition(Init)
			return
		}

		time.Sleep(time.Millisecond * 100)
	}
}
func (b *BossHuntTask) checkCarlonApear() bool {
	if !b.findCarlen && b.GC.IsMonsterInRange("卡仑") {
		b.findCarlen = true
		b.logger.Infof("发现卡仑，开始狩猎")
		b.fightCancel()
		b.fightMonstStar("卡仑", 0)
		b.transition(HUNT_CARLEN)
		return true
	} else {
		return false
	}
}

func (b *BossHuntTask) transition(SwitchState Work) {
	if b.workState != SwitchState {
		b.workState = SwitchState
	}
}

func (b *BossHuntTask) useElementStone() {
	item := b.GC.FindPackItemByName("风灵原石", Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		b.logger.Warnf("风灵原石没有找到")
	} else {
		b.logger.Infof("使用风灵原石")
		b.GC.UseItem(item.GetBase().GetGuid(), 1)
		time.Sleep(time.Millisecond * 1000)
	}
}

func (b *BossHuntTask) useElementArrow(arrowType gameTypes.ElementArrowType) {
	item := b.GC.FindPackItemByName(string(arrowType), Cmd.EPackType_EPACKTYPE_MAIN)
	if item == nil {
		b.logger.Warnf("%s没有找到", arrowType)
	} else {
		if item.GetBase().GetIsactive() {
			b.logger.Infof("%s已装备", arrowType)
			return
		}
		b.logger.Infof("使用%s", arrowType)
		b.GC.UseItem(item.GetBase().GetGuid(), 0)
		time.Sleep(time.Millisecond * 1000)
	}
}
