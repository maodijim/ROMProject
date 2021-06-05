package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math"
	"sort"
	"strconv"
	"time"
)

var (
	CampsFriend        = "Friend"
	CampsEnemy         = "Enemy"
	DefaultTargetRange = float64(9999)
)

func (g *GameConnection) SkillCmd(skillId uint32, data *Cmd.PhaseData, random1 bool) {
	skillItem := g.SkillItems[skillId]
	if skillItem.NameZh != "普通攻击" {
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

func (g *GameConnection) AttackTarget(skillId uint32, target *Cmd.MapNpc) {
	if g.MapNpcs[target.GetId()] != nil {
		hitType := int32(2)
		damage := int32(1)
		hitTargets := []*Cmd.HitedTarget{
			&Cmd.HitedTarget{
				Charid: target.Id,
				Type:   &hitType,
				Damage: &damage,
			},
		}
		if g.SkillItems[skillId].Range != "" {
			DmgRange, _ := strconv.ParseFloat(g.SkillItems[skillId].Range, 64)
			targetDict, targetRange := g.GetTargetByRange([]string{"all"}, target.GetPos(), DmgRange)
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
		dir := int32(utils.CalcDir(utils.GetAngleByAxisY(g.Role.Pos, target.GetPos())))
		pData := &Cmd.PhaseData{
			Number:       &num,
			Pos:          target.GetPos(),
			HitedTargets: hitTargets,
			Dir:          &dir,
		}
		var delay float64
		if g.SkillItems[skillId].NameZh == "普通攻击" {
			delay = 1 / (float64(g.GetAtkSpd()) / 1000)
			//delay = 1
		} else {
			delay, _ = strconv.ParseFloat(g.SkillItems[skillId].DelayCd, 64)
		}
		cd, _ := strconv.ParseFloat(g.SkillItems[skillId].CD, 64)
		if cd > delay {
			g.Role.CDs[skillId] = time.Now().Add(time.Duration(cd) * time.Second)
		}
		g.SkillCmd(skillId, pData, false)
		maxDelay := math.Max(delay, 0.25)
		time.Sleep(time.Duration(maxDelay*1000) * time.Millisecond)
	}
}

func (g *GameConnection) GetTargetByRange(monsterName []string, srcPos *Cmd.ScenePos, targetRange float64) (distDict map[float64]uint64, distanceList []float64) {
	distDict = map[float64]uint64{}
	g.Mutex.RLock()
	for _, npc := range g.MapNpcs {
		if npc.GetOwner() != 0 {
			continue
		}
		if (utils.StrSliceContain(monsterName, "all") || utils.StrSliceContain(monsterName, npc.GetName())) && len(npc.GetAttrs()) != 1 {
			distance := utils.GetDistanceXZ(srcPos, npc.GetPos())
			if distance <= targetRange*utils.Scale {
				distanceList = append(distanceList, distance)
				distDict[distance] = npc.GetId()
			}
		}
	}
	g.Mutex.RUnlock()
	sort.Float64s(distanceList)
	return distDict, distanceList
}

func (g *GameConnection) copyTarget(org *Cmd.MapNpc) *Cmd.MapNpc {
	orgX := org.GetPos().GetX()
	orgY := org.GetPos().GetY()
	orgZ := org.GetPos().GetZ()
	orgId := org.GetId()
	orgName := org.GetName()
	target := &Cmd.MapNpc{
		Pos: &Cmd.ScenePos{
			X: &orgX,
			Y: &orgY,
			Z: &orgZ,
		},
		Id:   &orgId,
		Name: &orgName,
	}
	return target
}

func (g *GameConnection) AttackClosestByName(skillId uint32, monsterName []string) {
	distDict, distanceList := g.GetTargetByRange(monsterName, g.Role.Pos, DefaultTargetRange)
	if len(distanceList) > 0 {
		distance := distanceList[0]
		closestId := distDict[distanceList[0]]
		target := g.copyTarget(g.MapNpcs[closestId])
		skillRange, _ := strconv.ParseFloat(g.SkillItems[skillId].LaunchRange, 64)
		launchSkillDis := skillRange * utils.Scale
		launchSkillPos := utils.GetPosAwayFromTarget(g.Role.Pos, target.GetPos(), launchSkillDis)
		targetDis := utils.GetDistanceXZ(g.Role.Pos, target.GetPos())
		lastPrint := time.Now().Add(-5 * time.Second)

		if targetDis >= launchSkillDis {
			g.MoveChart(launchSkillPos)
			staleMoveCount := 0
			x := g.Role.Pos.GetX()
			y := g.Role.Pos.GetY()
			z := g.Role.Pos.GetZ()
			prePos := Cmd.ScenePos{
				X: &x,
				Y: &y,
				Z: &z,
			}
		moveToTargetLoop:
			for {
				select {
				case <-g.quit:
					return
				default:
					if staleMoveCount > 10 {
						break moveToTargetLoop
					}
					if g.MapNpcs[closestId] != nil && distance <= launchSkillDis {
						break moveToTargetLoop
					} else {
						if _, ok := g.MapNpcs[closestId]; !ok {
							log.Infof("没有找到怪物id %d %s", closestId, monsterName)
							return
						}
						distance = utils.GetDistanceXZ(g.Role.Pos, target.GetPos())
						if time.Since(lastPrint) > time.Second*5 {
							lastPrint = time.Now()
							log.Infof("%s 跑路中 怪物id: %d 名字: %s 血量: %d 位置: %v 距离: %f 攻击距离 %f",
								g.Role.GetRoleName(), closestId, target.GetName(), utils.GetNpcAttrValByType(target.GetAttrs(), Cmd.EAttrType_EATTRTYPE_HP), target.GetPos(), distance, launchSkillDis)
						}
						target = g.copyTarget(g.MapNpcs[closestId])
						launchSkillPos = utils.GetPosAwayFromTarget(g.Role.Pos, target.GetPos(), launchSkillDis)
						if staleMoveCount > 5 {
							g.MoveChart(target.GetPos())
						} else if distance <= launchSkillDis {
							g.MoveChart(launchSkillPos)
						}
						time.Sleep(300 * time.Millisecond)
					}
					if prePos.GetX() == g.Role.Pos.GetX() && prePos.GetY() == g.Role.Pos.GetY() {
						staleMoveCount += 1
					} else {
						staleMoveCount = 0
						x = g.Role.Pos.GetX()
						y = g.Role.Pos.GetY()
						z = g.Role.Pos.GetZ()
						prePos = Cmd.ScenePos{
							X: &x,
							Y: &y,
							Z: &z,
						}
					}
				}
			}
		}
		g.AttackTarget(skillId, target)
	} else {
		time.Sleep(200 * time.Millisecond)
	}
}

func (g *GameConnection) EnableAutoAttack(monsterList []string, disable chan *bool) {
	go func() {
		for {
			select {
			case <-disable:
				return
			case <-g.quit:
				return
			default:
				autoSkills := g.GetAutoSkills()
				for _, skill := range autoSkills {
					select {
					case <-disable:
						return
					case <-g.quit:
						return
					default:
						skillItem := g.SkillItems[skill.GetId()]
						log.Debugf("自动技能位置: %d, 技能id: %d, 技能名字: %s",
							skill.GetShortcuts()[len(skill.GetShortcuts())-1].GetPos(), skill.GetId(), skillItem.NameZh)
						cd, _ := strconv.ParseFloat(skillItem.CD, 64)
						if time.Since(g.Role.CDs[skill.GetId()]) < time.Duration(cd) {
							log.Infof("技能CD中:%s", skillItem.NameZh)
							continue
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
											if curHpPer > g.GetHpPer() || g.GetBuffNameByRegex("原地休息") == "" || (g.GetHpPer() > 0.95 && g.GetSpPer() > 0.95) {
												break
											}
											time.Sleep(6 * time.Second)
										}
									}
								}
							} else if buff != "" {
								log.Debugf("找到技能buff: %s -> %s", skillItem.NameZh, buff)
								continue
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
							//这是攻击技能
							g.AttackClosestByName(skill.GetId(), monsterList)
						}
						time.Sleep(150 * time.Millisecond)
					}
				}
			}
		}
	}()
}
