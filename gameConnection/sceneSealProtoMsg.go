package gameConnection

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"google.golang.org/protobuf/proto"
)

func (g *GameConnection) HandleSceneSealProtoMsg(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.SealParam_value["SEALPARAM_UPDATESEAL"]:
		param = &Cmd.UpdateSeal{}
		err = utils.ParseCmd(rawData, param)
		updateS := param.(*Cmd.UpdateSeal)
		for _, delSeal := range updateS.GetDeldata() {
			var newSealData []*Cmd.SealData
			for i, val := range g.Role.SealData {
				if delSeal.GetMapid() != val.GetMapid() {
					newSealData = append(newSealData, g.Role.SealData[i])
				}
			}
			g.Role.SealData = newSealData
		}
		for _, addSeal := range updateS.GetNewdata() {
			found := false
			for i, val := range g.Role.SealData {
				if addSeal.GetMapid() == val.GetMapid() {
					g.Role.SealData[i] = addSeal
					found = true
					break
				}
			}
			if !found {
				g.Role.SealData = append(g.Role.SealData, addSeal)
			}
		}

	case Cmd.SealParam_value["SEALPARAM_ACCEPTSEAL"]:
		param = &Cmd.SealAcceptCmd{}
		err = utils.ParseCmd(rawData, param)
		g.Mutex.Lock()
		g.Role.AcceptSeal = param.(*Cmd.SealAcceptCmd)
		g.Mutex.Unlock()
		g.SendToNotifier(gameTypes.NtfType_SealParamAcceptSeal, param)

	case Cmd.SealParam_value["SEALPARAM_QUERYLIST"]:
		param = &Cmd.SealQueryList{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_SealParamQueryList, param)
	}
	return param, err
}
