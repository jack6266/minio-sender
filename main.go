package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var logger *Logger

func main() {
	// 初始化日志记录器
	var err error
	logger, err = NewLogger()
	if err != nil {
		fmt.Printf("初始化日志记录器失败: %v\n", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Println("开始执行文件上传程序...")

	// 检查mc工具是否存在
	if !CheckMcExists() {
		logger.Println("未找到 mc 工具，开始下载...")
		if err := DownloadMc(); err != nil {
			logger.Printf("下载 mc 工具失败: %v\n", err)
			os.Exit(1)
		}
	}

	// 获取MinIO配置
	config := NewMinioConfig()
	client := NewMinioClient(config)

	// 配置mc工具
	if err := client.Configure(); err != nil {
		logger.Printf("配置失败: %v\n", err)
		os.Exit(1)
	}

	// 检查命令行参数
	if len(os.Args) != 2 {
		logger.Println("使用方法: minio-uploader <paths.txt>")
		os.Exit(1)
	}

	pathsFile := os.Args[1]

	// 读取paths文件
	paths, err := ReadPathsFile(pathsFile)
	if err != nil {
		logger.Printf("读取paths文件失败: %v\n", err)
		os.Exit(1)
	}

	// 验证源文件是否都存在
	if err := ValidateSourcePaths(paths); err != nil {
		logger.Printf("验证源文件失败: %v\n", err)
		os.Exit(1)
	}

	logger.Println("\n开始上传文件...")
	logger.Println(strings.Repeat("-", 50))

	// 上传所有文件
	for _, path := range paths {
		// 构建完整的目标路径
		targetPath := filepath.Join(config.Alias, path.Bucket, path.TargetPath)
		logger.Printf("\n正在上传: %s -> %s\n", path.SourcePath, targetPath)

		// 上传文件
		if err := client.UploadFile(path.SourcePath, targetPath); err != nil {
			logger.Printf("上传失败 %s: %v\n", path.SourcePath, err)
			continue
		}
	}

	logger.Println("\n所有文件处理完成！")
	logger.Printf("日志文件保存在: %s\n", logger.GetLogPath())
}
