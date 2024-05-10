package main

type MonitorData struct {
	DeviceID   string `json:"DeviceID"`
	CPUModel   string `json:"CPUModel"`
	CPUNum     string `json:"CPUNum"`
	CPUFreq    string `json:"CPUFreq"`
	CPUUsage   string `json:"CPUUsage"`
	OSName     string `json:"OSName"`
	MemSize    string `json:"MemSize"`
	MemUsed    string `json:"MemUsed"`
	NumProcess string `json:"NumProcess"`
	DiskName   string `json:"DiskName"`
	DiskUsage  string `json:"DiskUsage"`
	DiskSize   string `json:"DiskSize"`
	Uptime     string `json:"Uptime"`
	IORead     string `json:"IORead"`
	IOWrite    string `json:"IOWrite"`
	Timestamp  string `json:"timestamp"`
}
