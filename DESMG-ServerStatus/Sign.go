package main

import (
	"crypto/sha512"
	"fmt"
	"net/url"
	"strings"
)

func signData(key string, data MonitorData) string {
	dataToBeSign := httpBuildQuery(data) + key
	if debug == "true" {
		println(dataToBeSign)
	}
	sign := strings.ToUpper(fmt.Sprintf("%x", sha512.Sum512([]byte(dataToBeSign))))
	return sign
}

func httpBuildQuery(data MonitorData) string {
	v := url.Values{}
	v.Set("DeviceID", data.DeviceID)
	v.Set("CPUModel", data.CPUModel)
	v.Set("CPUNum", data.CPUNum)
	v.Set("CPUFreq", data.CPUFreq)
	v.Set("CPUUsage", data.CPUUsage)
	v.Set("OSName", data.OSName)
	v.Set("MemSize", data.MemSize)
	v.Set("MemUsed", data.MemUsed)
	v.Set("NumProcess", data.NumProcess)
	v.Set("DiskName", data.DiskName)
	v.Set("DiskUsage", data.DiskUsage)
	v.Set("DiskSize", data.DiskSize)
	v.Set("Uptime", data.Uptime)
	v.Set("IORead", data.IORead)
	v.Set("IOWrite", data.IOWrite)
	v.Set("Timestamp", data.Timestamp)
	encoded := v.Encode()
	rfc3986 := strings.ReplaceAll(encoded, "+", "%20")
	return rfc3986
}
