#!/bin/bash

npm run tailwind

PLATFORMS=("darwin" "linux" "windows")

mkdir -p dist

for PLATFORM in "${PLATFORMS[@]}"; do
    GOOS="$PLATFORM"
    echo "Building for $GOOS..."

    export GOOS

    BIN_EXT=""
    if [[ "$GOOS" == "windows" ]]; then
        BIN_EXT=".exe"
    fi

    go build -o "dist/homepage-$GOOS$BIN_EXT"
done

echo "Build completed."