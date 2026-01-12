# Claude Skills Collection

个人 Claude AI 技能集合，包含 33 个实用技能，涵盖开发、设计、文档、测试等多个领域。

## 📊 技能统计

- **总技能数**: 33 个
- **创意与设计**: 6 个
- **文档与内容**: 7 个
- **开发与工程**: 12 个
- **代码审查与协作**: 4 个
- **技能与工具**: 3 个

## 🚀 快速开始

### 环境要求

- Claude Code CLI
- Python 3.x（用于 daily-report 脚本）

### 安装

1. 克隆此仓库到本地：
```bash
git clone https://github.com/YOUR_USERNAME/skills.git
cd skills
```

2. 确保 Claude Code 已配置技能路径：
```bash
# 在 Claude Code 配置中添加此目录
```

## 📁 目录结构

```
skills/
├── skills/                 # 所有技能文件
│   ├── SKILLS_INDEX.md     # 技能索引
│   ├── daily-report/       # 周报生成技能
│   ├── api-design-principles/  # API 设计最佳实践
│   ├── postgresql/         # PostgreSQL 数据库设计
│   └── ...
├── daily/                  # 打卡记录文件目录
└── README.md
```

## 🎯 主要技能

### 文档与内容
- **daily-report**: 自动从 daily/ 目录读取打卡记录，生成结构化周报
- **docx**: Word 文档创建、编辑和分析
- **pdf**: PDF 操作工具包
- **pptx**: PowerPoint 演示文稿创建和编辑

### 开发与工程
- **api-design-principles**: RESTful API 设计最佳实践
- **postgresql**: PostgreSQL 数据库设计和优化
- **frontend-design**: 创建高质量前端界面
- **webapp-testing**: Web 应用测试（使用 Playwright）
- **mcp-builder**: 构建 MCP 服务器

### 创意与设计
- **algorithmic-art**: 算法艺术生成，使用 p5.js 创建生成艺术
- **canvas-design**: 创建视觉艺术作品（.png 和 .pdf）
- **frontend-design**: 创建高质量前端界面

## 💡 使用示例

### 生成周报

```bash
# 1. 在 daily/ 目录下创建打卡文件
echo "WCS:
1、完成系统发布
2、优化监控系统" > daily/$(date +%Y%m%d)

# 2. 使用 daily-report 技能生成周报
# 在 Claude Code 中：使用 daily-report 技能
```

### API 设计

```markdown
# 在 Claude Code 中：
使用 api-design-principles 技能设计 RESTful API
```

## 📦 技能来源

- **Anthropic 官方技能**: [skills-main](https://github.com/anthropics/skills)
- **第三方技能**:
  - [wshobson/agents](https://github.com/wshobson/agents) - api-design-principles, postgresql
- **自定义技能**: daily-report

## 🔧 自定义技能

### 创建新技能

1. 在 `skills/` 目录下创建新文件夹
2. 创建 `SKILL.md` 文件，格式如下：

```markdown
---
name: your-skill-name
description: 技能描述
---

# 技能标题

技能详细说明...
```

3. 更新 `skills/SKILLS_INDEX.md`

## 📝 更新日志

### 2026-01-12
- ✅ 新增 api-design-principles 技能
- ✅ 新增 postgresql 技能
- ✅ 优化 daily-report 技能，添加格式 B-2
- ✅ 清理无用文件和目录

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 👤 作者

[Vigil](https://github.com/YOUR_USERNAME)

---

**注意**: 请确保不要在 `daily/` 目录中提交敏感信息。
