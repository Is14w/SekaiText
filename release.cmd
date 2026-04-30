@echo off
chcp 65001 >nul

set /p VERSION="Version (e.g. 0.2.0): "

echo.
echo Updating version...
call npm version %VERSION% --no-git-tag

echo.
echo Syncing version to tauri...
call node scripts/sync-version.mjs

echo.
echo Committing and pushing tag...
git add package.json src-tauri/tauri.conf.json src-tauri/Cargo.toml
git commit -m "v%VERSION%"
git tag v%VERSION%
git push origin master
git push origin v%VERSION%

echo.
echo Done! CI will build and create the release.
pause
