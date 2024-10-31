#!/bin/bash

# using CGO_ENABLED=0 here to get a statically linked binary 
# which is more ideal for cross-compiling with a more portable binary

# using -ldflags="-s -w" here as:
# -s strips the symbol table from the binary 
# -w strips DWARF debugging information

VERSION_SHA=$1

declare -a platforms=(
  "linux arm64"
  "linux amd64"
  "windows arm64"
  "windows amd64"
)

cd ./src

# Build for each platform
for platform in "${platforms[@]}"; do
  read -r OS ARCH <<< "$platform"

  echo "Building binary action-$OS-$ARCH-$VERSION_SHA"
  CGO_ENABLED=0 GOOS=$OS GOARCH=$ARCH go build \
    -ldflags="-s -w" \
    -trimpath \
    -o ../binaries/action-$OS-$ARCH-$VERSION_SHA \
    ./cmd/diff
done

cd -
