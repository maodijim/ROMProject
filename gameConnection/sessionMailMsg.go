package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	"github.com/mohae/deepcopy"
)

func (g *GameConnection) HandleSessionMailMsg(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.MailParam_value["MAILPARAM_QUERYALLMAIL"]:
		param = &Cmd.QueryAllMail{}
		err = utils.ParseCmd(rawData, param)
		mails, ok := param.(*Cmd.QueryAllMail)
		if ok {
			g.Mutex.Lock()
			defer g.Mutex.Unlock()
			g.mails = mails.GetDatas()
		}

	case Cmd.MailParam_value["MAILPARAM_UPDATE"]:
		param = &Cmd.MailUpdate{}
		err = utils.ParseCmd(rawData, param)
		g.Mutex.Lock()
		defer g.Mutex.Unlock()
		newMails := deepcopy.Copy(g.mails).([]*Cmd.MailData)
		for _, mailId := range param.(*Cmd.MailUpdate).GetDels() {
			for _, mail := range g.mails {
				if mail.GetId() != mailId {
					newMails = append(newMails, mail)
				}
			}
		}
		for _, mail := range param.(*Cmd.MailUpdate).GetUpdates() {
			newMails = append(newMails, mail)
		}
		g.mails = newMails
	}
	return param, err
}
