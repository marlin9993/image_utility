package main

import (
	"flag" // 用于处理命令行参数
	"fmt"
	"image"
	// 确保支持常见的图片格式
	_ "image/gif"  // 支持 GIF 格式
	_ "image/jpeg" // 支持 JPEG 格式
	"image/png"    // 用于 PNG 编码及解码
	"log"
	"os"
)

// processImage 处理图片的加载、获取信息及保存为PNG。
// inputPath: 输入图片路径
// outputPath: 输出图片路径
func processImage(inputPath string, outputPath string) error {
	// 1. 打开输入图片文件
	// 使用 os.Open 打开文件进行读取
	inputFile, err := os.Open(inputPath)
	if err != nil {
		// 如果打开失败，返回包含原始错误信息的错误
		return fmt.Errorf("打开图片 '%s' 失败: %w", inputPath, err)
	}
	// 使用 defer 确保在函数退出时关闭文件
	defer inputFile.Close()

	// 2. 解码图片
	// image.Decode 会自动识别图片格式并解码
	// 返回 image.Image 接口（包含图片数据），格式字符串，以及可能的错误
	img, format, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("解码图片 '%s' 失败: %w", inputPath, err)
	}

	// 从解码后的图片获取尺寸信息
	bounds := img.Bounds() // Bounds() 返回一个 image.Rectangle 表示图片的边界
	width := bounds.Dx()   // Dx() 返回宽度
	height := bounds.Dy()  // Dy() 返回高度
	log.Printf("图片加载成功: 格式=%s, 宽=%d, 高=%d\n", format, width, height)

	// 3. 创建输出图片文件
	// 使用 os.Create 创建文件用于写入，如果文件已存在则会清空
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("创建输出文件 '%s' 失败: %w", outputPath, err)
	}
	// 同样使用 defer 确保输出文件关闭
	defer outputFile.Close()

	// 4. 将图片编码为 PNG 格式并保存
	// png.Encode 将 image.Image 对象编码为 PNG 格式并写入到指定的 io.Writer
	if err := png.Encode(outputFile, img); err != nil {
		return fmt.Errorf("保存 PNG 到 '%s' 失败: %w", outputPath, err)
	}

	log.Printf("图片处理完成并保存到: %s\n", outputPath)
	return nil // 表示成功
}

func main() {
	// 定义命令行参数
	// flag.String 返回一个指向字符串的指针
	var (
		input  string
		output string
	)

	flag.StringVar(&input, "input", "input.png", "输入图片文件的路径")
	flag.StringVar(&output, "output", "crop_optimized.png", "处理后输出图片文件的路径")

	// 解析命令行参数
	flag.Parse()

	// 检查输入文件是否存在，如果不存在，则创建一个虚拟的PNG图片用于演示
	// 这个逻辑现在会基于命令行传入的 inputFilename
	if _, err := os.Stat(input); os.IsNotExist(err) {
		log.Printf("警告: 输入文件 '%s' 不存在", input)
		return
	}

	// 调用处理图片的函数
	log.Printf("开始处理图片: 输入='%s', 输出='%s'", input, output)
	if err := processImage(input, output); err != nil {
		// 如果处理过程中发生错误，记录致命错误并退出程序
		log.Fatalf("处理图片时发生错误: %v", err)
	}
}
