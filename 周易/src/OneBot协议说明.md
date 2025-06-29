# OneBot协议说明 - 周易占卜系统

## 概述

本系统已集成OneBot v11协议支持，可以作为OneBot实现与任何支持OneBot协议的机器人框架进行通信。OneBot是一个聊天机器人应用接口标准，旨在统一不同聊天平台上的机器人应用开发接口，使开发者只需编写一次业务逻辑代码即可应用到多种机器人平台。

## 系统架构

### 核心组件

1. **OneBot类型定义层** (`onebot_types.go`)
   - 完整的OneBot v11协议数据结构
   - 事件、动作、消息段类型定义
   - 错误码和状态常量

2. **协议处理层** (`onebot_handler.go`)
   - 动作请求处理器
   - 事件发送管理器
   - WebSocket连接管理
   - 统计信息收集

3. **WebSocket通信层** (`websocket.go`)
   - 独立的OneBot WebSocket端点
   - 协议兼容的消息处理
   - 连接生命周期管理

4. **HTTP API层** (`server.go`)
   - OneBot状态查询接口
   - 测试页面服务
   - 路由配置管理

## 连接信息

### WebSocket连接地址
```
ws://localhost:8080/onebot/ws
```
（端口号根据配置文件中的设置可能不同）

### HTTP API状态查询
```
GET http://localhost:8080/api/onebot/status
```

### 测试页面
```
http://localhost:8080/onebot/test
```

## 支持的动作（Actions）

### 基础信息获取

#### 1. 获取登录信息
```json
{
    "action": "get_login_info",
    "params": {},
    "echo": "test1"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "user_id": 12345678,
        "nickname": "周易占卜机器人"
    },
    "echo": "test1"
}
```

#### 2. 获取状态信息
```json
{
    "action": "get_status",
    "params": {},
    "echo": "test2"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "app_initialized": true,
        "app_enabled": true,
        "plugins_good": true,
        "app_good": true,
        "online": true,
        "good": true,
        "stat": {
            "packet_received": 10,
            "packet_sent": 15,
            "packet_lost": 0,
            "message_received": 5,
            "message_sent": 8,
            "disconnect_times": 0,
            "lost_times": 0
        }
    },
    "echo": "test2"
}
```

#### 3. 获取版本信息
```json
{
    "action": "get_version_info",
    "params": {},
    "echo": "test3"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "app_name": "yijing-onebot",
        "app_version": "1.0.0",
        "protocol_version": "11",
        "onebot_version": "11"
    },
    "echo": "test3"
}
```

### 功能查询

#### 4. 查询是否可以发送图片
```json
{
    "action": "can_send_image",
    "params": {},
    "echo": "test4"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "yes": true
    },
    "echo": "test4"
}
```

#### 5. 查询是否可以发送语音
```json
{
    "action": "can_send_record",
    "params": {},
    "echo": "test5"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "yes": false
    },
    "echo": "test5"
}
```

### 消息发送（模拟实现）

**重要说明：** 当前版本的消息发送功能为模拟实现，主要用于测试OneBot协议的兼容性。实际的消息发送需要根据具体的聊天平台（如QQ、微信等）进行实现。

#### 6. 发送私聊消息
```json
{
    "action": "send_private_msg",
    "params": {
        "user_id": 123456,
        "message": "Hello, this is a private message!"
    },
    "echo": "test6"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "message_id": 1640995200
    },
    "echo": "test6"
}
```

#### 7. 发送群消息
```json
{
    "action": "send_group_msg",
    "params": {
        "group_id": 123456,
        "message": "Hello, this is a group message!"
    },
    "echo": "test7"
}
```

**响应示例：**
```json
{
    "status": "success",
    "retcode": 0,
    "data": {
        "message_id": 1640995201
    },
    "echo": "test7"
}
```

#### 8. 通用发送消息
```json
{
    "action": "send_msg",
    "params": {
        "message_type": "private",
        "user_id": 123456,
        "message": "Hello via send_msg!"
    },
    "echo": "test8"
}
```

## 支持的事件（Events）

### 元事件

#### 1. 生命周期事件
当客户端连接或断开时，会发送生命周期事件：

```json
{
    "time": 1640995200,
    "self_id": 12345678,
    "post_type": "meta_event",
    "meta_event_type": "lifecycle",
    "sub_type": "connect"
}
```

**支持的生命周期类型：**
- `connect`: 连接建立
- `enable`: 启用
- `disable`: 禁用

#### 2. 心跳事件
每5秒发送一次心跳事件，用于保持连接活跃和状态同步：

```json
{
    "time": 1640995200,
    "self_id": 12345678,
    "post_type": "meta_event",
    "meta_event_type": "heartbeat",
    "status": {
        "app_initialized": true,
        "app_enabled": true,
        "plugins_good": true,
        "app_good": true,
        "online": true,
        "good": true,
        "stat": {
            "packet_received": 10,
            "packet_sent": 15,
            "packet_lost": 0,
            "message_received": 5,
            "message_sent": 8,
            "disconnect_times": 0,
            "lost_times": 0
        }
    },
    "interval": 5000
}
```

### 消息事件（扩展功能）

系统可以发送占卜结果作为消息事件，模拟机器人发送占卜信息。这为集成聊天机器人功能提供了基础。

#### 群消息事件示例
```json
{
    "time": 1640995200,
    "self_id": 12345678,
    "post_type": "message",
    "message_type": "group",
    "sub_type": "normal",
    "message_id": 1001,
    "user_id": 12345678,
    "group_id": 123456,
    "message": [
        {
            "type": "text",
            "data": {
                "text": "今日卦象：乾为天"
            }
        },
        {
            "type": "image",
            "data": {
                "file": "http://localhost:8080/output/today_gua_20231201.png"
            }
        },
        {
            "type": "text",
            "data": {
                "text": "日期：2023-12-01"
            }
        }
    ],
    "raw_message": "今日卦象：乾为天[CQ:image,file=http://localhost:8080/output/today_gua_20231201.png]日期：2023-12-01",
    "font": 0,
    "sender": {
        "user_id": 12345678,
        "nickname": "周易占卜机器人",
        "card": "周易占卜机器人",
        "role": "member"
    }
}
```

#### 私聊消息事件示例
```json
{
    "time": 1640995200,
    "self_id": 12345678,
    "post_type": "message",
    "message_type": "private",
    "sub_type": "friend",
    "message_id": 1002,
    "user_id": 123456,
    "message": [
        {
            "type": "text",
            "data": {
                "text": "您的占卜结果已生成"
            }
        }
    ],
    "raw_message": "您的占卜结果已生成",
    "font": 0,
    "sender": {
        "user_id": 12345678,
        "nickname": "周易占卜机器人",
        "role": "friend"
    }
}
```

## 错误处理

### 错误响应格式

当动作执行失败时，会返回标准的错误响应：

```json
{
    "status": "failed",
    "retcode": 1404,
    "message": "不支持的动作",
    "wording": "动作 unknown_action 未实现",
    "echo": "test_error"
}
```

### 错误码定义

- **1400**: 请求参数错误，缺少必需参数或参数格式不正确
- **1404**: 动作不存在，请求的动作未实现
- **1500**: 内部服务器错误，系统处理出现异常

### 常见错误场景

1. **缺少参数**: 发送消息时未提供user_id或group_id
2. **参数格式错误**: 提供的JSON格式不正确
3. **不支持的动作**: 请求未实现的OneBot动作
4. **网络连接问题**: WebSocket连接异常断开

## 消息段类型支持

系统支持OneBot v11标准的消息段类型：

### 基础消息段

- **text**: 纯文本消息
- **image**: 图片消息（支持URL和本地文件）
- **at**: @某人
- **reply**: 回复特定消息

### 消息段构造示例

```json
[
    {
        "type": "text",
        "data": {
            "text": "你好！今日卦象为："
        }
    },
    {
        "type": "image",
        "data": {
            "file": "http://example.com/image.jpg"
        }
    },
    {
        "type": "at",
        "data": {
            "qq": 123456
        }
    },
    {
        "type": "reply",
        "data": {
            "id": 12345
        }
    }
]
```

### CQ码支持

系统支持CQ码格式的消息解析和生成：

```
[CQ:text,text=你好]
[CQ:image,file=http://example.com/image.jpg]
[CQ:at,qq=123456]
[CQ:reply,id=12345]
```

## 技术实现细节

### 并发安全

- 使用读写锁保护客户端连接映射
- 通道机制确保消息传递的线程安全
- 统计信息使用原子操作更新

### 内存管理

- 连接断开时自动清理相关资源
- 心跳检测机制及时发现无效连接
- 消息缓冲区大小限制防止内存泄漏

### 性能优化

- WebSocket连接池管理
- 消息序列化缓存
- 异步事件处理

## 使用示例

### JavaScript WebSocket客户端

```javascript
const ws = new WebSocket('ws://localhost:8080/onebot/ws');

ws.onopen = function() {
    console.log('OneBot连接已建立');
    
    // 发送获取登录信息请求
    ws.send(JSON.stringify({
        action: 'get_login_info',
        params: {},
        echo: 'login_info_' + Date.now()
    }));
};

ws.onmessage = function(event) {
    const data = JSON.parse(event.data);
    console.log('收到响应或事件:', data);
    
    // 处理不同类型的消息
    if (data.status) {
        // 动作响应
        console.log('动作响应:', data);
    } else if (data.post_type) {
        // 事件上报
        handleOneBotEvent(data);
    }
};

ws.onerror = function(error) {
    console.error('WebSocket错误:', error);
};

ws.onclose = function(event) {
    console.log('连接已关闭:', event.code, event.reason);
};

function handleOneBotEvent(event) {
    switch (event.post_type) {
        case 'meta_event':
            if (event.meta_event_type === 'heartbeat') {
                console.log('收到心跳:', event.status);
            }
            break;
        case 'message':
            console.log('收到消息事件:', event);
            break;
    }
}
```

### Python客户端示例

```python
import asyncio
import websockets
import json
import logging

# 配置日志
logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

class OneBotClient:
    def __init__(self, uri):
        self.uri = uri
        self.websocket = None
        
    async def connect(self):
        """建立WebSocket连接"""
        try:
            self.websocket = await websockets.connect(self.uri)
            logger.info("OneBot连接已建立")
            return True
        except Exception as e:
            logger.error(f"连接失败: {e}")
            return False
    
    async def send_action(self, action, params=None, echo=None):
        """发送OneBot动作"""
        if not self.websocket:
            logger.error("WebSocket未连接")
            return None
            
        request = {
            "action": action,
            "params": params or {},
            "echo": echo or f"{action}_{int(time.time())}"
        }
        
        try:
            await self.websocket.send(json.dumps(request))
            logger.info(f"发送动作: {action}")
            return True
        except Exception as e:
            logger.error(f"发送失败: {e}")
            return False
    
    async def listen(self):
        """监听消息"""
        if not self.websocket:
            logger.error("WebSocket未连接")
            return
            
        try:
            async for message in self.websocket:
                data = json.loads(message)
                await self.handle_message(data)
        except websockets.exceptions.ConnectionClosed:
            logger.info("连接已关闭")
        except Exception as e:
            logger.error(f"监听错误: {e}")
    
    async def handle_message(self, data):
        """处理收到的消息"""
        if 'status' in data:
            # 动作响应
            logger.info(f"动作响应: {data}")
        elif 'post_type' in data:
            # 事件上报
            await self.handle_event(data)
    
    async def handle_event(self, event):
        """处理OneBot事件"""
        post_type = event.get('post_type')
        
        if post_type == 'meta_event':
            meta_type = event.get('meta_event_type')
            if meta_type == 'heartbeat':
                logger.debug("收到心跳事件")
            elif meta_type == 'lifecycle':
                logger.info(f"生命周期事件: {event.get('sub_type')}")
        elif post_type == 'message':
            logger.info(f"收到消息事件: {event}")

async def main():
    client = OneBotClient("ws://localhost:8080/onebot/ws")
    
    if await client.connect():
        # 发送一些测试动作
        await client.send_action("get_login_info")
        await client.send_action("get_status")
        await client.send_action("get_version_info")
        
        # 开始监听
        await client.listen()

if __name__ == "__main__":
    import time
    asyncio.run(main())
```

### Go客户端示例

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"
    
    "github.com/gorilla/websocket"
)

type OneBotClient struct {
    conn *websocket.Conn
    url  string
}

func NewOneBotClient(url string) *OneBotClient {
    return &OneBotClient{url: url}
}

func (c *OneBotClient) Connect() error {
    var err error
    c.conn, _, err = websocket.DefaultDialer.Dial(c.url, nil)
    if err != nil {
        return fmt.Errorf("连接失败: %v", err)
    }
    log.Println("OneBot连接已建立")
    return nil
}

func (c *OneBotClient) SendAction(action string, params map[string]interface{}, echo string) error {
    if echo == "" {
        echo = fmt.Sprintf("%s_%d", action, time.Now().Unix())
    }
    
    request := map[string]interface{}{
        "action": action,
        "params": params,
        "echo":   echo,
    }
    
    return c.conn.WriteJSON(request)
}

func (c *OneBotClient) Listen() {
    for {
        var msg map[string]interface{}
        err := c.conn.ReadJSON(&msg)
        if err != nil {
            log.Printf("读取消息失败: %v", err)
            break
        }
        
        c.handleMessage(msg)
    }
}

func (c *OneBotClient) handleMessage(msg map[string]interface{}) {
    if status, ok := msg["status"]; ok {
        // 动作响应
        log.Printf("动作响应: %v", msg)
        return
    }
    
    if postType, ok := msg["post_type"]; ok {
        // 事件上报
        log.Printf("收到事件: %s", postType)
    }
}

func main() {
    client := NewOneBotClient("ws://localhost:8080/onebot/ws")
    
    if err := client.Connect(); err != nil {
        log.Fatal(err)
    }
    defer client.conn.Close()
    
    // 发送测试动作
    client.SendAction("get_login_info", nil, "")
    client.SendAction("get_status", nil, "")
    
    // 开始监听
    client.Listen()
}
```

## 配置和部署

### 系统要求

- Go 1.21 或更高版本
- 支持WebSocket的现代浏览器（用于测试页面）
- 网络访问权限（用于万年历API调用）

### 配置文件

OneBot功能集成在主系统配置中，无需额外配置文件。相关设置包括：

```json
{
    "server": {
        "port": "8080"
    }
}
```

### 部署步骤

1. **克隆项目**
```bash
git clone https://github.com/judgeLC/zhouyi-divination.git
cd zhouyi-divination/周易/src
```

2. **安装依赖**
```bash
go mod tidy
```

3. **编译运行**
```bash
go build -o yijing.exe
./yijing.exe
```

4. **验证部署**
- 访问 `http://localhost:8080/onebot/test` 测试OneBot功能
- 查看 `http://localhost:8080/api/onebot/status` 确认状态

### Docker部署

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY 周易/src/ .
RUN go mod tidy && go build -o yijing .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/yijing .
COPY --from=builder /app/config.json .
COPY --from=builder /app/ttf/ ./ttf/

EXPOSE 8080
CMD ["./yijing"]
```

### 配置机器人ID

默认机器人ID为 `12345678`，可以在代码中修改：

```go
// 在 websocket.go 的 handleOneBotWSConnection 函数中
oneBotClient := NewOneBotClient(wsClient, 你的机器人ID)
```

### 性能调优

#### 连接数限制
```go
// 在 onebot_handler.go 中调整
const maxConnections = 100
```

#### 心跳间隔
```go
// 在 onebot_handler.go 中调整
heartbeatTicker := time.NewTicker(5 * time.Second)
```

#### 消息缓冲区
```go
// 在 websocket.go 中调整
Send: make(chan WSMessage, 256)
```

## 测试和调试

### 使用测试页面

访问 `http://localhost:8080/onebot/test` 使用内置的测试页面：

1. 点击"连接"建立WebSocket连接
2. 选择要测试的动作类型
3. 输入相应的参数（JSON格式）
4. 点击"发送动作"查看响应
5. 观察消息日志中的事件流

### 常见问题排查

#### 连接失败
- 检查服务器是否正常启动
- 确认端口号配置正确
- 检查防火墙设置

#### 动作无响应
- 查看服务器日志输出
- 确认JSON格式正确
- 检查必需参数是否完整

#### 心跳中断
- 检查网络连接稳定性
- 调整心跳间隔设置
- 查看WebSocket连接状态

### 日志级别配置

```go
// 在 main.go 中设置日志级别
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

## 扩展开发

### 添加新动作

1. **定义动作常量**
```go
// 在 onebot_types.go 中添加
const ActionCustomAction = "custom_action"
```

2. **实现处理函数**
```go
// 在 onebot_handler.go 中添加
func (c *OneBotClient) handleCustomAction(action *OneBotAction) *OneBotActionResponse {
    // 处理逻辑
    return &OneBotActionResponse{
        Status:  "success",
        Retcode: 0,
        Data:    responseData,
        Echo:    action.Echo,
    }
}
```

3. **注册处理器**
```go
// 在 registerActionHandlers 方法中添加
c.Handlers[ActionCustomAction] = func(client *OneBotClient, action *OneBotAction) *OneBotActionResponse {
    return client.handleCustomAction(action)
}
```

### 添加新事件类型

1. **定义事件结构**
```go
// 在 onebot_types.go 中添加
type CustomEvent struct {
    OneBotEvent
    CustomData string `json:"custom_data"`
}
```

2. **实现发送方法**
```go
// 在 onebot_handler.go 中添加
func (c *OneBotClient) SendCustomEvent(data string) error {
    event := &CustomEvent{
        OneBotEvent: OneBotEvent{
            Time:     time.Now().Unix(),
            SelfId:   c.BotId,
            PostType: "custom",
        },
        CustomData: data,
    }
    return c.SendEvent(event)
}
```

### 集成实际聊天平台

为了将系统集成到实际的聊天平台（如QQ、微信等），需要：

1. **实现平台特定的消息发送**
```go
func (c *OneBotClient) handleSendPrivateMsg(action *OneBotAction) *OneBotActionResponse {
    // 替换为实际的平台API调用
    realMessageId, err := sendToRealPlatform(params)
    if err != nil {
        return &OneBotActionResponse{
            Status:  "failed",
            Retcode: 1500,
            Message: err.Error(),
            Echo:    action.Echo,
        }
    }
    
    return &OneBotActionResponse{
        Status:  "success",
        Retcode: 0,
        Data:    map[string]interface{}{"message_id": realMessageId},
        Echo:    action.Echo,
    }
}
```

2. **监听平台消息并转换为OneBot事件**
```go
func convertPlatformMessageToOneBotEvent(platformMsg PlatformMessage) *GroupMessageEvent {
    return &GroupMessageEvent{
        OneBotEvent: OneBotEvent{
            Time:        platformMsg.Timestamp,
            SelfId:      botId,
            PostType:    EventTypeMessage,
            MessageType: MessageTypeGroup,
            // ... 其他字段映射
        },
    }
}
```

### 中间件支持

```go
type OneBotMiddleware func(*OneBotClient, *OneBotAction) (*OneBotActionResponse, bool)

func (c *OneBotClient) Use(middleware OneBotMiddleware) {
    // 添加中间件逻辑
}

// 示例：认证中间件
func AuthMiddleware(client *OneBotClient, action *OneBotAction) (*OneBotActionResponse, bool) {
    // 验证逻辑
    if !isAuthorized(action) {
        return &OneBotActionResponse{
            Status:  "failed",
            Retcode: 1403,
            Message: "未授权的操作",
            Echo:    action.Echo,
        }, false
    }
    return nil, true
}
```

## 兼容性说明

### OneBot版本支持

- **OneBot v11**: 完全支持
- **OneBot v12**: 计划支持
- **向后兼容**: 支持旧版本客户端

### 平台兼容性

- **go-cqhttp**: 完全兼容
- **nonebot**: 兼容
- **mirai**: 兼容
- **其他OneBot实现**: 理论兼容

### 浏览器支持

测试页面支持的浏览器：
- Chrome 80+
- Firefox 76+
- Safari 13+
- Edge 80+

## 安全注意事项

### 访问控制

- 建议在生产环境中启用身份验证
- 限制WebSocket连接的来源IP
- 使用HTTPS/WSS加密连接

### 数据保护

- 不在日志中记录敏感信息
- 定期清理过期的连接和缓存
- 验证输入参数防止注入攻击

### 配置安全

```go
// 建议的安全配置
const (
    MaxConnections = 50
    MaxMessageSize = 1024 * 1024 // 1MB
    ReadTimeout    = 60 * time.Second
    WriteTimeout   = 10 * time.Second
)
```

## 社区和支持

### 相关资源

- **OneBot官方文档**: https://onebot.dev/
- **Go语言官方文档**: https://golang.org/doc/
- **WebSocket协议规范**: https://tools.ietf.org/html/rfc6455

### 贡献指南

欢迎提交Issue和Pull Request：
1. Fork项目仓库
2. 创建功能分支
3. 提交代码和测试
4. 发起Pull Request

### 许可证

本项目遵循开源许可证，具体请查看项目根目录的LICENSE文件。 