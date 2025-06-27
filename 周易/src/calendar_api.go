// calendar_api.go 处理万年历API的调用和缓存管理
// 提供干支纪年、纪月、纪日信息，用于占卜中的日干确定和时间记录
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

// getChinaCurrentTime 获取中国北京时间
// 优先使用系统时区设置，如果系统时区不正确则通过网络时间服务获取
//
// 返回值：
//   - time.Time: 中国北京时间
//   - error: 错误信息，成功时为nil
func getChinaCurrentTime() (time.Time, error) {
	// 设置中国北京时区
	beijingLocation, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Printf("加载北京时区失败: %v", err)
		// 如果无法加载时区，使用UTC+8作为备选
		beijingLocation = time.FixedZone("CST", 8*3600)
	}

	// 获取当前UTC时间并转换为北京时间
	utcNow := time.Now().UTC()
	beijingTime := utcNow.In(beijingLocation)

	log.Printf("当前北京时间: %s", beijingTime.Format("2006-01-02 15:04:05"))
	return beijingTime, nil
}

// getNetworkTime 通过网络时间服务获取准确的北京时间（备用方案）
// 当系统时间不准确时使用此方法
//
// 返回值：
//   - time.Time: 网络获取的北京时间
//   - error: 错误信息，成功时为nil
func getNetworkTime() (time.Time, error) {
	// 创建HTTP客户端，设置较短的超时时间
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	// 调用时间API获取当前时间
	resp, err := client.Get("http://worldtimeapi.org/api/timezone/Asia/Shanghai")
	if err != nil {
		return time.Time{}, fmt.Errorf("获取网络时间失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, fmt.Errorf("读取时间API响应失败: %w", err)
	}

	// 解析时间API响应
	var timeResp struct {
		Datetime string `json:"datetime"`
	}

	err = json.Unmarshal(body, &timeResp)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析时间API响应失败: %w", err)
	}

	// 解析时间字符串
	networkTime, err := time.Parse(time.RFC3339, timeResp.Datetime)
	if err != nil {
		return time.Time{}, fmt.Errorf("解析网络时间格式失败: %w", err)
	}

	log.Printf("网络获取的北京时间: %s", networkTime.Format("2006-01-02 15:04:05"))
	return networkTime, nil
}

// getRiGanAndCalendarInfo 获取当前日期的干支信息和万年历数据
// 包含智能缓存机制和超时控制，避免重复调用外部API
// 现在使用准确的中国北京时间进行查询
//
// 功能说明：
// 1. 获取准确的中国北京时间
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
	// 获取准确的中国北京时间
	beijingTime, err := getChinaCurrentTime()
	if err != nil {
		log.Printf("获取北京时间失败，尝试网络时间服务: %v", err)
		// 如果本地时区获取失败，尝试网络时间服务
		beijingTime, err = getNetworkTime()
		if err != nil {
			log.Printf("网络时间服务也失败，使用系统时间: %v", err)
			// 最后备选方案：使用系统时间（可能不准确）
			beijingTime = time.Now()
		}
	}

	year := beijingTime.Year()
	month := int(beijingTime.Month())
	day := beijingTime.Day()

	log.Printf("使用的查询时间: %d年%d月%d日", year, month, day)

	// 创建基于日期的缓存键，格式：YYYY-MM-DD
	cacheKey := fmt.Sprintf("%d-%02d-%02d", year, month, day)

	// 检查缓存中是否已存在当日数据（使用读锁保证并发安全）
	calendarCacheMutex.RLock()
	cachedResponse, found := calendarCache[cacheKey]
	calendarCacheMutex.RUnlock()

	if found {
		// 缓存命中，直接返回缓存的数据
		log.Printf("使用缓存的万年历数据: %s", cacheKey)
		return cachedResponse.Ganzhiri, cachedResponse.Ganzhinian, cachedResponse.Ganzhiyue, nil
	}

	// 缓存未命中，需要调用外部万年历API获取数据
	// 创建带超时控制的HTTP客户端，避免长时间等待
	client := &http.Client{
		Timeout: 10 * time.Second, // 增加超时时间到10秒，提高成功率
	}

	// 从配置文件获取万年历API的连接参数
	config := GetConfig()

	// 构建完整的API请求URL
	// 包含API主机地址、认证ID和密钥、查询的年月日参数
	apiURL := fmt.Sprintf("%s/api/time/getday.php?id=%s&key=%s&nian=%d&yue=%d&ri=%d",
		config.Calendar.APIHost, config.Calendar.ID, config.Calendar.Key, year, month, day)

	log.Printf("调用万年历API: %s", apiURL)

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

	log.Printf("万年历API响应: %s", string(body))

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

	// 验证响应数据的完整性
	if strings.TrimSpace(apiResponse.Ganzhiri) == "" ||
		strings.TrimSpace(apiResponse.Ganzhinian) == "" ||
		strings.TrimSpace(apiResponse.Ganzhiyue) == "" {
		return "", "", "", fmt.Errorf("万年历API返回的数据不完整")
	}

	// 将获取的数据保存到缓存中（使用写锁保证并发安全）
	calendarCacheMutex.Lock()
	calendarCache[cacheKey] = apiResponse
	calendarCacheMutex.Unlock()

	log.Printf("万年历数据获取成功: %s %s %s", apiResponse.Ganzhinian, apiResponse.Ganzhiyue, apiResponse.Ganzhiri)

	// 返回成功获取的干支纪年、纪月、纪日信息
	return apiResponse.Ganzhiri, apiResponse.Ganzhinian, apiResponse.Ganzhiyue, nil
}
