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
	_ = g.sendProtoCmd(
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
	_ = g.sendProtoCmd(
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
			if math.Max(distanceXY, distanceXZ) <= 400 {
				arrived = true
				time.Sleep(time.Second)
				break loop
			} else if count > 150 {
				break loop
			} else {
				count += 1
				_ = g.sendProtoCmd(
					cmd,
					Cmd.Command_value["SCENE_USER_PROTOCMD"],
					Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
				)
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
	g.enteringMap = false
	g.inMap = false
	// If not moved strange things will happen
	g.MoveChart(g.Role.GetPos())
}

func (g *GameConnection) ChangeMapWithPos(mId uint32, pos Cmd.ScenePos) {
	cmd := &Cmd.ChangeSceneUserCmd{
		MapID: &mId,
		Pos:   &pos,
	}
	log.Infof("%s is sending change scene cmd: %v", g.Role.GetRoleName(), cmd)
	_ = g.sendProtoCmd(cmd, 5, 23)
	g.enteringMap = false
	g.inMap = false
	// If not moved strange things will happen
	g.MoveChart(g.Role.GetPos())
}

func (g *GameConnection) ExitMap(targetMapId uint32) {
	cmd := &Cmd.GoToExitPosUserCmd{
		Mapid: &targetMapId,
	}
	g.inMap = false
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

func (g *GameConnection) ExitMapPos(targetMapId uint32, exitId uint32, exitPos Cmd.ScenePos) {
	cmd := &Cmd.ExitPosUserCmd{
		Mapid:  g.Role.MapId,
		Exitid: &exitId,
		Pos:    &exitPos,
	}
	g.inMap = false
	_ = g.sendProtoCmd(cmd, sceneUser2CmdId, Cmd.User2Param_value["USER2PARAM_EXIT_POS"])
	time.Sleep(3 * time.Second)
	g.ChangeMapWithPos(targetMapId, g.ParsePos(0, 0, 0))
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

func (g *GameConnection) MoveToNpcIdWait(npcId uint32) error {
	npcs := g.GetMapNpcs()
	for _, npc := range npcs {
		if npc.GetNpcID() == npcId {
			g.MoveChartWait(*npc.GetPos())
			return nil
		}
	}
	return fmt.Errorf("npc %d not found", npcId)
}

func (g *GameConnection) GoToGear(mapId uint32) error {
	cmd := &Cmd.GoToGearUserCmd{
		Mapid: &mapId,
	}
	err := g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER2_PROTOCMD"],
		Cmd.User2Param_value["USER2PARAM_GOTO_GEAR"],
	)
	if err != nil {
		return err
	}
	return nil
}
