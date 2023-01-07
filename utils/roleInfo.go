package utils

import (
	"time"

	Cmd "ROMProject/Cmds"
)

type MatchDetail struct {
	MatchInfo     *Cmd.NtfMatchInfoCCmd
	TeamPrepInfos *Cmd.TeamPwsPreInfoMatchCCmd
	PrepedMember  []*uint64
	CreatedAt     time.Time
}

type RoleInfo struct {
	AcceptAllTeamInvite bool
	RoleId              *uint64
	RoleName            *string
	MapId               *uint32
	MapName             *string
	Pos                 *Cmd.ScenePos
	AuthConfirm         *bool
	Authed              *bool
	RoleSelected        *bool
	LoginResult         *uint32
	InGame              *bool
	Silver              *uint64
	PackItems           map[Cmd.EPackType]map[string]*Cmd.ItemData
	SkillItems          map[uint32]*Cmd.SkillItem
	Sequence            *uint32
	AutoSkills          map[uint32]*Cmd.SkillItem
	TeamExpFubenInfo    *Cmd.TeamExpQueryInfoFubenCmd
	MatchInfos          map[Cmd.EPvpType]*MatchDetail
	UserAttrs           []*Cmd.UserAttr
	UserDatas           []*Cmd.UserData
	Buffs               map[uint32]*Cmd.BufferData
	TeamData            *Cmd.TeamData
	TeamApply           []*Cmd.TeamApply
	TeamMemberPos       map[uint64]*Cmd.MemberPosUpdate
	CDs                 map[uint32]time.Time
	UserVars            map[Cmd.EVarType]*Cmd.Var
	AccVars             map[Cmd.EAccVarType]*Cmd.AccVar
	QuestList           map[Cmd.EQuestList]*Cmd.QuestList
	UserTowerInfo       *Cmd.UserTowerInfo
	FollowUserId        uint64
	DailySignIn         *Cmd.SignInNtfUserCmd
}

func (r *RoleInfo) GetProfession() string {
	val := GetNpcDataValByType(r.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_PROFESSION)
	return Cmd.EProfession_name[int32(val)]
}

func (r *RoleInfo) GetSkillPoint() int32 {
	val := GetNpcDataValByType(r.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_SKILL_POINT)
	return int32(val)
}

func (r *RoleInfo) GetTotalPoint() int32 {
	val := GetNpcDataValByType(r.UserDatas, Cmd.EUserDataType_EUSERDATATYPE_TOTALPOINT)
	return int32(val)
}

func (r *RoleInfo) GetPackItems() map[Cmd.EPackType]map[string]*Cmd.ItemData {
	return r.PackItems
}

func (r *RoleInfo) GetSilver() uint64 {
	if r.Silver != nil {
		return *r.Silver
	}
	return 0
}

func (r *RoleInfo) GetRoleName() string {
	if r.RoleName != nil {
		return *r.RoleName
	} else {
		return ""
	}
}

func (r *RoleInfo) GetSequence() uint32 {
	return *r.Sequence
}

func (r *RoleInfo) SetRoleName(roleName string) {
	r.RoleName = &roleName
}

func (r *RoleInfo) SetSequence(sequence uint32) {
	r.Sequence = &sequence
}

func (r *RoleInfo) SetMapId(mapId uint32) {
	r.MapId = &mapId
}

func (r *RoleInfo) GetMapId() uint32 {
	if r.MapId != nil {
		return *r.MapId
	}
	return 0
}

func (r *RoleInfo) GetRoleId() uint64 {
	if r.RoleId != nil {
		return *r.RoleId
	}
	return 0
}

func (r *RoleInfo) SetRoleId(roleId uint64) {
	r.RoleId = &roleId
}

func (r *RoleInfo) GetPos() Cmd.ScenePos {
	return *r.Pos
}

func (r *RoleInfo) GetJobLevel() uint64 {
	for _, data := range r.UserDatas {
		if data.GetType() == Cmd.EUserDataType_EUSERDATATYPE_JOBLEVEL {
			return data.GetValue()
		}
	}
	return 0
}

func (r *RoleInfo) GetRoleLevel() uint64 {
	for _, data := range r.UserDatas {
		if data.GetType() == Cmd.EUserDataType_EUSERDATATYPE_ROLELEVEL {
			return data.GetValue()
		}
	}
	return 0
}

func (r *RoleInfo) GetRoleSelected() bool {
	if r.RoleSelected != nil {
		return *r.RoleSelected
	}
	return false
}

func (r *RoleInfo) SetRoleSelected(roleSelected bool) {
	r.RoleSelected = &roleSelected
}

func (r *RoleInfo) GetLoginResult() uint32 {
	if r.LoginResult != nil {
		return *r.LoginResult
	}
	return 1
}

func (r *RoleInfo) SetLoginResult(result uint32) {
	r.LoginResult = &result
}

func (r *RoleInfo) GetRolePos() *Cmd.ScenePos {
	if r.LoginResult != nil {
		return r.Pos
	}
	return nil
}

func (r *RoleInfo) SetRolePos(pos *Cmd.ScenePos) {
	r.Pos = pos
}

func (r *RoleInfo) GetInGame() bool {
	if r.InGame != nil {
		return *r.InGame
	}
	return false
}

func (r *RoleInfo) GetAuthConfirm() bool {
	if r.AuthConfirm != nil {
		return *r.AuthConfirm
	}
	return false
}

func (r *RoleInfo) GetQuestList(questType Cmd.EQuestList) (questList *Cmd.QuestList) {
	if v, ok := r.QuestList[questType]; ok {
		questList = v
	}
	return questList
}

func (r *RoleInfo) GetSkillShortCut() (shortCutItems map[uint32][]*Cmd.SkillShortcut) {
	shortCutItems = map[uint32][]*Cmd.SkillShortcut{}
	for _, item := range r.SkillItems {
		if item.GetShortcuts() == nil {
			continue
		}
		shortCutItems[item.GetId()] = item.GetShortcuts()
	}
	return shortCutItems
}

func (r *RoleInfo) GetSkillAuto() (autoItems map[uint32]uint32) {
	autoItems = map[uint32]uint32{}
	for skillId, shortcuts := range r.GetSkillShortCut() {
		for _, shortcut := range shortcuts {
			if shortcut.GetType() == Cmd.ESkillShortcut_ESKILLSHORTCUT_AUTO {
				autoItems[skillId] = shortcut.GetPos()
			}
		}
	}
	return autoItems
}

func (r *RoleInfo) SetMapName(mapName string) {
	r.MapName = &mapName
}

func (r *RoleInfo) GetMapName() (mapName string) {
	if r.MapName != nil {
		return *r.MapName
	}
	return mapName
}

func NewRole() *RoleInfo {
	role := &RoleInfo{
		PackItems:     map[Cmd.EPackType]map[string]*Cmd.ItemData{},
		SkillItems:    map[uint32]*Cmd.SkillItem{},
		MatchInfos:    map[Cmd.EPvpType]*MatchDetail{},
		AutoSkills:    map[uint32]*Cmd.SkillItem{},
		Buffs:         map[uint32]*Cmd.BufferData{},
		CDs:           map[uint32]time.Time{},
		TeamMemberPos: map[uint64]*Cmd.MemberPosUpdate{},
		UserVars:      map[Cmd.EVarType]*Cmd.Var{},
		AccVars:       map[Cmd.EAccVarType]*Cmd.AccVar{},
		QuestList:     map[Cmd.EQuestList]*Cmd.QuestList{},
	}
	return role
}
