package main

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

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
	{
		Name:         "自动附魔",
		Desc:         "Description of Feature A",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "AutoEnchant",
		execFunc:     backendTasks.NewAutoEnchantTask,
	},
	{
		Name:         "自动MVP",
		Desc:         "Description of Feature B",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "AutoMVP",
		execFunc:     backendTasks.NewAutoBossHuntingTask,
	},
	{
		Name:         "自动挂机打怪",
		Desc:         "Description of Feature C",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "AutoHunt",
		execFunc:     backendTasks.NewAutoHuntingTask,
	},
	{
		Name:         "自动跟随定位",
		Desc:         "Description of Feature D",
		Actions:      []string{"Start", "Stop", "Log"},
		FunctionName: "AutoFollowPosition",
		execFunc:     backendTasks.NewAutoFollowPositionTask,
	},
	{
		Name:         "交易所监控",
		Desc:         "Description of Feature E",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "MarketMonitor",
		execFunc:     backendTasks.NewTradeMonitorTask,
	},
	{
		Name:         "自动抽奖",
		Desc:         "Description of Feature F",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "AutoLottery",
		execFunc:     backendTasks.NewLotteryTask,
	},
	{
		Name:         "每日日常任务",
		Desc:         "Description of Feature G",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "DailyTask",
		execFunc:     backendTasks.NewDailyTask,
	},
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
			defaultConfig := usersSpace.Configs[username].EnchantConfig.GetDefault()
			var existingConfig gameConfig.EnchantConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].EnchantConfig
				if existingConfig.Condition == nil || len(existingConfig.Condition) == 0 {
					existingConfig.Condition = defaultConfig.Condition
				}
			}
			config = existingConfig
		case "AutoMVP":
			defaultConfig := usersSpace.Configs[username].HuntConfig.HuntBossConfig.GetDefault()
			var existingConfig gameConfig.HuntConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].HuntConfig
				if existingConfig.HuntBossConfig == nil {
					existingConfig.HuntBossConfig = &defaultConfig
				}
			} else {
				existingConfig.HuntBossConfig = &defaultConfig
			}
			existingConfig.HuntMonsterConfig = nil

			config = existingConfig
		case "MarketMonitor":
			defaultConfig := usersSpace.Configs[username].TradeMonitorConfig.GetDefault()
			var existingConfig gameConfig.TradeMonitorConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].TradeMonitorConfig
				if len(existingConfig.GetWatchItems()) == 0 {
					existingConfig.WatchItems = defaultConfig.WatchItems
				}
				if len(existingConfig.GetWatchCategories()) == 0 {
					existingConfig.WatchCategories = defaultConfig.WatchCategories
				}
				if existingConfig.MonitorInterval == 0 {
					existingConfig.MonitorInterval = defaultConfig.MonitorInterval
				}
			}

			if existingConfig.BuyItems == nil || len(existingConfig.BuyItems) == 0 {
				existingConfig.BuyItems = make([]gameConfig.PurchaseItem, 1)
			}

			config = existingConfig
		case "AutoHunt":
			defaultConfig := usersSpace.Configs[username].HuntConfig.HuntMonsterConfig.GetDefault()
			var existingConfig gameConfig.HuntConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].HuntConfig
				if existingConfig.HuntMonsterConfig == nil {
					existingConfig.HuntMonsterConfig = &defaultConfig
				}
			} else {
				existingConfig.HuntMonsterConfig = &defaultConfig
			}
			existingConfig.HuntBossConfig = nil

			config = existingConfig

		case "AutoLottery":
			cfg := usersSpace.Configs[username].LotteryConfig.GetDefault()
			defaultConfig := cfg.(gameConfig.LotteryConfig)
			var existingConfig gameConfig.LotteryConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].LotteryConfig
				if existingConfig.DrawCount == 0 {
					existingConfig.DrawCount = defaultConfig.DrawCount
				}
				if len(existingConfig.LotteryType) == 0 {
					existingConfig.LotteryType = defaultConfig.LotteryType
				}
			} else {
				existingConfig = defaultConfig
			}
			config = existingConfig
		case "DailyTask":
			cfg := usersSpace.Configs[username].DailyTaskConfig.GetDefault()
			defaultConfig := cfg.(gameConfig.DailyTaskConfig)
			var existingConfig gameConfig.DailyTaskConfig
			// merge with existing config if any
			if usersSpace.Configs[username] != nil {
				existingConfig = usersSpace.Configs[username].DailyTaskConfig
			} else {
				existingConfig = defaultConfig
			}
			config = existingConfig
		default:
			config = map[string]interface{}{}
		}
	}

	labeled := StructToLabeledJSON(config)

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Data:    labeled,
	})
}

// 自动将 struct 转成带有 label 的 JSON 格式
func StructToLabeledJSON(data interface{}) map[string]interface{} {
	return convertValue(reflect.ValueOf(data))
}

func convertValue(v reflect.Value) map[string]interface{} {
	// 解指针
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	// 必须是 struct 才能继续
	if v.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]interface{})
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fv := v.Field(i)

		// 读取 JSON 字段名
		jsonTag := field.Tag.Get("json")
		jsonKey := strings.Split(jsonTag, ",")[0]
		if jsonKey == "" {
			jsonKey = field.Name
		}

		// label
		label := field.Tag.Get("label")

		// nil 检查
		if (fv.Kind() == reflect.Ptr ||
			fv.Kind() == reflect.Map ||
			fv.Kind() == reflect.Slice ||
			fv.Kind() == reflect.Interface) &&
			fv.IsNil() {
			continue
		}

		var packedValue interface{}

		switch fv.Kind() {

		// ⭐ struct → 递归处理
		case reflect.Struct:
			packedValue = convertValue(fv)

		// ⭐ 指针 → 递归
		case reflect.Ptr:
			packedValue = convertValue(fv)

		// ⭐ slice → 循环处理
		case reflect.Slice:
			arr := make([]interface{}, 0)

			for j := 0; j < fv.Len(); j++ {
				elem := fv.Index(j)
				if elem.Kind() == reflect.Struct || elem.Kind() == reflect.Ptr {
					arr = append(arr, convertValue(elem))
				} else {
					arr = append(arr, elem.Interface())
				}
			}
			packedValue = arr

		// ⭐ map → 保留原样或处理 struct
		case reflect.Map:
			mapResult := map[string]interface{}{}
			for _, key := range fv.MapKeys() {
				item := fv.MapIndex(key)
				if item.Kind() == reflect.Struct || item.Kind() == reflect.Ptr {
					mapResult[key.String()] = convertValue(item)
				} else {
					mapResult[key.String()] = item.Interface()
				}
			}
			packedValue = mapResult

		default:
			packedValue = fv.Interface()
		}

		// 最终包装格式
		result[jsonKey] = map[string]interface{}{
			"value": packedValue,
			"label": label,
		}
	}

	return result
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

		cleanConfig := StripLabelRecursive(req.Config)

		// 把 cleanConfig 转成 JSON
		jsonBytes, _ := json.Marshal(cleanConfig)

		// 塞进你的 config struct
		if err := json.Unmarshal(jsonBytes, &config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Config parse error",
			})
			return
		}

		userConfigs.EnchantConfig = config
	case "AutoHunt", "AutoMVP":
		config := userConfigs.HuntConfig

		cleanConfig := StripLabelRecursive(req.Config)

		// 把 cleanConfig 转成 JSON
		jsonBytes, _ := json.Marshal(cleanConfig)

		// 塞进你的 config struct
		if err := json.Unmarshal(jsonBytes, &config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Config parse error",
			})
			return
		}

		userConfigs.HuntConfig.Merge(config)
	case "MarketMonitor":
		config := userConfigs.TradeMonitorConfig

		cleanConfig := StripLabelRecursive(req.Config)

		// 把 cleanConfig 转成 JSON
		jsonBytes, _ := json.Marshal(cleanConfig)

		// 塞进你的 config struct
		if err := json.Unmarshal(jsonBytes, &config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Config parse error",
			})
			return
		}

		userConfigs.TradeMonitorConfig = config

	case "AutoLottery":
		config := userConfigs.LotteryConfig

		cleanConfig := StripLabelRecursive(req.Config)

		// 把 cleanConfig 转成 JSON
		jsonBytes, _ := json.Marshal(cleanConfig)

		// 塞进你的 config struct
		if err := json.Unmarshal(jsonBytes, &config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Config parse error",
			})
			return
		}

		userConfigs.LotteryConfig = config

	case "DailyTask":
		config := userConfigs.DailyTaskConfig

		cleanConfig := StripLabelRecursive(req.Config)

		// 把 cleanConfig 转成 JSON
		jsonBytes, _ := json.Marshal(cleanConfig)

		// 塞进你的 config struct
		if err := json.Unmarshal(jsonBytes, &config); err != nil {
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Config parse error",
			})
			return
		}

		userConfigs.DailyTaskConfig = config

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

// StripLabelRecursive 递归解析前端传来的 config，移除所有 label 并提取 value
func StripLabelRecursive(v interface{}) interface{} {

	// 如果是 map
	if m, ok := v.(map[string]interface{}); ok {

		// case1: 这是 {value, label} 结构
		if val, exists := m["value"]; exists {
			return StripLabelRecursive(val) // 继续往下拆
		}

		// case2: 普通 map，需要继续递归处理每个字段
		cleaned := make(map[string]interface{})
		for k, v2 := range m {
			cleaned[k] = StripLabelRecursive(v2)
		}
		return cleaned
	}

	// 如果是 array/slice
	if arr, ok := v.([]interface{}); ok {
		newArr := make([]interface{}, len(arr))
		for i, item := range arr {
			newArr[i] = StripLabelRecursive(item)
		}
		return newArr
	}

	// 否则是基本类型（string, int, bool...）直接返回
	return v
}

func handleOptions() http.Handler {
	optionMux := http.NewServeMux()
	optionMux.HandleFunc("/api/options/mini", GetMiniList)
	optionMux.HandleFunc("/api/options/mvp", GetMVPList)
	optionMux.HandleFunc("/api/options/hmvp", GetHMVPList)
	optionMux.HandleFunc("/api/options/map", GetMAPList)
	optionMux.HandleFunc("/api/options/naturetype", GetNatureList)
	optionMux.HandleFunc("/api/options/enchantequippos", GetEnchantEquipPosList)
	optionMux.HandleFunc("/api/options/enchanttype", GetEnchantTypeList)
	optionMux.HandleFunc("/api/options/extras", GetExtrasList)
	optionMux.HandleFunc("/api/options/tradeaction", GetTradeActionList)
	optionMux.HandleFunc("/api/options/watchcategories", GetTradeZhCategoriesList)
	optionMux.HandleFunc("/api/options/lotterytype", GetLotteryTypeList)
	optionMux.HandleFunc("/api/options/dailytaskliefengtype", GetDailyTaskLieFengTypeList)
	return optionMux
}
