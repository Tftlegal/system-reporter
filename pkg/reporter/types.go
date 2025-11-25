package reporter

import (
	"time"
)

// Конфигурация репортера
type Config struct {
	APIBaseURL     string
	ReportEndpoint string
	Timeout        time.Duration
	AgentName      string
}

// Структуры для JSON отчета
type SystemReport struct {
	APIVersion string    `json:"api_version"`
	Generated  time.Time `json:"generated"`
	TotalHosts int       `json:"total_hosts,omitempty"`
	Reports    []Report  `json:"reports"`
}

type Report struct {
	HostID       string            `json:"host_id"`
	ReportNumber int               `json:"report_number"`
	Timestamp    time.Time         `json:"timestamp"`
	Sections     map[string]Section `json:"sections"`
}

type Section struct {
	Title string      `json:"title"`
	Data  interface{} `json:"data"`
}

// Структуры для данных разделов
type HostInfo struct {
	Hostname string    `json:"hostname"`
	OS       string    `json:"os"`
	Kernel   string    `json:"kernel"`
	Uptime   UptimeInfo `json:"uptime"`
}

type UptimeInfo struct {
	Hours    uint64    `json:"hours"`
	BootTime time.Time `json:"boot_time"`
}

type CPUInfo struct {
	Model        string   `json:"model"`
	Cores        int32    `json:"cores"`
	Threads      int      `json:"threads"`
	UsagePercent float64  `json:"usage_percent"`
	LoadAverage  LoadAvg  `json:"load_average"`
}

type LoadAvg struct {
	Load1  float64 `json:"1min"`
	Load5  float64 `json:"5min"`
	Load15 float64 `json:"15min"`
}

type MemoryInfo struct {
	RAM  RAMInfo  `json:"ram"`
	Swap SwapInfo `json:"swap"`
}

type RAMInfo struct {
	TotalGB      float64 `json:"total_gb"`
	AvailableGB  float64 `json:"available_gb"`
	UsedGB       float64 `json:"used_gb"`
	UsedPercent  float64 `json:"used_percent"`
	FreeGB       float64 `json:"free_gb"`
	CachedGB     float64 `json:"cached_gb"`
	BuffersMB    float64 `json:"buffers_mb"`
}

type SwapInfo struct {
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	UsedPercent float64 `json:"used_percent"`
}

type DiskInfo struct {
	Device      string  `json:"device"`
	Mountpoint  string  `json:"mountpoint"`
	Filesystem  string  `json:"filesystem"`
	TotalGB     float64 `json:"total_gb"`
	UsedGB      float64 `json:"used_gb"`
	UsedPercent float64 `json:"used_percent"`
	FreeGB      float64 `json:"free_gb"`
}

type NetworkInfo struct {
	Interfaces []InterfaceInfo `json:"interfaces"`
}

type InterfaceInfo struct {
	Name       string         `json:"name"`
	MAC        string         `json:"mac"`
	IPs        []string       `json:"ips"`
	Statistics InterfaceStats `json:"statistics"`
}

type InterfaceStats struct {
	SentGB     float64 `json:"sent_gb"`
	ReceivedGB float64 `json:"received_gb"`
}

type ProcessInfo struct {
	PID        int32   `json:"pid"`
	Name       string  `json:"name"`
	MemoryMB   float64 `json:"memory_mb"`
	CPUPercent float64 `json:"cpu_percent"`
}

type DockerContainer struct {
	ContainerID string `json:"container_id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Status      string `json:"status"`
	Uptime      string `json:"uptime"`
}

type SecurityStatus struct {
	Fail2ban            string `json:"fail2ban"`
	UfwStatus           string `json:"ufw_status"`
	LastUpdates         string `json:"last_updates"`
	SSHFailedAttempts   int    `json:"ssh_failed_attempts"`
}

// Структура для отправки отчета на API
type APIReportRequest struct {
	Agent  string                 `json:"agent"`
	Report map[string]interface{} `json:"report"`
}
