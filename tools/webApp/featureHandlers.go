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
		Desc:         "Description of Feature D",
		Actions:      []string{"Start", "Stop", "Configure", "Log"},
		FunctionName: "MarketMonitor",
		execFunc:     backendTasks.NewTradeMonitorTask,
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

			viewConfig := map[string]interface{}{
				"GameDuration": existingConfig.GameDuration,
				"CarryTeam":    existingConfig.CarryTeam,
				"PrepEliteCD":  existingConfig.PrepEliteCD,
				"Mini":         existingConfig.Mini,
				"MVP":          existingConfig.MVP,
				"HMVP":         existingConfig.HMVP,
			}
			config = viewConfig
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
			config = existingConfig
		case "AutoHunt":
			huntconfig := usersSpace.Configs[username].HuntConfig

			if huntconfig.TargetMonsters == nil {
				huntconfig.TargetMonsters = []string{}
			}
			if huntconfig.TargetItems == nil {
				huntconfig.TargetItems = []string{}
			}

			viewConfig := map[string]interface{}{
				"TimerFly":       huntconfig.TimerFly,
				"PrepEliteCD":    huntconfig.PrepEliteCD,
				"TargetMonsters": huntconfig.TargetMonsters,
				"TargetItems":    huntconfig.TargetItems,
				"Map":            huntconfig.Map,
				"UseDoubleEXP":   huntconfig.UseDoubleEXP,
			}
			config = viewConfig
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
		config.BossInfoParseFromInterface(req.Config)
		userConfigs.HuntConfig = config
	case "MarketMonitor":
		config := userConfigs.TradeMonitorConfig
		config.ParseFromInterface(req.Config)
		userConfigs.TradeMonitorConfig = config
	case "AutoHunt":
		config := userConfigs.HuntConfig
		config.HuntInfoParseFromInterface(req.Config)
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

func GetMiniList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"狸猫",
		"蓝疯兔",
		"波利之王",
		"摇滚蝗虫",
		"幽灵波利",
		"蛙王",
		"直升机哥布灵",
		"龙蝇",
		"流浪之狼",
		"枯树精",
		"狮鹫兽",
		"安毕斯",
		"妖君",
		"兽人婴儿",
		"南瓜先生",
		"半龙人",
		"草精",
		"鹗枭首领",
		"爱丽丝女仆",
		"艾斯恩魔女",
		"弑神者",
		"迷幻之王",
		"大笨钟",
		"钟塔守护者",
		"魔灵娃娃",
		"炎之小魔女",
		"吹笛人",
		"银月魔女",
		"暗赛尼亚",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetMVPList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"天使波利",
		"黄金虫",
		"恶魔波利",
		"海盗之王",
		"海神",
		"哥布灵首领",
		"蜂后",
		"蚁后",
		"皮里恩",
		"虎王",
		"俄塞里斯",
		"月夜猫",
		"兽人英雄",
		"犬妖首领",
		"死灵",
		"阿特罗斯",
		"兽人酋长",
		"迪塔勒泰晤勒",
		"鹗枭男爵",
		"血腥骑士",
		"巴风特",
		"黑暗之王",
		"时间管理人",
		"斯佩夏尔",
		"冰暴骑士",
		"炎之领主卡浩",
		"圣天使波利",
		"凯特莉娜",
		"艾勒梅斯",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetHMVPList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"卡仑",
		"狼外婆",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}

func GetHMAPList(w http.ResponseWriter, r *http.Request) {
	list := []string{
		"普隆德拉",
		"普隆德拉南门",
		"普隆德拉西门",
		"迷藏森林",
		"伊斯鲁得岛",
		"沉船",
		"幽灵船",
		"海底洞窟岛",
		"海底神殿",
		"吉芬",
		"妙勒尼山脉",
		"摩洛克",
		"金字塔 1F",
		"斐扬",
		"斐扬南门",
		"兽人村落",
		"古城郊外",
		"古城",
		"斐扬森林",
		"哥布林森林",
		"科德森林",
		"苏克拉特沙漠",
		"普隆德拉北门",
		"阿尔德巴朗",
		"普隆德拉大厅 1F",
		"姜饼城",
		"玩具工厂 1F",
		"波利岛",
		"斐扬森林南部",
		"兽人村落南部",
		"古城外围",
		"天水之国·安塔修",
		"尤诺",
		"边境检查站",
		"艾因布洛克原野",
		"熔岩洞窟 1F",
		"熔岩洞窟 2F",
		"熔岩洞窟 3F",
		"莱斯特灯塔",
		"尼芙海姆",
		"迷雾森林",
		"骷髅洞穴",
		"哈姆林",
		"乌帕拉",
		"里希塔尔岑",
		"里希塔尔岑平原",
		"生体地下1F",
		"生体地下2F",
		"生体地下3F",
		"伊达平原",
		"瑞秋",
		"拉萨尼亚",
		"多拉多岛",
		"意大利饺森林",
		"洛阳",
		"夕阳海岸",
		"荒境",
		"月之湖",
		"伊克莱基",
		"时间花园",
		"克雷普特学院",
		"星泪森林",
		"绽放之地",
		"风之森",
		"科摩多",
		"可可蒙海滩",
		"流星森林",
		"阿尔贝塔",
		"海龟岛",
		"古城之泪",
		"副本·极限挑战",
		"深渊之湖",

		// 特殊地图
		"高级房间",
		"皇家料理间",
		"公会领地",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    list,
	})
}
