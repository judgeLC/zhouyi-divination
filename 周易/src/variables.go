package main

import (
	"image"
	"math/rand"
	"sync"
	"time"
)

// 文本缓存映射
var textCacheMap = make(map[string]*TextCache)

// 背景图片缓存
var cachedBackground *image.NRGBA

// 全局随机数生成器相关变量
var (
	globalRand     *rand.Rand
	globalRandOnce sync.Once
)

// 背景图片缓存相关变量
var (
	cachedBackgroundOnce sync.Once
)

// 字体缓存相关变量
var (
	cachedFontBytes     []byte
	cachedFontBytesOnce sync.Once
)

// 日历信息缓存
var (
	calendarCache      map[string]CalendarAPIResponse = make(map[string]CalendarAPIResponse)
	calendarCacheMutex sync.RWMutex
)

// 优化随机数生成器获取
func getGlobalRand() *rand.Rand {
	globalRandOnce.Do(func() {
		globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return globalRand
}
