#!/bin/sh
# Compila gitsync para todos los sistemas y genera dist/checksums.txt.
#   scripts/build.sh v0.2.0
set -eu
VERSION="${1:-dev}"
cd "$(dirname "$0")/.."
rm -rf dist
mkdir -p dist

for target in windows/amd64 windows/arm64 darwin/amd64 darwin/arm64 linux/amd64 linux/arm64; do
  os="${target%/*}"
  arch="${target#*/}"
  out="dist/gitsync-$os-$arch"
  [ "$os" = windows ] && out="$out.exe"
  echo "build $out"
  CGO_ENABLED=0 GOOS="$os" GOARCH="$arch" go build -trimpath \
    -ldflags "-s -w -X main.version=$VERSION" -o "$out" ./cmd/gitsync
done

cd dist
if command -v sha256sum >/dev/null 2>&1; then
  sha256sum gitsync-* > checksums.txt
else
  shasum -a 256 gitsync-* > checksums.txt
fi
cat checksums.txt
