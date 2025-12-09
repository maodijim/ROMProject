package gameConnection

import (
	"context"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	notifier "ROMProject/gameConnection/types"
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
)

type AttackMonsterStat struct {
	lastAttack      time.Time
	lock            sync.RWMutex
	Standstill      bool
	CurrentTargetId uint64
	IsAutoAttacking bool
}

type TargetScore struct {
	Id      uint64
	Dist2   int64 // 距離平方（不用 sqrt）
	Density int   // 半徑 r 內怪物數
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

func (a *AttackMonsterStat) SetCurrentTargetId(targetId uint64) {
	a.lock.Lock()
	defer a.lock.Unlock()
	a.CurrentTargetId = targetId
}

func (a *AttackMonsterStat) GetCurrentTargetId() uint64 {
	a.lock.RLock()
	defer a.lock.RUnlock()
	return a.CurrentTargetId
}

func (g *GameConnection) SkillCmd(skillId uint32, data *Cmd.PhaseData, random1 bool) {
	if skillItem, ok := g.SkillItems[skillId]; ok && skillItem.NameZh != "普通攻击" && skillItem.NameZh != "扩散攻击" {
		log.Infof("%s 释放技能 %d %s", g.Role.GetRoleName(), skillId, skillItem.NameZh)
	}
	random := uint32(1)
	if !random1 {
		random = uint32(utils.GetRandom(0, 100))
	}
	cmd := &Cmd.SkillBroadcastUserCmd{
		Charid:  g.Role.RoleId,
		SkillID: &skillId,
		Random:  &random,
	}
	if data != nil {
		cmd.Data = data
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
	// 判断技能范围伤害目标
	if skillItem.Range != "" && (skillItem.Logic == attackLogic["SkillLockedTarget"] || skillItem.Logic == attackLogic["SkillPointRange"]) {
		DmgRange, _ := strconv.ParseFloat(skillItem.Range, 64)
		targetDict, targetRange := g.GetTargetByRange([]string{"all"}, *target.GetPos(), DmgRange)
		for _, t := range targetRange {
			// 跳过第一个目标（已经包含在hitTargets中了）
			if targetDict[t] == target.GetId() {
				continue
			}
			newTarget := targetDict[t]
			newHitedTarget := &Cmd.HitedTarget{
				Charid: &newTarget,
				Type:   &hitType,
				Damage: &damage,
			}
			hitTargets = append(hitTargets, newHitedTarget)
			// 判断技能范围伤害目标数量是否超过上限
			if skillItem.GetRangeNum() > 0 && len(hitTargets) > skillItem.GetRangeNum() {
				break
			}
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
	g.AtkStat.SetCurrentTargetId(target.GetId())
}

// GetTargetByRange returns a map of distance to target ID and a sorted list of distances
func (g *GameConnection) GetTargetByRange(monsterList []string, srcPos Cmd.ScenePos, targetRange float64) (distDict map[float64]uint64, distanceList []float64) {
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
		if (utils.Contains(monsterList, "all") || utils.Contains(monsterList, npc.GetName())) && len(npc.GetAttrs()) != 1 {
			if npc.GetPos() == nil {
				continue
			}
			distance := utils.GetDistanceXZ(srcPos, *npc.GetPos())
			if distance <= targetRange*utils.AtkRangeScale {
				distanceList = append(distanceList, distance)
				distDict[distance] = npc.GetId()
			}
		}
	}
	sort.Float64s(distanceList)
	return distDict, distanceList
}
func DistSquaredXZ(a, b Cmd.ScenePos) int64 {
	dx := int64(a.GetX() - b.GetX())
	dz := int64(a.GetZ() - b.GetZ())
	return dx*dx + dz*dz
}

// Selects the best target by balancing distance and cluster density.
func (g *GameConnection) GetTargetByDensitySameReturn(
	monsterList []string,
	srcPos Cmd.ScenePos,
	densityRange float64, // ✅ 只用來算密集度
) (distDict map[float64]uint64, distanceList []float64) {

	distDict = make(map[float64]uint64)
	distanceList = make([]float64, 0)

	type blockInfo struct {
		Count     int
		MinDist2  float64
		NearestId uint64
	}

	blocks := make(map[[2]int]*blockInfo)

	srcX := float64(srcPos.GetX())
	srcZ := float64(srcPos.GetZ())
	blockSize := densityRange

	// ✅ 记录全场最近
	globalMinDist2 := math.MaxFloat64
	var globalNearestBlock *blockInfo

	// ✅ 记录所有目标（给 else 用）
	type allTarget struct {
		Dist2 float64
		Id    uint64
	}
	allTargets := make([]allTarget, 0)

	// ✅ 一、分区块 + 密集度 + 最近距离 + 全目标收集
	for _, npc := range g.GetMapNpcs() {

		if npc.GetOwner() != 0 || npc.GetId() < 10000 {
			continue
		}
		if !(utils.Contains(monsterList, "all") || utils.Contains(monsterList, npc.GetName())) {
			continue
		}
		if npc.GetPos() == nil || len(npc.GetAttrs()) == 1 {
			continue
		}

		pos := npc.GetPos()
		x := float64(pos.GetX())
		z := float64(pos.GetZ())

		dx := x - srcX
		dz := z - srcZ
		dist2 := dx*dx + dz*dz

		// ✅ 记录全部目标（给 else 排序用）
		allTargets = append(allTargets, allTarget{
			Dist2: dist2,
			Id:    npc.GetId(),
		})

		bx := int(math.Floor(x / blockSize))
		bz := int(math.Floor(z / blockSize))
		key := [2]int{bx, bz}

		if _, ok := blocks[key]; !ok {
			blocks[key] = &blockInfo{
				Count:    0,
				MinDist2: math.MaxFloat64,
			}
		}

		info := blocks[key]
		info.Count++

		// ✅ 区块内最近
		if dist2 < info.MinDist2 {
			info.MinDist2 = dist2
			info.NearestId = npc.GetId()
		}

		// ✅ 全场最近
		if dist2 < globalMinDist2 {
			globalMinDist2 = dist2
			globalNearestBlock = info
		}
	}

	// ✅ 二、找最密集区块
	var bestDenseBlock *blockInfo
	for _, b := range blocks {
		if bestDenseBlock == nil || b.Count > bestDenseBlock.Count {
			bestDenseBlock = b
		}
	}

	// ✅ 三、是否触发 2 倍密集度规则
	if bestDenseBlock != nil &&
		globalNearestBlock != nil &&
		bestDenseBlock.Count >= globalNearestBlock.Count*2 {

		// ✅ 仅回传「最密集区块的最近目标」
		finalDist := math.Sqrt(bestDenseBlock.MinDist2)
		distDict[finalDist] = bestDenseBlock.NearestId
		distanceList = append(distanceList, finalDist)

		return
	}

	// ✅ ✅ ✅ else：回传「附近所有目标 → 依距离排序」

	sort.Slice(allTargets, func(i, j int) bool {
		return allTargets[i].Dist2 < allTargets[j].Dist2
	})

	for _, t := range allTargets {
		dist := math.Sqrt(t.Dist2)
		distDict[dist] = t.Id
		distanceList = append(distanceList, dist)
	}

	return
}

func (g *GameConnection) IsMonsterInRange(monsterList ...string) bool {
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

		if utils.Contains(monsterList, "all") || (utils.Contains(monsterList, npc.GetName()) && len(npc.GetAttrs()) != 1) {
			return true
		}
	}
	return false
}

func (g *GameConnection) IsMonsterInDistance(distance int, monsterList ...string) bool {
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
		if utils.Contains(monsterList, npc.GetName()) && len(npc.GetAttrs()) != 1 {
			if utils.GetDistanceXYZ(
				g.Role.GetPos(),
				*npc.GetPos(),
			) > float64(distance) {
				return false
			}
			return true
		}
	}
	return false
}

func (g *GameConnection) copyTarget(org *Cmd.MapNpc) *Cmd.MapNpc {
	g.Mutex.RLock()
	target := deepcopy.Copy(org).(*Cmd.MapNpc)
	g.Mutex.RUnlock()
	return target
}

func (g *GameConnection) AttackClosestByName(skillId uint32, monsterName []string) {
	var (
		distDict     map[float64]uint64
		distanceList []float64
	)
	/*skillItem, ok := g.SkillItems[skillId]

	IsRangeSkill := ok && skillItem.Range != "" && (skillItem.Logic == attackLogic["SkillLockedTarget"] || skillItem.Logic == attackLogic["SkillPointRange"])

	if IsRangeSkill {
		distDict, distanceList = g.GetTargetByDensitySameReturn(monsterName, g.Role.GetPos(), 20000)
	} else {*/
	distDict, distanceList = g.GetTargetByRange(monsterName, g.Role.GetPos(), DefaultTargetRange)
	//}

	if len(distanceList) > 0 {
		distance := distanceList[0]
		closestId := distDict[distanceList[0]]
		target, ok := g.GetMapNpcs()[closestId]
		if !ok {
			g.AtkStat.SetCurrentTargetId(0)
			return
		}
		g.AtkStat.SetCurrentTargetId(target.GetId())
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
			} else {
				g.logger.Debugf("%s 跑向目标 怪物id: %d 名字: %s 血量: %d 位置: %v 距离: %f 攻击距离 %f 角度 %f 预计攻击位置: %v",
					g.Role.GetRoleName(),
					closestId,
					target.GetName(),
					utils.GetNpcAttrValByType(target.GetAttrs(), Cmd.EAttrType_EATTRTYPE_HP),
					target.GetPos(),
					distance,
					launchSkillDis,
					utils.GetAngleByAxisY(g.Role.GetPos(), *target.GetPos()),
					&launchSkillPos,
				)
				g.MoveChart(launchSkillPos)
			}
			lastPos := g.Role.GetPos()
			after := time.After(50 * time.Millisecond)
			check := time.NewTicker(100 * time.Millisecond)
			launchPosCheck := time.NewTicker(200 * time.Millisecond)
			newLaunchSkillDis := launchSkillDis
			defer check.Stop()
			defer launchPosCheck.Stop()
		moveToTargetLoop:
			for {
				select {
				case <-launchPosCheck.C:
					curPos := g.Role.GetPos()
					if lastPos.X == curPos.X && lastPos.Z == curPos.Z {
						// 卡住了
						newLaunchSkillDis = newLaunchSkillDis * 0.85
						if launchSkillDis < 2000 {
							g.MoveChart(*target.GetPos())
						} else {
							g.logger.Debugf("卡住了调整位置, 攻击距离: %f", newLaunchSkillDis)
							launchSkillPos = utils.GetPosAwayFromTarget(g.Role.GetPos(), *target.GetPos(), newLaunchSkillDis)
							g.MoveChart(launchSkillPos)
						}
					}
					lastPos = g.Role.GetPos()
				case <-check.C:
					target, ok = g.GetMapNpcs()[closestId]
					// 寻路时如果有更近的目标自动切换
					distDict, distanceList = g.GetTargetByRange(monsterName, g.Role.GetPos(), DefaultTargetRange)

					if len(distanceList) > 0 {
						closestId := distDict[distanceList[0]]
						newtarget, ok2 := g.GetMapNpcs()[closestId]
						if ok2 && newtarget.Id != target.Id {
							break moveToTargetLoop
						}
					}

					if !ok {
						break moveToTargetLoop
					}
					distance = utils.GetDistanceXYZ(g.Role.GetPos(), *target.GetPos())
					if distance <= launchSkillDis {
						check.Stop()
						break moveToTargetLoop
					}
				case <-after:
					// oldDistance := distance
					target, ok = g.GetMapNpcs()[closestId]
					if !ok {
						break moveToTargetLoop
					}
					if time.Since(lastPrint) > time.Second*5 {
						lastPrint = time.Now()
						g.logger.Infof("%s 跑路中 怪物id: %d 名字: %s 血量: %d 位置: %v 距离: %f 攻击距离 %f 角度 %f",
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
				}
			}
		} else {
			// 幽灵波利拍照显形
			if target.GetNpcID() == 20004 {
				attrs := target.GetAttrs()
				for _, a := range attrs {
					if a.GetType() == Cmd.EAttrType_EATTRTYPE_HIDE && a.GetValue() == 1 {
						g.logger.Infof("幽灵波利隐身中，拍照显形")
						g.TakePhotoSkill(&Cmd.CameraFocus{
							Targets: []uint64{target.GetId()},
						}, *target.GetPos(), []Cmd.MapNpc{target})
						time.Sleep(time.Millisecond * 1000)
					}
				}
			}
			g.AttackTarget(skillId, target)
			if g.GetMapNpcs()[closestId].Id == nil {
				log.Warnf("target %s is killed", target.GetName())
				g.AtkStat.SetCurrentTargetId(0)
				return
			}
		}
	}
}

func (g *GameConnection) EnableAutoAttack(ctx context.Context, monsterList ...string) {
	if g.AtkStat.IsAutoAttacking == true {
		log.Warnf("auto attack is already enabled")
		return
	}
	g.AtkStat.IsAutoAttacking = true
	var attackCtx context.Context
	attackCtx, g.cancelAtkCtx = context.WithCancel(ctx)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 75)

		defer func() {
			g.logger.Infof("stop auto attack")
			ticker.Stop()
			g.AtkStat.IsAutoAttacking = false
		}()
		for {
			select {
			case <-attackCtx.Done():
				g.logger.Infof("stop auto attack")
				g.AtkStat.SetCurrentTargetId(0)
				ticker.Stop()
				return
			case <-g.quit:
				g.logger.Infof("stop auto attack")
				g.AtkStat.SetCurrentTargetId(0)
				ticker.Stop()
				return
			default:
				autoSkills := g.GetAutoSkills()
			skillLoop:
				for _, skill := range autoSkills {
					select {
					case <-attackCtx.Done():
						log.Debugf("stop auto attack skill loop")
						g.AtkStat.SetCurrentTargetId(0)
						ticker.Stop()
						return
					case <-g.quit:
						log.Debugf("stop auto attack skill loop")
						g.AtkStat.SetCurrentTargetId(0)
						ticker.Stop()
						return
					case <-ticker.C:
						skillItem := g.SkillItems[skill.GetId()]
						log.Debugf("自动技能位置: %d, 技能id: %d, 技能名字: %s",
							skill.GetShortcuts()[len(skill.GetShortcuts())-1].GetPos(), skill.GetId(), skillItem.NameZh)
						cd, _ := strconv.ParseFloat(skillItem.CD, 64)
						if time.Since(g.Role.GetSkillCd(skill.GetId())) < time.Duration(cd) {
							if skill.GetId() != 50057001 {
								g.logger.Infof("技能CD中:%s", skillItem.NameZh)
							}
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
									lastPrint = time.Now().Add(10 * time.Second)
									for startTime := time.Now(); time.Since(startTime) < 50*time.Second; {
										select {
										case <-attackCtx.Done():
											ticker.Stop()
											return
										case <-g.quit:
											ticker.Stop()
											return
										default:
											if time.Since(lastPrint) > 10*time.Second {
												lastPrint = time.Now()
												g.logger.Infof("%s 装死中 血量:%d SP:%d", g.Role.GetRoleName(), g.GetCurrentHp(), g.GetCurrentSp())
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
							} else if skill.GetId() == 50057001 {
								// 这是备战精英
								if time.Since(g.Role.GetSkillCd(skill.GetId())) < time.Duration(cd) {
									log.Tracef("备战精英CD中:%s", skillItem.NameZh)
									continue skillLoop
								}
								g.SkillCmd(skill.GetId(), nil, true)
								g.Role.SetSkillCd(skill.GetId(), time.Now().Add(time.Duration(cd)*time.Second))
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
						}
						if skillItem.Camps == CampsEnemy {
							// 这是攻击技能
							if skillItem.NameZh == "普通攻击" {
								g.AttackClosestByName(g.ChangeAttackID(skill.GetId()), monsterList)
							} else {
								g.AttackClosestByName(skill.GetId(), monsterList)
							}
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
		if g.Role.GetBuffById(131080) != nil {
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
	<-g.Notifier(notifier.NtfType_AddAttributePoint)
	g.RemoveNotifier(notifier.NtfType_AddAttributePoint)
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

func (g *GameConnection) GetBuffByID(BuffID uint32) utils.BuffItem {
	g.Role.Mutex.RLock()
	defer g.Role.Mutex.RUnlock()
	for _, v := range g.Role.Buffs {
		if BuffID == v.GetId() {
			return g.BuffItems[v.GetId()]
		}
	}
	return utils.BuffItem{}
}

func (g *GameConnection) DelBuffByName(name string) {
	g.Role.Mutex.Lock()
	defer g.Role.Mutex.Unlock()
	if buffs, ok := g.BuffItemsByName[name]; ok {
		for _, buff := range buffs.Items {
			id, _ := buff.Id.Int64()
			if _, ok := g.Role.Buffs[uint32(id)]; !ok {
				continue
			}
			delete(g.Role.Buffs, uint32(id))
		}
	}
}

func (g *GameConnection) ChangeAttackID(skillID uint32) uint32 {

	Profession := g.Role.GetProfession()

	if Profession >= Cmd.EProfession_EPROFESSION_ARCHER && Profession <= Cmd.EProfession_EPROFESSION_RANGER {
		if g.Role.GetBuffById(131070) != nil {
			return 252001
		} else {
			return 300001
		}
	} else if Profession >= Cmd.EProfession_EPROFESSION_PRIEST && Profession <= Cmd.EProfession_EPROFESSION_ARCHBISHOP {
		if g.Role.GetBuffById(129040) != nil {
			return 406001
		} else {
			return 143001
		}
	}

	return skillID
}
