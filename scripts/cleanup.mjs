import { execSync } from 'child_process';

const ports = process.argv.slice(2).map(Number);

for (const port of ports) {
  try {
    const out = execSync(`netstat -ano | findstr "127.0.0.1:${port}"`, {
      encoding: 'utf8',
      timeout: 3000,
    });
    const pids = new Set(
      out
        .split('\n')
        .filter((l) => l.includes('LISTENING'))
        .map((l) => l.trim().split(/\s+/).at(-1))
        .filter(Boolean)
    );
    for (const pid of pids) {
      try {
        execSync(`taskkill /F /PID ${pid}`, { stdio: 'ignore', timeout: 2000 });
      } catch {}
    }
  } catch {
    // port not in use
  }
}
