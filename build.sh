#!/bin/bash

# exit when any command fails
set -e

declare -a GOOS=("darwin" "linux" "windows")
declare -a GOARCH=("386" "amd64")

for goos in "${GOOS[@]}"; do
  for goarch in "${GOARCH[@]}"; do
    bin="bin/copy-basta.$goos-$goarch"
    if [[ "$bin" == *"windows"* ]]; then
      bin="$bin.exe"
    fi
    cmd="GOOS=$goos GOARCH=$goarch go build -a -o $bin ./cmd/"
    echo "$cmd"
    eval "$cmd"
  done
done
