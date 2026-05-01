<!-- markdownlint-disable-next-line MD036 -->
_✨ Project Sekai 剧情翻译编辑器 ✨_

<!-- prettier-ignore-start -->
<!-- markdownlint-disable-next-line MD036 -->

<p align="center">
  <a href="https://github.com/Is14w/SekaiText/blob/master/README.md">简体中文</a> |
  <a href="https://github.com/Is14w/SekaiText/blob/master/README.en.md">English</a>
</p>

<p align="center">
  <img alt="SekaiText" src="public/app-icon.png" width="128" height="128" />
</p>

<p align="center">
  <img alt="License" src="https://img.shields.io/github/license/Is14w/SekaiText?style=flat-square&color=ff69b4" />
  <img alt="Version" src="https://img.shields.io/badge/version-0.1.1-blue?style=flat-square" />
  <img alt="Tauri" src="https://img.shields.io/badge/Tauri-v2-ffc131?style=flat-square&logo=tauri&logoColor=white" />
  <img alt="Vue" src="https://img.shields.io/badge/Vue-3.5-4fc08d?style=flat-square&logo=vuedotjs&logoColor=white" />
  <img alt="Go" src="https://img.shields.io/badge/Go-1.24-00add8?style=flat-square&logo=go&logoColor=white" />
</p>

<p align="center">
  <img alt="Platforms" src="https://img.shields.io/badge/platform-Windows%20%7C%20macOS%20%7C%20Linux-6c5ce7?style=flat-square" />
</p>

<!-- prettier-ignore-end -->

> **注意**：本项目是一个 **Vibe Coding 产物**，目前处于早期 alpha 阶段。欢迎贡献和反馈！

SekaiText 是一款桌面应用程序，用于浏览、翻译和校对《世界计划 彩色舞台 feat. 初音未来》的剧情文本。它结合了 Vue 3 前端与 Go 后端，打包为 Tauri 桌面应用。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) + [Vite](https://vitejs.dev/) + [Tailwind CSS v4](https://tailwindcss.com/) |
| 状态管理 | [Pinia](https://pinia.vuejs.org/) |
| 路由 | [Vue Router v5](https://router.vuejs.org/) |
| 后端 | [Go 1.24](https://go.dev/) + [chi router](https://github.com/go-chi/chi) |
| 桌面壳 | [Tauri v2](https://v2.tauri.app/) (Rust) |
| 图标 | [Lucide](https://lucide.dev/) via `lucide-vue-next` |
| 表格 | [TanStack Table](https://tanstack.com/table/v8) + [TanStack Virtual](https://tanstack.com/virtual/v3) |

## 项目结构

```
sekaitext/
├── backend/                 # Go 后端服务
│   ├── cmd/server/main.go   # 入口（端口 9800）
│   └── internal/
│       ├── api/             # HTTP 路由与处理器
│       ├── config/          # 应用配置
│       ├── model/           # 数据类型
│       └── service/         # 业务逻辑
├── src/                     # Vue 3 前端
│   ├── api/                 # API 客户端
│   ├── components/          # Vue 组件
│   ├── composables/         # 通用组合式函数
│   ├── pages/               # 路由页面
│   ├── stores/              # Pinia 状态
│   └── types/               # TypeScript 类型
├── src-tauri/               # Tauri 桌面壳 (Rust)
│   ├── capabilities/        # Tauri v2 权限
│   └── src/                 # Rust 源码
├── scripts/                 # 工具脚本
├── public/                  # 静态资源
└── package.json
```

## 快速开始

### 环境要求

- [Node.js](https://nodejs.org/) >= 20
- [Go](https://go.dev/) >= 1.24
- [Tauri 依赖](https://v2.tauri.app/start/prerequisites/)（Rust 工具链、系统依赖）

### 安装

```bash
# 安装前端依赖
npm install

# 安装 Tauri CLI
npm install -D @tauri-apps/cli
```

### 开发

```bash
# 同时启动 Go 后端 + Vite 开发服务器
npm run dev:web

# 或使用 Tauri 桌面窗口启动
npm run dev:tauri
```

Go API 服务运行在 `http://localhost:9800`，Vite 开发服务器运行在 `http://localhost:5173`。

### 构建

```bash
# 构建 Tauri 桌面应用
npm run build:tauri
```

## 配置

应用内通过侧栏设置页面（齿轮图标）调整：

- **字号** — 编辑器文本显示大小（10–48px）
- **索引排序** — 故事索引下拉列表的排列顺序
- **保存 \\N 换行符** — 翻译文件中保留 `\N` 换行标记
- **保存语音文件** — 下载并保存语音文件到本地
- **SSL 验证** — 禁用 SSL 证书验证（某些网络环境需要）
- **暗色主题** — 切换亮色 / 暗色显示
- **调试日志** — 在编辑器底部显示调试日志面板

## 版本

当前版本：**0.1.1** (alpha)

## 开源协议

本项目采用 [MIT License](LICENSE) 协议。

*《世界计划 彩色舞台 feat. 初音未来》是 SEGA / Colorful Palette 的商标。本项目仅用于学习和粉丝翻译目的。*
