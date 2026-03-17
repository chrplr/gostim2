#!/bin/bash
set -e

# Get the current git tag and commit hash
VERSION=$(git describe --tags --always)
COMMIT=$(git rev-parse HEAD)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Inject them into the binary
LDFLAGS="-X 'gostim2/internal/version.Version=$VERSION' \
         -X 'gostim2/internal/version.GitCommit=$COMMIT' \
         -X 'gostim2/internal/version.BuildTime=$BUILD_TIME'"

echo "Building gostim2 and gostim2-gui with version $VERSION ($COMMIT)..."

# Build CLI version
go build -ldflags "$LDFLAGS" -o gostim2 ./cmd/gostim2

# Build GUI version
go build -ldflags "$LDFLAGS" -o gostim2-gui ./cmd/gostim2-gui

echo "Done."
