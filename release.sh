#!/usr/bin/env bash
set -e

read -p "Version (e.g. 0.2.0): " VERSION

echo ""
echo "Updating version..."
npm version "$VERSION" --no-git-tag

echo ""
echo "Syncing version to tauri..."
node scripts/sync-version.mjs

echo ""
echo "Committing and pushing tag..."
git add package.json src-tauri/tauri.conf.json src-tauri/Cargo.toml
git commit -m "v$VERSION"
git tag "v$VERSION"
git push origin master
git push origin "v$VERSION"

echo ""
echo "Done! CI will build and create the release."
