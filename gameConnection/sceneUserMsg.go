package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
	notifier "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

var (
	sceneUserCmdId = Cmd.Command_value["SCENE_USER_PROTOCMD"]
)

func (g *GameConnection) HandleSceneUserProtoCmd(cmdParamId int32, rawData []byte) (err error) {
	var param proto.Message
	switch cmdParamId {
	case Cmd.CmdParam_value["GOTO_USER_CMD"]:
		param = &Cmd.GoToUserCmd{}
		err = utils.ParseCmd(rawData, param)
		goTo := param.(*Cmd.GoToUserCmd)
		if goTo.GetPos() != nil {
			g.Role.Pos = goTo.GetPos()
			log.Debugf("%s moved to %v", g.Role.GetRoleName(), g.Role.Pos)
		}

	case Cmd.CmdParam_value["DELETE_ENTRY_USER_CMD"]:
		param = &Cmd.DeleteEntryUserCmd{}
		err = utils.ParseCmd(rawData, param)
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
		err = utils.ParseCmd(rawData, param)
		cmd := param.(*Cmd.RetMoveUserCmd)
		if g.Role.FollowUserId != 0 && g.Role.FollowUserId == cmd.GetCharid() {
			go func() {
				log.Debugf("following user %d to %v", cmd.GetCharid(), cmd.GetPos())
				if cmd.GetPos() != nil {
					g.MoveChart(*cmd.GetPos())
				}
			}()
		}
		if cmd.GetCharid() == g.Role.GetRoleId() {
			if time.Now().Second()%5 == 0 {
				log.Debugf(
					"Moving charater %s to position: %v",
					g.Role.GetRoleName(),
					param.(*Cmd.RetMoveUserCmd).GetPos(),
				)
			}
			g.Role.SetRolePos(cmd.GetPos())
		} else if g.MapNpcs[cmd.GetCharid()] != nil {
			g.Mutex.Lock()
			g.MapNpcs[cmd.GetCharid()].Pos = cmd.GetPos()
			g.Mutex.Unlock()
		} else if g.MapUsers[cmd.GetCharid()] != nil {
			g.Mutex.Lock()
			g.MapUsers[cmd.GetCharid()].Pos = cmd.GetPos()
			g.Mutex.Unlock()
		}

	case Cmd.CmdParam_value["USERPARAM_USERSYNC"]:
		param = &Cmd.UserSyncCmd{}
		err = utils.ParseCmd(rawData, param)
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
		if g.Notifier(notifier.NtfType_AddAttributePoint) != nil {
			g.Notifier(notifier.NtfType_AddAttributePoint) <- true
		}

	case Cmd.CmdParam_value["CHANGE_SCENE_USER_CMD"]:
		param = &Cmd.ChangeSceneUserCmd{}
		err = utils.ParseCmd(rawData, param)
		changeScene := param.(*Cmd.ChangeSceneUserCmd)
		if changeScene.GetMapID() != 0 {
			log.Infof("Moving %s to map %v", g.Role.GetRoleName(), changeScene)
			g.MapNpcs = map[uint64]*Cmd.MapNpc{}
			g.MapUsers = map[uint64]*Cmd.MapUser{}
			g.Role.SetMapId(changeScene.GetMapID())
			g.Role.SetMapName(changeScene.GetMapName())
			g.inMap = true
		}
		if changeScene != nil {
			g.Role.SetRolePos(changeScene.GetPos())
			inGame := true
			g.Role.InGame = &inGame
			if g.ShouldChangeScene {
				go func() {
					g.ChangeMap(g.Role.GetMapId())
				}()

			}
		}
	}
	return err
}
