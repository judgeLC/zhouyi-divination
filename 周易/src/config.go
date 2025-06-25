package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Config 应用配置结构体
type Config struct {
	Server   ServerConfig   `json:"server"`   // 服务器配置
	Calendar CalendarConfig `json:"calendar"` // 万年历配置
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string `json:"port"` // HTTP服务器端口
}

// CalendarConfig 万年历配置
type CalendarConfig struct {
	APIHost string `json:"api_host"` // 万年历API地址
	ID      string `json:"id"`       // API ID
	Key     string `json:"key"`      // API密钥
}

// 全局配置变量
var appConfig *Config

// 默认配置
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8090",
		},
		Calendar: CalendarConfig{
			APIHost: "https://cn.apihz.cn",
			ID:      "88888888",
			Key:     "88888888",
		},
	}
}

// 初始化配置
func initConfig() error {
	configPath := "config.json"

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 配置文件不存在，创建默认配置文件
		log.Printf("配置文件不存在，创建默认配置文件: %s", configPath)
		if err := createDefaultConfigFile(configPath); err != nil {
			return fmt.Errorf("创建默认配置文件失败: %v", err)
		}
	}

	// 读取配置文件
	if err := loadConfig(configPath); err != nil {
		return fmt.Errorf("加载配置文件失败: %v", err)
	}

	log.Printf("配置加载成功 - 端口: %s, 万年历API: %s",
		appConfig.Server.Port, appConfig.Calendar.APIHost)

	return nil
}

// 创建默认配置文件
func createDefaultConfigFile(configPath string) error {
	defaultConfig := getDefaultConfig()

	// 将配置结构体转换为JSON
	configJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %v", err)
	}

	// 写入配置文件
	if err := ioutil.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	log.Printf("已创建默认配置文件: %s", configPath)
	return nil
}

// 加载配置文件
func loadConfig(configPath string) error {
	// 读取配置文件内容
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON配置
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置有效性
	if err := validateConfig(&config); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 设置全局配置
	appConfig = &config
	return nil
}

// 验证配置有效性
func validateConfig(config *Config) error {
	// 验证端口配置
	if config.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	// 验证万年历API配置
	if config.Calendar.APIHost == "" {
		return fmt.Errorf("万年历API地址不能为空")
	}

	if config.Calendar.ID == "" {
		return fmt.Errorf("万年历API ID不能为空")
	}

	if config.Calendar.Key == "" {
		return fmt.Errorf("万年历API密钥不能为空")
	}

	return nil
}

// GetConfig 获取当前配置
func GetConfig() *Config {
	if appConfig == nil {
		// 如果配置未初始化，返回默认配置
		log.Println("警告: 配置未初始化，使用默认配置")
		return getDefaultConfig()
	}
	return appConfig
}

// 重新加载配置
func ReloadConfig() error {
	log.Println("重新加载配置文件...")
	return initConfig()
}
