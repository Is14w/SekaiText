#!/usr/bin/env node
// Kill processes listening on specified ports.
// Usage: node scripts/cleanup.mjs <port1> <port2> ...

import { execSync } from 'child_process'

const ports = process.argv.slice(2)
if (ports.length === 0) {
  process.exit(0)
}

for (const port of ports) {
  try {
    const cmd =
      process.platform === 'win32'
        ? `netstat -ano | findstr :${port}`
        : `lsof -ti:${port}`

    const stdout = execSync(cmd, { encoding: 'utf8', timeout: 5000 })
    const lines = stdout.trim().split('\n')

    for (const line of lines) {
      let pid
      if (process.platform === 'win32') {
        if (!line.includes('LISTENING')) continue
        const parts = line.trim().split(/\s+/)
        pid = parts[parts.length - 1]
      } else {
        pid = line.trim()
      }
      if (pid && /^\d+$/.test(pid)) {
        try {
          const killCmd =
            process.platform === 'win32'
              ? `taskkill /F /PID ${pid}`
              : `kill -9 ${pid}`
          execSync(killCmd, { stdio: 'ignore', timeout: 3000 })
        } catch {
          // process already gone
        }
      }
    }
  } catch {
    // no process on this port
  }
}
