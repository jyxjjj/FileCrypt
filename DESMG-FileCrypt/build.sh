#!/bin/bash
GOOS=darwin GOARCH=arm64 go build -o ../dist/macos-arm64/filecrypt FileCrypt.go
GOOS=linux GOARCH=amd64 go build -o ../dist/linux-x64/filecrypt FileCrypt.go
GOOS=windows GOARCH=amd64 go build -o ../dist/windows-x64/filecrypt.exe FileCrypt.go
