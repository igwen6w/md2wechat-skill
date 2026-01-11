# md2wechat

> **Markdown 转微信公众号的专业工具** — 不懂编程也能轻松使用！

[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue)](LICENSE)

---

## 🚀 5分钟快速上手（新手必看）

> **你不需要懂编程！** 按照下面的步骤操作即可。

### 第一步：下载软件

| 你的系统 | 下载链接 |
|----------|----------|
| ![Windows](https://img.shields.io/badge/🪟-Windows-blue) | [下载 .exe](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-windows-amd64.exe) |
| ![Mac Intel](https://img.shields.io/badge/🍎-Mac%20Intel-grey) | [下载](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-darwin-amd64) |
| ![Mac M1/M2](https://img.shields.io/badge/🍎-Mac%20ARM-grey) | [下载](https://github.com/geekjourneyx/md2wechat-skill/releases/latest/download/md2wechat-darwin-arm64) |

> 💡 **Mac 用户**：下载后如果提示「无法打开」，右键点击 → 打开 → 仍要打开

### 第二步：一键安装

**Windows**：把下载的文件放到 `C:\Windows\System32\` 文件夹

**Mac/Linux**：打开终端，运行
```bash
chmod +x ~/Downloads/md2wechat
sudo mv ~/Downloads/md2wechat /usr/local/bin/
```

### 第三步：配置微信（只需1次）

```bash
md2wechat config init
```

用记事本打开生成的 `md2wechat.yaml`，填入你的微信公众号 AppID 和 Secret

> 📍 在哪里获取？登录 [mp.weixin.qq.com](https://mp.weixin.qq.com) → 设置与开发 → 基本配置

### 第四步：开始使用

```bash
# 预览文章
md2wechat convert 我的文章.md --preview

# 发送到微信草稿箱
md2wechat convert 我的文章.md --draft
```

---

## ✨ 功能特性

| 功能 | 说明 |
|------|------|
| 🎨 **精美主题** | 秋日暖光、春日清新、深海静谧 |
| 🖼️ **自动处理图片** | 上传、压缩、AI 生成 |
| 📝 **一键发草稿** | 直接推送到微信后台 |
| ⚙️ **简单配置** | 配置一次，永久使用 |
| 🔄 **双模式** | API 模式（稳定）+ AI 模式（精美） |

---

## 📚 核心文档

| 文档 | 说明 |
|------|------|
| [新手入门指南](QUICKSTART.md) | **强烈推荐！** 5分钟上手教程 |
| [故障排查](docs/TROUBLESHOOTING.md) | 遇到问题看这里 |
| [常见问题](docs/FAQ.md) | 20+ 常见问题解答 |
| [使用教程](docs/USAGE.md) | 完整功能说明 |

---

## 🎯 常用命令

```bash
# 预览效果
md2wechat convert 文章.md --preview

# 使用精美主题（AI 模式）
md2wechat convert 文章.md --mode ai --theme autumn-warm --preview

# 上传图片 + 创建草稿
md2wechat convert 文章.md --upload --draft

# 查看配置
md2wechat config show

# 验证配置
md2wechat config validate
```

---

## 🎨 精美主题预览

| 主题 | 命令 | 效果 |
|------|------|------|
| 🟠 秋日暖光 | `--mode ai --theme autumn-warm` | 温暖治愈，橙色调 |
| 🟢 春日清新 | `--mode ai --theme spring-fresh` | 清新自然，绿色调 |
| 🔵 深海静谧 | `--mode ai --theme ocean-calm` | 理性专业，蓝色调 |

---

## 💡 使用场景

我是**内容创作者**，我想：
- ✅ 用 Markdown 写文章
- ✅ 一键转换成公众号排版
- ✅ 自动上传图片
- ✅ 直接发到草稿箱

我是**产品经理**，我想：
- ✅ 快速发布产品公告
- ✅ 使用专业的排版样式
- ✅ 不需要懂 HTML/CSS

---

## 🔧 高级安装（开发者）

### 使用 Go 工具链

```bash
go install github.com/geekjourneyx/md2wechat-skill/cmd/md2wechat@latest
```

### 使用安装脚本

**Mac/Linux**：
```bash
curl -fsSL https://raw.githubusercontent.com/geekjourneyx/md2wechat-skill/main/scripts/install.sh | bash
```

**Windows PowerShell**：
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; iex ((New-Object System.Net.WebClient).DownloadString('https://raw.githubusercontent.com/geekjourneyx/md2wechat-skill/main/scripts/install.ps1'))
```

### 从源码编译

```bash
git clone https://github.com/geekjourneyx/md2wechat-skill.git
cd md2wechat-skill
make build
```

---

## 📂 项目结构

```
md2wechat-skill/
├── cmd/            # 命令行工具
├── internal/       # 核心功能
│   ├── converter/  # 转换器
│   ├── image/      # 图片处理
│   ├── draft/      # 草稿服务
│   └── wechat/     # 微信 API
├── docs/           # 文档
├── examples/       # 示例文件
└── scripts/        # 安装脚本
```

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

## 📄 许可证

MIT License
