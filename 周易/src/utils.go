// utils.go 提供了系统中常用的工具函数
// 包括文件操作、目录管理、路径处理、日志初始化等基础功能
package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// getCurrentDir 获取当前程序的执行目录
// 用于确定相对路径的基准位置，包括配置文件、字体文件、图片文件等的路径计算
//
// 返回值：当前执行目录的绝对路径，失败时返回"."作为默认值
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录失败: %v", err)
		return "." // 获取失败时返回当前目录标识符
	}
	return dir
}

// fileExists 检查指定的文件或目录是否存在
// 常用于配置文件检查、资源文件验证等场景
//
// 参数：
//   - path: 要检查的文件或目录的完整路径
//
// 返回值：存在返回true，不存在返回false
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil // 没有错误表示文件存在
}

// ensureDir 确保指定目录存在，如果不存在则创建
// 支持创建多级目录结构，权限设置为0755（所有者全权限，其他用户读执行权限）
//
// 参数：
//   - dirPath: 要创建的目录路径
//
// 返回值：成功返回nil，失败返回具体错误信息
func ensureDir(dirPath string) error {
	if !fileExists(dirPath) {
		// 使用MkdirAll创建多级目录，权限设置为0755
		return os.MkdirAll(dirPath, 0755)
	}
	return nil // 目录已存在，直接返回成功
}

// extractRiGan 从干支日字符串中提取日干
// 干支日格式通常为"甲子日"、"乙丑日"等，需要提取第一个字符作为日干
// 日干用于确定六神的起始位置
//
// 参数：
//   - ganzhiri: 干支日字符串，如"甲子日"
//
// 返回值：日干字符，如"甲"，解析失败时返回默认值"甲"
func extractRiGan(ganzhiri string) string {
	// 使用rune切片处理中文字符，确保正确截取
	if len(ganzhiri) > 0 {
		runes := []rune(ganzhiri)
		if len(runes) > 0 {
			return string(runes[:1]) // 返回第一个字符（日干）
		}
	}
	return "甲" // 解析失败时返回默认日干
}

// 爻相等 判断两个爻组（三爻或六爻）是否完全相同
// 用于卦象比较和识别，比较每个爻位的阴阳性质
//
// 参数：
//   - 爻1: 第一个爻组数组
//   - 爻2: 第二个爻组数组
//
// 返回值：完全相同返回true，否则返回false
func 爻相等(爻1, 爻2 []int) bool {
	// 长度不同直接返回false
	if len(爻1) != len(爻2) {
		return false
	}
	// 逐个比较每个爻位
	for i := range 爻1 {
		if 爻1[i] != 爻2[i] {
			return false
		}
	}
	return true
}

// hasChangingYao 检查卦象是否包含动爻（变爻）
// 动爻表示会发生变化的爻位，影响卦象的发展趋势
//
// 参数：
//   - 变爻: 变爻标记数组，true表示该爻位为动爻
//
// 返回值：存在动爻返回true，否则返回false
func hasChangingYao(变爻 []bool) bool {
	for _, 变 := range 变爻 {
		if 变 {
			return true // 发现任意一个动爻即返回true
		}
	}
	return false // 没有动爻
}

// max 返回两个整数中的较大值
// 通用的数学工具函数，用于布局计算等场景
//
// 参数：
//   - a: 第一个整数
//   - b: 第二个整数
//
// 返回值：两者中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getYaoWeiName 根据爻位和阴阳性质生成爻位名称
// 按照传统易学命名规则：初二三四五上 + 六（阴爻）/九（阳爻）
//
// 参数：
//   - yaoWei: 爻位序号（1-6，从下到上）
//   - 爻: 爻的阴阳性质（0为阴爻，1为阳爻）
//
// 返回值：标准爻位名称，如"初六"、"九二"等
func getYaoWeiName(yaoWei int, 爻 int) string {
	var name string

	// 根据爻位确定位置名称
	switch yaoWei {
	case 1:
		name = "初" // 第一爻称为"初"
	case 2:
		name = "二" // 第二爻
	case 3:
		name = "三" // 第三爻
	case 4:
		name = "四" // 第四爻
	case 5:
		name = "五" // 第五爻
	case 6:
		name = "上" // 第六爻称为"上"
	}

	// 根据阴阳性质添加数字
	if 爻 == 0 {
		name += "六" // 阴爻用"六"表示
	} else {
		name += "九" // 阳爻用"九"表示
	}

	return name
}

// saveImageToPath 保存图像到指定路径（废弃函数，仅保留接口兼容性）
// 该函数已被saveImageToPathFixed替代，仅用于向后兼容
// 实际的图像保存逻辑现在在image_generator.go中实现
//
// 参数：
//   - img: 要保存的图像对象
//   - fileName: 保存的文件名
//
// 返回值：保存路径和错误信息
func saveImageToPath(img interface{}, fileName string) (string, error) {
	// 尝试多个可能的保存位置，优先级从高到低
	possibleDirs := []string{
		filepath.Join(getCurrentDir(), "photos"), // 首选photos目录
		filepath.Join(getCurrentDir(), "output"), // 备选output目录
		getCurrentDir(),                          // 当前目录
		".",                                      // 相对当前目录
		"./photos",                               // 相对photos目录
		"./output",                               // 相对output目录
	}

	// 确保至少一个目录存在且可写
	var savePath string
	var dirExists bool
	for _, dir := range possibleDirs {
		if err := ensureDir(dir); err == nil {
			savePath = filepath.Join(dir, fileName)
			dirExists = true
			break
		}
	}

	if !dirExists {
		return "", fmt.Errorf("无法创建任何保存目录")
	}

	log.Printf("正在保存图像到: %s", savePath)
	// 注意：这里只是路径处理，实际的图像保存在具体实现中通过imaging.Save完成

	return savePath, nil
}

// initLoggerSimple 初始化简单的日志系统
// 设置日志同时输出到控制台和按日期分割的日志文件
// 日志文件按天创建，格式为app_YYYYMMDD.log
func initLoggerSimple() {
	// 确保log目录存在，用于存储日志文件
	logDir := filepath.Join(getCurrentDir(), "log")
	ensureDir(logDir)

	// 获取当前日期，用于日志文件命名
	today := time.Now().Format("20060102") // YYYYMMDD格式

	// 构建日志文件完整路径
	logPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", today))

	// 创建或打开日志文件，支持追加写入
	// 权限设置为0666，允许所有用户读写
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("创建日志文件失败: %v", err)
		return // 失败时仍然可以输出到控制台
	}

	// 创建多重写入器，同时写入控制台和文件
	// 这样既可以在控制台看到实时日志，又可以保存到文件中
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 配置默认日志记录器
	log.SetOutput(multiWriter)                   // 设置输出目标
	log.SetFlags(log.LstdFlags | log.Lshortfile) // 设置日志格式：时间戳 + 短文件名

	log.Printf("日志系统初始化成功，日志保存到: %s", logPath)
}
