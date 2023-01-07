package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

// MoveChart Request to move current character to position
func (g *GameConnection) MoveChart(pos *Cmd.ScenePos) {
	cmd := &Cmd.ReqMoveUserCmd{
		Target: pos,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
	)
}

func (g *GameConnection) ParsePos(x, y, z int32) *Cmd.ScenePos {
	return &Cmd.ScenePos{
		X: &x,
		Y: &y,
		Z: &z,
	}
}

// MoveChartWait Return until character move to target position or timed out
func (g *GameConnection) MoveChartWait(pos *Cmd.ScenePos) bool {
	// orgPos := g.Role.GetPos()
	cmd := &Cmd.ReqMoveUserCmd{
		Target: pos,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
	)
	// dir := uint32(utils.CalcDir(utils.GetAngleByAxisY(&orgPos, pos)))
	// cmd2 := &Cmd.SetDirection{
	// 	Dir: &dir,
	// }
	// g.sendProtoCmd(
	// 	cmd2,
	// 	sceneUser2CmdId,
	// 	Cmd.User2Param_value["USER2PARAM_SET_DIRECTION"],
	// )
	count := 0
	arrived := false
loop:
	for {
		select {
		case <-time.After(1 * time.Minute):
			log.Warnf("MoveChartWait timed out")
			break loop
		case <-time.Tick(1 * time.Second):
			curPos := g.Role.GetPos()
			if utils.GetDistanceXZ(&curPos, pos) < 50 {
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
