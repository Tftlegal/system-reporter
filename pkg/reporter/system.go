package reporter

import (
	"fmt"
	"runtime"
	"sort"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// Вспомогательные функции
func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

func bytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

// GetHostID возвращает host_id (hostname) системы
func GetHostID() string {
	hostInfo, err := host.Info()
	if err != nil {
		return "unknown-host"
	}
	return hostInfo.Hostname
}

// GenerateSystemReport генерирует полный системный отчет
func GenerateSystemReport() (*SystemReport, error) {
	hostID := GetHostID()

	report := &SystemReport{
		APIVersion: "1.0",
		Generated:  time.Now(),
		TotalHosts: 1,
		Reports: []Report{
			{
				HostID:       hostID,
				ReportNumber: 1,
				Timestamp:    time.Now(),
				Sections:     make(map[string]Section),
			},
		},
	}

	// Собираем данные по секциям
	sections := map[string]func() (interface{}, error){
		"1": getHostInformation,
		"2": getCPUInformation,
		"3": getMemoryInformation,
		"4": getDiskInformation,
		"5": getNetworkInformation,
		"6": getTopProcessesByMemory,
		"7": getDockerContainers,
		"8": getSecurityStatus,
	}

	titles := map[string]string{
		"1": "HOST INFORMATION",
		"2": "CPU INFORMATION", 
		"3": "MEMORY INFORMATION",
		"4": "DISK INFORMATION",
		"5": "NETWORK INFORMATION",
		"6": "TOP PROCESSES BY MEMORY",
		"7": "DOCKER CONTAINERS",
		"8": "SECURITY STATUS",
	}

	for key, fn := range sections {
		data, err := fn()
		if err != nil {
			fmt.Printf("Warning: failed to get %s: %v\n", titles[key], err)
			continue
		}
		report.Reports[0].Sections[key] = Section{
			Title: titles[key],
			Data:  data,
		}
	}

	return report, nil
}

func getHostInformation() (*HostInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	return &HostInfo{
		Hostname: hostInfo.Hostname,
		OS:       fmt.Sprintf("%s %s %s", hostInfo.OS, hostInfo.Platform, hostInfo.PlatformVersion),
		Kernel:   hostInfo.KernelVersion,
		Uptime: UptimeInfo{
			Hours:    hostInfo.Uptime / 3600,
			BootTime: time.Unix(int64(hostInfo.BootTime), 0),
		},
	}, nil
}

func getCPUInformation() (*CPUInfo, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	loadAvg, err := load.Avg()
	if err != nil {
		return nil, err
	}

	model := ""
	cores := int32(0)
	if len(cpuInfo) > 0 {
		model = cpuInfo[0].ModelName
		cores = cpuInfo[0].Cores
	}

	return &CPUInfo{
		Model:        model,
		Cores:        cores,
		Threads:      runtime.NumCPU(),
		UsagePercent: percent[0],
		LoadAverage: LoadAvg{
			Load1:  loadAvg.Load1,
			Load5:  loadAvg.Load5,
			Load15: loadAvg.Load15,
		},
	}, nil
}

func getMemoryInformation() (*MemoryInfo, error) {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryInfo{
		RAM: RAMInfo{
			TotalGB:      bytesToGB(vmem.Total),
			AvailableGB:  bytesToGB(vmem.Available),
			UsedGB:       bytesToGB(vmem.Used),
			UsedPercent:  vmem.UsedPercent,
			FreeGB:       bytesToGB(vmem.Free),
			CachedGB:     bytesToGB(vmem.Cached),
			BuffersMB:    bytesToMB(vmem.Buffers),
		},
		Swap: SwapInfo{
			TotalGB:     bytesToGB(swap.Total),
			UsedGB:      bytesToGB(swap.Used),
			UsedPercent: swap.UsedPercent,
		},
	}, nil
}

func getDiskInformation() ([]DiskInfo, error) {
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var disks []DiskInfo
	for _, partition := range partitions {
		usage, err := disk.Usage(partition.Mountpoint)
		if err != nil {
			continue
		}

		disks = append(disks, DiskInfo{
			Device:      partition.Device,
			Mountpoint:  partition.Mountpoint,
			Filesystem:  partition.Fstype,
			TotalGB:     bytesToGB(usage.Total),
			UsedGB:      bytesToGB(usage.Used),
			UsedPercent: usage.UsedPercent,
			FreeGB:      bytesToGB(usage.Free),
		})
	}

	return disks, nil
}

func getNetworkInformation() (*NetworkInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	ioMap := make(map[string]net.IOCountersStat)
	for _, io := range ioCounters {
		ioMap[io.Name] = io
	}

	var ifaceList []InterfaceInfo
	for _, iface := range interfaces {
		if len(iface.Addrs) == 0 {
			continue
		}

		var ips []string
		for _, addr := range iface.Addrs {
			ips = append(ips, addr.Addr)
		}

		stats := InterfaceStats{}
		if io, exists := ioMap[iface.Name]; exists {
			stats.SentGB = bytesToGB(io.BytesSent)
			stats.ReceivedGB = bytesToGB(io.BytesRecv)
		}

		ifaceList = append(ifaceList, InterfaceInfo{
			Name:       iface.Name,
			MAC:        iface.HardwareAddr,
			IPs:        ips,
			Statistics: stats,
		})
	}

	return &NetworkInfo{Interfaces: ifaceList}, nil
}

func getTopProcessesByMemory() ([]ProcessInfo, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var procList []ProcessInfo
	for _, p := range processes {
		if len(procList) >= 20 {
			break
		}

		name, err := p.Name()
		if err != nil {
			continue
		}

		memInfo, err := p.MemoryInfo()
		if err != nil || memInfo == nil {
			continue
		}

		cpuPercent, err := p.CPUPercent()
		if err != nil {
			cpuPercent = 0
		}

		procList = append(procList, ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			MemoryMB:   bytesToMB(memInfo.RSS),
			CPUPercent: cpuPercent,
		})
	}

	sort.Slice(procList, func(i, j int) bool {
		return procList[i].MemoryMB > procList[j].MemoryMB
	})

	if len(procList) > 10 {
		procList = procList[:10]
	}

	return procList, nil
}

func getDockerContainers() ([]DockerContainer, error) {
	return []DockerContainer{}, nil
}

func getSecurityStatus() (*SecurityStatus, error) {
	return &SecurityStatus{
		Fail2ban:          "unknown",
		UfwStatus:         "unknown", 
		LastUpdates:       time.Now().Format("2006-01-02"),
		SSHFailedAttempts: 0,
	}, nil
}
