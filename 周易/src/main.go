// Package main 实现了一个基于Go语言的周易占卜系统
// 该系统提供HTTP API和WebSocket接口，支持六爻占卜、卦象图片生成等功能
//
// 主要功能：
// - 传统六爻占卜算法实现
// - 卦象图片自动生成（1200x900像素）
// - HTTP RESTful API接口
// - WebSocket实时通信支持
// - 万年历API集成获取干支信息
// - 完整的配置管理和日志记录
//
// 作者: 周易占卜系统开发团队
// 版本: v1.0.0
package main

import (
	"log"
	"sync"
)

// main 是程序的入口函数，负责系统初始化和服务启动
//
// 执行流程：
// 1. 初始化日志系统，设置日志输出格式和文件存储
// 2. 加载配置文件，如果不存在则创建默认配置
// 3. 预热系统缓存，包括字体文件和背景图片的预加载
// 4. 启动HTTP服务器和WebSocket服务
//
// 使用并发方式预加载资源以提高启动速度
func main() {
	// 初始化日志系统
	// 创建日志目录，设置日志文件按日期分割，同时输出到控制台和文件
	initLoggerSimple()

	// 初始化配置系统
	// 从config.json文件加载配置，如果文件不存在则创建默认配置
	// 包括服务器端口、万年历API配置等关键参数
	if err := initConfig(); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 预热系统缓存，提高首次请求的响应速度
	// 使用WaitGroup确保所有预加载任务完成后再启动服务器
	var wg sync.WaitGroup
	wg.Add(2)

	// 并行预加载字体文件
	// 尝试加载项目字体文件或系统字体，为卦象图片生成做准备
	go func() {
		defer wg.Done()
		_, err := loadFontFile()
		if err != nil {
			log.Printf("预加载字体失败: %v", err)
		} else {
			log.Printf("字体文件预加载完成")
		}
	}()

	// 并行预加载背景图片
	// 加载或生成默认的渐变背景图片，用于卦象图片的背景
	go func() {
		defer wg.Done()
		getBackground(ImageWidth, ImageHeight, "images/background.png")
		log.Printf("背景图片预加载完成")
	}()

	// 等待所有预加载任务完成
	wg.Wait()
	log.Printf("系统预热完成，所有缓存资源已就绪")

	// 启动HTTP服务器
	// 包括API路由设置、WebSocket服务初始化、静态文件服务等
	startServer()
}
