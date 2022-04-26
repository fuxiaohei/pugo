#!/bin/bash

VERSION=$(cat ./version)

mkdir -p bin/release
go build -v -ldflags "-s -w -X=main.version=$VERSION" -o bin/release/pugo ./cmd/pugo

cd bin/release

OS=$(uname -s | awk '{print tolower($0)}')
ARCH=$(uname -m)

PUGO_VERSION=$(./pugo version)
PUGO_VERSION="${PUGO_VERSION// /_}""_$OS""_$ARCH"
tar -czf $PUGO_VERSION.tar.gz pugo
echo $PUGO_VERSION.tar.gz