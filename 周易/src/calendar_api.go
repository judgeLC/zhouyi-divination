// calendar_api.go 处理万年历API的调用和缓存管理
// 提供干支纪年、纪月、纪日信息，用于占卜中的日干确定和时间记录
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// getRiGanAndCalendarInfo 获取当前日期的干支信息和万年历数据
// 包含智能缓存机制和超时控制，避免重复调用外部API
//
// 功能说明：
// 1. 获取当前系统时间的年月日
// 2. 检查缓存中是否已有当日数据
// 3. 如无缓存则调用外部万年历API
// 4. 解析API响应并缓存结果
// 5. 返回干支纪年、纪月、纪日信息
//
// 返回值：
//   - string: 干支日，如"乙巳日"
//   - string: 干支年，如"甲辰年"
//   - string: 干支月，如"丙寅月"
//   - error: 错误信息，成功时为nil
func getRiGanAndCalendarInfo() (string, string, string, error) {
	// 获取当前系统时间
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()

	// 创建基于日期的缓存键，格式：YYYY-MM-DD
	cacheKey := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	// 检查缓存中是否已存在当日数据（使用读锁保证并发安全）
	calendarCacheMutex.RLock()
	cachedResponse, found := calendarCache[cacheKey]
	calendarCacheMutex.RUnlock()

	if found {
		// 缓存命中，直接返回缓存的数据
		return cachedResponse.Ganzhiri, cachedResponse.Ganzhinian, cachedResponse.Ganzhiyue, nil
	}

	// 缓存未命中，需要调用外部万年历API获取数据
	// 创建带超时控制的HTTP客户端，避免长时间等待
	client := &http.Client{
		Timeout: 5 * time.Second, // 5秒超时，平衡响应速度和可靠性
	}

	// 从配置文件获取万年历API的连接参数
	config := GetConfig()

	// 构建完整的API请求URL
	// 包含API主机地址、认证ID和密钥、查询的年月日参数
	apiURL := fmt.Sprintf("%s/api/time/getday.php?id=%s&key=%s&nian=%d&yue=%d&ri=%d",
		config.Calendar.APIHost, config.Calendar.ID, config.Calendar.Key, year, month, day)

	// 发送HTTP GET请求到万年历API
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", "", "", fmt.Errorf("调用万年历API失败: %w", err)
	}
	defer resp.Body.Close() // 确保响应体被正确关闭

	// 读取API响应的原始数据
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("读取API响应数据失败: %w", err)
	}

	// 将JSON格式的响应解析为结构体
	var apiResponse CalendarAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", "", "", fmt.Errorf("解析万年历API响应JSON失败: %w", err)
	}

	// 验证API响应的状态码
	if apiResponse.Code != 200 {
		return "", "", "", fmt.Errorf("万年历API返回错误状态码: %d", apiResponse.Code)
	}

	// 将获取的数据保存到缓存中（使用写锁保证并发安全）
	calendarCacheMutex.Lock()
	calendarCache[cacheKey] = apiResponse
	calendarCacheMutex.Unlock()

	// 返回成功获取的干支纪年、纪月、纪日信息
	return apiResponse.Ganzhiri, apiResponse.Ganzhinian, apiResponse.Ganzhiyue, nil
}
