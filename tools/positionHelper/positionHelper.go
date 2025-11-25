package main

import (
	"flag"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

var ver = "v0.1.0"

func main() {
	log.Infof("自动附魔版本 %s", ver)
	configPath := flag.String("config", "config.yml", "配置文件路径")
	flag.Parse()
	items := utils.NewItemsLoader("", "", "")
	conf := config.NewServerConfigs(*configPath)
	skills := utils.NewSkillParser("")
	g := gameConnection.NewConnection(conf, skills, items).LoadMonster("")

	g.ShouldChangeScene = true
	g.GameServerLogin()

	for {
		if g.GetTeamLeaderName(true) != "" {
			leader := g.GetTeamLeaderData(true)
			if utils.GetMemberDataByType(leader.GetDatas(), Cmd.EMemberData_EMEMBERDATA_MAPID) != uint64(g.Role.GetMapId()) {
				log.Infof("队长 %s (ID: %d) 不在当前地图", leader.GetName(), leader.GetGuid())
				g.Role.FollowUserId = 0
			} else if g.Role.FollowUserId != 0 {
				pos := g.Role.GetPos()
				log.Infof("正在跟随队长 %s (ID: %d) 坐标: x:%d y:%d z:%d", leader.GetName(), leader.GetGuid(), pos.GetX(), pos.GetY(), pos.GetZ())
			} else {
				log.Infof("开始跟随队长 %s (ID: %d)", leader.GetName(), leader.GetGuid())
				g.FollowUser(leader.GetGuid())
			}
		} else if g.GetTeamLeaderName(true) == "" {
			log.Infof("丢失队长，等待中...")
			g.Role.FollowUserId = 0
		} else {
			log.Infof("当前没有队长，等待中...")
		}
		time.Sleep(time.Second * 5)
	}
}
