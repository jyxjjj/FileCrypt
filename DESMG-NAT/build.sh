#!/bin/bash
GOOS=darwin GOARCH=arm64 go build -o ../dist/macos-arm64/nattester NAT.go
GOOS=linux GOARCH=arm64 go build -o ../dist/linux-arm64/nattester NAT.go
GOOS=windows GOARCH=amd64 go build -o ../dist/windows-x64/nattester.exe NAT.go
