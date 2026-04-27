package backendTasks

import (
	"context"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/AutoBossHunting"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type AutoBossHuntingTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (a *AutoBossHuntingTask) GetTaskName() string {
	return "自动MVP"
}

func (a *AutoBossHuntingTask) StartTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		return
	}
	a.running = true
	oneTick := time.After(time.Millisecond * 50)
	go func() {
		autoMvpTask := AutoBossHunting.NewBossHuntTask(a.ctx, a.gameConn)
		for {
			select {
			case <-a.ctx.Done():
				autoMvpTask.Stop()
				return
			case <-oneTick:
				autoMvpTask.Start()
			}
		}
	}()
}

func (a *AutoBossHuntingTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	a.cancel()
	a.running = false
}

func (a *AutoBossHuntingTask) IsRunning() bool {
	return a.running
}

func (a *AutoBossHuntingTask) GetLogs() []string {
	return a.gameConn.GetLogs()
}

func (a *AutoBossHuntingTask) GetLogStream() chan string {
	return a.gameConn.LogNotify
}

func (a *AutoBossHuntingTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

func (a *AutoBossHuntingTask) GetConfig() TaskConfig {
	return usersSpace.Configs[a.username].HuntConfig
}

func (a *AutoBossHuntingTask) UpdateConfig(config TaskConfig) error {
	g := a.GetGameConnection()
	if g == nil {
		return nil
	}
	if config == nil {
		return nil
	}
	cfg, ok := config.(gameConfig.HuntConfig)
	if !ok {
		return nil
	}
	g.Configs.HuntConfig = cfg
	return nil
}

// NewAutoBossHuntingTask 创建一个新的自动MVP任务实例
func NewAutoBossHuntingTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &AutoBossHuntingTask{
		username: username,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}
