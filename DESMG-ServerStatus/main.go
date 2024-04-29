package main

import (
	"encoding/json"
	"os"
)

func main() {
	surl := os.Getenv("SERVER_URL")
	if surl == "" {
		panic("SERVER_URL has not been set")
	}
	key := os.Getenv("SIGN_KEY")
	if key == "" {
		panic("SIGN_KEY has not been set")
	}
	debug := os.Getenv("DEBUG")
	if key == "" {
		debug = "false"
	}
	data := run()
	data.Timestamp = getTimestamp(data.DeviceID)
	sign := signData(key, data)
	signed := SignedData{
		Data: data,
		Sign: sign,
	}
	signedData, err := json.Marshal(signed)
	if err != nil {
		panic(err)
	}
	if debug == "true" {
		println(string(signedData))
	} else {
		send(surl, data.DeviceID, signedData)
	}
}

func run() MonitorData {
	iorw := GetIORW()
	data := MonitorData{
		DeviceID:   GetDeviceID(),
		CPUModel:   GetCPUModel(),
		CPUNum:     GetCPUNum(),
		CPUFreq:    GetCPUFreq(),
		CPUUsage:   GetCPUUsage(),
		OSName:     GetOSName(),
		MemSize:    GetMemSize(),
		MemUsed:    GetMemUsed(),
		NumProcess: GetNumProcess(),
		DiskName:   GetDiskName(),
		DiskUsage:  GetDiskUsage(),
		DiskSize:   GetDiskSize(),
		Uptime:     GetUptime(),
		IORead:     iorw[0],
		IOWrite:    iorw[1],
	}
	return data
}
