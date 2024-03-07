package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	jsonData := run()
	send(jsonData)
}

func send(jsonData string) {
	url := os.Getenv("SERVER_URL")
	if url == "" {
		println("SERVER_URL has not been set")
		return
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err2 := client.Do(req)
	if err2 != nil {
		return
	}
}

type JSONData struct {
	DeviceID   string `json:"DeviceID"`
	CPUModel   string `json:"CPUModel"`
	CPUNum     int64  `json:"CPUNum"`
	CPUFreq    string `json:"CPUFreq"`
	CPUUsage   int64  `json:"CPUUsage"`
	OSName     string `json:"OSName"`
	MemSize    string `json:"MemSize"`
	MemUsed    int64  `json:"MemUsed"`
	NumProcess int64  `json:"NumProcess"`
	DiskName   string `json:"DiskName"`
	DiskUsage  int64  `json:"DiskUsage"`
	DiskSize   string `json:"DiskSize"`
	Uptime     string `json:"Uptime"`
	IORead     int64  `json:"IORead"`
	IOWrite    int64  `json:"IOWrite"`
}

func run() string {
	iorw := GetIORW()
	data := JSONData{
		DeviceID:   GetDeviceID(),
		CPUModel:   GetCPUModel(),
		CPUNum:     GetCPUNum(),
		CPUFreq:    GetCPUFreq(),
		CPUUsage:   GetCPUUsage(),
		OSName:     GetOSName(),
		MemSize:    GetMemSize(),
		MemUsed:    GetMemUsed(),
		NumProcess: GetNumProcess(),

		DiskName:  GetDiskName(),
		DiskUsage: GetDiskUsage(),
		DiskSize:  GetDiskSize(),
		Uptime:    GetUptime(),
		IORead:    iorw[0],
		IOWrite:   iorw[1],
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return string(jsonData)
}

func GetDeviceID() string {
	id, err := host.HostID()
	if err != nil {
		return ""
	}
	return id
}

func GetCPUModel() string {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return "Unknown"
	}
	return cpuInfo[0].ModelName
}

func GetCPUNum() int64 {
	cpuNum, err := cpu.Counts(true)
	if err != nil {
		return -1
	}
	return int64(cpuNum)
}

func GetCPUFreq() string {
	cpuFreq, err := cpu.Info()
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%.2f", cpuFreq[0].Mhz/1000)
}

func GetCPUUsage() int64 {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return -1
	}
	return int64(percent[0])
}

func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetOSName() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		hostInfo, err := host.Info()
		if err != nil {
			return "Unknown"
		}
		return FirstUpper(hostInfo.OS) + " " + hostInfo.PlatformVersion + " " + FirstUpper(hostInfo.Platform) + " " + hostInfo.KernelArch + " " + hostInfo.KernelVersion + " " + hostInfo.Hostname
	}
	defer func(file *os.File) {
		ferr := file.Close()
		if ferr != nil {
			panic(ferr)
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "PRETTY_NAME") {
			return strings.Split(line, "=")[1]
		}
	}
	return "Unknown"
}

func GetMemSize() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%d", memInfo.Total/1024/1024/1024)
}

func GetMemUsed() int64 {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return -1
	}
	return int64(memInfo.UsedPercent)
}

func GetNumProcess() int64 {
	processes, err := process.Processes()
	if err != nil {
		return -1
	}
	return int64(len(processes))
}

func GetDiskName() string {
	disks, err := disk.Partitions(false)
	if err != nil {
		return "Unknown"
	}
	for _, mainDisk := range disks {
		if mainDisk.Mountpoint == "/" {
			devicePath := mainDisk.Device
			name, nameerr := disk.Label("/")
			if nameerr != nil {
				return "Unknown Disk (" + devicePath + ")"
			}
			return name + " (" + devicePath + ")"
		}
	}
	return "Unknown"
}

func GetDiskUsage() int64 {
	disks, err := disk.Usage("/")
	if err != nil {
		return -1
	}
	return int64(disks.UsedPercent)
}

func GetDiskSize() string {
	disks, err := disk.Usage("/")
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%.2f", float64(disks.Total/1024/1024/1024))
}

func GetUptime() string {
	uptime, err := host.Uptime()
	if err != nil {
		return "-1"
	}
	day := uptime / 86400
	hour := uptime % 86400 / 3600
	minute := uptime % 3600 / 60
	second := uptime % 60
	return fmt.Sprintf("%dD%dH%dM%dS", day, hour, minute, second)
}

func GetIORW() []int64 {
	disks, err := disk.Partitions(false)
	if err != nil {
		return []int64{-1, -1}
	}
	for _, mainDisk := range disks {
		if mainDisk.Mountpoint == "/" {
			devicePath := mainDisk.Device
			deviceBase := filepath.Base(devicePath)
			io0, ioerr0 := disk.IOCounters(devicePath)
			if ioerr0 != nil {
				return []int64{-1, -1}
			}
			time.Sleep(time.Second)
			io1, ioerr1 := disk.IOCounters(devicePath)
			if ioerr1 != nil {
				return []int64{-1, -1}
			}
			read := io1[deviceBase].ReadBytes - io0[deviceBase].ReadBytes
			write := io1[deviceBase].WriteBytes - io0[deviceBase].WriteBytes
			return []int64{int64(read), int64(write)}

		}
	}
	return []int64{-1, -1}
}
