# 🔮 周易WebSocket反向接口说明

## 📖 概述
周易占卜系统现已支持WebSocket反向接口，可以实现实时双向通信，支持客户端连接后接收服务器主动推送的卦象信息。

## 🚀 接口地址
- **WebSocket连接地址**: `ws://localhost:8090/ws`
- **状态查询API**: `http://localhost:8090/api/ws/status`
- **测试页面**: `http://localhost:8090/test`

## 📡 WebSocket消息格式

### 发送消息格式
```json
{
    "type": "消息类型",
    "data": {
        // 消息数据
    }
}
```

### 接收消息格式
```json
{
    "type": "消息类型", 
    "data": {
        // 响应数据
    }
}
```

## 🎯 支持的消息类型

### 1. 连接事件 (connect)
**方向**: 服务器 → 客户端  
**触发时机**: 客户端成功连接到WebSocket服务器时  
**数据格式**:
```json
{
    "type": "connect",
    "data": {
        "id": "client_1234567890"
    }
}
```

### 2. 断开连接事件 (disconnect)
**方向**: 服务器 → 客户端  
**触发时机**: 连接断开时

### 3. 占卜请求 (divine)
**方向**: 双向  

**客户端发送**:
```json
{
    "type": "divine",
    "data": {
        "request_time": "2024-01-01T12:00:00Z"
    }
}
```

**注意**: `imagepath` 字段现在返回完整的HTTP URL，可直接在浏览器中访问或用于图片显示。

**服务器响应**:
```json
{
    "type": "divine",
    "data": {
        "id": "divine_1234567890",
        "date": "2024-01-01",
        "imagepath": "http://localhost:8090/photos/卜卦_20240101102217.png",
        "created_at": 1704110400,
        "ganzhinian": "甲辰年",
        "ganzhiyue": "乙亥月",
        "ganzhiri": "丙子日",
        "bengua": "乾为天",
        "benguadesc": "乾卦描述...",
        "biangua": "变卦名",
        "bianguadesc": "变卦描述...",
        "hasdonyao": true
    }
}
```

### 4. 心跳检测 (heartbeat)
**方向**: 双向  
**触发时机**: 每30秒自动发送，也可手动发送

**客户端发送**:
```json
{
    "type": "heartbeat",
    "data": {
        "timestamp": 1704110400
    }
}
```

**服务器响应**:
```json
{
    "type": "heartbeat",
    "data": {
        "timestamp": 1704110400,
        "server": "yijing-api"
    }
}
```

### 5. 错误信息 (error)
**方向**: 服务器 → 客户端  
**触发时机**: 处理请求出错时

```json
{
    "type": "error",
    "data": {
        "error": "具体错误信息",
        "message": "生成卦象失败"
    }
}
```

## 🔄 反向推送特性

当有新的卦象通过HTTP API (`/api/divine`) 生成时，服务器会自动向所有连接的WebSocket客户端广播卦象结果，实现真正的"反向推送"功能。

## 📊 状态查询API

**请求**: `GET /api/ws/status`

**响应**:
```json
{
    "code": 200,
    "message": "成功",
    "data": {
        "connected_clients": 3,
        "server_time": 1704110400,
        "websocket_enabled": true
    }
}
```

## 🛠️ 使用示例

### JavaScript示例
```javascript
// 连接WebSocket
const ws = new WebSocket('ws://localhost:8090/ws');

// 连接成功
ws.onopen = function(event) {
    console.log('连接成功');
    
    // 发送占卜请求
    ws.send(JSON.stringify({
        type: 'divine',
        data: { request_time: new Date().toISOString() }
    }));
};

// 接收消息
ws.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log('收到消息:', message);
    
    switch(message.type) {
        case 'connect':
            console.log('客户端ID:', message.data.id);
            break;
        case 'divine':
            console.log('卦象结果:', message.data);
            break;
        case 'heartbeat':
            console.log('心跳:', message.data);
            break;
        case 'error':
            console.error('错误:', message.data);
            break;
    }
};

// 连接关闭
ws.onclose = function(event) {
    console.log('连接关闭:', event.code, event.reason);
};

// 连接错误 
ws.onerror = function(error) {
    console.error('连接错误:', error);
};
```

### Python示例
```python
import asyncio
import websockets
import json

async def client():
    uri = "ws://localhost:8090/ws"
    
    async with websockets.connect(uri) as websocket:
        # 发送占卜请求
        await websocket.send(json.dumps({
            "type": "divine",
            "data": {"request_time": "2024-01-01T12:00:00Z"}
        }))
        
        # 接收消息
        async for message in websocket:
            data = json.loads(message)
            print(f"收到消息: {data}")

# 运行客户端
asyncio.run(client())
```

## 🧪 测试工具

系统提供了内置的WebSocket测试页面，访问 `http://localhost:8090/test` 即可：

- ✅ 可视化连接状态
- 📤 发送各类测试消息  
- 📥 实时显示接收的消息
- 💓 心跳检测测试
- 📊 服务器状态查询

## ⚠️ 注意事项

1. **连接超时**: WebSocket连接会在120秒无活动后自动断开
2. **写入超时**: 消息写入超时设置为30秒
3. **心跳机制**: 服务器每30秒发送心跳包，客户端应响应pong帧
4. **消息缓冲**: 每个客户端的发送缓冲区限制为256条消息
5. **并发连接**: 支持多个客户端同时连接
6. **错误处理**: 自动检测和清理失效连接
7. **安全性**: 当前允许所有来源连接，生产环境需要设置合适的`CheckOrigin`

## 🔧 配置选项

WebSocket相关配置可以通过修改代码中的常量进行调整：

- 心跳间隔: 30秒
- 读取超时: 120秒 (已优化)
- 写入超时: 30秒 (已优化)
- 消息缓冲区: 256条
- Ping间隔: 54秒
- 连接清理: 自动检测失效连接

## 🚀 启动服务

```bash
cd 周易/src
go run .
```

启动后会看到类似输出：
```
启动HTTP服务器，监听端口 8090...
API接口路径: http://localhost:8090/api/divine
WebSocket接口路径: ws://localhost:8090/ws  
WebSocket状态查询: http://localhost:8090/api/ws/status
WebSocket测试页面: http://localhost:8090/test
```

## 📈 应用场景

1. **实时卦象推送**: 多个客户端可同时接收新生成的卦象
2. **状态监控**: 实时监控服务器状态和连接数
3. **交互式占卜**: 支持客户端主动请求卦象
4. **系统集成**: 可与其他系统进行实时数据交换

## 🎯 扩展建议

- 添加用户认证机制
- 实现房间/频道功能
- 添加消息持久化
- 支持二进制消息传输
- 集成消息队列系统 