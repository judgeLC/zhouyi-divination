package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// 获取当前执行目录
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("获取当前目录失败: %v", err)
		return "."
	}
	return dir
}

// 检查文件是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 创建或确保目录存在
func ensureDir(dirPath string) error {
	if !fileExists(dirPath) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// 提取日干
func extractRiGan(ganzhiri string) string {
	// 假设ganzhiri的格式是 "甲子日"，提取第一个字 "甲"
	if len(ganzhiri) > 0 {
		return string([]rune(ganzhiri)[:1])
	}
	return "甲" // 默认返回"甲"
}

// 辅助函数：判断两个爻组是否相等
func 爻相等(爻1, 爻2 []int) bool {
	if len(爻1) != len(爻2) {
		return false
	}
	for i := range 爻1 {
		if 爻1[i] != 爻2[i] {
			return false
		}
	}
	return true
}

// 检查卦象是否包含变爻
func hasChangingYao(变爻 []bool) bool {
	for _, 变 := range 变爻 {
		if 变 {
			return true
		}
	}
	return false
}

// 返回两个整数中的较大值
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// 获取爻位名称
func getYaoWeiName(yaoWei int, 爻 int) string {
	var name string
	switch yaoWei {
	case 1:
		name = "初"
	case 2:
		name = "二"
	case 3:
		name = "三"
	case 4:
		name = "四"
	case 5:
		name = "五"
	case 6:
		name = "上"
	}

	if 爻 == 0 {
		name += "六"
	} else {
		name += "九"
	}

	return name
}

// 保存图像到指定路径
func saveImageToPath(img interface{}, fileName string) (string, error) {
	// 尝试多个可能的保存位置
	possibleDirs := []string{
		filepath.Join(getCurrentDir(), "photos"),
		filepath.Join(getCurrentDir(), "output"),
		getCurrentDir(),
		".",
		"./photos",
		"./output",
	}

	// 确保至少一个目录存在
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
	// 这里需要根据实际图像类型进行保存，在具体实现中会调用imaging.Save

	return savePath, nil
}

// 简单的日志初始化函数
func initLoggerSimple() {
	// 确保log目录存在
	logDir := filepath.Join(getCurrentDir(), "log")
	ensureDir(logDir)

	// 获取当前日期
	today := time.Now().Format("20060102")

	// 创建日志文件路径
	logPath := filepath.Join(logDir, fmt.Sprintf("app_%s.log", today))

	// 创建或打开日志文件
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("创建日志文件失败: %v", err)
		return
	}

	// 创建多重写入器，同时输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	// 设置默认日志记录器输出到文件和控制台
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("日志系统初始化成功，日志保存到: %s", logPath)
}
