# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Claude Code Skill support with `.claude-plugin/plugin.json`
- Claude Code Skill in `skill/md2wechat/` directory for distribution
- New API themes: `bytedance`, `apple`, `sports`, `chinese`, `cyber`
- Comprehensive troubleshooting guide in SKILL.md
- API theme selection section in README.md

### Changed
- Updated themes.md with complete API and AI theme documentation
- Enhanced HTML guide with AI theme specific requirements
- Improved SKILL.md with detailed error handling and FAQ
- Updated command help text to reflect all available themes
- Enhanced FAQ.md with IP whitelist configuration guide

### Removed
- `leo` theme (deprecated)

---

## [1.1.0] - 2025-01-11

### Added
- YAML-based theme system for AI mode
- AI themes: `autumn-warm`, `spring-fresh`, `ocean-calm`
- Custom theme support with `custom.yaml`
- Theme configuration with color schemes and style info
- Reference documentation for themes, HTML guide, image syntax, and WeChat API

### Changed
- Refactored theme system to use YAML files instead of hardcoded prompts
- Updated SKILL.md with new AI theme workflow
- Enhanced themes.md with detailed style specifications

---

## [1.0.1] - 2025-01-11

### Fixed
- Mermaid diagrams rendering for GitHub documentation

---

## [1.0.0] - 2025-01-11

### Added
- Initial release of md2wechat
- API mode conversion using md2wechat.cn API
- AI mode conversion with Claude AI
- WeChat draft upload functionality
- Image upload to WeChat material library
- Configuration management with YAML support
- Command-line interface with cobra
- Multi-platform binary support (Windows, macOS, Linux)
- Comprehensive documentation (README, FAQ, USAGE)
- Test draft command for HTML validation

### Features
- Convert Markdown to WeChat Official Account formatted HTML
- Support for local images, online images, and AI-generated images
- Automatic image compression and optimization
- Draft creation with cover image support
- Environment variable and config file support
- Claude Code Skill integration

---

## Version History Summary

| Version | Date | Description |
|---------|------|-------------|
| [Unreleased] | TBD | Claude Code plugin support, new API themes, enhanced documentation |
| [1.1.0] | 2025-01-11 | YAML theme system, AI themes (autumn-warm, spring-fresh, ocean-calm) |
| [1.0.1] | 2025-01-11 | Fixed Mermaid diagrams for GitHub rendering |
| [1.0.0] | 2025-01-11 | Initial release with full md2wechat functionality |

---

## Upgrade Guide

### From v1.1.0 to Unreleased

**New API Themes Available:**
```bash
# New themes available
md2wechat convert article.md --theme bytedance
md2wechat convert article.md --theme apple
md2wechat convert article.md --theme sports
md2wechat convert article.md --theme chinese
md2wechat convert article.md --theme cyber
```

**Claude Code Skill Integration:**
```bash
# Install as Claude Code Skill
cp -r skill/md2wechat ~/.claude/skills/
```

### From v1.0.0 to v1.1.0

**Theme System Migration:**
- Old theme names have been updated
- AI themes now use YAML configuration files
- Update your commands to use new theme names:
  - `autumn-warm` instead of `elegant`
  - `spring-fresh` instead of `minimal`
  - `ocean-calm` instead of `tech`

---

## Links

- [GitHub Repository](https://github.com/geekjourneyx/md2wechat-skill)
- [Documentation](README.md)
- [Issues](https://github.com/geekjourneyx/md2wechat-skill/issues)
