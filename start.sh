#!/bin/bash

echo "设置 MinIO 环境变量..."

# 设置 MinIO 服务器配置
export MINIO_ENDPOINT="http://localhost:9000"
export MINIO_ACCESS_KEY="minioadmin"
export MINIO_SECRET_KEY="minioadmin"
export MINIO_ALIAS="myminio"
export MINIO_BUCKET="mybucket"

# 显示配置信息
echo "MinIO 配置信息："
echo "服务器地址: $MINIO_ENDPOINT"
echo "Access Key: $MINIO_ACCESS_KEY"
echo "Secret Key: $MINIO_SECRET_KEY"
echo "别名: $MINIO_ALIAS"
echo "存储桶: $MINIO_BUCKET"
echo

# 检查 paths.txt 是否存在
if [ ! -f paths.txt ]; then
    echo "错误: paths.txt 文件不存在！"
    echo "请确保 paths.txt 文件在当前目录中。"
    read -p "按回车键退出..."
    exit 1
fi

echo "开始上传文件..."
echo

# 运行上传程序
./minio-uploader-linux paths.txt

# 检查程序执行结果
if [ $? -eq 0 ]; then
    echo
    echo "程序执行完成！"
else
    echo
    echo "程序执行出现错误！"
fi

# 等待用户按键
read -p "按回车键退出..." 