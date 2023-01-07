package gameConnection

import (
	"time"

	Cmd "ROMProject/Cmds"
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
		g.Mutex.Lock()
		if g.Role.FollowUserId != 0 && g.Role.FollowUserId == cmd.GetCharid() {
			go func() {
				log.Debugf("following user %d to %v", cmd.GetCharid(), cmd.GetPos())
				g.MoveChart(cmd.GetPos())
			}()
		}
		if cmd.GetCharid() == *g.Role.RoleId {
			log.Infof(
				"Moving charater %s to position: %v",
				g.Role.GetRoleName(),
				param.(*Cmd.RetMoveUserCmd).GetPos(),
			)
			g.Role.Pos = cmd.GetPos()
		} else if g.MapNpcs[cmd.GetCharid()] != nil {
			g.MapNpcs[cmd.GetCharid()].Pos = cmd.GetPos()
		} else if g.MapUsers[cmd.GetCharid()] != nil {
			g.MapUsers[cmd.GetCharid()].Pos = cmd.GetPos()
		}
		g.Mutex.Unlock()

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
		if g.notifier["AddAttrPoint"] != nil {
			g.notifier["AddAttrPoint"] <- true
		}

	case Cmd.CmdParam_value["CHANGE_SCENE_USER_CMD"]:
		param = &Cmd.ChangeSceneUserCmd{}
		err = utils.ParseCmd(rawData, param)
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
	return err
}
