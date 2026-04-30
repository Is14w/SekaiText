@echo off
chcp 65001 >nul

set /p VERSION="Version (e.g. 0.2.0): "

:: Check if version already matches package.json
for /f "delims=" %%v in ('node -p "require('./package.json').version"') do set CURRENT_VERSION=%%v

if "%CURRENT_VERSION%"=="%VERSION%" (
  echo.
  echo Version already %VERSION%, skipping bump.
) else (
  echo.
  echo Updating version...
  call npm version %VERSION% --no-git-tag
)

echo.
echo Syncing version to tauri...
call node scripts/sync-version.mjs

echo.
echo Committing and pushing tag...
git add package.json src-tauri/tauri.conf.json src-tauri/Cargo.toml

:: Check if there are staged changes before committing
git diff --cached --quiet
if errorlevel 1 (
  git commit -m "v%VERSION%"
) else (
  echo No changes to commit, reusing existing commit.
)

:: Force-recreate the tag in case the previous build failed
git tag -f v%VERSION%

git push origin master

:: Force-push the tag (may exist from a failed previous run)
git push -f origin v%VERSION%

echo.
echo Done! CI will build and create the release.
pause
