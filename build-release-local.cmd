@echo off
chcp 65001 >nul

echo Cleaning Go sidecar binaries...
if exist "src-tauri\binaries\sekaitext-backend-*" (
    del /q "src-tauri\binaries\sekaitext-backend-*"
    echo   Deleted.
)

echo.
echo Building release...
call npx tauri build

echo.
echo Done!
pause
