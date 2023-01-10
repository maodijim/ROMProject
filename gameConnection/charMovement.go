package gameConnection

import (
	"fmt"
	"math"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

// MoveChart Request to move current character to position
func (g *GameConnection) MoveChart(pos Cmd.ScenePos) {
	cmd := &Cmd.ReqMoveUserCmd{
		Target: &pos,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
	)
}

func (g *GameConnection) ParsePos(x, y, z int32) Cmd.ScenePos {
	return Cmd.ScenePos{
		X: &x,
		Y: &y,
		Z: &z,
	}
}

// MoveChartWait Return until character move to target position or timed out
func (g *GameConnection) MoveChartWait(pos Cmd.ScenePos) bool {
	// orgPos := g.Role.GetPos()
	cmd := &Cmd.ReqMoveUserCmd{
		Target: &pos,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
	)
	count := 0
	arrived := false
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-time.After(1 * time.Minute):
			log.Warnf("MoveChartWait timed out")
			break loop
		case <-ticker.C:
			curPos := g.Role.GetPos()
			distanceXY := utils.GetDistanceXY(curPos, pos)
			distanceXZ := utils.GetDistanceXZ(curPos, pos)
			if math.Max(distanceXY, distanceXZ) < 200 {
				arrived = true
				break loop
			} else if count > 200 {
				break loop
			} else {
				count += 1
			}
		}
	}
	return arrived
}

func (g *GameConnection) GoToMapGear(mapId uint32) {
	cmd := &Cmd.GoToGearUserCmd{
		Mapid: &mapId,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER2_PROTOCMD"],
		Cmd.User2Param_value["USER2PARAM_GOTO_GEAR"],
	)
}

func (g *GameConnection) ChangeMap(mId uint32) {
	cmd := &Cmd.ChangeSceneUserCmd{
		MapID: &mId,
	}
	log.Infof("%s is sending change scene cmd: %v", g.Role.GetRoleName(), cmd)
	_ = g.sendProtoCmd(cmd, 5, 23)
	g.enteringMap = true
	g.inMap = true
}

func (g *GameConnection) ExitMap(targetMapId uint32) {
	cmd := &Cmd.GoToExitPosUserCmd{
		Mapid: &targetMapId,
	}
	_ = g.sendProtoCmd(cmd, sceneUserCmdId, Cmd.CmdParam_value["GOTO_EXIT_POS_USER_CMD"])
	time.Sleep(2 * time.Second)
	g.ChangeMap(targetMapId)
}

func (g *GameConnection) ExitMapWait(mapId uint32) {
	g.ExitMap(mapId)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			if g.Role.GetMapId() != mapId {
				continue
			}
			g.ChangeMap(mapId)
			return
		}
	}
}

func (g *GameConnection) MoveToNpcWait(npcName string) error {
	npcs := g.GetMapNpcs()
	for _, npc := range npcs {
		if npc.GetName() == npcName {
			g.MoveChartWait(*npc.GetPos())
			return nil
		}
	}
	return fmt.Errorf("npc %s not found", npcName)
}
