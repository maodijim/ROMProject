package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

func (g *GameConnection) TakePhoto(focus *Cmd.CameraFocus, pos Cmd.ScenePos) {
	photoState := Cmd.ECreatureStatus_ECREATURESTATUS_SELF_PHOTO
	g.StateChangeCmd(&photoState)
	g.TakePhotoCmd(focus, pos)
	g.StateChangeCmd(nil)
}

func (g *GameConnection) TakePhotoCmd(focus *Cmd.CameraFocus, pos Cmd.ScenePos) {
	if focus != nil {
		_ = g.sendProtoCmd(
			focus,
			sceneUser2CmdId,
			Cmd.User2Param_value["USER2PARAM_CAMERAFOCUS"],
		)
	}
	number := int32(1)
	data := Cmd.PhaseData{
		Number: &number,
		Pos:    &pos,
	}
	g.SkillCmd(20004001, &data, true)
}

func (g *GameConnection) StateChangeCmd(state *Cmd.ECreatureStatus) {
	cmd := Cmd.StateChange{
		Status: state,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sceneUser2CmdId,
		Cmd.User2Param_value["USER2PARAM_STATECHANGE"],
	)
}

func (g *GameConnection) SceneryCmd(sceneryId uint32) {
	cmd := Cmd.SceneryUserCmd{
		Scenerys: []*Cmd.Scenery{},
	}
	cmd.Scenerys = append(cmd.Scenerys, &Cmd.Scenery{Sceneryid: &sceneryId})
	_ = g.sendProtoCmd(
		&cmd,
		sceneUser2CmdId,
		Cmd.User2Param_value["USER2PARAM_SCENERY"],
	)
}
