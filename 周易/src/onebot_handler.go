// onebot_handler.go OneBot协议处理器
// 实现OneBot v11标准的动作处理和事件发送功能
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// OneBot客户端扩展结构
type OneBotClient struct {
	*WSClient
	BotId    int64                                                               // 机器人ID
	Stats    *Stats                                                              // 统计信息
	Status   *OneBotStatus                                                       // 状态信息
	Handlers map[string]func(*OneBotClient, *OneBotAction) *OneBotActionResponse // 动作处理器
}

// OneBot管理器
type OneBotManager struct {
	*WSManager
	clients map[string]*OneBotClient // OneBot客户端映射
}

var oneBotManager *OneBotManager

// 初始化OneBot管理器
func initOneBotManager() {
	oneBotManager = &OneBotManager{
		WSManager: wsManager,
		clients:   make(map[string]*OneBotClient),
	}

	// 注册基本的动作处理器
	registerDefaultActionHandlers()
}

// 创建OneBot客户端
func NewOneBotClient(wsClient *WSClient, botId int64) *OneBotClient {
	client := &OneBotClient{
		WSClient: wsClient,
		BotId:    botId,
		Stats: &Stats{
			PacketReceived:  0,
			PacketSent:      0,
			PacketLost:      0,
			MessageReceived: 0,
			MessageSent:     0,
			DisconnectTimes: 0,
			LostTimes:       0,
		},
		Status: &OneBotStatus{
			AppInitialized: true,
			AppEnabled:     true,
			PluginsGood:    true,
			AppGood:        true,
			Online:         true,
			Good:           true,
		},
		Handlers: make(map[string]func(*OneBotClient, *OneBotAction) *OneBotActionResponse),
	}

	// 注册客户端特定的处理器
	client.registerActionHandlers()

	return client
}

// 注册动作处理器
func (c *OneBotClient) registerActionHandlers() {
	c.Handlers[ActionSendPrivateMsg] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleSendPrivateMsg(action)
	}
	c.Handlers[ActionSendGroupMsg] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleSendGroupMsg(action)
	}
	c.Handlers[ActionSendMsg] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleSendMsg(action)
	}
	c.Handlers[ActionGetLoginInfo] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleGetLoginInfo(action)
	}
	c.Handlers[ActionGetStatus] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleGetStatus(action)
	}
	c.Handlers[ActionGetVersionInfo] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleGetVersionInfo(action)
	}
	c.Handlers[ActionCanSendImage] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleCanSendImage(action)
	}
	c.Handlers[ActionCanSendRecord] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
		return client.handleCanSendRecord(action)
	}
}

// 处理OneBot消息
func (c *OneBotClient) handleOneBotMessage(rawMessage []byte) {
	var action OneBotAction
	if err := json.Unmarshal(rawMessage, &action); err != nil {
		log.Printf("OneBot客户端 %s 解析动作失败: %v", c.ID, err)
		return
	}

	c.Stats.PacketReceived++

	// 查找并执行动作处理器
	if handler, exists := c.Handlers[action.Action]; exists {
		response := handler(c, &action)
		c.sendActionResponse(response)
	} else {
		// 未知动作
		response := &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1404,
			Message: "不支持的动作",
			Wording: fmt.Sprintf("动作 %s 未实现", action.Action),
			Echo:    action.Echo,
		}
		c.sendActionResponse(response)
	}
}

// 发送动作响应
func (c *OneBotClient) sendActionResponse(response *OneBotActionResponse) {
	responseData, err := json.Marshal(response)
	if err != nil {
		log.Printf("序列化动作响应失败: %v", err)
		return
	}

	// 发送响应
	c.Conn.WriteMessage(1, responseData) // 1 表示文本消息
	c.Stats.PacketSent++
}

// 发送OneBot事件
func (c *OneBotClient) SendEvent(event interface{}) error {
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("序列化事件失败: %v", err)
	}

	err = c.Conn.WriteMessage(1, eventData)
	if err != nil {
		c.Stats.PacketLost++
		return fmt.Errorf("发送事件失败: %v", err)
	}

	c.Stats.PacketSent++
	return nil
}

// 处理发送私聊消息
func (c *OneBotClient) handleSendPrivateMsg(action *OneBotAction) *OneBotActionResponse {
	params := action.Params

	userId, ok := params["user_id"]
	if !ok {
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "缺少参数 user_id",
			Echo:    action.Echo,
		}
	}

	message, ok := params["message"]
	if !ok {
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "缺少参数 message",
			Echo:    action.Echo,
		}
	}

	// 这里应该实现实际的消息发送逻辑
	log.Printf("发送私聊消息给 %v: %v", userId, message)

	// 生成消息ID
	messageId := time.Now().Unix()

	c.Stats.MessageSent++

	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: map[string]interface{}{
			"message_id": messageId,
		},
		Echo: action.Echo,
	}
}

// 处理发送群消息
func (c *OneBotClient) handleSendGroupMsg(action *OneBotAction) *OneBotActionResponse {
	params := action.Params

	groupId, ok := params["group_id"]
	if !ok {
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "缺少参数 group_id",
			Echo:    action.Echo,
		}
	}

	message, ok := params["message"]
	if !ok {
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "缺少参数 message",
			Echo:    action.Echo,
		}
	}

	// 这里应该实现实际的消息发送逻辑
	log.Printf("发送群消息到 %v: %v", groupId, message)

	// 生成消息ID
	messageId := time.Now().Unix()

	c.Stats.MessageSent++

	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: map[string]interface{}{
			"message_id": messageId,
		},
		Echo: action.Echo,
	}
}

// 处理发送消息（通用）
func (c *OneBotClient) handleSendMsg(action *OneBotAction) *OneBotActionResponse {
	params := action.Params

	messageType, ok := params["message_type"].(string)
	if !ok {
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "缺少参数 message_type",
			Echo:    action.Echo,
		}
	}

	switch messageType {
	case MessageTypePrivate:
		return c.handleSendPrivateMsg(action)
	case MessageTypeGroup:
		return c.handleSendGroupMsg(action)
	default:
		return &OneBotActionResponse{
			Status:  "failed",
			Retcode: 1400,
			Message: "不支持的消息类型",
			Echo:    action.Echo,
		}
	}
}

// 处理获取登录信息
func (c *OneBotClient) handleGetLoginInfo(action *OneBotAction) *OneBotActionResponse {
	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: &LoginInfo{
			UserId:   c.BotId,
			Nickname: "周易占卜机器人",
		},
		Echo: action.Echo,
	}
}

// 处理获取状态
func (c *OneBotClient) handleGetStatus(action *OneBotAction) *OneBotActionResponse {
	c.Status.Stat = c.Stats

	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data:    c.Status,
		Echo:    action.Echo,
	}
}

// 处理获取版本信息
func (c *OneBotClient) handleGetVersionInfo(action *OneBotAction) *OneBotActionResponse {
	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: &VersionInfo{
			AppName:         OneBotImpl,
			AppVersion:      "1.0.0",
			ProtocolVersion: OneBotVersion,
			OneBotVersion:   OneBotVersion,
		},
		Echo: action.Echo,
	}
}

// 处理是否可以发送图片
func (c *OneBotClient) handleCanSendImage(action *OneBotAction) *OneBotActionResponse {
	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: map[string]interface{}{
			"yes": true,
		},
		Echo: action.Echo,
	}
}

// 处理是否可以发送语音
func (c *OneBotClient) handleCanSendRecord(action *OneBotAction) *OneBotActionResponse {
	return &OneBotActionResponse{
		Status:  "success",
		Retcode: 0,
		Data: map[string]interface{}{
			"yes": false,
		},
		Echo: action.Echo,
	}
}

// 注册默认动作处理器
func registerDefaultActionHandlers() {
	// 这里可以注册全局的动作处理器
}

// 发送生命周期事件
func (c *OneBotClient) SendLifecycleEvent(subType string) error {
	event := &LifecycleEvent{
		OneBotEvent: OneBotEvent{
			Time:     time.Now().Unix(),
			SelfId:   c.BotId,
			PostType: EventTypeMeta,
		},
		SubType: subType,
	}

	return c.SendEvent(event)
}

// 发送心跳事件
func (c *OneBotClient) SendHeartbeatEvent() error {
	event := &HeartbeatEvent{
		OneBotEvent: OneBotEvent{
			Time:     time.Now().Unix(),
			SelfId:   c.BotId,
			PostType: EventTypeMeta,
		},
		Status:   c.Status,
		Interval: 5000, // 5秒间隔
	}

	return c.SendEvent(event)
}

// 发送占卜结果作为群消息事件
func (c *OneBotClient) SendDivineResultAsGroupMessage(groupId int64, result *DivineResult) error {
	// 创建消息段
	var message Message
	message = append(message, NewTextSegment(fmt.Sprintf("今日卦象：%s", result.BenGua)))
	if result.ImagePath != "" {
		message = append(message, NewImageSegment(result.ImagePath))
	}
	message = append(message, NewTextSegment(fmt.Sprintf("日期：%s", result.Date)))

	event := &GroupMessageEvent{
		OneBotEvent: OneBotEvent{
			Time:        time.Now().Unix(),
			SelfId:      c.BotId,
			PostType:    EventTypeMessage,
			MessageType: MessageTypeGroup,
			SubType:     "normal",
			MessageId:   int32(time.Now().Unix()),
			UserId:      c.BotId,
			GroupId:     groupId,
			Message:     message,
			RawMessage:  message.String(),
			Font:        0,
			Sender: &Sender{
				UserId:   c.BotId,
				Nickname: "周易占卜机器人",
				Card:     "周易占卜机器人",
				Role:     "member",
			},
		},
	}

	c.Stats.MessageSent++
	return c.SendEvent(event)
}

// 发送占卜结果作为私聊消息事件
func (c *OneBotClient) SendDivineResultAsPrivateMessage(userId int64, result *DivineResult) error {
	// 创建消息段
	var message Message
	message = append(message, NewTextSegment(fmt.Sprintf("今日卦象：%s", result.BenGua)))
	if result.ImagePath != "" {
		message = append(message, NewImageSegment(result.ImagePath))
	}
	message = append(message, NewTextSegment(fmt.Sprintf("日期：%s", result.Date)))

	event := &PrivateMessageEvent{
		OneBotEvent: OneBotEvent{
			Time:        time.Now().Unix(),
			SelfId:      c.BotId,
			PostType:    EventTypeMessage,
			MessageType: MessageTypePrivate,
			SubType:     "friend",
			MessageId:   int32(time.Now().Unix()),
			UserId:      userId,
			Message:     message,
			RawMessage:  message.String(),
			Font:        0,
			Sender: &Sender{
				UserId:   c.BotId,
				Nickname: "周易占卜机器人",
				Role:     "friend",
			},
		},
	}

	c.Stats.MessageSent++
	return c.SendEvent(event)
}

// 广播OneBot事件到所有客户端
func BroadcastOneBotEvent(event interface{}) {
	if oneBotManager == nil {
		return
	}

	for _, client := range oneBotManager.clients {
		if err := client.SendEvent(event); err != nil {
			log.Printf("向客户端 %s 发送事件失败: %v", client.ID, err)
		}
	}
}

// 获取OneBot客户端数量
func GetOneBotClientsCount() int {
	if oneBotManager == nil {
		return 0
	}
	return len(oneBotManager.clients)
}

// OneBot读取协程
func (c *OneBotClient) readOneBotPump() {
	defer func() {
		// 从OneBot管理器中移除客户端
		if oneBotManager != nil {
			delete(oneBotManager.clients, c.ID)
		}
		c.Conn.Close()
		log.Printf("OneBot客户端 %s 读取协程已退出", c.ID)
	}()

	// 设置读取超时
	c.Conn.SetReadDeadline(time.Now().Add(120 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(120 * time.Second))
		return nil
	})

	// 发送连接成功事件
	c.SendLifecycleEvent("connect")

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("OneBot客户端 %s WebSocket异常关闭: %v", c.ID, err)
			} else {
				log.Printf("OneBot客户端 %s 读取消息失败: %v", c.ID, err)
			}
			break
		}

		// 处理OneBot消息
		c.handleOneBotMessage(message)
	}
}

// OneBot写入协程
func (c *OneBotClient) writeOneBotPump() {
	ticker := time.NewTicker(54 * time.Second)
	heartbeatTicker := time.NewTicker(5 * time.Second) // OneBot心跳间隔
	defer func() {
		ticker.Stop()
		heartbeatTicker.Stop()
		c.Conn.Close()
		log.Printf("OneBot客户端 %s 写入协程已退出", c.ID)
	}()

	for {
		select {
		case <-ticker.C:
			// 发送WebSocket ping
			c.Conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("OneBot客户端 %s 发送ping失败: %v", c.ID, err)
				return
			}

		case <-heartbeatTicker.C:
			// 发送OneBot心跳事件
			if err := c.SendHeartbeatEvent(); err != nil {
				log.Printf("OneBot客户端 %s 发送心跳事件失败: %v", c.ID, err)
			}
		}
	}
}
