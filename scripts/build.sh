#!/bin/bash
cd "$(dirname "$0")/.." || exit 1

source settings.cfg

PLATFORM="$(uname | tr "[:upper:]" "[:lower:]")"
FILE_NAME="$BOT_NAME-$PLATFORM-$BOT_VERSION"

echo -n "Cleaning old files... "
rm -rf ./dist
echo "done"

echo -n "Building main.go... "
go build -o "dist/$FILE_NAME" ./main.go || (echo "Failed to build"; exit 1)
echo "done"

echo
echo "You can now find a built executable in dist/$FILE_NAME"
echo
