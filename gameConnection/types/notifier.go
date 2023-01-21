package gameConnection

type NotifierType string

const (
	NtfType_BossWorldNtf               NotifierType = "BOSS_WORLD_NTF"
	NtfType_AddAttributePoint          NotifierType = "AddAttrPoint"
	NtfType_TeamExpQueryInfo           NotifierType = "TEAMEXP_QUERY_INFO"
	NtfType_TeamParamQueryUserTeamInfo NotifierType = "TEAMPARAM_QUERYUSERTEAMINFO"
	NtfType_TeamParamQueryTeamInfo     NotifierType = "TEAMPARAM_QUERYTEAMINFO"
	NtfType_QuestList                  NotifierType = "QUESTPARAM_QUESTLIST"
	NtfType_PetAdventureQueryList      NotifierType = "PETPARAM_ADVENTURE_QUERYLIST"
	NtfType_PetQueryWorkData           NotifierType = "PETPARAM_WORK_QUERYWORKDATA"
	NtfType_PetQueryBattlePet          NotifierType = "PETPARAM_ADVENTURE_QUERYBATTLEPET"
	NtfType_SocialityFindUser          NotifierType = "SOCIALITYPARAM_FINDUSER"
	NtfType_TowerUserInfo              NotifierType = "ETOWERPARAM_USERTOWERINFO"
	NtfType_TowerTeamSummary           NotifierType = "ETOWERPARAM_TEAMTOWERSUMMARY"
	NtfType_InterviewQuestion          NotifierType = "INTER_QUESTION"
	NtfType_ShopQueryShopConfig        NotifierType = "SHOPPARAM_QUERY_SHOP_CONFIG"
	NtfType_ShopBuyItem                NotifierType = "SHOPPARAM_BUYITEM"
	NtfType_UserActionDialog           NotifierType = "EUSERACTIONTYPE_DIALOG"
	NtfType_User2QueryZoneStatus       NotifierType = "USER2PARAM_QUERY_ZONESTATUS"
	NtfType_UserItemPickup             NotifierType = "ITEM_PICKUP"
)
