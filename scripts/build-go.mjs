import { execSync } from 'child_process'
import { existsSync, mkdirSync } from 'fs'
import { platform, arch } from 'os'
import { join } from 'path'

// Map target triple (TAURI_TARGET env) -> Go GOOS/GOARCH
const TARGET_MAP = {
  'x86_64-pc-windows-msvc':    { goos: 'windows', goarch: 'amd64', ext: '.exe' },
  'x86_64-unknown-linux-gnu':  { goos: 'linux',   goarch: 'amd64', ext: '' },
  'x86_64-apple-darwin':       { goos: 'darwin',  goarch: 'amd64', ext: '' },
  'aarch64-apple-darwin':      { goos: 'darwin',  goarch: 'arm64', ext: '' },
}

const AUTO_DETECT = {
  'win32-x64':   'x86_64-pc-windows-msvc',
  'linux-x64':   'x86_64-unknown-linux-gnu',
  'darwin-x64':  'x86_64-apple-darwin',
  'darwin-arm64':'aarch64-apple-darwin',
}

const target = process.env.TAURI_TARGET || AUTO_DETECT[`${platform()}-${arch()}`] || null

if (!target) {
  console.error('Cannot determine target triple. Set TAURI_TARGET env var.')
  process.exit(1)
}

const info = TARGET_MAP[target]
if (!info) {
  console.error(`Unknown target triple: ${target}`)
  process.exit(1)
}

const binaryName = `sekaitext-backend-${target}${info.ext}`
const binariesDir = join(import.meta.dirname, '..', 'src-tauri', 'binaries')
const outPath = join(binariesDir, binaryName)

if (existsSync(outPath)) {
  console.log(`[build-go] ${binaryName} already exists, skipping. (Delete it to force rebuild.)`)
  process.exit(0)
}

mkdirSync(binariesDir, { recursive: true })

const ldFlags = target.includes('windows')
  ? '-s -w -H windowsgui'
  : '-s -w'

console.log(`[build-go] Compiling Go backend for ${target} (${info.goos}/${info.goarch})...`)

execSync(
  `go build -ldflags="${ldFlags}" -o "${outPath}" ./cmd/sekaitext/`,
  {
    cwd: join(import.meta.dirname, '..', 'backend'),
    env: { ...process.env, GOOS: info.goos, GOARCH: info.goarch, CGO_ENABLED: '0' },
    stdio: 'inherit',
  },
)

console.log(`[build-go] -> ${outPath}`)
