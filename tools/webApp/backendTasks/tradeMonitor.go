package backendTasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/tradeMonitor"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type TradeMonitorTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (t *TradeMonitorTask) GetTaskName() string {
	return "交易所监控"
}

func (t *TradeMonitorTask) StartTask() {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()
	if t.running {
		return
	}
	t.running = true
	go func() {
		task := tradeMonitor.NewTradeMonitorTask(t.ctx, t.gameConn)
		oneTick := time.After(time.Millisecond * 50)
		for {
			select {
			case <-t.ctx.Done():
				task.Stop()
				return
			case <-oneTick:
				task.Start()
			}
		}
	}()
}

func (t *TradeMonitorTask) StopTask() {
	t.taskMutex.Lock()
	defer t.taskMutex.Unlock()
	t.cancel()
	t.running = false
}

func (t *TradeMonitorTask) IsRunning() bool {
	return t.running
}

func (t *TradeMonitorTask) GetLogs() []string {
	return t.gameConn.GetLogs()
}

func (t *TradeMonitorTask) GetLogStream() chan string {
	return t.gameConn.LogNotify
}

func (t *TradeMonitorTask) GetGameConnection() *gameConnection.GameConnection {
	return t.gameConn
}

func (t *TradeMonitorTask) GetConfig() TaskConfig {
	return usersSpace.Configs[t.username].TradeMonitorConfig
}

func (t *TradeMonitorTask) UpdateConfig(config TaskConfig) error {
	g := t.GetGameConnection()
	if g == nil {
		return nil
	}
	if config == nil {
		return nil
	}
	// 类型断言
	cfg, ok := config.(gameConfig.TradeMonitorConfig)
	if !ok {
		return fmt.Errorf("invalid config type")
	}
	g.Configs.TradeMonitorConfig = cfg
	return nil
}

func NewTradeMonitorTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &TradeMonitorTask{
		username: username,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}
