# 周易占卜系统 (Zhouyi Divination System)

一个基于Go语言开发的传统周易六爻占卜系统，提供RESTful API接口和卦象图片生成功能。系统结合传统易学理论与现代Web技术，为用户提供准确、美观的占卜服务。

## ✨ 功能特性

- 🎯 **六爻占卜**: 基于传统周易理论的六爻占卜算法
- 🖼️ **卦象可视化**: 生成包含完整信息的卦象图片
- 🌐 **RESTful API**: 提供标准HTTP接口服务
- 📅 **万年历集成**: 自动获取干支纪年信息
- ⚙️ **配置管理**: 支持自定义端口和API配置
- 🔮 **纳甲起卦**: 支持传统纳甲起卦方法
- 📊 **六神配置**: 完整的六神、六亲、五行配置

## 🚀 快速开始

### 环境要求
- **操作系统**: Windows 7/8/10/11
- **网络连接**: 需要访问万年历API
- **权限**: 文件读写权限

### 1. 下载和运行
```bash
# 下载项目
git clone https://github.com/judgeLC/zhouyi-divination.git
cd zhouyi-divination/周易/src

# 直接运行
./Yijing.exe
```

### 2. 配置系统
首次运行会自动生成 `config.json` 配置文件：
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

### 3. 测试API接口
```bash
curl -X POST http://localhost:8090/api/divine \
  -H "Content-Type: application/json" \
  -d '{"method": "today", "params": {}}'
```

## 📡 API 接口文档

### 占卜接口
- **URL**: `POST /api/divine`
- **Content-Type**: `application/json`

#### 请求示例
```json
{
    "method": "today",
    "params": {}
}
```

#### 响应示例
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

#### 图片访问
```
http://localhost:8090/photos/卜卦_20231231154000.png
```

**详细API文档**: 请查看 [API接口文档.md](API接口文档.md)

## 📁 项目结构

```
zhouyi-divination/
├── 周易/                          # 主项目目录
│   └── src/                      # 源代码目录
│       ├── Yijing.exe           # 可执行文件 (带太极图标)
│       ├── config.json          # 配置文件
│       ├── 配置说明.md            # 配置使用说明
│       ├── 万年历API说明.txt       # API专门说明
│       │
│       ├── main.go              # 程序入口
│       ├── server.go            # HTTP服务器
│       ├── config.go            # 配置管理
│       ├── calendar_api.go      # 万年历API
│       ├── divine_generator.go  # 占卜生成器
│       ├── image_generator.go   # 图片生成器
│       ├── gua_logic.go         # 卦象逻辑
│       ├── gua_data.go          # 卦象数据
│       ├── najia.go             # 纳甲功能
│       ├── constants.go         # 常量定义
│       ├── types.go             # 数据类型
│       ├── utils.go             # 工具函数
│       ├── variables.go         # 变量定义
│       │
│       ├── go.mod               # Go模块文件
│       ├── go.sum               # 依赖校验
│       │
│       ├── images/              # 资源图片
│       │   └── background.png   # 背景图片
│       ├── ttf/                 # 字体文件
│       │   └── simkai.ttf       # 楷体字体
│       ├── cache/               # 缓存目录 (运行时)
│       ├── log/                 # 日志目录 (运行时)
│       ├── output/              # 输出目录 (运行时)
│       └── photos/              # 卦象图片 (运行时)
│
├── API接口文档.md                  # 详细API文档
├── README.md                     # 项目说明
└── .gitignore                   # Git忽略配置
```

## ⚙️ 配置说明

### 自定义端口
编辑 `config.json` 修改服务端口：
```json
{
    "server": {
        "port": "9000"
    }
}
```

### 万年历API配置
⚠️ **重要**: 默认使用公共测试密钥，有频次限制。建议获取个人密钥：

1. 访问 [接口盒子官网](https://www.apihz.cn) 注册账号
2. 获取个人ID和KEY
3. 修改配置文件中的 `id` 和 `key` 字段

详细配置说明请查看：`周易/src/配置说明.md`

## 🎨 卦象图片说明

### 图片特点
- **格式**: PNG (1200x900像素)
- **内容**: 完整六爻卦象信息
  - 主卦和变卦
  - 六爻详细信息 (爻位、五行、六亲、六神)
  - 干支纪年信息
  - 占卜时间

### 图片示例
生成的卦象图片包含：
- 🎯 卦名和卦象符号
- 📊 六爻排盘信息
- 🔮 五行生克关系
- 📅 干支日期信息

## 🛠️ 开发信息

### 技术栈
- **语言**: Go 1.19+
- **架构**: RESTful API
- **图片处理**: Go标准库
- **配置管理**: JSON格式
- **外部依赖**: 万年历API

### 编译说明
```bash
# 编译命令
go build -o Yijing.exe .

# 带图标编译
rsrc -arch amd64 -ico taiji_icon.ico
go build -o Yijing.exe .
```

## 📞 技术支持

### 常见问题
1. **端口被占用**: 修改 `config.json` 中的端口配置
2. **API频次限制**: 使用个人万年历API密钥
3. **图片生成失败**: 检查字体文件和权限设置

### 日志查看
- 应用日志: `log/app_YYYYMMDD.log`
- 错误日志: `log/error_YYYYMMDD.log`

### 配置文档
- 📖 [配置说明.md](周易/src/配置说明.md) - 详细配置指南
- 📋 [万年历API说明.txt](周易/src/万年历API说明.txt) - API使用说明
- 📡 [API接口文档.md](API接口文档.md) - 完整接口文档

## 📄 许可证

本项目仅供学习和研究使用。

## 🔗 相关链接

- [万年历API服务](https://www.apihz.cn) - 数据源提供商
- [项目仓库](https://github.com/judgeLC/zhouyi-divination) - GitHub仓库 