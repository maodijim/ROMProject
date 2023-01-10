package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

var (
	sceneUserSkillCmdId = Cmd.Command_value["SCENE_USER_SKILL_PROTOCMD"]
)

func (g *GameConnection) LevelUpSkill(skillIds []uint32, skillType Cmd.ELevelupType) {
	cmd := Cmd.LevelupSkill{
		Skillids: skillIds,
		Type:     &skillType,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sceneUserSkillCmdId,
		Cmd.SkillParam_value["SKILLPARAM_LEVELUPSKILL"],
	)
}

func (g *GameConnection) GetSkillIdByName(skillName string, level uint32) uint32 {
	if skills, ok := g.SkillItemsByName[skillName]; ok && len(skills) > 0 {
		baseId, _ := skills[0].Id.Int64()
		return uint32(baseId/100)*100 + level
	}
	return 0
}

func (g *GameConnection) IsSkillLearnedByName(skillName string, level uint32) bool {
	skillId := g.GetSkillIdByName(skillName, level)
	if skillId == 0 {
		return false
	}
	return g.Role.IsSkillLearned(skillId)
}

func (g *GameConnection) GetLearnedSkillLevelByName(skillName string) uint32 {
	lSkill := g.Role.GetLearnedSkill()
	skills, ok := g.SkillItemsByName[skillName]
	if !ok {
		return 0
	}
	skillBaseId, _ := skills[0].Id.Int64()
	skillBaseId = skillBaseId / 100
	for _, skill := range lSkill {
		if skill.GetId()/100 == uint32(skillBaseId) {
			return skill.GetId() % uint32(skillBaseId)
		}
	}
	return 0
}
