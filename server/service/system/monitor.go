package system

import (
	"context"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	redisInit "server/setup/redis"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/mem"
)

// processStartTime 在进程启动时固定，用于计算服务运行时长。
var processStartTime = time.Now()

// HostInfo 主机基础信息。
type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       string `json:"os"`       // 如 darwin / linux
	Platform string `json:"platform"` // 如 macOS 15.x / Ubuntu 22.04
	Arch     string `json:"arch"`
	Uptime   uint64 `json:"uptime"` // 主机开机时长（秒）。
}

// CPUInfo CPU 信息。
type CPUInfo struct {
	Cores       int     `json:"cores"`
	ModelName   string  `json:"modelName"`
	UsedPercent float64 `json:"usedPercent"`
}

// MemoryInfo 内存信息，单位字节。
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// DiskInfo 磁盘信息（根分区），单位字节。
type DiskInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

// GoRuntimeInfo Go 进程运行时信息。
type GoRuntimeInfo struct {
	GoVersion  string `json:"goVersion"`
	Goroutines int    `json:"goroutines"`
	HeapAlloc  uint64 `json:"heapAlloc"` // 当前堆内存占用（字节）。
	GCCount    uint32 `json:"gcCount"`
	PID        int    `json:"pid"`
	StartTime  string `json:"startTime"`
	Uptime     int64  `json:"uptime"` // 服务运行时长（秒）。
}

// RedisInfo Redis 服务信息，来自 INFO 命令。
type RedisInfo struct {
	Version          string `json:"version"`
	Mode             string `json:"mode"`
	ConnectedClients string `json:"connectedClients"`
	UsedMemoryHuman  string `json:"usedMemoryHuman"`
	UptimeSeconds    string `json:"uptimeSeconds"`
	TotalCommands    string `json:"totalCommands"`
	KeyCount         int64  `json:"keyCount"`
	Available        bool   `json:"available"`
}

// ServerMonitor 服务监控接口的完整返回。
type ServerMonitor struct {
	Host    HostInfo      `json:"host"`
	CPU     CPUInfo       `json:"cpu"`
	Memory  MemoryInfo    `json:"memory"`
	Disk    DiskInfo      `json:"disk"`
	Runtime GoRuntimeInfo `json:"runtime"`
	Redis   RedisInfo     `json:"redis"`
}

// ServerMonitorInfo 采集主机、进程和 Redis 的监控快照。
// 单项采集失败不影响整体返回（比如容器里拿不到磁盘信息），失败项保持零值。
func ServerMonitorInfo() *ServerMonitor {
	result := &ServerMonitor{}

	if info, err := host.Info(); err == nil {
		result.Host = HostInfo{
			Hostname: info.Hostname,
			OS:       info.OS,
			Platform: info.Platform + " " + info.PlatformVersion,
			Arch:     info.KernelArch,
			Uptime:   info.Uptime,
		}
	}

	result.CPU.Cores = runtime.NumCPU()
	if infos, err := cpu.Info(); err == nil && len(infos) > 0 {
		result.CPU.ModelName = infos[0].ModelName
	}
	// Percent 需要一个采样窗口，200ms 在"看个大概"的监控页里够用且不拖慢接口。
	if percents, err := cpu.Percent(200*time.Millisecond, false); err == nil && len(percents) > 0 {
		result.CPU.UsedPercent = round2(percents[0])
	}

	if vm, err := mem.VirtualMemory(); err == nil {
		result.Memory = MemoryInfo{Total: vm.Total, Used: vm.Used, UsedPercent: round2(vm.UsedPercent)}
	}

	if usage, err := disk.Usage("/"); err == nil {
		result.Disk = DiskInfo{Total: usage.Total, Used: usage.Used, UsedPercent: round2(usage.UsedPercent)}
	}

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	result.Runtime = GoRuntimeInfo{
		GoVersion:  runtime.Version(),
		Goroutines: runtime.NumGoroutine(),
		HeapAlloc:  memStats.HeapAlloc,
		GCCount:    memStats.NumGC,
		PID:        os.Getpid(),
		StartTime:  processStartTime.Format("2006-01-02 15:04:05"),
		Uptime:     int64(time.Since(processStartTime).Seconds()),
	}

	result.Redis = redisMonitorInfo()
	return result
}

// redisMonitorInfo 通过 INFO 命令采集 Redis 状态，Redis 不可用时 Available=false。
func redisMonitorInfo() RedisInfo {
	info := RedisInfo{}
	client := redisInit.Redis.Get()
	if client == nil {
		return info
	}
	ctx := context.Background()

	raw, err := client.Info(ctx).Result()
	if err != nil {
		return info
	}
	fields := parseRedisInfo(raw)
	info.Available = true
	info.Version = fields["redis_version"]
	info.Mode = fields["redis_mode"]
	info.ConnectedClients = fields["connected_clients"]
	info.UsedMemoryHuman = fields["used_memory_human"]
	info.UptimeSeconds = fields["uptime_in_seconds"]
	info.TotalCommands = fields["total_commands_processed"]

	if size, err := client.DBSize(ctx).Result(); err == nil {
		info.KeyCount = size
	}
	return info
}

// parseRedisInfo 把 INFO 命令的 "key:value\r\n" 文本解析成 map。
func parseRedisInfo(raw string) map[string]string {
	fields := make(map[string]string)
	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if key, value, found := strings.Cut(line, ":"); found {
			fields[key] = value
		}
	}
	return fields
}

// round2 保留两位小数，避免前端拿到 63.3333333... 这种展示噪音。
func round2(value float64) float64 {
	rounded, err := strconv.ParseFloat(strconv.FormatFloat(value, 'f', 2, 64), 64)
	if err != nil {
		return value
	}
	return rounded
}
