package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
)

func (g *GameConnection) HandleSceneUserMapProtoCmd(cmdParamId int32, rawData []byte) (err error) {
	var param proto.Message
	switch cmdParamId {
	case Cmd.MapParam_value["MAPPARAM_ADDMAPUSER"]:
		param = &Cmd.AddMapUser{}
		err = utils.ParseCmd(rawData, param)
		addUsers := param.(*Cmd.AddMapUser)
		g.Mutex.Lock()
		for _, user := range addUsers.GetUsers() {
			g.MapUsers[user.GetGuid()] = user
			if user.GetGuid() == g.Role.GetRoleId() {
				for _, buff := range user.GetBuffs() {
					g.Role.Buffs[buff.GetId()] = buff
				}
				g.UpdateUserParams(user.GetDatas(), user.GetAttrs())
			}
		}
		g.Mutex.Unlock()

	case Cmd.MapParam_value["MAPPARAM_ADDMAPNPC"]:
		param = &Cmd.AddMapNpc{}
		err = utils.ParseCmd(rawData, param)
		addNpcs := param.(*Cmd.AddMapNpc)
		g.Mutex.Lock()
		for _, npc := range addNpcs.GetNpcs() {
			g.MapNpcs[npc.GetId()] = npc
		}
		g.Mutex.Unlock()

	case Cmd.MapParam_value["MAPPARAM_MAP_CMD_END"]:
		param = &Cmd.MapCmdEnd{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.MapParam_value["MAPPARAM_ADDMAPITEM"]:
		param = &Cmd.AddMapItem{}
		err = utils.ParseCmd(rawData, param)
		// pickup item own by user
		g.PickupMapItem(param.(*Cmd.AddMapItem))
	}
	return err
}
