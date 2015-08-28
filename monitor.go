package vpanel

import (
	"fmt"
	"github.com/cloudfoundry/gosigar"
	"math"
	"os"
	"time"
)

type DiskStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type MemoryStats struct {
	Utilization float64 `json:"utilization"`
	Free        string  `json:"free"`
	Total       string  `json:"total"`
}

type HostStats struct {
	Hostname string               `json:"hostname"`
	Cpu      float64              `json:"cpu"`
	Memory   MemoryStats          `json:"memory"`
	Disk     map[string]DiskStats `json:"disk"`
	Uptime   string               `json:"uptime"`
}

type HostMonitor struct {
	interval  time.Duration
	stopCh    chan bool
	manager   *Manager
	HostStats HostStats
}

func NewHostMonitor(m *Manager, interval time.Duration) *HostMonitor {
	monitor := HostMonitor{
		interval: interval,
		manager:  m,
	}
	return &monitor
}

func round(value float64) float64 {
	return math.Floor(value*100+0.5) / 100
}

func utilize(value, max uint64) float64 {
	return round(100 * float64(value) / float64(max))
}

func humanize(value uint64) string {
	var megabyte = 1024.0
	var gigabyte = megabyte * 1024.0
	var terabyte = gigabyte * 1024.0

	v := float64(value)
	if v >= terabyte {
		return fmt.Sprintf("%.1fT", v/terabyte)
	} else if v >= gigabyte {
		return fmt.Sprintf("%.1fG", v/gigabyte)
	} else if v >= megabyte {
		return fmt.Sprintf("%.1fH", v/megabyte)
	} else {
		return fmt.Sprintf("%.1f", value)
	}
}

func (m *HostMonitor) monitorLoop() {
	cpu := sigar.Cpu{}
	cpu.Get()

	memory := sigar.Mem{}
	fslist := sigar.FileSystemList{}
	uptime := sigar.Uptime{}

	ticker := time.NewTicker(m.interval)

	for {
		select {
		case <-ticker.C:
			oldCpu := cpu
			cpu.Get()
			deltaCpu := cpu.Delta(oldCpu)

			memory.Get()
			uptime.Get()
			fslist.Get()

			hostStats := HostStats{}
			hostStats.Hostname, _ = os.Hostname()
			hostStats.Cpu = utilize(deltaCpu.Total()-deltaCpu.Idle, deltaCpu.Total())
			hostStats.Memory.Free = humanize(memory.Free)
			hostStats.Memory.Total = humanize(memory.Total)
			hostStats.Memory.Utilization = utilize(memory.Total-memory.Free, memory.Total)
			hostStats.Uptime = uptime.Format()

			hostStats.Disk = make(map[string]DiskStats)
			for _, fs := range fslist.List {
				usage := sigar.FileSystemUsage{}
				usage.Get(fs.DirName)
				hostStats.Disk[fs.DirName] = DiskStats{
					Utilization: utilize(usage.Total-usage.Free, usage.Total),
					Free:        fmt.Sprintf("%d", usage.Free),
					Total:       fmt.Sprintf("%d", usage.Total),
				}
			}

			m.HostStats = hostStats
		case <-m.stopCh:
			ticker.Stop()
			m.stopCh <- true
			return
		}
	}
}

func (m *HostMonitor) Start() {
	m.stopCh = make(chan bool)
	go m.monitorLoop()
}

func (m *HostMonitor) Stop() {
	m.stopCh <- true
	<-m.stopCh
	close(m.stopCh)
}
