package main

import (
	"log"
	"sync"
)

func main() {
	// 初始化日志系统
	initLoggerSimple()

	// 预热缓存
	var wg sync.WaitGroup
	wg.Add(2)

	// 并行预加载字体和背景
	go func() {
		defer wg.Done()
		_, err := loadFontFile()
		if err != nil {
			log.Printf("预加载字体失败: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		getBackground(ImageWidth, ImageHeight, "images/background.png")
	}()

	// 等待预加载完成
	wg.Wait()

	// 启动HTTP服务器
	startServer()
}
