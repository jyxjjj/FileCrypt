package main

type SignedData struct {
	Data MonitorData `json:"data"`
	Sign string      `json:"sign"`
}
