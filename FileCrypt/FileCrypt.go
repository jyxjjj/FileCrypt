package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"os"
)

var (
	method   string
	filename string
	key      string
)

func main() {
	args := os.Args
	if args == nil || len(args) != 4 {
		println("Invalid command")
		println("Usage:")
		println("-encrypt [filename] [key]")
		println("-decrypt [filename] [key]")
		os.Exit(-1)
	}
	method = args[1]
	filename = args[2]
	key = args[3]

	if fileExists(filename) == false {
		println("File not exists.")
		os.Exit(-2)
	}
	if method == "-encrypt" {
		AES256GCMEncrypt()
		os.Exit(0)
	} else if method == "-decrypt" {
		AES256GCMDecrypt()
		os.Exit(0)
	} else {
		println("Invalid method")
		println("Usage:")
		println("-encrypt [filename] [key]")
		println("-decrypt [filename] [key]")
		os.Exit(-1)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func bin2hex(bin []byte) string {
	return hex.EncodeToString(bin)
}

func hex2bin(hexstr string) []byte {
	bin, _ := hex.DecodeString(hexstr)
	return bin
}

func AES256GCMEncrypt() {
	data, err := os.ReadFile(filename)
	if err != nil {
		println(err.Error())
		println("File read failed.")
		os.Exit(-2)
	}
	if len(key) != 32 {
		println("Invaild Key")
		os.Exit(-3)
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		println(err.Error())
		println("Invaild Key")
		os.Exit(-3)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		println(err.Error())
		println("Invaild Length")
		os.Exit(-4)
	}
	nonce := make([]byte, 12) //12byte = 24hex (255="FF", 1 byte = 2 hex)
	_, err = rand.Read(nonce)
	if err != nil {
		println(err.Error())
		println("Random Failed")
		os.Exit(-5)
	}
	encryptedData := gcm.Seal(nil, nonce, data, nil)
	//encryptedLength := len(encryptedData)
	//tagLength := aesgcm.Overhead()
	//println(tagLength)//16byte = 32hex (255="FF", 1 byte = 2 hex)
	err = os.WriteFile(filename+".enc", encryptedData, 0755)
	if err != nil {
		println(err.Error())
		println("Encrypted File write failed.")
		os.Exit(-6)
	}
	err = os.WriteFile(filename+".enc.nonce", nonce, 0755)
	if err != nil {
		println(err.Error())
		println("Nonce File write failed.")
		println("Encrypted File will not be deleted, pls delete them manually.")
		os.Exit(-7)
	}
	println("nonce:" + bin2hex(nonce))
	println("Finished")
}

func AES256GCMDecrypt() {
	if filename[len(filename)-4:] != ".enc" {
		println("In File name not '.enc'.")
		println("Encrypted File read canceled.")
		os.Exit(-2)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		println(err.Error())
		println("Encrypted File read failed.")
		os.Exit(-2)
	}
	nonce, err := os.ReadFile(filename + ".nonce")
	if err != nil {
		println(err.Error())
		println("Nonce File read failed.")
		os.Exit(-2)
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		println(err.Error())
		println("Invaild Key")
		os.Exit(-3)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		println(err.Error())
		println("Invaild Length")
		os.Exit(-4)
	}
	decryptedData, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		println(err.Error())
		println("Decrypt Failed")
		os.Exit(-6)
	}
	if fileExists(filename[:len(filename)-4]) == true {
		println("Decrypted File name duplicated.")
		println("Decrypted File write canceled.")
		os.Exit(-8)
	}
	err = os.WriteFile(filename[:len(filename)-4], decryptedData, 0755)
	if err != nil {
		println(err.Error())
		println("Decrypted File write failed.")
		os.Exit(-6)
	}
	println("Finished")
}
