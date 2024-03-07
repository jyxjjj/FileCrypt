#!/bin/bash
GOOS=darwin GOARCH=arm64 go build -o ../dist/macos-arm64/serialnumber SerialNumber.go
GOOS=linux GOARCH=amd64 go build -o ../dist/linux-x64/serialnumber SerialNumber.go
GOOS=windows GOARCH=amd64 go build -o ../dist/windows-x64/serialnumber.exe SerialNumber.go
