package gameConnection

import (
	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
)

var (
	sessionUserGuildCmdId = Cmd.Command_value["SESSION_USER_GUILD_PROTOCMD"]
)

func (g *GameConnection) GetGuildName() string {
	return g.Role.GuildData.GetName()
}

func (g *GameConnection) GuildDonateList() (res *Cmd.DonateListGuildCmd, err error) {
	cmd := Cmd.DonateListGuildCmd{}
	g.AddNotifier(gameTypes.NtfType_GuildParamDonateList)
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserGuildCmdId,
		Cmd.GuildParam_value["GUILDPARAM_DONATELIST"],
	)
	r, err := g.waitForResponse(gameTypes.NtfType_GuildParamDonateList)
	if err != nil {
		return nil, err
	}
	res = r.(*Cmd.DonateListGuildCmd)
	return res, nil
}

func (g *GameConnection) GuildDonate(donateItem *Cmd.DonateItem) {
	t := donateItem.GetTime()
	c := donateItem.GetConfigid()
	cmd := Cmd.DonateGuildCmd{
		Configid: &c,
		Time:     &t,
	}
	_ = g.sendProtoCmd(
		&cmd,
		sessionUserGuildCmdId,
		Cmd.GuildParam_value["GUILDPARAM_DONATE"],
	)
}
