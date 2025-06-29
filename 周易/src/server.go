package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

// API处理函数 - 处理"今日卦象"请求
func handleDivineRequest(w http.ResponseWriter, r *http.Request) {
	// 1. 打印请求体 (可选，但有助于调试)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "无法读取请求体", http.StatusBadRequest)
		return
	}
	defer r.Body.Close() // 确保关闭请求体，避免资源泄漏
	fmt.Printf("Request Body: %s\n", string(bodyBytes))

	var req DivineRequest
	decoder := json.NewDecoder(io.NopCloser(bytes.NewBuffer(bodyBytes))) // 使用 io.NopCloser 避免重复关闭

	// 2. JSON 解码和错误处理
	err = decoder.Decode(&req)
	if err != nil {
		// 更细致的错误处理，根据错误类型返回不同的 HTTP 状态码
		switch {
		case err == io.EOF:
			http.Error(w, "请求体为空", http.StatusBadRequest)
		case err.(*json.SyntaxError) != nil:
			http.Error(w, fmt.Sprintf("JSON 语法错误: %s", err), http.StatusBadRequest)
		case err.(*json.UnmarshalTypeError) != nil:
			http.Error(w, fmt.Sprintf("类型不匹配: %s", err), http.StatusBadRequest)
		default:
			log.Printf("JSON 解码错误: %v", err) // 记录错误到日志
			http.Error(w, "服务器内部错误", http.StatusInternalServerError)
		}
		return
	}
	fmt.Printf("Received request body: %+v\n", req) // 打印解码后的数据

	// 生成卦象图片
	imagePath, err := generateTodayGua()
	if err != nil {
		http.Error(w, "生成卦象失败: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 构建完整的图片URL
	config := GetConfig()
	port := config.Server.Port
	fullImageURL := fmt.Sprintf("http://localhost:%s/%s", port, imagePath)

	// 创建响应
	now := time.Now()
	divineResult := DivineResult{
		ID:        fmt.Sprintf("divine_%d", now.UnixNano()),
		Date:      now.Format("2006-01-02"),
		ImagePath: fullImageURL, // 返回完整的图片URL
		CreatedAt: now.Unix(),
	}

	response := ApiResponse{
		Code:    200,
		Message: "成功",
		Data:    divineResult,
	}

	// 广播到WebSocket客户端
	BroadcastMessage(WSEventDivine, divineResult)

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// 新增API路由处理
func setupAPIRoutes() {
	http.HandleFunc("/api/divine", handleDivineRequest)
	http.HandleFunc("/ws", handleWSConnection)                // WebSocket连接端点
	http.HandleFunc("/onebot/ws", handleOneBotWSConnection)   // OneBot WebSocket连接端点
	http.HandleFunc("/api/ws/status", handleWSStatus)         // WebSocket状态查询
	http.HandleFunc("/api/onebot/status", handleOneBotStatus) // OneBot状态查询
	http.HandleFunc("/test", serveWebSocketTestPage)          // WebSocket测试页面
	http.HandleFunc("/onebot/test", serveOneBotTestPage)      // OneBot测试页面
}

// 提供WebSocket测试页面
func serveWebSocketTestPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "websocket_test.html")
}

// 提供OneBot测试页面
func serveOneBotTestPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "onebot_test.html")
}

// OneBot状态API
func handleOneBotStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"connected_clients": GetOneBotClientsCount(),
		"server_time":       time.Now().Unix(),
		"onebot_enabled":    true,
		"onebot_version":    OneBotVersion,
		"implementation":    OneBotImpl,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Code:    200,
		Message: "成功",
		Data:    status,
	})
}

// 提供图像访问
func serveImageFiles() {
	photosDir := filepath.Join(getCurrentDir(), "photos")
	if err := ensureDir(photosDir); err == nil {
		fs := http.FileServer(http.Dir(photosDir))
		http.Handle("/photos/", http.StripPrefix("/photos/", fs))
	}

	outputDir := filepath.Join(getCurrentDir(), "output")
	if err := ensureDir(outputDir); err == nil {
		fs := http.FileServer(http.Dir(outputDir))
		http.Handle("/output/", http.StripPrefix("/output/", fs))
	}
}

// 启动HTTP服务器
func startServer() {
	// 初始化WebSocket管理器
	initWSManager()

	// 初始化OneBot管理器
	initOneBotManager()

	// 设置API路由
	setupAPIRoutes()

	// 提供图像文件访问
	serveImageFiles()

	// 从配置文件获取端口
	config := GetConfig()
	port := config.Server.Port
	log.Printf("启动HTTP服务器，监听端口 %s...", port)
	log.Printf("API接口路径: http://localhost:%s/api/divine", port)
	log.Printf("WebSocket接口路径: ws://localhost:%s/ws", port)
	log.Printf("OneBot WebSocket接口路径: ws://localhost:%s/onebot/ws", port)
	log.Printf("WebSocket状态查询: http://localhost:%s/api/ws/status", port)
	log.Printf("OneBot状态查询: http://localhost:%s/api/onebot/status", port)
	log.Printf("WebSocket测试页面: http://localhost:%s/test", port)
	log.Printf("OneBot测试页面: http://localhost:%s/onebot/test", port)

	// 启动HTTP服务器
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("启动HTTP服务器失败: %v", err)
	}
}
