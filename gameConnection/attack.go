package gameConnection

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/mohae/deepcopy"
	log "github.com/sirupsen/logrus"
)

var (
	CampsFriend        = "Friend"
	CampsEnemy         = "Enemy"
	DefaultTargetRange = float64(9999)
	attackLogic        = map[string]string{
		"SkillLockedTarget": "SkillLockedTarget",
		"SkillPointRange":   "SkillPointRange",
		"SkillSelfRange":    "SkillSelfRange",
		"SkillNone":         "SkillNone",
		"SkillForwardRect":  "SkillForwardRect",
	}
	lastPrint = time.Now()
	lastMove  = time.Now()
)

type AttackMonsterStat struct {
	lastAttack time.Time
	lock       sync.RWMutex
	Standstill bool
}

func (a *AttackMonsterStat) IsStandstill() bool {
	return a.Standstill
}

func (a *AttackMonsterStat) SetStandstill(standstill bool) {
	a.Standstill = standstill
}

func (a *AttackMonsterStat) GetLastAttack() time.Time {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.lastAttack
}

func (a *AttackMonsterStat) SetLastAttack(t time.Time) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.lastAttack = t
}

func (g *GameConnection) SkillCmd(skillId uint32, data *Cmd.PhaseData, random1 bool) {
	if skillItem, ok := g.SkillItems[skillId]; ok && skillItem.NameZh != "普通攻击" {
		log.Infof("%s 释放技能 %d %s", g.Role.GetRoleName(), skillId, skillItem.NameZh)
	}
	random := uint32(1)
	if !random1 {
		random = uint32(utils.GetRandom(0, 100))
	}
	cmd := &Cmd.SkillBroadcastUserCmd{
		Charid:  g.Role.RoleId,
		SkillID: &skillId,
		Data:    data,
		Random:  &random,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["SKILL_BROADCAST_USER_CMD"],
	)
}

func (g *GameConnection) CalDmgTargets() []*Cmd.HitedTarget {
	return nil
}

func (g *GameConnection) AttackTarget(skillId uint32, target Cmd.MapNpc) {
	skillItem := g.SkillItems[skillId]
	g.Mutex.Lock()
	if g.MapNpcs[target.GetId()] == nil {
		g.Mutex.Unlock()
		return
	}
	g.Mutex.Unlock()

	hitType := int32(2)
	damage := int32(1)
	hitTargets := []*Cmd.HitedTarget{
		&Cmd.HitedTarget{
			Charid: target.Id,
			Type:   &hitType,
			Damage: &damage,
		},
	}
	if skillItem.Range != "" && skillItem.Logic == attackLogic["SkillLockedTarget"] {
		DmgRange, _ := strconv.ParseFloat(skillItem.Range, 64)
		targetDict, targetRange := g.GetTargetByRange([]string{"all"}, *target.GetPos(), DmgRange)
		for _, t := range targetRange {
			newTarget := targetDict[t]
			newHitedTarget := &Cmd.HitedTarget{
				Charid: &newTarget,
				Type:   &hitType,
				Damage: &damage,
			}
			hitTargets = append(hitTargets, newHitedTarget)
		}
	}
	num := int32(1)
	// dir := int32(utils.CalcDir(utils.GetAngleByAxisY(g.Role.Pos, target.GetPos())))
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    target.GetPos(),
		// Dir:    &dir,
	}
	if skillItem.Logic == attackLogic["SkillLockedTarget"] {
		pData.HitedTargets = hitTargets
	} else if skillItem.Logic == attackLogic["SkillPointRange"] {
		num = int32(0)
		pData.Number = &num
	}

	// Calculate Skill Delay & CD
	var delay float64
	if skillItem.NameZh == "普通攻击" {
		delay = 1 / (float64(g.GetAtkSpd()) / 1000 * (1 + float64(g.getAtkSpdPer())/1000))
		// delay = 1
	} else {
		delay, _ = strconv.ParseFloat(skillItem.DelayCd, 64)
	}
	cd, _ := strconv.ParseFloat(skillItem.CD, 64)
	if cd > delay {
		g.Role.SetSkillCd(skillId, time.Now().Add(time.Duration(cd)*time.Second))
	}
	maxDelay := math.Max(delay, 0.1)
	if time.Since(g.AtkStat.GetLastAttack()) >= time.Duration(maxDelay*500)*time.Millisecond {
		g.SkillCmd(skillId, pData, false)
		g.AtkStat.SetLastAttack(time.Now())
	}
}

func (g *GameConnection) GetTargetByRange(monsterName []string, srcPos Cmd.ScenePos, targetRange float64) (distDict map[float64]uint64, distanceList []float64) {
	distDict = map[float64]uint64{}
	g.Mutex.RLock()
	mapNpcs := deepcopy.Copy(g.MapNpcs).(map[uint64]*Cmd.MapNpc)
	g.Mutex.RUnlock()
	for _, npc := range mapNpcs {
		if npc.GetOwner() != 0 {
			continue
		}
		// This is not a monster
		if npc.GetId() < 10000 {
			continue
		}
		if (utils.StrSliceContain(monsterName, "all") || utils.StrSliceContain(monsterName, npc.GetName())) && len(npc.GetAttrs()) != 1 {
			if npc.GetPos() == nil {
				continue
			}
			distance := utils.GetDistanceXYZ(srcPos, *npc.GetPos())
			if distance <= targetRange*utils.AtkRangeScale {
				distanceList = append(distanceList, distance)
				distDict[distance] = npc.GetId()
			}
		}
	}
	sort.Float64s(distanceList)
	return distDict, distanceList
}

func (g *GameConnection) copyTarget(org *Cmd.MapNpc) *Cmd.MapNpc {
	g.Mutex.RLock()
	target := deepcopy.Copy(org).(*Cmd.MapNpc)
	g.Mutex.RUnlock()
	return target
}

func (g *GameConnection) AttackClosestByName(skillId uint32, monsterName []string) {
	distDict, distanceList := g.GetTargetByRange(monsterName, g.Role.GetPos(), DefaultTargetRange)
	if len(distanceList) > 0 {
		distance := distanceList[0]
		closestId := distDict[distanceList[0]]
		target, ok := g.GetMapNpcs()[closestId]
		if !ok {
			return
		}
		skillRange := g.GetAttackRange(skillId)
		launchSkillDis := skillRange * utils.AtkRangeScale
		var launchSkillPos Cmd.ScenePos
		if launchSkillDis <= 1500 {
			launchSkillPos = *target.GetPos()
		} else {
			launchSkillPos = utils.GetPosAwayFromTarget(g.Role.GetPos(), *target.GetPos(), launchSkillDis)
		}
		targetDis := utils.GetDistanceXZ(g.Role.GetPos(), *target.GetPos())

		if targetDis >= launchSkillDis {
			if g.AtkStat.IsStandstill() {
				log.Warnf(
					"attack mode is standstill but monster %s distance is %f greater than range %f, skip attack",
					target.GetName(),
					targetDis,
					launchSkillDis,
				)
				return
			} else if time.Since(lastMove) > time.Millisecond*150 {
				lastMove = time.Now()
				g.MoveChart(launchSkillPos)
			}
		moveToTargetLoop:
			for {
				select {
				case <-time.After(100 * time.Millisecond):
					// oldDistance := distance
					target, ok = g.GetMapNpcs()[closestId]
					if !ok {
						break moveToTargetLoop
					}
					if time.Since(lastPrint) > time.Second*5 {
						lastPrint = time.Now()
						log.Infof("%s 跑路中 怪物id: %d 名字: %s 血量: %d 位置: %v 距离: %f 攻击距离 %f 角度 %f",
							g.Role.GetRoleName(),
							closestId,
							target.GetName(),
							utils.GetNpcAttrValByType(target.GetAttrs(), Cmd.EAttrType_EATTRTYPE_HP),
							target.GetPos(),
							distance,
							launchSkillDis,
							utils.GetAngleByAxisY(g.Role.GetPos(), *target.GetPos()),
						)
					}
					if g.GetMapNpcs()[closestId].Id == nil {
						log.Warnf("target %s is dead, skip attack", target.GetName())
						return
					} else if distance <= launchSkillDis {
						break moveToTargetLoop
					}
					break moveToTargetLoop
				}
			}
		} else {
			g.AttackTarget(skillId, target)
			if g.GetMapNpcs()[closestId].Id == nil {
				log.Warnf("target %s is killed", target.GetName())
				return
			}
		}
	}
}

func (g *GameConnection) EnableAutoAttack(monsterList []string, disable chan bool) {
	go func() {
		ticker := time.NewTicker(time.Millisecond * 100)
		defer ticker.Stop()
		for {
			select {
			case <-disable:
				return
			case <-g.quit:
				return
			default:
				autoSkills := g.GetAutoSkills()
			skillLoop:
				for _, skill := range autoSkills {
					select {
					case <-disable:
						return
					case <-g.quit:
						return
					case <-ticker.C:
						skillItem := g.SkillItems[skill.GetId()]
						log.Debugf("自动技能位置: %d, 技能id: %d, 技能名字: %s",
							skill.GetShortcuts()[len(skill.GetShortcuts())-1].GetPos(), skill.GetId(), skillItem.NameZh)
						cd, _ := strconv.ParseFloat(skillItem.CD, 64)
						if time.Since(g.Role.GetSkillCd(skill.GetId())) < time.Duration(cd) {
							log.Infof("技能CD中:%s", skillItem.NameZh)
							continue skillLoop
						}
						// 这是buff
						if skillItem.Camps == CampsFriend {
							buff := g.GetBuffNameByRegex(fmt.Sprintf("%s.*", skillItem.NameZh))
							if skillItem.NameZh == "装死" {
								maxHp := g.GetMaxHp()
								curHpPer := g.GetHpPer()
								per := 0.1
								if maxHp > 0 && g.GetSpPer() < per || g.GetHpPer() < per {
									num := int32(1)
									dir := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
									pData := &Cmd.PhaseData{
										Number: &num,
										Pos:    g.Role.Pos,
										Dir:    &dir,
									}
									g.SkillCmd(skill.GetId(), pData, true)
									// block other action until 装死 ended'
									lastPrint := time.Now().Add(10 * time.Second)
									for startTime := time.Now(); time.Since(startTime) < 50*time.Second; {
										select {
										case <-disable:
											return
										case <-g.quit:
											return
										default:
											if time.Since(lastPrint) > 10*time.Second {
												lastPrint = time.Now()
												log.Infof("%s 装死中 血量:%d SP:%d", g.Role.GetRoleName(), g.GetCurrentHp(), g.GetCurrentSp())
											}
											if curHpPer > g.GetHpPer() ||
												g.GetBuffNameByRegex("原地休息") == "" ||
												(g.GetHpPer() > 0.95 && g.GetSpPer() > 0.95) {
												break
											}
											time.Sleep(5 * time.Second)
										}
									}
								}
							} else if skillItem.SkillType == "Heal" {
								if g.GetHpPer() > 0.65 {
									continue skillLoop
								}
							} else if skillItem.SkillType == "Reborn" {
								continue skillLoop
							} else if buff != "" {
								log.Debugf("找到技能buff: %s -> %s", skillItem.NameZh, buff)
								continue skillLoop
							} else {
								log.Debugf("没有找到技能buff %s", skillItem.NameZh)
								num := int32(1)
								dir := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DIR))
								pData := &Cmd.PhaseData{
									Number: &num,
									Pos:    g.Role.Pos,
									Dir:    &dir,
								}
								g.SkillCmd(skill.GetId(), pData, true)
								var delay float64
								if skillItem.DelayCd != "" {
									delay, _ = strconv.ParseFloat(skillItem.DelayCd, 64)
								} else {
									delay = 0
								}
								time.Sleep(time.Duration(math.Max(delay, 0.1)*1000) * time.Millisecond)
							}
						} else if skillItem.Camps == CampsEnemy {
							// 这是攻击技能
							g.AttackClosestByName(skill.GetId(), monsterList)
						}
					}
				}
			}
		}
	}()
}

func (g *GameConnection) GetAttackRange(skillId uint32) (atkRange float64) {
	skillItem := g.SkillItems[skillId]
	atkRange, _ = strconv.ParseFloat(skillItem.LaunchRange, 64)
	if skillItem.NameZh == "普通攻击" {
		// 无限星辰
		if g.Role.Buffs[131080] != nil {
			atkRange += float64(g.Role.SkillItems[13234].GetExtralv()) * 0.1
		}
	}
	atkPer := utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_ATKDISTANCEPER)
	if atkPer > 0 {
		atkRange = atkRange * (1 + float64(atkPer)/1000)
	}
	return atkRange
}

func (g *GameConnection) AddAttrPoint(s, a, v, i, d, l uint32) ([]int32, error) {
	pt := g.Role.GetTotalPoint()
	sNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_STRPOINT))
	aNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_AGIPOINT))
	vNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_VITPOINT))
	iNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_INTPOINT))
	dNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DEXPOINT))
	lNow := int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_LUKPOINT))
	attrs := []int32{sNow, aNow, vNow, iNow, dNow, lNow}
	errMsg := "Not enough point for %d %s need %d more point"
	if s > 0 {
		if int32(s)*utils.GetAttrPointReq(sNow) > pt {
			msg := fmt.Sprintf(errMsg, s, "strenth", int32(s)*utils.GetAttrPointReq(sNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	if a > 0 {
		if int32(a)+utils.GetAttrPointReq(aNow) > pt {
			msg := fmt.Sprintf(errMsg, a, "agility", int32(a)*utils.GetAttrPointReq(aNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	if v > 0 {
		if int32(v)+utils.GetAttrPointReq(vNow) > pt {
			msg := fmt.Sprintf(errMsg, v, "vitality", int32(v)*utils.GetAttrPointReq(vNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	if i > 0 {
		if int32(i)+utils.GetAttrPointReq(iNow) > pt {
			msg := fmt.Sprintf(errMsg, i, "intelligence", int32(i)*utils.GetAttrPointReq(iNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	if d > 0 {
		if int32(d)+utils.GetAttrPointReq(dNow) > pt {
			msg := fmt.Sprintf(errMsg, d, "dexterity", int32(d)*utils.GetAttrPointReq(dNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	if l > 0 {
		if int32(l)+utils.GetAttrPointReq(lNow) > pt {
			msg := fmt.Sprintf(errMsg, l, "luck", int32(l)*utils.GetAttrPointReq(lNow)-pt)
			return attrs, errors.New(msg)
		}
	}
	attType := Cmd.PointType_POINTTYPE_ADD
	cmd := &Cmd.AddAttrPoint{
		Type:     &attType,
		Strpoint: &s,
		Agipoint: &a,
		Vitpoint: &v,
		Intpoint: &i,
		Dexpoint: &d,
		Lukpoint: &l,
	}
	g.AddNotifier("AddAttrPoint")
	_ = g.sendProtoCmd(
		cmd,
		sceneUser2CmdId,
		Cmd.User2Param_value["USER2PARAM_ADDATTRPOINT"],
	)
	<-g.notifier["AddAttrPoint"]
	g.removeNotifier("AddAttrPoint")
	attrs = []int32{
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_STRPOINT)),
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_AGIPOINT)),
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_VITPOINT)),
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_INTPOINT)),
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_DEXPOINT)),
		int32(utils.GetNpcDataValByType(g.Role.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_LUKPOINT)),
	}
	return attrs, nil
}

func (g *GameConnection) GetBuffByName(name string) utils.BuffItem {
	g.Role.Mutex.RLock()
	defer g.Role.Mutex.RUnlock()
	for _, v := range g.Role.Buffs {
		if name == g.BuffItems[v.GetId()].BuffName {
			return g.BuffItems[v.GetId()]
		}
	}
	return utils.BuffItem{}
}
