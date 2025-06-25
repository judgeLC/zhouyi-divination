# 周易占卜项目 (Zhouyi Divination)

这是一个基于Go语言开发的周易占卜系统，提供卦象生成、图片生成和Web服务功能。

## 项目结构

```
goproject/
├── 周易/                    # 周易占卜主项目
│   ├── src/                # 主要源代码目录
│   │   ├── main.go         # 程序入口文件
│   │   ├── server.go       # Web服务器
│   │   ├── gua_data.go     # 卦象数据定义
│   │   ├── gua_logic.go    # 卦象逻辑处理
│   │   ├── divine_generator.go  # 占卜生成器
│   │   ├── image_generator.go   # 图片生成器
│   │   ├── calendar_api.go # 日历API接口
│   │   ├── najia.go        # 纳甲相关功能
│   │   ├── types.go        # 数据类型定义
│   │   ├── constants.go    # 常量定义
│   │   ├── variables.go    # 变量定义
│   │   ├── utils.go        # 工具函数
│   │   ├── go.mod          # Go模块依赖
│   │   ├── go.sum          # 依赖校验文件
│   │   ├── images/         # 资源图片目录
│   │   │   └── background.png
│   │   ├── ttf/            # 字体文件目录
│   │   │   └── simkai.ttf
│   │   ├── cache/          # 缓存目录 (运行时生成)
│   │   ├── log/            # 日志目录 (运行时生成)
│   │   ├── output/         # 输出目录 (运行时生成)
│   │   └── photos/         # 生成的卦图目录 (运行时生成)
│   └── ttf/                # 额外字体文件目录
│       └── simkai.ttf
└── .gitignore             # Git忽略文件配置
```

## 功能特性

- 🎯 **卦象生成**: 根据时间和用户输入生成六爻卦象
- 🖼️ **图片生成**: 将卦象结果生成为可视化图片
- 🌐 **Web服务**: 提供HTTP接口供前端调用
- 📅 **日历API**: 集成日历功能
- 🔮 **纳甲预测**: 支持纳甲起卦方法

## 快速开始

### 环境要求
- Go 1.19+
- Git

### 运行周易占卜系统
```bash
cd 周易/src
go mod tidy
go run main.go
```

## 注意事项

- `cache/`, `log/`, `output/`, `photos/` 目录在程序运行时自动创建
- 生成的图片文件保存在相应的photos目录中
- 所有编译文件(.exe)和临时文件已通过.gitignore排除

## 许可证

本项目仅供学习和研究使用。 