package main

import (
	"flag"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	teamTowerVer = "0.0.3"
	maxTowerTime = time.Minute * 75
)

var (
	teamLeaderName = ""
)

func init() {
	// log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
}

func worker(wg *sync.WaitGroup, cPath string, skills map[uint32]utils.SkillItem, items *utils.ItemsLoader, enableDebug bool) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	conf := config.NewServerConfigs(cPath)
	if teamLeaderName != "" {
		conf.SetTeamLeader(teamLeaderName)
		conf.SetFollowTeamLeader(true)
	}
	gameConnect := gameConnection.NewConnection(conf, skills, items)
	gameConnect.ShouldChangeScene = true
	if enableDebug {
		gameConnect.DebugMsg = enableDebug
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("worker for %s ready", cPath)
	gameConnect.GameServerLogin()
	_ = gameConnect.GetMainPackItems()
	_ = gameConnect.GetTempMainPackItems()
	gameConnect.DailySignIn()
	startTime := time.Now()
	for time.Since(startTime) < maxTowerTime {
		if gameConnect.Role.GetInGame() {
			roleName := gameConnect.Role.GetRoleName()
			if gameConnect.Role.TeamData != nil {
				if !strings.HasPrefix(gameConnect.GetTeamLeaderName(false), gameConnect.Configs.TeamConfig.GetLeaderName()) {
					log.Infof("队长不在队伍里 退出队伍")
					gameConnect.ExitTeam()
					time.Sleep(5 * time.Second)
					continue
				} else {
					if gameConnect.Role.UserTowerInfo != nil {
						curLayer := gameConnect.Role.UserTowerInfo.GetCurmaxlayer()
						if curLayer >= 100 {
							log.Infof("完成了100层塔")
							gameConnect.QuickSellItems()
							time.Sleep(time.Second * 2)
							log.Infof("回收 %d 件在临时背包道具", len(gameConnect.Role.PackItems[Cmd.EPackType_EPACKTYPE_TEMP_MAIN]))
							gameConnect.GetTempItems()
							time.Sleep(30 * time.Second)
							return
						}
					}
				}
				if strings.HasPrefix(gameConnect.Role.GetMapName(), "无限塔") || strings.HasPrefix(gameConnect.Role.GetMapName(), "恩德勒斯塔") {
					log.Infof("%s 在塔里 %s", roleName, gameConnect.Role.GetMapName())
					if gameConnect.Configs.TeamConfig.FollowTeamLeader && gameConnect.Role.FollowUserId != gameConnect.GetTeamLeader(false) {
						gameConnect.Role.FollowUserId = gameConnect.GetTeamLeader(false)
					}
					// leaderMapId := utils.GetMemberDataByType(gameConnect.GetTeamLeaderData(false).GetDatas(), Cmd.EMemberData_EMEMBERDATA_MAPID)
					// if uint32(leaderMapId) != gameConnect.Role.GetMapId() {
					// 	log.Infof("队长离开了地图 %s 离开副本", roleName)
					// 	gameConnect.ExitMapWait(31)
					// }
					time.Sleep(15 * time.Second)
					continue
				}
			} else if gameConnect.Role.TeamData == nil {
				log.Infof("%s 申请加入%s队伍", roleName, gameConnect.Configs.TeamConfig.GetLeaderName())
				gameConnect.AutoCreateJoinTeam(gameConnect.Configs.TeamConfig.GetLeaderName())
				time.Sleep(10 * time.Second)
			}
		} else {
			for {
				select {
				case <-time.After(15 * time.Second):
					log.Warnf("进入角色超时")
				}
			}
		}
	}
}

func main() {
	log.Infof("ROM team tower version %s", teamTowerVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	exchangeItemFile := flag.String("exchangeItemPath", "", "Exchange Item Json Path")
	itemFile := flag.String("itemPath", "", "All Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	skillJson := flag.String("skillJson", "", "json file of skills")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	teamLeader := flag.String("teamLeader", "", "Team Leader Name")
	flag.Parse()

	teamLeaderName = *teamLeader
	items := utils.NewItemsLoader(*exchangeItemFile, *buffFile, *itemFile)
	skills := utils.NewSkillParser(*skillJson)

	fi, _ := os.Stat(*confFile)
	var wg sync.WaitGroup
	switch mode := fi.Mode(); {
	case mode.IsDir():
		teamFolders, err := os.ReadDir(*confFile)
		if err != nil {
			log.Fatalf("failed to read directory %s", *confFile)
		}
		var matchedConfig []string
		for _, cFile := range teamFolders {
			configPath := path.Join(*confFile, cFile.Name())
			log.Infof("found configuration for team %s conf file is %s", cFile.Name(), configPath)
			matchedConfig = append(matchedConfig, configPath)
			wg.Add(1)
		}
		for _, cPath := range matchedConfig {
			go func(cPath string) {
				worker(&wg, cPath, skills, items, *enableDebug)
				log.Infof("worker completed")
			}(cPath)
		}
		wg.Wait()
	case mode.IsRegular():
		log.Errorf("Expect %s to be a folder contains all config files", *confFile)
	}
}
