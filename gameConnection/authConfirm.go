package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"time"
)

func (g *GameConnection) AuthConfirm(pass string) {
	cmd := &Cmd.ConfirmAuthorizeUserCmd{
		Password: &pass,
	}
	g.sendProtoCmd(
		cmd,
		Cmd.Command_value["LOGIN_USER_PROTOCMD"],
		Cmd.LoginCmdParam_value["CONFIRM_AUTHORIZE_USER_CMD"],
	)
	time.Sleep(time.Second)
}
