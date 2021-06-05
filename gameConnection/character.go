package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
	"regexp"
)

var (
	sceneUser2CmdId = Cmd.Command_value["SCENE_USER2_PROTOCMD"]
)

func (g *GameConnection) GetAutoSkills() []*Cmd.SkillItem {
	var sList []*Cmd.SkillItem
	g.Mutex.RLock()
	for _, skill := range g.Role.AutoSkills {
		sList = append(sList, skill)
	}
	g.Mutex.RUnlock()
	return sList
}

func (g *GameConnection) updateAutoSkill(skill *Cmd.SkillItem) {
	for _, s := range skill.GetShortcuts() {
		if s.GetType() == Cmd.ESkillShortcut_ESKILLSHORTCUT_AUTO {
			g.Role.AutoSkills[s.GetPos()] = skill
		}
	}
}

func (g *GameConnection) GetAtkRange(skillId uint32) uint32 {

	//g.SkillItems[]
	return 0
}

func (g *GameConnection) GetBuffNameByRegex(searchRegex string) (buffName string) {
	g.Mutex.RLock()
	for _, curBuff := range g.Role.Buffs {
		re := regexp.MustCompile(searchRegex)
		search := re.Find([]byte(g.BuffItems[curBuff.GetId()].BuffName))
		if len(search) > 0 {
			buffName = string(search)
		}
	}
	g.Mutex.RUnlock()
	return buffName
}

func (g *GameConnection) GetAtkSpd() int32 {
	return utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_ATKSPD)
}

func (g *GameConnection) GetCurrentHp() int32 {
	return utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_HP)
}

func (g *GameConnection) GetMaxHp() int32 {
	return utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_MAXHP)
}

func (g *GameConnection) GetCurrentSp() int32 {
	return utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_SP)
}

func (g *GameConnection) GetMaxSp() int32 {
	return utils.GetNpcAttrValByType(g.Role.UserAttrs, Cmd.EAttrType_EATTRTYPE_MAXSP)
}

func (g *GameConnection) GetSpPer() float64 {
	return float64(g.GetCurrentSp()) / float64(g.GetMaxSp())
}

func (g *GameConnection) GetHpPer() float64 {
	return float64(g.GetCurrentHp()) / float64(g.GetMaxHp())
}

func (g *GameConnection) GoToMap(mapId uint32) {
	for _, goToMapId := range g.GotoList.GetMapid() {
		if goToMapId == mapId {
			cmd := &Cmd.GoToGearUserCmd{
				Mapid: &mapId,
			}
			g.sendProtoCmd(cmd,
				sceneUser2CmdId,
				Cmd.User2Param_value["USER2PARAM_GOTO_GEAR"],
			)
			return
		}
	}
	log.Warnf("mapId: %d is not in map goto list %v", mapId, g.GotoList)
}

// 领取执事奖励
func (g *GameConnection) takeServantReward(wid uint32) {
	cmd := &Cmd.ReceiveServantUserCmd{
		Dwid: &wid,
	}
	g.sendProtoCmd(cmd,
		sceneUser2CmdId,
		Cmd.User2Param_value["USER2PARAM_SERVANT_RECEIVE"],
	)
}
