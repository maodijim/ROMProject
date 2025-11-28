package main

import (
	"context"
	"flag"
	"io"
	"os"
	"time"

	"ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/autoEnchant"
	"ROMProject/utils"

	log "github.com/sirupsen/logrus"
)

const (
	ver     = "0.2.1"
	logFile = "autoEnchant.log"
)

var (
	mw io.Writer
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	f, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw = io.MultiWriter(os.Stdout, f)
	} else {
		log.Warnf("无法写入日志文件 %s，使用默认输出", logFile)
		mw = os.Stdout
	}
	log.SetOutput(mw)
}

func main() {
	log.Infof("自动附魔版本 %s", ver)
	configPath := flag.String("config", "config.yml", "配置文件路径")
	enableDebug := flag.Bool("debug", false, "是否开启调试模式")
	speed := flag.Uint("speed", 850, "附魔速度，单位毫秒")
	enchantCount := flag.Uint("count", 1, "每次附魔次数")
	flag.Parse()
	items := utils.NewItemsLoader("", "", "")
	conf := config.NewServerConfigs(*configPath)
	conf.EnchantConfig.EnchantCount = uint32(*enchantCount)
	skills := utils.NewSkillParser("")
	g := gameConnection.NewConnection(conf, skills, items).LoadMonster("")
	if *enableDebug {
		g.DebugMsg = true
		log.SetLevel(log.DebugLevel)
	}

	ctx, _ := context.WithCancel(context.Background())
	task := autoEnchant.NewEnchantTask(ctx, g, *speed)
	task.SetLogger(mw)
	task.Start()
	<-ctx.Done()
	time.Sleep(time.Second * 1)
}
