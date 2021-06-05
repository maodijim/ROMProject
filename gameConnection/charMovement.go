package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"time"
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

// MoveChartWait Return until character move to target position or timed out
func (g *GameConnection) MoveChartWait(pos *Cmd.ScenePos) (curPos *Cmd.ScenePos) {
	orgPos := &g.Role.Pos
	cmd := &Cmd.ReqMoveUserCmd{
		Target: pos,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["SCENE_USER_PROTOCMD"],
		Cmd.CmdParam_value["REQ_MOVE_USER_CMD"],
	)
	count := 0
	for {
		if utils.GetDistanceXZ(*orgPos, g.Role.Pos) < 1500 || count > 200 {
			break
		} else {
			count += 1
			time.Sleep(500 * time.Millisecond)
		}
	}
	return g.Role.Pos
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
