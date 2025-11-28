package main

import (
	"context"
	"flag"

	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/positionHelper"
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
	ctx, _ := context.WithCancel(context.Background())
	positionHelper.NewPositionTask(ctx, g).Start()
	<-ctx.Done()
}
