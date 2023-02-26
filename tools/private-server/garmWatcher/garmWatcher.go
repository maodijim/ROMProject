package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	"github.com/olivere/elastic/v7"
	log "github.com/sirupsen/logrus"
)

var (
	g         *gameConnection.GameConnection
	esClient  *elastic.Client
	dialogNtf chan interface{}
	msgNtf    chan interface{}
	nameMap   = map[string]uint32{
		"卡仑":   gameTypes.MapId_GingerbreadCity.Uint32(),
		"嗜血怪人": gameTypes.MapId_Skellington.Uint32(),
		"死灵骑士": gameTypes.MapId_Niflheim.Uint32(),
		"狼外婆":  gameTypes.MapId_MistyForest.Uint32(),
	}
	names = []string{"卡仑", "嗜血怪人", "死灵骑士", "狼外婆"}
)

type killInfo struct {
	Killer    string
	MvpName   string
	Event     string
	KillTime  time.Time
	ShowTime  time.Time
	TimeStamp time.Time
	Reporter  string
}

func main() {
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	hiddenMvp := flag.String("hiddenMvp", "卡仑", "隐藏boss的名字: 卡仑,嗜血怪人,死灵骑士,狼外婆")
	flag.Parse()
	if !utils.Contains(names, *hiddenMvp) {
		log.Fatalf("隐藏boss的名字不正确: %s, 必须是卡仑,嗜血怪人,死灵骑士,狼外婆 其一", *hiddenMvp)
	}
	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	_ = createEsClient(conf.EsConfig)
	g.GameServerLogin()
	if g.Role.GetMapId() != nameMap[(*hiddenMvp)] {
		log.Infof("不在目标地图 瞬移过去中")
		g.ExitMapWait(nameMap[(*hiddenMvp)])
		time.Sleep(2 * time.Second)
	}
	g.AddNotifier(gameTypes.NtfType_UserActionDialog)
	g.AddNotifier(gameTypes.NtfType_SysMsg)
	dialogNtf = g.Notifier(gameTypes.NtfType_UserActionDialog)
	msgNtf = g.Notifier(gameTypes.NtfType_SysMsg)
	for {
		info := getNotification()
		if info.Killer != "" {
			log.Infof(
				"MVP 杀手: %s, MVP: %s, 杀死时间: %s",
				info.Killer,
				info.MvpName,
				info.KillTime.UTC().Format("2006-01-02 15:04:05"),
			)
		}
		if !info.ShowTime.IsZero() {
			log.Infof(
				"MVP 刷新时间: %s",
				info.ShowTime.UTC().Format("2006-01-02 15:04:05"),
			)
		}
		info.TimeStamp = time.Now().UTC()
		info.Reporter = g.Role.GetRoleName()
		saveMvp(info)
	}
}

func getNotification() killInfo {
	ticker := time.NewTicker(30 * time.Second)
	for {
		select {
		case sNtf := <-msgNtf:
			ntf, ok := sNtf.(*Cmd.SysMsg)
			// 4000 is mvp killed sysmsg
			if ok && ntf.GetId() == 4000 {
				params := ntf.GetParams()
				killer := params[0].GetParam()
				mvpName := params[1].GetParam()
				return killInfo{Killer: killer, MvpName: mvpName, KillTime: time.Now().UTC()}
			}
		case dNtf := <-dialogNtf:
			ntf, ok := dNtf.(*Cmd.UserActionNtf)
			if !ok {
				log.Warnf("未知的对话通知: %v", dNtf)
				continue
			}
			switch ntf.GetValue() {
			case 65100:
				return killInfo{ShowTime: time.Now().UTC(), Event: "狼外婆刷新"}
			case 55435:
				return killInfo{ShowTime: time.Now().UTC(), Event: "卡仑刷新"}
			case 53170:
				return killInfo{ShowTime: time.Now().UTC(), Event: "玩具士兵星刷新"}
			case 65073:
				return killInfo{ShowTime: time.Now().UTC(), Event: "小恶魔星刷新"}
			case 65101:
				return killInfo{ShowTime: time.Now().UTC(), Event: "嗜血怪人刷新"}
			case 65102:
				return killInfo{ShowTime: time.Now().UTC(), Event: "死灵骑士刷新"}
			case 65103:
				return killInfo{ShowTime: time.Now().UTC(), Event: "银月魔女刷新"}
			default:
				return killInfo{ShowTime: time.Now().UTC(), Event: fmt.Sprintf("未知刷新通知ID: %d", ntf.GetValue())}
			}
		case <-ticker.C:
			log.Infof("等待MVP刷新")
		}
	}
}

func saveMvp(killInfo killInfo) {
	log.Infof("保存MVP数据")
	req := elastic.NewBulkIndexRequest().Index(
		getIndexName(),
	).Doc(
		killInfo,
	)
	res, err := esClient.Bulk().Add(req).Do(context.Background())
	if err != nil {
		log.Errorf("Failed to save mvp data: %v", err)
	} else {
		log.Infof("response from elasticsearch: %d failed", len(res.Failed()))
	}
}

func createEsClient(config config.EsConfig) (err error) {
	log.Infof("Establishing connection to elastic search: %v", config.Urls)
	esClient, err = elastic.NewClient(
		elastic.SetURL(config.Urls...),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false),
	)
	if err != nil {
		log.Fatalf("Failed to connect to elastic search: %v", err)
	}
	return err
}

func getIndexName() string {
	return fmt.Sprintf("private-rom-mvp-%s", time.Now().UTC().Format("2006-01-02"))
}
