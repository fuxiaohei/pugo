#!/bin/bash

VERSION=$(git describe --abbrev=0)
WORKDIR=$(pwd)

# npx update theme css 
if which npx >/dev/null; then
    cd ./themes/default
    echo "Updating theme css..."
    npx tailwindcss -c tailwind.config.js -i static/css/input.css -o static/css/style.css --minify
else
    echo "npx could not be found, skipping theme css update"
fi

echo "Building release..."
cd $WORKDIR
mkdir -p bin/release
go build -v -ldflags "-s -w -X=main.version=$VERSION" -o bin/release/pugo ./cmd/pugo

cd bin/release

OS=$(uname -s | awk '{print tolower($0)}')
ARCH=$(uname -m)

echo "Creating archive..."
PUGO_VERSION=$(./pugo version)
PUGO_VERSION="${PUGO_VERSION// /-}""-$OS""-$ARCH"
tar -czf $PUGO_VERSION.tar.gz pugo
echo $PUGO_VERSION.tar.gz

rm -rf ./pugo