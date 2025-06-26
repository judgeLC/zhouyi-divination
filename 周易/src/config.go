// config.go 实现了应用程序的配置管理功能
// 支持JSON格式的配置文件，自动创建默认配置，配置验证等功能
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// Config 应用程序主配置结构体
// 包含服务器配置和万年历API配置两个主要部分
type Config struct {
	Server   ServerConfig   `json:"server"`   // HTTP服务器相关配置
	Calendar CalendarConfig `json:"calendar"` // 万年历API相关配置
}

// ServerConfig HTTP服务器配置结构体
// 定义服务器运行的基本参数
type ServerConfig struct {
	Port string `json:"port"` // HTTP服务器监听端口，默认为8090
}

// CalendarConfig 万年历API配置结构体
// 用于获取干支纪年、干支纪月、干支纪日等信息
type CalendarConfig struct {
	APIHost string `json:"api_host"` // 万年历API服务器地址
	ID      string `json:"id"`       // API访问ID，用于身份认证
	Key     string `json:"key"`      // API访问密钥，用于身份认证
}

// appConfig 全局配置变量，存储当前应用程序的配置信息
// 通过initConfig()函数初始化，通过GetConfig()函数获取
var appConfig *Config

// getDefaultConfig 返回系统默认配置
// 当配置文件不存在或配置初始化失败时使用此默认配置
//
// 默认配置说明：
// - 服务器端口：8090
// - 万年历API：使用测试API地址和默认密钥
//
// 返回值：包含默认设置的Config结构体指针
func getDefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8090", // 默认HTTP服务端口
		},
		Calendar: CalendarConfig{
			APIHost: "https://cn.apihz.cn", // 万年历API服务地址
			ID:      "88888888",            // 测试用API ID
			Key:     "88888888",            // 测试用API密钥
		},
	}
}

// initConfig 初始化应用程序配置系统
//
// 执行流程：
// 1. 检查config.json配置文件是否存在
// 2. 如果不存在，创建包含默认配置的config.json文件
// 3. 加载并解析配置文件
// 4. 验证配置参数的有效性
// 5. 设置全局配置变量
//
// 返回值：成功返回nil，失败返回具体错误信息
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

	// 读取并解析配置文件
	if err := loadConfig(configPath); err != nil {
		return fmt.Errorf("加载配置文件失败: %v", err)
	}

	// 输出配置加载成功信息，便于调试和运维
	log.Printf("配置加载成功 - 服务器端口: %s, 万年历API地址: %s",
		appConfig.Server.Port, appConfig.Calendar.APIHost)

	return nil
}

// createDefaultConfigFile 创建包含默认配置的JSON配置文件
// 当配置文件不存在时自动调用，确保程序能正常启动
//
// 参数：
//   - configPath: 配置文件的完整路径
//
// 返回值：成功返回nil，失败返回具体错误信息
func createDefaultConfigFile(configPath string) error {
	// 获取默认配置结构体
	defaultConfig := getDefaultConfig()

	// 将配置结构体转换为格式化的JSON字符串
	// 使用4个空格缩进，提高可读性
	configJSON, err := json.MarshalIndent(defaultConfig, "", "    ")
	if err != nil {
		return fmt.Errorf("序列化默认配置失败: %v", err)
	}

	// 写入配置文件，设置文件权限为0644（所有者可读写，其他用户只读）
	if err := ioutil.WriteFile(configPath, configJSON, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	log.Printf("已创建默认配置文件: %s", configPath)
	return nil
}

// loadConfig 从指定路径加载并解析配置文件
// 执行配置文件的读取、JSON解析、配置验证等完整流程
//
// 参数：
//   - configPath: 配置文件的完整路径
//
// 返回值：成功返回nil，失败返回具体错误信息
func loadConfig(configPath string) error {
	// 读取配置文件的原始内容
	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON格式的配置数据到Config结构体
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 验证配置参数的有效性和完整性
	if err := validateConfig(&config); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 将验证通过的配置设置为全局配置
	appConfig = &config
	return nil
}

// validateConfig 验证配置参数的有效性和完整性
// 确保所有必需的配置项都已正确设置
//
// 参数：
//   - config: 待验证的配置结构体指针
//
// 返回值：验证通过返回nil，失败返回具体错误信息
func validateConfig(config *Config) error {
	// 验证服务器端口配置
	if config.Server.Port == "" {
		return fmt.Errorf("服务器端口不能为空")
	}

	// 验证万年历API服务器地址
	if config.Calendar.APIHost == "" {
		return fmt.Errorf("万年历API地址不能为空")
	}

	// 验证万年历API访问ID
	if config.Calendar.ID == "" {
		return fmt.Errorf("万年历API ID不能为空")
	}

	// 验证万年历API访问密钥
	if config.Calendar.Key == "" {
		return fmt.Errorf("万年历API密钥不能为空")
	}

	return nil
}

// GetConfig 获取当前应用程序的配置信息
// 线程安全的配置获取函数，如果配置未初始化则返回默认配置
//
// 返回值：当前配置的Config结构体指针
func GetConfig() *Config {
	if appConfig == nil {
		// 配置未初始化的情况下，返回默认配置并记录警告日志
		log.Println("警告: 配置未初始化，使用默认配置")
		return getDefaultConfig()
	}
	return appConfig
}

// ReloadConfig 重新加载配置文件
// 用于运行时动态更新配置，无需重启程序
//
// 返回值：成功返回nil，失败返回具体错误信息
func ReloadConfig() error {
	log.Println("重新加载配置文件...")
	return initConfig()
}
