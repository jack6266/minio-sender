package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// MinioConfig 配置结构
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Alias     string
	Bucket    string
}

// NewMinioConfig 从环境变量获取MinIO配置
func NewMinioConfig() MinioConfig {
	return MinioConfig{
		Endpoint:  getEnvOrDefault("MINIO_ENDPOINT", "http://localhost:9000"),
		AccessKey: getEnvOrDefault("MINIO_ACCESS_KEY", "erdcloud"),
		SecretKey: getEnvOrDefault("MINIO_SECRET_KEY", "Pw!123456"),
		Alias:     getEnvOrDefault("MINIO_ALIAS", "myminio"),
		Bucket:    getEnvOrDefault("MINIO_BUCKET", "plat"),
	}
}

// getEnvOrDefault 获取环境变量，如果不存在则返回默认值
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// MinioClient MinIO客户端结构
type MinioClient struct {
	config MinioConfig
}

// NewMinioClient 创建新的MinIO客户端
func NewMinioClient(config MinioConfig) *MinioClient {
	return &MinioClient{
		config: config,
	}
}

// checkSystemMc 检查系统中是否存在mc工具
func checkSystemMc() (string, bool) {
	mcPath, err := exec.LookPath("mc")
	if err == nil {
		return mcPath, true
	}
	return "", false
}

// getMcPath 获取mc工具的完整路径
func getMcPath() string {
	// 首先检查系统中是否存在mc
	if systemMcPath, exists := checkSystemMc(); exists {
		return systemMcPath
	}

	// 如果系统中不存在，则使用本地bin目录
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	// bin目录路径
	binDir := filepath.Join(execDir, "bin")

	// mc可执行文件名
	mcName := "mc"
	if runtime.GOOS == "windows" {
		mcName = "mc.exe"
	}

	return filepath.Join(binDir, mcName)
}

// CheckMcExists 检查mc工具是否存在
func CheckMcExists() bool {
	// 首先检查系统中是否存在mc
	if _, exists := checkSystemMc(); exists {
		fmt.Println("检测到系统已安装mc工具")
		return true
	}

	// 如果系统中不存在，检查本地bin目录
	mcPath := getMcPath()
	if _, err := os.Stat(mcPath); err == nil {
		fmt.Println("检测到本地bin目录中存在mc工具")
		return true
	}

	fmt.Println("未检测到mc工具")
	return false
}

// DownloadMc 下载mc工具
func DownloadMc() error {
	var downloadURL string
	switch runtime.GOOS {
	case "linux":
		downloadURL = "https://dl.min.io/client/mc/release/linux-amd64/mc"
	case "darwin":
		downloadURL = "https://dl.min.io/client/mc/release/darwin-amd64/mc"
	case "windows":
		downloadURL = "https://dl.min.io/client/mc/release/windows-amd64/mc.exe"
	default:
		return fmt.Errorf("不支持的操作系统: %s", runtime.GOOS)
	}

	fmt.Println("正在下载 mc 工具...")

	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()

	// 确保bin目录存在
	mcPath := getMcPath()
	binDir := filepath.Dir(mcPath)
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("创建bin目录失败: %v", err)
	}

	// 创建文件
	out, err := os.OpenFile(mcPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer out.Close()

	// 写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("保存文件失败: %v", err)
	}

	// 设置执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(mcPath, 0755); err != nil {
			return fmt.Errorf("设置执行权限失败: %v", err)
		}
	}

	fmt.Println("mc 工具下载完成")
	return nil
}

// Configure 配置mc工具并测试连接
func (c *MinioClient) Configure() error {
	mcPath := getMcPath()

	// 配置别名
	configCmd := exec.Command(mcPath, "alias", "set", c.config.Alias, c.config.Endpoint, c.config.AccessKey, c.config.SecretKey)
	output, err := configCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("配置mc失败: %v\n输出: %s", err, string(output))
	}
	fmt.Println("mc配置成功")

	// 测试连接
	fmt.Printf("正在测试连接 MinIO 服务器 (%s)...\n", c.config.Endpoint)
	testCmd := exec.Command(mcPath, "ls", c.config.Alias)
	output, err = testCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("连接MinIO服务器失败: %v\n输出: %s", err, string(output))
	}

	// 测试存储桶访问
	bucketCmd := exec.Command(mcPath, "ls", fmt.Sprintf("%s/%s", c.config.Alias, c.config.Bucket))
	output, err = bucketCmd.CombinedOutput()
	if err != nil {
		// 如果存储桶不存在，尝试创建
		createBucketCmd := exec.Command(mcPath, "mb", fmt.Sprintf("%s/%s", c.config.Alias, c.config.Bucket))
		output, err = createBucketCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("创建存储桶失败: %v\n输出: %s", err, string(output))
		}
		fmt.Printf("成功创建存储桶: %s\n", c.config.Bucket)
	}

	fmt.Printf("\n连接成功！\n")
	fmt.Printf("MinIO服务器: %s\n", c.config.Endpoint)
	fmt.Printf("存储桶: %s\n", c.config.Bucket)
	fmt.Println(strings.Repeat("-", 50))

	return nil
}

// UploadFile 上传文件到MinIO，指定目标路径
func (c *MinioClient) UploadFile(sourcePath, targetPath string) error {
	mcPath := getMcPath()

	cmd := exec.Command(mcPath, "cp", sourcePath, targetPath)
	fmt.Printf("cmd: %s\n", cmd.String())

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("上传失败: %v\n输出: %s", err, string(output))
	}

	// 获取文件大小
	fileInfo, err := os.Stat(sourcePath)
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 打印成功消息
	fmt.Printf("\n上传成功！\n")
	fmt.Printf("源文件: %s\n", sourcePath)
	fmt.Printf("目标路径: %s\n", targetPath)
	fmt.Printf("文件大小: %.2f MB\n", float64(fileInfo.Size())/(1024*1024))
	fmt.Printf("mc输出: %s\n", string(output))
	fmt.Println(strings.Repeat("-", 50))

	return nil
}
