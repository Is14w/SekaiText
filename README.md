# SekaiText

> Project Sekai story translation editor

![SekaiText](public/app-icon.png)

SekaiText is a desktop application for browsing, translating, and proofreading Project Sekai: Colorful Stage! story scenarios. It combines a Vue 3 frontend with a Go backend, packaged as a Tauri desktop app.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | [Vue 3](https://vuejs.org/) + [TypeScript](https://www.typescriptlang.org/) + [Vite](https://vitejs.dev/) + [Tailwind CSS v4](https://tailwindcss.com/) |
| State | [Pinia](https://pinia.vuejs.org/) |
| Routing | [Vue Router v5](https://router.vuejs.org/) |
| Backend | [Go 1.24](https://go.dev/) + [chi router](https://github.com/go-chi/chi) |
| Desktop Shell | [Tauri v2](https://v2.tauri.app/) (Rust) |
| Icons | [Lucide](https://lucide.dev/) via `lucide-vue-next` |
| Tables | [TanStack Table](https://tanstack.com/table/v8) + [TanStack Virtual](https://tanstack.com/virtual/v3) |

## Project Structure

```
sekaitext/
├── backend/                 # Go backend server
│   ├── cmd/server/main.go   # Entry point (port 9800)
│   └── internal/
│       ├── api/             # HTTP handlers + router
│       ├── config/          # App configuration
│       ├── model/           # Data types
│       └── service/         # Business logic
├── src/                     # Vue 3 frontend
│   ├── api/                 # API client
│   ├── components/          # Vue components
│   ├── composables/         # Shared composables
│   ├── pages/               # Route pages
│   ├── stores/              # Pinia stores
│   └── types/               # TypeScript types
├── src-tauri/               # Tauri desktop shell (Rust)
│   ├── capabilities/        # Tauri v2 permission capabilities
│   └── src/                 # Rust source
├── scripts/                 # Utility scripts
├── public/                  # Static assets
└── package.json
```

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) >= 20
- [Go](https://go.dev/) >= 1.24
- [Tauri prerequisites](https://v2.tauri.app/start/prerequisites/) (Rust toolchain, system dependencies)

### Install

```bash
# Install frontend dependencies
npm install

# Install Tauri CLI
npm install -D @tauri-apps/cli
```

### Development

```bash
# Start both Go backend and Vite dev server concurrently
npm run dev:web

# Or start with Tauri desktop window
npm run dev:tauri
```

The Go API server runs on `http://localhost:9800` and the Vite dev server on `http://localhost:5173`.

### Build

```bash
# Build the Tauri desktop application
npx tauri build
```

## Configuration

Settings are available in-app via the Settings page (gear icon in the sidebar):

- **Font Size** — Editor text display size (10–48px)
- **Download Source** — Story JSON data source
- **Save \\N** — Preserve `\N` line break markers in translation files
- **Save Voice** — Download and save voice files locally
- **SSL Verification** — Disable SSL certificate verification (needed in some network environments)
- **Index Order** — Story index dropdown sort order (ascending / descending)
- **Dark Mode** — Toggle between light and dark themes
- **Debug Log** — Show debug log panel at the bottom of the editor

## Version

Current version: **0.1.0** (alpha)

## License

This project is licensed under the [MIT License](LICENSE).

*Project Sekai is a trademark of SEGA / Colorful Palette. This project is for educational and fan-translation purposes only.*
