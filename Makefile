main: sn fc

clean:
	rm -rf dist

sn: DESMG-SerialNumber/SerialNumber.go
	GOOS=darwin GOARCH=arm64 go build -o dist/macos-arm64/serialnumber DESMG-SerialNumber/SerialNumber.go
	GOOS=linux GOARCH=amd64 go build -o dist/linux-x64/serialnumber DESMG-SerialNumber/SerialNumber.go
	GOOS=windows GOARCH=amd64 go build -o dist/windows-x64/serialnumber.exe DESMG-SerialNumber/SerialNumber.go

fc: FileCrypt/FileCrypt.go
	GOOS=darwin GOARCH=arm64 go build -o dist/macos-arm64/filecrypt FileCrypt/FileCrypt.go
	GOOS=linux GOARCH=amd64 go build -o dist/linux-x64/filecrypt FileCrypt/FileCrypt.go
	GOOS=windows GOARCH=amd64 go build -o dist/windows-x64/filecrypt.exe FileCrypt/FileCrypt.go

release:
	tar -czvf dist.tar.gz dist
	gh release create "$(shell date +'%Y%m%d%H%M%S')" -t "$(shell date +'%Y-%m-%d %H:%M:%S')" -n "" ./dist.tar.gz
	rm -rf dist.tar.gz
