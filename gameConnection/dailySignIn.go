package gameConnection

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	user2ParamSignInId = Cmd.User2Param_value["USER2PARAM_SIGNIN"]
)

func (g *GameConnection) DailySignIn() {
	if g.Role.DailySignIn != nil {
		if g.Role.DailySignIn.GetIssign() == 0 && g.Role.DailySignIn.GetType() == Cmd.ESignInType_ESIGNINTYPE_ACTIVITY {
			log.Infof("daily signIn %d days", g.Role.DailySignIn.GetCount())
			signType := Cmd.ESignInType_ESIGNINTYPE_ACTIVITY
			cmd := &Cmd.SignInUserCmd{
				Type: &signType,
			}
			g.sendProtoCmd(
				cmd,
				sceneUser2CmdId,
				user2ParamSignInId,
			)
		} else if g.Role.DailySignIn.GetIssign() == 1 {
			log.Infof("daily signin already done")
		} else {
			log.Warnf("unknown daily signin type %s, day count %d, is sign %d",
				g.Role.DailySignIn.GetType(),
				g.Role.DailySignIn.GetCount(),
				g.Role.DailySignIn.GetIssign(),
			)
		}
	}
}
