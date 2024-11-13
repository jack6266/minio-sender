package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// PathInfo 存储解析后的路径信息
type PathInfo struct {
	StoreID    string
	SourcePath string
	TargetPath string
}

// ReadPathsFile 读取并解析paths文件
func ReadPathsFile(filename string) ([]PathInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	var paths []PathInfo
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行
		if line == "" {
			continue
		}

		// 跳过注释行
		if strings.HasPrefix(line, "#") {
			continue
		}

		// 使用分号分割
		parts := strings.Split(line, ";")
		if len(parts) != 3 {
			return nil, fmt.Errorf("第 %d 行格式错误: %s (应该包含一个分号)", lineNum, line)
		}

		// 清理路径中的空白字符
		sourcePath := strings.TrimSpace(parts[1])
		targetPath := strings.TrimSpace(parts[2])

		// 验证路径不为空
		if sourcePath == "" || targetPath == "" {
			return nil, fmt.Errorf("第 %d 行包含空路径: %s", lineNum, line)
		}

		paths = append(paths, PathInfo{
			SourcePath: sourcePath,
			TargetPath: targetPath,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件时发生错误: %v", err)
	}

	if len(paths) == 0 {
		return nil, fmt.Errorf("文件 %s 中没有有效的路径配置", filename)
	}

	return paths, nil
}

// ValidateSourcePaths 验证所有源文件是否存在
func ValidateSourcePaths(paths []PathInfo) error {
	for i, path := range paths {
		if _, err := os.Stat(path.SourcePath); os.IsNotExist(err) {
			return fmt.Errorf("源文件不存在 (第 %d 项): %s", i+1, path.SourcePath)
		}
	}
	return nil
}
