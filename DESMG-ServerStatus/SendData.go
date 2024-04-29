package main

import (
	"bytes"
	"net/http"
)

func send(surl string, deviceId string, data []byte) {
	req, err := http.NewRequest("POST", surl, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Set("User-Agent", "DESMG-ServerStatus")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Device-ID", deviceId)
	client := &http.Client{}
	_, err2 := client.Do(req)
	if err2 != nil {
		panic(err)
	}
}
