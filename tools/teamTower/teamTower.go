package main

import (
	"flag"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameConnection2 "ROMProject/gameConnection/types"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

const (
	teamTowerVer = "0.0.4"
	maxTowerTime = time.Minute * 75
)

var (
	teamLeaderName = ""
	targetZone     = uint32(0)
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
	conf := config.NewServerConfigs(cPath)
	if teamLeaderName != "" {
		conf.SetTeamLeader(teamLeaderName)
		conf.SetFollowTeamLeader(true)
	}
	g := gameConnection.NewConnection(conf, skills, items)
	g.ShouldChangeScene = true
	if enableDebug {
		g.DebugMsg = enableDebug
		log.SetLevel(log.DebugLevel)
	}
	log.Infof("worker for %s ready", cPath)
	g.GameServerLogin()
	g.DailySignIn()
	mails := g.GetMails()
	log.Infof("邮件数量: %d", len(mails))
	// 获取邮件
	for _, mail := range mails {
		if len(mail.GetAttach().GetAttachs()) > 0 {
			attachs := mail.GetAttach().GetAttachs()
			log.Infof("邮件有附件: %s", attachs)
			log.Infof("收取邮件 标题：%s 发送人：%s 内容：%s", mail.GetTitle(), mail.GetSender(), mail.GetMsg())
			g.GetMailAttachment(mail.GetId())
		}
	}

	if targetZone != 0 && !IsInTower(g) {
		log.Infof("换线中...")
		if g.Role.GetMapId() != gameConnection2.MapId_IzludeIsland.Uint32() {
			g.ExitMapWait(gameConnection2.MapId_IzludeIsland.Uint32())
			g.Reconnect()
			g.ChangeMap(gameConnection2.MapId_IzludeIsland.Uint32())
		}
		g.MoveChartWait(g.ParsePos(4182, 7086, 10633))
		err := g.MoveToNpcWait("世界线传送师")
		if err != nil {
			log.Errorf("换线失败: %s", err)
			return
		}
		_, err = g.VisitNpcByName("世界线传送师")
		if err != nil {
			log.Errorf("换线失败: %s", err)
		}
		g.QueryZoneStatus()
		g.JumpZone(targetZone, 0)
	}
	startTime := time.Now()
	for time.Since(startTime) < maxTowerTime {
		roleName := g.Role.GetRoleName()
		if g.Role.UserTowerInfo != nil {
			curLayer := g.Role.UserTowerInfo.GetCurmaxlayer()
			if curLayer >= 100 {
				log.Infof("%s 完成了100层塔", g.Role.GetRoleName())
				g.QuickSellItems()
				time.Sleep(time.Second * 2)
				log.Infof("回收 %d 件在临时背包道具", len(g.Role.PackItems[Cmd.EPackType_EPACKTYPE_TEMP_MAIN]))
				g.GetTempItems()
				time.Sleep(time.Second * 2)
				buyTreeBranchBag(g)
				time.Sleep(15 * time.Second)
				return
			}
		}
		if g.Role.TeamData != nil {
			var leaderNameInTeam bool
			if g.Configs.TeamConfig.GetLeaderName() != "" {
				leaderNameInTeam = strings.HasPrefix(g.GetTeamLeaderName(false), g.Configs.TeamConfig.GetLeaderName())
			}
			leaderIdInTeam := false
			if g.GetTeamLeaderData(false) != nil {
				leaderIdInTeam = g.GetTeamLeaderData(false).GetGuid() == *g.Configs.TeamConfig.GetLeaderId()
			}
			if !leaderNameInTeam && !leaderIdInTeam && !IsInTower(g) {
				log.Infof("队长不在队伍里 退出队伍")
				g.ExitTeam()
				time.Sleep(5 * time.Second)
				continue
			}
			if IsInTower(g) {
				log.Infof("%s 在塔里 %s", roleName, g.Role.GetMapName())
				if g.Configs.TeamConfig.FollowTeamLeader {
					if getCurrentLayer(g.Role.GetMapName())%10 == 0 {
						log.Infof("当前层是10的倍数，停止跟随队长%s", g.Role.GetRoleName())
						g.Role.FollowUserId = 0
					} else if getCurrentLayer(g.Role.GetMapName())%10 == 1 {
						log.Infof("当前层是10的倍数+1，移动到传送点%s", g.Role.GetRoleName())
						g.Role.FollowUserId = 0
						// 等待客户端东西加载完毕
						time.Sleep(4 * time.Second)
						// 无限塔传送点
						g.MoveChart(g.ParsePos(-58187, 7981, 12800))
						time.Sleep(5 * time.Second)
						g.Role.FollowUserId = g.GetTeamLeader(false)
					} else if g.Role.FollowUserId != g.GetTeamLeader(false) {
						g.Role.FollowUserId = g.GetTeamLeader(false)
					}
				}
				// leaderMapId := utils.GetMemberDataByType(gameConnect.GetTeamLeaderData(false).GetDatas(), Cmd.EMemberData_EMEMBERDATA_MAPID)
				// if uint32(leaderMapId) != gameConnect.Role.GetMapId() {
				// 	log.Infof("队长离开了地图 %s 离开副本", roleName)
				// 	gameConnect.ExitMapWait(31)
				// }
				time.Sleep(5 * time.Second)
				continue
			}
		} else if g.Role.TeamData == nil {
			log.Infof("%s 申请加入%s队伍", roleName, g.Configs.TeamConfig.GetLeaderName())
			g.AutoCreateJoinTeam(g.Configs.TeamConfig)
			time.Sleep(10 * time.Second)
		}
	}
}

func IsInTower(gameConnect *gameConnection.GameConnection) bool {
	return strings.HasPrefix(gameConnect.Role.GetMapName(), "无限塔") || strings.HasPrefix(gameConnect.Role.GetMapName(), "恩德勒斯塔")
}

func getCurrentLayer(mapName string) uint32 {
	re := regexp.MustCompile("[0-9]+")
	layer := re.FindAllString(mapName, -1)
	if len(layer) == 0 {
		return 0
	}
	parseUint, err := strconv.ParseUint(layer[0], 10, 32)
	if err != nil {
		return 0
	}
	return uint32(parseUint)
}

func buyTreeBranchBag(g *gameConnection.GameConnection) {
	if g.Role.GetLottery() < 40000 {
		log.Infof("猫币不足买树枝礼包")
		return
	}
	shopItems, err := g.QueryShopConfig(gameConnection2.ShopType_Lottery, 1)
	for {
		if err != nil {
			log.Infof("查询商店失败重试")
			time.Sleep(2 * time.Second)
			shopItems, err = g.QueryShopConfig(gameConnection2.ShopType_Lottery, 1)
		} else {
			break
		}
	}
	for _, shopItem := range shopItems.GetGoods() {
		if g.FindItemNameById(shopItem.GetItemid()) == "神秘树枝福袋" {
			log.Infof("买树枝福袋")
			time.Sleep(1 * time.Second)
			g.BuyShopItem(shopItem, shopItem.GetMaxcount())
		} else if g.FindItemNameById(shopItem.GetItemid()) == "每日礼包" {
			log.Infof("买每日礼包")
			time.Sleep(1 * time.Second)
			g.BuyShopItem(shopItem, 1)
			item := g.FindPackItemByName("每日礼包", Cmd.EPackType_EPACKTYPE_MAIN)
			if item != nil {
				log.Infof("使用每日礼包")
				g.UseItem(item.GetBase().GetGuid(), item.GetBase().GetCount())
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
	zoneId := flag.Uint("zoneId", 0, "Zone Id")
	flag.Parse()
	targetZone = uint32(*zoneId)
	teamLeaderName = *teamLeader
	items := utils.NewItemsLoader(*exchangeItemFile, *buffFile, *itemFile)
	skills := utils.NewSkillParser(*skillJson)

	fi, err := os.Stat(*confFile)
	if err != nil {
		log.Fatalf("failed to read configuration file %s: %s", *confFile, err)
	}
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
