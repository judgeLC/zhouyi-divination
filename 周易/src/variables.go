// variables.go 定义了系统运行时使用的全局变量和缓存
// 包括文本渲染缓存、图片缓存、随机数生成器、万年历缓存等
package main

import (
	"image"
	"math/rand"
	"sync"
	"time"
)

// textCacheMap 文本渲染缓存映射表
// 缓存已计算的文本宽度和渲染器，避免重复计算提高图片生成性能
// 键为文本内容，值为TextCache结构体，包含文本宽度、字体面等信息
var textCacheMap = make(map[string]*TextCache)

// cachedBackground 背景图片缓存
// 存储预加载的背景图片，避免每次生成卦象时重复读取磁盘文件
var cachedBackground *image.NRGBA

// 全局随机数生成器相关变量
// 使用单例模式确保整个程序使用同一个随机数生成器
var (
	globalRand     *rand.Rand // 全局随机数生成器实例
	globalRandOnce sync.Once  // 确保随机数生成器只初始化一次
)

// 背景图片缓存控制变量
// 使用sync.Once确保背景图片只加载一次，提高性能
var (
	cachedBackgroundOnce sync.Once // 确保背景图片只加载一次
)

// 字体文件缓存变量
// 字体文件较大，缓存可以显著提高图片生成速度
var (
	cachedFontBytes     []byte    // 缓存的字体文件二进制数据
	cachedFontBytesOnce sync.Once // 确保字体文件只加载一次
)

// 万年历API结果缓存
// 由于万年历数据按日期变化，使用日期作为键进行缓存
// 避免重复调用外部API，提高响应速度和降低网络开销
var (
	calendarCache      map[string]CalendarAPIResponse = make(map[string]CalendarAPIResponse) // 万年历数据缓存映射表
	calendarCacheMutex sync.RWMutex                                                          // 万年历缓存读写锁，保证并发安全
)

// getGlobalRand 获取全局随机数生成器
// 使用单例模式，确保整个程序使用同一个随机数生成器
// 基于当前纳秒时间戳作为种子，保证随机性
//
// 返回值：全局随机数生成器实例
func getGlobalRand() *rand.Rand {
	globalRandOnce.Do(func() {
		// 使用当前纳秒时间戳作为随机种子，确保每次程序启动时的随机性
		globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	return globalRand
}
