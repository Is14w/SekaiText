import { execSync } from 'child_process';
import { platform } from 'os';

const ports = process.argv.slice(2).map(Number);
const isWindows = platform() === 'win32';

for (const port of ports) {
  try {
    const pids = isWindows ? getWindowsPids(port) : getUnixPids(port);
    for (const pid of pids) {
      try {
        execSync(isWindows ? `taskkill /F /PID ${pid}` : `kill -9 ${pid}`, {
          stdio: 'ignore',
          timeout: 2000,
        });
      } catch {}
    }
  } catch {
    // port not in use
  }
}

function getWindowsPids(port) {
  const out = execSync(`netstat -ano | findstr "127.0.0.1:${port}"`, {
    encoding: 'utf8',
    timeout: 3000,
  });

  return new Set(
    out
      .split('\n')
      .filter((line) => line.includes('LISTENING'))
      .map((line) => line.trim().split(/\s+/).at(-1))
      .filter(Boolean)
  );
}

function getUnixPids(port) {
  const out = execSync(`lsof -ti tcp:${port}`, {
    encoding: 'utf8',
    timeout: 3000,
  });

  return new Set(out.split('\n').map((line) => line.trim()).filter(Boolean));
}
