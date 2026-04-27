package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	TowerProtoCmdId = Cmd.Command_value["INFINITE_TOWER_PROTOCMD"]
)

func (g *GameConnection) TowerReplyAgree() {
	agree := Cmd.ETowerReply_ETOWERREPLY_AGREE
	cmd := &Cmd.TeamTowerReplyCmd{
		EReply: &agree,
		Userid: g.Role.RoleId,
	}
	_ = g.sendProtoCmd(cmd,
		TowerProtoCmdId,
		Cmd.TowerParam_value["ETOWERPARAM_REPLY"],
	)
}

func (g *GameConnection) GetTeamTowerSummary() (towerSummary *Cmd.TeamTowerSummary) {
	cmd := &Cmd.TeamTowerInfoCmd{
		Teamid: g.Role.TeamData.Guid,
	}
	g.AddNotifier("ETOWERPARAM_TEAMTOWERSUMMARY")
	_ = g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TowerParam_value["ETOWERPARAM_TEAMTOWERINFO"],
	)
	res, err := g.waitForResponse("ETOWERPARAM_TEAMTOWERSUMMARY")
	if err != nil {
		log.Errorf("failed to get team tower summary: %v", err)
	}
	if res != nil {
		towerSummary = res.(*Cmd.TeamTowerSummary)
	}
	return towerSummary
}

func (g *GameConnection) ExitTower() {
	num := int32(1)
	dir := int32(0)
	pos := g.Role.GetPos()
	pData := &Cmd.PhaseData{
		Number: &num,
		Pos:    &pos,
		Dir:    &dir,
	}
	g.SkillCmd(20002001, pData, true)
	time.Sleep(3 * time.Second)
}
