#!/bin/bash

# 创建 build 目录
mkdir -p build

# 构建 Linux 版本
echo "构建 Linux 版本..."
GOOS=linux GOARCH=amd64 go build -o build/linux/minio-uploader-linux
mkdir -p build/linux/bin
cp start.sh build/linux/
cp paths.txt build/linux/
chmod +x build/linux/start.sh
chmod +x build/linux/minio-uploader-linux

# 构建 MacOS 版本
echo "构建 MacOS 版本..."
GOOS=darwin GOARCH=amd64 go build -o build/mac/minio-uploader-mac
mkdir -p build/mac/bin
cp start.sh build/mac/
cp paths.txt build/mac/
chmod +x build/mac/start.sh
chmod +x build/mac/minio-uploader-mac

# 构建 Windows 版本
echo "构建 Windows 版本..."
GOOS=windows GOARCH=amd64 go build -o build/windows/minio-uploader.exe
mkdir -p build/windows/bin
cp start.bat build/windows/
cp paths.txt build/windows/

echo "构建完成！"
echo "可执行文件位于 build 目录中" 