package backendTasks

import (
	"context"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/AutoHunting"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type AutoHuntingTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (a *AutoHuntingTask) GetTaskName() string {
	return "自动挂机打怪"
}

func (a *AutoHuntingTask) StartTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		return
	}
	a.running = true
	oneTick := time.After(time.Millisecond * 50)
	go func() {
		autoHuntTask := AutoHunting.NewHuntTask(a.ctx, a.gameConn)
		for {
			select {
			case <-a.ctx.Done():
				autoHuntTask.Stop()
				return
			case <-oneTick:
				autoHuntTask.Start()
			}
		}
	}()
}

func (a *AutoHuntingTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	a.cancel()
	a.running = false
}

func (a *AutoHuntingTask) IsRunning() bool {
	return a.running
}

func (a *AutoHuntingTask) GetLogs() []string {
	return a.gameConn.GetLogs()
}

func (a *AutoHuntingTask) GetLogStream() chan string {
	return a.gameConn.LogNotify
}

func (a *AutoHuntingTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

func (a *AutoHuntingTask) GetConfig() TaskConfig {
	return usersSpace.Configs[a.username].HuntConfig
}

func (a *AutoHuntingTask) UpdateConfig(config TaskConfig) error {
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

// NewAutoHuntTask 创建一个新的自动狩猎任务实例
func NewAutoHuntingTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &AutoHuntingTask{
		username: username,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}
