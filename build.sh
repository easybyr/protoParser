#!/bin/bash
OS=$1
ARCH=$2
if [ -n "$OS" -a -n "$ARCH" ]; then
    echo  "specific build, os: $OS; arch: $ARCH"
    GOOS=$OS GOARCH=$ARCH go build -o protoParser main.go
else
    echo "default build..."
    go build -o protoParser main.go
fi

# ./build.sh linux amd64
