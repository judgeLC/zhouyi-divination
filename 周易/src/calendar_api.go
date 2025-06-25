package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// 获取日干和万年历信息，带缓存和超时控制
func getRiGanAndCalendarInfo() (string, string, string, error) {
	// 获取当前日期
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	day := now.Day()

	// 创建缓存键
	cacheKey := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	// 检查缓存
	calendarCacheMutex.RLock()
	cachedResponse, found := calendarCache[cacheKey]
	calendarCacheMutex.RUnlock()

	if found {
		return cachedResponse.Ganzhiri, cachedResponse.Ganzhinian, cachedResponse.Ganzhiyue, nil
	}

	// 使用带超时的HTTP客户端
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// 构建 API URL
	apiURL := fmt.Sprintf("http://101.35.2.25/api/time/getday.php?id=10005598&key=12345678910&nian=%d&yue=%d&ri=%d", year, month, day)

	// 调用 API 接口
	resp, err := client.Get(apiURL)
	if err != nil {
		return "", "", "", fmt.Errorf("调用 API 失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取 API 响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", fmt.Errorf("读取 API 响应失败: %w", err)
	}

	// 解析 JSON 数据
	var apiResponse CalendarAPIResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return "", "", "", fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 检查 API 状态码
	if apiResponse.Code != 200 {
		return "", "", "", fmt.Errorf("API 返回错误状态码: %d", apiResponse.Code)
	}

	// 保存到缓存
	calendarCacheMutex.Lock()
	calendarCache[cacheKey] = apiResponse
	calendarCacheMutex.Unlock()

	// 返回日干和万年历信息
	return apiResponse.Ganzhiri, apiResponse.Ganzhinian, apiResponse.Ganzhiyue, nil
}
