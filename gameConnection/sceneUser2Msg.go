package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func (g *GameConnection) HandleSceneUser2ProtoCmd(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.User2Param_value["USER2PARAM_CDTIME"]:
		// TODO handles items cd time
		param = &Cmd.CDTimeUserCmd{}
		err = utils.ParseCmd(rawData, param)

	case Cmd.User2Param_value["USER2PARAM_SIGNIN_NTF"]:
		param = &Cmd.SignInNtfUserCmd{}
		err = utils.ParseCmd(rawData, param)
		dailySign := param.(*Cmd.SignInNtfUserCmd)
		g.Role.DailySignIn = dailySign

	case Cmd.User2Param_value["USER2PARAM_SERVANT_RECOMMEND"]:
		param = &Cmd.RecommendServantUserCmd{}
		err = utils.ParseCmd(rawData, param)
		recommendServant := param.(*Cmd.RecommendServantUserCmd)
		for _, i := range recommendServant.GetItems() {
			if i.GetStatus() == Cmd.ERecommendStatus_ERECOMMEND_STATUS_RECEIVE {
				go func() {
					time.Sleep(time.Second)
					g.TakeServantReward(i.GetDwid())
				}()
			}
		}

	case Cmd.User2Param_value["USER2PARAM_INVITEFOLLOW"]:
		param = &Cmd.InviteFollowUserCmd{}

	case Cmd.User2Param_value["USER2PARAM_BUFFERSYNC"]:
		param = &Cmd.UserBuffNineSyncCmd{}
		err = utils.ParseCmd(rawData, param)
		buffSync := param.(*Cmd.UserBuffNineSyncCmd)
		g.Mutex.Lock()
		if buffSync.GetGuid() == g.Role.GetRoleId() {
			for _, updateBuff := range buffSync.GetUpdates() {
				g.Role.Buffs[updateBuff.GetId()] = updateBuff
			}

			for _, delBuff := range buffSync.GetDels() {
				if g.Role.Buffs[delBuff] != nil {
					delete(g.Role.Buffs, delBuff)
				}
			}
		}
		g.Mutex.Unlock()

	case Cmd.User2Param_value["USER2PARAM_VAR"]:
		param = &Cmd.VarUpdate{}
		err = utils.ParseCmd(rawData, param)
		userVar := param.(*Cmd.VarUpdate)
		g.Mutex.Lock()
		for _, uv := range userVar.GetVars() {
			g.Role.UserVars[uv.GetType()] = uv
		}
		for _, av := range userVar.GetAccvars() {
			g.Role.AccVars[av.GetType()] = av
		}
		g.Mutex.Unlock()

	case Cmd.User2Param_value["USER2PARAM_GOTO_LIST"]:
		param = &Cmd.GoToListUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.GotoList = param.(*Cmd.GoToListUserCmd)

	case Cmd.User2Param_value["USER2PARAM_READYTOMAP"]:
		param = &Cmd.ReadyToMapUserCmd{}
		err = utils.ParseCmd(rawData, param)
		rMap := param.(*Cmd.ReadyToMapUserCmd)
		if rMap.GetMapID() != 0 {
			log.Debugf("Ready to move to map ID: %d", rMap.GetMapID())
			g.SetEnteringMap()
			g.Role.MapId = rMap.MapID
		}

	case Cmd.User2Param_value["USER2PARAM_NPCDATASYNC"]:
		param = &Cmd.NpcDataSync{}
		err = utils.ParseCmd(rawData, param)
		dataSync := param.(*Cmd.NpcDataSync)
		if dataSync != nil {
			g.Mutex.Lock()
			var userDatas []*Cmd.UserData
			var userAttrs []*Cmd.UserAttr
			if g.MapNpcs[dataSync.GetGuid()] != nil {
				userDatas = g.MapNpcs[dataSync.GetGuid()].GetDatas()
				userAttrs = g.MapNpcs[dataSync.GetGuid()].GetAttrs()
			} else if g.MapUsers[dataSync.GetGuid()] != nil {
				userDatas = g.MapUsers[dataSync.GetGuid()].GetDatas()
				userAttrs = g.MapUsers[dataSync.GetGuid()].GetAttrs()
			}
			// update NPC data
			for _, ds := range dataSync.GetDatas() {
				for _, data := range userDatas {
					if data.GetType() == ds.GetType() {
						data.Data = ds.GetData()
					}
				}
			}
			// update NPC attr
			for _, as := range dataSync.GetAttrs() {
				for _, attr := range userAttrs {
					if attr.GetType() == as.GetType() {
						attr.Value = as.Value
					}
				}
			}
			g.Mutex.Unlock()
		}

	case Cmd.User2Param_value["USER2PARAM_QUERY_MAPAREA"]:
		param = &Cmd.QueryMapArea{}
		err = utils.ParseCmd(rawData, param)
		if len(param.(*Cmd.QueryMapArea).Areas) > 0 {
			g.Role.SetMapId(param.(*Cmd.QueryMapArea).GetAreas()[0])
		}
	}
	return param, err
}
