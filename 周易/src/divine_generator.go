package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"strings"
	"time"

	"golang.org/x/image/font"
)

// 新增：生成今日卦象并返回图像路径
func generateTodayGua() (string, error) {
	// 获取信号量，限制并发图片生成数量
	imageGenerationSem <- struct{}{}
	defer func() { <-imageGenerationSem }()

	// 加锁保证图片生成过程的串行化，防止并发冲突
	imageGenerationMutex.Lock()
	defer imageGenerationMutex.Unlock()

	log.Printf("开始生成卦象图片...")

	// 获取日干和万年历信息
	ganzhiri, ganzhinian, ganzhiyue, err := getRiGanAndCalendarInfo()
	if err != nil {
		log.Printf("获取日干和万年历信息失败: %v\n", err)
		// 使用默认值以便测试
		ganzhiri = "乙巳日"
		ganzhinian = "甲辰年"
		ganzhiyue = "王中月"
	}

	// 提取日干
	日干 := extractRiGan(ganzhiri)

	// 生成卦象
	本卦, 变卦, 变爻标记 := generateGua()

	// 检测是否有动爻
	有动爻 := hasChangingYao(变爻标记)

	// 根据卦象确定卦名
	本卦名 := guaToName(本卦)
	变卦名 := guaToName(变卦)

	// 初始化布局
	layout := initLayout(ImageWidth, 有动爻)

	// 获取背景
	dst := getBackground(ImageWidth, ImageHeight, "images/background.png")

	// 加载字体文件
	fontBytes, err := loadFontFile()
	if err != nil {
		return "", fmt.Errorf("加载字体文件失败: %v", err)
	}

	// 创建字体面
	titleFace, normalFace, smallFace, err := createFontFaces(fontBytes)
	if err != nil {
		return "", fmt.Errorf("创建字体失败: %v", err)
	}
	defer titleFace.Close()
	defer normalFace.Close()
	defer smallFace.Close()

	// 清空文本缓存
	textCacheMap = make(map[string]*TextCache)

	// 绘制图像内容
	err = drawGuaImage(dst, layout, 日干, 本卦名, 变卦名, 本卦, 变卦, 变爻标记, 有动爻, ganzhinian, ganzhiyue, ganzhiri, titleFace, normalFace, smallFace)
	if err != nil {
		return "", fmt.Errorf("绘制卦象图像失败: %v", err)
	}

	// 保存图像
	fileName := fmt.Sprintf("卜卦_%s.png", time.Now().Format("20060102150405"))
	savePath, err := saveImageToPathFixed(dst, fileName)
	if err != nil {
		return "", fmt.Errorf("保存图像失败: %v", err)
	}

	// 输出卦象信息到日志
	log.Printf("%s，%s，%s", ganzhinian, ganzhiyue, ganzhiri)
	log.Printf("本卦：%s %s", 本卦名, guaXiang[本卦名].FullName)
	if 有动爻 {
		log.Printf("变卦：%s %s", 变卦名, guaXiang[变卦名].FullName)
	} else {
		log.Printf("无动爻，无变卦")
	}

	log.Printf("卦象图片生成完成: %s", savePath)
	return savePath, nil
}

// 绘制卦象图像
func drawGuaImage(dst interface{}, layout *Layout, 日干, 本卦名, 变卦名 string, 本卦, 变卦 []int, 变爻标记 []bool, 有动爻 bool, ganzhinian, ganzhiyue, ganzhiri string, titleFace, normalFace, smallFace interface{}) error {
	img := dst.(*image.NRGBA)

	// 绘制标题（年月日）- 使用优化的居中文本绘制
	titleText := ganzhinian + " " + ganzhiyue + " " + ganzhiri
	drawCenteredText(img, titleText, ImageWidth/2, 70, titleFace.(font.Face))

	// 绘制本卦和变卦信息
	if 有动爻 {
		// 有动爻，显示双卦标题
		leftInfoX := layout.左卦中心X - 60
		rightInfoX := layout.右卦中心X - 45
		drawCachedText(img, guaXiang[本卦名].FullName, leftInfoX, 210, normalFace.(font.Face))
		drawCachedText(img, "("+guaXiang[本卦名].GuaGong+")", leftInfoX, 250, normalFace.(font.Face))
		drawCachedText(img, guaXiang[变卦名].FullName, rightInfoX, 210, normalFace.(font.Face))
		drawCachedText(img, "("+guaXiang[变卦名].GuaGong+")", rightInfoX, 250, normalFace.(font.Face))
	} else {
		// 无动爻，只显示单卦标题并居中
		titleX := layout.左卦中心X
		drawCenteredText(img, guaXiang[本卦名].FullName, titleX, 210, normalFace.(font.Face))
		drawCenteredText(img, "("+guaXiang[本卦名].GuaGong+")", titleX, 250, normalFace.(font.Face))
	}

	// 绘制卦象主体
	err := drawGuaBody(img, layout, 日干, 本卦名, 变卦名, 本卦, 变卦, 变爻标记, 有动爻, normalFace.(font.Face), smallFace.(font.Face))
	if err != nil {
		return err
	}

	// 绘制爻辞
	drawYaoCi(img, layout, 本卦名, 变卦名, 本卦, 变卦, 有动爻, smallFace.(font.Face))

	return nil
}

// 绘制卦象主体
func drawGuaBody(img *image.NRGBA, layout *Layout, 日干, 本卦名, 变卦名 string, 本卦, 变卦 []int, 变爻标记 []bool, 有动爻 bool, normalFace, smallFace font.Face) error {
	// 预先计算六神排序
	起始位置 := 日干六神[日干]
	六神排序 := make([]string, 6)
	for i := 0; i < 6; i++ {
		索引 := (起始位置 + i) % 6
		六神排序[i] = 六神[索引]
	}

	// 绘制基本布局
	yaoColor := color.RGBA{139, 69, 19, 255} // 棕色

	// 绘制六神和爻循环
	for i := 0; i < 6; i++ {
		rowY := layout.基础Y + i*layout.爻间距
		文字Y := rowY + layout.文字基线偏移 // 文字Y位置，基线对齐

		// 六神
		drawCachedText(img, 六神排序[i], layout.六神X, 文字Y, normalFace)

		// 本卦爻
		爻Y := rowY - layout.爻高度/2
		if 本卦[5-i] == 1 { // 阳爻
			drawYangYao(img, layout.左卦中心X-layout.爻宽度/2, 爻Y, layout.爻宽度, layout.爻高度, yaoColor)
		} else { // 阴爻
			drawYinYao(img, layout.左卦中心X-layout.爻宽度/2, 爻Y, layout.爻宽度, layout.爻高度, yaoColor)
		}

		// 本卦六亲信息
		_, 干支五行 := naJia(guaXiang[本卦名].GuaGong, 6-i, 本卦)
		五行部分 := strings.Split(干支五行, "(")[1]
		五行部分 = strings.TrimSuffix(五行部分, ")")
		六亲 := getLiuQin(guaXiang[本卦名].GuaGong, 五行部分)
		干支部分 := strings.Split(干支五行, " ")[0]
		drawCachedText(img, 六亲+干支部分+五行部分, layout.左卦中心X+layout.爻宽度/2+10, 文字Y, smallFace)

		// 只在有动爻情况下绘制变卦
		if 有动爻 {
			// 变卦爻
			if 变卦[5-i] == 1 { // 阳爻
				drawYangYao(img, layout.右卦中心X-layout.爻宽度/2, 爻Y, layout.爻宽度, layout.爻高度, yaoColor)
			} else { // 阴爻
				drawYinYao(img, layout.右卦中心X-layout.爻宽度/2, 爻Y, layout.爻宽度, layout.爻高度, yaoColor)
			}

			// 变卦六亲信息
			_, 干支五行 = naJia(guaXiang[变卦名].GuaGong, 6-i, 变卦)
			五行部分 = strings.Split(干支五行, "(")[1]
			五行部分 = strings.TrimSuffix(五行部分, ")")
			六亲 = getLiuQin(guaXiang[变卦名].GuaGong, 五行部分)
			干支部分 = strings.Split(干支五行, " ")[0]
			drawCachedText(img, 六亲+干支部分+五行部分, layout.右卦中心X+layout.爻宽度/2+10, 文字Y, smallFace)
		}

		// 动爻判定
		if 变爻标记[5-i] {
			动爻X := layout.左卦中心X + layout.爻宽度/2 + 150
			drawCachedText(img, "● 动爻", 动爻X, 文字Y, normalFace)
		}
	}

	// 绘制底部标签
	if 有动爻 {
		drawCachedText(img, "主卦", layout.左卦中心X-20, layout.基础Y+6*layout.爻间距+10, normalFace)
		drawCachedText(img, "变卦", layout.右卦中心X-20, layout.基础Y+6*layout.爻间距+10, normalFace)
	} else {
		drawCachedText(img, "主卦", layout.左卦中心X-20, layout.基础Y+6*layout.爻间距+10, normalFace)
	}

	return nil
}

// 绘制爻辞
func drawYaoCi(img *image.NRGBA, layout *Layout, 本卦名, 变卦名 string, 本卦, 变卦 []int, 有动爻 bool, smallFace font.Face) {
	if 有动爻 {
		// 双卦爻辞显示
		const (
			本卦爻辞X = 80
			变卦爻辞X = 620
			爻辞宽度  = 480
			行间距   = 25
		)

		nextY := layout.爻辞Y

		// 绘制爻辞标题
		drawCachedText(img, "爻辞：", 本卦爻辞X, nextY, smallFace)
		drawCachedText(img, "爻辞：", 变卦爻辞X, nextY, smallFace)
		nextY += 行间距

		// 双卦爻辞绘制
		for i := 0; i < 6; i++ {
			爻位 := 6 - i    // 从上爻到初爻
			爻辞索引 := 爻位 - 1 // 数组索引从0开始

			// 获取爻位名称
			爻位名称 := getYaoWeiName(爻位, 本卦[爻辞索引])

			// 获取对应的爻辞
			爻辞 := "无爻辞" // 默认值
			if gua, exists := guaXiang[本卦名]; exists && 爻辞索引 < len(gua.YaoCi) {
				爻辞 = gua.YaoCi[爻辞索引]
			}

			// 绘制本卦爻辞（带换行）
			完整爻辞 := fmt.Sprintf("%s: %s", 爻位名称, 爻辞)
			本爻辞Y := drawWrappedText(img, 完整爻辞, 本卦爻辞X, nextY, 爻辞宽度, 行间距, smallFace)

			// 获取变卦爻辞
			变爻位名称 := getYaoWeiName(爻位, 变卦[爻辞索引])
			变爻辞 := "无爻辞" // 默认值
			if gua, exists := guaXiang[变卦名]; exists && 爻辞索引 < len(gua.YaoCi) {
				变爻辞 = gua.YaoCi[爻辞索引]
			}

			// 绘制变卦爻辞（带换行）
			完整变爻辞 := fmt.Sprintf("%s: %s", 变爻位名称, 变爻辞)
			变爻辞Y := drawWrappedText(img, 完整变爻辞, 变卦爻辞X, nextY, 爻辞宽度, 行间距, smallFace)

			// 取两侧爻辞高度的较大值作为下一个爻辞的起始 Y 坐标
			nextY = max(本爻辞Y, 变爻辞Y) + 10 // 添加额外的10像素间距
		}
	} else {
		// 单卦爻辞显示（居中）
		const (
			爻辞宽度 = 600 // 单卦爻辞更宽一些
			行间距  = 25
		)

		本卦爻辞X := (ImageWidth - 爻辞宽度) / 2 // 居中显示爻辞
		nextY := layout.爻辞Y

		// 绘制爻辞标题
		drawCenteredText(img, "爻辞", ImageWidth/2, nextY, smallFace)
		nextY += 行间距

		// 单卦爻辞绘制
		for i := 0; i < 6; i++ {
			爻位 := 6 - i    // 从上爻到初爻
			爻辞索引 := 爻位 - 1 // 数组索引从0开始

			// 获取爻位名称
			爻位名称 := getYaoWeiName(爻位, 本卦[爻辞索引])

			// 获取对应的爻辞
			爻辞 := "无爻辞" // 默认值
			if gua, exists := guaXiang[本卦名]; exists && 爻辞索引 < len(gua.YaoCi) {
				爻辞 = gua.YaoCi[爻辞索引]
			}

			// 绘制本卦爻辞（带换行）
			完整爻辞 := fmt.Sprintf("%s: %s", 爻位名称, 爻辞)
			nextY = drawWrappedText(img, 完整爻辞, 本卦爻辞X, nextY, 爻辞宽度, 行间距, smallFace)
			nextY += 10 // 添加额外的间距
		}
	}
}
