package backendTasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/autoEnchant"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

// 自动附魔 -------------------------------------------------------------------------

// AutoEnchantTask 自动附魔任务
type AutoEnchantTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (a *AutoEnchantTask) GetLogs() []string {
	return a.gameConn.GetLogs()
}

func (a *AutoEnchantTask) GetConfig() TaskConfig {
	return usersSpace.Configs[a.username].EnchantConfig
}

func (a *AutoEnchantTask) UpdateConfig(config TaskConfig) error {
	g := a.GetGameConnection()
	if g == nil {
		return nil
	}
	if config == nil {
		return nil
	}
	// 类型断言
	cfg, ok := config.(gameConfig.EnchantConfig)
	if !ok {
		return fmt.Errorf("invalid config type")
	}
	g.Configs.EnchantConfig = cfg
	return nil
}

func (a *AutoEnchantTask) GetTaskName() string {
	return "自动附魔"
}

func (a *AutoEnchantTask) StartTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		return
	}
	a.running = true
	oneTick := time.After(time.Millisecond * 50)
	// 启动自动附魔逻辑的协程
	go func() {
		enchantTask := autoEnchant.NewEnchantTask(a.ctx, a.gameConn, 850)
		// 模拟自动附魔过程
		for {
			select {
			case <-a.ctx.Done():
				enchantTask.Stop()
				return
			case <-oneTick:
				// 这里放置自动附魔的具体实现逻辑
				enchantTask.Start()
			}
		}
	}()
}

func (a *AutoEnchantTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	a.cancel()
	a.running = false
}

func (a *AutoEnchantTask) IsRunning() bool {
	return a.running
}

func (a *AutoEnchantTask) GetLogStream() chan string {
	return a.gameConn.LogNotify
}

func (a *AutoEnchantTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

// NewAutoEnchantTask 自动附魔
func NewAutoEnchantTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &AutoEnchantTask{
		username: username,
		running:  false,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}

// End -------------------------------------------------------------------------
