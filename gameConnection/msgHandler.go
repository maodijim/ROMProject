package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	queryOnce = false
)

func (g *GameConnection) waitForResponse(notifierType string) (res interface{}, err error) {
	start := time.Now()
	for time.Since(start) < queryTimeout {
		select {
		case res = <-g.notifier[notifierType]:
			g.removeNotifier(notifierType)
			return res, err
		default:
			time.Sleep(time.Second)
		}
	}
	if time.Since(start) > queryTimeout {
		err = ErrQueryTimeout
	}
	g.removeNotifier(notifierType)
	return res, err
}

func (g *GameConnection) HandleMsg(output [][]byte) {
	if g.Role.GetRoleId() != 0 && g.Role.GetRoleName() != "" && g.Role.GetAuthenticated() && g.conn != nil && !g.Role.GetRoleSelected() {
		g.SelectRole()
	}
	if g.conn != nil && g.Role.GetMapId() != 0 && g.Role.GetInGame() && !g.enteringMap && g.Role.GetLoginResult() == 0 {
		g.enterGameMap()
	}

	for _, o := range output {
		if len(o) < 2 {
			log.Warn("result is empty")
			continue
		}
		cmdId := int32(o[0])
		cmdName := Cmd.Command_name[cmdId]
		cmdParamId := int32(o[1])
		cmdParamName := utils.NameParamMap[cmdName][cmdParamId]
		var param proto.Message
		var err error
		switch cmdId {
		case Cmd.Command_value["LOGIN_USER_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.LoginCmdParam_value["LOGIN_RESULT_USER_CMD"]:
				param = &Cmd.LoginResultUserCmd{}
				err = utils.ParseCmd(o, param)
				g.Role.SetLoginResult(param.(*Cmd.LoginResultUserCmd).GetRet())

			case Cmd.LoginCmdParam_value["REAL_AUTHORIZE_USER_CMD"]:
				param = &Cmd.RealAuthorizeUserCmd{}
				err = proto.Unmarshal(o[2:], param)
				g.Role.SetAuthenticated(param.(*Cmd.RealAuthorizeUserCmd).GetAuthorized())

			case Cmd.LoginCmdParam_value["HEART_BEAT_USER_CMD"]:
				param = &Cmd.HeartBeatUserCmd{}
				err = utils.ParseCmd(o, param)
				if g.conn != nil {
					g.currentIndex = 1
					g.ShouldHeartBeat = true
					g.lastHeartBeat = time.Now()
				}

			case Cmd.LoginCmdParam_value["SNAPSHOT_USER_CMD"]:
				param = &Cmd.SnapShotUserCmd{}
				err = utils.ParseCmd(o, param)
				roleData := param.(*Cmd.SnapShotUserCmd)
				roleId := roleData.GetData()[0].GetId()
				if len(roleData.GetData()) >= int(g.Configs.Char) {
					roleId = roleData.GetData()[g.Configs.Char-1].GetId()
				}
				g.Role.SetRoleId(roleId)
				log.Infof("setting role id to %d", g.Role.GetRoleId())
				for _, char := range param.(*Cmd.SnapShotUserCmd).GetData() {
					if char.Id != nil && g.Role.GetRoleId() == *char.Id {
						g.Role.SetRoleName(*char.Name)
						log.Infof("setting role name to %s", g.Role.GetRoleName())
					}
				}

			case Cmd.LoginCmdParam_value["SERVERTIME_USER_CMD"]:
				param = &Cmd.ServerTimeUserCmd{}
				err = utils.ParseCmd(o, param)
				g.currentIndex = 1
				if g.conn != nil {
					g.sendServerTimeUserCmd(0)
				}

			case Cmd.LoginCmdParam_value["CONFIRM_AUTHORIZE_USER_CMD"]:
				param = &Cmd.ConfirmAuthorizeUserCmd{}
				err = utils.ParseCmd(o, param)
				if g.conn != nil {
					g.sendServerTimeUserCmd(0)
					g.Role.AuthConfirm = param.(*Cmd.ConfirmAuthorizeUserCmd).Success
					if *g.Role.AuthConfirm == false {
						log.Warn("account is NOT authorized to trade and perform certain actions")
					} else {
						log.Info("account is authorized to trade and perform certain actions")
					}
				}

			default:
				continue
			}

		case Cmd.Command_value["SCENE_USER2_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.User2Param_value["USER2PARAM_SERVANT_RECOMMEND"]:
				param = &Cmd.RecommendServantUserCmd{}
				err = utils.ParseCmd(o, param)
				recommendServant := param.(*Cmd.RecommendServantUserCmd)
				for _, i := range recommendServant.GetItems() {
					if i.GetStatus() == Cmd.ERecommendStatus_ERECOMMEND_STATUS_RECEIVE {
						go func() {
							time.Sleep(time.Second)
							g.takeServantReward(i.GetDwid())
						}()
					}
				}

			case Cmd.User2Param_value["USER2PARAM_INVITEFOLLOW"]:
				param = &Cmd.InviteFollowUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_BUFFERSYNC"]:
				param = &Cmd.UserBuffNineSyncCmd{}
				err = utils.ParseCmd(o, param)
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
				err = utils.ParseCmd(o, param)
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
				err = utils.ParseCmd(o, param)
				g.GotoList = param.(*Cmd.GoToListUserCmd)

			case Cmd.User2Param_value["USER2PARAM_READYTOMAP"]:
				param = &Cmd.ReadyToMapUserCmd{}
				err = utils.ParseCmd(o, param)
				rMap := param.(*Cmd.ReadyToMapUserCmd)
				if rMap.GetMapID() != 0 {
					log.Debugf("Ready to move to map ID: %d", rMap.GetMapID())
					g.enteringMap = true
					g.Role.MapId = rMap.MapID
				}

			case Cmd.User2Param_value["USER2PARAM_NPCDATASYNC"]:
				param = &Cmd.NpcDataSync{}
				err = utils.ParseCmd(o, param)
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
				err = utils.ParseCmd(o, param)
				if len(param.(*Cmd.QueryMapArea).Areas) > 0 {
					g.Role.SetMapId(param.(*Cmd.QueryMapArea).GetAreas()[0])
				}
			}

		case Cmd.Command_value["SCENE_USER_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.CmdParam_value["DELETE_ENTRY_USER_CMD"]:
				param = &Cmd.DeleteEntryUserCmd{}
				err = utils.ParseCmd(o, param)
				entries := param.(*Cmd.DeleteEntryUserCmd)
				for _, entryId := range entries.GetList() {
					g.Mutex.Lock()
					if g.MapNpcs[entryId] != nil {
						log.Debugf("delete entry: id: %d, name: %s",
							entryId,
							g.MapNpcs[entryId].GetName(),
						)
						delete(g.MapNpcs, entryId)
					} else if g.MapUsers[entryId] != nil {
						log.Debugf("delete entry: id: %d, name: %s",
							entryId,
							g.MapNpcs[entryId].GetName(),
						)
						delete(g.MapUsers, entryId)
					}
					g.Mutex.Unlock()
				}

			case Cmd.CmdParam_value["RET_MOVE_USER_CMD"]:
				param = &Cmd.RetMoveUserCmd{}
				err = utils.ParseCmd(o, param)
				cmd := param.(*Cmd.RetMoveUserCmd)
				g.Mutex.Lock()
				if g.Role.FollowUserId != 0 && g.Role.FollowUserId == cmd.GetCharid() {
					go func() {
						log.Debugf("following user %d to %v", cmd.GetCharid(), cmd.GetPos())
						g.MoveChart(cmd.GetPos())
					}()
				}
				if cmd.GetCharid() == *g.Role.RoleId {
					log.Infof("Moving charater to position: %v", param.(*Cmd.RetMoveUserCmd).GetPos())
					g.Role.Pos = cmd.GetPos()
				} else if g.MapNpcs[cmd.GetCharid()] != nil {
					g.MapNpcs[cmd.GetCharid()].Pos = cmd.GetPos()
				} else if g.MapUsers[cmd.GetCharid()] != nil {
					g.MapUsers[cmd.GetCharid()].Pos = cmd.GetPos()
				}
				g.Mutex.Unlock()

			case Cmd.CmdParam_value["USERPARAM_USERSYNC"]:
				param = &Cmd.UserSyncCmd{}
				err = utils.ParseCmd(o, param)
				if param.(*Cmd.UserSyncCmd).GetType() == Cmd.EUserSyncType_EUSERSYNCTYPE_SYNC {
					datas := param.(*Cmd.UserSyncCmd).GetDatas()
					attrs := param.(*Cmd.UserSyncCmd).GetAttrs()
					g.Mutex.Lock()
					g.UpdateUserParams(datas, attrs)
					g.Mutex.Unlock()
				} else if param.(*Cmd.UserSyncCmd).GetType() == Cmd.EUserSyncType_EUSERSYNCTYPE_INIT {
					datas := param.(*Cmd.UserSyncCmd).GetDatas()
					attrs := param.(*Cmd.UserSyncCmd).GetAttrs()
					g.Mutex.Lock()
					g.UpdateUserParams(datas, attrs)
					g.Mutex.Unlock()
				}

			case Cmd.CmdParam_value["CHANGE_SCENE_USER_CMD"]:
				param = &Cmd.ChangeSceneUserCmd{}
				err = utils.ParseCmd(o, param)
				changeScene := param.(*Cmd.ChangeSceneUserCmd)
				if changeScene.GetMapID() != 0 {
					log.Infof("Moving %s to %v", g.Role.GetRoleName(), changeScene)
					g.MapNpcs = map[uint64]*Cmd.MapNpc{}
					g.MapUsers = map[uint64]*Cmd.MapUser{}
					g.Role.SetMapId(changeScene.GetMapID())
					g.Role.SetMapName(changeScene.GetMapName())
				}
				if changeScene != nil {
					g.Role.SetRolePos(changeScene.GetPos())
					inGame := true
					g.Role.InGame = &inGame
					if g.ShouldChangeScene {
						go func() {
							time.Sleep(time.Second)
							g.ChangeMap(g.Role.GetMapId())
						}()

					}
				}
			}

		case Cmd.Command_value["SCENE_USER_MAP_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.MapParam_value["MAPPARAM_ADDMAPUSER"]:
				param = &Cmd.AddMapUser{}
				err = utils.ParseCmd(o, param)
				addUsers := param.(*Cmd.AddMapUser)
				g.Mutex.Lock()
				for _, user := range addUsers.GetUsers() {
					g.MapUsers[user.GetGuid()] = user
				}
				g.Mutex.Unlock()

			case Cmd.MapParam_value["MAPPARAM_ADDMAPNPC"]:
				param = &Cmd.AddMapNpc{}
				err = utils.ParseCmd(o, param)
				addNpcs := param.(*Cmd.AddMapNpc)
				g.Mutex.Lock()
				for _, npc := range addNpcs.GetNpcs() {
					g.MapNpcs[npc.GetId()] = npc
				}
				g.Mutex.Unlock()

			case Cmd.MapParam_value["MAPPARAM_MAP_CMD_END"]:
				param = &Cmd.MapCmdEnd{}
				err = utils.ParseCmd(o, param)
			}

		case Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.RecordUserTradeParam_value["MY_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.MyPendingListRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				pRes := param.(*Cmd.MyPendingListRecordTradeCmd)
				if len(pRes.Lists) > 0 {
					g.Mutex.Lock()
					g.pendingSells = pRes
					g.Mutex.Unlock()
				}

			case Cmd.RecordUserTradeParam_value["SELL_ITEM_RECORDTRADE"]:
				param = &Cmd.SellItemRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				sellRes := param.(*Cmd.SellItemRecordTradeCmd)
				g.Mutex.Lock()
				g.sellItem[sellRes.GetItemInfo().GetItemid()] = sellRes
				g.Mutex.Unlock()

			case Cmd.RecordUserTradeParam_value["REQ_SERVER_PRICE_RECORDTRADE"]:
				param = &Cmd.ReqServerPriceRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				reqServerPrice := param.(*Cmd.ReqServerPriceRecordTradeCmd)
				if reqServerPrice.GetPrice() > 0 {
					g.Mutex.Lock()
					g.reqServerPrice[reqServerPrice.GetItemData().GetBase().GetId()] = reqServerPrice
					g.Mutex.Unlock()
				}

			case Cmd.RecordUserTradeParam_value["BUY_ITEM_RECORDTRADE"]:
				param = &Cmd.BuyItemRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				buyRes := param.(*Cmd.BuyItemRecordTradeCmd)
				g.Mutex.Lock()
				g.buyItem[buyRes.ItemInfo.GetItemid()] = buyRes
				g.Mutex.Unlock()

			case Cmd.RecordUserTradeParam_value["MY_TRADE_LOG_LIST_RECORDTRADE"]:
				param = &Cmd.MyTradeLogRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				history := param.(*Cmd.MyTradeLogRecordTradeCmd)
				if len(history.GetLogList()) > 0 {
					g.Mutex.Lock()
					g.tradeHistory = history
					g.Mutex.Unlock()
				}

			case Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.DetailPendingListRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				detail := param.(*Cmd.DetailPendingListRecordTradeCmd)
				if detail.GetSearchCond() != nil {
					g.Mutex.Lock()
					g.tradeDetail[detail.GetSearchCond().GetItemId()] = detail
					g.Mutex.Unlock()
				}

			case Cmd.RecordUserTradeParam_value["BRIEF_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.BriefPendingListRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				brief := param.(*Cmd.BriefPendingListRecordTradeCmd)
				if brief.GetCategory() != 0 {
					g.Mutex.Lock()
					g.tradeBrief[brief.GetCategory()] = brief
					g.Mutex.Unlock()
				}

			case Cmd.RecordUserTradeParam_value["ITEM_SELL_INFO_RECORDTRADE"]:
				param = &Cmd.ItemSellInfoRecordTradeCmd{}
				err = utils.ParseCmd(o, param)
				sellInfo := param.(*Cmd.ItemSellInfoRecordTradeCmd)
				if sellInfo.GetItemid() != 0 {
					g.Mutex.Lock()
					g.sellInfo[sellInfo.GetItemid()] = sellInfo
					g.Mutex.Unlock()
				}
			}

		case Cmd.Command_value["ERROR_USER_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.ErrCmdParam_value["REG_ERR_USER_CMD"]:
				param = &Cmd.RegErrUserCmd{}
				err = utils.ParseCmd(o, param)
				if param.(*Cmd.RegErrUserCmd).GetRet() == Cmd.RegErrRet_REG_ERR_DUPLICATE_LOGIN {
					log.Warnf("%s Account has been logged in on another device: %v", g.Role.GetRoleName(), param)
					g.quit <- true
				} else {
					log.Errorf("Server return err: %v", param)
					g.Close()
					return
				}

			case Cmd.ErrCmdParam_value["MAINTAIN_USER_CMD"]:
				param = &Cmd.MaintainUserCmd{}
				err = utils.ParseCmd(o, param)
				log.Warnf("Server is under maintanence: %v", param)
				os.Exit(1)
			}
		case Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEUPDATE"]:
				param = &Cmd.PackageUpdate{}
				err = utils.ParseCmd(o, param)
				packUpdate := param.(*Cmd.PackageUpdate)
				packType := packUpdate.GetType()
				g.Mutex.Lock()
				for _, item := range packUpdate.GetUpdateItems() {
					guid := item.GetBase().GetGuid()
					if g.Role.PackItems[packType] == nil {
						g.Role.PackItems[packType] = map[string]*Cmd.ItemData{}
					}
					g.Role.PackItems[packType][guid] = item
				}
				for _, item := range packUpdate.GetDelItems() {
					guid := item.GetBase().GetGuid()
					delete(g.Role.PackItems[packType], guid)
				}
				g.Mutex.Unlock()

			case Cmd.ItemParam_value["ITEMPARAM_BROWSEPACK"]:
				param = &Cmd.BrowsePackage{}
				err = utils.ParseCmd(o, param)

			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEITEM"]:
				param = &Cmd.PackageItem{}
				err = utils.ParseCmd(o, param)
				g.Mutex.Lock()
				items := param.(*Cmd.PackageItem)
				if len(items.GetData()) == 0 {
					continue
				} else {
					if g.Role.PackItems[items.GetType()] == nil {
						g.Role.PackItems[items.GetType()] = map[string]*Cmd.ItemData{}
					}
					for _, data := range items.GetData() {
						g.Role.PackItems[items.GetType()][data.GetBase().GetGuid()] = data
					}
				}
				g.Mutex.Unlock()
			}

		case Cmd.Command_value["SCENE_USER_SKILL_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.SkillParam_value["SKILLPARAM_SKILLITEM"]:
				param = &Cmd.ReqSkillData{}
				err = utils.ParseCmd(o, param)
				skillItems := param.(*Cmd.ReqSkillData)
				if len(skillItems.GetData()) > 0 {
					g.Mutex.Lock()
					for _, skillData := range skillItems.GetData() {
						for _, skillItem := range skillData.GetItems() {
							g.Role.SkillItems[skillItem.GetId()] = skillItem
							g.updateAutoSkill(skillItem)
						}
					}
					g.Mutex.Unlock()
				}

			case Cmd.SkillParam_value["SKILLPARAM_SKILLUPDATE"]:
				param = &Cmd.SkillUpdate{}
				err = utils.ParseCmd(o, param)
				skillUpdate := param.(*Cmd.SkillUpdate)
				//g.Mutex.Lock()
				for _, skillData := range skillUpdate.GetUpdate() {
					for _, newSkillItem := range skillData.GetItems() {
						if g.Role.SkillItems[newSkillItem.GetId()] != nil {
							g.Role.SkillItems[newSkillItem.GetId()] = newSkillItem
							g.updateAutoSkill(newSkillItem)
						}
					}
				}
				//g.Mutex.Unlock()

			default:
				continue
			}
		default:
			continue

		case Cmd.Command_value["FUBEN_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.FuBenParam_value["TEAMEXP_RAID_REPORT"]:
				param = &Cmd.TeamExpReportFubenCmd{}
				err = utils.ParseCmd(o, param)
				//report := param.(*Cmd.TeamExpReportFubenCmd)
				go func() {
					time.Sleep(10 * time.Second)
					g.ExitTeamExpFuben()
					time.Sleep(30 * time.Second)
					g.InviteTeamExpFuben()
				}()

			case Cmd.FuBenParam_value["TEAMEXP_QUERY_INFO"]:
				param = &Cmd.TeamExpQueryInfoFubenCmd{}
				err = utils.ParseCmd(o, param)
				queryInfo := param.(*Cmd.TeamExpQueryInfoFubenCmd)
				if g.notifier["TEAMEXP_QUERY_INFO"] != nil {
					g.notifier["TEAMEXP_QUERY_INFO"] <- queryInfo
				} else {
					g.Mutex.Lock()
					g.Role.TeamExpFubenInfo = queryInfo
					g.Mutex.Unlock()
				}

			default:
				continue
			}

		case Cmd.Command_value["MATCHC_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.MatchCParam_value["MATCHCPARAM_NTF_MATCHINFO"]:
				param = &Cmd.NtfMatchInfoCCmd{}
				err = utils.ParseCmd(o, param)
				// create team match info
				matchInfo := param.(*Cmd.NtfMatchInfoCCmd)
				if matchInfo.GetEtype() != Cmd.EPvpType_EPVPTYPE_MIN {
					tInfo := &Cmd.TeamPwsPreInfoMatchCCmd{}
					var pMember []*uint64
					ca := time.Now()
					if g.Role.MatchInfos[matchInfo.GetEtype()] != nil &&
						time.Since(g.Role.MatchInfos[matchInfo.GetEtype()].CreatedAt) < 60*time.Second {
						tInfo = g.Role.MatchInfos[matchInfo.GetEtype()].TeamPrepInfos
						pMember = g.Role.MatchInfos[matchInfo.GetEtype()].PrepedMember
						ca = g.Role.MatchInfos[matchInfo.GetEtype()].CreatedAt
					}

					detail := &utils.MatchDetail{
						MatchInfo:     matchInfo,
						TeamPrepInfos: tInfo,
						PrepedMember:  pMember,
						CreatedAt:     ca,
					}
					g.Role.MatchInfos[matchInfo.GetEtype()] = detail
					// Remove timed out invite after 60 seconds
					go func(matchType Cmd.EPvpType) {
						time.Sleep(60 * time.Second)
						if g.Role.MatchInfos[matchType] != nil && g.Role.MatchInfos[matchType].CreatedAt == ca {
							delete(g.Role.MatchInfos, matchType)
						}
					}(matchInfo.GetEtype())
				} else if matchInfo.GetEtype() == 108 && !matchInfo.GetIsmatch() {
					g.Role.MatchInfos = map[Cmd.EPvpType]*utils.MatchDetail{}
				}

			case Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_PREPARE_LIST"]:
				param = &Cmd.TeamPwsPreInfoMatchCCmd{}
				err = utils.ParseCmd(o, param)
				// create team prepare list
				prepList := param.(*Cmd.TeamPwsPreInfoMatchCCmd)
				if len(prepList.GetTeaminfos()) > 0 && g.Role.MatchInfos[prepList.GetEtype()] != nil {
					g.Role.MatchInfos[prepList.GetEtype()].TeamPrepInfos = prepList
				}

			case Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_PREPARE_UPDATE"]:
				param = &Cmd.UpdatePreInfoMatchCCmd{}
				err = utils.ParseCmd(o, param)
				// update team prepare list
				prepMember := param.(*Cmd.UpdatePreInfoMatchCCmd)
				g.Mutex.Lock()
				if prepMember.GetCharid() != 0 && g.Role.MatchInfos[prepMember.GetEtype()] != nil {
					g.Role.MatchInfos[prepMember.GetEtype()].PrepedMember = append(
						g.Role.MatchInfos[prepMember.GetEtype()].PrepedMember,
						prepMember.Charid,
					)
				}
				g.Mutex.Unlock()

			default:
				continue
			}

		case Cmd.Command_value["CHAT_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.ChatParam_value["CHATPARAM_CHAT_RET"]:
				param = &Cmd.ChatRetCmd{}
				err = utils.ParseCmd(o, param)
				chatRet := param.(*Cmd.ChatRetCmd)
				log.Infof("Receive chat from channel: %s, sender: %s content: %s", chatRet.GetChannel().String(), chatRet.GetName(), chatRet.GetStr())

			default:
				continue
			}

		case Cmd.Command_value["SESSION_USER_TEAM_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.TeamParam_value["TEAMPARAM_APPLYUPDATE"]:
				param = &Cmd.TeamApplyUpdate{}
				err = utils.ParseCmd(o, param)
				applyList := param.(*Cmd.TeamApplyUpdate)
				g.Mutex.Lock()
			applyLoop:
				for _, apply := range applyList.GetUpdates() {
					for _, cur := range g.Role.TeamApply {
						if apply.GetGuid() == cur.GetGuid() {
							cur = apply
							g.AcceptTeamApply(apply.GetGuid())
							continue applyLoop
						}
					}
					g.Role.TeamApply = append(g.Role.TeamApply, apply)
				}
				for _, del := range applyList.GetDeletes() {
					for i, cur := range g.Role.TeamApply {
						if cur.GetGuid() == del {
							g.Role.TeamApply = append(g.Role.TeamApply[:i], g.Role.TeamApply[i+1:]...)
						}
					}
				}
				g.Mutex.Unlock()

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERPOSUPDATE"]:
				param = &Cmd.MemberPosUpdate{}
				err = utils.ParseCmd(o, param)
				memberPos := param.(*Cmd.MemberPosUpdate)
				if memberPos.GetId() != 0 {
					g.Role.TeamMemberPos[memberPos.GetId()] = memberPos
				}

			case Cmd.TeamParam_value["TEAMPARAM_INVITEMEMBER"]:
				param = &Cmd.InviteMember{}
				err = utils.ParseCmd(o, param)
				invite := param.(*Cmd.InviteMember)
				if invite.GetUserguid() != 0 {
					go func() {
						time.Sleep(3 * time.Second)
						if len(g.GetCurrentTeamMembers()) == 1 {
							log.Infof("only one memeber in team exiting")
							g.ExitTeam()
						}
						log.Infof("accepting team invite from %s team name: %s",
							invite.GetUsername(), invite.GetTeamname())
						g.AcceptTeamInvite(invite.Userguid)
					}()
				}

			case Cmd.TeamParam_value["TEAMPARAM_ENTERTEAM"]:
				param = &Cmd.EnterTeam{}
				err = utils.ParseCmd(o, param)
				enterT := param.(*Cmd.EnterTeam)
				g.Mutex.Lock()
				g.Role.TeamData = enterT.Data
				g.Mutex.Unlock()

			case Cmd.TeamParam_value["TEAMPARAM_QUERYUSERTEAMINFO"]:
				param = &Cmd.QueryUserTeamInfoTeamCmd{}
				err = utils.ParseCmd(o, param)
				if _, ok := g.notifier["TEAMPARAM_QUERYUSERTEAMINFO"]; ok {
					g.notifier["TEAMPARAM_QUERYUSERTEAMINFO"] <- param.(*Cmd.QueryUserTeamInfoTeamCmd)
				}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERDATAUPDATE"]:
				param = &Cmd.MemberDataUpdate{}
				err = utils.ParseCmd(o, param)
				mData := param.(*Cmd.MemberDataUpdate)
				// Update Team member data
				g.Mutex.Lock()
				g.updateTeamMemberDatas(mData)
				g.Mutex.Unlock()

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERUPDATE"]:
				param = &Cmd.TeamMemberUpdate{}
				err = utils.ParseCmd(o, param)
				teamUpdate := param.(*Cmd.TeamMemberUpdate)
				g.Mutex.Lock()
				// Add new team member
				if len(teamUpdate.GetUpdates()) > 0 {
					for _, member := range teamUpdate.GetUpdates() {
						g.addTeamMember(member)
					}
				}
				// Delete team member
				if len(teamUpdate.GetDeletes()) > 0 {
					for _, del := range teamUpdate.GetDeletes() {
						g.removeTeamMember(del)
					}
				}
				g.Mutex.Unlock()

			default:
				continue
			}

		case Cmd.Command_value["SCENE_USER_QUEST_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.QuestParam_value["QUESTPARAM_QUESTLIST"]:
				param = &Cmd.QuestList{}
				err = utils.ParseCmd(o, param)
				ql := param.(*Cmd.QuestList)
				if g.notifier["QUESTPARAM_QUESTLIST"] != nil {
					g.notifier["QUESTPARAM_QUESTLIST"] <- ql
				} else {
					g.Role.QuestList = ql
				}

			default:
				continue
			}
		case Cmd.Command_value["SCENE_USER_PET_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYLIST"]:
				param = &Cmd.QueryPetAdventureListPetCmd{}
				err = utils.ParseCmd(o, param)
				if g.notifier["PETPARAM_ADVENTURE_QUERYLIST"] != nil {
					ql := param.(*Cmd.QueryPetAdventureListPetCmd)
					g.notifier["PETPARAM_ADVENTURE_QUERYLIST"] <- ql
				}

			case Cmd.PetParam_value["PETPARAM_WORK_QUERYWORKDATA"]:
				param = &Cmd.QueryPetWorkDataPetCmd{}
				err = utils.ParseCmd(o, param)
				if g.notifier["PETPARAM_WORK_QUERYWORKDATA"] != nil {
					workData := param.(*Cmd.QueryPetWorkDataPetCmd)
					g.notifier["PETPARAM_WORK_QUERYWORKDATA"] <- workData
				}

			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYBATTLEPET"]:
				param = &Cmd.QueryBattlePetCmd{}
				err = utils.ParseCmd(o, param)
				if g.notifier["PETPARAM_ADVENTURE_QUERYBATTLEPET"] != nil {
					battlePet := param.(*Cmd.QueryBattlePetCmd)
					g.notifier["PETPARAM_ADVENTURE_QUERYBATTLEPET"] <- battlePet
				}

			default:
				continue
			}
		case Cmd.Command_value["SESSION_USER_SOCIALITY_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.SocialityParam_value["SOCIALITYPARAM_FINDUSER"]:
				param = &Cmd.FindUser{}
				err = utils.ParseCmd(o, param)
				if _, ok := g.notifier["SOCIALITYPARAM_FINDUSER"]; ok {
					g.notifier["SOCIALITYPARAM_FINDUSER"] <- param.(*Cmd.FindUser)
				}

			default:
				continue
			}

		case Cmd.Command_value["INFINITE_TOWER_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.TowerParam_value["ETOWERPARAM_USERTOWERINFO"]:
				param = &Cmd.UserTowerInfoCmd{}
				err = utils.ParseCmd(o, param)
				towerInfo := param.(*Cmd.UserTowerInfoCmd)
				if g.notifier["ETOWERPARAM_USERTOWERINFO"] != nil {
					g.notifier["ETOWERPARAM_USERTOWERINFO"] <- towerInfo
				}
				g.Mutex.Lock()
				g.Role.UserTowerInfo = towerInfo.GetUsertower()
				g.Mutex.Unlock()

			case Cmd.TowerParam_value["ETOWERPARAM_TEAMTOWERSUMMARY"]:
				param = &Cmd.TeamTowerSummary{}
				err = utils.ParseCmd(o, param)
				if g.notifier["ETOWERPARAM_TEAMTOWERSUMMARY"] != nil {
					g.notifier["ETOWERPARAM_TEAMTOWERSUMMARY"] <- param.(*Cmd.TeamTowerSummary)
				}
			case Cmd.TowerParam_value["ETOWERPARAM_INVITE"]:
				go func() {
					time.Sleep(2 * time.Second)
					g.TowerReplyAgree()
				}()
			default:
				continue
			}
		}

		if err != nil && param == nil {
			log.Errorf("while handling command %s param %s encounter err: %s", cmdName, cmdParamName, err)
		} else {
			log.Debugf("cmd: %s, param: %s", cmdName, cmdParamName)
		}
	}
}

func (g *GameConnection) UpdateUserParams(datas []*Cmd.UserData, attrs []*Cmd.UserAttr) {
	for _, data := range datas {
		addData := false
		if data.GetType() == Cmd.EUserDataType_EUSERDATATYPE_SILVER {
			silver := data.GetValue()
			g.Role.Silver = &silver
			log.Infof("%s has %d zeny", g.Role.GetRoleName(), silver)
		}
		for _, d := range g.Role.UserDatas {
			if d.GetType() == data.GetType() {
				val := data.GetValue()
				d.Value = &val
				addData = true
			}
		}
		if !addData {
			g.Role.UserDatas = append(g.Role.UserDatas, data)
		}
	}
	for _, attr := range attrs {
		addData := false
		for _, a := range g.Role.UserAttrs {
			if a.GetType() == attr.GetType() {
				val := attr.GetValue()
				a.Value = &val
				addData = true
			}
		}
		if !addData {
			g.Role.UserAttrs = append(g.Role.UserAttrs, attr)
		}
	}
}
