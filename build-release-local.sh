#!/usr/bin/env bash
set -e

echo "Cleaning Go sidecar binaries..."
rm -f src-tauri/binaries/sekaitext-backend-*
echo "  Deleted."

echo ""
echo "Building release..."
npx tauri build

echo ""
echo "Done!"
