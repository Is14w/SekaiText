import { readFileSync, writeFileSync } from 'fs'
import { join } from 'path'

const root = join(import.meta.dirname, '..')

// Read version from package.json
const pkg = JSON.parse(readFileSync(join(root, 'package.json'), 'utf8'))
const version = pkg.version

// Sync to tauri.conf.json
const tauriConfPath = join(root, 'src-tauri', 'tauri.conf.json')
const tauriConf = JSON.parse(readFileSync(tauriConfPath, 'utf8'))
tauriConf.version = version
writeFileSync(tauriConfPath, JSON.stringify(tauriConf, null, 2) + '\n')

// Sync to Cargo.toml
const cargoPath = join(root, 'src-tauri', 'Cargo.toml')
let cargo = readFileSync(cargoPath, 'utf8')
cargo = cargo.replace(/^version\s*=\s*".*"/m, `version = "${version}"`)
writeFileSync(cargoPath, cargo)

console.log(`[sync-version] Synced version ${version} to tauri.conf.json and Cargo.toml`)
