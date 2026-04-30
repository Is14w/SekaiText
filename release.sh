#!/usr/bin/env bash
set -e

read -p "Version (e.g. 0.2.0): " VERSION

CURRENT_VERSION=$(node -p "require('./package.json').version")

if [ "$CURRENT_VERSION" != "$VERSION" ]; then
  echo ""
  echo "Updating version..."
  npm version "$VERSION" --no-git-tag
else
  echo ""
  echo "Version already $VERSION, skipping bump."
fi

echo ""
echo "Syncing version to tauri..."
node scripts/sync-version.mjs

echo ""
echo "Committing and pushing tag..."
git add package.json src-tauri/tauri.conf.json src-tauri/Cargo.toml

if ! git diff --cached --quiet; then
  git commit -m "v$VERSION"
else
  echo "No changes to commit, reusing existing commit."
fi

# Force-recreate the tag in case the previous build failed
git tag -f "v$VERSION"

git push origin master

# Force-push the tag (may exist from a failed previous run)
git push -f origin "v$VERSION"

echo ""
echo "Done! CI will build and create the release."
