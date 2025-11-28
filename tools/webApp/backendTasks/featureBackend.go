package backendTasks

import (
	"sync"

	"ROMProject/gameConnection"
)

var (
	FeatureTasks       = make(map[string]Task)
	FeatureBackendLock = sync.Mutex{}
)

type TaskConfig any

type Task interface {
	GetTaskName() string
	StartTask()
	StopTask()
	IsRunning() bool
	GetLogs() []string
	GetLogStream() chan string
	GetGameConnection() *gameConnection.GameConnection
	GetConfig() TaskConfig
	UpdateConfig(config TaskConfig) error
}

func RegisterFeatureTask(username string, task Task) {
	FeatureBackendLock.Lock()
	defer FeatureBackendLock.Unlock()
	FeatureTasks[username] = task
}

func GetFeatureTask(username string) Task {
	FeatureBackendLock.Lock()
	defer FeatureBackendLock.Unlock()
	return FeatureTasks[username]
}

func RemoveFeatureTask(username string) {
	task := GetFeatureTask(username)
	if task != nil {
		task.StopTask()
	}
	FeatureBackendLock.Lock()
	defer FeatureBackendLock.Unlock()
	delete(FeatureTasks, username)
}
