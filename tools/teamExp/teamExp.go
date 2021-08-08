package main

import (
	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/utils"
	"flag"
	"fmt"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	autoScriptVer = "0.0.3"
)

var (
	maxExpFubenTime = time.Minute * 40
)

func init() {
	//log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	logPath := path.Join(filepath.Dir(ex), "teamExp.log")
	f, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", logPath, "%Y-%m-%d"),
		rotatelogs.WithRotationCount(10),
		rotatelogs.WithMaxAge(time.Hour*24),
	)
	if err != nil {
		log.Errorf("failed to open teamExp.log: %s", err)
	} else {
		log.Infof("saving log to %s", logPath)
		mw := io.MultiWriter(os.Stdout, f)
		log.SetOutput(mw)
	}
}

func worker(wg *sync.WaitGroup, completeFuben chan bool, cPath string, skills map[uint32]utils.SkillItem, items *utils.ItemsLoader, enableDebug bool) {
	defer wg.Done()
	time.Sleep(5 * time.Second)
	conf := config.NewServerConfigs(cPath)
	gameConnect := gameConnection.NewConnection(conf, skills, items)
	if enableDebug {
		gameConnect.DebugMsg = enableDebug
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("worker for %s ready", cPath)

	// 开始副本
	gameConnect.GameServerLogin()
	quit := make(chan bool)
	gameConnect.CheckForFubenInviteInBackground(quit)
	disable := make(chan *bool)
	gameConnect.EnableAutoAttack([]string{"all"}, disable)
	gameConnect.InviteTeamExpFuben()
	gameConnect.AutoSubmitWantedQuest()
	go func() {
		time.Sleep(10 * time.Second)
		tradeHistory, _ := gameConnect.QueryTradeHistoryLog(0)
		gameConnect.HandleTradeHistory(tradeHistory)
		gameConnect.ShouldChangeScene = true
		gameConnect.ChangeMap(gameConnect.Role.GetMapId())
		if !gameConnect.IsCookWorkSpaceInUse() {
			log.Infof("没有宠物烹饪打工, 开始任命一个")
			gameConnect.BattlePetToWork()
		} else {
			gameConnect.GetCookWorkSpaceReward()
		}
		gameConnect.DailySignIn()
	}()

	startTime := time.Now()
	waitCount := 0
	for time.Since(startTime) < maxExpFubenTime {
		select {
		case <-completeFuben:
			_, _ = gameConnect.GetTeamEXPQueryInfo()
			log.Infof("%s:研究所副本次数 剩余%d/%d",
				gameConnect.Role.GetRoleName(),
				gameConnect.Role.TeamExpFubenInfo.GetRewardtimes(),
				gameConnect.Role.TeamExpFubenInfo.GetTotaltimes(),
			)
			log.Infof("队长完成副本 %s 退出", gameConnect.Role.GetRoleName())
			disableAuto := true
			disable <- &disableAuto
			gameConnect.Close()
			return
		default:
			if gameConnect.Role.GetInGame() {
				time.Sleep(15 * time.Second)
				if gameConnect.Role.TeamExpFubenInfo == nil {
					_, _ = gameConnect.GetTeamEXPQueryInfo()
				}
				if gameConnect.IsTeamLeader(gameConnect.Role.GetRoleId(), false) && gameConnect.Role.TeamExpFubenInfo.GetRewardtimes() == 0 {
					log.Infof("%s 队长生态副本已完成发送退出消息", gameConnect.Role.GetRoleName())
					for i := 0; i < len(gameConnect.GetOnlineMemebers()); i++ {
						go func() {
							completeFuben <- true
						}()
					}
					disableAuto := true
					disable <- &disableAuto
					gameConnect.Close()
					return
				}
			} else {
				if waitCount > 3 {
					return
				}
				waitCount += 1
				log.Infof("等待角色进入游戏")
				time.Sleep(10 * time.Second)
			}
		}
	}
	if time.Since(startTime) > maxExpFubenTime {
		log.Warnf("研究所副本用时超时: %d 分钟", time.Since(startTime)/time.Minute)
		return
	}
}

func main() {
	log.Infof("ROM auto script version %s", autoScriptVer)
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	scriptYml := flag.String("script", "script.yml", "yaml file of the script actions")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	skillJson := flag.String("skillJson", "", "json file of skills")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()

	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	scriptActions := config.ScriptParser(*scriptYml)
	skills := utils.NewSkillParser(*skillJson)
	log.Infof("%v", scriptActions)

	fi, err := os.Stat(*confFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Errorf("failed to find configuration file")
		} else {
			log.Error(err)
		}
		log.Exit(1)
	}
	var wg sync.WaitGroup
	switch mode := fi.Mode(); {
	case mode.IsDir():
		teamFolders, err := ioutil.ReadDir(*confFile)
		if err != nil {
			log.Fatalf("failed to read directory %s", *confFile)
		}
		for _, teamFolder := range teamFolders {
			if teamFolder.IsDir() {
				tfPath := path.Join(*confFile, teamFolder.Name())
				configFiles, _ := ioutil.ReadDir(tfPath)
				completeFuben := make(chan bool)
				var matchedConfig []string
				for _, cFile := range configFiles {
					if !cFile.IsDir() &&
						(strings.HasSuffix(cFile.Name(), "yml") ||
							strings.HasSuffix(cFile.Name(), "yaml")) {
						configPath := path.Join(tfPath, cFile.Name())
						log.Infof("found configuration for team %s conf file is %s", teamFolder.Name(), configPath)
						matchedConfig = append(matchedConfig, configPath)
						wg.Add(1)
					}
				}
				for _, cPath := range matchedConfig {
					go func(cPath string, completeFuben chan bool) {
						worker(&wg, completeFuben, cPath, skills, items, *enableDebug)
						log.Infof("worker completed")
					}(cPath, completeFuben)
				}
				wg.Wait()
			} else {
				log.Warnf("%s is not a folder", teamFolder.Name())
			}
		}

	case mode.IsRegular():

	}
}
