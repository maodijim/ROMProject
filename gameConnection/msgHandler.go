package gameConnection

import (
	"os"
	"time"

	Cmd "ROMProject/Cmds"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	queryOnce = false
)

func (g *GameConnection) waitForResponse(notifierType gameTypes.NotifierType) (res interface{}, err error) {
	start := time.Now()
	for time.Since(start) < queryTimeout {
		select {
		case res = <-g.Notifier(notifierType):
			g.RemoveNotifier(notifierType)
			return res, err
		default:
			time.Sleep(time.Second)
		}
	}
	if time.Since(start) > queryTimeout {
		err = ErrQueryTimeout
	}
	g.RemoveNotifier(notifierType)
	return res, err
}

func (g *GameConnection) HandleMsg(output [][]byte) {
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
		case Cmd.Command_value["SESSION_USER_MAIL_PROTOCMD"]:
			_, _ = g.HandleSessionMailMsg(cmdParamId, o)

		case Cmd.Command_value["CENE_BOSS_PROTOCMD"]:
			_, _ = g.HandleSceneBossMsg(cmdParamId, o)

		case Cmd.Command_value["LOGIN_USER_PROTOCMD"]:
			_, _ = g.HandleLoginUserCmd(cmdParamId, o)

		case Cmd.Command_value["SCENE_USER2_PROTOCMD"]:
			_, _ = g.HandleSceneUser2ProtoCmd(cmdParamId, o)

		case Cmd.Command_value["SCENE_USER_PROTOCMD"]:
			_ = g.HandleSceneUserProtoCmd(cmdParamId, o)

		case Cmd.Command_value["SCENE_USER_MAP_PROTOCMD"]:
			_ = g.HandleSceneUserMapProtoCmd(cmdParamId, o)

		case Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"]:
			_ = g.HandleRecordUserTradeProtoCmd(cmdParamId, o)

		case Cmd.Command_value["ERROR_USER_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.ErrCmdParam_value["REG_ERR_USER_CMD"]:
				param = &Cmd.RegErrUserCmd{}
				err = utils.ParseCmd(o, param)
				returnCode := param.(*Cmd.RegErrUserCmd).GetRet()
				if returnCode == Cmd.RegErrRet_REG_ERR_DUPLICATE_LOGIN {
					log.Warnf("%s Account has been logged in on another device: %v", g.Role.GetRoleName(), param)
					g.Close()
				} else if returnCode == Cmd.RegErrRet_REG_ERR_ACC_FORBID {
					log.Errorf("%s Account forbidden", g.Role.GetRoleName())
					g.Close()
				} else {
					log.Errorf("Server return err: %v", param)
					g.Close()
				}
				return

			case Cmd.ErrCmdParam_value["MAINTAIN_USER_CMD"]:
				param = &Cmd.MaintainUserCmd{}
				err = utils.ParseCmd(o, param)
				log.Warnf("Server is under maintanence: %v", param)
				os.Exit(1)
			}
		case Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.ItemParam_value["ITEMPARAM_NTF_HIGHTREFINE_DATA"]:
				param = &Cmd.NtfHighRefineDataCmd{}
				err = utils.ParseCmd(o, param)

			case Cmd.ItemParam_value["ITEMPARAM_PACKSLOTNTF"]:
				param = &Cmd.PackSlotNtfItemCmd{}
				err = utils.ParseCmd(o, param)

			case Cmd.ItemParam_value["ITEMPARAM_ITEMSHOW"]:
				param = &Cmd.ItemShow{}
				err = utils.ParseCmd(o, param)

			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEUPDATE"]:
				param = &Cmd.PackageUpdate{}
				err = utils.ParseCmd(o, param)
				packUpdate := param.(*Cmd.PackageUpdate)
				packType := packUpdate.GetType()
				g.Role.Mutex.Lock()
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
				g.Role.Mutex.Unlock()

			case Cmd.ItemParam_value["ITEMPARAM_BROWSEPACK"]:
				param = &Cmd.BrowsePackage{}
				err = utils.ParseCmd(o, param)

			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEITEM"]:
				param = &Cmd.PackageItem{}
				err = utils.ParseCmd(o, param)
				g.Role.Mutex.Lock()
				items := param.(*Cmd.PackageItem)
				if len(items.GetData()) == 0 {
					g.Role.Mutex.Unlock()
					continue
				} else {
					if g.Role.PackItems[items.GetType()] == nil {
						g.Role.PackItems[items.GetType()] = map[string]*Cmd.ItemData{}
					}
					for _, data := range items.GetData() {
						g.Role.PackItems[items.GetType()][data.GetBase().GetGuid()] = data
					}
				}
				g.Role.Mutex.Unlock()
			}

		case Cmd.Command_value["SCENE_USER_SKILL_PROTOCMD"]:
			switch cmdParamId {

			case Cmd.SkillParam_value["SKILLPARAM_SKILLITEM"]:
				param = &Cmd.ReqSkillData{}
				err = utils.ParseCmd(o, param)
				skillItems := param.(*Cmd.ReqSkillData)
				if len(skillItems.GetData()) > 0 {
					g.Role.Mutex.Lock()
					for _, skillData := range skillItems.GetData() {
						for _, skillItem := range skillData.GetItems() {
							g.Role.SkillItems[skillItem.GetId()] = skillItem
							g.updateAutoSkill(skillItem)
						}
					}
					g.Role.Mutex.Unlock()
				}

			case Cmd.SkillParam_value["SKILLPARAM_SKILLUPDATE"]:
				param = &Cmd.SkillUpdate{}
				err = utils.ParseCmd(o, param)
				skillUpdate := param.(*Cmd.SkillUpdate)
				g.Role.Mutex.Lock()
				for _, skillData := range skillUpdate.GetUpdate() {
					for _, newSkillItem := range skillData.GetItems() {
						if g.Role.SkillItems[newSkillItem.GetId()] != nil {
							g.Role.SkillItems[newSkillItem.GetId()] = newSkillItem
							g.updateAutoSkill(newSkillItem)
						}
					}
				}
				g.Role.Mutex.Unlock()

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
				// report := param.(*Cmd.TeamExpReportFubenCmd)
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
				if g.Notifier(gameTypes.NtfType_TeamExpQueryInfo) != nil {
					g.Notifier(gameTypes.NtfType_TeamExpQueryInfo) <- queryInfo
				} else {
					g.Role.Mutex.Lock()
					g.Role.TeamExpFubenInfo = queryInfo
					g.Role.Mutex.Unlock()
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
				g.Role.Mutex.Lock()
				if prepMember.GetCharid() != 0 && g.Role.MatchInfos[prepMember.GetEtype()] != nil {
					g.Role.MatchInfos[prepMember.GetEtype()].PrepedMember = append(
						g.Role.MatchInfos[prepMember.GetEtype()].PrepedMember,
						prepMember.Charid,
					)
				}
				g.Role.Mutex.Unlock()

			default:
				continue
			}

		case Cmd.Command_value["CHAT_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.ChatParam_value["CHATPARAM_CHAT_RET"]:
				param = &Cmd.ChatRetCmd{}
				err = utils.ParseCmd(o, param)
				chatRet := param.(*Cmd.ChatRetCmd)
				log.Infof(
					"Receive chat from channel: %s, sender id: %d sender: %s content: %s",
					chatRet.GetChannel().String(),
					chatRet.GetId(),
					chatRet.GetName(),
					chatRet.GetStr(),
				)

			default:
				continue
			}

		case Cmd.Command_value["SESSION_USER_TEAM_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.TeamParam_value["TEAMPARAM_APPLYUPDATE"]:
				param = &Cmd.TeamApplyUpdate{}
				err = utils.ParseCmd(o, param)
				applyList := param.(*Cmd.TeamApplyUpdate)
				g.Role.Mutex.Lock()
				for _, apply := range applyList.GetUpdates() {
					for _, cur := range g.Role.TeamApply {
						if apply.GetGuid() == cur.GetGuid() {
							cur = apply
							g.AcceptTeamApply(apply.GetGuid())
						} else if g.Role.AcceptAllTeamInvite {
							go func() {
								time.Sleep(time.Second * 5)
								g.AcceptTeamApply(apply.GetGuid())
							}()
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
				g.Role.Mutex.Unlock()

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
				g.Role.Mutex.Lock()
				g.Role.TeamData = enterT.Data
				g.Role.Mutex.Unlock()

			case Cmd.TeamParam_value["TEAMPARAM_QUERYUSERTEAMINFO"]:
				param = &Cmd.QueryUserTeamInfoTeamCmd{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_TeamParamQueryUserTeamInfo) != nil {
					g.Notifier(gameTypes.NtfType_TeamParamQueryUserTeamInfo) <- param.(*Cmd.QueryUserTeamInfoTeamCmd)
				}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERDATAUPDATE"]:
				param = &Cmd.MemberDataUpdate{}
				err = utils.ParseCmd(o, param)
				mData := param.(*Cmd.MemberDataUpdate)
				// Update Team member data
				g.updateTeamMemberDatas(mData)

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERUPDATE"]:
				param = &Cmd.TeamMemberUpdate{}
				err = utils.ParseCmd(o, param)
				teamUpdate := param.(*Cmd.TeamMemberUpdate)
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

			default:
				continue
			}

		case Cmd.Command_value["SCENE_USER_QUEST_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.QuestParam_value["QUESTPARAM_QUESTUPDATE"]:
				param = &Cmd.QuestUpdate{}
				err = utils.ParseCmd(o, param)
				questUpdate := param.(*Cmd.QuestUpdate)
				g.Role.Mutex.Lock()
				for _, quests := range questUpdate.GetItems() {
					tempQuestList := g.Role.GetQuestList(quests.GetType()).GetList()
					for _, update := range quests.GetUpdate() {
						for i, quest := range g.Role.GetQuestList(quests.GetType()).GetList() {
							if quest.GetId() == update.GetId() {
								tempQuestList[i] = update
							}
						}
					}
					for _, del := range quests.GetDel() {
						for i, quest := range g.Role.GetQuestList(quests.GetType()).GetList() {
							if quest.GetId() == del {
								tempQuestList = append(tempQuestList[:i], tempQuestList[i+1:]...)
							}
						}
					}
					if g.Role.QuestList[quests.GetType()] == nil {
						g.Role.QuestList[quests.GetType()] = &Cmd.QuestList{}
						g.Role.QuestList[quests.GetType()].List = tempQuestList
					} else {
						g.Role.QuestList[quests.GetType()].List = tempQuestList
					}
				}
				g.Role.Mutex.Unlock()

			case Cmd.QuestParam_value["QUESTPARAM_QUESTLIST"]:
				param = &Cmd.QuestList{}
				err = utils.ParseCmd(o, param)
				ql := param.(*Cmd.QuestList)
				if g.Notifier(gameTypes.NtfType_QuestList) != nil {
					g.Notifier(gameTypes.NtfType_QuestList) <- ql
				} else {
					g.Role.QuestList[ql.GetType()] = ql
				}

			default:
				continue
			}
		case Cmd.Command_value["SCENE_USER_PET_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYLIST"]:
				param = &Cmd.QueryPetAdventureListPetCmd{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_PetAdventureQueryList) != nil {
					ql := param.(*Cmd.QueryPetAdventureListPetCmd)
					g.Notifier(gameTypes.NtfType_PetAdventureQueryList) <- ql
				}

			case Cmd.PetParam_value["PETPARAM_WORK_QUERYWORKDATA"]:
				param = &Cmd.QueryPetWorkDataPetCmd{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_PetQueryWorkData) != nil {
					workData := param.(*Cmd.QueryPetWorkDataPetCmd)
					g.Notifier(gameTypes.NtfType_PetQueryWorkData) <- workData
				}

			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYBATTLEPET"]:
				param = &Cmd.QueryBattlePetCmd{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_PetQueryBattlePet) != nil {
					battlePet := param.(*Cmd.QueryBattlePetCmd)
					g.Notifier(gameTypes.NtfType_PetQueryBattlePet) <- battlePet
				}

			default:
				continue
			}
		case Cmd.Command_value["SESSION_USER_SOCIALITY_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.SocialityParam_value["SOCIALITYPARAM_FINDUSER"]:
				param = &Cmd.FindUser{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_SocialityFindUser) != nil {
					g.Notifier(gameTypes.NtfType_SocialityFindUser) <- param.(*Cmd.FindUser)
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
				if g.Notifier(gameTypes.NtfType_TowerUserInfo) != nil {
					g.Notifier(gameTypes.NtfType_TowerUserInfo) <- towerInfo
				}
				g.Role.Mutex.Lock()
				g.Role.UserTowerInfo = towerInfo.GetUsertower()
				g.Role.Mutex.Unlock()

			case Cmd.TowerParam_value["ETOWERPARAM_TEAMTOWERSUMMARY"]:
				param = &Cmd.TeamTowerSummary{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_TowerTeamSummary) != nil {
					g.Notifier(gameTypes.NtfType_TowerTeamSummary) <- param.(*Cmd.TeamTowerSummary)
				}
			case Cmd.TowerParam_value["ETOWERPARAM_INVITE"]:
				go func() {
					time.Sleep(2 * time.Second)
					g.TowerReplyAgree()
				}()
			default:
				continue
			}
		case Cmd.Command_value["SCENE_USER_INTER_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.InterParam_value["INTERPARAM_NEWINTERLOCUTION"]:
				param = &Cmd.NewInter{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_InterviewQuestion) != nil {
					go func() {
						g.Notifier(gameTypes.NtfType_InterviewQuestion) <- param.(*Cmd.NewInter)
					}()
				}
			}
		case Cmd.Command_value["SESSION_USER_SHOP_PROTOCMD"]:
			switch cmdParamId {
			case Cmd.ShopParam_value["SHOPPARAM_QUERY_SHOP_CONFIG"]:
				param = &Cmd.QueryShopConfigCmd{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_ShopQueryShopConfig) != nil {
					g.Notifier(gameTypes.NtfType_ShopQueryShopConfig) <- param.(*Cmd.QueryShopConfigCmd)
				}
			case Cmd.ShopParam_value["SHOPPARAM_BUYITEM"]:
				param = &Cmd.BuyShopItem{}
				err = utils.ParseCmd(o, param)
				if g.Notifier(gameTypes.NtfType_ShopBuyItem) != nil {
					g.Notifier(gameTypes.NtfType_ShopBuyItem) <- param.(*Cmd.BuyShopItem)
				}
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
		} else if data.GetType() == Cmd.EUserDataType_EUSERDATATYPE_LOTTERY {
			lottery := data.GetValue()
			g.Role.Lottery = &lottery
			log.Infof("%s has %d lottery", g.Role.GetRoleName(), lottery)
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

func (g *GameConnection) setLastHeartBeat() {
	if g.conn != nil {
		g.currentIndex = 1
		g.ShouldHeartBeat = true
		g.lastHeartBeat = time.Now()
	}
}

func (g *GameConnection) setCurrentMsgIndex(index uint32) {
	g.currentIndex = index
}
