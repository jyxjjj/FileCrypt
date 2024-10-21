#!/bin/bash
GOOS=windows GOARCH=amd64 go build -o ../dist/windows-x64/upnptest.exe UPnP.go
