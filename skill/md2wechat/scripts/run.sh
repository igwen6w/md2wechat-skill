#!/bin/bash
#
# md2wechat 统一入口脚本
# 自动检测平台并调用对应的二进制文件
#

set -e

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 检测操作系统和架构
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# 根据平台选择二进制文件
case "$OS" in
    darwin)
        # macOS
        if [ "$ARCH" = "arm64" ]; then
            BINARY="$SCRIPT_DIR/bin/md2wechat_darwin_arm64"
        else
            BINARY="$SCRIPT_DIR/bin/md2wechat_darwin_amd64"
        fi
        ;;
    linux)
        BINARY="$SCRIPT_DIR/bin/md2wechat_linux_amd64"
        ;;
    msys*|mingw*|cygwin*)
        BINARY="$SCRIPT_DIR/bin/md2wechat_windows_amd64.exe"
        ;;
    *)
        echo "Error: Unsupported OS: $OS" >&2
        echo "Supported platforms: macOS (x64/arm64), Linux (x64), Windows (x64)" >&2
        exit 1
        ;;
esac

# 检查二进制文件是否存在
if [ ! -f "$BINARY" ]; then
    echo "Error: Binary not found: $BINARY" >&2
    echo "Please build the project first or download the correct binary for your platform." >&2
    exit 1
fi

# 检查二进制文件是否可执行
if [ ! -x "$BINARY" ]; then
    chmod +x "$BINARY" 2>/dev/null || true
fi

# 执行二进制文件
exec "$BINARY" "$@"
