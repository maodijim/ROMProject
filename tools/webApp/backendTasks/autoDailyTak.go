package backendTasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
	"ROMProject/tools/private-server/autoDailyTask"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type DailyTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (d *DailyTask) GetTaskName() string {
	return "每日日常任务"
}

func (d *DailyTask) StartTask() {
	d.taskMutex.Lock()
	defer d.taskMutex.Unlock()
	if d.running {
		return
	}
	oneTick := time.After(time.Millisecond * 50)
	d.running = true
	go func() {
		task := autoDailyTask.NewDailyTask(d.ctx, d.gameConn)
		for {
			select {
			case <-d.ctx.Done():
				task.Stop()
				return
			case <-oneTick:
				task.Start()
			}
		}
	}()
}

func (d *DailyTask) StopTask() {
	d.taskMutex.Lock()
	defer d.taskMutex.Unlock()
	d.cancel()
	d.gameConn.Close()
	d.running = false
}

func (d *DailyTask) IsRunning() bool {
	return d.running
}

func (d *DailyTask) GetLogs() []string {
	return d.gameConn.GetLogs()
}

func (d *DailyTask) GetLogStream() chan string {
	return d.gameConn.LogNotify
}

func (d *DailyTask) GetGameConnection() *gameConnection.GameConnection {
	return d.gameConn
}

func (d *DailyTask) GetConfig() TaskConfig {
	return usersSpace.Configs[d.username].DailyTaskConfig
}

func (d *DailyTask) UpdateConfig(config TaskConfig) error {
	g := d.GetGameConnection()
	if g == nil {
		return nil
	}
	if config == nil {
		return nil
	}
	// 类型断言
	cfg, ok := config.(gameConfig.DailyTaskConfig)
	if !ok {
		return fmt.Errorf("invalid config type")
	}
	g.Configs.DailyTaskConfig = cfg
	return nil
}

func NewDailyTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	newTask := &DailyTask{
		username: username,
		running:  false,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, newTask)
	return newTask
}
