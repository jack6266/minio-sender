@echo off
echo 设置 MinIO 环境变量...

:: 设置 MinIO 服务器配置
set MINIO_ENDPOINT=http://localhost:9000
set MINIO_ACCESS_KEY=minioadmin
set MINIO_SECRET_KEY=minioadmin
set MINIO_ALIAS=myminio
set MINIO_BUCKET=mybucket

echo MinIO 配置信息：
echo 服务器地址: %MINIO_ENDPOINT%
echo Access Key: %MINIO_ACCESS_KEY%
echo Secret Key: %MINIO_SECRET_KEY%
echo 别名: %MINIO_ALIAS%
echo 存储桶: %MINIO_BUCKET%
echo.

:: 检查 paths.txt 是否存在
if not exist paths.txt (
    echo 错误: paths.txt 文件不存在！
    echo 请确保 paths.txt 文件在当前目录中。
    pause
    exit /b 1
)

echo 开始上传文件...
echo.

:: 运行上传程序
minio-uploader.exe paths.txt

echo.
if errorlevel 1 (
    echo 程序执行出现错误！
) else (
    echo 程序执行完成！
)

pause 