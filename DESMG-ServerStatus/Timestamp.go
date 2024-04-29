package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getTimestamp(deviceId string) string {
	tsUrl := "https://www.desmg.com/api/timestamp"
	req, err1 := http.NewRequest("GET", tsUrl, nil)
	if err1 != nil {
		panic(err1)
	}
	req.Header.Set("User-Agent", "DESMG-ServerStatus")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Device-ID", deviceId)
	client := &http.Client{}
	resp, err2 := client.Do(req)
	if err2 != nil {
		panic(err2)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	body, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		panic(err3)
	}
	type TimestampResponse struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
		Data int64  `json:"data"`
	}
	var ts TimestampResponse
	err4 := json.Unmarshal(body, &ts)
	if err4 != nil {
		panic(err4)
	}
	return fmt.Sprintf("%d", ts.Data)
}
