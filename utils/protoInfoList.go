package utils

import (
	Cmd "ROMProject/Cmds"
)

var (
	NameParamMap = map[string]map[int32]string{
		"SESSION_USER_MAIL_PROTOCMD":      Cmd.MailParam_name,
		"SCENE_USER_TIP_PROTOCMD":         Cmd.TipParam_name,
		"MATCHC_PROTOCMD":                 Cmd.MatchCParam_name,
		"FUBEN_PROTOCMD":                  Cmd.FuBenParam_name,
		"USERSHOW_PROTOCMD":               Cmd.EUserShowParam_name,
		"INTERACT_PROTOCMD":               Cmd.InteractParam_name,
		"SCENE_BOSS_PROTOCMD":             Cmd.BossParam_name,
		"TECHTREE_PROTOCMD":               Cmd.TechTreeParam_name,
		"SCENE_USER_ASTROLABE_PROTOCMD":   Cmd.AstrolabeParam_name,
		"SCENE_USER_PET_PROTOCMD":         Cmd.PetParam_name,
		"ACTHITPOLLY_PROTOCMD":            Cmd.ActHiyPollyParam_name,
		"SCENE_USER_TUTOR_PROTOCMD":       Cmd.TutorParam_name,
		"SCENE_USER_BEING_PROTOCMD":       Cmd.BeingParam_name,
		"SESSION_USER_AUTHORIZE_PROTOCMD": Cmd.AuthorizeParam_name,
		"SESSION_USER_TEAM_PROTOCMD":      Cmd.TeamParam_name,
		"LOGIN_USER_PROTOCMD":             Cmd.LoginCmdParam_name,
		"SCENE_USER_ITEM_PROTOCMD":        Cmd.ItemParam_name,
		"SCENE_USER_MANUAL_PROTOCMD":      Cmd.ManualParam_name,
		"SCENE_USER_CHAT_PROTOCMD":        Cmd.ChatParam_name,
		"USER_EVENT_PROTOCMD":             Cmd.EventParam_name,
		"SCENE_USER_PROTOCMD":             Cmd.CmdParam_name,
		"SCENE_USER2_PROTOCMD":            Cmd.User2Param_name,
		"RECORD_USER_TRADE_PROTOCMD":      Cmd.RecordUserTradeParam_name,
		"SESSION_USER_SOCIALITY_PROTOCMD": Cmd.SocialityParam_name,
		"SCENE_USER_CARRIER_PROTOCMD":     Cmd.CarrierParam_name,
		"SCENE_USER_SKILL_PROTOCMD":       Cmd.SkillParam_name,
		"INFINITE_TOWER_PROTOCMD":         Cmd.TowerParam_name,
		"CHAT_PROTOCMD":                   Cmd.ChatParam_name,
		"ROGUELIKE_PROTOCMD":              Cmd.RoguelikeParam_name,
		"ACTIVITY_EVENT_PROTOCMD":         Cmd.ActivityEventParam_name,
		"SCENE_USER_ACHIEVE_PROTOCMD":     Cmd.AchieveParam_name,
		"SCENE_USER_PHOTO_PROTOCMD":       Cmd.PhotoParam_name,
		"ERROR_USER_PROTOCMD":             Cmd.ErrCmdParam_name,
		"SCENE_USER_QUEST_PROTOCMD":       Cmd.QuestParam_name,
		"SCENE_USER_MAP_PROTOCMD":         Cmd.MapParam_name,
		"SESSION_USER_SHOP_PROTOCMD":      Cmd.ShopParam_name,
		"ACTIVITY_PROTOCMD":               Cmd.ActivityParam_name,
	}
)
