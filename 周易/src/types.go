package main

import (
	"golang.org/x/image/font"
)

// 卦象结构体
type Gua struct {
	FullName string   // 卦的全称，如 乾为天
	ShangGua string   // 上卦名
	XiaGua   string   // 下卦名
	GuaGong  string   // 卦宫
	YaoCi    []string // 爻辞
}

// 文本缓存结构
type TextCache struct {
	Text   string
	Width  int
	Face   font.Face
	Drawer *font.Drawer
}

// 算卦结果结构
type DivineResult struct {
	ID          string `json:"id"`          // 唯一ID
	Date        string `json:"date"`        // 日期
	Ganzhinian  string `json:"ganzhinian"`  // 干支年
	Ganzhiyue   string `json:"ganzhiyue"`   // 干支月
	Ganzhiri    string `json:"ganzhiri"`    // 干支日
	BenGua      string `json:"bengua"`      // 本卦名
	BenGuaDesc  string `json:"benguadesc"`  // 本卦描述
	BianGua     string `json:"biangua"`     // 变卦名
	BianGuaDesc string `json:"bianguadesc"` // 变卦描述
	HasDongYao  bool   `json:"hasdonyao"`   // 是否有动爻
	ImagePath   string `json:"imagepath"`   // 图片路径
	CreatedAt   int64  `json:"created_at"`  // 创建时间
}

// 请求参数结构
type DivineRequest struct {
	Type string `json:"type"` // 注意：这里是小写 "type"
}

// 响应结构
type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 万年历API响应结构
type CalendarAPIResponse struct {
	Code       int    `json:"code"`
	Ganzhinian string `json:"ganzhinian"`
	Ganzhiyue  string `json:"ganzhiyue"`
	Ganzhiri   string `json:"ganzhiri"`
}

// 布局结构体
type Layout struct {
	六神X    int
	左卦中心X  int
	右卦中心X  int
	基础Y    int
	爻间距    int
	爻高度    int
	爻宽度    int
	文字基线偏移 int
	爻辞Y    int
}
