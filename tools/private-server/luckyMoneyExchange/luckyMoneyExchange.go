package main

import (
	"flag"
	"math"

	Cmd "ROMProject/Cmds"
	"ROMProject/config"
	"ROMProject/gameConnection"
	gameTypes "ROMProject/gameConnection/types"
	"ROMProject/utils"
	log "github.com/sirupsen/logrus"
)

var (
	g *gameConnection.GameConnection
)

func main() {
	confFile := flag.String("configPath", "config.yml", "Game Server Configuration Yaml Path")
	itemFile := flag.String("itemPath", "", "Exchange Item Json Path")
	buffFile := flag.String("buffPath", "", "Buff Json Path")
	monsterFile := flag.String("monsterPath", "", "Monster Json Path")
	skillFile := flag.String("skillPath", "", "Skill Json Path")
	enableDebug := flag.Bool("debug", false, "Enable Debugging")
	flag.Parse()
	items := utils.NewItemsLoader(*itemFile, *buffFile, "")
	conf := config.NewServerConfigs(*confFile)
	skills := utils.NewSkillParser(*skillFile)
	g = gameConnection.NewConnection(conf, skills, items).LoadMonster(*monsterFile)
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}
	g.GameServerLogin()
	// 找有没有人礼盒
	// 有的话就打开
	if item := g.FindPackItemByName("节日礼包2.0", Cmd.EPackType_EPACKTYPE_MAIN); item != nil {
		log.Infof("找到礼盒2.0, 打开")
		g.UseItem(item.GetBase().GetGuid(), 1)
	} else {
		log.Infof("没有找到礼盒2.0")
	}

	mails := g.GetMails()
	for _, mail := range mails {
		if len(mail.GetAttach().GetAttachs()) > 0 {
			attachs := mail.GetAttach().GetAttachs()
			log.Infof("邮件有附件: %s", attachs)
			log.Infof("收取邮件 标题：%s 发送人：%s 内容：%s", mail.GetTitle(), mail.GetSender(), mail.GetMsg())
			g.GetMailAttachment(mail.GetId())
		}
	}

	getLuckyMoneyCount()
	shopItems, _ := g.QueryShopConfig(gameTypes.ShopType_Item, 10)
	// targetItem := "熔岩宝石"
	// targetItem := "黯魂粉尘"
	// targetItem := "圣诞红色转蛋"
	targetItem := "焰之余烬"
	pocketName := "B格猫皇家红包1.0"
	count := uint32(20)
	if shopItems != nil {
		for _, item := range shopItems.Goods {
			name := g.FindItemNameById(item.GetItemid())
			if name == pocketName {
				if item := g.FindPackItemByName("B格猫皇家红包2.0", Cmd.EPackType_EPACKTYPE_MAIN); item != nil {
					log.Infof("找到%s, 购买%d个", pocketName, 1)
				}
				g.BuyShopItem(item, 1)
			}
			if name == targetItem {
				max := maxPurchase(item.GetMoneycount())
				if max == 0 {
					log.Infof("红包不够, 不能购买")
					continue
				} else if max < count {
					log.Infof("红包不够, 最多购买%d个%s", max, targetItem)
					count = max
				}
				log.Infof("找到%s, 购买%d", targetItem, count)
				g.BuyShopItem(item, count)
			}
		}
	}
	getLuckyMoneyCount()
}

func getLuckyMoneyCount() (count uint32) {
	item := g.FindPackItemByName("B格猫皇家红包1.0", Cmd.EPackType_EPACKTYPE_MAIN)
	if item != nil {
		log.Infof("红包有%d个", item.GetBase().GetCount())
		count = item.GetBase().GetCount()
	} else {
		log.Infof("没有找到红包")
	}
	return count
}

func maxPurchase(cost uint32) (count uint32) {
	lCount := getLuckyMoneyCount()
	if lCount > 0 {
		count = uint32(math.Floor(float64(lCount) / float64(cost)))
	}
	return count
}
