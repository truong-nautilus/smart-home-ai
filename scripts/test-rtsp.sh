#!/bin/bash

# Script để test RTSP video analyzer

echo "🎬 Test RTSP Video Analyzer"
echo "=============================="
echo ""

# Kiểm tra RTSP URL trong .env
if [ -f .env ]; then
    source .env
    if [ -n "$RTSP_URL" ]; then
        echo "✅ Tìm thấy RTSP_URL trong .env: $RTSP_URL"
    else
        echo "⚠️  Không tìm thấy RTSP_URL trong .env, sử dụng mặc định"
    fi
else
    echo "⚠️  Không tìm thấy file .env"
fi

echo ""
echo "🏗️  Building test program..."
go build -o test-rtsp cmd/test-rtsp/main.go

if [ $? -ne 0 ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build thành công!"
echo ""
echo "🚀 Chạy test (phân tích 3 frames, mỗi 5 giây)..."
echo "🛑 Nhấn Ctrl+C để dừng sớm"
echo ""

./test-rtsp

# Cleanup
rm -f test-rtsp

echo ""
echo "🧹 Đã dọn dẹp file binary"
echo "✨ Test hoàn tất!"
