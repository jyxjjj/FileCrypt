package main

import (
	"bufio"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

func GetCPUNum() string {
	cpuNum, err := cpu.Counts(true)
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%d", cpuNum)
}

func GetCPUFreq() string {
	cpuFreq, err := cpu.Info()
	if err != nil {
		return "-1"
	}
	mhz := cpuFreq[0].Mhz
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.2f", mhz))
}

func GetCPUUsage() string {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return "-1"
	}
	pct := percent[0]
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.2f", pct))
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
			return strings.ReplaceAll(strings.Split(line, "=")[1], "\"", "")
		}
	}
	return "Unknown"
}

func GetMemSize() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "-1"
	}
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.3f", float64(memInfo.Total)/1024/1024/1024))
}

func GetMemUsed() string {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return "-1"
	}
	pct := memInfo.UsedPercent
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.2f", pct))
}

func GetNumProcess() string {
	processes, err := process.Processes()
	if err != nil {
		return "-1"
	}
	return fmt.Sprintf("%d", len(processes))
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
			if name != "" {
				return name + " (" + devicePath + ")"
			} else {
				return "No Label (" + devicePath + ")"
			}
		}
	}
	return "Unknown"
}

func GetDiskUsage() string {
	disks, err := disk.Usage("/")
	if err != nil {
		return "-1"
	}
	pct := disks.UsedPercent
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.2f", pct))
}

func GetDiskSize() string {
	disks, err := disk.Usage("/")
	if err != nil {
		return "-1"
	}
	total := float64(disks.Total) / 1024 / 1024 / 1024
	return removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.3f", total))
}

func GetUptime() string {
	uptime, err := host.Uptime()
	if err != nil {
		return "0D0H0M0S"
	}
	day := uptime / 86400
	hour := uptime % 86400 / 3600
	minute := uptime % 3600 / 60
	second := uptime % 60
	return fmt.Sprintf("%dD%dH%dM%dS", day, hour, minute, second)
}

func GetIORW() []string {
	defaultVal := []string{"-1", "-1"}
	disks, err := disk.Partitions(false)
	if err != nil {
		return defaultVal
	}
	for _, mainDisk := range disks {
		if mainDisk.Mountpoint == "/" {
			devicePath := mainDisk.Device
			deviceBase := filepath.Base(devicePath)
			io0, ioerr0 := disk.IOCounters(devicePath)
			if ioerr0 != nil {
				return defaultVal
			}
			time.Sleep(time.Second)
			io1, ioerr1 := disk.IOCounters(devicePath)
			if ioerr1 != nil {
				return defaultVal
			}
			read := io1[deviceBase].ReadBytes - io0[deviceBase].ReadBytes
			write := io1[deviceBase].WriteBytes - io0[deviceBase].WriteBytes
			read = read / 1024
			readStr := removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.3f", float64(read)))
			write = write / 1024
			writeStr := removeAllRightZeroAndPointForFloatString(fmt.Sprintf("%.3f", float64(write)))
			return []string{readStr, writeStr}

		}
	}
	return defaultVal
}
