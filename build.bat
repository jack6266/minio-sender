@echo off
echo 开始构建 MinIO 上传工具...

:: 创建必要的目录
mkdir build 2>nul
mkdir build\bin 2>nul
mkdir build\logs 2>nul

:: 构建程序
go build -o build\minio-uploader.exe

:: 复制必要文件
copy paths.txt build\ 2>nul
copy bin\mc.exe build\bin\ 2>nul
copy start.bat build\ 2>nul

echo.
echo 构建完成！
echo 可执行文件位于 build 目录中
echo.
echo 请修改 start.bat 中的 MinIO 配置后运行程序