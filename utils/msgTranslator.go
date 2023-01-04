package utils

import (
	"strings"

	Cmd "ROMProject/Cmds"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func TranslateMsg(output [][]byte) {
	for _, o := range output {
		log.Printf("%s", strings.Repeat("-", 50))
		if len(o) < 2 {
			log.Printf("result is empty")
			continue
		}
		cmdId := int32(o[0])
		cmdName := Cmd.Command_name[cmdId]
		cmdParamId := int32(o[1])
		cmdParamName := NameParamMap[cmdName][cmdParamId]
		log.Printf("command id: %d; command name: %s; command param id: %d; param name: %s",
			cmdId, cmdName, cmdParamId, cmdParamName,
		)
		var param proto.Message
		var err error
		if cmdId == Cmd.Command_value["LOGIN_USER_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.LoginCmdParam_value["LOGIN_RESULT_USER_CMD"]:
				param = &Cmd.LoginResultUserCmd{}

			case Cmd.LoginCmdParam_value["CLIENT_FRAME_USER_CMD"]:
				param = &Cmd.ClientFrameUserCmd{}

			case Cmd.LoginCmdParam_value["CLIENT_INFO_USER_CMD"]:
				param = &Cmd.ClientInfoUserCmd{}

			case Cmd.LoginCmdParam_value["SELECT_ROLE_USER_CMD"]:
				param = &Cmd.SelectRoleUserCmd{}

			case Cmd.LoginCmdParam_value["REQ_LOGIN_USER_CMD"]:
				param = &Cmd.ReqLoginUserCmd{}

			case Cmd.LoginCmdParam_value["SERVERTIME_USER_CMD"]:
				param = &Cmd.ServerTimeUserCmd{}

			case Cmd.LoginCmdParam_value["REAL_AUTHORIZE_USER_CMD"]:
				param = &Cmd.RealAuthorizeUserCmd{}

			case Cmd.LoginCmdParam_value["HEART_BEAT_USER_CMD"]:
				param = &Cmd.HeartBeatUserCmd{}

			case Cmd.LoginCmdParam_value["SNAPSHOT_USER_CMD"]:
				param = &Cmd.SnapShotUserCmd{}

			case Cmd.LoginCmdParam_value["SAFE_DEVICE_USER_CMD"]:
				param = &Cmd.SafeDeviceUserCmd{}

			case Cmd.LoginCmdParam_value["CONFIRM_AUTHORIZE_USER_CMD"]:
				param = &Cmd.ConfirmAuthorizeUserCmd{}

			case Cmd.LoginCmdParam_value["REQ_LOGIN_PARAM_USER_CMD"]:
				param = &Cmd.ReqLoginParamUserCmd{}

			case Cmd.LoginCmdParam_value["CREATE_CHAR_USER_CMD"]:
				param = &Cmd.CreateCharUserCmd{}

			case Cmd.LoginCmdParam_value["DELETE_CHAR_USER_CMD"]:
				param = &Cmd.DeleteCharUserCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}

			if param != nil {
				err = ParseCmd(o, param)
			}
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["RECORD_USER_TRADE_PROTOCMD"] {
			switch cmdParamId {

			case Cmd.RecordUserTradeParam_value["TAKE_LOG_TRADE_PARAM"]:
				param = &Cmd.TakeLogCmd{}

			case Cmd.RecordUserTradeParam_value["NTF_CAN_TAKE_COUNT_TRADE_PARAM"]:
				param = &Cmd.NtfCanTakeCountTradeCmd{}

			case Cmd.RecordUserTradeParam_value["LIST_NTF_RECORDTRADE"]:
				param = &Cmd.ListNtfRecordTrade{}

			case Cmd.RecordUserTradeParam_value["BRIEF_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.BriefPendingListRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["ITEM_SELL_INFO_RECORDTRADE"]:
				param = &Cmd.ItemSellInfoRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["REQ_SERVER_PRICE_RECORDTRADE"]:
				param = &Cmd.ReqServerPriceRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["PANEL_RECORDTRADE"]:
				param = &Cmd.PanelRecordTrade{}

			case Cmd.RecordUserTradeParam_value["MY_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.MyPendingListRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["MY_TRADE_LOG_LIST_RECORDTRADE"]:
				param = &Cmd.MyTradeLogRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["DETAIL_PENDING_LIST_RECORDTRADE"]:
				param = &Cmd.DetailPendingListRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["BUY_ITEM_RECORDTRADE"]:
				param = &Cmd.BuyItemRecordTradeCmd{}

			case Cmd.RecordUserTradeParam_value["SELL_ITEM_RECORDTRADE"]:
				param = &Cmd.SellItemRecordTradeCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SCENE_USER2_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.User2Param_value["USER2PARAM_SIGNIN"]:
				param = &Cmd.SignInUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SERVANT_RECEIVE"]:
				param = &Cmd.ReceiveServantUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SERVANT_REWARD_STATUS"]:
				param = &Cmd.ServantRewardStatusUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SERVANT_RECOMMEND"]:
				param = &Cmd.RecommendServantUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_INVITEFOLLOW"]:
				param = &Cmd.InviteFollowUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_FOLLOWER"]:
				param = &Cmd.FollowerUser{}

			case Cmd.User2Param_value["USER2PARAM_GOMAP_FOLLOW"]:
				param = &Cmd.GoMapFollowUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_GAMEHEALTH_UPDATE"]:
				param = &Cmd.UpdateGameHealthLevelUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_PRESETCHATMSG"]:
				param = &Cmd.PresetMsgCmd{}

			case Cmd.User2Param_value["USER2PARAM_PRESTIGE_NTF"]:
				param = &Cmd.PrestigeNtfUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_GOTO_GEAR"]:
				param = &Cmd.GoToGearUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SET_DIRECTION"]:
				param = &Cmd.SetDirection{}

			case Cmd.User2Param_value["USER2PARAM_GAMETIME"]:
				param = &Cmd.GameTimeCmd{}

			case Cmd.User2Param_value["USER2PARAM_EXIT_POS"]:
				param = &Cmd.ExitPosUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_CHEAT_TAG_STAT"]:
				param = &Cmd.CheatTagStatUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_READYTOMAP"]:
				param = &Cmd.ReadyToMapUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_CDTIME"]:
				param = &Cmd.CDTimeUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SETOPTION"]:
				param = &Cmd.SetOptionUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_QUERYPORTRAITLIST"]:
				param = &Cmd.QueryPortraitListUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_GOTO_LIST"]:
				param = &Cmd.GoToListUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_QUERYFIGHTERINFO"]:
				param = &Cmd.QueryFighterInfo{}

			case Cmd.User2Param_value["USER2PARAM_SERVANT_STATISTICS"]:
				param = &Cmd.UpdateBranchInfoUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_UPDATE_BRANCH_INFO"]:
				param = &Cmd.UpdateBranchInfoUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_UPDATE_RECORD_INFO"]:
				param = &Cmd.UpdateRecordInfoUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_QUERY_TRACE_LIST"]:
				param = &Cmd.QueryTraceList{}

			case Cmd.User2Param_value["USER2PARAM_VAR"]:
				param = &Cmd.VarUpdate{}

			case Cmd.User2Param_value["USER2PARAM_MENU"]:
				param = &Cmd.MenuList{}

			case Cmd.User2Param_value["USER2PARAM_BUFFERSYNC"]:
				param = &Cmd.UserBuffNineSyncCmd{}

			case Cmd.User2Param_value["USER2PARAM_QUERY_ACTION"]:
				param = &Cmd.QueryShow{}

			case Cmd.User2Param_value["USER2PARAM_BATTLE_TIMELEN_USER_CMD"]:
				param = &Cmd.BattleTimelenUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_SYSMSG"]:
				param = &Cmd.SysMsg{}

			case Cmd.User2Param_value["USER2PARAM_SERVER_INFO_NTF"]:
				param = &Cmd.ServerInfoNtf{}

			case Cmd.User2Param_value["USER2PARAM_SIGNIN_NTF"]:
				param = &Cmd.SignInNtfUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_ACTION"]:
				param = &Cmd.UserActionNtf{}

			case Cmd.User2Param_value["USER2PARAM_QUERY_MAPAREA"]:
				param = &Cmd.QueryMapArea{}

			case Cmd.User2Param_value["USER2PARAM_NPCDATASYNC"]:
				param = &Cmd.NpcDataSync{}

			case Cmd.User2Param_value["USER2PARAM_QUERYSHOPGOTITEM"]:
				param = &Cmd.QueryShopGotItem{}

			case Cmd.User2Param_value["USER2PARAM_QUERYSHORTCUT"]:
				param = &Cmd.QueryShortcut{}

			case Cmd.User2Param_value["USER2PARAM_USERNINESYNC"]:
				param = &Cmd.UserNineSyncCmd{}

			case Cmd.User2Param_value["USER2PARAM_NEW_SET_OPTION"]:
				param = &Cmd.NewSetOptionUserCmd{}

			case Cmd.User2Param_value["USER2PARAM_TALKINFO"]:
				param = &Cmd.TalkInfo{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SCENE_USER_ITEM_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ItemParam_value["ITEMPARAM_EQUIP"]:
				param = &Cmd.Equip{}

			case Cmd.ItemParam_value["ITEMPARAM_SELLITEM"]:
				param = &Cmd.SellItem{}

			case Cmd.ItemParam_value["ITEMPARAM_QUICK_SELLITEM"]:
				param = &Cmd.QuickSellItemCmd{}

			case Cmd.User2Param_value["USER2PARAM_SERVANT_GROWTH"]:
				param = &Cmd.GrowthServantUserCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_EQUIPSTRENGTH"]:
				param = &Cmd.EquipStrength{}

			case Cmd.ItemParam_value["ITEMPARAM_ITEMUSE"]:
				param = &Cmd.ItemUse{}

			case Cmd.ItemParam_value["ITEMPARAM_GETCOUNT"]:
				param = &Cmd.GetCountItemCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_ONOFFSTORE"]:
				param = &Cmd.OnOffStoreItemCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_NTF_HIGHTREFINE_DATA"]:
				param = &Cmd.NtfHighRefineDataCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_EQUIPPOSDATA_UPDATE"]:
				param = &Cmd.EquipPosDataUpdate{}

			case Cmd.ItemParam_value["ITEMPARAM_PACKSLOTNTF"]:
				param = &Cmd.PackSlotNtfItemCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_QUERY_ITEMDEBT"]:
				param = &Cmd.QueryDebtItemCmd{}

			case Cmd.ItemParam_value["ITEMPARAM_BROWSEPACK"]:
				param = &Cmd.BrowsePackage{}

			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEITEM"]:
				param = &Cmd.PackageItem{}

			case Cmd.ItemParam_value["ITEMPARAM_PACKAGEUPDATE"]:
				param = &Cmd.PackageUpdate{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SCENE_USER_MANUAL_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ManualParam_value["MANUALPARAM_UPDATE"]:
				param = &Cmd.ManualUpdate{}

			case Cmd.ManualParam_value["MANUALPARAM_LEVELSYNC"]:
				param = &Cmd.LevelSync{}

			case Cmd.ManualParam_value["MANUALPARAM_SKILLPOINTSYNC"]:
				param = &Cmd.SkillPointSync{}

			case Cmd.ManualParam_value["MANUALPARAM_QUERYVERSION"]:
				param = &Cmd.QueryVersion{}

			case Cmd.ManualParam_value["MANUALPARAM_NPCZONE"]:
				param = &Cmd.NpcZoneDataManualCmd{}

			case Cmd.ManualParam_value["MANUALPARAM_QUERYDATA"]:
				param = &Cmd.QueryManualData{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SCENE_USER_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.CmdParam_value["MAP_OBJECT_DATA"]:
				param = &Cmd.MapObjectData{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["GOTO_USER_CMD"]:
				param = &Cmd.GoToUserCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["SKILL_BROADCAST_USER_CMD"]:
				param = &Cmd.SkillBroadcastUserCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["REQ_MOVE_USER_CMD"]:
				param = &Cmd.ReqMoveUserCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["CHANGE_SCENE_USER_CMD"]:
				param = &Cmd.ChangeSceneUserCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["USERPARAM_USERSYNC"]:
				param = &Cmd.UserSyncCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["RET_MOVE_USER_CMD"]:
				param = &Cmd.RetMoveUserCmd{}
				err = ParseCmd(o, param)

			case Cmd.CmdParam_value["DELETE_ENTRY_USER_CMD"]:
				param = &Cmd.DeleteEntryUserCmd{}
				err = ParseCmd(o, param)

			default:
				log.Infof("没有parsing")
				continue
			}

			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["ERROR_USER_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ErrCmdParam_value["REG_ERR_USER_CMD"]:
				param = &Cmd.ChangeSceneUserCmd{}
				err = proto.Unmarshal(o[2:], param)

			case Cmd.ErrCmdParam_value["MAINTAIN_USER_CMD"]:
				param = &Cmd.MaintainUserCmd{}
				err = ParseCmd(o, param)

			default:
				log.Infof("没有parsing")
				continue
			}

			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SCENE_USER_SKILL_PROTOCMD"] {
			switch cmdParamId {

			case Cmd.SkillParam_value["SKILLPARAM_SKILLVALIDPOS"]:
				param = &Cmd.SkillValidPos{}

			case Cmd.SkillParam_value["SKILLPARAM_EQUIPSKILL"]:
				param = &Cmd.EquipSkill{}

			case Cmd.SkillParam_value["SKILLPARAM_SKILLITEM"]:
				param = &Cmd.ReqSkillData{}

			case Cmd.SkillParam_value["SKILLPARAM_SPEC_SKILL_INFO"]:
				param = &Cmd.UpSkillInfoSkillCmd{}

			case Cmd.SkillParam_value["SKILLPARAM_SKILLUPDATE"]:
				param = &Cmd.SkillUpdate{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SCENE_USER_ACHIEVE_PROTOCMD"] {
			switch cmdParamId {
			default:
				log.Infof("没有parsing")
				continue

			case Cmd.AchieveParam_value["ACHIEVEPARAM_NEW_ACHNTF"]:
				param = &Cmd.QueryAchieveDataAchCmd{}
				err = ParseCmd(o, param)

			case Cmd.AchieveParam_value["ACHIEVEPARAM_QUERY_ACHDATA"]:
				param = &Cmd.QueryAchieveDataAchCmd{}
				err = ParseCmd(o, param)
			}

			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SCENE_USER_MAP_PROTOCMD"] {
			switch cmdParamId {
			default:
				log.Infof("没有parsing")
				continue

			case Cmd.MapParam_value["MAPPARAM_ADDMAPITEM"]:
				param = &Cmd.AddMapItem{}

			case Cmd.MapParam_value["MAPPARAM_ADDMAPTRAP"]:
				param = &Cmd.AddMapTrap{}

			case Cmd.MapParam_value["MAPPARAM_EXIT_POINT_STATE"]:
				param = &Cmd.ExitPointState{}

			case Cmd.MapParam_value["MAPPARAM_PICKUPITEM"]:
				param = &Cmd.PickupItem{}

			case Cmd.MapParam_value["MAPPARAM_ADDMAPUSER"]:
				param = &Cmd.AddMapUser{}

			case Cmd.MapParam_value["MAPPARAM_ADDMAPNPC"]:
				param = &Cmd.AddMapNpc{}

			case Cmd.MapParam_value["MAPPARAM_MAP_CMD_END"]:
				param = &Cmd.MapCmdEnd{}

			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["USER_EVENT_PROTOCMD"] {
			switch cmdParamId {
			default:
				log.Infof("没有parsing")
				continue
			case Cmd.EventParam_value["USER_EVENT_QUERY_CHARGE_CNT"]:
				param = &Cmd.QueryChargeCnt{}

			case Cmd.EventParam_value["USER_EVENT_AUTOBATTLE"]:
				param = &Cmd.SwitchAutoBattleUserEvent{}

			case Cmd.EventParam_value["USER_EVENT_NTF_VERSION_CARD"]:
				param = &Cmd.NtfVersionCardInfo{}

			case Cmd.EventParam_value["USER_EVENT_ATTACK_NPC"]:
				param = &Cmd.DamageNpcUserEvent{}

			case Cmd.EventParam_value["USER_EVENT_FIRST_ACTION"]:
				param = &Cmd.FirstActionUserEvent{}

			case Cmd.EventParam_value["USER_EVENT_UPDATE_RANDOM"]:
				param = &Cmd.UpdateRandomUserEvent{}

			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SESSION_USER_SHOP_PROTOCMD"] {
			switch cmdParamId {
			default:
				log.Infof("没有parsing")
				continue

			case Cmd.ShopParam_value["SHOPPARAM_BUYITEM"]:
				param = &Cmd.BuyShopItem{}

			case Cmd.ShopParam_value["SHOPPARAM_QUERY_SHOP_CONFIG"]:
				param = &Cmd.QueryShopConfigCmd{}

			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["CHAT_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ChatParam_value["CHATPARAM_SYSTEM_BARRAGE"]:
				param = &Cmd.SystemBarrageChatCmd{}

			case Cmd.ChatParam_value["CHATPARAM_CHAT_RET"]:
				param = &Cmd.ChatRetCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SESSION_USER_TEAM_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.TeamParam_value["TEAMPARAM_MEMBERAPPLY"]:
				param = &Cmd.TeamMemberApply{}

			case Cmd.TeamParam_value["TEAMPARAM_PROCESSAPPLY"]:
				param = &Cmd.ProcessTeamApply{}

			case Cmd.TeamParam_value["TEAMPARAM_APPLYUPDATE"]:
				param = &Cmd.TeamApplyUpdate{}

			case Cmd.TeamParam_value["TEAMPARAM_TEAMLIST"]:
				param = &Cmd.TeamList{}

			case Cmd.TeamParam_value["TEAMPARAM_CREATETEAM"]:
				param = &Cmd.CreateTeam{}

			case Cmd.TeamParam_value["TEAMPARAM_EXCHANGELEADER"]:
				param = &Cmd.ExchangeLeader{}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERPOSUPDATE"]:
				param = &Cmd.MemberPosUpdate{}

			case Cmd.TeamParam_value["TEAMPARAM_PROCESSINVITE"]:
				param = &Cmd.ProcessTeamInvite{}

			case Cmd.TeamParam_value["TEAMPARAM_QUERYMEMBERCAT"]:
				param = &Cmd.QueryMemberCatTeamCmd{}

			case Cmd.TeamParam_value["TEAMPARAM_ENTERTEAM"]:
				param = &Cmd.EnterTeam{}

			case Cmd.TeamParam_value["TEAMPARAM_EXITTEAM"]:
				param = &Cmd.ExitTeam{}

			case Cmd.TeamParam_value["TEAMPARAM_INVITEMEMBER"]:
				param = &Cmd.InviteMember{}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERUPDATE"]:
				param = &Cmd.TeamMemberUpdate{}

			case Cmd.TeamParam_value["TEAMPARAM_QUERYUSERTEAMINFO"]:
				param = &Cmd.QueryUserTeamInfoTeamCmd{}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERCAT_UPDATE"]:
				param = &Cmd.MemberCatUpdateTeam{}

			case Cmd.TeamParam_value["TEAMPARAM_DATAUPDATE"]:
				param = &Cmd.TeamDataUpdate{}

			case Cmd.TeamParam_value["TEAMPARAM_MEMBERDATAUPDATE"]:
				param = &Cmd.MemberDataUpdate{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SESSION_USER_AUTHORIZE_PROTOCMD"] {
			switch cmdParamId {

			case Cmd.AuthorizeParam_value["SET_AUTHORIZE_USER_CMD"]:
				param = &Cmd.SetAuthorizeUserCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SCENE_USER_QUEST_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.QuestParam_value["QUESTPARAM_QUESTACTION"]:
				param = &Cmd.QuestAction{}

			case Cmd.QuestParam_value["QUESTPARAM_QUESTSTEPUPDATE"]:
				param = &Cmd.QuestStepUpdate{}

			case Cmd.QuestParam_value["QUESTPARAM_QUESTLIST"]:
				param = &Cmd.QuestList{}

			case Cmd.QuestParam_value["QUESTPARAM_QUESTUPDATE"]:
				param = &Cmd.QuestUpdate{}

			case Cmd.QuestParam_value["QUESTPARAM_VISIT_NPC"]:
				param = &Cmd.VisitNpcUserCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["FUBEN_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.FuBenParam_value["MONSTER_COUNT_USER_CMD"]:
				param = &Cmd.MonsterCountUserCmd{}

			case Cmd.FuBenParam_value["BEGIN_FIRE_FUBENCMD"]:
				param = &Cmd.BeginFireFubenCmd{}

			case Cmd.FuBenParam_value["START_STAGE_USER_CMD"]:
				param = &Cmd.StartStageUserCmd{}

			case Cmd.FuBenParam_value["EXIT_RAID_CMD"]:
				param = &Cmd.ExitMapFubenCmd{}

			case Cmd.FuBenParam_value["TEAMEXP_RAID_REPORT"]:
				param = &Cmd.TeamExpReportFubenCmd{}

			case Cmd.FuBenParam_value["FUBEN_STEP_SYNC"]:
				param = &Cmd.FubenStepSyncCmd{}

			case Cmd.FuBenParam_value["TRACK_FUBEN_USER_CMD"]:
				param = &Cmd.TrackFuBenUserCmd{}

			case Cmd.FuBenParam_value["TEAMEXP_QUERY_INFO"]:
				param = &Cmd.TeamExpQueryInfoFubenCmd{}

			case Cmd.FuBenParam_value["FUBEN_CLEAR_SYNC"]:
				param = &Cmd.FuBenClearInfoCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["MATCHC_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_QUERY_TEAMINFO"]:
				param = &Cmd.QueryTeamPwsTeamInfoMatchCCmd{}

			case Cmd.MatchCParam_value["MATCHCPARAM_TUTOR_MATCHNTF"]:
				param = &Cmd.TutorMatchResultNtfMatchCCmd{}

			case Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_PREPARE_UPDATE"]:
				param = &Cmd.UpdatePreInfoMatchCCmd{}

			case Cmd.MatchCParam_value["MATCHCPARAM_JOIN_ROOM"]:
				param = &Cmd.JoinRoomCCmd{}

			case Cmd.MatchCParam_value["MATCHCPARAM_NTF_MATCHINFO"]:
				param = &Cmd.NtfMatchInfoCCmd{}

			case Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_PREPARE_LIST"]:
				param = &Cmd.TeamPwsPreInfoMatchCCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["ACTIVITY_EVENT_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ActivityEventParam_value["ACTIVITYEVENTPARAM_USER_DATA_NTF"]:
				param = &Cmd.ActivityEventUserDataNtf{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["ACTIVITY_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.ActivityParam_value["ACTIVITYPARAM_GLOBAL_ACT_START"]:
				param = &Cmd.StartGlobalActCmd{}

			case Cmd.ActivityParam_value["ACTIVITYPARAM_ACT_START"]:
				param = &Cmd.StartActCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}

			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)

		} else if cmdId == Cmd.Command_value["SESSION_USER_SOCIALITY_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.SocialityParam_value["SOCIALITYPARAM_FINDUSER"]:
				param = &Cmd.FindUser{}

			case Cmd.SocialityParam_value["SOCIALITYPARAM_UPDATEDATA"]:
				param = &Cmd.SocialDataUpdate{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SCENE_USER_PET_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.PetParam_value["PETPARAM_WORK_STARTWORK"]:
				param = &Cmd.StartWorkPetCmd{}

			case Cmd.PetParam_value["PETPARAM_RESTORE_EGG"]:
				param = &Cmd.EggRestorePetCmd{}

			case Cmd.PetParam_value["PETPARAM_CAT_SKILLOPTION"]:
				param = &Cmd.CatSkillOptionPetCmd{}

			case Cmd.PetParam_value["PETPARAM_WORK_GETREWARD"]:
				param = &Cmd.GetPetWorkRewardPetCmd{}

			case Cmd.PetParam_value["PETPARAM_WORK_SPACEUPDATE"]:
				param = &Cmd.WorkSpaceUpdate{}

			case Cmd.PetParam_value["PETPARAM_WORK_QUERYWORKDATA"]:
				param = &Cmd.QueryPetWorkDataPetCmd{}

			case Cmd.PetParam_value["PETPARAM_PETINFO_UPDATE"]:
				param = &Cmd.PetInfoUpdatePetCmd{}

			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYBATTLEPET"]:
				param = &Cmd.QueryBattlePetCmd{}

			case Cmd.PetParam_value["PETPARAM_ADVENTURE_QUERYLIST"]:
				param = &Cmd.QueryPetAdventureListPetCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SCENE_USER_TIP_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.TipParam_value["TIPPARAM_RED"]:
				param = &Cmd.GameTipCmd{}
			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["INFINITE_TOWER_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.TowerParam_value["ETOWERPARAM_LAYER_SYNC"]:
				param = &Cmd.TowerLayerSyncTowerCmd{}

			case Cmd.TowerParam_value["ETOWERPARAM_USERTOWERINFO"]:
				param = &Cmd.UserTowerInfoCmd{}

			case Cmd.TowerParam_value["ETOWERPARAM_TEAMTOWERINFO"]:
				param = &Cmd.TeamTowerInfoCmd{}

			case Cmd.TowerParam_value["ETOWERPARAM_TEAMTOWERSUMMARY"]:
				param = &Cmd.TeamTowerSummaryCmd{}

			case Cmd.TowerParam_value["ETOWERPARAM_REPLY"]:
				param = &Cmd.TeamTowerReplyCmd{}

			case Cmd.TowerParam_value["ETOWERPARAM_INVITE"]:
				param = &Cmd.TeamTowerInviteCmd{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		} else if cmdId == Cmd.Command_value["SESSION_USER_MAIL_PROTOCMD"] {
			switch cmdParamId {
			case Cmd.MailParam_value["MAILPARAM_UPDATE"]:
				param = &Cmd.MailUpdate{}

			case Cmd.MapParam_value["MAILPARAM_READ"]:
				param = &Cmd.MailRead{}

			case Cmd.MailParam_value["MAILPARAM_GETATTACH"]:
				param = &Cmd.MailAttach{}

			default:
				log.Infof("没有parsing")
				continue
			}
			err = ParseCmd(o, param)
			PrintTranslateMsgResult(cmdParamName, err, param)
		}
	}
}
