package main

import (
	"fmt"
	"sync"

	gameConfig "ROMProject/config"
	"ROMProject/gameConnection"
)

var (
	featureTasks       = make(map[string]Task)
	featureBackendLock = sync.Mutex{}
)

type TaskConfig interface{}

type Task interface {
	GetTaskName() string
	StartTask()
	StopTask()
	IsRunning() bool
	GetLogStream() chan string
	GetGameConnection() *gameConnection.GameConnection
	GetConfig() TaskConfig
	UpdateConfig(config TaskConfig) error
}

func registerFeatureTask(username string, task Task) {
	featureBackendLock.Lock()
	defer featureBackendLock.Unlock()
	featureTasks[username] = task
}

func getFeatureTask(username string) Task {
	featureBackendLock.Lock()
	defer featureBackendLock.Unlock()
	return featureTasks[username]
}

func removeFeatureTask(username string) {
	task := getFeatureTask(username)
	if task != nil {
		task.StopTask()
	}
	featureBackendLock.Lock()
	defer featureBackendLock.Unlock()
	delete(featureTasks, username)
}

// 自动附魔 -------------------------------------------------------------------------

// AutoEnchantTask 自动附魔任务
type AutoEnchantTask struct {
	username  string
	running   bool
	logStream chan string
	gameConn  *gameConnection.GameConnection
	stopChan  chan struct{}
	taskMutex sync.RWMutex
}

func (a *AutoEnchantTask) GetConfig() TaskConfig {
	return configs[a.username].EnchantConfig
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

	// 启动自动附魔逻辑的协程
	go func() {
		// 模拟自动附魔过程
		for {
			select {
			case <-a.stopChan:
				a.logStream <- "自动附魔任务已停止"
				return
			default:
				// 这里放置自动附魔的具体实现逻辑
				a.logStream <- "正在执行自动附魔..."
				// 模拟工作负载
				// time.Sleep(2 * time.Second)
			}
		}
	}()
}

func (a *AutoEnchantTask) StopTask() {
	a.taskMutex.Lock()
	defer a.taskMutex.Unlock()
	if a.running {
		close(a.stopChan)
		a.running = false
	}
}

func (a *AutoEnchantTask) IsRunning() bool {
	return a.running
}

func (a *AutoEnchantTask) GetLogStream() chan string {
	// TODO implement me
	panic("implement me")
}

func (a *AutoEnchantTask) GetGameConnection() *gameConnection.GameConnection {
	return a.gameConn
}

// NewAutoEnchantTask 自动附魔
func NewAutoEnchantTask(username string) Task {
	newTask := &AutoEnchantTask{
		username:  username,
		running:   true,
		logStream: make(chan string),
		gameConn:  nil,
		stopChan:  make(chan struct{}),
	}
	newTask.StartTask()
	registerFeatureTask(username, newTask)
	return newTask
}

// End -------------------------------------------------------------------------
