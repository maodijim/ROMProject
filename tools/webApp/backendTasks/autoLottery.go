package backendTasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/autoLottery"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type AutoLotteryTask struct {
	// Define the fields for the AutoLotteryTask struct
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (a *AutoLotteryTask) GetTaskName() string {
	return "自动抽奖"
}

func (a *AutoLotteryTask) StartTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		return
	}
	a.running = true
	oneTick := time.After(time.Millisecond * 50)
	// 启动自动抽奖逻辑的协程
	go func() {
		lotteryTask := autoLottery.NewLotteryTask(a.ctx, a.gameConn)
		// 模拟自动抽奖过程
		for {
			select {
			case <-a.ctx.Done():
				lotteryTask.Stop()
				return
			case <-oneTick:
				// 这里放置自动抽奖的具体实现逻辑
				lotteryTask.Start()
			}
		}
	}()
}

func (a *AutoLotteryTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	a.cancel()
	a.gameConn.Close()
	a.running = false
}

func (a *AutoLotteryTask) IsRunning() bool {
	return a.running
}

func (a *AutoLotteryTask) GetLogs() []string {
	return a.gameConn.GetLogs()
}

func (a *AutoLotteryTask) GetLogStream() chan string {
	return a.gameConn.LogNotify
}

func (a *AutoLotteryTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

func (a *AutoLotteryTask) GetConfig() TaskConfig {
	return usersSpace.Configs[a.username].LotteryConfig
}

func (a *AutoLotteryTask) UpdateConfig(config TaskConfig) error {
	g := a.GetGameConnection()
	if g == nil {
		return nil
	}
	if config == nil {
		return nil
	}
	// 类型断言
	cfg, ok := config.(gameConfig.LotteryConfig)
	if !ok {
		return fmt.Errorf("invalid config type")
	}
	g.Configs.LotteryConfig = cfg
	return nil
}

// NewLotteryTask creates a new AutoLotteryTask for the given username
func NewLotteryTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &AutoLotteryTask{
		username: username,
		running:  false,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}
