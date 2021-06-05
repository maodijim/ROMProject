package gameConnection

import (
	Cmd "ROMProject/Cmds"
	log "github.com/sirupsen/logrus"
)

var (
	SocialityProtoCmdId = Cmd.Command_value["SESSION_USER_SOCIALITY_PROTOCMD"]
)

func (g *GameConnection) FindUser(userName string) (result *Cmd.FindUser, err error) {
	cmd := &Cmd.FindUser{
		Keyword: &userName,
	}
	g.addNotifier("SOCIALITYPARAM_FINDUSER")
	g.sendProtoCmd(cmd,
		SocialityProtoCmdId,
		Cmd.SocialityParam_value["SOCIALITYPARAM_FINDUSER"],
	)
	res, err := g.waitForResponse("SOCIALITYPARAM_FINDUSER")
	if err != nil {
		log.Errorf("failed to find user %s: %s", userName, err)
	}
	if res != nil {
		result = res.(*Cmd.FindUser)
	}
	return result, err
}
