package main

import "fmt"

// 纳甲
func naJia(guaGong string, yaoWei int, 本卦 []int) (string, string) {
	干支, 五行 := naJia干支五行(guaGong, yaoWei)
	爻位名称 := getYaoWeiName(yaoWei, 本卦[yaoWei-1])
	return 爻位名称, fmt.Sprintf("%s (%s)", 干支, 五行)
}

// 纳甲干支五行
func naJia干支五行(guaGong string, yaoWei int) (string, string) {
	switch guaGong {
	case "乾宫":
		switch yaoWei {
		case 1:
			return "甲子", "水"
		case 2:
			return "甲寅", "木"
		case 3:
			return "甲辰", "土"
		case 4:
			return "壬午", "火"
		case 5:
			return "壬申", "金"
		case 6:
			return "壬戌", "土"
		default:
			return "未知", "未知"
		}
	case "坤宫":
		switch yaoWei {
		case 1:
			return "乙未", "土"
		case 2:
			return "乙巳", "火"
		case 3:
			return "乙卯", "木"
		case 4:
			return "癸丑", "土"
		case 5:
			return "癸亥", "水"
		case 6:
			return "癸酉", "金"
		default:
			return "未知", "未知"
		}
	case "震宫":
		switch yaoWei {
		case 1:
			return "庚子", "水"
		case 2:
			return "庚寅", "木"
		case 3:
			return "庚辰", "土"
		case 4:
			return "庚午", "火"
		case 5:
			return "庚申", "金"
		case 6:
			return "庚戌", "土"
		default:
			return "未知", "未知"
		}
	case "巽宫":
		switch yaoWei {
		case 1:
			return "辛丑", "土"
		case 2:
			return "辛亥", "水"
		case 3:
			return "辛酉", "金"
		case 4:
			return "辛卯", "木"
		case 5:
			return "辛巳", "火"
		case 6:
			return "辛未", "土"
		default:
			return "未知", "未知"
		}
	case "坎宫":
		switch yaoWei {
		case 1:
			return "戊寅", "木"
		case 2:
			return "戊辰", "土"
		case 3:
			return "戊午", "火"
		case 4:
			return "戊申", "金"
		case 5:
			return "戊戌", "土"
		case 6:
			return "戊子", "水"
		default:
			return "未知", "未知"
		}
	case "离宫":
		switch yaoWei {
		case 1:
			return "己卯", "木"
		case 2:
			return "己巳", "火"
		case 3:
			return "己未", "土"
		case 4:
			return "己酉", "金"
		case 5:
			return "己亥", "水"
		case 6:
			return "己丑", "土"
		default:
			return "未知", "未知"
		}
	case "艮宫":
		switch yaoWei {
		case 1:
			return "丙辰", "土"
		case 2:
			return "丙午", "火"
		case 3:
			return "丙申", "金"
		case 4:
			return "丙戌", "土"
		case 5:
			return "丙子", "水"
		case 6:
			return "丙寅", "木"
		default:
			return "未知", "未知"
		}
	case "兑宫":
		switch yaoWei {
		case 1:
			return "丁巳", "火"
		case 2:
			return "丁未", "土"
		case 3:
			return "丁酉", "金"
		case 4:
			return "丁亥", "水"
		case 5:
			return "丁丑", "土"
		case 6:
			return "丁卯", "木"
		default:
			return "未知", "未知"
		}
	default:
		return "未知", "未知"
	}
}
