# 周易64卦占卜系统

## 项目介绍

本项目是一个基于传统周易理论的64卦占卜模拟系统，采用Go语言开发。系统严格按照传统六爻占卜理论实现，包含完整的卦象生成、纳甲配置、六亲关系、六神安排等功能。

项目主要功能包括：
- 传统六爻占卜算法实现
- 64卦卦象自动识别和解析
- 卦象图片自动生成（1200x900像素）
- 干支纪年月日信息获取
- 爻辞查询和显示
- 动爻变卦处理

系统提供HTTP API接口和WebSocket实时通信两种方式，支持多用户并发使用，适合学习传统易学理论或进行占卜娱乐。

## 技术说明

### 目录结构

```
goproject/
├── 周易/
│   └── src/
│       ├── main.go              # 程序入口，系统初始化
│       ├── config.go            # 配置管理，JSON配置文件处理
│       ├── server.go            # HTTP服务器，API路由处理
│       ├── types.go             # 数据结构定义
│       ├── constants.go         # 易学常量定义（天干地支、五行等）
│       ├── variables.go         # 全局变量和缓存管理
│       ├── gua_logic.go         # 卦象生成和识别逻辑
│       ├── gua_data.go          # 64卦数据和爻辞
│       ├── najia.go             # 纳甲理论实现
│       ├── divine_generator.go  # 占卜结果生成
│       ├── image_generator.go   # 卦象图片生成
│       ├── calendar_api.go      # 万年历API调用
│       ├── websocket.go         # WebSocket通信处理
│       ├── utils.go             # 工具函数
│       ├── config.json          # 配置文件
│       ├── go.mod               # Go模块依赖
│       ├── go.sum               # 依赖校验文件
│       ├── ttf/                 # 字体文件目录
│       ├── images/              # 背景图片目录
│       ├── photos/              # 生成的卦象图片目录
│       ├── output/              # 输出文件目录
│       └── log/                 # 日志文件目录
├── API接口文档.md               # API接口详细说明
├── 配置说明.md                 # 配置文件说明
└── README.md                   # 项目说明文档
```

### API接口说明

#### HTTP API接口

**占卜接口**
- **URL**: `POST /api/divine`
- **功能**: 生成今日卦象
- **请求参数**:
  ```json
  {
    "type": "today"
  }
  ```
- **响应格式**:
  ```json
  {
    "code": 200,
    "message": "成功",
    "data": {
      "id": "divine_1703123456789",
      "date": "2023-12-21",
      "ganzhinian": "甲辰年",
      "ganzhiyue": "丙寅月", 
      "ganzhiri": "乙巳日",
      "bengua": "乾",
      "benguadesc": "乾为天",
      "biangua": "坤",
      "bianguadesc": "坤为地",
      "hasdonyao": true,
      "imagepath": "http://localhost:8090/photos/卜卦_20231221123456.png",
      "created_at": 1703123456
    }
  }
  ```

**WebSocket状态查询**
- **URL**: `GET /api/ws/status`
- **功能**: 查询当前WebSocket连接状态
- **响应**: 返回当前活跃连接数等信息

**静态文件访问**
- **卦象图片**: `GET /photos/{filename}`
- **输出文件**: `GET /output/{filename}`

#### WebSocket接口

**连接地址**: `ws://localhost:8090/ws`

**消息类型**:
- `connect`: 客户端连接成功
- `disconnect`: 客户端断开连接  
- `divine`: 占卜结果推送
- `heartbeat`: 心跳检测
- `error`: 错误信息

**测试页面**: `http://localhost:8090/test`

### 万年历API说明

系统集成了外部万年历API来获取准确的干支纪年月日信息，这对于传统占卜的准确性至关重要。

**API提供商**: cn.apihz.cn  
**接口地址**: `https://cn.apihz.cn/api/time/getday.php`  
**请求方式**: GET  
**请求参数**:
- `id`: API访问ID
- `key`: API访问密钥  
- `nian`: 查询年份
- `yue`: 查询月份
- `ri`: 查询日期

**响应格式**:
```json
{
  "code": 200,
  "ganzhinian": "甲辰年",
  "ganzhiyue": "丙寅月",
  "ganzhiri": "乙巳日"
}
```

**缓存机制**: 系统按日期缓存万年历数据，避免重复调用外部API，提高响应速度。

**容错处理**: 当万年历API不可用时，系统会使用默认的干支信息确保程序正常运行。

**使用说明**: 默认的API ID与KEY为公共ID与KEY，共享每分钟调用频次限制。接口本身免费，请使用自己的ID与KEY，独享每分钟调用频次。每日调用无上限。建议访问 https://www.apihz.cn 注册获取个人专属密钥。

### 技术架构

- **开发语言**: Go 1.21+
- **Web框架**: 标准库 net/http
- **WebSocket**: gorilla/websocket
- **图像处理**: disintegration/imaging
- **字体渲染**: golang.org/x/image/font
- **并发控制**: sync包，读写锁保证线程安全
- **缓存策略**: 内存缓存，提高性能
- **日志系统**: 标准库log，支持文件和控制台双输出

## 配置和部署

### 配置文件说明

系统使用JSON格式的配置文件 `config.json`，首次运行时会自动创建默认配置：

```json
{
    "server": {
        "port": "8090"
    },
    "calendar": {
        "api_host": "https://cn.apihz.cn",
        "id": "88888888",
        "key": "88888888"
    }
}
```

**配置项说明**:
- `server.port`: HTTP服务器监听端口
- `calendar.api_host`: 万年历API服务器地址
- `calendar.id`: 万年历API访问ID（需要注册获取）
- `calendar.key`: 万年历API访问密钥（需要注册获取）

### 部署步骤

#### 1. 环境准备
```bash
# 安装Go语言环境（1.21+）
# 下载项目源码
git clone <repository-url>
cd goproject/周易/src
```

#### 2. 依赖安装
```bash
# 下载依赖包
go mod tidy
```

#### 3. 配置设置
- 编辑 `config.json` 文件
- 配置万年历API的ID和密钥（可选，有默认值）
- 根据需要修改服务器端口

#### 4. 编译运行
```bash
# 编译生成可执行文件
go build -o Yijing.exe

# 运行程序
./Yijing.exe
```

#### 5. 验证部署
- 访问 `http://localhost:8090/test` 测试WebSocket功能
- 调用 `POST http://localhost:8090/api/divine` 测试占卜接口
- 检查 `log/` 目录下的日志文件

### 生产环境部署建议

1. **反向代理**: 使用Nginx等反向代理服务器
2. **进程管理**: 使用systemd或supervisor管理进程
3. **日志轮转**: 配置logrotate进行日志文件管理
4. **资源监控**: 监控CPU、内存使用情况
5. **备份策略**: 定期备份配置文件和重要数据

### 系统要求

- **操作系统**: Windows/Linux/macOS
- **内存**: 最低128MB，推荐256MB以上
- **磁盘空间**: 最低50MB（用于程序文件、字体、日志等）
- **网络**: 需要访问外部万年历API（可选）

### 故障排除

**常见问题**:
1. **端口占用**: 修改config.json中的端口配置
2. **字体缺失**: 确保ttf目录下有字体文件，或系统有中文字体
3. **万年历API失败**: 检查网络连接和API配置，系统会使用默认值继续运行
4. **图片生成失败**: 检查photos目录权限，确保程序有写入权限 