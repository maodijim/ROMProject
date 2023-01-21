package gameConnection

import (
	"fmt"
	"strings"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	log "github.com/sirupsen/logrus"
)

var (
	TeamProtoCmdId   = Cmd.Command_value["SESSION_USER_TEAM_PROTOCMD"]
	DefaultTeamType  = uint32(10010)
	DefaultMinTeamLv = uint32(1)
	DefaultMaxTeamLv = uint32(170)
)

func (g *GameConnection) getCurTeamMemberDatas(userId uint64) (datas []*Cmd.MemberData) {
	if g.Role.TeamData == nil {
		return datas
	}
	for _, member := range g.Role.TeamData.GetMembers() {
		if member.GetGuid() == userId {
			datas = member.GetDatas()
			break
		}
	}
	return datas
}

func (g *GameConnection) GetMemberIdByName(userName string) (userId uint64) {
	if g.Role.TeamData != nil {
		for _, member := range g.Role.TeamData.GetMembers() {
			if member.GetName() == userName {
				userId = member.GetGuid()
				break
			}
		}
	}
	return userId
}

func (g *GameConnection) GetMemberNameById(userId uint64) (mName string) {
	if g.Role.TeamData != nil {
		for _, member := range g.Role.TeamData.GetMembers() {
			if member.GetGuid() == userId {
				mName = member.GetName()
				break
			}
		}
	}
	return mName
}

func (g *GameConnection) GetCurrentTeamName() (teamName string) {
	if g.Role.TeamData != nil {
		teamName = g.Role.TeamData.GetName()
	}
	return teamName
}

func (g *GameConnection) removeTeamMember(userId uint64) {
	deleteIndex := -1
	if g.Role.TeamData != nil {
		g.Role.Mutex.Lock()
		for i, member := range g.Role.TeamData.GetMembers() {
			if member.GetGuid() == userId {
				deleteIndex = i
			}
		}
		if deleteIndex > -1 {
			g.Role.TeamData.Members = append(g.Role.TeamData.Members[:deleteIndex], g.Role.TeamData.Members[deleteIndex+1:]...)
		}
		g.Role.Mutex.Unlock()
	}
}

func (g *GameConnection) updateTeamMember(member *Cmd.TeamMember) {
	if g.Role.TeamData != nil {
		for _, m := range g.Role.TeamData.GetMembers() {
			if m.GetGuid() == member.GetGuid() {
				m = member
			}
		}
	}
}

func (g *GameConnection) updateTeamMemberDatas(newDatas *Cmd.MemberDataUpdate) {
	if g.Role.TeamData != nil {
		g.Role.Mutex.Lock()
		for _, m := range g.Role.TeamData.GetMembers() {
			if m.GetGuid() == newDatas.GetId() {
				for _, dNew := range newDatas.GetMembers() {
					isUpdate := false
					for _, dOld := range m.GetDatas() {
						if dNew.GetType() == dOld.GetType() {
							dOld.Value = dNew.Value
							isUpdate = true
							break
						}
					}
					if !isUpdate {
						m.Datas = append(m.Datas, dNew)
					}
				}
			}
		}
		g.Role.Mutex.Unlock()
	} else {
		g.Role.TeamData = &Cmd.TeamData{}
		g.updateTeamMemberDatas(newDatas)
	}
}

func (g *GameConnection) addTeamMember(member *Cmd.TeamMember) {
	if g.Role.TeamData != nil {
		g.Role.Mutex.Lock()
		g.Role.TeamData.Members = append(g.Role.TeamData.Members, member)
		g.Role.Mutex.Unlock()
	}
}

func (g *GameConnection) GetCurrentTeamMembers() (teamMembers []*Cmd.TeamMember) {
	if g.Role.TeamData != nil {
		teamMembers = g.Role.TeamData.GetMembers()
	}
	return teamMembers
}

func (*GameConnection) IsMemberOnline(member *Cmd.TeamMember) bool {
	for _, d := range member.GetDatas() {
		if d.GetType() == Cmd.EMemberData_EMEMBERDATA_OFFLINE && d.GetValue() == 0 {
			return true
		}
	}
	return false
}

func (g *GameConnection) AcceptTeamInvite(userGuid *uint64) {
	t := Cmd.ETeamInviteType_ETEAMINVITETYPE_AGREE
	cmd := &Cmd.ProcessTeamInvite{
		Userguid: userGuid,
		Type:     &t,
	}
	g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TeamParam_value["TEAMPARAM_PROCESSINVITE"],
	)
}

func (g *GameConnection) AcceptTeamFollow(userGuid *uint64) {
	t := Cmd.ETeamInviteType_ETEAMINVITETYPE_AGREE
	cmd := &Cmd.ProcessTeamInvite{
		Userguid: userGuid,
		Type:     &t,
	}
	g.sendProtoCmd(cmd,
		Cmd.Command_value["SESSION_USER_TEAM_PROTOCMD"],
		Cmd.TeamParam_value["TEAMPARAM_PROCESSINVITE"])
}

func (g *GameConnection) GetMemberDataByType(userId uint64, memberDataType Cmd.EMemberData) uint64 {
	for _, d := range g.getCurTeamMemberDatas(userId) {
		if memberDataType == d.GetType() {
			return d.GetValue()
		}
	}
	return 0
}

func (g *GameConnection) GetTeamMemberPos(userId uint64) *Cmd.ScenePos {
	for _, member := range g.Role.TeamData.GetMembers() {
		if member.GetGuid() == userId &&
			uint32(g.GetMemberDataByType(member.GetGuid(), Cmd.EMemberData_EMEMBERDATA_MAPID)) == g.Role.GetMapId() {
			return g.Role.TeamMemberPos[member.GetGuid()].GetPos()
		}
	}
	log.Warnf("Either member not found or member not in the current map: %d", g.Role.GetMapId())
	return nil
}

func (g *GameConnection) GetTeamLeaderData(getTempLeader bool) (member *Cmd.TeamMember) {
	for _, member = range g.Role.TeamData.GetMembers() {
		if uint32(g.GetMemberDataByType(member.GetGuid(), Cmd.EMemberData_EMEMBERDATA_JOB)) == 1 && g.IsMemberOnline(member) {
			return member
		} else if getTempLeader && uint32(g.GetMemberDataByType(member.GetGuid(), Cmd.EMemberData_EMEMBERDATA_JOB)) == 4 && g.IsMemberOnline(member) {
			return member
		}
	}
	return member
}

func (g *GameConnection) GetTeamLeader(getTempLeader bool) uint64 {
	leaderData := g.GetTeamLeaderData(getTempLeader)
	if leaderData != nil {
		return leaderData.GetGuid()
	}
	return g.Role.GetRoleId()
}

func (g *GameConnection) GetTeamLeaderName(getTempLeader bool) string {
	leaderData := g.GetTeamLeaderData(getTempLeader)
	if leaderData != nil {
		return leaderData.GetName()
	}
	return ""
}

func (g *GameConnection) IsTeamLeader(userId uint64, getTempLeader bool) bool {
	if userId == g.GetTeamLeader(getTempLeader) {
		return true
	}
	return false
}

func (g *GameConnection) AllTeamMemberOnline() bool {
	for _, member := range g.Role.TeamData.GetMembers() {
		if g.GetMemberDataByType(member.GetGuid(), Cmd.EMemberData_EMEMBERDATA_OFFLINE) == 1 {
			return false
		}
	}
	return true
}

func (g *GameConnection) GetOnlineMemebers() []string {
	var om []string
	g.Mutex.RLock()
	for _, member := range g.Role.TeamData.GetMembers() {
		if g.IsMemberOnline(member) {
			om = append(om, member.GetName())
		}
	}
	g.Mutex.RUnlock()
	return om
}

func (g *GameConnection) GetOfflineMemebers() []string {
	var om []string
	g.Mutex.RLock()
	for _, member := range g.Role.TeamData.GetMembers() {
		if !g.IsMemberOnline(member) {
			om = append(om, member.GetName())
		}
	}
	g.Mutex.RUnlock()
	return om
}

func (g *GameConnection) ExitTeam() {
	log.Infof("退出队伍")
	cmd := &Cmd.ExitTeam{}
	g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TeamParam_value["TEAMPARAM_EXITTEAM"])
	g.Role.TeamData = nil
}

func (g *GameConnection) AcceptTeamApply(userId uint64) {
	aType := Cmd.ETeamApplyType_ETEAMAPPLYTYPE_AGREE
	cmd := &Cmd.ProcessTeamApply{
		Type:     &aType,
		Userguid: &userId,
	}
	g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TeamParam_value["TEAMPARAM_PROCESSAPPLY"],
	)
}

func (g *GameConnection) CreateTeam(teamType uint32) {
	if g.Role.TeamData == nil {
		if teamType == 0 {
			teamType = DefaultTeamType
		}
		desc := "自由队伍"
		teamName := fmt.Sprintf("%s_的队伍", g.Role.GetRoleName())
		accept := Cmd.EAutoType_EAUTOTYPE_GUILDFRIEND
		cmd := &Cmd.CreateTeam{
			Minlv:      &DefaultMinTeamLv,
			Maxlv:      &DefaultMaxTeamLv,
			Type:       &teamType,
			Autoaccept: &accept,
			Name:       &teamName,
			Desc:       &desc,
		}
		g.sendProtoCmd(cmd,
			TeamProtoCmdId,
			Cmd.TeamParam_value["TEAMPARAM_CREATETEAM"],
		)
	}
}

func (g *GameConnection) QueryTeamInfo(charId uint64) (teamInfo *Cmd.QueryUserTeamInfoTeamCmd) {
	cmd := &Cmd.QueryUserTeamInfoTeamCmd{
		Charid: &charId,
	}
	g.AddNotifier("TEAMPARAM_QUERYUSERTEAMINFO")
	g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TeamParam_value["TEAMPARAM_QUERYUSERTEAMINFO"],
	)
	res, err := g.waitForResponse("TEAMPARAM_QUERYUSERTEAMINFO")
	if err != nil {
		log.Errorf("failed to query team info: %s", err)
	}
	if res != nil {
		teamInfo = res.(*Cmd.QueryUserTeamInfoTeamCmd)
	}
	return teamInfo
}

func (g *GameConnection) TeamMemberApply(guid uint64) {
	cmd := &Cmd.TeamMemberApply{
		Guid: &guid,
	}
	g.sendProtoCmd(cmd,
		TeamProtoCmdId,
		Cmd.TeamParam_value["TEAMPARAM_MEMBERAPPLY"],
	)
}

func (g *GameConnection) AutoCreateJoinTeam(teamConfig config.TeamConfig) {
	if g.Role.TeamData != nil || (teamConfig.GetLeaderName() == "" && *teamConfig.GetLeaderId() == 0) {
		return
	}
	var userSocData *Cmd.SocialData
	if strings.Contains(g.Role.GetRoleName(), teamConfig.GetLeaderName()) {
		log.Infof("创建新队伍")
		g.CreateTeam(DefaultTeamType)
	} else {
		time.Sleep(3 * time.Second)
		if teamConfig.GetLeaderName() != "" {
			res, _ := g.FindUser(teamConfig.GetLeaderName())
			if len(res.GetDatas()) > 0 {
				userSocData = res.GetDatas()[0]
				log.Infof("尝试加入%s队伍", teamConfig.GetLeaderName())
			} else {
				log.Warnf("user %s not found", teamConfig.GetLeaderName())
				if *teamConfig.GetLeaderId() != 0 {
					log.Infof("尝试加入用户ID%d队伍", teamConfig.GetLeaderId())
					userSocData = &Cmd.SocialData{
						Guid: teamConfig.GetLeaderId(),
					}
				}
			}
		}
		teamInfo := g.QueryTeamInfo(userSocData.GetGuid())
		time.Sleep(2 * time.Second)
		g.TeamMemberApply(teamInfo.GetTeamid())

	}
}
