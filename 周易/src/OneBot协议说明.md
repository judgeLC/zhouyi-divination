# OneBot协议说明 - 周易占卜系统

## 概述

本系统已集成OneBot v11协议支持，可以作为OneBot实现与任何支持OneBot协议的机器人框架进行通信。

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

### 消息发送（模拟）

> **注意：** 当前版本的消息发送功能为模拟实现，主要用于测试OneBot协议的兼容性。实际的消息发送需要根据具体的聊天平台进行实现。

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

#### 2. 心跳事件
每5秒发送一次心跳事件：

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

系统可以发送占卜结果作为消息事件，模拟机器人发送占卜信息：

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

## 错误响应

当动作执行失败时，会返回错误响应：

```json
{
    "status": "failed",
    "retcode": 1404,
    "message": "不支持的动作",
    "wording": "动作 unknown_action 未实现",
    "echo": "test_error"
}
```

### 常见错误码

- `1400`: 请求参数错误
- `1404`: 动作不存在
- `1500`: 内部服务器错误

## 消息段类型支持

系统支持以下消息段类型：

- `text`: 纯文本
- `image`: 图片
- `at`: @某人
- `reply`: 回复

**消息段示例：**
```json
[
    {
        "type": "text",
        "data": {
            "text": "你好！"
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
    }
]
```

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
};

ws.onerror = function(error) {
    console.error('WebSocket错误:', error);
};
```

### Python客户端示例

```python
import asyncio
import websockets
import json

async def onebot_client():
    uri = "ws://localhost:8080/onebot/ws"
    
    async with websockets.connect(uri) as websocket:
        print("OneBot连接已建立")
        
        # 发送获取状态请求
        request = {
            "action": "get_status",
            "params": {},
            "echo": "status_check"
        }
        await websocket.send(json.dumps(request))
        
        # 接收响应
        response = await websocket.recv()
        data = json.loads(response)
        print("收到响应:", data)

# 运行客户端
asyncio.run(onebot_client())
```

## 配置机器人ID

默认机器人ID为 `12345678`，可以在代码中修改：

```go
// 在 handleOneBotWSConnection 函数中
oneBotClient := NewOneBotClient(wsClient, 你的机器人ID)
```

## 测试建议

1. 使用提供的测试页面 `http://localhost:8080/onebot/test` 进行基础功能测试
2. 通过WebSocket连接测试各种动作的调用和响应
3. 观察心跳事件和生命周期事件的正常发送
4. 测试错误处理机制

## 注意事项

1. 当前实现主要用于OneBot协议兼容性测试
2. 消息发送功能为模拟实现，实际使用时需要根据具体平台进行扩展
3. 心跳间隔为5秒，可根据需要调整
4. 系统支持多个WebSocket客户端同时连接
5. 所有时间戳使用Unix时间戳格式

## 扩展开发

如需添加新的动作或事件类型，请参考：

1. `onebot_types.go` - 添加新的数据结构和常量
2. `onebot_handler.go` - 添加新的动作处理器
3. 在 `registerActionHandlers()` 方法中注册新的处理器

## 技术支持

如有问题或建议，请通过以下方式联系：
- 查看项目文档
- 提交Issue
- 参与讨论 