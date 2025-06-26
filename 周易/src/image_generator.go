package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// 生成渐变背景 - 优化版本，直接操作像素数据
func generateGradientBackground(width, height int) *image.NRGBA {
	// 直接创建NRGBA图像避免额外的内存分配
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	// 添加渐变效果
	for y := 0; y < height; y++ {
		factor := float64(y) / float64(height)
		r := uint8(248 - int(20*factor))
		g := uint8(240 - int(15*factor))
		b := uint8(224 - int(10*factor))
		// 优化:一次填充一整行
		for x := 0; x < width; x++ {
			pos := y*img.Stride + x*4
			img.Pix[pos] = r
			img.Pix[pos+1] = g
			img.Pix[pos+2] = b
			img.Pix[pos+3] = 255
		}
	}
	return img
}

// 加载或获取缓存的背景图片
func getBackground(width, height int, bgPath string) *image.NRGBA {
	cachedBackgroundOnce.Do(func() {
		// 打印当前状态信息
		log.Printf("正在加载背景图片，尺寸: %dx%d", width, height)
		// 创建默认背景（米色）
		bg := imaging.New(width, height, color.NRGBA{248, 240, 224, 255})

		// 可能的背景图路径列表
		paths := []string{
			bgPath,
			"background.png",
			"images/background.png",
			"../images/background.png",
			filepath.Join(getCurrentDir(), "images", "background.png"),
		}

		// 确保images目录存在
		imagesDir := filepath.Join(getCurrentDir(), "images")
		if err := ensureDir(imagesDir); err == nil {
			// 如果没有找到背景图，创建一个默认的渐变背景
			defaultBgPath := filepath.Join(imagesDir, "background.png")
			if !fileExists(defaultBgPath) {
				defaultBg := generateGradientBackground(width, height)
				err := imaging.Save(defaultBg, defaultBgPath)
				if err == nil {
					log.Printf("已创建默认背景图: %s", defaultBgPath)
					// 添加到搜索路径的开头
					paths = append([]string{defaultBgPath}, paths...)
				} else {
					log.Printf("创建默认背景失败: %v", err)
				}
			}
		}

		// 尝试加载背景图片
		var loadedBg bool
		for _, path := range paths {
			if fileExists(path) {
				log.Printf("尝试加载背景图: %s", path)
				bgImage, err := imaging.Open(path)
				if err == nil {
					// 调整大小以适应目标尺寸
					bgImage = imaging.Resize(bgImage, width, height, imaging.Lanczos)
					bg = imaging.Paste(bg, bgImage, image.Pt(0, 0))
					log.Printf("成功加载背景图片: %s", path)
					loadedBg = true
					break
				} else {
					log.Printf("加载背景图片失败: %v", err)
				}
			}
		}

		if !loadedBg {
			log.Printf("无法加载任何背景图片，使用生成的渐变背景")
			bg = generateGradientBackground(width, height)
		}

		cachedBackground = bg
	})

	// 返回背景副本，避免并发修改
	return imaging.Clone(cachedBackground)
}

// 加载字体文件，带缓存
func loadFontFile() ([]byte, error) {
	cachedFontBytesOnce.Do(func() {
		// 可能的字体路径
		paths := []string{
			"simkai.ttf",
			"ttf/simkai.ttf",
			"../ttf/simkai.ttf",
			"./ttf/simkai.ttf",
			filepath.Join(getCurrentDir(), "ttf", "simkai.ttf"),
			// Windows系统字体
			"C:\\Windows\\Fonts\\simkai.ttf", // 楷体
			"C:\\Windows\\Fonts\\simhei.ttf", // 黑体
			"C:\\Windows\\Fonts\\simsun.ttc", // 宋体
			// Linux系统字体
			"/usr/share/fonts/truetype/arphic/ukai.ttc",
			"/usr/share/fonts/truetype/wqy/wqy-microhei.ttc",
			// macOS系统字体
			"/Library/Fonts/Arial Unicode.ttf",
			"/System/Library/Fonts/PingFang.ttc",
		}

		// 确保ttf目录存在
		ttfDir := filepath.Join(getCurrentDir(), "ttf")
		ensureDir(ttfDir)

		log.Printf("正在寻找字体文件...")
		for _, path := range paths {
			if fileExists(path) {
				log.Printf("尝试加载字体: %s", path)
				data, err := os.ReadFile(path)
				if err == nil {
					log.Printf("成功加载字体: %s", path)
					// 如果加载的是系统字体，复制到项目目录
					if !strings.Contains(path, getCurrentDir()) {
						destPath := filepath.Join(ttfDir, "simkai.ttf")
						if err := os.WriteFile(destPath, data, 0644); err == nil {
							log.Printf("已将字体复制到: %s", destPath)
						}
					}
					cachedFontBytes = data
					return
				}
			}
		}
	})

	if cachedFontBytes == nil {
		return nil, fmt.Errorf("找不到可用的字体文件")
	}

	return cachedFontBytes, nil
}

// 预初始化布局常量
func initLayout(imageWidth int, 有动爻 bool) *Layout {
	var 左卦中心X int
	var 右卦中心X int
	var 六神X int

	if 有动爻 {
		// 有动爻，显示双卦
		左卦中心X = 300
		右卦中心X = 800
		六神X = 140
	} else {
		// 无动爻，只显示单卦并居中
		左卦中心X = imageWidth / 2 // 居中显示
		右卦中心X = -1             // 不显示右卦
		六神X = (imageWidth / 2) - 180
	}

	return &Layout{
		六神X:    六神X,
		左卦中心X:  左卦中心X,
		右卦中心X:  右卦中心X,
		基础Y:    300,
		爻间距:    40,
		爻高度:    20,
		爻宽度:    160,
		文字基线偏移: 10,
		爻辞Y:    300 + 6*50 + 20,
	}
}

// 优化后的绘制阳爻函数 - 一次性填充矩形区域
func drawYangYao(img *image.NRGBA, x, y, width, height int, color color.RGBA) {
	// 创建一行预填充的颜色数据
	row := make([]uint8, width*4) // RGBA 每像素4字节
	for i := 0; i < width; i++ {
		idx := i * 4
		row[idx] = color.R
		row[idx+1] = color.G
		row[idx+2] = color.B
		row[idx+3] = color.A
	}

	// 使用更高效的方式填充区域
	for j := 0; j < height; j++ {
		startPos := (y+j)*img.Stride + x*4
		endPos := startPos + width*4
		if startPos >= 0 && endPos <= len(img.Pix) {
			copy(img.Pix[startPos:endPos], row)
		}
	}
}

// 优化后的绘制阴爻函数 - 一次性填充矩形区域
func drawYinYao(img *image.NRGBA, x, y, width, height int, color color.RGBA) {
	gap := width / 3
	leftWidth := (width - gap) / 2
	rightWidth := width - gap - leftWidth

	// 创建左侧一行预填充的颜色数据
	leftRow := make([]uint8, leftWidth*4)
	for i := 0; i < leftWidth; i++ {
		idx := i * 4
		leftRow[idx] = color.R
		leftRow[idx+1] = color.G
		leftRow[idx+2] = color.B
		leftRow[idx+3] = color.A
	}

	// 创建右侧一行预填充的颜色数据
	rightRow := make([]uint8, rightWidth*4)
	for i := 0; i < rightWidth; i++ {
		idx := i * 4
		rightRow[idx] = color.R
		rightRow[idx+1] = color.G
		rightRow[idx+2] = color.B
		rightRow[idx+3] = color.A
	}

	// 直接复制到每一行，避免重复设置像素
	for j := 0; j < height; j++ {
		leftStartPos := (y+j)*img.Stride + x*4
		leftEndPos := leftStartPos + leftWidth*4
		rightStartPos := (y+j)*img.Stride + (x+leftWidth+gap)*4
		rightEndPos := rightStartPos + rightWidth*4

		// 确保在有效范围内，避免越界
		if leftStartPos >= 0 && leftEndPos <= len(img.Pix) {
			copy(img.Pix[leftStartPos:leftEndPos], leftRow)
		}
		if rightStartPos >= 0 && rightEndPos <= len(img.Pix) {
			copy(img.Pix[rightStartPos:rightEndPos], rightRow)
		}
	}
}

// 优化的文本绘制函数
func drawCachedText(img *image.NRGBA, text string, x, y int, face font.Face) {
	// 为每个独特的文本创建缓存键
	cacheKey := fmt.Sprintf("%s_%p", text, face)

	// 检查缓存
	cache, exists := textCacheMap[cacheKey]
	if !exists {
		// 如果缓存不存在，创建新的
		width := font.MeasureString(face, text).Round()
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}), // 黑色
			Face: face,
		}
		cache = &TextCache{
			Text:   text,
			Width:  width,
			Face:   face,
			Drawer: drawer,
		}
		// 存储到缓存
		textCacheMap[cacheKey] = cache
	} else {
		// 更新目标图像，以防图像已变更
		cache.Drawer.Dst = img
	}

	// 设置位置并绘制
	cache.Drawer.Dot = fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	cache.Drawer.DrawString(text)
}

// 优化版的居中绘制文本
func drawCenteredText(img *image.NRGBA, text string, centerX, y int, face font.Face) {
	cacheKey := fmt.Sprintf("%s_%p", text, face)
	cache, exists := textCacheMap[cacheKey]
	if !exists {
		width := font.MeasureString(face, text).Round()
		drawer := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(color.RGBA{0, 0, 0, 255}),
			Face: face,
		}
		cache = &TextCache{
			Text:   text,
			Width:  width,
			Face:   face,
			Drawer: drawer,
		}
		textCacheMap[cacheKey] = cache
	} else {
		// 更新目标图像
		cache.Drawer.Dst = img
	}

	x := centerX - cache.Width/2
	cache.Drawer.Dot = fixed.Point26_6{
		X: fixed.Int26_6(x * 64),
		Y: fixed.Int26_6(y * 64),
	}
	cache.Drawer.DrawString(text)
}

// 绘制自动换行的文本
func drawWrappedText(img *image.NRGBA, text string, x, y, maxWidth, lineHeight int, face font.Face) int {
	words := []rune(text)
	if len(words) == 0 {
		return y
	}

	// 当前行文本
	currentLine := ""
	currentY := y

	for i := 0; i < len(words); i++ {
		// 尝试添加一个字符
		testLine := currentLine + string(words[i])
		// 测量当前行宽度
		width := font.MeasureString(face, testLine).Round()
		// 如果超过最大宽度，绘制当前行并换行
		if width > maxWidth && len(currentLine) > 0 {
			drawCachedText(img, currentLine, x, currentY, face)
			currentY += lineHeight
			currentLine = string(words[i])
		} else {
			currentLine = testLine
		}
	}

	// 绘制最后一行
	if len(currentLine) > 0 {
		drawCachedText(img, currentLine, x, currentY, face)
		currentY += lineHeight
	}

	// 返回下一行的 Y 坐标
	return currentY
}

// 保存图像到指定路径 - 修正版本，确保图片完全写入
func saveImageToPathFixed(img *image.NRGBA, fileName string) (string, error) {
	// 优先保存到photos目录
	photosDir := filepath.Join(getCurrentDir(), "photos")
	if err := ensureDir(photosDir); err != nil {
		// 如果photos目录创建失败，尝试output目录
		outputDir := filepath.Join(getCurrentDir(), "output")
		if err := ensureDir(outputDir); err != nil {
			return "", fmt.Errorf("无法创建保存目录: %v", err)
		}
		photosDir = outputDir
	}

	// 完整的文件系统路径
	fullPath := filepath.Join(photosDir, fileName)

	// 保存图像
	log.Printf("正在保存图像到: %s", fullPath)
	err := imaging.Save(img, fullPath)
	if err != nil {
		return "", fmt.Errorf("保存图像失败: %v", err)
	}

	// 确保文件完全写入磁盘
	// 检查文件是否存在且大小合理
	if !fileExists(fullPath) {
		return "", fmt.Errorf("图片保存后文件不存在: %s", fullPath)
	}

	// 获取文件信息，确保文件大小合理
	fileInfo, err := os.Stat(fullPath)
	if err != nil {
		return "", fmt.Errorf("无法获取保存文件信息: %v", err)
	}

	// 检查文件大小是否合理（至少10KB，避免空文件或损坏文件）
	if fileInfo.Size() < 10*1024 {
		// 删除可能损坏的文件
		os.Remove(fullPath)
		return "", fmt.Errorf("生成的图片文件过小，可能未完全渲染: %d bytes", fileInfo.Size())
	}

	log.Printf("图片保存成功，文件大小: %d bytes", fileInfo.Size())

	// 返回用于HTTP访问的相对路径
	// 确定是保存在photos还是output目录
	baseDir := filepath.Base(photosDir)
	relativePath := baseDir + "/" + fileName

	log.Printf("返回相对路径: %s", relativePath)
	return relativePath, nil
}

// 创建字体面
func createFontFaces(fontBytes []byte) (font.Face, font.Face, font.Face, error) {
	f, err := opentype.Parse(fontBytes)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("解析字体文件失败: %v", err)
	}

	// 创建各种大小的字体
	titleFace, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    36, // 标题大号字体
		DPI:     78,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("创建标题字体失败: %v", err)
	}

	normalFace, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    30, // 增大正常字体确保可见性
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		titleFace.Close()
		return nil, nil, nil, fmt.Errorf("创建正常字体失败: %v", err)
	}

	smallFace, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    26, // 增大小号字体
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		titleFace.Close()
		normalFace.Close()
		return nil, nil, nil, fmt.Errorf("创建小号字体失败: %v", err)
	}

	return titleFace, normalFace, smallFace, nil
}
