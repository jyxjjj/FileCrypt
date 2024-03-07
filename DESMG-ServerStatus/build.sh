#!/bin/bash
GOOS=linux GOARCH=amd64 go build -o ../dist/linux-x64/serverstatus ServerStatus.go
