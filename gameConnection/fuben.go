package gameConnection

import (
	"math/rand"
	"strings"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

var (
	FubenProtoCmdId   = Cmd.Command_value["FUBEN_PROTOCMD"]
	MatchProtoCmdId   = Cmd.Command_value["MATCHC_PROTOCMD"]
	TeamExpFubenMapId = uint32(50103)
)

// GetTeamEXPQueryInfo 查询怪物研究所副本次数
func (g *GameConnection) GetTeamEXPQueryInfo() (info *Cmd.TeamExpQueryInfoFubenCmd, err error) {
	cmd := &Cmd.TeamExpQueryInfoFubenCmd{}
	g.AddNotifier("TEAMEXP_QUERY_INFO")
	g.sendProtoCmd(
		cmd,
		FubenProtoCmdId,
		Cmd.FuBenParam_value["TEAMEXP_QUERY_INFO"],
	)
	res, err := g.waitForResponse("TEAMEXP_QUERY_INFO")
	if res != nil {
		g.Mutex.Lock()
		info = res.(*Cmd.TeamExpQueryInfoFubenCmd)
		g.Role.TeamExpFubenInfo = info
		g.Mutex.Unlock()
	}
	return info, err
}

// AcceptFubenTeamInvite 接受组队副本邀请
func (g *GameConnection) AcceptFubenTeamInvite(eType Cmd.EPvpType) {
	cmd := &Cmd.TeamPwsPreInfoMatchCCmd{
		Etype: &eType,
	}
	g.sendProtoCmd(
		cmd,
		MatchProtoCmdId,
		Cmd.MatchCParam_value["MATCHCPARAM_TEAMPWS_PREPARE_UPDATE"],
	)
}

func (g *GameConnection) CheckForFubenInviteInBackground(quit chan bool) {
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				if g.Role.MatchInfos[Cmd.EPvpType_EPVPTYPE_TEAMEXP] != nil {
					log.Infof("%s 找到生态研究所组队邀请", g.Role.GetRoleName())
					time.Sleep(time.Second)
					accepted := utils.Uint64SliceContains(g.Role.MatchInfos[Cmd.EPvpType_EPVPTYPE_TEAMEXP].PrepedMember, g.Role.GetRoleId())
					if !accepted {
						log.Infof("接受生态研究所组队邀请")
						g.AcceptFubenTeamInvite(Cmd.EPvpType_EPVPTYPE_TEAMEXP)
					}
					if len(g.Role.MatchInfos[Cmd.EPvpType_EPVPTYPE_TEAMEXP].TeamPrepInfos.GetTeaminfos()) > 0 {
						teammemberCount := len(g.Role.MatchInfos[Cmd.EPvpType_EPVPTYPE_TEAMEXP].TeamPrepInfos.GetTeaminfos()[0].GetCharids())
						acceptedCount := len(g.Role.MatchInfos[Cmd.EPvpType_EPVPTYPE_TEAMEXP].PrepedMember)
						if teammemberCount <= acceptedCount {
							log.Infof("所有队员接受了邀请 删除组队邀请记录")
							delete(g.Role.MatchInfos, Cmd.EPvpType_EPVPTYPE_TEAMEXP)
							// 分散站位
							go func() {
								time.Sleep(time.Second * 10)
								g.Role.GetRolePos()
								randomX := g.Role.GetRolePos().GetX() + (rand.Int31n(20000) - 10000)
								randomY := g.Role.GetRolePos().GetY() + (rand.Int31n(20000) - 10000)
								z := g.Role.GetRolePos().GetZ()
								newPos := &Cmd.ScenePos{
									X: &randomX,
									Y: &randomY,
									Z: &z,
								}
								g.MoveChart(newPos)
							}()
						}
					}
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()
}

func (g *GameConnection) JoinMatchRoom(roomId uint64, roomType Cmd.EPvpType, isQuick bool, teamExpType Cmd.ERewardTeamExpType) {
	cmd := &Cmd.JoinRoomCCmd{
		Roomid:      &roomId,
		Type:        &roomType,
		Isquick:     &isQuick,
		Teamexptype: &teamExpType,
	}
	g.sendProtoCmd(cmd,
		MatchProtoCmdId,
		Cmd.MatchCParam_value["MATCHCPARAM_JOIN_ROOM"],
	)
}

func (g *GameConnection) InviteTeamExpFuben() {
	go func() {
		waitForFullMemberCount := 0
		g.Role.AcceptAllTeamInvite = true
		for {
			select {
			case <-g.quit:
				return
			default:
				if g.Role.GetInGame() {
					if len(g.Role.TeamData.GetMembers()) == 1 && !strings.Contains(g.Role.GetRoleName(), g.Configs.TeamConfig.GetLeaderName()) {
						g.ExitTeam()
						time.Sleep(2 * time.Second)
						continue
					}
					if g.Role.TeamData == nil {
						log.Warnf("%s没有组队 跳过申请生态研究所副本", g.Role.GetRoleName())
						g.AutoCreateJoinTeam(g.Configs.TeamConfig.GetLeaderName())
						time.Sleep(15 * time.Second)
						continue
					}
					if !strings.HasPrefix(g.GetTeamLeaderName(false), g.Configs.TeamConfig.GetLeaderName()) {
						log.Infof("队长不在队伍里 退出队伍")
						g.ExitTeam()
						time.Sleep(3 * time.Second)
						continue
					}
					if !g.IsTeamLeader(g.Role.GetRoleId(), false) {
						log.Warnf("%s 不是队长 跳过申请生态研究所副本", g.Role.GetRoleName())
						return
					}
					if len(g.Role.TeamData.GetMembers()) < 6 && waitForFullMemberCount <= 5 {
						log.Infof("队伍%s只有%d人等待队员加入 %d/%d次",
							g.Role.TeamData.GetName(),
							len(g.Role.TeamData.GetMembers()),
							waitForFullMemberCount,
							5,
						)
						for _, apply := range g.Role.TeamApply {
							log.Infof("同意 %s 进队申请", apply.GetName())
							g.AcceptTeamApply(apply.GetGuid())
						}
						time.Sleep(15 * time.Second)
						waitForFullMemberCount += 1
						continue
					}
					fubenInfo, err := g.GetTeamEXPQueryInfo()
					if err != nil {
						log.Errorf("failed to get team exp query info: %v", err)
						return
					}
					if fubenInfo.GetRewardtimes() == 0 {
						log.Warnf("%s 生态研究所副本次数已用完 %d/%d", g.Role.GetRoleName(), fubenInfo.GetRewardtimes(), fubenInfo.GetTotaltimes())
						return
					}
					if !g.AllTeamMemberOnline() {
						if !g.IsTeamLeader(g.Role.GetRoleId(), false) {
							return
						}
						log.Infof("%s 生态研究所:等待所有队员上线; 离线队员: %v", g.Role.GetRoleName(), g.GetOfflineMemebers())
						time.Sleep(10 * time.Second)
						continue
					}
					log.Infof("%s 申请进入生态研究所副本", g.Role.GetRoleName())
					g.JoinMatchRoom(50103, Cmd.EPvpType_EPVPTYPE_TEAMEXP, true, Cmd.ERewardTeamExpType_REWARD_TEAM_EXP_ITEM)
					go func() {
						timedOut := 300 * time.Second
						for start := time.Now(); time.Since(start) < timedOut; {
							if g.Role.GetMapId() == TeamExpFubenMapId {
								time.Sleep(5 * time.Second)
								log.Infof("%s starting team exp fuben", g.Role.GetRoleName())
								g.StartTeamExpFuben()
								return
							}
							time.Sleep(15 * time.Second)
						}
					}()
					if g.Role.GetMapId() != TeamExpFubenMapId {
						time.Sleep(10 * time.Second)
						continue
					}
					return
				} else {
					time.Sleep(15 * time.Second)
				}
			}
		}
	}()
}

func (g *GameConnection) StartTeamExpFuben() {
	go func() {
		for {
			select {
			case <-g.quit:
				return
			default:
				if g.Role.GetInGame() {
					if !g.IsTeamLeader(g.Role.GetRoleId(), false) {
						log.Warnf("不是队长 跳过开启副本")
						return
					}
					if g.Role.GetMapId() != TeamExpFubenMapId {
						log.Warnf("不在地图 魔物研究所·岭之间")
						time.Sleep(20 * time.Second)
						continue
					}
					cmd := &Cmd.BeginFireFubenCmd{}
					g.sendProtoCmd(cmd,
						FubenProtoCmdId,
						Cmd.FuBenParam_value["BEGIN_FIRE_FUBENCMD"])
					return
				} else {
					log.Infof("角色还没进入游戏 等待")
				}
				time.Sleep(time.Second * 10)
			}
		}
	}()
}

func (g *GameConnection) ExitTeamExpFuben() {
	if g.Role.GetMapId() != TeamExpFubenMapId {
		return
	}
	cmd := &Cmd.ExitMapFubenCmd{}
	g.sendProtoCmd(cmd,
		FubenProtoCmdId,
		Cmd.FuBenParam_value["EXIT_RAID_CMD"])
}
