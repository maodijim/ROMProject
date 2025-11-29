package main

import (
	"encoding/json"
	"net/http"

	gameConfig "ROMProject/config"
	"ROMProject/tools/webApp/backendTasks"
	"ROMProject/tools/webApp/usersSpace"
)

type Feature struct {
	Name         string                                  `json:"name"`
	Desc         string                                  `json:"desc"`
	Actions      []string                                `json:"actions,omitempty"`
	FunctionName string                                  `json:"functionName,omitempty"`
	execFunc     func(username string) backendTasks.Task `json:"-"`
}

type RunningTaskInfoResponse struct {
	Username    string `json:"username"`
	FeatureName string `json:"featureName"`
	Status      string `json:"status"`
}

var features = []Feature{
	{Name: "自动附魔", Desc: "Description of Feature A", Actions: []string{"Start", "Stop", "Configure", "Log"}, FunctionName: "AutoEnchant", execFunc: backendTasks.NewAutoEnchantTask},
	{Name: "自动MVP", Desc: "Description of Feature B", Actions: []string{"Start", "Stop", "Configure", "Log"}, FunctionName: "AutoMVP", execFunc: backendTasks.NewAutoBossHuntingTask},
	{Name: "自动跟随定位", Desc: "Description of Feature C", Actions: []string{"Start", "Stop", "Log"}, FunctionName: "AutoFollowPosition", execFunc: backendTasks.NewAutoFollowPositionTask},
}

func handleGetFeatures(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}
	response := Response{Success: true, Data: features}
	json.NewEncoder(w).Encode(response)
}

func handleGetRunningTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	runningTasks := make([]RunningTaskInfoResponse, 0)
	for username, task := range backendTasks.FeatureTasks {
		if task.IsRunning() {
			runningTasks = append(runningTasks, RunningTaskInfoResponse{
				Username:    username,
				FeatureName: task.GetTaskName(),
				Status:      "运行中",
			})
		}
	}
	backendTasks.FeatureBackendLock.Unlock()

	response := Response{Success: true, Data: runningTasks}
	json.NewEncoder(w).Encode(response)
}

func handleGetUserRunningTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	username := r.URL.Query().Get("username")
	if username == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username parameter is required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil || !task.IsRunning() {
		json.NewEncoder(w).Encode(Response{
			Success: true,
			Data:    nil,
			Message: "No running task found",
		})
		return
	}

	runningTask := RunningTaskInfoResponse{
		Username:    username,
		FeatureName: task.GetTaskName(),
		Status:      "运行中",
	}

	response := Response{Success: true, Data: runningTask}
	json.NewEncoder(w).Encode(response)
}

func handleStartFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Username    string `json:"username"`
		FeatureName string `json:"featureName"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.FeatureName == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username and featureName are required",
		})
		return
	}

	// Check if task is already running
	backendTasks.FeatureBackendLock.Lock()
	existingTask := backendTasks.FeatureTasks[req.Username]
	if existingTask != nil && existingTask.IsRunning() {
		backendTasks.FeatureBackendLock.Unlock()
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "账户已有任务在运行中，请先停止当前任务",
		})
		return
	}
	backendTasks.FeatureBackendLock.Unlock()

	// Find the feature
	var selectedFeature *Feature
	for _, feature := range features {
		if feature.Name == req.FeatureName {
			selectedFeature = &feature
			break
		}
	}

	if selectedFeature == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Feature not found",
		})
		return
	}

	// Create and start task based on feature
	var task backendTasks.Task
	switch selectedFeature.execFunc {
	case nil:
		// TODO: Implement task creation
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Feature not implemented yet",
		})
		return
	default:
		task = selectedFeature.execFunc(req.Username)
		task.StartTask()
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Feature started successfully",
	})
}

func handleStopFeature(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username is required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[req.Username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "No task found for user",
		})
		return
	}

	if !task.IsRunning() {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Task is not running",
		})
		return
	}

	backendTasks.RemoveFeatureTask(req.Username)

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Feature stopped successfully",
	})
}

func handleGetFeatureConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	username := r.URL.Query().Get("username")
	functionName := r.URL.Query().Get("functionName")

	if username == "" || functionName == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username and functionName are required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[username]
	backendTasks.FeatureBackendLock.Unlock()

	var config backendTasks.TaskConfig
	if task != nil && task.GetTaskName() == functionName {
		config = task.GetConfig()
	} else {
		// Return default config based on feature type
		switch functionName {
		case "AutoEnchant":
			config = usersSpace.Configs[username].EnchantConfig
		case "AutoMVP":
			defaultConfig := usersSpace.Configs[username].HuntConfig.GetDefault()
			var existingConfig gameConfig.HuntConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].HuntConfig
				if existingConfig.PrepEliteCD == 0 {
					existingConfig.PrepEliteCD = defaultConfig.PrepEliteCD
				}
				if len(existingConfig.MVP) == 0 {
					existingConfig.MVP = defaultConfig.MVP
				}
				if len(existingConfig.Mini) == 0 {
					existingConfig.Mini = defaultConfig.Mini
				}
			} else {
				existingConfig = defaultConfig
			}
			config = existingConfig
		default:
			config = map[string]interface{}{}
		}
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    config,
	})
}

func handleUpdateFeatureConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Username    string                 `json:"username"`
		FeatureName string                 `json:"featureName"`
		Config      map[string]interface{} `json:"config"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.FeatureName == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username and featureName are required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[req.Username]
	backendTasks.FeatureBackendLock.Unlock()

	if usersSpace.Configs[req.Username] == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "No configuration found for user",
		})
		return
	}

	// Update the config in files
	userConfigs := usersSpace.Configs[req.Username]
	switch req.FeatureName {
	case "AutoEnchant":
		config := userConfigs.EnchantConfig
		config.ParseFromInterface(req.Config)
		userConfigs.EnchantConfig = config
	case "AutoMVP":
		config := userConfigs.HuntConfig
		config.ParseFromInterface(req.Config)
		userConfigs.HuntConfig = config
	default:
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Unknown feature",
		})
		return
	}

	err := saveConfigs()
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	if task != nil && task.GetTaskName() == req.FeatureName {
		if err := task.UpdateConfig(req.Config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to update task config: " + err.Error(),
			})
			return
		}
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Configuration updated successfully",
	})
}

// Add this handler to featureHandlers.go

func handleGetFeatureTaskLog(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	username := r.URL.Query().Get("username")
	featureName := r.URL.Query().Get("featureName")
	if username == "" || featureName == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username and featureName are required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil || task.GetTaskName() != featureName {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Task not found",
		})
		return
	}

	// Assume Task has GetLogs() []string
	type LogResponse struct {
		Logs []string `json:"logs"`
	}
	logs := task.GetLogs()

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    LogResponse{Logs: logs},
	})
}

func handleGetFeatureTaskChatHistory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}
	username := r.URL.Query().Get("username")
	featureName := r.URL.Query().Get("featureName")
	if username == "" || featureName == "" {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username and featureName are required",
		})
		return
	}

	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil || task.GetTaskName() != featureName {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Task not found",
		})
		return
	}

	// Assume Task has GetGameConnection() *gameConnection.GameConnection
	conn := task.GetGameConnection()
	if conn == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Game connection not found",
		})
		return
	}

	// Get chat history from the game connection
	chatHistory := conn.GetChatHistory()
	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    chatHistory,
	})
}

func handleSendChatMsg(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var req struct {
		Username    string  `json:"username"`
		FeatureName string  `json:"featureName"`
		Message     string  `json:"message"`
		DestId      float64 `json:"destId"`
		ChannelId   float64 `json:"channelId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}
	if req.Username == "" || req.FeatureName == "" || req.Message == "" || req.ChannelId == 0 {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Username, featureName, message and channelId are required",
		})
		return
	}
	backendTasks.FeatureBackendLock.Lock()
	task := backendTasks.FeatureTasks[req.Username]
	backendTasks.FeatureBackendLock.Unlock()

	if task == nil || task.GetTaskName() != req.FeatureName {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Task not found",
		})
		return
	}

	conn := task.GetGameConnection()
	if conn == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Game connection not found",
		})
		return
	}

	err := conn.SentChatMessage(int32(req.ChannelId), req.Message, uint64(req.DestId))
	if err != nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to send chat message",
		})
		return
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: "Chat message sent successfully",
	})
}
