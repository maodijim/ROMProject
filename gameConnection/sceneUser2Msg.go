package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"

	"google.golang.org/protobuf/proto"
)

func (g *GameConnection) HandleSceneUser2ProtoCmd(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.User2Param_value["USER2PARAM_SERVERTIME"]:
		param = &Cmd.ServerTime{}
		err = utils.ParseCmd(rawData, param)
		fixedSkillCDSubtract = param.(*Cmd.ServerTime).GetTime() - uint64(time.Now().UnixMilli())

	case Cmd.User2Param_value["USER2PARAM_GAMETIME"]:
		param = &Cmd.GameTimeCmd{}
		err = utils.ParseCmd(rawData, param)
		g.Role.GameTime = param.(*Cmd.GameTimeCmd)
		g.SendToNotifier(gameTypes.NtfType_GameTime, param)

	case Cmd.User2Param_value["USER2PARAM_TALKINFO"]:
		param = &Cmd.TalkInfo{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_TalkInfo, param)

	case Cmd.User2Param_value["USER2PARAM_SYSMSG"]:
		param = &Cmd.SysMsg{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_SysMsg, param)

	case Cmd.User2Param_value["USER2PARAM_QUERY_ZONESTATUS"]:
		param = &Cmd.QueryZoneStatusUserCmd{}
		err = utils.ParseCmd(rawData, param)

		g.SendToNotifier(gameTypes.NtfType_User2QueryZoneStatus, param)

	case Cmd.User2Param_value["USER2PARAM_ACTION"]:
		param = &Cmd.UserActionNtf{}
		err = utils.ParseCmd(rawData, param)
		if err == nil {
			userActionNtf := param.(*Cmd.UserActionNtf)
			if userActionNtf.GetType() == Cmd.EUserActionType_EUSERACTIONTYPE_DIALOG {
				g.SendToNotifier(gameTypes.NtfType_UserActionDialog, userActionNtf)
			}
		}

	case Cmd.User2Param_value["USER2PARAM_CDTIME"]:
		// TODO handles items cd time
		param = &Cmd.CDTimeUserCmd{}
		err = utils.ParseCmd(rawData, param)
		cdTime := param.(*Cmd.CDTimeUserCmd)
		if cdTime.GetList() != nil {
			for _, cd := range cdTime.GetList() {
				if cd.GetType() == Cmd.CD_TYPE_CD_TYPE_SKILL {
					if cd.GetId() == 50057001 {
						// ignore 备战精英
						continue
					}
					if skillItem, ok := g.SkillItems[cd.GetId()]; ok && skillItem.NameZh == "普通攻击" {
						// ignore 普通攻击
						continue
					}
					now := time.Now().UnixMilli()
					calculatedCd := int64(cd.GetTime() - fixedSkillCDSubtract)
					if calculatedCd < now {
						continue
					}
					nowFormatted := time.UnixMilli(now).Format(time.RFC3339)
					calculatedCdFormatted := time.UnixMilli(calculatedCd).Format(time.RFC3339)
					cdTimeFormatted := time.UnixMilli(int64(cd.GetTime())).Format(time.RFC3339)
					g.logger.Infof(
						"Skill CD received: SkillID=%d, CDTime=%d (%s), Now=%s, CalculatedCD=%d (%s)",
						cd.GetId(),
						cd.GetTime(),
						cdTimeFormatted,
						nowFormatted,
						calculatedCd,
						calculatedCdFormatted)
					g.Role.SetSkillCd(cd.GetId(), time.UnixMilli(calculatedCd))
				}
			}
		}

	// case Cmd.User2Param_value["USER2PARAM_SIGNIN_NTF"]:
	// 	param = &Cmd.SignInNtfUserCmd{}
	// 	err = utils.ParseCmd(rawData, param)
	// 	dailySign := param.(*Cmd.SignInNtfUserCmd)
	// 	g.Role.DailySignIn = dailySign

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
		if buffSync.GetGuid() == g.Role.GetRoleId() {
			for _, updateBuff := range buffSync.GetUpdates() {
				g.Role.Mutex.Lock()
				g.Role.Buffs[updateBuff.GetId()] = updateBuff
				g.Role.Mutex.Unlock()
			}

			for _, delBuff := range buffSync.GetDels() {
				g.Role.Mutex.Lock()
				if g.Role.Buffs[delBuff] != nil {
					delete(g.Role.Buffs, delBuff)
				}
				g.Role.Mutex.Unlock()
			}
		}

	case Cmd.User2Param_value["USER2PARAM_VAR"]:
		param = &Cmd.VarUpdate{}
		err = utils.ParseCmd(rawData, param)
		userVar := param.(*Cmd.VarUpdate)
		g.Role.Mutex.Lock()
		for _, uv := range userVar.GetVars() {
			g.Role.UserVars[uv.GetType()] = uv
		}
		// for _, av := range userVar.GetAccvars() {
		// 	g.Role.AccVars[av.GetType()] = av
		// }
		g.Role.Mutex.Unlock()

	case Cmd.User2Param_value["USER2PARAM_GOTO_LIST"]:
		param = &Cmd.GoToListUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.GotoList = param.(*Cmd.GoToListUserCmd)

	// case Cmd.User2Param_value["USER2PARAM_READYTOMAP"]:
	// 	param = &Cmd.ReadyToMapUserCmd{}
	// 	err = utils.ParseCmd(rawData, param)
	// 	rMap := param.(*Cmd.ReadyToMapUserCmd)
	// 	if rMap.GetMapID() != 0 {
	// 		log.Debugf("Ready to move to map ID: %d", rMap.GetMapID())
	// 		g.SetEnteringMap()
	// 		g.Role.MapId = rMap.MapID
	// 	}

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

	case Cmd.User2Param_value["USER2PARAM_EFFECT"]:
		param = &Cmd.EffectUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.SendToNotifier(gameTypes.NtfType_EffectUser, param)
	}

	return param, err
}
