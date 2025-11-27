package main

import (
	"encoding/json"
	"net/http"
)

type Feature struct {
	Name         string   `json:"name"`
	Desc         string   `json:"desc"`
	Actions      []string `json:"actions,omitempty"`
	FunctionName string   `json:"functionName,omitempty"`
}

type RunningTaskInfoResponse struct {
	Username    string `json:"username"`
	FeatureName string `json:"featureName"`
	Status      string `json:"status"`
}

var features = []Feature{
	{Name: "自动附魔", Desc: "Description of Feature A", Actions: []string{"Start", "Stop", "Configure", "Log"}, FunctionName: "AutoEnchant"},
	{Name: "自动MVP", Desc: "Description of Feature B", Actions: []string{"Start", "Stop", "Configure", "Log"}, FunctionName: "AutoMVP"},
	{Name: "FeatureC", Desc: "Description of Feature C", Actions: []string{"Start", "Stop", "Configure", "Log"}, FunctionName: "FeatureC"},
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

	featureBackendLock.Lock()
	runningTasks := make([]RunningTaskInfoResponse, 0)
	for username, task := range featureTasks {
		if task.IsRunning() {
			runningTasks = append(runningTasks, RunningTaskInfoResponse{
				Username:    username,
				FeatureName: task.GetTaskName(),
				Status:      "运行中",
			})
		}
	}
	featureBackendLock.Unlock()

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

	featureBackendLock.Lock()
	task := featureTasks[username]
	featureBackendLock.Unlock()

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
	featureBackendLock.Lock()
	existingTask := featureTasks[req.Username]
	if existingTask != nil && existingTask.IsRunning() {
		featureBackendLock.Unlock()
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "User already has a running task",
		})
		return
	}
	featureBackendLock.Unlock()

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
	var task Task
	switch selectedFeature.FunctionName {
	case "AutoEnchant":
		task = NewAutoEnchantTask(req.Username)
		json.NewEncoder(w).Encode(Response{
			Success: task.IsRunning(),
			Message: "Feature started successfully",
		})
		return
	case "AutoMVP":
		// task = NewAutoMVPTask(req.Username)
		// TODO: Implement task creation
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Feature not implemented yet",
		})
		return
	default:
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Unknown feature",
		})
		return
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

	featureBackendLock.Lock()
	task := featureTasks[req.Username]
	featureBackendLock.Unlock()

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

	removeFeatureTask(req.Username)

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

	featureBackendLock.Lock()
	task := featureTasks[username]
	featureBackendLock.Unlock()

	var config TaskConfig
	if task != nil && task.GetTaskName() == functionName {
		config = task.GetConfig()
	} else {
		// Return default config based on feature type
		switch functionName {
		case "AutoEnchant":
			config = configs[username].EnchantConfig
		case "AutoMVP":
			config = configs[username].HuntConfig
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

	featureBackendLock.Lock()
	task := featureTasks[req.Username]
	featureBackendLock.Unlock()

	if configs[req.Username] == nil {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "No configuration found for user",
		})
		return
	}

	// Update the config in files
	userConfigs := configs[req.Username]
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
