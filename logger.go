package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger 日志记录器
type Logger struct {
	file    *os.File
	logPath string
}

// NewLogger 创建新的日志记录器
func NewLogger() (*Logger, error) {
	// 创建logs目录
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 生成日志文件名
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logFileName := fmt.Sprintf("minio_upload_%s.log", timestamp)
	logPath := filepath.Join(logsDir, logFileName)

	// 创建日志文件
	file, err := os.Create(logPath)
	if err != nil {
		return nil, fmt.Errorf("创建日志文件失败: %v", err)
	}

	return &Logger{
		file:    file,
		logPath: logPath,
	}, nil
}

// Printf 格式化输出日志
func (l *Logger) Printf(format string, args ...interface{}) {
	// 添加时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	// 输出到控制台
	fmt.Print(logMessage)

	// 输出到文件
	if l.file != nil {
		fmt.Fprint(l.file, logMessage)
		l.file.Sync() // 确保写入磁盘
	}
}

// Println 输出一行日志
func (l *Logger) Println(args ...interface{}) {
	// 添加时间戳
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintln(args...)
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)

	// 输出到控制台
	fmt.Print(logMessage)

	// 输出到文件
	if l.file != nil {
		fmt.Fprint(l.file, logMessage)
		l.file.Sync() // 确保写入磁盘
	}
}

// Close 关闭日志文件
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// GetLogPath 获取日志文件路径
func (l *Logger) GetLogPath() string {
	return l.logPath
}
