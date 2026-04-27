package gameConnection

import (
	Cmd "ROMProject/Cmds"
)

var (
	mailCmdId = Cmd.Command_value["SESSION_USER_MAIL_PROTOCMD"]
)

func (g *GameConnection) GetMails() []*Cmd.MailData {
	return g.mails
}

func (g *GameConnection) GetMailAttachment(id uint64) {
	_ = g.sendProtoCmd(
		&Cmd.GetMailAttach{
			Id: &id,
		},
		mailCmdId,
		Cmd.MailParam_value["MAILPARAM_GETATTACH"],
	)
}
