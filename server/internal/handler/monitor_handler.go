package handler

import (
	"math"
	"runtime"

	"go-admin/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
)

type MonitorHandler struct{}

func NewMonitorHandler() *MonitorHandler { return &MonitorHandler{} }

func round1(f float64) float64 { return math.Round(f*10) / 10 }

// Server 返回主机/CPU/内存/磁盘/运行时等监控指标。单项指标采集失败时降级为 0，不影响整体响应。
func (h *MonitorHandler) Server(c *gin.Context) {
	// host
	hostname, osName, platform := "", "", ""
	var uptime uint64
	if hi, err := host.Info(); err == nil {
		hostname = hi.Hostname
		osName = hi.OS
		platform = hi.Platform
		uptime = hi.Uptime
	}

	// cpu
	var cpuPercent float64
	if p, err := cpu.Percent(0, false); err == nil && len(p) > 0 {
		cpuPercent = round1(p[0])
	}
	cpuCores := 0
	if n, err := cpu.Counts(true); err == nil {
		cpuCores = n
	}

	// memory
	var memUsedMB, memTotalMB uint64
	var memPercent float64
	if vm, err := mem.VirtualMemory(); err == nil {
		memUsedMB = vm.Used / 1024 / 1024
		memTotalMB = vm.Total / 1024 / 1024
		memPercent = round1(vm.UsedPercent)
	}

	// disk
	var diskUsedGB, diskTotalGB float64
	var diskPercent float64
	if du, err := disk.Usage("/"); err == nil {
		diskUsedGB = round1(float64(du.Used) / 1024 / 1024 / 1024)
		diskTotalGB = round1(float64(du.Total) / 1024 / 1024 / 1024)
		diskPercent = round1(du.UsedPercent)
	}

	// runtime
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	response.OK(c, gin.H{
		"host": gin.H{
			"hostname": hostname,
			"os":       osName,
			"platform": platform,
			"arch":     runtime.GOARCH,
		},
		"cpu": gin.H{
			"percent": cpuPercent,
			"cores":   cpuCores,
		},
		"memory": gin.H{
			"usedMB":  memUsedMB,
			"totalMB": memTotalMB,
			"percent": memPercent,
		},
		"disk": gin.H{
			"usedGB":  diskUsedGB,
			"totalGB": diskTotalGB,
			"percent": diskPercent,
		},
		"runtime": gin.H{
			"goVersion":  runtime.Version(),
			"goroutines": runtime.NumGoroutine(),
			"numCPU":     runtime.NumCPU(),
			"allocMB":    ms.Alloc / 1024 / 1024,
		},
		"uptimeSeconds": uptime,
	})
}
