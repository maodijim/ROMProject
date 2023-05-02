package gameConnection

import (
	Cmd "ROMProject/Cmds"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func (g *GameConnection) HandleLoginUserCmd(cmdParamId int32, rawData []byte) (param proto.Message, err error) {
	switch cmdParamId {
	case Cmd.LoginCmdParam_value["LOGIN_RESULT_USER_CMD"]:
		param = &Cmd.LoginResultUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.Role.SetLoginResult(param.(*Cmd.LoginResultUserCmd).GetRet())

	case Cmd.LoginCmdParam_value["REAL_AUTHORIZE_USER_CMD"]:
		param = &Cmd.RealAuthorizeUserCmd{}
		err = proto.Unmarshal(rawData[2:], param)
		g.SetAuthed(param.(*Cmd.RealAuthorizeUserCmd).GetAuthorized())
		// g.SetAuthed(true)

	case Cmd.LoginCmdParam_value["HEART_BEAT_USER_CMD"]:
		param = &Cmd.HeartBeatUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.setLastHeartBeat()

	case Cmd.LoginCmdParam_value["SNAPSHOT_USER_CMD"]:
		param = &Cmd.SnapShotUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.Mutex.Lock()
		for _, char := range param.(*Cmd.SnapShotUserCmd).GetData() {
			if char.GetId() == 0 {
				continue
			}
			roleOption := utils.RoleTeamOption(g.Configs.TeamConfig)
			role := utils.NewRole(roleOption)
			role.SetRoleId(char.GetId())
			role.SetRoleName(char.GetName())
			role.SetSequence(char.GetSequence())
			g.AvailableRoles[char.GetSequence()] = role
			log.Infof("received role %d name %s", char.GetSequence(), role.GetRoleName())
		}
		g.Mutex.Unlock()

	case Cmd.LoginCmdParam_value["SERVERTIME_USER_CMD"]:
		param = &Cmd.ServerTimeUserCmd{}
		err = utils.ParseCmd(rawData, param)
		g.setCurrentMsgIndex(1)
		if g.IsTCPConnected() {
			g.SendServerTimeUserCmd(0)
		}

	case Cmd.LoginCmdParam_value["CONFIRM_AUTHORIZE_USER_CMD"]:
		param = &Cmd.ConfirmAuthorizeUserCmd{}
		err = utils.ParseCmd(rawData, param)
		if g.IsTCPConnected() {
			g.SendServerTimeUserCmd(0)
			g.Role.AuthConfirm = param.(*Cmd.ConfirmAuthorizeUserCmd).Success
			if *g.Role.AuthConfirm == false {
				log.Warn("account is NOT authorized to trade and perform certain actions")
			} else {
				log.Info("account is authorized to trade and perform certain actions")
			}
		}

	case Cmd.LoginCmdParam_value["REQ_LOGIN_PARAM_USER_CMD"]:
		param = &Cmd.ReqLoginParamUserCmd{}
		err = utils.ParseCmd(rawData, param)
		if g.IsTCPConnected() {
			g.Mutex.Lock()
			g.Configs.Sha1Str = *param.(*Cmd.ReqLoginParamUserCmd).Sha1
			g.Mutex.Unlock()
			g.SendServerTimeUserCmd(0)
			g.SendReqUserLoginCmd(*param.(*Cmd.ReqLoginParamUserCmd).Timestamp)
		}
	}
	return param, err
}
