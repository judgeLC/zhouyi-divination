package main

import (
	"fmt"
	"log"
)

func init() {
	// 初始化卦象映射表
	for 卦名, 卦 := range guaXiang {
		key := 卦.ShangGua + "_" + 卦.XiaGua
		卦象映射表[key] = 卦名
	}
}

// 模拟摇一次铜钱，返回阴（0）或阳（1），以及是否为变爻
func yaoQian() (int, bool) {
	r := getGlobalRand() // 使用优化后的全局随机数生成器

	// 直接计算正面次数，避免循环
	正面次数 := r.Intn(4) // 0-3之间的随机数，直接模拟三次投掷的结果

	// 调试输出保持不变
	fmt.Printf("正面次数=%d ", 正面次数)

	switch 正面次数 {
	case 0: // 老阴，变为阳爻
		fmt.Printf("结果=0 (阴), 变爻=true\n")
		return 0, true
	case 3: // 老阳，变为阴爻
		fmt.Printf("结果=1 (阳), 变爻=true\n")
		return 1, true
	default:
		结果 := 正面次数 % 2
		fmt.Printf("结果=%d (%s), 变爻=false\n", 结果, map[int]string{0: "阴", 1: "阳"}[结果])
		return 结果, false // 少阴少阳，不变爻
	}
}

// 测试用卦象生成函数，确保显示阴爻和动爻
func generateTestGua() ([]int, []int, []bool) {
	本卦 := []int{1, 0, 1, 0, 1, 0}                          // 手动构造阴阳相间的卦象
	变卦 := []int{1, 1, 0, 0, 1, 0}                          // 第二爻发生变化，0->1
	变爻标记 := []bool{false, true, true, false, false, false} // 标记变化的爻
	return 本卦, 变卦, 变爻标记
}

// 生成卦象，优化版本
func generateGua() ([]int, []int, []bool) {
	本卦 := make([]int, 6)
	变卦 := make([]int, 6)
	变爻标记 := make([]bool, 6) // 记录每一爻是否为变爻

	// 一次性生成所有爻，避免多次设置种子
	for i := 0; i < 6; i++ {
		fmt.Printf("生成第%d爻: ", i+1)
		爻, 是否变爻 := yaoQian()
		本卦[i] = 爻
		变爻标记[i] = 是否变爻
		// 计算变卦
		if 是否变爻 {
			变卦[i] = 1 - 爻 // 如果是变爻，阴变阳，阳变阴
		} else {
			变卦[i] = 爻 // 如果不是变爻，保持不变
		}
	}

	return 本卦, 变卦, 变爻标记
}

// 将卦象转换为名称 - 优化版本，避免临时切片创建
func guaToName(gua []int) string {
	var shangGuaName, xiaGuaName string

	// 使用预定义的上下卦爻索引来避免创建临时切片
	xiaGuaYao := [3]int{gua[0], gua[1], gua[2]}
	shangGuaYao := [3]int{gua[3], gua[4], gua[5]}

	// 识别上下卦
	for name, yaoPattern := range 卦爻映射 {
		if yaoPattern[0] == xiaGuaYao[0] && yaoPattern[1] == xiaGuaYao[1] && yaoPattern[2] == xiaGuaYao[2] {
			xiaGuaName = name
		}
		if yaoPattern[0] == shangGuaYao[0] && yaoPattern[1] == shangGuaYao[1] && yaoPattern[2] == shangGuaYao[2] {
			shangGuaName = name
		}
	}

	// 使用映射表直接查询卦名
	key := shangGuaName + "_" + xiaGuaName
	if guaName, ok := 卦象映射表[key]; ok {
		return guaName
	}

	// 如果找不到，返回未知卦
	log.Printf("警告: 未找到匹配的卦名 (上卦:%s, 下卦:%s)\n", shangGuaName, xiaGuaName)
	return "未知卦"
}

// 定世爻
func dingShiYao(guaGong string) int {
	return guaGongShiYao[guaGong]
}

// 定应爻
func dingYingYao(shiYao int) int {
	return 7 - shiYao
}

// 获取六亲
func getLiuQin(guaGong string, yaoWuXing string) string {
	guaWuXing := 卦宫五行[guaGong]
	return 六亲关系[guaWuXing][yaoWuXing]
}

// 安六神
func anLiuShen(日干 string, yaoWei int) string {
	起始位置 := 日干六神[日干]
	六神索引 := (起始位置 + yaoWei - 1) % 6 // 计算六神索引
	return 六神[六神索引]
}

// 获取特定卦的爻辞
func getYaoCi(guaName string, yaoIndex int) string {
	gua, exists := guaXiang[guaName]
	if !exists || yaoIndex < 0 || yaoIndex >= len(gua.YaoCi) {
		return "无爻辞"
	}
	return gua.YaoCi[yaoIndex]
}
