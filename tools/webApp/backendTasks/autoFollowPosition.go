package backendTasks

import (
	"context"
	"sync"
	"time"

	"ROMProject/gameConnection"
	"ROMProject/tools/positionHelper"
	"ROMProject/tools/webApp/usersSpace"
	"ROMProject/utils"
)

type AutoFollowPositionTask struct {
	username  string
	running   bool
	gameConn  *gameConnection.GameConnection
	ctx       context.Context
	cancel    context.CancelFunc
	taskMutex sync.RWMutex
}

func (a *AutoFollowPositionTask) GetTaskName() string {
	return "自动跟随定位"
}

func (a *AutoFollowPositionTask) StartTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		return
	}
	a.running = true

	go func() {
		followTask := positionHelper.NewPositionTask(a.ctx, a.gameConn)
		oneTick := time.After(time.Millisecond * 50)
		for {
			select {
			case <-a.ctx.Done():
				followTask.Stop()
				return
			case <-oneTick:
				// 这里放置自动附魔的具体实现逻辑
				followTask.Start()
			}
		}
	}()
}

func (a *AutoFollowPositionTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	a.cancel()
	a.running = false
}

func (a *AutoFollowPositionTask) IsRunning() bool {
	return a.running
}

func (a *AutoFollowPositionTask) GetLogs() []string {
	return a.gameConn.GetLogs()
}

func (a *AutoFollowPositionTask) GetLogStream() chan string {
	return a.gameConn.LogNotify
}

func (a *AutoFollowPositionTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

func (a *AutoFollowPositionTask) GetConfig() TaskConfig {
	return map[string]any{}
}

func (a *AutoFollowPositionTask) UpdateConfig(config TaskConfig) error {
	return nil
}

func NewAutoFollowPositionTask(username string) Task {
	g := gameConnection.NewConnection(usersSpace.Configs[username], utils.NewSkillParser(""), utils.NewItemsLoader("", "", "")).LoadMonster("")
	ctx, cancel := context.WithCancel(context.Background())
	t := &AutoFollowPositionTask{
		username: username,
		gameConn: g,
		ctx:      ctx,
		cancel:   cancel,
	}
	RegisterFeatureTask(username, t)
	return t
}
