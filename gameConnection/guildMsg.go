package gameConnection

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"google.golang.org/protobuf/proto"
)

func (g *GameConnection) HandleGuildMsg(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.GuildParam_value["GUILDPARAM_ENTERGUILD"]:
		param = &Cmd.EnterGuildGuildCmd{}
		err = utils.ParseCmd(rawData, param)
		g.Role.GuildData = param.(*Cmd.EnterGuildGuildCmd).GetData()

	case Cmd.GuildParam_value["GUILDPARAM_QUERYPACK"]:
		param = &Cmd.QueryPackGuildCmd{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_GuildParamQueryPack, param)

	case Cmd.GuildParam_value["GUILDPARAM_DONATELIST"]:
		param = &Cmd.DonateListGuildCmd{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_GuildParamDonateList, param)

	}
	return param, err
}
