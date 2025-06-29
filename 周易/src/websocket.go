package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// WebSocket升级器
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源，生产环境应该严格验证
		return true
	},
}

// WebSocket连接管理器
type WSManager struct {
	clients    map[string]*WSClient // 客户端映射
	broadcast  chan WSMessage
	register   chan *WSClient
	unregister chan *WSClient
	mu         sync.RWMutex
}

// WebSocket客户端
type WSClient struct {
	ID   string
	Conn *websocket.Conn
	Send chan WSMessage
}

var wsManager *WSManager

// 初始化WebSocket管理器
func initWSManager() {
	wsManager = &WSManager{
		clients:    make(map[string]*WSClient),
		broadcast:  make(chan WSMessage, 256),
		register:   make(chan *WSClient),
		unregister: make(chan *WSClient),
	}
	go wsManager.run()
}

// 运行WebSocket管理器
func (m *WSManager) run() {
	ticker := time.NewTicker(30 * time.Second) // 心跳检测
	defer ticker.Stop()

	for {
		select {
		case client := <-m.register:
			m.mu.Lock()
			m.clients[client.ID] = client
			m.mu.Unlock()

			log.Printf("WebSocket客户端已连接: %s", client.ID)

			// 发送欢迎消息
			welcomeMsg := WSMessage{
				Type: WSEventConnect,
				Data: WSConnection{
					ID: client.ID,
				},
			}
			select {
			case client.Send <- welcomeMsg:
			default:
				close(client.Send)
			}

		case client := <-m.unregister:
			m.mu.Lock()
			if _, ok := m.clients[client.ID]; ok {
				delete(m.clients, client.ID)
				close(client.Send)
				log.Printf("WebSocket客户端已断开: %s", client.ID)
			}
			m.mu.Unlock()

		case message := <-m.broadcast:
			m.mu.RLock()
			clientsToRemove := make([]string, 0)

			for id, client := range m.clients {
				select {
				case client.Send <- message:
					// 消息发送成功
				default:
					// 客户端发送通道已满或已关闭，标记为需要移除
					log.Printf("客户端 %s 发送通道阻塞，将断开连接", id)
					clientsToRemove = append(clientsToRemove, id)
				}
			}
			m.mu.RUnlock()

			// 移除问题客户端
			if len(clientsToRemove) > 0 {
				m.mu.Lock()
				for _, id := range clientsToRemove {
					if client, ok := m.clients[id]; ok {
						client.Conn.Close()
						close(client.Send)
						delete(m.clients, id)
						log.Printf("已移除问题客户端: %s", id)
					}
				}
				m.mu.Unlock()
			}

		case <-ticker.C:
			// 定期清理无效连接和发送心跳包
			m.mu.RLock()
			activeClients := len(m.clients)
			m.mu.RUnlock()

			log.Printf("心跳检测：当前活跃客户端数量: %d", activeClients)

			// 发送心跳包
			heartbeat := WSMessage{
				Type: WSEventHeartbeat,
				Data: map[string]interface{}{
					"timestamp": time.Now().Unix(),
					"server":    "yijing-api",
					"clients":   activeClients,
				},
			}
			m.broadcast <- heartbeat
		}
	}
}

// WebSocket连接处理器
func handleWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v", err)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("client_%d", time.Now().UnixNano())

	client := &WSClient{
		ID:   clientID,
		Conn: conn,
		Send: make(chan WSMessage, 256),
	}

	// 注册客户端
	wsManager.register <- client

	// 启动读写协程
	go client.writePump()
	go client.readPump()
}

// OneBot WebSocket连接处理器
func handleOneBotWSConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("OneBot WebSocket升级失败: %v", err)
		return
	}

	// 生成客户端ID
	clientID := fmt.Sprintf("onebot_%d", time.Now().UnixNano())

	// 创建基础WebSocket客户端
	wsClient := &WSClient{
		ID:   clientID,
		Conn: conn,
		Send: make(chan WSMessage, 256),
	}

	// 创建OneBot客户端（使用默认机器人ID）
	oneBotClient := NewOneBotClient(wsClient, 12345678)

	// 注册到OneBot管理器
	if oneBotManager != nil {
		oneBotManager.clients[clientID] = oneBotClient
	}

	log.Printf("OneBot客户端已连接: %s, BotId: %d", clientID, oneBotClient.BotId)

	// 发送连接生命周期事件
	oneBotClient.SendLifecycleEvent("connect")

	// 启动协程处理OneBot消息
	go oneBotClient.readOneBotPump()
	go oneBotClient.writeOneBotPump()
}

// 读取消息
func (c *WSClient) readPump() {
	defer func() {
		wsManager.unregister <- c
		c.Conn.Close()
		log.Printf("客户端 %s 读取协程已退出", c.ID)
	}()

	// 设置读取超时为120秒，比之前的60秒更宽松
	c.Conn.SetReadDeadline(time.Now().Add(120 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		return nil
	})

	for {
		var msg WSMessage
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("客户端 %s WebSocket异常关闭: %v", c.ID, err)
			} else {
				log.Printf("客户端 %s 读取消息失败: %v", c.ID, err)
			}
			break
		}

		// 处理客户端消息
		c.handleMessage(msg)
	}
}

// 写入消息
func (c *WSClient) writePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
		log.Printf("客户端 %s 写入协程已退出", c.ID)
	}()

	for {
		select {
		case message, ok := <-c.Send:
			// 设置写入超时为30秒，比之前的10秒更宽松
			c.Conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteJSON(message); err != nil {
				log.Printf("客户端 %s 写入消息失败: %v", c.ID, err)
				return
			}

		case <-ticker.C:
			// 设置ping超时为30秒
			c.Conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("客户端 %s 发送ping失败: %v", c.ID, err)
				return
			}
		}
	}
}

// 处理客户端消息
func (c *WSClient) handleMessage(msg WSMessage) {
	switch msg.Type {
	case WSEventDivine:
		// 处理占卜请求
		go c.handleDivineRequest(msg)
	case WSEventHeartbeat:
		// 心跳响应
		response := WSMessage{
			Type: WSEventHeartbeat,
			Data: map[string]interface{}{
				"timestamp": time.Now().Unix(),
				"client_id": c.ID,
			},
		}
		c.Send <- response
	default:
		log.Printf("未知消息类型: %s", msg.Type)
	}
}

// 处理占卜请求
func (c *WSClient) handleDivineRequest(msg WSMessage) {
	// 生成卦象图片
	imagePath, err := generateTodayGua()
	if err != nil {
		errorMsg := WSMessage{
			Type: WSEventError,
			Data: map[string]interface{}{
				"error":   err.Error(),
				"message": "生成卦象失败",
			},
		}
		c.Send <- errorMsg
		return
	}

	// 构建完整的图片URL
	config := GetConfig()
	port := config.Server.Port
	fullImageURL := fmt.Sprintf("http://localhost:%s/%s", port, imagePath)

	// 创建响应
	now := time.Now()
	result := DivineResult{
		ID:        fmt.Sprintf("divine_%d", now.UnixNano()),
		Date:      now.Format("2006-01-02"),
		ImagePath: fullImageURL, // 返回完整的图片URL
		CreatedAt: now.Unix(),
	}

	response := WSMessage{
		Type: WSEventDivine,
		Data: result,
	}

	c.Send <- response
}

// 广播消息到所有连接的客户端
func BroadcastMessage(msgType string, data interface{}) {
	if wsManager != nil {
		msg := WSMessage{
			Type: msgType,
			Data: data,
		}
		wsManager.broadcast <- msg
	}
}

// 获取连接的客户端数量
func GetConnectedClientsCount() int {
	if wsManager == nil {
		return 0
	}
	wsManager.mu.RLock()
	defer wsManager.mu.RUnlock()
	return len(wsManager.clients)
}

// WebSocket状态API
func handleWSStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"connected_clients": GetConnectedClientsCount(),
		"server_time":       time.Now().Unix(),
		"websocket_enabled": true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ApiResponse{
		Code:    200,
		Message: "成功",
		Data:    status,
	})
}
