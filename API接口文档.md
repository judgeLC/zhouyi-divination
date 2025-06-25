# 周易占卜系统 API 接口文档

## 📋 概述

周易占卜系统提供RESTful API接口，支持通过HTTP POST请求生成卦象图片。系统基于Go语言开发，采用传统周易六爻占卜算法，结合万年历数据生成准确的占卜结果。

## 🌐 服务器配置

### 默认访问地址
- **协议**: HTTP
- **主机**: localhost
- **端口**: 8090 (可在config.json中配置)
- **基础URL**: `http://localhost:8090`

### 自定义端口配置
修改 `config.json` 文件中的端口设置：
```json
{
    "server": {
        "port": "9000"
    }
}
```

## 🔌 API 接口详情

### 1. 占卜卦象生成接口

#### 基本信息
- **接口路径**: `/api/divine`
- **请求方法**: `POST`
- **Content-Type**: `application/json`
- **响应格式**: `JSON`

#### 完整URL
```
POST http://localhost:8090/api/divine
```

## 📨 请求说明

### 请求头 (Headers)
```http
Content-Type: application/json
Accept: application/json
```

### 请求体 (Request Body)
```json
{
    "method": "today",
    "params": {}
}
```

#### 请求参数说明
| 参数名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| method | string | 是 | 占卜方法，固定值: "today" |
| params | object | 是 | 参数对象，当前为空对象 |

### 请求示例

#### cURL 示例
```bash
curl -X POST http://localhost:8090/api/divine \
  -H "Content-Type: application/json" \
  -d '{
    "method": "today",
    "params": {}
  }'
```

#### JavaScript 示例
```javascript
fetch('http://localhost:8090/api/divine', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    method: "today",
    params: {}
  })
})
.then(response => response.json())
.then(data => console.log(data));
```

#### Python 示例
```python
import requests
import json

url = "http://localhost:8090/api/divine"
headers = {
    "Content-Type": "application/json"
}
data = {
    "method": "today",
    "params": {}
}

response = requests.post(url, headers=headers, data=json.dumps(data))
result = response.json()
print(result)
```

#### PowerShell 示例
```powershell
$headers = @{
    "Content-Type" = "application/json"
}
$body = @{
    method = "today"
    params = @{}
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8090/api/divine" -Method POST -Headers $headers -Body $body
```

## 📤 响应说明

### 成功响应格式
```json
{
    "code": 200,
    "message": "成功",
    "data": {
        "id": "divine_1640995200000000000",
        "date": "2023-12-31",
        "image_path": "photos/卜卦_20231231154000.png",
        "created_at": 1640995200
    }
}
```

#### 响应字段说明
| 字段名 | 类型 | 说明 |
|--------|------|------|
| code | number | 状态码，200表示成功 |
| message | string | 响应消息 |
| data | object | 响应数据对象 |
| data.id | string | 占卜记录唯一标识 |
| data.date | string | 占卜日期 (YYYY-MM-DD格式) |
| data.image_path | string | 生成的卦象图片相对路径 |
| data.created_at | number | 创建时间戳 (Unix时间戳) |

### 图片访问
生成的卦象图片可通过以下URL访问：
```
http://localhost:8090/photos/卜卦_20231231154000.png
```

### 错误响应格式
```json
{
    "code": 400,
    "message": "请求参数错误",
    "data": null
}
```

#### 常见错误码
| 错误码 | 说明 | 解决方案 |
|--------|------|----------|
| 400 | 请求参数错误 | 检查请求体格式和参数 |
| 405 | 请求方法不允许 | 确保使用POST方法 |
| 500 | 服务器内部错误 | 检查服务器日志 |

## 🖼️ 卦象图片说明

### 图片特点
- **格式**: PNG
- **尺寸**: 1200x900 像素
- **内容**: 包含完整的六爻卦象信息
  - 卦名和卦象符号
  - 六爻详细信息 (爻位、五行、六亲、六神)
  - 日期干支信息
  - 占卜时间

### 图片存储位置
- **服务器路径**: `photos/` 目录
- **访问路径**: `/photos/` URL路径
- **命名规则**: `卜卦_YYYYMMDDHHMMSS.png`

## 🔧 系统配置

### 万年历API配置
系统依赖万年历API获取干支信息，配置位于 `config.json`：
```json
{
    "calendar": {
        "api_host": "https://cn.apihz.cn",
        "id": "88888888",
        "key": "88888888"
    }
}
```

⚠️ **注意**: 默认使用公共测试密钥，建议获取个人密钥以避免频次限制。

## 🚀 快速开始

### 1. 启动服务
```bash
# Windows
.\Yijing.exe

# 服务启动后会显示：
# 启动HTTP服务器，监听端口 8090...
# API接口路径: http://localhost:8090/api/divine
```

### 2. 测试接口
```bash
curl -X POST http://localhost:8090/api/divine \
  -H "Content-Type: application/json" \
  -d '{"method": "today", "params": {}}'
```

### 3. 查看生成的图片
在响应中获取 `image_path`，然后访问：
```
http://localhost:8090/photos/卜卦_YYYYMMDDHHMMSS.png
```

## 🛠️ 故障排除

### 常见问题

#### 1. 连接被拒绝
- **原因**: 服务未启动或端口被占用
- **解决**: 检查服务状态，确认端口配置

#### 2. 万年历API错误
- **原因**: 网络问题或API频次限制
- **解决**: 检查网络连接，考虑使用个人API密钥

#### 3. 图片生成失败
- **原因**: 字体文件缺失或权限问题
- **解决**: 检查 `ttf/` 目录和文件权限

### 日志查看
系统运行日志保存在 `log/` 目录：
- `app_YYYYMMDD.log` - 应用日志
- `error_YYYYMMDD.log` - 错误日志

## 📞 技术支持

### 系统要求
- **操作系统**: Windows 7/8/10/11
- **网络**: 需要访问万年历API
- **权限**: 需要文件读写权限

### 配置文件
- `config.json` - 主配置文件
- `配置说明.md` - 详细配置指南
- `万年历API说明.txt` - API专门说明

### 开发信息
- **语言**: Go 1.19+
- **架构**: RESTful API
- **依赖**: 标准库 + 万年历API 