package main

import (
	"encoding/json"
	"flag"
	"os"
)

var surl, key, debug string

func init() {
	flag.StringVar(&surl, "s", "", "Server URL")
	flag.StringVar(&key, "k", "", "Sign Key")
	flag.StringVar(&debug, "d", "false", "Debug mode")
	flag.Parse()

	if surl == "" {
		surl = os.Getenv("SERVER_URL")
		if surl == "" {
			panic("SERVER_URL has not been set")
		}
	}
	if key == "" {
		key = os.Getenv("SIGN_KEY")
		if key == "" {
			panic("SIGN_KEY has not been set")
		}
	}
	if debug == "false" || debug == "" {
		debug = os.Getenv("DEBUG")
		if debug == "false" || debug == "" {
			debug = "false"
		}
	}
}

func main() {
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
