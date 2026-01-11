# md2wechat Makefile
# 适用于开发者和高级用户

.PHONY: all build clean test install help lint fmt vet

# 默认目标
all: build

# 构建所有平台的二进制文件
build-all:
	@echo "构建所有平台..."
	@echo "Building for Linux amd64..."
	GOOS=linux GOARCH=amd64 go build -o dist/md2wechat-linux-amd64 ./cmd/md2wechat
	@echo "Building for Linux arm64..."
	GOOS=linux GOARCH=arm64 go build -o dist/md2wechat-linux-arm64 ./cmd/md2wechat
	@echo "Building for macOS amd64..."
	GOOS=darwin GOARCH=amd64 go build -o dist/md2wechat-darwin-amd64 ./cmd/md2wechat
	@echo "Building for macOS arm64..."
	GOOS=darwin GOARCH=arm64 go build -o dist/md2wechat-darwin-arm64 ./cmd/md2wechat
	@echo "Building for Windows amd64..."
	GOOS=windows GOARCH=amd64 go build -o dist/md2wechat-windows-amd64.exe ./cmd/md2wechat
	@echo "Done! Binaries in dist/"

# 构建当前平台
build:
	@echo "构建当前平台..."
	go build -o md2wechat ./cmd/md2wechat
	@echo "构建完成: ./md2wechat"

# 清理
clean:
	@echo "清理..."
	rm -f md2wechat
	rm -rf dist/
	rm -f *.log

# 运行测试
test:
	@echo "运行测试..."
	go test -v ./...

# 代码检查
lint:
	@echo "代码检查..."
	golangci-lint run ./...

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...
	gofmt -w .

# 静态分析
vet:
	@echo "静态分析..."
	go vet ./...

# 安装到 GOPATH/bin
install:
	@echo "安装到 $(GOPATH)/bin..."
	go install ./cmd/md2wechat

# 下载依赖
deps:
	@echo "下载依赖..."
	go mod download
	go mod tidy

# 帮助
help:
	@echo "md2wechat Makefile 命令:"
	@echo ""
	@echo "  make build       - 构建当前平台二进制"
	@echo "  make build-all   - 构建所有平台二进制"
	@echo "  make clean       - 清理构建文件"
	@echo "  make test        - 运行测试"
	@echo "  make fmt         - 格式化代码"
	@echo "  make vet         - 静态分析"
	@echo "  make install     - 安装到 GOPATH/bin"
	@echo "  make deps        - 下载依赖"
	@echo ""
	@echo "用户快速安装:"
	@echo "  go install github.com/geekjourneyx/md2wechat-skill/cmd/md2wechat@latest"
