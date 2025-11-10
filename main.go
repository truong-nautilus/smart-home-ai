package main

// Đây là điểm vào cũ.
// Ứng dụng đã được tái cấu trúc sử dụng Clean Architecture.
// Vui lòng sử dụng: go run cmd/assistant/main.go

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("⚠️  Điểm vào này đã lỗi thời.")
	fmt.Println("📦 Dự án đã được tái cấu trúc sử dụng Clean Architecture.")
	fmt.Println("🚀 Vui lòng chạy: go run cmd/assistant/main.go")
	fmt.Println()
	fmt.Println("Hoặc build file thực thi:")
	fmt.Println("   go build -o tro-ly-thong-minh cmd/assistant/main.go")
	fmt.Println("   ./tro-ly-thong-minh")
	os.Exit(0)
}
