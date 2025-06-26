// najia.go 实现传统六爻占卜中的纳甲理论
// 纳甲是将天干地支配置到卦象各爻位上的方法，用于确定每个爻的五行属性
// 这是六爻占卜中确定六亲关系的重要理论基础
package main

import "fmt"

// naJia 纳甲主函数，获取指定爻位的完整纳甲信息
// 将卦宫、爻位和卦象信息结合，生成该爻的爻位名称和干支五行信息
//
// 纳甲理论说明：
// 纳甲是将60甲子（天干地支组合）分配到八卦的各个爻位上
// 不同卦宫有不同的纳甲规律，每个爻位对应特定的干支和五行
// 这些信息用于计算六亲关系和判断吉凶
//
// 参数：
//   - guaGong: 卦宫名称，如"乾宫"、"坤宫"等
//   - yaoWei: 爻位序号（1-6，从下到上）
//   - 本卦: 卦象数组，用于确定该爻的阴阳性质
//
// 返回值：
//   - string: 爻位名称，如"初六"、"九二"等
//   - string: 干支五行信息，格式为"干支 (五行)"，如"甲子 (水)"
func naJia(guaGong string, yaoWei int, 本卦 []int) (string, string) {
	// 获取该爻位对应的干支和五行
	干支, 五行 := naJia干支五行(guaGong, yaoWei)

	// 根据爻位和阴阳性质生成标准爻位名称
	爻位名称 := getYaoWeiName(yaoWei, 本卦[yaoWei-1])

	// 返回爻位名称和格式化的干支五行信息
	return 爻位名称, fmt.Sprintf("%s (%s)", 干支, 五行)
}

// naJia干支五行 根据卦宫和爻位确定对应的干支和五行属性
// 这是纳甲理论的核心查询函数，每个卦宫都有固定的纳甲规律
//
// 参数：
//   - guaGong: 卦宫名称
//   - yaoWei: 爻位序号（1-6）
//
// 返回值：
//   - string: 对应的干支，如"甲子"、"乙丑"等
//   - string: 对应的五行，如"水"、"火"、"木"、"金"、"土"
func naJia干支五行(guaGong string, yaoWei int) (string, string) {
	switch guaGong {
	case "乾宫": // 乾宫纳甲：内卦配甲，外卦配壬
		switch yaoWei {
		case 1: // 初爻
			return "甲子", "水" // 甲配子，子属水
		case 2: // 二爻
			return "甲寅", "木" // 甲配寅，寅属木
		case 3: // 三爻
			return "甲辰", "土" // 甲配辰，辰属土
		case 4: // 四爻
			return "壬午", "火" // 壬配午，午属火
		case 5: // 五爻
			return "壬申", "金" // 壬配申，申属金
		case 6: // 上爻
			return "壬戌", "土" // 壬配戌，戌属土
		default:
			return "未知", "未知"
		}
	case "坤宫": // 坤宫纳甲：内卦配乙，外卦配癸
		switch yaoWei {
		case 1: // 初爻
			return "乙未", "土" // 乙配未，未属土
		case 2: // 二爻
			return "乙巳", "火" // 乙配巳，巳属火
		case 3: // 三爻
			return "乙卯", "木" // 乙配卯，卯属木
		case 4: // 四爻
			return "癸丑", "土" // 癸配丑，丑属土
		case 5: // 五爻
			return "癸亥", "水" // 癸配亥，亥属水
		case 6: // 上爻
			return "癸酉", "金" // 癸配酉，酉属金
		default:
			return "未知", "未知"
		}
	case "震宫": // 震宫纳甲：内外卦均配庚
		switch yaoWei {
		case 1: // 初爻
			return "庚子", "水" // 庚配子，子属水
		case 2: // 二爻
			return "庚寅", "木" // 庚配寅，寅属木
		case 3: // 三爻
			return "庚辰", "土" // 庚配辰，辰属土
		case 4: // 四爻
			return "庚午", "火" // 庚配午，午属火
		case 5: // 五爻
			return "庚申", "金" // 庚配申，申属金
		case 6: // 上爻
			return "庚戌", "土" // 庚配戌，戌属土
		default:
			return "未知", "未知"
		}
	case "巽宫": // 巽宫纳甲：内外卦均配辛
		switch yaoWei {
		case 1: // 初爻
			return "辛丑", "土" // 辛配丑，丑属土
		case 2: // 二爻
			return "辛亥", "水" // 辛配亥，亥属水
		case 3: // 三爻
			return "辛酉", "金" // 辛配酉，酉属金
		case 4: // 四爻
			return "辛卯", "木" // 辛配卯，卯属木
		case 5: // 五爻
			return "辛巳", "火" // 辛配巳，巳属火
		case 6: // 上爻
			return "辛未", "土" // 辛配未，未属土
		default:
			return "未知", "未知"
		}
	case "坎宫": // 坎宫纳甲：内外卦均配戊
		switch yaoWei {
		case 1: // 初爻
			return "戊寅", "木" // 戊配寅，寅属木
		case 2: // 二爻
			return "戊辰", "土" // 戊配辰，辰属土
		case 3: // 三爻
			return "戊午", "火" // 戊配午，午属火
		case 4: // 四爻
			return "戊申", "金" // 戊配申，申属金
		case 5: // 五爻
			return "戊戌", "土" // 戊配戌，戌属土
		case 6: // 上爻
			return "戊子", "水" // 戊配子，子属水
		default:
			return "未知", "未知"
		}
	case "离宫": // 离宫纳甲：内外卦均配己
		switch yaoWei {
		case 1: // 初爻
			return "己卯", "木" // 己配卯，卯属木
		case 2: // 二爻
			return "己巳", "火" // 己配巳，巳属火
		case 3: // 三爻
			return "己未", "土" // 己配未，未属土
		case 4: // 四爻
			return "己酉", "金" // 己配酉，酉属金
		case 5: // 五爻
			return "己亥", "水" // 己配亥，亥属水
		case 6: // 上爻
			return "己丑", "土" // 己配丑，丑属土
		default:
			return "未知", "未知"
		}
	case "艮宫": // 艮宫纳甲：内外卦均配丙
		switch yaoWei {
		case 1: // 初爻
			return "丙辰", "土" // 丙配辰，辰属土
		case 2: // 二爻
			return "丙午", "火" // 丙配午，午属火
		case 3: // 三爻
			return "丙申", "金" // 丙配申，申属金
		case 4: // 四爻
			return "丙戌", "土" // 丙配戌，戌属土
		case 5: // 五爻
			return "丙子", "水" // 丙配子，子属水
		case 6: // 上爻
			return "丙寅", "木" // 丙配寅，寅属木
		default:
			return "未知", "未知"
		}
	case "兑宫": // 兑宫纳甲：内外卦均配丁
		switch yaoWei {
		case 1: // 初爻
			return "丁巳", "火" // 丁配巳，巳属火
		case 2: // 二爻
			return "丁未", "土" // 丁配未，未属土
		case 3: // 三爻
			return "丁酉", "金" // 丁配酉，酉属金
		case 4: // 四爻
			return "丁亥", "水" // 丁配亥，亥属水
		case 5: // 五爻
			return "丁丑", "土" // 丁配丑，丑属土
		case 6: // 上爻
			return "丁卯", "木" // 丁配卯，卯属木
		default:
			return "未知", "未知"
		}
	default: // 未知卦宫的默认处理
		return "未知", "未知"
	}
}
