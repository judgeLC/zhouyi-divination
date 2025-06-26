// types.go 定义了整个周易占卜系统中使用的所有数据结构和常量
// 包括卦象数据、API接口、WebSocket通信、图像渲染等相关类型定义
package main

import (
	"golang.org/x/image/font"
)

// Gua 卦象数据结构体
// 存储单个卦象的完整信息，包括卦名、构成、归宫和爻辞等
type Gua struct {
	FullName string   // 卦的完整名称，如"乾为天"、"坤为地"等
	ShangGua string   // 上卦名称（上三爻组成的卦）
	XiaGua   string   // 下卦名称（下三爻组成的卦）
	GuaGong  string   // 所属卦宫，用于确定纳甲和六亲关系
	YaoCi    []string // 爻辞数组，从初爻到上爻的六条爻辞
}

// TextCache 文本渲染缓存结构体
// 用于缓存已渲染的文本，避免重复计算文本宽度和创建渲染器
type TextCache struct {
	Text   string       // 缓存的文本内容
	Width  int          // 文本渲染宽度（像素）
	Face   font.Face    // 使用的字体面
	Drawer *font.Drawer // 文本绘制器
}

// DivineResult 占卜结果数据结构体
// 存储一次完整占卜的所有信息，用于API响应和数据存储
type DivineResult struct {
	ID          string `json:"id"`          // 占卜结果的唯一标识符
	Date        string `json:"date"`        // 占卜日期，格式：YYYY-MM-DD
	Ganzhinian  string `json:"ganzhinian"`  // 干支纪年，如"甲辰年"
	Ganzhiyue   string `json:"ganzhiyue"`   // 干支纪月，如"丙寅月"
	Ganzhiri    string `json:"ganzhiri"`    // 干支纪日，如"乙巳日"
	BenGua      string `json:"bengua"`      // 本卦名称
	BenGuaDesc  string `json:"benguadesc"`  // 本卦完整描述
	BianGua     string `json:"biangua"`     // 变卦名称（如果有动爻）
	BianGuaDesc string `json:"bianguadesc"` // 变卦完整描述（如果有动爻）
	HasDongYao  bool   `json:"hasdonyao"`   // 是否存在动爻（变爻）
	ImagePath   string `json:"imagepath"`   // 生成的卦象图片完整URL路径
	CreatedAt   int64  `json:"created_at"`  // 创建时间戳（Unix时间戳）
}

// DivineRequest 占卜请求参数结构体
// 客户端发送占卜请求时使用的参数格式
type DivineRequest struct {
	Type string `json:"type"` // 占卜类型，目前支持"today"（今日卦象）等
}

// ApiResponse 统一API响应格式结构体
// 所有HTTP API接口都使用此格式返回数据，确保响应格式的统一性
type ApiResponse struct {
	Code    int         `json:"code"`    // 响应状态码，200表示成功
	Message string      `json:"message"` // 响应消息，成功时为"成功"，失败时为错误描述
	Data    interface{} `json:"data"`    // 响应数据，具体内容根据接口而定
}

// CalendarAPIResponse 万年历API响应结构体
// 外部万年历API返回的数据格式，用于获取干支纪年月日信息
type CalendarAPIResponse struct {
	Code       int    `json:"code"`       // API状态码，200表示成功
	Ganzhinian string `json:"ganzhinian"` // 干支纪年，如"甲辰年"
	Ganzhiyue  string `json:"ganzhiyue"`  // 干支纪月，如"丙寅月"
	Ganzhiri   string `json:"ganzhiri"`   // 干支纪日，如"乙巳日"
}

// Layout 卦象图片布局参数结构体
// 定义卦象图片各个元素的位置和尺寸参数
type Layout struct {
	六神X    int // 六神文字的X坐标位置
	左卦中心X  int // 左侧卦象（本卦）中心X坐标
	右卦中心X  int // 右侧卦象（变卦）中心X坐标，-1表示不显示
	基础Y    int // 爻的基础Y坐标位置
	爻间距    int // 相邻两爻之间的垂直间距
	爻高度    int // 单个爻的高度
	爻宽度    int // 单个爻的宽度
	文字基线偏移 int // 文字相对于爻中心的基线偏移量
	爻辞Y    int // 爻辞文字的Y坐标位置
}

// WSMessage WebSocket消息结构体
// 定义WebSocket通信中使用的消息格式
type WSMessage struct {
	Type string      `json:"type"` // 消息类型，如"connect"、"divine"等
	Data interface{} `json:"data"` // 消息数据，根据消息类型而定
}

// WSConnection WebSocket连接信息结构体
// 存储单个WebSocket连接的基本信息
type WSConnection struct {
	ID   string `json:"id"`             // 连接的唯一标识符
	Name string `json:"name,omitempty"` // 连接名称（可选）
}

// WebSocket事件类型常量定义
// 定义WebSocket通信中使用的各种事件类型
const (
	WSEventConnect    = "connect"    // 客户端连接事件
	WSEventDisconnect = "disconnect" // 客户端断开连接事件
	WSEventDivine     = "divine"     // 占卜结果推送事件
	WSEventHeartbeat  = "heartbeat"  // 心跳检测事件
	WSEventError      = "error"      // 错误事件
)
