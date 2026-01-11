package converter

import (
	"fmt"
	"strings"
)

// Theme 主题定义
type Theme struct {
	Name     string // 主题名称
	Type     string // "api" | "ai"
	APITheme string // API 模式使用的主题名
	AIPrompt string // AI 模式使用的提示词
}

// ThemeManager 主题管理器
type ThemeManager struct {
	themes map[string]Theme
}

// NewThemeManager 创建主题管理器
func NewThemeManager() *ThemeManager {
	tm := &ThemeManager{
		themes: make(map[string]Theme),
	}
	tm.initBuiltInThemes()
	return tm
}

// initBuiltInThemes 初始化内置主题
func (tm *ThemeManager) initBuiltInThemes() {
	// ==================== API 模式主题 ====================
	tm.themes["default"] = Theme{
		Name:     "default",
		Type:     "api",
		APITheme: "default",
	}
	tm.themes["leo"] = Theme{
		Name:     "leo",
		Type:     "api",
		APITheme: "leo",
	}

	// ==================== AI 模式主题 ====================

	// 秋日暖光美学
	tm.themes["autumn-warm"] = Theme{
		Name: "autumn-warm",
		Type: "ai",
		AIPrompt: `【终极指令 V4.0】秋日暖光美学兼容性网页设计提示词

指令：
你是一位世界顶级的网页设计师和提示词工程师，专精于温暖治愈和文艺美学，并对代码在不同平台（特别是微信公众号编辑器）的兼容性有深刻理解。你的任务是根据以下经过多轮优化的风格指南和技术要求，创建一个完整、纯粹使用HTML内联样式的单页式网页模板。

核心主题与愿景 (Core Theme & Vision):
创造一个沉浸式、充满治愈感、被秋日暖光浸染的文艺世界。最终成品应如同精致的艺术博客或个人作品集，充满了自然感、柔和光效和清晰的视觉层次。它既要传达信息，本身也要成为一件充满美学价值的数字艺术品。

第一部分：【兼容性优先】结构与技术要求 (Structural & Technical Requirements)

【关键】主容器结构 (Main Container):
- 必须在 <body> 标签之后立即创建一个主 <div> 容器来包裹所有内容
- 所有全局样式（特别是 background-color, padding, display: flex, letter-spacing 等布局样式）必须应用在这个主 <div> 上，而不是 <body>，以确保在微信等环境中背景和布局不丢失
- 主容器 padding 精确设置为 40px 10px

【关键】样式实现 (Styling Implementation):
- 必须使用纯HTML内联样式，禁止使用 <style> 标签或任何外部CSS文件
- 必须为每一个 <p> 标签明确地添加 color: #4a413d; 样式，以防止被微信编辑器强制重置为黑色

模块化与间距 (Modularity & Spacing):
- 内容的核心载体是 <section> 模块（卡片）
- 卡片之间的垂直间距 gap 固定为 40px

第二部分：设计美学与风格指南 (Aesthetics & Style Guide)

色彩方案 (Color Palette):
- 暖白背景: #faf9f5 (应用于主容器)
- 主文字体: #4a413d
- 秋日暖橙 (主强调色): #d97758
- 橙红高亮 (副强调色): #c06b4d
- 引用背景: #fef4e7

卡片式布局 (Card Layout):
- 最大宽度: max-width: 800px
- 内部边距: padding: 25px
- 背景: 必须结合使用 background-color: #ffffff; 和 background-image: linear-gradient(rgba(0,0,0,0.02) 1px, transparent 1px), linear-gradient(90deg, rgba(0,0,0,0.02) 1px, transparent 1px); background-size: 20px 20px; 来实现带有精致米白方格纹理的背景效果
- 边框: border: 1px solid rgba(0, 0, 0, 0.05);
- 暖光阴影: box-shadow: 0 10px 30px rgba(0, 0, 0, 0.04), 0 0 15px rgba(217, 119, 88, 0.4);
- 圆角: border-radius: 18px

第三部分：排版与元素特效 (Typography & Element Effects)

字体 (Font):
- 字体族: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif
- 正文字号: font-size: 16px
- 行高: line-height: 1.75
- 字间距: letter-spacing: 0.5px (应用于主容器)

一级标题 (<h2>) 特效:
- 结构: 必须由两个 <span> 构成：一个用于 ▶ 符号，另一个用于标题文本
- ▶ 符号 <span>: 应用 color: #d97758; 和 text-shadow: 0 0 12px rgba(217, 119, 88, 0.5);
- 标题文本 <span>: 必须应用纯色 color: #d97758;，禁止使用任何渐变色
- 下划线: border-bottom: 1px dashed rgba(74, 65, 61, 0.3);

二级标题 (<h3>) 特效:
- 样式: 必须应用纯色 color: #d97758;，并使用 border-bottom: 2px solid #d97758; 来创建短实线，长度与文字对齐
- 禁止为文字本身添加 text-shadow

加粗/高亮 (<strong>):
- 效果: 文字颜色设为 color: #c06b4d;，禁止附带任何 text-shadow 效果

引用 (<blockquote>):
- 背景: background-color: #fef4e7;
- 左边框: border-left: 5px solid #d97758;
- 阴影: box-shadow: inset 0 0 15px rgba(217, 119, 88, 0.1);
- 禁止为引用内的文字添加 text-shadow

分割线 (<hr>):
- 样式: border: none; height: 1px; background-color: rgba(74, 65, 61, 0.1);

第四部分：最终交付要求 (Final Delivery Requirements)

输出格式: 提供一个完整的、独立的 HTML 内容
代码封装: 将完整的HTML代码包裹在Markdown的代码块中
无外部依赖: 确保代码自包含，字体通过CDN链接，无本地图片

重要补充规则:
1. 图片使用占位符格式：<!-- IMG:index -->，例如第一张图用 <!-- IMG:0 -->
2. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
3. 返回完整的 HTML，不需要其他说明文字

请转换以下 Markdown 内容：`,
	}

	// 春日清新自然
	tm.themes["spring-fresh"] = Theme{
		Name: "spring-fresh",
		Type: "ai",
		AIPrompt: `【终极指令 V4.0】春日清新自然兼容性网页设计提示词

指令：
你是一位世界顶级的网页设计师和提示词工程师，专精于清新自然和生机美学，并对代码在不同平台（特别是微信公众号编辑器）的兼容性有深刻理解。你的任务是根据以下经过多轮优化的风格指南和技术要求，创建一个完整、纯粹使用HTML内联样式的单页式网页模板。

核心主题与愿景 (Core Theme & Vision):
创造一个沉浸式、充满清新感的春日花园世界。最终成品应如同精致的园艺博客或自然杂志，充满了生机感、绿意盎然和清晰的视觉层次。它既要传达信息，本身也要成为一件充满美学价值的数字艺术品。

第一部分：【兼容性优先】结构与技术要求 (Structural & Technical Requirements)

【关键】主容器结构 (Main Container):
- 必须在 <body> 标签之后立即创建一个主 <div> 容器来包裹所有内容
- 所有全局样式（特别是 background-color, padding, display: flex, letter-spacing 等布局样式）必须应用在这个主 <div> 上，而不是 <body>
- 主容器 padding 精确设置为 40px 10px

【关键】样式实现 (Styling Implementation):
- 必须使用纯HTML内联样式，禁止使用 <style> 标签或任何外部CSS文件
- 必须为每一个 <p> 标签明确地添加 color: #3d4a3d; 样式，以防止被微信编辑器强制重置为黑色

模块化与间距 (Modularity & Spacing):
- 内容的核心载体是 <section> 模块（卡片）
- 卡片之间的垂直间距 gap 固定为 40px

第二部分：设计美学与风格指南 (Aesthetics & Style Guide)

色彩方案 (Color Palette):
- 淡绿背景: #f5f8f5 (应用于主容器)
- 主文字体: #3d4a3d (深绿灰)
- 春日嫩绿 (主强调色): #6b9b7a
- 草地翠绿 (副强调色): #4a8058
- 引用背景: #e8f0e8

卡片式布局 (Card Layout):
- 最大宽度: max-width: 800px
- 内部边距: padding: 25px
- 背景: 必须结合使用 background-color: #ffffff; 和 background-image: radial-gradient(circle at 1px 1px, rgba(107, 155, 122, 0.08) 1px, transparent 0); background-size: 20px 20px; 来实现带有清新点状纹理的背景效果
- 边框: border: 1px solid rgba(107, 155, 122, 0.1);
- 清新阴影: box-shadow: 0 8px 24px rgba(74, 128, 88, 0.08), 0 0 12px rgba(107, 155, 122, 0.2);
- 圆角: border-radius: 16px

第三部分：排版与元素特效 (Typography & Element Effects)

字体 (Font):
- 字体族: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif
- 正文字号: font-size: 16px
- 行高: line-height: 1.8
- 字间距: letter-spacing: 0.3px (应用于主容器)

一级标题 (<h2>) 特效:
- 结构: 必须由两个 <span> 构成：一个用于 ❀ 符号，另一个用于标题文本
- ❀ 符号 <span>: 应用 color: #6b9b7a; 和 text-shadow: 0 0 10px rgba(107, 155, 122, 0.4);
- 标题文本 <span>: 必须应用纯色 color: #4a8058;
- 下划线: border-bottom: 1px dashed rgba(74, 128, 88, 0.25);

二级标题 (<h3>) 特效:
- 样式: 必须应用纯色 color: #4a8058;，并使用 border-bottom: 2px solid #6b9b7a; 来创建短实线
- 禁止为文字本身添加 text-shadow

加粗/高亮 (<strong>):
- 效果: 文字颜色设为 color: #4a8058;，禁止附带任何 text-shadow 效果

引用 (<blockquote>):
- 背景: background-color: #e8f0e8;
- 左边框: border-left: 5px solid #6b9b7a;
- 阴影: box-shadow: inset 0 0 12px rgba(107, 155, 122, 0.1);
- 禁止为引用内的文字添加 text-shadow

分割线 (<hr>):
- 样式: border: none; height: 1px; background: linear-gradient(90deg, transparent, rgba(107, 155, 122, 0.3), transparent);

第四部分：最终交付要求 (Final Delivery Requirements)

输出格式: 提供一个完整的、独立的 HTML 内容
代码封装: 将完整的HTML代码包裹在Markdown的代码块中
无外部依赖: 确保代码自包含，字体通过CDN链接，无本地图片

重要补充规则:
1. 图片使用占位符格式：<!-- IMG:index -->
2. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
3. 返回完整的 HTML，不需要其他说明文字

请转换以下 Markdown 内容：`,
	}

	// 深海静谧冷静
	tm.themes["ocean-calm"] = Theme{
		Name: "ocean-calm",
		Type: "ai",
		AIPrompt: `【终极指令 V4.0】深海静谧冷静兼容性网页设计提示词

指令：
你是一位世界顶级的网页设计师和提示词工程师，专精于深邃静谧和理性美学，并对代码在不同平台（特别是微信公众号编辑器）的兼容性有深刻理解。你的任务是根据以下经过多轮优化的风格指南和技术要求，创建一个完整、纯粹使用HTML内联样式的单页式网页模板。

核心主题与愿景 (Core Theme & Vision):
创造一个沉浸式、充满静谧感的深海世界。最终成品应如同精致的专业期刊或学术博客，充满了理性感、深邃蓝调和清晰的视觉层次。它既要传达信息，本身也要成为一件充满美学价值的数字艺术品。

第一部分：【兼容性优先】结构与技术要求 (Structural & Technical Requirements)

【关键】主容器结构 (Main Container):
- 必须在 <body> 标签之后立即创建一个主 <div> 容器来包裹所有内容
- 所有全局样式（特别是 background-color, padding, display: flex, letter-spacing 等布局样式）必须应用在这个主 <div> 上，而不是 <body>
- 主容器 padding 精确设置为 40px 10px

【关键】样式实现 (Styling Implementation):
- 必须使用纯HTML内联样式，禁止使用 <style> 标签或任何外部CSS文件
- 必须为每一个 <p> 标签明确地添加 color: #3a4150; 样式，以防止被微信编辑器强制重置为黑色

模块化与间距 (Modularity & Spacing):
- 内容的核心载体是 <section> 模块（卡片）
- 卡片之间的垂直间距 gap 固定为 40px

第二部分：设计美学与风格指南 (Aesthetics & Style Guide)

色彩方案 (Color Palette):
- 淡蓝背景: #f0f4f8 (应用于主容器)
- 主文字体: #3a4150 (深蓝灰)
- 深海蔚蓝 (主强调色): #4a7c9b
- 静谧石蓝 (副强调色): #3d6a8a
- 引用背景: #e8f0f8

卡片式布局 (Card Layout):
- 最大宽度: max-width: 800px
- 内部边距: padding: 25px
- 背景: 必须结合使用 background-color: #ffffff; 和 background-image: linear-gradient(rgba(74, 124, 155, 0.03) 1px, transparent 1px), linear-gradient(90deg, rgba(74, 124, 155, 0.03) 1px, transparent 1px); background-size: 24px 24px; 来实现带有淡蓝网格纹理的背景效果
- 边框: border: 1px solid rgba(74, 124, 155, 0.08);
- 深海阴影: box-shadow: 0 8px 28px rgba(58, 65, 80, 0.06), 0 0 16px rgba(74, 124, 155, 0.15);
- 圆角: border-radius: 14px

第三部分：排版与元素特效 (Typography & Element Effects)

字体 (Font):
- 字体族: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif
- 正文字号: font-size: 16px
- 行高: line-height: 1.8
- 字间距: letter-spacing: 0.2px (应用于主容器)

一级标题 (<h2>) 特效:
- 结构: 必须由两个 <span> 构成：一个用于 ◆ 符号，另一个用于标题文本
- ◆ 符号 <span>: 应用 color: #4a7c9b; 和 text-shadow: 0 0 10px rgba(74, 124, 155, 0.4);
- 标题文本 <span>: 必须应用纯色 color: #3d6a8a;
- 下划线: border-bottom: 1px dashed rgba(74, 124, 155, 0.3);

二级标题 (<h3>) 特效:
- 样式: 必须应用纯色 color: #3d6a8a;，并使用 border-bottom: 2px solid #4a7c9b; 来创建短实线
- 禁止为文字本身添加 text-shadow

加粗/高亮 (<strong>):
- 效果: 文字颜色设为 color: #3d6a8a;，禁止附带任何 text-shadow 效果

引用 (<blockquote>):
- 背景: background-color: #e8f0f8;
- 左边框: border-left: 5px solid #4a7c9b;
- 阴影: box-shadow: inset 0 0 12px rgba(74, 124, 155, 0.08);
- 禁止为引用内的文字添加 text-shadow

分割线 (<hr>):
- 样式: border: none; height: 1px; background: linear-gradient(90deg, transparent, rgba(74, 124, 155, 0.25), transparent);

第四部分：最终交付要求 (Final Delivery Requirements)

输出格式: 提供一个完整的、独立的 HTML 内容
代码封装: 将完整的HTML代码包裹在Markdown的代码块中
无外部依赖: 确保代码自包含，字体通过CDN链接，无本地图片

重要补充规则:
1. 图片使用占位符格式：<!-- IMG:index -->
2. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
3. 返回完整的 HTML，不需要其他说明文字

请转换以下 Markdown 内容：`,
	}

	// 通用 AI 主题（灵活定制）
	tm.themes["custom"] = Theme{
		Name: "custom",
		Type: "ai",
		AIPrompt: `你是一个专业的微信公众号排版助手。请将以下 Markdown 内容转换为微信公众号兼容的 HTML。

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->
5. 返回完整的 HTML，不需要其他说明文字

请转换以下 Markdown 内容：`,
	}
}

// GetTheme 获取主题
func (tm *ThemeManager) GetTheme(name string) (*Theme, error) {
	theme, ok := tm.themes[name]
	if !ok {
		return nil, fmt.Errorf("theme not found: %s", name)
	}
	return &theme, nil
}

// ListThemes 列出所有主题
func (tm *ThemeManager) ListThemes() []string {
	var names []string
	for name := range tm.themes {
		names = append(names, name)
	}
	return names
}

// ListAIThemes 列出所有 AI 主题
func (tm *ThemeManager) ListAIThemes() []string {
	var names []string
	for name, theme := range tm.themes {
		if theme.Type == "ai" {
			names = append(names, name)
		}
	}
	return names
}

// GetAPITheme 获取 API 模式的主题名
func (tm *ThemeManager) GetAPITheme(name string) (string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}
	if theme.Type != "api" {
		return "", fmt.Errorf("theme '%s' is not an API theme", name)
	}
	return theme.APITheme, nil
}

// GetAIPrompt 获取 AI 模式的提示词
func (tm *ThemeManager) GetAIPrompt(name string) (string, error) {
	theme, err := tm.GetTheme(name)
	if err != nil {
		return "", err
	}
	if theme.Type != "ai" {
		return "", fmt.Errorf("theme '%s' is not an AI theme", name)
	}
	return theme.AIPrompt, nil
}

// BuildCustomAIPrompt 构建自定义 AI 提示词
func BuildCustomAIPrompt(customPrompt string) string {
	if customPrompt == "" {
		return customPrompt
	}

	// 确保包含基本规则
	baseRules := `

## 重要规则
1. 所有 CSS 必须使用内联 style 属性
2. 不使用外部样式表或 <style> 标签
3. 只使用安全的 HTML 标签（section, p, span, strong, em, a, h1-h6, ul, ol, li, blockquote, pre, code, table, img, br, hr）
4. 图片使用占位符格式：<!-- IMG:index -->
5. 返回完整的 HTML，不需要其他说明文字

`

	if !strings.Contains(customPrompt, "重要规则") && !strings.Contains(customPrompt, "规则") {
		customPrompt += baseRules
	}

	if !strings.Contains(customPrompt, "请转换") {
		customPrompt += "\n\n请转换以下 Markdown 内容："
	}

	return customPrompt
}

// IsAPITheme 检查是否是 API 主题
func (tm *ThemeManager) IsAPITheme(name string) bool {
	theme, ok := tm.themes[name]
	if !ok {
		return false
	}
	return theme.Type == "api"
}

// IsAITheme 检查是否是 AI 主题
func (tm *ThemeManager) IsAITheme(name string) bool {
	theme, ok := tm.themes[name]
	if !ok {
		return false
	}
	return theme.Type == "ai"
}

// GetThemeDescription 获取主题描述
func (tm *ThemeManager) GetThemeDescription(name string) string {
	descriptions := map[string]string{
		"default":      "API 默认主题",
		"leo":          "API leo 主题",
		"autumn-warm":  "【秋日暖光】温暖治愈，橙色调，文艺美学",
		"spring-fresh": "【春日清新】清新自然，绿色调，生机盎然",
		"ocean-calm":   "【深海静谧】深邃冷静，蓝色调，理性专业",
		"custom":       "自定义主题，使用您自己的提示词",
	}
	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return "未知主题"
}
