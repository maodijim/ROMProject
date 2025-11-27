package main

import (
	"encoding/json"
	"net/http"
	"os"

	"ROMProject/config"

	"gopkg.in/yaml.v3"
)

var (
	configs = make(map[string]*config.ServerConfigs)
)

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	RoleNum  int    `json:"roleNum"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func handleUserAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	var userReq UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Invalid request body",
		})
		return
	}

	if r.Method != http.MethodDelete {
		if userReq.RoleNum < 1 || userReq.RoleNum > 3 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Role number must be between 1 and 3",
			})
			return
		}
	}

	// create empty config.yml if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		emptyFile, err := os.Create(configPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to create config file",
			})
			return
		}
		defer emptyFile.Close()
		initConfig := make(map[string]*config.ServerConfigs) // Initialize an empty map
		data, err := yaml.Marshal(&initConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to initialize config file",
			})
			return
		}
		if _, err := emptyFile.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to write to config file",
			})
			return
		}
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to read config file",
		})
		return
	}

	if err := yaml.Unmarshal(data, &configs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to parse config file",
		})
		return
	}

	if len(configs) == 0 {
		sc := config.NewServerConfigs("config.yml")
		sc.Username = userReq.Username
		sc.Password = userReq.Password
		sc.Char = uint(userReq.RoleNum)
		configs[userReq.Username] = sc
	} else if r.Method == http.MethodPost {
		if _, exists := configs[userReq.Username]; exists {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "User already exists",
			})
			return
		}
		sc := config.NewServerConfigs("config.yml")
		sc.Username = userReq.Username
		sc.Password = userReq.Password
		sc.Char = uint(userReq.RoleNum)
		configs[userReq.Username] = sc
	} else if r.Method == http.MethodPut {
		if _, exists := configs[userReq.Username]; !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "User not found",
			})
			return
		}
		configs[userReq.Username].Password = userReq.Password
		configs[userReq.Username].Char = uint(userReq.RoleNum)
	} else if r.Method == http.MethodDelete {
		if _, exists := configs[userReq.Username]; !exists {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "User not found",
			})
			return
		}
		task := getFeatureTask(userReq.Username)
		if task != nil {
			task.StopTask()
		}
		delete(configs, userReq.Username)
	}

	err = saveConfigs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	message := "User added successfully"
	if r.Method == http.MethodPut {
		message = "User updated successfully"
	}
	if r.Method == http.MethodDelete {
		message = "User deleted successfully"
	}

	json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: message,
		Data: map[string]interface{}{
			"username": userReq.Username,
			"roleNum":  userReq.RoleNum,
		},
	})
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	// create empty config.yml if not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		emptyFile, err := os.Create(configPath)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to create config file",
			})
			return
		}
		defer emptyFile.Close()
		initConfig := make(map[string]*config.ServerConfigs) // Initialize an empty map
		data, err := yaml.Marshal(&initConfig)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to initialize config file",
			})
			return
		}
		if _, err := emptyFile.Write(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(Response{
				Success: false,
				Message: "Failed to write to config file",
			})
			return
		}
	}

	// Read the config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to read config file",
		})
		return
	}

	if err := yaml.Unmarshal(data, &configs); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(Response{
			Success: false,
			Message: "Failed to parse config file",
		})
		return
	}

	res := Response{
		Success: true,
		Message: "User retrieved successfully",
	}
	userList := make([]map[string]interface{}, 0)
	for username, conf := range configs {
		userList = append(userList, map[string]interface{}{
			"username": username,
			"roleNum":  conf.Char,
		})
	}
	res.Data = userList

	if len(configs) == 0 {
		res.Data = []interface{}{}
	}
	json.NewEncoder(w).Encode(res)
}
