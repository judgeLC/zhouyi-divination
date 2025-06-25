package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 完整卦象结构喵
type Hexagram struct {
	Name       string  // 卦名
	Palace     string  // 卦宫
	Lines      [6]Line // 六爻
	WorldLine  int     // 世爻位置
	DateGanZhi string  // 起卦日干支
}

type Line struct {
	YaoType  string // 爻象
	GanZhi   string // 干支
	Relation string // 六亲
	God      string // 六神
	IsMoving bool   // 是否动爻
}

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("叮铃铃~依凌开始起卦喵！(ฅ•ω•ฅ)♡")

	// 1. 生成六爻
	hex := generateHexagram()

	// 2. 显示完整卦象
	printHexagram(hex)

	fmt.Println("\n锵锵~卦象生成完毕喵！(举起爪爪)")
	fmt.Println("主人要请依凌吃三文鱼寿司才给解卦喵~ (⁄ ⁄•⁄ω⁄•⁄ ⁄)")
}

func generateHexagram() Hexagram {
	// 当前时间作为起卦时间喵
	now := time.Now()
	dateGan := getDayGan(now)

	h := Hexagram{
		DateGanZhi: getGanZhi(now),
	}

	// 生成六爻喵
	for i := 0; i < 6; i++ {
		coins := rand.Intn(2) + rand.Intn(2) + rand.Intn(2)
		h.Lines[i].IsMoving = coins == 0 || coins == 3

		switch coins {
		case 3:
			h.Lines[i].YaoType = "——  老阳"
		case 2:
			h.Lines[i].YaoType = "——  少阳"
		case 1:
			h.Lines[i].YaoType = "— — 少阴"
		case 0:
			h.Lines[i].YaoType = "— — 老阴"
		}
	}

	// 计算上下卦喵
	upper, lower := calculateTrigrams(h.Lines)
	h.Name = getHexagramName(upper, lower)

	// 定卦宫喵
	h.Palace = getPalace(upper, lower)

	// 定世爻喵
	h.WorldLine = getWorldLine(upper, lower)

	// 纳干支喵
	applyGanZhi(&h)

	// 定六亲喵
	applyRelations(&h)

	// 安六神喵
	applyGods(&h, dateGan)

	return h
}

func printHexagram(hex Hexagram) {
	fmt.Printf("\n【%s】%s宫 世在%d爻\n", hex.Name, hex.Palace, hex.WorldLine+1)
	fmt.Printf("起卦时间: %s %s\n", time.Now().Format("2006-01-02 15:04"), hex.DateGanZhi)

	for i := 5; i >= 0; i-- {
		line := hex.Lines[i]
		movingMark := ""
		if line.IsMoving {
			movingMark = " (动爻)"
		}
		if i == hex.WorldLine {
			movingMark += "(世)"
		}
		fmt.Printf("%s %s %s %s%s\n", line.YaoType, line.God, line.Relation, line.GanZhi, movingMark)
	}
}

// 其他辅助函数实现喵
func calculateTrigrams(lines [6]Line) (int, int) {
	// 实际实现卦象计算喵
	return 7, 7 // 示例值喵
}

func getHexagramName(upper, lower int) string {
	// 实际卦名映射喵
	return "乾为天"
}

func getPalace(upper, lower int) string {
	// 实际卦宫算法喵
	return "乾"
}

func getWorldLine(upper, lower int) int {
	// 实际世爻定位喵
	return 5
}

func applyGanZhi(hex *Hexagram) {
	// 实际纳甲算法喵
	hex.Lines[5].GanZhi = "壬戌"
	hex.Lines[4].GanZhi = "壬申"
	hex.Lines[3].GanZhi = "壬午"
	hex.Lines[2].GanZhi = "甲辰"
	hex.Lines[1].GanZhi = "甲寅"
	hex.Lines[0].GanZhi = "甲子"
}

func applyRelations(hex *Hexagram) {
	// 实际六亲算法喵
	relations := []string{"父母", "兄弟", "官鬼", "妻财", "官鬼", "子孙"}
	for i := range hex.Lines {
		hex.Lines[i].Relation = relations[i]
	}
}

func applyGods(hex *Hexagram, dayGan string) {
	// 六神排列规则喵
	godsOrder := map[string][]string{
		"甲": {"青龙", "朱雀", "勾陈", "螣蛇", "白虎", "玄武"},
		"乙": {"青龙", "朱雀", "勾陈", "螣蛇", "白虎", "玄武"},
		"丙": {"朱雀", "勾陈", "螣蛇", "白虎", "玄武", "青龙"},
		//...其他天干映射
	}

	if order, ok := godsOrder[dayGan]; ok {
		for i := range hex.Lines {
			hex.Lines[i].God = order[i]
		}
	}
}

func getDayGan(t time.Time) string {
	// 简化版日干计算喵
	gans := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	return gans[t.Day()%10]
}

func getGanZhi(t time.Time) string {
	// 简化版干支计算喵
	return "壬寅年 丙午月 戊戌日"
}
